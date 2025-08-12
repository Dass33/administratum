package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type createBranchParam struct {
	Name        string `json:"name"`
	IsProtected bool   `json:"is_protected"`
	TableId     string `json:"table_id"`
}

func (cfg *apiConfig) createBranchHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := createBranchParam{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	tableId, err := uuid.Parse(params.TableId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the table id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	createBranchParams := database.CreateBranchParams{
		Name:        params.Name,
		IsProtected: params.IsProtected,
		TableID:     tableId,
	}

	branch, err := cfg.db.CreateBranch(r.Context(), createBranchParams)
	if err != nil {
		msg := fmt.Sprintf("Could not create branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	oldestBranch, err := cfg.db.GetOldestBranchFromTable(r.Context(), tableId)
	if err != nil {
		if err != sql.ErrNoRows {
			msg := fmt.Sprintf("Could not get oldest branch: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		if oldestBranch.ID != branch.ID {
			err = cfg.copyBranchSheetsWithTransaction(r.Context(), oldestBranch.ID, branch.ID)
			if err != nil {
				msg := fmt.Sprintf("Could not copy sheets to new branch: %s", err)
				respondWithError(w, http.StatusInternalServerError, msg)
				return
			}
		}
	}

	cfg.switchBranch(w, r, branch.ID, userId, http.StatusCreated)
}

func (cfg *apiConfig) copyBranchSheetsWithTransaction(ctx context.Context, sourceBranchId, targetBranchId uuid.UUID) error {
	tx, err := cfg.rawDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txQueries := cfg.db.WithTx(tx)

	err = cfg.copyBranchSheetsInTx(ctx, tx, txQueries, sourceBranchId, targetBranchId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (cfg *apiConfig) copyBranchSheetsInTx(ctx context.Context, tx *sql.Tx, txQueries *database.Queries, sourceBranchId, targetBranchId uuid.UUID) error {
	dbSheets, err := txQueries.GetSheetsFromBranch(ctx, sourceBranchId)
	if err != nil {
		return fmt.Errorf("could not get source branch sheets: %w", err)
	}

	for i := range dbSheets {
		createSheetParams := database.CreateSheetParams{
			BranchID:      targetBranchId,
			Name:          dbSheets[i].Name,
			Type:          dbSheets[i].Type,
			SourceSheetID: sql.NullString{String: dbSheets[i].ID.String(), Valid: true},
		}
		sheet, err := txQueries.CreateSheet(ctx, createSheetParams)
		if err != nil {
			return fmt.Errorf("could not create sheet: %w", err)
		}

		err = cfg.copySheetColumnsInTx(ctx, tx, txQueries, dbSheets[i].ID, sheet.ID)
		if err != nil {
			return fmt.Errorf("could not copy columns for sheet %s: %w", sheet.Name, err)
		}
	}
	return nil
}

func (cfg *apiConfig) copySheetColumnsInTx(ctx context.Context, tx *sql.Tx, txQueries *database.Queries, sourceSheetId, targetSheetId uuid.UUID) error {
	columns, err := cfg.GetColumnsWithTx(txQueries, sourceSheetId, ctx)
	if err != nil {
		return fmt.Errorf("could not get columns: %w", err)
	}

	for e := range columns {
		addColumnParams := database.AddColumnParams{
			Name:           columns[e].Name,
			Type:           columns[e].Type,
			Required:       columns[e].Required,
			SheetID:        targetSheetId,
			SourceColumnID: sql.NullString{String: columns[e].ID.String(), Valid: true},
		}
		newColumn, err := txQueries.AddColumn(ctx, addColumnParams)
		if err != nil {
			return fmt.Errorf("could not add column: %w", err)
		}

		if len(columns[e].Data) > 0 {
			err = cfg.copyColumnDataBulk(ctx, tx, columns[e].Data, newColumn.ID)
			if err != nil {
				return fmt.Errorf("could not copy data for column %s: %w", columns[e].Name, err)
			}
		}
	}
	return nil
}

// copyColumnDataBulk performs a bulk insert of column data using raw SQL for performance.
// This uses a single multi-row INSERT statement instead of individual SQLC calls
// to significantly improve performance when copying large datasets.
func (cfg *apiConfig) copyColumnDataBulk(ctx context.Context, tx *sql.Tx, data []ColumnData, columnId uuid.UUID) error {
	if len(data) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*4)

	for _, cell := range data {
		valueStrings = append(valueStrings, "(gen_random_uuid(), ?, ?, ?, ?, datetime('now'), datetime('now'))")
		valueArgs = append(valueArgs, cell.Idx, cell.Value, cell.Type, columnId)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO column_data (id, idx, value, type, column_id, created_at, updated_at)
		VALUES %s
	`, strings.Join(valueStrings, ", "))

	_, err := tx.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("could not bulk insert %d column data rows: %w", len(data), err)
	}

	return nil
}

type BranchData struct {
	Branch Branch
	Sheet  Sheet
}

func (cfg *apiConfig) switchBranch(w http.ResponseWriter, r *http.Request, branchId, userId uuid.UUID, code int) {
	optionalBranchId := uuid.NullUUID{
		UUID:  branchId,
		Valid: true,
	}
	branch, err := cfg.GetBranch(userId, optionalBranchId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	dbSheets, err := cfg.db.GetSheetsFromBranch(r.Context(), branchId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	sheetId := uuid.UUID{}
	if len(dbSheets) == 0 {
		sheetId, err = cfg.createMapSheet(r.Context(), "config", branchId)
		if err != nil {
			msg := fmt.Sprintf("Could not create map sheet: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		sheetId = dbSheets[0].ID
	}

	optionalSheetId := uuid.NullUUID{
		UUID:  sheetId,
		Valid: true,
	}
	sheet, err := cfg.GetSheet(optionalSheetId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	setOpenedSheetParams := database.SetOpenedSheetParams{
		ID:          userId,
		OpenedSheet: optionalSheetId,
	}
	err = cfg.db.SetOpenedSheet(r.Context(), setOpenedSheetParams)
	if err != nil {
		msg := fmt.Sprintf("Could not set opened sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	data := BranchData{
		Branch: branch,
		Sheet:  sheet,
	}
	respondWithJSON(w, code, data)
}
