package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getJsonHandler(w http.ResponseWriter, r *http.Request) {
	branchIdStr := chi.URLParam(r, "branch_id")
	branchId, err := uuid.Parse(branchIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id from url: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	sheetsDb, err := cfg.db.GetSheetsFromBranch(r.Context(), branchId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from branch id: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	data := make([]any, 0, len(sheetsDb))
	for _, sheet := range sheetsDb {
		if sheet.Type == SheetTypeMap {
			row, err := cfg.getMapSheetJson(sheet, r.Context())
			if err != nil {
				msg := fmt.Sprintf("Could not get row from map sheet: %s", err)
				respondWithError(w, http.StatusInternalServerError, msg)
				return
			}
			data = append(data, row)
			continue
		}

		rows, err := cfg.getListSheetJson(sheet, r.Context())
		if err != nil {
			msg := fmt.Sprintf("Could not get rows from list sheet: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}

		data = append(data, rows)
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *apiConfig) getColumnsWitRowCount(sheetId uuid.UUID, ctx context.Context) ([]Column, int64, error) {
	columns, err := cfg.GetColumns(sheetId, ctx)
	if err != nil {
		return nil, 0, errors.New("Could not get columns with given sheet id")
	}

	var rowCount int64 = 0
	for i := range columns {
		currLen := int64(len(columns[i].Data))
		if currLen > rowCount {
			rowCount = currLen
		}
	}
	return columns, rowCount, nil
}

func (cfg *apiConfig) getMapSheetJson(sheet database.Sheet, ctx context.Context) (map[string]any, error) {
	columns, rowCount, err := cfg.getColumnsWitRowCount(sheet.ID, ctx)
	if err != nil {
		return nil, err
	}

	if len(columns) < 2 {
		return nil, errors.New("There are not enough columns")
	}

	row := make(map[string]any)

	for i := range rowCount {
		nameCell, ok := getDataAtColIdx(columns[0].Data, i)
		if !ok || !nameCell.Value.Valid {
			continue
		}
		valCell, ok := getDataAtColIdx(columns[1].Data, i)
		if !ok || !valCell.Value.Valid || !valCell.Type.Valid {
			continue
		}

		val, err := ParseValue(valCell.Value.String, valCell.Type.String)
		if err != nil {
			return nil, err
		}
		row[nameCell.Value.String] = val
	}
	return row, nil
}

func (cfg *apiConfig) getListSheetJson(sheet database.Sheet, ctx context.Context) ([]map[string]any, error) {
	columns, rowCount, err := cfg.getColumnsWitRowCount(sheet.ID, ctx)
	if err != nil {
		return nil, err
	}

	rows := make([]map[string]any, 0, rowCount)

	for i := range rowCount {
		row := make(map[string]any)

		for e := range columns {
			col := &columns[e]
			if cell, ok := getDataAtColIdx(col.Data, i); ok {
				val, err := ParseValue(cell.Value.String, col.Type)
				if err != nil {
					return nil, err
				}
				row[col.Name] = val
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// the column data has to be sorted ascending by their index
func getDataAtColIdx(data []ColumnData, idx int64) (ColumnData, bool) {
	i := sort.Search(len(data), func(i int) bool {
		return data[i].Idx >= idx
	})

	if i >= len(data) || data[i].Idx != idx {
		return ColumnData{}, false
	}

	return data[i], true
}

func ParseValue(input string, valueType string) (any, error) {
	switch strings.ToLower(valueType) {
	case "text", "string":
		return input, nil

	case "number", "int", "float":
		return parseNumber(input)

	case "array":
		return parseArray(input)

	case "bool", "boolean":
		return strings.ToLower(input) == "true", nil

	default:
		return nil, fmt.Errorf("unsupported type: %s", valueType)
	}
}

func parseNumber(input string) (any, error) {
	cleaned := strings.ReplaceAll(input, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	cleaned = strings.TrimSpace(cleaned)

	if cleaned == "" {
		return nil, fmt.Errorf("empty number string")
	}
	if intVal, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
		return intVal, nil
	}
	if floatVal, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return floatVal, nil
	}

	return nil, fmt.Errorf("cannot parse '%s' as number", input)
}

func parseArray(input string) (any, error) {
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		var result []any
		if err := json.Unmarshal([]byte(input), &result); err != nil {
			return nil, fmt.Errorf("invalid JSON array: %v", err)
		}
		return result, nil
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{}, nil
	}

	return result, nil
}
