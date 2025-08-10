package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	err = cfg.db.DeleteBranch(r.Context(), branchId)
	if err != nil {
		msg := fmt.Sprintf("branch branch not be deleted: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
