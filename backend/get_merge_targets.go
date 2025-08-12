package main

import (
	"net/http"

	"github.com/google/uuid"
)

type MergeTarget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetMergeTargetsResponse struct {
	ValidTargets []MergeTarget `json:"valid_targets"`
	TargetBranch MergeTarget   `json:"target_branch"`
}

func (cfg *apiConfig) getMergeTargetsHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	tableIDStr := r.URL.Query().Get("table_id")
	if tableIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing table_id parameter")
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table_id format")
		return
	}

	ctx := r.Context()

	allBranches, err := cfg.db.GetBranchesFromTable(ctx, tableID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get table branches")
		return
	}

	oldestBranch, err := cfg.db.GetOldestBranchFromTable(ctx, tableID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find oldest branch in table")
		return
	}

	if !cfg.checkBranchPermission(userId, oldestBranch.ID, "write", ctx) {
		respondWithError(w, http.StatusForbidden, "No write permission on target branch")
		return
	}

	validSources := make([]MergeTarget, 0)
	for _, branch := range allBranches {
		if branch.ID == oldestBranch.ID {
			continue
		}

		if !cfg.checkBranchPermission(userId, branch.ID, "read", ctx) {
			continue
		}

		validSources = append(validSources, MergeTarget{
			ID:   branch.ID.String(),
			Name: branch.Name,
		})
	}

	response := GetMergeTargetsResponse{
		ValidTargets: validSources,
		TargetBranch: MergeTarget{
			ID:   oldestBranch.ID.String(),
			Name: oldestBranch.Name,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}
