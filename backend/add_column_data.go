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
	Col      Column     `json:"column"`
	Sheet_id uuid.UUID  `json:"sheet_id"`
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
		Idx:     params.Data.Idx,
		Value:   params.Data.Value,
		Name:    params.Col.Name,
		SheetID: params.Sheet_id,
	}
	_, err = cfg.db.AddColumnData(r.Context(), addColumnDataParams)
	if err != nil {
		msg := fmt.Sprintf("Column data could not be updated: %s", err)
		respondWithError(w, 500, msg)
	}

	updateSheetRowCountParams := database.UpdateSheetRowCountParams{
		RowCount:   params.Data.Idx + 1,
		RowCount_2: params.Data.Idx + 1,
		ID:         params.Sheet_id,
	}
	err = cfg.db.UpdateSheetRowCount(r.Context(), updateSheetRowCountParams)
	if err != nil {
		msg := fmt.Sprintf("Could not update the row count: %s", err)
		respondWithError(w, 500, msg)
	}
	respondWithJSON(w, 200, "")
}
