package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateColumnDataHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	colData := ColumnData{}

	err := decoder.Decode(&colData)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, 400, msg)
		return
	}
	updateColumnDataParams := database.UpdateColumnDataParams{
		ID:    colData.ID,
		Value: colData.Value,
	}
	err = cfg.db.UpdateColumnData(r.Context(), updateColumnDataParams)
	if err != nil {
		msg := fmt.Sprintf("column data could not be updated: %s", err)
		respondWithError(w, 500, msg)
		return
	}
	respondWithJSON(w, 200, "")
}
