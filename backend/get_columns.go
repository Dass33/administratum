package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Column struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Required bool      `json:"required"`
	Data     []any     `json:"data"`
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
		vals := make([]any, 0)
		for e, item := range columns_data {
			if int64(e) != item.Idx {
				continue
			}
			vals = append(vals, item.Value)
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
