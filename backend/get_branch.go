package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Enum struct {
	Name    string   `json:"name"`
	SheetID string   `json:"sheet_id"`
	Vals    []string `json:"vals"`
}

type Branch struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsProtected bool      `json:"is_protected"`
	Enums       []Enum    `json:"enums"`
}

func (cfg *apiConfig) getBranchHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	branchIdStr := chi.URLParam(r, "branch_id")
	branchId, err := uuid.Parse(branchIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the branch from url: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !cfg.checkBranchPermission(userId, branchId, "read", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient read permissions")
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

	enums, err := cfg.getEnumsForBranch(branch_id, ctx)
	if err != nil {
		return Branch{}, errors.New("Could not get enums for branch")
	}

	data := Branch{
		ID:          branch_id,
		Name:        branch.Name,
		IsProtected: branch.IsProtected,
		Enums:       enums,
	}
	return data, nil
}

func (cfg *apiConfig) getEnumsForBranch(branchID uuid.UUID, ctx context.Context) ([]Enum, error) {
	sheets, err := cfg.db.GetSheetsFromBranch(ctx, branchID)
	if err != nil {
		return nil, err
	}

	var enums []Enum

	for _, sheet := range sheets {
		if sheet.Type == "enums" {
			columns, err := cfg.db.GetColumnsFromSheet(ctx, sheet.ID)
			if err != nil {
				continue
			}

			if len(columns) == 0 {
				continue
			}

			firstColumn := columns[0]
			columnData, err := cfg.db.GetColumnsData(ctx, firstColumn.ID)
			if err != nil {
				continue
			}

			var vals []string
			for _, data := range columnData {
				if data.Value.Valid && data.Value.String != "" {
					vals = append(vals, data.Value.String)
				}
			}

			enum := Enum{
				Name:    sheet.Name,
				SheetID: sheet.ID.String(),
				Vals:    vals,
			}
			enums = append(enums, enum)
		}
	}

	return enums, nil
}
