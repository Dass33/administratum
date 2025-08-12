package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type deleteBranchParams struct {
	BranchId string `json:"branch_id"`
}

func (cfg *apiConfig) deleteBranchHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := deleteBranchParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	branchId, err := uuid.Parse(params.BranchId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the branch id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	rowsAffected, err := cfg.db.DeleteBranchWithPermissionCheck(r.Context(), database.DeleteBranchWithPermissionCheckParams{
		ID:     branchId,
		UserID: userId,
	})
	if err != nil {
		msg := fmt.Sprintf("branch could not be deleted: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	
	if rowsAffected == 0 {
		respondWithError(w, http.StatusForbidden, "Branch not found or insufficient permissions")
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
