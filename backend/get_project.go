package main

import (
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectData struct {
	Table TableData `table:"table"`
	Sheet Sheet     `sheet:"sheet"`
}

func (cfg *apiConfig) getProjectHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	tableIdStr := chi.URLParam(r, "table_id")
	tableId, err := uuid.Parse(tableIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the table id from url: %s", err)
		respondWithError(w, 400, msg)
	}
	optionalTableId := uuid.NullUUID{
		UUID:  tableId,
		Valid: true,
	}
	table, err := cfg.GetTable(userId, optionalTableId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get table: %s", err)
		respondWithError(w, 500, msg)
	}

	sheets, err := cfg.db.GetSheetsFromTable(r.Context(), tableId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from table: %s", err)
		respondWithError(w, 500, msg)
	}

	sheetId := uuid.UUID{}
	if len(sheets) == 0 {
		if len(table.BranchesNames) == 0 {
			createBranchParams := database.CreateBranchParams{
				Name:    "main",
				TableID: tableId,
			}
			branch, err := cfg.db.CreateBranch(r.Context(), createBranchParams)
			if err != nil {
				msg := fmt.Sprintf("Could not create a main branch: %s", err)
				respondWithError(w, 500, msg)
			}
			branchIdName := IdName{ID: branch.ID, Name: branch.Name}
			table.BranchesNames = append(table.BranchesNames, branchIdName)
		}
		createSheetParams := database.CreateSheetParams{
			Name:     "config",
			RowCount: 0,
			BranchID: table.BranchesNames[0].ID,
		}
		dbSheet, err := cfg.db.CreateSheet(r.Context(), createSheetParams)
		if err != nil {
			msg := fmt.Sprintf("Could not create a config sheet: %s", err)
			respondWithError(w, 500, msg)
		}
		sheetId = dbSheet.ID
	} else {
		sheetId = sheets[0].ID
	}

	optionalSheetId := uuid.NullUUID{
		UUID:  sheetId,
		Valid: true,
	}
	sheet, err := cfg.GetSheet(optionalSheetId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get sheet: %s", err)
		respondWithError(w, 500, msg)
	}

	setOpenedSheetParams := database.SetOpenedSheetParams{
		ID:          userId,
		OpenedSheet: optionalSheetId,
	}
	err = cfg.db.SetOpenedSheet(r.Context(), setOpenedSheetParams)
	if err != nil {
		msg := fmt.Sprintf("Could not set opened sheet: %s", err)
		respondWithError(w, 500, msg)
	}
	data := ProjectData{
		Table: table,
		Sheet: sheet,
	}
	respondWithJSON(w, 200, data)
}
