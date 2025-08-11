package main

import (
	"context"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) checkBranchPermission(userId, branchId uuid.UUID, permType string, ctx context.Context) bool {
	branch, err := cfg.db.GetBranch(ctx, branchId)
	if err != nil {
		return false
	}

	userTable, err := cfg.db.GetUserTables(ctx, database.GetUserTablesParams{
		UserID:  userId,
		TableID: branch.TableID,
	})
	if err != nil {
		return false
	}

	if userTable.Permission == OwnerPermission {
		return true
	}

	if userTable.Permission == ContributorPermission && (permType == "read" || permType == "write") {
		return true
	}

	return false
}
