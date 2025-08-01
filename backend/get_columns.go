package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type ColumnData struct {
	ID    uuid.UUID      `json:"id"`
	Idx   int64          `json:"idx"`
	Value sql.NullString `json:"value"`
}

type Column struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Required bool         `json:"required"`
	Data     []ColumnData `json:"data"`
}

func (cfg *apiConfig) GetColumns(sheet_id uuid.UUID, ctx context.Context) ([]Column, error) {

	columns, err := cfg.db.GetColumnsFromSheet(ctx, sheet_id)
	if err != nil {
		return nil, errors.New("Could not get columns with given sheet id")
	}

	data := make([]Column, 0, len(columns))

	for i := range columns {
		columns_data, err := cfg.db.GetColumnsData(ctx, columns[i].ID)
		if err != nil {
			return nil, errors.New("Could not get columns data with given column id")
		}
		vals := make([]ColumnData, 0)
		for e, _ := range columns_data {
			if int64(e) != columns_data[e].Idx {
				continue
			}
			item := ColumnData{
				ID:    columns_data[e].ID,
				Idx:   columns_data[e].Idx,
				Value: columns_data[e].Value,
			}
			vals = append(vals, item)
		}
		col := Column{
			ID:       columns[i].ID,
			Name:     columns[i].Name,
			Type:     columns[i].Type,
			Required: columns[i].Required,
			Data:     vals,
		}
		data = append(data, col)
	}

	return data, nil
}
