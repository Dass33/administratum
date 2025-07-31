package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type TableData struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	GameUrl       sql.NullString `json:"game_url"`
	Permision     string         `json:"permision"`
	BranchesNames []IdName       `json:"branches_names"`
}

func (cfg *apiConfig) GetTable(user_id uuid.UUID, table_id_ptr *uuid.UUID, ctx context.Context) (TableData, error) {
	if table_id_ptr == nil {
		return TableData{}, errors.New("Table id not present")
	}
	table_id := uuid.UUID(table_id_ptr.NodeID())

	table, err := cfg.db.GetTable(ctx, table_id)
	if err != nil {
		return TableData{}, errors.New("Could not get table with given id")
	}

	userTablesParams := database.GetUserTablesParams{
		UserID:  user_id,
		TableID: table_id,
	}
	userTables, err := cfg.db.GetUserTables(ctx, userTablesParams)
	if err != nil {
		return TableData{}, errors.New("Could not get user tables with given id")
	}

	branches, err := cfg.db.GetBranchesFromTable(ctx, table_id)
	if err != nil {
		return TableData{}, errors.New("Could not get branches from table")
	}
	branchNames := make([]IdName, 0, len(branches))

	for i := range branches {
		item := IdName{
			ID:   branches[i].ID,
			Name: branches[i].Name,
		}
		branchNames = append(branchNames, item)
	}

	data := TableData{
		ID:            table_id,
		Name:          table.Name,
		GameUrl:       table.GameUrl,
		Permision:     userTables.Permission,
		BranchesNames: branchNames,
	}

	return data, nil
}
