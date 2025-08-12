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
	Resolutions    []MergeResolution `json:"resolutions"`
}

type MergeExecuteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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

	// Get the oldest branch from the same table as target
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

	// No hierarchical validation needed - we always merge to oldest branch (main)

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
	err = cfg.executeMerge(sourceData, targetData, conflicts, resolutionMap, targetBranch.ID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Merge execution failed: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Merge failed: %v", err))
		return
	}

	// Create new sheets that exist in source but not in target
	fmt.Printf("Creating new sheets from source branch...\n")
	err = cfg.createNewSheets(sourceData, targetData, targetBranch.ID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Failed to create new sheets: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new sheets: %v", err))
		return
	}

	// Refresh target data after creating new sheets so columns can find their target sheets
	fmt.Printf("Refreshing target data after sheet creation...\n")
	targetData, err = cfg.db.GetBranchDataForMerge(ctx, targetBranch.ID)
	if err != nil {
		fmt.Printf("Failed to refresh target data: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to refresh target data: %v", err))
		return
	}

	// Create new columns that exist in source but not in target
	fmt.Printf("Creating new columns from source branch...\n")
	err = cfg.createNewColumns(sourceData, targetData, targetBranch.ID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Failed to create new columns: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create new columns: %v", err))
		return
	}

	// Copy data to newly created columns
	fmt.Printf("Copying data to newly created columns...\n")
	err = cfg.copyDataToNewColumns(sourceData, targetBranch.ID, sourceBranch.CreatedAt, ctx)
	if err != nil {
		fmt.Printf("Failed to copy data to new columns: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to copy data to new columns: %v", err))
		return
	}

	fmt.Printf("Merge execution completed successfully\n")

	// Delete the source branch after successful merge
	fmt.Printf("Deleting source branch %s after successful merge...\n", req.SourceBranchID)
	err = cfg.db.DeleteBranch(ctx, req.SourceBranchID)
	if err != nil {
		fmt.Printf("Warning: Failed to delete source branch: %v\n", err)
		// Don't fail the merge because of this - branch deletion is cleanup
	} else {
		fmt.Printf("Source branch deleted successfully\n")
	}

	response := MergeExecuteResponse{
		Success: true,
		Message: "Merge completed successfully and source branch deleted",
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
				// Parse conflict ID to get the matching key used during conflict detection
				// Format: "cell-{sheet_key}-{column_key}-{row_idx}"
				conflictKey := conflict.ID[5:] // Remove "cell-" prefix
				
				// Find the source row using the conflict key
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
				
				// Find the corresponding target cell using the same matching logic
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
		
		// Check if this is data that was updated after branch creation
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnDataCreatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Valid &&
			sourceRow.ColumnDataUpdatedAt.Time.After(branchCreatedAt) {
			
			fmt.Printf("Row %d matches criteria for non-conflicting change\n", i)
			
			// Check if this change was handled as a conflict by building the same key
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
			fmt.Printf("Row %d wasConflict: %v\n", i, wasConflict)
			
			if !wasConflict {
				fmt.Printf("Looking for target row to match: sheet_key=%s, column_key=%s, idx=%d\n", 
					sourceSheetKey, sourceColumnKey, sourceRow.ColumnDataIdx.Int64)
				
				// Find the corresponding target cell 
				// The key insight: source B references target A directly
				// So we look for a target cell where the target cell's own IDs match source B's references
				found := false
				for j, targetRow := range targetData {
					if targetRow.ColumnID.Valid && targetRow.ColumnDataID.Valid &&
						targetRow.SheetID.String() == sourceSheetKey &&
						targetRow.ColumnID.UUID.String() == sourceColumnKey &&
						targetRow.ColumnDataIdx.Int64 == sourceRow.ColumnDataIdx.Int64 {
						
						fmt.Printf("Found matching target row %d: sheet_id=%s, column_id=%s, idx=%d\n", 
							j, targetRow.SheetID.String(), targetRow.ColumnID.UUID.String(), targetRow.ColumnDataIdx.Int64)
						
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

func (cfg *apiConfig) validateHierarchicalMerge(ctx context.Context, sourceBranchID, targetBranchID uuid.UUID) (bool, error) {
	// Get source branch data to check if it references the target branch
	sourceSheets, err := cfg.db.GetSheetsFromBranch(ctx, sourceBranchID)
	if err != nil {
		return false, fmt.Errorf("could not get source branch sheets: %w", err)
	}

	// Get target branch data to compare sheet IDs
	targetSheets, err := cfg.db.GetSheetsFromBranch(ctx, targetBranchID)
	if err != nil {
		return false, fmt.Errorf("could not get target branch sheets: %w", err)
	}

	// Create a map of target sheet IDs for quick lookup
	targetSheetIDs := make(map[string]bool)
	for _, sheet := range targetSheets {
		targetSheetIDs[sheet.ID.String()] = true
	}

	// Check if source sheets reference target sheets
	// At least one source sheet must have source_sheet_id pointing to a target sheet
	validReferences := 0
	for _, sourceSheet := range sourceSheets {
		fmt.Printf("DEBUG: Source sheet %s (%s) has source_sheet_id: %v (%s)\n", 
			sourceSheet.ID.String(), sourceSheet.Name, sourceSheet.SourceSheetID.Valid, sourceSheet.SourceSheetID.String)
		if sourceSheet.SourceSheetID.Valid {
			if targetSheetIDs[sourceSheet.SourceSheetID.String] {
				fmt.Printf("DEBUG: Valid reference found: source sheet %s references target sheet %s\n",
					sourceSheet.ID.String(), sourceSheet.SourceSheetID.String)
				validReferences++
			}
		}
	}

	fmt.Printf("DEBUG: Found %d valid references out of %d source sheets\n", validReferences, len(sourceSheets))
	
	// If source has sheets but none reference target, it's not a direct child
	if len(sourceSheets) > 0 && validReferences == 0 {
		fmt.Printf("DEBUG: Rejecting merge - no valid references found\n")
		return false, nil
	}

	// Additional validation: check if source was created after target
	// This prevents merging in wrong direction
	sourceBranch, err := cfg.db.GetBranch(ctx, sourceBranchID)
	if err != nil {
		return false, fmt.Errorf("could not get source branch: %w", err)
	}

	targetBranch, err := cfg.db.GetBranch(ctx, targetBranchID)
	if err != nil {
		return false, fmt.Errorf("could not get target branch: %w", err)
	}

	fmt.Printf("DEBUG: Source branch %s (%s) created at: %v\n", sourceBranch.ID.String(), sourceBranch.Name, sourceBranch.CreatedAt)
	fmt.Printf("DEBUG: Target branch %s (%s) created at: %v\n", targetBranch.ID.String(), targetBranch.Name, targetBranch.CreatedAt)

	// Source must be created after target for hierarchical relationship
	if !sourceBranch.CreatedAt.After(targetBranch.CreatedAt) {
		fmt.Printf("DEBUG: Rejecting merge - source not created after target\n")
		return false, nil
	}

	fmt.Printf("DEBUG: Hierarchical merge validation passed\n")

	return true, nil
}

func (cfg *apiConfig) createNewSheets(sourceData, targetData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, branchCreatedAt time.Time, ctx context.Context) error {
	// Get unique sheets from source that are truly new (no source reference)
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
			
			fmt.Printf("Creating new sheet: name=%s, type=%s, source_id=%s\n", 
				sourceSheet.SheetName, sourceSheet.SheetType, sourceSheet.SheetID.String())
			
			_, err := cfg.db.CreateSheet(ctx, createSheetParams)
			if err != nil {
				return fmt.Errorf("failed to create sheet %s: %v", sourceSheet.SheetName, err)
			}
			sheetCount++
		}
	}
	
	fmt.Printf("Created %d new sheets\n", sheetCount)
	return nil
}

func (cfg *apiConfig) createNewColumns(sourceData, targetData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, branchCreatedAt time.Time, ctx context.Context) error {
	// Get source columns that are truly new (no source reference)
	sourceColumns := make(map[uuid.UUID]database.GetBranchDataForMergeRow)
	for _, row := range sourceData {
		if row.ColumnID.Valid && !row.SourceColumnID.Valid {
			sourceColumns[row.ColumnID.UUID] = row
		}
	}

	// Get existing target columns to avoid duplicates
	targetColumns := make(map[string]bool)
	for _, row := range targetData {
		if row.ColumnID.Valid {
			targetColumns[row.ColumnID.UUID.String()] = true
		}
	}

	// Get target sheets for sheet ID mapping (we need to know which target sheet to add columns to)
	targetSheetMap := make(map[string]uuid.UUID) // source sheet ID -> target sheet ID
	for _, row := range targetData {
		if row.SourceSheetID.Valid {
			targetSheetMap[row.SourceSheetID.String] = row.SheetID
		} else {
			targetSheetMap[row.SheetID.String()] = row.SheetID
		}
	}

	columnCount := 0
	for _, sourceColumn := range sourceColumns {
		// Check if target already has this column by column ID
		shouldCreate := !targetColumns[sourceColumn.ColumnID.UUID.String()]
		
		if shouldCreate {
			// Find target sheet - if source sheet has a reference, use that; otherwise use sheet ID
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
				
				fmt.Printf("Creating new column: name=%s, type=%s, sheet_id=%s, source_id=%s\n", 
					sourceColumn.ColumnName.String, sourceColumn.ColumnType.String, 
					targetSheetID.String(), sourceColumn.ColumnID.UUID.String())
				
				_, err := cfg.db.AddColumn(ctx, addColumnParams)
				if err != nil {
					return fmt.Errorf("failed to create column %s: %v", sourceColumn.ColumnName.String, err)
				}
				columnCount++
			} else {
				fmt.Printf("WARNING: Could not find target sheet for new column %s\n", sourceColumn.ColumnName.String)
			}
		}
	}
	
	fmt.Printf("Created %d new columns\n", columnCount)
	return nil
}

func (cfg *apiConfig) copyDataToNewColumns(sourceData []database.GetBranchDataForMergeRow, targetBranchID uuid.UUID, branchCreatedAt time.Time, ctx context.Context) error {
	// Get the updated target data after new columns have been created
	targetData, err := cfg.db.GetBranchDataForMerge(ctx, targetBranchID)
	if err != nil {
		return fmt.Errorf("could not get updated target branch data: %v", err)
	}

	// Create a map of newly created target columns by their source_column_id
	newTargetColumns := make(map[string]uuid.UUID) // source column ID -> new target column ID
	for _, row := range targetData {
		if row.ColumnID.Valid && row.SourceColumnID.Valid {
			newTargetColumns[row.SourceColumnID.String] = row.ColumnID.UUID
		}
	}

	// Copy data from source columns that belong to newly created columns
	copiedDataCount := 0
	for _, sourceRow := range sourceData {
		if sourceRow.ColumnDataID.Valid && sourceRow.ColumnID.Valid {
			// Check if this source column has a corresponding newly created target column
			if targetColumnID, exists := newTargetColumns[sourceRow.ColumnID.UUID.String()]; exists {
				// Check if target already has this data at this index to avoid duplicates
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
					fmt.Printf("Copying data to new column: target_column=%s, idx=%d, value='%s'\n", 
						targetColumnID.String(), sourceRow.ColumnDataIdx.Int64, sourceRow.ColumnDataValue.String)
					
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
					copiedDataCount++
				}
			}
		}
	}
	
	fmt.Printf("Copied %d data entries to new columns\n", copiedDataCount)
	return nil
}