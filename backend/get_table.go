package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type Enum struct {
	Name    string    `json:"name"`
	SheetId uuid.UUID `json:"sheet_id"`
	Vals    []string  `json:"vals"`
}

type TableData struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	GameUrl         sql.NullString `json:"game_url"`
	Permision       string         `json:"permision"`
	BranchesIdNames []IdName       `json:"branches_id_names"`
}

func (cfg *apiConfig) GetTable(user_id uuid.UUID, optional_table_id uuid.NullUUID, ctx context.Context) (TableData, error) {
	if !optional_table_id.Valid {
		return TableData{}, errors.New("Table id not present")
	}
	table_id := optional_table_id.UUID

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
		ID:              table_id,
		Name:            table.Name,
		GameUrl:         table.GameUrl,
		Permision:       userTables.Permission,
		BranchesIdNames: branchNames,
	}

	return data, nil
}
