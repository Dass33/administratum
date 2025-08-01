package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) UpdateColumn(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	col := Column{}

	err := decoder.Decode(&col)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, 400, msg)
		return
	}
	updateColumnParams := database.UpdateColumnParams{
		Name:     col.Name,
		Type:     col.Type,
		Required: col.Required,
		ID:       col.ID,
	}
	err = cfg.db.UpdateColumn(r.Context(), updateColumnParams)
	if err != nil {
		msg := fmt.Sprintf("column could not be updated: %s", err)
		respondWithError(w, 500, msg)
	}
	respondWithJSON(w, 200, "")
}
