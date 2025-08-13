package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type MergeResolution struct {
	ConflictID   string `json:"conflict_id"`
	ChosenSource string `json:"chosen_source"`
}

type MergeExecuteRequest struct {
	SourceBranchID uuid.UUID         `json:"source_branch_id"`
	Resolutions    []MergeResolution `json:"resolutions"`
}

type MergeExecuteResponse struct {
	Success        bool      `json:"success"`
	Message        string    `json:"message"`
	TargetBranchID uuid.UUID `json:"target_branch_id"`
}

func (cfg *apiConfig) mergeExecuteHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var req MergeExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()

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

	resolutionMap := make(map[string]string)
	for _, resolution := range req.Resolutions {
		resolutionMap[resolution.ConflictID] = resolution.ChosenSource
	}

	if len(conflicts) != len(req.Resolutions) {
		respondWithError(w, http.StatusBadRequest, "All conflicts must be resolved")
		return
	}

	err = cfg.executeMerge(sourceData, targetData, conflicts, resolutionMap, sourceBranch.CreatedAt, ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Merge failed: %v", err))
		return
	}

	err = cfg.createNewSheets(sourceData, targetData, targetBranch.ID, ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new sheets: %v", err))
		return
	}

	targetData, err = cfg.db.GetBranchDataForMerge(ctx, targetBranch.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to refresh target data: %v", err))
		return
	}

	err = cfg.createNewColumns(sourceData, targetData, ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new columns: %v", err))
		return
	}

	err = cfg.copyDataToNewColumns(sourceData, targetBranch.ID, ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to copy data to new columns: %v", err))
		return
	}

	err = cfg.handleDeletions(sourceData, targetData, ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to handle deletions: %v", err))
		return
	}

	cfg.db.DeleteBranch(ctx, req.SourceBranchID)

	response := MergeExecuteResponse{
		Success:        true,
		Message:        "Merge completed successfully and source branch deleted",
		TargetBranchID: targetBranch.ID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) executeMerge(
	sourceData, targetData []database.GetBranchDataForMergeRow,
	conflicts []MergeConflict,
	resolutions map[string]string,
	branchCreatedAt time.Time,
	ctx context.Context,
) error {
	err := cfg.processConflictResolutions(sourceData, targetData, conflicts, resolutions, ctx)
	if err != nil {
		return fmt.Errorf("failed to process conflict resolutions: %v", err)
	}

	err = cfg.updateNonConflictingData(sourceData, targetData, conflicts, branchCreatedAt, ctx)
	if err != nil {
		return fmt.Errorf("failed to update non-conflicting data: %v", err)
	}

	err = cfg.createNewCellData(sourceData, targetData, branchCreatedAt, ctx)
	if err != nil {
		return fmt.Errorf("failed to create new cell data: %v", err)
	}

	return nil
}

func (cfg *apiConfig) processConflictResolutions(
	sourceData, targetData []database.GetBranchDataForMergeRow,
	conflicts []MergeConflict,
	resolutions map[string]string,
	ctx context.Context,
) error {
	for _, conflict := range conflicts {
		resolution, exists := resolutions[conflict.ID]
		if !exists {
			return fmt.Errorf("no resolution provided for conflict %s", conflict.ID)
		}

		if resolution != "source" && resolution != "target" {
			return fmt.Errorf("invalid resolution %s for conflict %s", resolution, conflict.ID)
		}

		if resolution == "source" {
			switch conflict.Type {
			case "cell_data":
				err := cfg.resolveCellDataConflict(sourceData, targetData, conflict, ctx)
				if err != nil {
					return err
				}
			case "column_property":
				err := cfg.resolveColumnPropertyConflict(conflict, ctx)
				if err != nil {
					return err
				}
			case "sheet_property":
				err := cfg.resolveSheetPropertyConflict(conflict, ctx)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (cfg *apiConfig) resolveCellDataConflict(
	sourceData, targetData []database.GetBranchDataForMergeRow,
	conflict MergeConflict,
	ctx context.Context,
) error {
	conflictKey := conflict.ID[5:] // Remove "cell-" prefix

	var sourceRow database.GetBranchDataForMergeRow
	var sourceExists bool
	for _, row := range sourceData {
		if row.ColumnDataID.Valid {
			sheetKey := row.SourceSheetID.String
			columnKey := row.SourceColumnID.String
			key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, row.ColumnDataIdx.Int64)
			if key == conflictKey {
				sourceRow = row
				sourceExists = true
				break
			}
		}
	}

	for _, targetRow := range targetData {
		if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid {
			sheetKey := targetRow.SourceSheetID.String
			columnKey := targetRow.SourceColumnID.String
			key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, targetRow.ColumnDataIdx.Int64)

			if key == conflictKey {
				var valueToUse sql.NullString
				if sourceExists {
					valueToUse = sourceRow.ColumnDataValue
				} else {
					valueToUse = sql.NullString{String: conflict.SourceValue, Valid: true}
				}

				updateParams := database.UpdateColumnDataParams{
					ID:    targetRow.ColumnDataID.UUID,
					Value: valueToUse,
				}
				err := cfg.db.UpdateColumnData(ctx, updateParams)
				if err != nil {
					return fmt.Errorf("failed to update cell data: %v", err)
				}
				break
			}
		}
	}
	return nil
}

func (cfg *apiConfig) resolveColumnPropertyConflict(conflict MergeConflict, ctx context.Context) error {
	if conflict.Property == "name" {
		err := cfg.db.UpdateColumn(ctx, database.UpdateColumnParams{
			ID:   conflict.ColumnID,
			Name: conflict.SourceValue,
		})
		if err != nil {
			return fmt.Errorf("failed to update column name: %v", err)
		}
	}
	return nil
}

func (cfg *apiConfig) resolveSheetPropertyConflict(conflict MergeConflict, ctx context.Context) error {
	if conflict.Property == "name" {
		err := cfg.db.RenameSheet(ctx, database.RenameSheetParams{
			ID:   conflict.SheetID,
			Name: conflict.SourceValue,
		})
		if err != nil {
			return fmt.Errorf("failed to update sheet name: %v", err)
		}
	}
	return nil
}

func (cfg *apiConfig) updateNonConflictingData(
	sourceData, targetData []database.GetBranchDataForMergeRow,
	conflicts []MergeConflict,
	branchCreatedAt time.Time,
	ctx context.Context,
) error {
	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {

			wasConflict := false
			sourceSheetKey := sourceRow.SourceSheetID.String
			sourceColumnKey := sourceRow.SourceColumnID.String
			sourceKey := fmt.Sprintf("%s-%s-%d", sourceSheetKey, sourceColumnKey, sourceRow.ColumnDataIdx.Int64)
			expectedConflictID := fmt.Sprintf("cell-%s", sourceKey)

			for _, conflict := range conflicts {
				if conflict.Type == "cell_data" && conflict.ID == expectedConflictID {
					wasConflict = true
					break
				}
			}

			if !wasConflict {
				for _, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid &&
						targetRow.SheetID.String() == sourceSheetKey &&
						targetRow.ColumnID.UUID.String() == sourceColumnKey &&
						targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {

						updateParams := database.UpdateColumnDataParams{
							ID:    targetRow.ColumnDataID.UUID,
							Value: sourceRow.ColumnDataValue,
						}
						err := cfg.db.UpdateColumnData(ctx, updateParams)
						if err != nil {
							return fmt.Errorf("failed to update non-conflicting cell data: %v", err)
						}
						break
					}
				}
			}
		}
	}
	return nil
}

func (cfg *apiConfig) createNewCellData(
	sourceData, targetData []database.GetBranchDataForMergeRow,
	branchCreatedAt time.Time,
	ctx context.Context,
) error {
	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid {

			var sourceSheetKey, sourceColumnKey string
			if sourceRow.SourceSheetID.Valid {
				sourceSheetKey = sourceRow.SourceSheetID.String
			} else {
				sourceSheetKey = sourceRow.SheetID.String()
			}
			if sourceRow.SourceColumnID.Valid {
				sourceColumnKey = sourceRow.SourceColumnID.String
			} else {
				sourceColumnKey = sourceRow.ColumnID.UUID.String()
			}

			found := false
			var targetColumnID uuid.UUID
			for _, targetRow := range targetData {
				if targetRow.ColumnID.Valid {
					var targetSheetKey, targetColumnKey string
					targetSheetKey = targetRow.SheetID.String()
					targetColumnKey = targetRow.ColumnID.UUID.String()

					if targetSheetKey == sourceSheetKey && targetColumnKey == sourceColumnKey {
						targetColumnID = targetRow.ColumnID.UUID

						if targetRow.ColumnDataID.Valid && targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {
							found = true
							break
						}
					}
				}
			}

			if !found && targetColumnID != uuid.Nil {
				shouldMerge := false

				if sourceRow.ColumnDataCreatedAt.Time.After(branchCreatedAt) {
					shouldMerge = true
				}

				if sourceRow.ColumnDataUpdatedAt.Valid &&
					sourceRow.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {
					shouldMerge = true
				}

				if shouldMerge {
					maxIdx := int64(-1)
					for _, targetRow := range targetData {
						if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid &&
							targetRow.ColumnID.UUID == targetColumnID &&
							targetRow.ColumnDataIdx.Int64 > maxIdx {
							maxIdx = targetRow.ColumnDataIdx.Int64
						}
					}
					nextIdx := maxIdx + 1

					params := database.CreateColumnDataParams{
						Idx:      nextIdx,
						Value:    sourceRow.ColumnDataValue,
						Type:     sql.NullString{Valid: false},
						ColumnID: targetColumnID,
					}
					_, err := cfg.db.CreateColumnData(ctx, params)
					if err != nil {
						return fmt.Errorf("failed to create new column data: %v", err)
					}
				}
			}
		}
	}
	return nil
}

func (cfg *apiConfig) createNewSheets(sourceData, targetData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, ctx context.Context) error {
	sourceSheets := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		if !row.SourceSheetID.Valid {
			sourceSheets[row.SheetID] = row
		}
	}

	targetSheets := make(map[string]bool)
	for _, row := range targetData {
		targetSheets[row.SheetID.String()] = true
	}

	for _, sourceSheet := range sourceSheets {

		if targetSheets[sourceSheet.SheetID.String()] {
			continue
		}
		createSheetParams := database.CreateSheetParams{
			BranchID:      targetBranchID,
			Name:          sourceSheet.SheetName,
			Type:          sourceSheet.SheetType,
			SourceSheetID: sql.NullString{String: sourceSheet.SheetID.String(), Valid: true},
		}

		_, err := cfg.db.CreateSheet(ctx, createSheetParams)
		if err != nil {
			return fmt.Errorf("failed to create sheet %s: %v", sourceSheet.SheetName, err)
		}
	}

	return nil
}

func (cfg *apiConfig) createNewColumns(sourceData, targetData []database.GetBranchDataForMergeRow, ctx context.Context) error {
	sourceColumns := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		if row.ColumnID.Valid && !row.SourceColumnID.Valid {
			sourceColumns[row.ColumnID.UUID] = row
		}
	}

	targetColumns := make(map[string]bool)
	for _, row := range targetData {
		if row.ColumnID.Valid {
			targetColumns[row.ColumnID.UUID.String()] = true
		}
	}

	targetSheetMap := make(map[string]uuid.UUID) // source sheet ID -> target sheet ID
	for _, row := range targetData {
		targetSheetMap[row.SheetID.String()] = row.SheetID
	}

	for _, sourceColumn := range sourceColumns {
		if targetColumns[sourceColumn.ColumnID.UUID.String()] {
			continue
		}

		var targetSheetID uuid.UUID
		var found bool

		if sourceColumn.SourceSheetID.Valid {
			targetSheetID, found = targetSheetMap[sourceColumn.SourceSheetID.String]
		} else {
			targetSheetID, found = targetSheetMap[sourceColumn.SheetID.String()]
		}

		if found {
			addColumnParams := database.AddColumnParams{
				Name:           sourceColumn.ColumnName.String,
				Type:           sourceColumn.ColumnType.String,
				Required:       sourceColumn.ColumnRequired.Bool,
				SheetID:        targetSheetID,
				SourceColumnID: sql.NullString{String: sourceColumn.ColumnID.UUID.String(), Valid: true},
			}
			_, err := cfg.db.AddColumn(ctx, addColumnParams)
			if err != nil {
				return fmt.Errorf("failed to create column %s: %v", sourceColumn.ColumnName.String, err)
			}
		}
	}

	return nil
}

func (cfg *apiConfig) copyDataToNewColumns(sourceData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, ctx context.Context) error {
	targetData, err := cfg.db.GetBranchDataForMerge(ctx, targetBranchID)
	if err != nil {
		return fmt.Errorf("could not get updated target branch data: %v", err)
	}

	newTargetColumns := make(map[string]uuid.UUID) // source column ID -> new target column ID
	for _, row := range targetData {
		if row.ColumnID.Valid && row.SourceColumnID.Valid {
			newTargetColumns[row.SourceColumnID.String] = row.ColumnID.UUID
		}
	}

	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnID.Valid {
			if targetColumnID, exists := newTargetColumns[sourceRow.ColumnID.UUID.String()]; exists {
				hasData := false
				for _, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid &&
						targetRow.ColumnID.UUID == targetColumnID &&
						targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {
						hasData = true
						break
					}
				}

				if !hasData {
					params := database.CreateColumnDataParams{
						Idx:      sourceRow.ColumnDataIdx.Int64,
						Value:    sourceRow.ColumnDataValue,
						Type:     sql.NullString{Valid: false},
						ColumnID: targetColumnID,
					}
					_, err := cfg.db.CreateColumnData(ctx, params)
					if err != nil {
						return fmt.Errorf("failed to copy column data: %v", err)
					}
				}
			}
		}
	}
	return nil
}

func (cfg *apiConfig) handleDeletions(sourceData, targetData []database.GetBranchDataForMergeRow, ctx context.Context) error {
	err := cfg.handleColumnDeletions(sourceData, targetData, ctx)
	if err != nil {
		return fmt.Errorf("failed to handle column deletions: %v", err)
	}

	err = cfg.handleSheetDeletions(sourceData, targetData, ctx)
	if err != nil {
		return fmt.Errorf("failed to handle sheet deletions: %v", err)
	}

	return nil
}

func (cfg *apiConfig) handleColumnDeletions(sourceData, targetData []database.GetBranchDataForMergeRow, ctx context.Context) error {
	referencedTargetColumns := make(map[string]bool)
	for _, row := range sourceData {
		if row.ColumnID.Valid && row.SourceColumnID.Valid {
			referencedTargetColumns[row.SourceColumnID.String] = true
		}
	}

	processedColumns := make(map[uuid.UUID]bool)
	for _, targetRow := range targetData {
		if targetRow.ColumnID.Valid && !processedColumns[targetRow.ColumnID.UUID] {
			processedColumns[targetRow.ColumnID.UUID] = true

			targetColumnId := targetRow.ColumnID.UUID.String()

			if !referencedTargetColumns[targetColumnId] {
				deleteParams := database.DeleteColumnParams{
					Name:    targetRow.ColumnName.String,
					SheetID: targetRow.SheetID,
				}
				err := cfg.db.DeleteColumn(ctx, deleteParams)
				if err != nil {
					return fmt.Errorf("failed to delete column %s: %v", targetRow.ColumnName.String, err)
				}
			}
		}
	}

	return nil
}

func (cfg *apiConfig) handleSheetDeletions(sourceData, targetData []database.GetBranchDataForMergeRow, ctx context.Context) error {
	referencedTargetSheets := make(map[string]bool)
	for _, row := range sourceData {
		if row.SourceSheetID.Valid {
			referencedTargetSheets[row.SourceSheetID.String] = true
		}
	}

	processedSheets := make(map[uuid.UUID]bool)
	for _, targetRow := range targetData {
		if processedSheets[targetRow.SheetID] {
			continue
		}
		processedSheets[targetRow.SheetID] = true

		targetSheetId := targetRow.SheetID.String()

		if !referencedTargetSheets[targetSheetId] {
			err := cfg.db.DeleteSheet(ctx, targetRow.SheetID)
			if err != nil {
				return fmt.Errorf("failed to delete sheet %s: %v", targetRow.SheetName, err)
			}
		}
	}

	return nil
}
