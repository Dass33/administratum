package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type MergePreviewRequest struct {
	SourceBranchID uuid.UUID `json:"source_branch_id"`
}

type MergeConflict struct {
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	SheetID         uuid.UUID `json:"sheet_id"`
	SheetName       string    `json:"sheet_name"`
	ColumnID        uuid.UUID `json:"column_id,omitempty"`
	ColumnName      string    `json:"column_name,omitempty"`
	RowIndex        *int64    `json:"row_index,omitempty"`
	Property        string    `json:"property,omitempty"`
	SourceValue     string    `json:"source_value"`
	TargetValue     string    `json:"target_value"`
	SourceUpdatedAt time.Time `json:"source_updated_at"`
	TargetUpdatedAt time.Time `json:"target_updated_at"`
}

type MergePreviewResponse struct {
	Conflicts []MergeConflict `json:"conflicts"`
}

func (cfg *apiConfig) mergePreviewHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var req MergePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()

	sourceBranch, err := cfg.db.GetBranch(ctx, req.SourceBranchID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Source branch not found")
		return
	}

	targetBranch, err := cfg.db.GetOldestBranchFromTable(ctx, sourceBranch.TableID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find target branch (oldest branch in table)")
		return
	}

	if !cfg.checkBranchPermission(userId, req.SourceBranchID, "read", ctx) {
		respondWithError(w, http.StatusForbidden, "No read permission on source branch")
		return
	}

	if !cfg.checkBranchPermission(userId, targetBranch.ID, "write", ctx) {
		respondWithError(w, http.StatusForbidden, "No write permission on target branch")
		return
	}

	sourceData, err := cfg.db.GetBranchDataForMerge(ctx, req.SourceBranchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get source branch data")
		return
	}

	targetData, err := cfg.db.GetBranchDataForMerge(ctx, targetBranch.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get target branch data")
		return
	}

	conflicts := cfg.detectMergeConflicts(sourceData, targetData, sourceBranch.CreatedAt)

	if conflicts == nil {
		conflicts = []MergeConflict{}
	}

	response := MergePreviewResponse{
		Conflicts: conflicts,
	}

	respondWithJSON(w, http.StatusOK, response)
}

type mergeData struct {
	sheets   map[uuid.UUID]database.GetBranchDataForMergeRow
	columns  map[uuid.UUID]database.GetBranchDataForMergeRow
	cellData map[string]database.GetBranchDataForMergeRow
}

func buildMergeData(data []database.GetBranchDataForMergeRow) mergeData {
	sheets := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	columns := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	cellData := make(map[string]database.GetBranchDataForMergeRow)

	for _, row := range data {
		sheets[row.SheetID] = row
		if row.ColumnID.Valid {
			columns[row.ColumnID.UUID] = row
			if row.ColumnDataID.Valid {
				sheetKey := row.SourceSheetID.String
				columnKey := row.SourceColumnID.String
				key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, row.ColumnDataIdx.Int64)
				cellData[key] = row
			}
		}
	}

	return mergeData{
		sheets:   sheets,
		columns:  columns,
		cellData: cellData,
	}
}

func detectSheetConflicts(sourceData, targetData mergeData, branchCreatedAt time.Time) []MergeConflict {
	var conflicts []MergeConflict

	for sheetId, sourceSheet := range sourceData.sheets {
		targetSheet, exists := targetData.sheets[sheetId]
		if !exists {
			continue
		}

		if sourceSheet.SheetUpdatedAt.After(branchCreatedAt) &&
			targetSheet.SheetUpdatedAt.After(branchCreatedAt) {
			conflicts = append(conflicts, MergeConflict{
				ID:              fmt.Sprintf("sheet-%s", sheetId.String()),
				Type:            "sheet_property",
				SheetID:         sheetId,
				SheetName:       sourceSheet.SheetName,
				Property:        "name",
				SourceValue:     sourceSheet.SheetName,
				TargetValue:     targetSheet.SheetName,
				SourceUpdatedAt: sourceSheet.SheetUpdatedAt,
				TargetUpdatedAt: targetSheet.SheetUpdatedAt,
			})
		}
	}

	return conflicts
}

func detectColumnConflicts(sourceData, targetData mergeData, branchCreatedAt time.Time) []MergeConflict {
	var conflicts []MergeConflict

	for columnId, sourceColumn := range sourceData.columns {
		targetColumn, exists := targetData.columns[columnId]
		if !exists {
			continue
		}

		if sourceColumn.ColumnUpdatedAt.Valid && targetColumn.ColumnUpdatedAt.Valid &&
			sourceColumn.ColumnUpdatedAt.Time.After(branchCreatedAt) &&
			targetColumn.ColumnUpdatedAt.Time.After(branchCreatedAt) {

			if sourceColumn.ColumnName.String != targetColumn.ColumnName.String {
				conflicts = append(conflicts, MergeConflict{
					ID:              fmt.Sprintf("column-%s-name", columnId.String()),
					Type:            "column_property",
					SheetID:         sourceColumn.SheetID,
					SheetName:       sourceColumn.SheetName,
					ColumnID:        columnId,
					ColumnName:      sourceColumn.ColumnName.String,
					Property:        "name",
					SourceValue:     sourceColumn.ColumnName.String,
					TargetValue:     targetColumn.ColumnName.String,
					SourceUpdatedAt: sourceColumn.ColumnUpdatedAt.Time,
					TargetUpdatedAt: targetColumn.ColumnUpdatedAt.Time,
				})
			}

			if sourceColumn.ColumnType.String != targetColumn.ColumnType.String {
				conflicts = append(conflicts, MergeConflict{
					ID:              fmt.Sprintf("column-%s-type", columnId.String()),
					Type:            "column_property",
					SheetID:         sourceColumn.SheetID,
					SheetName:       sourceColumn.SheetName,
					ColumnID:        columnId,
					ColumnName:      sourceColumn.ColumnName.String,
					Property:        "type",
					SourceValue:     sourceColumn.ColumnType.String,
					TargetValue:     targetColumn.ColumnType.String,
					SourceUpdatedAt: sourceColumn.ColumnUpdatedAt.Time,
					TargetUpdatedAt: targetColumn.ColumnUpdatedAt.Time,
				})
			}
		}
	}

	return conflicts
}

func detectCellDataConflicts(sourceData, targetData mergeData, branchCreatedAt time.Time) []MergeConflict {
	var conflicts []MergeConflict

	for key, sourceCell := range sourceData.cellData {
		if targetCell, exists := targetData.cellData[key]; exists {
			if sourceCell.ColumnDataUpdatedAt.Valid && targetCell.ColumnDataUpdatedAt.Valid &&
				sourceCell.ColumnDataUpdatedAt.Time.After(branchCreatedAt) &&
				targetCell.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {

				sourceValue := ""
				if sourceCell.ColumnDataValue.Valid {
					sourceValue = sourceCell.ColumnDataValue.String
				}

				targetValue := ""
				if targetCell.ColumnDataValue.Valid {
					targetValue = targetCell.ColumnDataValue.String
				}

				if sourceValue != targetValue {
					rowIndex := sourceCell.ColumnDataIdx.Int64
					conflicts = append(conflicts, MergeConflict{
						ID:              fmt.Sprintf("cell-%s", key),
						Type:            "cell_data",
						SheetID:         sourceCell.SheetID,
						SheetName:       sourceCell.SheetName,
						ColumnID:        sourceCell.ColumnID.UUID,
						ColumnName:      sourceCell.ColumnName.String,
						RowIndex:        &rowIndex,
						SourceValue:     sourceValue,
						TargetValue:     targetValue,
						SourceUpdatedAt: sourceCell.ColumnDataUpdatedAt.Time,
						TargetUpdatedAt: targetCell.ColumnDataUpdatedAt.Time,
					})
				}
			}
		}
	}

	return conflicts
}

func (cfg *apiConfig) detectMergeConflicts(sourceData, targetData []database.GetBranchDataForMergeRow, branchCreatedAt time.Time) []MergeConflict {
	sourceMergeData := buildMergeData(sourceData)
	targetMergeData := buildMergeData(targetData)

	var conflicts []MergeConflict
	conflicts = append(conflicts, detectSheetConflicts(sourceMergeData, targetMergeData, branchCreatedAt)...)
	conflicts = append(conflicts, detectColumnConflicts(sourceMergeData, targetMergeData, branchCreatedAt)...)
	conflicts = append(conflicts, detectCellDataConflicts(sourceMergeData, targetMergeData, branchCreatedAt)...)

	return conflicts
}
