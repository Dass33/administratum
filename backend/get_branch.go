package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Branch struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsProtected bool      `json:"is_protected"`
}

func (cfg *apiConfig) getBranchHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	branchIdStr := chi.URLParam(r, "branch_id")
	branchId, err := uuid.Parse(branchIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the branch from url: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	cfg.switchBranch(w, r, branchId, userId, http.StatusOK)
}

func (cfg *apiConfig) GetBranch(user_id uuid.UUID, optional_branch_id uuid.NullUUID, ctx context.Context) (Branch, error) {
	if !optional_branch_id.Valid {
		return Branch{}, errors.New("Table id not present")
	}
	branch_id := optional_branch_id.UUID

	branch, err := cfg.db.GetBranch(ctx, branch_id)
	if err != nil {
		return Branch{}, errors.New("Could not get table with given id")
	}

	data := Branch{
		ID:          branch_id,
		Name:        branch.Name,
		IsProtected: branch.IsProtected,
	}
	return data, nil
}
