package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type newColDataParams struct {
	Data     ColumnData `json:"data"`
	ColumnId uuid.UUID  `json:"column_id"`
}

func (cfg *apiConfig) AddColumnData(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := newColDataParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	addColumnDataParams := database.AddColumnDataParams{
		Idx:      params.Data.Idx,
		Value:    params.Data.Value,
		ColumnID: params.ColumnId,
	}
	_, err = cfg.db.AddColumnData(r.Context(), addColumnDataParams)
	if err != nil {
		msg := fmt.Sprintf("Column data could not be updated: %s", err)
		respondWithError(w, 500, msg)
	}
	respondWithJSON(w, 200, "")
}
