package main

import (
	"context"

	"github.com/google/uuid"
)

func (cfg *apiConfig) checkBranchPermission(userId, branchId uuid.UUID, permType string, ctx context.Context) bool {
	branch, err := cfg.db.GetBranch(ctx, branchId)
	if err != nil {
		return false
	}
	return cfg.checkTablePermission(userId, branch.TableID, permType, ctx)
}
