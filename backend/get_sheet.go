package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Sheet struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	RowCount      int64     `json:"row_count"`
	Type          string    `json:"type"`
	Columns       []Column  `json:"columns"`
	CurrBranch    Branch    `json:"curr_branch"`
	SheetsIdNames []IdName  `json:"sheets_id_names"`
}

func (cfg *apiConfig) getSheetHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	sheetIdStr := chi.URLParam(r, "sheet_id")
	sheetId, err := uuid.Parse(sheetIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id from url: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	// Get sheet and check permissions in one go
	dbSheet, err := cfg.db.GetSheet(r.Context(), sheetId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Sheet not found")
		return
	}

	if !cfg.checkBranchPermission(userId, dbSheet.BranchID, "read", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient read permissions")
		return
	}

	sheet, err := cfg.buildSheetResponse(dbSheet, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get sheet: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	optionalSheetId := uuid.NullUUID{
		UUID:  sheetId,
		Valid: true,
	}
	setOpenedSheetParams := database.SetOpenedSheetParams{
		ID:          userId,
		OpenedSheet: optionalSheetId,
	}
	err = cfg.db.SetOpenedSheet(r.Context(), setOpenedSheetParams)
	if err != nil {
		msg := fmt.Sprintf("Could not set opened sheet: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	respondWithJSON(w, 200, sheet)
}

func (cfg *apiConfig) GetSheet(optional_sheet_id uuid.NullUUID, ctx context.Context) (Sheet, error) {
	if !optional_sheet_id.Valid {
		return Sheet{}, errors.New("Table id not present")
	}
	sheet_id := optional_sheet_id.UUID

	sheet, err := cfg.db.GetSheet(ctx, sheet_id)
	if err != nil {
		return Sheet{}, errors.New("Could not get sheet with given id")
	}

	branch, err := cfg.db.GetBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get branch with given id")
	}
	enums, err := cfg.getEnumsForBranch(sheet.BranchID, ctx)
	if err != nil {
		return Sheet{}, errors.New("Could not get enums for branch")
	}

	currBranch := Branch{
		ID:          branch.ID,
		Name:        branch.Name,
		IsProtected: branch.IsProtected,
		Enums:       enums,
	}

	sheets, err := cfg.db.GetSheetsFromBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get sheets with given branch id")
	}
	sheetsIdNames := make([]IdName, 0, len(sheets))

	for i := range sheets {
		item := IdName{
			ID:   sheets[i].ID,
			Name: sheets[i].Name,
		}
		sheetsIdNames = append(sheetsIdNames, item)
	}

	columns, err := cfg.GetColumns(sheet_id, ctx)
	if err != nil {
		return Sheet{}, errors.New("Could not get columns with given sheet id")
	}

	rowCount := 0
	for i := range columns {
		currLen := len(columns[i].Data)
		if currLen > rowCount {
			rowCount = currLen
		}
	}

	data := Sheet{
		ID:            sheet_id,
		Name:          sheet.Name,
		RowCount:      int64(rowCount),
		Type:          sheet.Type,
		CurrBranch:    currBranch,
		SheetsIdNames: sheetsIdNames,
		Columns:       columns,
	}

	return data, nil
}

func (cfg *apiConfig) buildSheetResponse(sheet database.Sheet, ctx context.Context) (Sheet, error) {
	branch, err := cfg.db.GetBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get branch with given id")
	}
	enums, err := cfg.getEnumsForBranch(sheet.BranchID, ctx)
	if err != nil {
		return Sheet{}, errors.New("Could not get enums for branch")
	}

	currBranch := Branch{
		ID:          branch.ID,
		Name:        branch.Name,
		IsProtected: branch.IsProtected,
		Enums:       enums,
	}

	sheets, err := cfg.db.GetSheetsFromBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get sheets with given branch id")
	}
	sheetsIdNames := make([]IdName, 0, len(sheets))

	for i := range sheets {
		item := IdName{
			ID:   sheets[i].ID,
			Name: sheets[i].Name,
		}
		sheetsIdNames = append(sheetsIdNames, item)
	}

	columns, err := cfg.GetColumns(sheet.ID, ctx)
	if err != nil {
		return Sheet{}, errors.New("Could not get columns with given sheet id")
	}

	rowCount := 0
	for i := range columns {
		currLen := len(columns[i].Data)
		if currLen > rowCount {
			rowCount = currLen
		}
	}

	data := Sheet{
		ID:            sheet.ID,
		Name:          sheet.Name,
		RowCount:      int64(rowCount),
		Type:          sheet.Type,
		CurrBranch:    currBranch,
		SheetsIdNames: sheetsIdNames,
		Columns:       columns,
	}

	return data, nil
}
