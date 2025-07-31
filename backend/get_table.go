package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Data     any    `json:"data"`
}

type Sheet struct {
	Name             string   `json:"name"`
	Columns          []Column `json:"columns"`
	OpenedBranchName string   `json:"opened_branch_name"`
}

type TableData struct {
	GameUrl       sql.NullString `json:"game_url"`
	Permision     string         `json:"permision"`
	OpenedSheet   Sheet          `json:"opened_sheet"`
	SheetsNames   []string       `json:"sheets_names"`
	BranchesNames []string       `json:"branches_names"`
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

	// sheetsNames, err := cfg.db.GetSheets(ctx, )

	data := TableData{
		GameUrl:   table.GameUrl,
		Permision: userTables.Permission,
	}

	return data, nil
}
