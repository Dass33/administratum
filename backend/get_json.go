package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getJsonHandler(w http.ResponseWriter, r *http.Request) {
	branchIdStr := chi.URLParam(r, "branch_id")
	branchId, err := uuid.Parse(branchIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id from url: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	sheetsDb, err := cfg.db.GetSheetsFromBranch(r.Context(), branchId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from branch id: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	data := make([]any, 0, len(sheetsDb))
	for _, sheet := range sheetsDb {
		rows, err := cfg.stringMapsFromSheet(sheet, r.Context())
		if err != nil {
			msg := fmt.Sprintf("Could not get rows from sheet: %s", err)
			respondWithError(w, 500, msg)
			return
		}

		data = append(data, rows)
	}
	respondWithJSON(w, 200, data)
}

func (cfg *apiConfig) stringMapsFromSheet(sheet database.Sheet, ctx context.Context) ([]map[string]string, error) {
	columns, err := cfg.GetColumns(sheet.ID, ctx)
	if err != nil {
		return nil, errors.New("Could not get columns with given sheet id")
	}

	rows := make([]map[string]string, 0, sheet.RowCount)

	for i := range sheet.RowCount {
		row := make(map[string]string)

		for e := range columns {
			col := &columns[e]
			if val, ok := getDataAtColIdx(col.Data, i); ok {
				row[col.Name] = val.String
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// the column data has to be sorted ascending by their index
func getDataAtColIdx(data []ColumnData, idx int64) (sql.NullString, bool) {
	i := sort.Search(len(data), func(i int) bool {
		return data[i].Idx >= idx
	})

	if i >= len(data) || data[i].Idx != idx {
		return sql.NullString{}, false
	}

	return data[i].Value, true
}
