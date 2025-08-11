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
	TargetBranchID uuid.UUID `json:"target_branch_id"`
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

	_, err = cfg.db.GetBranch(ctx, req.TargetBranchID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Target branch not found")
		return
	}

	if !cfg.checkBranchPermission(userId, req.SourceBranchID, "read", ctx) {
		respondWithError(w, http.StatusForbidden, "No read permission on source branch")
		return
	}

	if !cfg.checkBranchPermission(userId, req.TargetBranchID, "write", ctx) {
		respondWithError(w, http.StatusForbidden, "No write permission on target branch")
		return
	}

	sourceData, err := cfg.db.GetBranchDataForMerge(ctx, req.SourceBranchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get source branch data")
		return
	}

	targetData, err := cfg.db.GetBranchDataForMerge(ctx, req.TargetBranchID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get target branch data")
		return
	}

	conflicts := cfg.detectMergeConflicts(sourceData, targetData, sourceBranch.CreatedAt)

	fmt.Printf("DEBUG merge_preview: Branch created at %v\n", sourceBranch.CreatedAt)
	fmt.Printf("DEBUG merge_preview: Source data rows: %d\n", len(sourceData))
	fmt.Printf("DEBUG merge_preview: Target data rows: %d\n", len(targetData))
	fmt.Printf("DEBUG merge_preview: Detected %d conflicts\n", len(conflicts))

	if conflicts == nil {
		conflicts = []MergeConflict{}
	}

	response := MergePreviewResponse{
		Conflicts: conflicts,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) detectMergeConflicts(sourceData, targetData []database.GetBranchDataForMergeRow, branchCreatedAt time.Time) []MergeConflict {
	conflicts := make([]MergeConflict, 0)

	sourceSheets := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	sourceColumns := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	sourceCellData := make(map[string]database.GetBranchDataForMergeRow)

	targetSheets := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	targetColumns := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	targetCellData := make(map[string]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		sourceSheets[row.SheetID] = row
		if row.ColumnID.Valid {
			sourceColumns[row.ColumnID.UUID] = row
			if row.ColumnDataID.Valid {
				// Use source IDs for matching if available, otherwise fall back to names
				var sheetKey, columnKey string
				if row.SourceSheetID.Valid {
					sheetKey = row.SourceSheetID.String
				} else {
					sheetKey = row.SheetID.String()
				}
				if row.SourceColumnID.Valid {
					columnKey = row.SourceColumnID.String
				} else {
					columnKey = row.ColumnID.UUID.String()
				}
				key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, row.ColumnDataIdx.Int64)
				sourceCellData[key] = row
			}
		}
	}

	for _, row := range targetData {
		targetSheets[row.SheetID] = row
		if row.ColumnID.Valid {
			targetColumns[row.ColumnID.UUID] = row
			if row.ColumnDataID.Valid {
				// For target branch: if it has source references, use them (it's a copy)
				// If not, use its own IDs since it's the original
				var sheetKey, columnKey string
				if row.SourceSheetID.Valid {
					sheetKey = row.SourceSheetID.String
				} else {
					// This is the original target data - use its own ID as the reference
					sheetKey = row.SheetID.String()
				}
				if row.SourceColumnID.Valid {
					columnKey = row.SourceColumnID.String
				} else {
					// This is the original target data - use its own ID as the reference
					columnKey = row.ColumnID.UUID.String()
				}
				key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, row.ColumnDataIdx.Int64)
				targetCellData[key] = row
			}
		}
	}

	for sheetId, sourceSheet := range sourceSheets {
		if targetSheet, exists := targetSheets[sheetId]; exists {
			// Check if both sheets were updated after branch creation
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
	}

	for columnId, sourceColumn := range sourceColumns {
		if targetColumn, exists := targetColumns[columnId]; exists {
			// Check if both columns were updated after branch creation
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
	}

	for key, sourceCell := range sourceCellData {
		if targetCell, exists := targetCellData[key]; exists {
			fmt.Printf("DEBUG: Checking cell %s: source_created=%v, target_created=%v, branch_created=%v\n",
				key, sourceCell.ColumnDataCreatedAt.Time, targetCell.ColumnDataCreatedAt.Time, branchCreatedAt)
			fmt.Printf("DEBUG: source_updated=%v, target_updated=%v\n",
				sourceCell.ColumnDataUpdatedAt.Time, targetCell.ColumnDataUpdatedAt.Time)
			fmt.Printf("DEBUG: source has refs: sheet=%v, column=%v\n", 
				sourceCell.SourceSheetID.Valid, sourceCell.SourceColumnID.Valid)
			fmt.Printf("DEBUG: target has refs: sheet=%v, column=%v\n", 
				targetCell.SourceSheetID.Valid, targetCell.SourceColumnID.Valid)

			// Check for conflicts: both cells modified after branch creation
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
