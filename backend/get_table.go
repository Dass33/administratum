package main

import (
	"context"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type Table struct {
}

func (cfg *apiConfig) GetTable(user_id, table_id uuid.UUID, ctx context.Context) (Table, error) {
	//todo
	userTablesParams := database.GetUserTablesParams{
		UserID:  user_id,
		TableID: table_id,
	}
	userTables, err := cfg.db.GetUserTables(ctx, userTablesParams)
	if err != nil {

	}

	return Table{}, nil
}
