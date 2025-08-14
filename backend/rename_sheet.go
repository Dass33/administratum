package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type renameSheetParams struct {
	Name    string `json:"Name"`
	SheetId string `json:"SheetId"`
}

func (cfg *apiConfig) renameSheetHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := renameSheetParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	sheetId, err := uuid.Parse(params.SheetId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !cfg.checkSheetPermission(userId, sheetId, "write", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient write permissions")
		return
	}

	renameSheetParams := database.RenameSheetParams{
		Name: params.Name,
		ID:   sheetId,
	}
	err = cfg.db.RenameSheet(r.Context(), renameSheetParams)
	if err != nil {
		msg := fmt.Sprintf("Sheet could not be renamed: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
