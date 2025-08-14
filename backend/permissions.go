package main

import (
	"context"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

const OwnerPermission string = "owner"
const ContributorPermission string = "contributor"

func PermissionWeight(s string) (int, bool) {
	switch s {
	case OwnerPermission:
		return 0, true
	case ContributorPermission:
		return 1, true
	}
	return -1, false
}

func (cfg *apiConfig) canAssignPermision(userId, tableId uuid.UUID, perm string, ctx context.Context) bool {
	recieverWeight, ok := PermissionWeight(perm)
	if !ok {
		return false
	}

	getUserTableParams := database.GetUserTablesParams{
		UserID:  userId,
		TableID: tableId,
	}
	userTable, err := cfg.db.GetUserTables(ctx, getUserTableParams)
	if err != nil {
		return false
	}

	granterWeight, ok := PermissionWeight(userTable.Permission)
	if !ok || granterWeight > recieverWeight {
		return false
	}

	return true
}

func (cfg *apiConfig) checkTablePermission(userId, tableId uuid.UUID, permType string, ctx context.Context) bool {
	userTable, err := cfg.db.GetUserTables(ctx, database.GetUserTablesParams{
		UserID:  userId,
		TableID: tableId,
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

func (cfg *apiConfig) checkSheetPermission(userId, sheetId uuid.UUID, permType string, ctx context.Context) bool {
	sheet, err := cfg.db.GetSheet(ctx, sheetId)
	if err != nil {
		return false
	}
	return cfg.checkBranchPermission(userId, sheet.BranchID, permType, ctx)
}

func (cfg *apiConfig) checkBranchPermission(userId, branchId uuid.UUID, permType string, ctx context.Context) bool {
	branch, err := cfg.db.GetBranch(ctx, branchId)
	if err != nil {
		return false
	}
	return cfg.checkTablePermission(userId, branch.TableID, permType, ctx)
}
