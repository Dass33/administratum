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
	ConflictID    string `json:"conflict_id"`
	ChosenSource  string `json:"chosen_source"`
}

type MergeExecuteRequest struct {
	SourceBranchID uuid.UUID         `json:"source_branch_id"`
	TargetBranchID uuid.UUID         `json:"target_branch_id"`
	Resolutions    []MergeResolution `json:"resolutions"`
}

type MergeExecuteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

	fmt.Printf("=== MERGE_EXECUTE DEBUG START ===\n")
	fmt.Printf("Source branch ID: %s\n", req.SourceBranchID)
	fmt.Printf("Target branch ID: %s\n", req.TargetBranchID)
	fmt.Printf("Source branch created at: %v\n", sourceBranch.CreatedAt)
	fmt.Printf("Source data rows: %d\n", len(sourceData))
	fmt.Printf("Target data rows: %d\n", len(targetData))
	fmt.Printf("Number of resolutions provided: %d\n", len(req.Resolutions))

	conflicts := cfg.detectMergeConflicts(sourceData, targetData, sourceBranch.CreatedAt)

	fmt.Printf("Found %d conflicts for merge\n", len(conflicts))
	for i, conflict := range conflicts {
		fmt.Printf("Conflict %d: ID=%s, Type=%s, SourceValue=%s, TargetValue=%s\n", 
			i+1, conflict.ID, conflict.Type, conflict.SourceValue, conflict.TargetValue)
	}

	resolutionMap := make(map[string]string)
	for _, resolution := range req.Resolutions {
		resolutionMap[resolution.ConflictID] = resolution.ChosenSource
	}

	if len(conflicts) != len(req.Resolutions) {
		respondWithError(w, http.StatusBadRequest, "All conflicts must be resolved")
		return
	}

	fmt.Printf("Starting merge execution...\n")
	err = cfg.executeMerge(sourceData, targetData, conflicts, resolutionMap, req.TargetBranchID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Merge execution failed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Merge failed: %v", err))
		return
	}
	fmt.Printf("Merge execution completed successfully\n")

	response := MergeExecuteResponse{
		Success: true,
		Message: "Merge completed successfully",
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
			key := fmt.Sprintf("%s-%d", row.ColumnID.UUID.String(), row.ColumnDataIdx.Int64)
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
				// Find the target cell and update it with the source value
				key := fmt.Sprintf("%s-%d", conflict.ColumnID.String(), *conflict.RowIndex)
				sourceRow, sourceExists := sourceMap[key]
				
				// Find the corresponding target cell
				for _, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid &&
						targetRow.ColumnID.UUID == conflict.ColumnID &&
						targetRow.ColumnDataIdx.Int64 == *conflict.RowIndex {
						
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

	// Apply non-conflicting changes from source branch
	fmt.Printf("Applying non-conflicting changes from source branch...\n")
	fmt.Printf("Branch created at: %v\n", branchCreatedAt)
	changeCount := 0
	for i, sourceRow := range sourceData {
		fmt.Printf("Source row %d: ColumnDataID.Valid=%v, CreatedAt=%v, UpdatedAt=%v\n", 
			i, sourceRow.ColumnDataID.Valid, sourceRow.ColumnDataCreatedAt.Time, sourceRow.ColumnDataUpdatedAt.Time)
		fmt.Printf("Source row %d: SourceSheetID.Valid=%v (%s), SourceColumnID.Valid=%v (%s)\n", 
			i, sourceRow.SourceSheetID.Valid, sourceRow.SourceSheetID.String, 
			sourceRow.SourceColumnID.Valid, sourceRow.SourceColumnID.String)
		
		// Check if this is copied data (has source references) that was updated after branch creation
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Valid &&
			sourceRow.SourceSheetID.Valid && sourceRow.SourceColumnID.Valid &&
			sourceRow.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {
			
			fmt.Printf("Row %d matches criteria for non-conflicting change\n", i)
			
			// Check if this change was handled as a conflict
			wasConflict := false
			for _, conflict := range conflicts {
				if conflict.Type == "cell_data" && 
					conflict.ColumnID == sourceRow.ColumnID.UUID &&
					conflict.RowIndex != nil && *conflict.RowIndex == sourceRow.ColumnDataIdx.Int64 {
					wasConflict = true
					break
				}
			}
			fmt.Printf("Row %d wasConflict: %v\n", i, wasConflict)
			
			if !wasConflict {
				// Build source key using source IDs for matching
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
				
				fmt.Printf("Looking for target row to match: sheet_key=%s, column_key=%s, idx=%d\n", 
					sourceSheetKey, sourceColumnKey, sourceRow.ColumnDataIdx.Int64)
				
				// Find the corresponding target cell by matching source reference IDs
				found := false
				for j, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid {
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
						
						fmt.Printf("Target row %d: sheet_key=%s, column_key=%s, idx=%d\n", 
							j, targetSheetKey, targetColumnKey, targetRow.ColumnDataIdx.Int64)
						
						if targetSheetKey == sourceSheetKey &&
							targetColumnKey == sourceColumnKey &&
							targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {
						
						fmt.Printf("Updating non-conflicting cell: column=%s, idx=%d, value=%s\n", 
							sourceRow.ColumnID.UUID, sourceRow.ColumnDataIdx.Int64, sourceRow.ColumnDataValue.String)
						
						updateParams := database.UpdateColumnDataParams{
							ID:    targetRow.ColumnDataID.UUID,
							Value: sourceRow.ColumnDataValue,
						}
						err := cfg.db.UpdateColumnData(ctx, updateParams)
						if err != nil {
							return fmt.Errorf("failed to update non-conflicting cell data: %v", err)
						}
						changeCount++
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("No matching target row found for source row %d\n", i)
				}
			}
		}
		}
	}
	fmt.Printf("Applied %d non-conflicting changes\n", changeCount)

	// Copy new data from source branch (data created after the branch was created)
	fmt.Printf("Copying new data from source branch...\n")
	newDataCount := 0
	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataCreatedAt.Time.After(branchCreatedAt) {

			fmt.Printf("Found source data created after branch: created=%v, branch=%v\n", 
				sourceRow.ColumnDataCreatedAt.Time, branchCreatedAt)

			// Build source key using source IDs for matching
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

					// Check if we found the matching target column
					if targetSheetKey == sourceSheetKey && targetColumnKey == sourceColumnKey {
						targetColumnID = targetRow.ColumnID.UUID
						
						// Check if this specific row already exists
						if targetRow.ColumnDataID.Valid && targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {
							found = true
							break
						}
					}
				}
			}

			if !found && targetColumnID != uuid.Nil {
				// Find the maximum index for this column in the target branch to avoid duplicates
				maxIdx := int64(-1)
				for _, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid && 
						targetRow.ColumnID.UUID == targetColumnID &&
						targetRow.ColumnDataIdx.Int64 > maxIdx {
						maxIdx = targetRow.ColumnDataIdx.Int64
					}
				}
				nextIdx := maxIdx + 1
				
				fmt.Printf("Creating new column data: target_column=%s, source_idx=%d, new_idx=%d, value='%s', valueValid=%v\n", 
					targetColumnID, sourceRow.ColumnDataIdx.Int64, nextIdx, sourceRow.ColumnDataValue.String, sourceRow.ColumnDataValue.Valid)
				params := database.CreateColumnDataParams{
					Idx:      nextIdx,
					Value:    sourceRow.ColumnDataValue,
					Type:     sql.NullString{Valid: false},
					ColumnID: targetColumnID,
				}
				_, err := cfg.db.CreateColumnData(ctx, params)
				if err != nil {
					fmt.Printf("Failed to create column data: %v\n", err)
					return fmt.Errorf("failed to create new column data: %v", err)
				}
				newDataCount++
			} else if targetColumnID == uuid.Nil {
				fmt.Printf("WARNING: Could not find target column for new data (source sheet=%s, column=%s)\n", 
					sourceSheetKey, sourceColumnKey)
			} else {
				fmt.Printf("Row already exists or no target column found for idx=%d\n", sourceRow.ColumnDataIdx.Int64)
			}
		}
	}
	fmt.Printf("Created %d new data entries\n", newDataCount)

	return nil
}