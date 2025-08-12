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

	fmt.Printf("=== MERGE_EXECUTE DEBUG START ===\n")
	fmt.Printf("Source branch ID: %s\n", req.SourceBranchID)
	fmt.Printf("Target branch ID: %s\n", targetBranch.ID)
	fmt.Printf("Source branch created at: %v\n", sourceBranch.CreatedAt)
	fmt.Printf("Source data rows: %d\n", len(sourceData))
	fmt.Printf("Target data rows: %d\n", len(targetData))
	fmt.Printf("Number of resolutions provided: %d\n\n", len(req.Resolutions))

	conflicts := cfg.detectMergeConflicts(sourceData, targetData, sourceBranch.CreatedAt)

	resolutionMap := make(map[string]string)
	for _, resolution := range req.Resolutions {
		resolutionMap[resolution.ConflictID] = resolution.ChosenSource
	}

	if len(conflicts) != len(req.Resolutions) {
		respondWithError(w, http.StatusBadRequest, "All conflicts must be resolved")
		return
	}

	fmt.Printf("Starting merge execution...\n")
	err = cfg.executeMerge(sourceData, targetData, conflicts, resolutionMap, targetBranch.ID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Merge execution failed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Merge failed: %v", err))
		return
	}

	fmt.Printf("Creating new sheets from source branch...\n")
	err = cfg.createNewSheets(sourceData, targetData, targetBranch.ID, ctx)
	if err != nil {
		fmt.Printf("Failed to create new sheets: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new sheets: %v", err))
		return
	}

	fmt.Printf("Refreshing target data after sheet creation...\n")
	targetData, err = cfg.db.GetBranchDataForMerge(ctx, targetBranch.ID)
	if err != nil {
		fmt.Printf("Failed to refresh target data: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to refresh target data: %v", err))
		return
	}

	fmt.Printf("Creating new columns from source branch...\n")
	err = cfg.createNewColumns(sourceData, targetData, ctx)
	if err != nil {
		fmt.Printf("Failed to create new columns: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new columns: %v", err))
		return
	}

	fmt.Printf("Copying data to newly created columns...\n")
	err = cfg.copyDataToNewColumns(sourceData, targetBranch.ID, ctx)
	if err != nil {
		fmt.Printf("Failed to copy data to new columns: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to copy data to new columns: %v", err))
		return
	}

	fmt.Printf("Merge execution completed successfully\n")

	fmt.Printf("Deleting source branch %s after successful merge...\n", req.SourceBranchID)
	err = cfg.db.DeleteBranch(ctx, req.SourceBranchID)
	if err != nil {
		fmt.Printf("Warning: Failed to delete source branch: %v\n", err)
	} else {
		fmt.Printf("Source branch deleted successfully\n")
	}

	response := MergeExecuteResponse{
		Success:        true,
		Message:        "Merge completed successfully and source branch deleted",
		TargetBranchID: targetBranch.ID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) executeMerge(sourceData, targetData []database.GetBranchDataForMergeRow, conflicts []MergeConflict, resolutions map[string]string, targetBranchID uuid.UUID, branchCreatedAt time.Time, ctx context.Context) error {
	fmt.Printf("=== EXECUTE_MERGE DEBUG START ===\n")
	fmt.Printf("Target branch ID: %s\n", targetBranchID)
	fmt.Printf("Branch created at: %v\n", branchCreatedAt)
	fmt.Printf("Processing %d conflicts\n", len(conflicts))
	fmt.Printf("Processing %d resolutions\n", len(resolutions))
	sourceMap := make(map[string]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		if row.ColumnDataID.Valid {
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
			sourceMap[key] = row
		}
	}

	for i, conflict := range conflicts {
		fmt.Printf("Processing conflict %d: %s\n", i+1, conflict.ID)
		resolution, exists := resolutions[conflict.ID]
		if !exists {
			fmt.Printf("ERROR: No resolution provided for conflict %s\n", conflict.ID)
			return fmt.Errorf("no resolution provided for conflict %s", conflict.ID)
		}

		fmt.Printf("Resolution for %s: %s\n", conflict.ID, resolution)

		if resolution != "source" && resolution != "target" {
			fmt.Printf("ERROR: Invalid resolution %s for conflict %s\n", resolution, conflict.ID)
			return fmt.Errorf("invalid resolution %s for conflict %s", resolution, conflict.ID)
		}

		if resolution == "source" {
			fmt.Printf("Applying source resolution for conflict type: %s\n", conflict.Type)
			switch conflict.Type {
			case "cell_data":
				conflictKey := conflict.ID[5:] // Remove "cell-" prefix

				var sourceRow database.GetBranchDataForMergeRow
				var sourceExists bool
				for _, row := range sourceData {
					if row.ColumnDataID.Valid {
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
						if key == conflictKey {
							sourceRow = row
							sourceExists = true
							break
						}
					}
				}

				for _, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid {
						var sheetKey, columnKey string
						if targetRow.SourceSheetID.Valid {
							sheetKey = targetRow.SourceSheetID.String
						} else {
							sheetKey = targetRow.SheetID.String()
						}
						if targetRow.SourceColumnID.Valid {
							columnKey = targetRow.SourceColumnID.String
						} else {
							columnKey = targetRow.ColumnID.UUID.String()
						}
						key := fmt.Sprintf("%s-%s-%d", sheetKey, columnKey, targetRow.ColumnDataIdx.Int64)

						if key == conflictKey {
							var valueToUse sql.NullString
							if sourceExists {
								valueToUse = sourceRow.ColumnDataValue
							} else {
								valueToUse = sql.NullString{String: conflict.SourceValue, Valid: true}
							}

							fmt.Printf("Updating target cell with source value: target_id=%s, source_value=%s\n",
								targetRow.ColumnDataID.UUID, valueToUse.String)

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

			case "column_property":
				if conflict.Property == "name" {
					err := cfg.db.UpdateColumn(ctx, database.UpdateColumnParams{
						ID:   conflict.ColumnID,
						Name: conflict.SourceValue,
					})
					if err != nil {
						return fmt.Errorf("failed to update column name: %v", err)
					}
				}

			case "sheet_property":
				if conflict.Property == "name" {
					err := cfg.db.RenameSheet(ctx, database.RenameSheetParams{
						ID:   conflict.SheetID,
						Name: conflict.SourceValue,
					})
					if err != nil {
						return fmt.Errorf("failed to update sheet name: %v", err)
					}
				}
			}
		}
	}

	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {

			wasConflict := false
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

	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataCreatedAt.Time.After(branchCreatedAt) {

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
					if targetRow.SourceSheetID.Valid {
						targetSheetKey = targetRow.SourceSheetID.String
					} else {
						targetSheetKey = targetRow.SheetID.String()
					}
					if targetRow.SourceColumnID.Valid {
						targetColumnKey = targetRow.SourceColumnID.String
					} else {
						targetColumnKey = targetRow.ColumnID.UUID.String()
					}

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
	return nil
}

func (cfg *apiConfig) createNewSheets(sourceData, targetData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, ctx context.Context) error {
	sourceSheets := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		if !row.SourceSheetID.Valid {
			sourceSheets[row.SheetID] = row
		}
	}

	// Get existing target sheets to avoid duplicates
	targetSheets := make(map[string]bool)
	for _, row := range targetData {
		targetSheets[row.SheetID.String()] = true
	}

	sheetCount := 0
	for _, sourceSheet := range sourceSheets {
		// Check if target already has this sheet by sheet ID
		shouldCreate := !targetSheets[sourceSheet.SheetID.String()]

		if shouldCreate {
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
			sheetCount++
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
		if row.SourceSheetID.Valid {
			targetSheetMap[row.SourceSheetID.String] = row.SheetID
		} else {
			targetSheetMap[row.SheetID.String()] = row.SheetID
		}
	}

	for _, sourceColumn := range sourceColumns {
		shouldCreate := !targetColumns[sourceColumn.ColumnID.UUID.String()]

		if shouldCreate {
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
