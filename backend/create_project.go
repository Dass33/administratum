package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type nameParam struct {
	Name string `json:"name"`
}

func (cfg *apiConfig) createProjectHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := nameParam{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	createTableParams := database.CreateTableParams{
		Name:    params.Name,
		GameUrl: sql.NullString{},
	}
	table, err := cfg.db.CreateTable(r.Context(), createTableParams)
	if err != nil {
		msg := fmt.Sprintf("Could not create table: %s", err)
		respondWithError(w, 500, msg)
		return
	}
	createUserTableParams := database.CreateUserTableParams{
		UserID:     userId,
		TableID:    table.ID,
		Permission: OwnerPermission,
	}
	_, err = cfg.db.CreateUserTable(r.Context(), createUserTableParams)
	if err != nil {
		msg := fmt.Sprintf("Could not create user table: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	cfg.switchProject(w, r, table.ID, userId, 201)
}

func (cfg *apiConfig) switchProject(w http.ResponseWriter, r *http.Request, tableId, userId uuid.UUID, code int) {
	optionalTableId := uuid.NullUUID{
		UUID:  tableId,
		Valid: true,
	}
	table, err := cfg.GetTable(userId, optionalTableId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get table: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	sheets, err := cfg.db.GetSheetsFromTable(r.Context(), tableId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from table: %s", err)
		respondWithError(w, 500, msg)
		return
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
				return
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
			return
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
		return
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
	data := ProjectData{
		Table: table,
		Sheet: sheet,
	}
	respondWithJSON(w, code, data)
}
