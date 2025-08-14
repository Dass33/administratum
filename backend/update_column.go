package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateColumnHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	col := Column{}

	err := decoder.Decode(&col)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	rowsAffected, err := cfg.db.UpdateColumnWithPermissionCheck(r.Context(), database.UpdateColumnWithPermissionCheckParams{
		Name:     col.Name,
		Type:     col.Type,
		Required: col.Required,
		ID:       col.ID,
		UserID:   id,
	})
	if err != nil {
		msg := fmt.Sprintf("column could not be updated: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	
	if rowsAffected == 0 {
		respondWithError(w, http.StatusForbidden, "Column not found or insufficient permissions")
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
