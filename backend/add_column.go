package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type ColParams struct {
	SheetId string `json:"sheet_id"`
	Col     Column `json:"col"`
}

type ColumnResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

func (cfg *apiConfig) addColumnHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := ColParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	sheet_id, err := uuid.Parse(params.SheetId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !cfg.checkSheetPermission(id, sheet_id, "write", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient write permissions")
		return
	}

	addColumnParams := database.AddColumnParams{
		Name:     params.Col.Name,
		Type:     params.Col.Type,
		Required: params.Col.Required,
		SheetID:  sheet_id,
	}
	newCol, err := cfg.db.AddColumn(r.Context(), addColumnParams)
	if err != nil {
		msg := fmt.Sprintf("Column could not be updated: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	response := ColumnResponse{
		ID:       newCol.ID.String(),
		Name:     newCol.Name,
		Type:     newCol.Type,
		Required: newCol.Required,
	}
	respondWithJSON(w, http.StatusCreated, response)
}
