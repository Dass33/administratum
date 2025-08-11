package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type createBranchParam struct {
	Name         string `json:"name"`
	IsProtected  bool   `json:"is_protected"`
	TableId      string `json:"table_id"`
	CurrBranchId string `json:"curr_branch_id"`
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
	currBranchId, err := uuid.Parse(params.CurrBranchId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the curr branch id: %s", err)
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

	err = cfg.copyBranchSheets(r.Context(), currBranchId, branch.ID)
	if err != nil {
		msg := fmt.Sprintf("Could not copy sheets to new branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	cfg.switchBranch(w, r, branch.ID, userId, http.StatusCreated)
}

func (cfg *apiConfig) copyBranchSheets(ctx context.Context, sourceBranchId, targetBranchId uuid.UUID) error {
	dbSheets, err := cfg.db.GetSheetsFromBranch(ctx, sourceBranchId)
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
		sheet, err := cfg.db.CreateSheet(ctx, createSheetParams)
		if err != nil {
			return fmt.Errorf("could not create sheet: %w", err)
		}

		err = cfg.copySheetColumns(ctx, dbSheets[i].ID, sheet.ID)
		if err != nil {
			return fmt.Errorf("could not copy columns for sheet %s: %w", sheet.Name, err)
		}
	}
	return nil
}

func (cfg *apiConfig) copySheetColumns(ctx context.Context, sourceSheetId, targetSheetId uuid.UUID) error {
	columns, err := cfg.GetColumns(sourceSheetId, ctx)
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
		_, err := cfg.db.AddColumn(ctx, addColumnParams)
		if err != nil {
			return fmt.Errorf("could not add column: %w", err)
		}

		err = cfg.copyColumnData(ctx, columns[e], targetSheetId)
		if err != nil {
			return fmt.Errorf("could not copy data for column %s: %w", columns[e].Name, err)
		}
	}
	return nil
}

func (cfg *apiConfig) copyColumnData(ctx context.Context, column Column, targetSheetId uuid.UUID) error {
	for j := range column.Data {
		cell := &column.Data[j]
		addColumnDataParams := database.AddColumnDataParams{
			Idx:     cell.Idx,
			Value:   cell.Value,
			Type:    cell.Type,
			Name:    column.Name,
			SheetID: targetSheetId,
		}
		_, err := cfg.db.AddColumnData(ctx, addColumnDataParams)
		if err != nil {
			return fmt.Errorf("could not add column data: %w", err)
		}
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
