package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		respondWithError(w, 400, msg)
		return
	}

	projectId, err := uuid.Parse(params.ProjectId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the project id: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	err = cfg.db.DeleteTable(r.Context(), projectId)
	if err != nil {
		msg := fmt.Sprintf("Project could not be deleted: %s", err)
		respondWithError(w, 500, msg)
		return
	}
	respondWithJSON(w, 204, "")
}
