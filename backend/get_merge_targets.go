package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type MergeTarget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetMergeTargetsResponse struct {
	ValidTargets   []MergeTarget `json:"valid_targets"`
	TargetBranch   MergeTarget   `json:"target_branch"`
}

func (cfg *apiConfig) getMergeTargetsHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	// Get table ID from query parameters (since we want to show all possible source branches)
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

	ctx := context.Background()

	// Get all branches from the table
	allBranches, err := cfg.db.GetBranchesFromTable(ctx, tableID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get table branches")
		return
	}

	// Find the oldest branch (target for all merges)
	oldestBranch, err := cfg.db.GetOldestBranchFromTable(ctx, tableID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find oldest branch in table")
		return
	}

	// Check if user has write permission on the target (oldest) branch
	if !cfg.checkBranchPermission(userId, oldestBranch.ID, "write", ctx) {
		respondWithError(w, http.StatusForbidden, "No write permission on target branch")
		return
	}

	// Return all branches except the oldest (target) as possible sources
	validSources := make([]MergeTarget, 0)
	for _, branch := range allBranches {
		// Skip the oldest branch (it's the target, can't merge into itself)
		if branch.ID == oldestBranch.ID {
			continue
		}

		// Check if user has read permission on this potential source
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