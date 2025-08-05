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
