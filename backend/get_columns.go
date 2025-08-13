package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type ColumnData struct {
	ID    uuid.UUID      `json:"id"`
	Idx   int64          `json:"idx"`
	Value sql.NullString `json:"value"`
	Type  sql.NullString `json:"type"`
}

type Column struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Required bool         `json:"required"`
	Data     []ColumnData `json:"data"`
}

func (cfg *apiConfig) GetColumnsWithTx(txQueries *database.Queries, sheet_id uuid.UUID, ctx context.Context) ([]Column, error) {
	rows, err := txQueries.GetColumnsWithDataBySheet(ctx, sheet_id)
	if err != nil {
		return nil, errors.New("Could not get columns with data for given sheet id")
	}

	columnMap := make(map[uuid.UUID]*Column)
	var columnOrder []uuid.UUID

	for _, row := range rows {
		columnID := row.ColumnID

		if _, exists := columnMap[columnID]; !exists {
			columnMap[columnID] = &Column{
				ID:       columnID,
				Name:     row.ColumnName,
				Type:     row.ColumnType,
				Required: row.ColumnRequired,
				Data:     make([]ColumnData, 0),
			}
			columnOrder = append(columnOrder, columnID)
		}

		if row.DataID.Valid {
			columnData := ColumnData{
				ID:    row.DataID.UUID,
				Idx:   row.DataIdx.Int64,
				Value: row.DataValue,
				Type:  row.DataType,
			}
			columnMap[columnID].Data = append(columnMap[columnID].Data, columnData)
		}
	}

	data := make([]Column, 0, len(columnOrder))
	for _, columnID := range columnOrder {
		data = append(data, *columnMap[columnID])
	}

	return data, nil
}

func (cfg *apiConfig) GetColumns(sheet_id uuid.UUID, ctx context.Context) ([]Column, error) {
	return cfg.GetColumnsWithTx(cfg.db, sheet_id, ctx)
}
