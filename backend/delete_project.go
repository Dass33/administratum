package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type deleteProjectParams struct {
	ProjectId string `json:"ProjectId"`
}

func (cfg *apiConfig) deleteProjectHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := deleteProjectParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	projectId, err := uuid.Parse(params.ProjectId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the project id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	rowsAffected, err := cfg.db.DeleteTableWithPermissionCheck(r.Context(), database.DeleteTableWithPermissionCheckParams{
		ID:     projectId,
		UserID: userId,
	})
	if err != nil {
		msg := fmt.Sprintf("Project could not be deleted: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	
	if rowsAffected == 0 {
		respondWithError(w, http.StatusForbidden, "Project not found or insufficient permissions")
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
