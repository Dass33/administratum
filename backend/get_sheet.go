package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Sheet struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Columns       []Column  `json:"columns"`
	BranchIdName  IdName    `json:"branch_id_name"`
	SheetsIdNames []IdName  `json:"sheets_id_names"`
}

func (cfg *apiConfig) GetSheet(optional_sheet_id uuid.NullUUID, ctx context.Context) (Sheet, error) {
	if !optional_sheet_id.Valid {
		return Sheet{}, errors.New("Table id not present")
	}
	sheet_id := optional_sheet_id.UUID

	sheet, err := cfg.db.GetSheet(ctx, sheet_id)
	if err != nil {
		return Sheet{}, errors.New("Could not get sheet with given id")
	}

	branch, err := cfg.db.GetBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get branch with given id")
	}
	branchIdName := IdName{
		ID:   branch.ID,
		Name: branch.Name,
	}

	sheets, err := cfg.db.GetSheetsFromBranch(ctx, sheet.BranchID)
	if err != nil {
		return Sheet{}, errors.New("Could not get sheets with given branch id")
	}
	sheetsIdNames := make([]IdName, 0, len(sheets))

	for i := range sheets {
		item := IdName{
			ID:   sheets[i].ID,
			Name: sheets[i].Name,
		}
		sheetsIdNames = append(sheetsIdNames, item)
	}

	columns, err := cfg.GetColumns(sheet_id, ctx)
	if err != nil {
		return Sheet{}, errors.New("Could not get columns with given sheet id")
	}

	data := Sheet{
		ID:            sheet_id,
		Name:          sheet.Name,
		BranchIdName:  branchIdName,
		SheetsIdNames: sheetsIdNames,
		Columns:       columns,
	}

	return data, nil
}
