package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type renameProjectParams struct {
	Name      string `json:"Name"`
	ProjectId string `json:"ProjectId"`
}

func (cfg *apiConfig) renemeProjectHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := renameProjectParams{}

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

	renameTableParams := database.RenameTableParams{
		Name: params.Name,
		ID:   projectId,
	}
	err = cfg.db.RenameTable(r.Context(), renameTableParams)
	if err != nil {
		msg := fmt.Sprintf("Project could not be renamed: %s", err)
		respondWithError(w, 500, msg)
		return
	}
	respondWithJSON(w, 200, "")
}
