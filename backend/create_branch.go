package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type createBranchParam struct {
	Name        string `json:"name"`
	IsProtected bool   `json:"is_protected"`
	TableId     string `json:"table_id"`
}

func (cfg *apiConfig) createBranchHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	startTime := time.Now()
	log.Printf("🚀 Starting branch creation: %s", time.Now().Format("15:04:05.000"))
	
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
	
	branchCreateStart := time.Now()
	branch, err := cfg.db.CreateBranch(r.Context(), createBranchParams)
	log.Printf("📝 Branch DB creation took: %v", time.Since(branchCreateStart))
	if err != nil {
		msg := fmt.Sprintf("Could not create branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	// Find the oldest branch in the table to copy from
	findBranchStart := time.Now()
	oldestBranch, err := cfg.db.GetOldestBranchFromTable(r.Context(), tableId)
	log.Printf("🔍 Finding oldest branch took: %v", time.Since(findBranchStart))
	
	if err != nil {
		// If no branches exist yet (sql.ErrNoRows), that's fine - create empty branch
		if err == sql.ErrNoRows {
			log.Printf("📋 No existing branches - creating empty branch")
		} else {
			msg := fmt.Sprintf("Could not get oldest branch: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		// Oldest branch exists, copy from it (unless it's the branch we just created)
		if oldestBranch.ID != branch.ID {
			copyStart := time.Now()
			log.Printf("📁 Starting to copy from branch: %s", oldestBranch.ID)
			err = cfg.copyBranchSheetsWithTransaction(r.Context(), oldestBranch.ID, branch.ID)
			log.Printf("📁 Branch copying took: %v", time.Since(copyStart))
			if err != nil {
				msg := fmt.Sprintf("Could not copy sheets to new branch: %s", err)
				respondWithError(w, http.StatusInternalServerError, msg)
				return
			}
		}
	}

	switchStart := time.Now()
	cfg.switchBranch(w, r, branch.ID, userId, http.StatusCreated)
	log.Printf("🔄 Branch switching took: %v", time.Since(switchStart))
	log.Printf("✅ Total branch creation time: %v", time.Since(startTime))
}

func (cfg *apiConfig) copyBranchSheetsWithTransaction(ctx context.Context, sourceBranchId, targetBranchId uuid.UUID) error {
	// Start a transaction to batch all operations
	txStart := time.Now()
	tx, err := cfg.rawDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	log.Printf("💳 Transaction start took: %v", time.Since(txStart))

	// Use transaction-enabled queries
	txQueries := cfg.db.WithTx(tx)
	
	err = cfg.copyBranchSheetsInTx(ctx, txQueries, sourceBranchId, targetBranchId)
	if err != nil {
		return err
	}

	// Commit the transaction
	commitStart := time.Now()
	err = tx.Commit()
	log.Printf("✅ Transaction commit took: %v", time.Since(commitStart))
	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (cfg *apiConfig) copyBranchSheetsInTx(ctx context.Context, txQueries *database.Queries, sourceBranchId, targetBranchId uuid.UUID) error {
	getSheetsStart := time.Now()
	dbSheets, err := txQueries.GetSheetsFromBranch(ctx, sourceBranchId)
	log.Printf("📊 Getting sheets from branch took: %v (found %d sheets)", time.Since(getSheetsStart), len(dbSheets))
	if err != nil {
		return fmt.Errorf("could not get source branch sheets: %w", err)
	}

	for i := range dbSheets {
		sheetStart := time.Now()
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
		log.Printf("📄 Created sheet '%s' in: %v", sheet.Name, time.Since(sheetStart))

		columnsCopyStart := time.Now()
		err = cfg.copySheetColumnsInTx(ctx, txQueries, dbSheets[i].ID, sheet.ID)
		log.Printf("🗂️  Copying columns for sheet '%s' took: %v", sheet.Name, time.Since(columnsCopyStart))
		if err != nil {
			return fmt.Errorf("could not copy columns for sheet %s: %w", sheet.Name, err)
		}
	}
	return nil
}

func (cfg *apiConfig) copySheetColumnsInTx(ctx context.Context, txQueries *database.Queries, sourceSheetId, targetSheetId uuid.UUID) error {
	getColumnsStart := time.Now()
	columns, err := cfg.GetColumnsWithTx(txQueries, sourceSheetId, ctx)
	log.Printf("📋 Getting columns took: %v (found %d columns)", time.Since(getColumnsStart), len(columns))
	if err != nil {
		return fmt.Errorf("could not get columns: %w", err)
	}

	for e := range columns {
		columnStart := time.Now()
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
		log.Printf("📝 Created column '%s' in: %v", columns[e].Name, time.Since(columnStart))

		if len(columns[e].Data) > 0 {
			dataStart := time.Now()
			err = cfg.copyColumnDataBatchInTx(ctx, txQueries, columns[e].Data, newColumn.ID)
			log.Printf("💾 Copying %d data rows for column '%s' took: %v", len(columns[e].Data), columns[e].Name, time.Since(dataStart))
			if err != nil {
				return fmt.Errorf("could not copy data for column %s: %w", columns[e].Name, err)
			}
		} else {
			log.Printf("📋 Column '%s' has no data to copy", columns[e].Name)
		}
	}
	return nil
}

func (cfg *apiConfig) copySheetColumns(ctx context.Context, sourceSheetId, targetSheetId uuid.UUID) error {
	getColumnsStart := time.Now()
	columns, err := cfg.GetColumns(sourceSheetId, ctx)
	log.Printf("📋 Getting columns took: %v (found %d columns)", time.Since(getColumnsStart), len(columns))
	if err != nil {
		return fmt.Errorf("could not get columns: %w", err)
	}

	for e := range columns {
		columnStart := time.Now()
		addColumnParams := database.AddColumnParams{
			Name:           columns[e].Name,
			Type:           columns[e].Type,
			Required:       columns[e].Required,
			SheetID:        targetSheetId,
			SourceColumnID: sql.NullString{String: columns[e].ID.String(), Valid: true},
		}
		newColumn, err := cfg.db.AddColumn(ctx, addColumnParams)
		if err != nil {
			return fmt.Errorf("could not add column: %w", err)
		}
		log.Printf("📝 Created column '%s' in: %v", columns[e].Name, time.Since(columnStart))

		if len(columns[e].Data) > 0 {
			dataStart := time.Now()
			err = cfg.copyColumnDataBatch(ctx, columns[e].Data, newColumn.ID)
			log.Printf("💾 Copying %d data rows for column '%s' took: %v", len(columns[e].Data), columns[e].Name, time.Since(dataStart))
			if err != nil {
				return fmt.Errorf("could not copy data for column %s: %w", columns[e].Name, err)
			}
		} else {
			log.Printf("📋 Column '%s' has no data to copy", columns[e].Name)
		}
	}
	return nil
}

func (cfg *apiConfig) copyColumnDataBatchInTx(ctx context.Context, txQueries *database.Queries, data []ColumnData, columnId uuid.UUID) error {
	// Batch insert data using transaction queries for better performance
	for j := range data {
		cell := &data[j]
		bulkParams := database.BulkAddColumnDataParams{
			Idx:      cell.Idx,
			Value:    cell.Value,
			Type:     cell.Type,
			ColumnID: columnId,
		}
		err := txQueries.BulkAddColumnData(ctx, bulkParams)
		if err != nil {
			return fmt.Errorf("could not add column data: %w", err)
		}
	}
	return nil
}

func (cfg *apiConfig) copyColumnDataBatch(ctx context.Context, data []ColumnData, columnId uuid.UUID) error {
	// Batch insert data using multiple single inserts in a transaction for better performance
	// This reduces overhead compared to individual calls while staying compatible with SQLite
	for j := range data {
		cell := &data[j]
		bulkParams := database.BulkAddColumnDataParams{
			Idx:      cell.Idx,
			Value:    cell.Value,
			Type:     cell.Type,
			ColumnID: columnId,
		}
		err := cfg.db.BulkAddColumnData(ctx, bulkParams)
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
	getBranchStart := time.Now()
	optionalBranchId := uuid.NullUUID{
		UUID:  branchId,
		Valid: true,
	}
	branch, err := cfg.GetBranch(userId, optionalBranchId, r.Context())
	log.Printf("🌿 Getting branch data took: %v", time.Since(getBranchStart))
	if err != nil {
		msg := fmt.Sprintf("Could not get branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	getSheetsStart := time.Now()
	dbSheets, err := cfg.db.GetSheetsFromBranch(r.Context(), branchId)
	log.Printf("📊 Getting sheets for switch took: %v", time.Since(getSheetsStart))
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from branch: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	sheetId := uuid.UUID{}
	if len(dbSheets) == 0 {
		createSheetStart := time.Now()
		sheetId, err = cfg.createMapSheet(r.Context(), "config", branchId)
		log.Printf("📄 Creating default map sheet took: %v", time.Since(createSheetStart))
		if err != nil {
			msg := fmt.Sprintf("Could not create map sheet: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		sheetId = dbSheets[0].ID
		log.Printf("📋 Using existing sheet: %s", sheetId)
	}

	getSheetStart := time.Now()
	optionalSheetId := uuid.NullUUID{
		UUID:  sheetId,
		Valid: true,
	}
	sheet, err := cfg.GetSheet(optionalSheetId, r.Context())
	log.Printf("📋 Getting sheet data took: %v", time.Since(getSheetStart))
	if err != nil {
		msg := fmt.Sprintf("Could not get sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	updateUserStart := time.Now()
	setOpenedSheetParams := database.SetOpenedSheetParams{
		ID:          userId,
		OpenedSheet: optionalSheetId,
	}
	err = cfg.db.SetOpenedSheet(r.Context(), setOpenedSheetParams)
	log.Printf("👤 Updating user opened sheet took: %v", time.Since(updateUserStart))
	if err != nil {
		msg := fmt.Sprintf("Could not set opened sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	
	responseStart := time.Now()
	data := BranchData{
		Branch: branch,
		Sheet:  sheet,
	}
	respondWithJSON(w, code, data)
	log.Printf("📤 Preparing response took: %v", time.Since(responseStart))
}
