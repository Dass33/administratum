package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type adjustBranch struct {
	Name        string `json:"name"`
	BranchId    string `json:"branch_id"`
	IsProtected bool   `json:"is_protected"`
}

func (cfg *apiConfig) updateBranchHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := adjustBranch{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding branch: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	branchId, err := uuid.Parse(params.BranchId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the branch id from url: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	updateBranchParams := database.UpdateBranchParams{
		Name:        params.Name,
		IsProtected: params.IsProtected,
		ID:          branchId,
	}
	err = cfg.db.UpdateBranch(r.Context(), updateBranchParams)
	if err != nil {
		msg := fmt.Sprintf("Branch could not be updated: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
