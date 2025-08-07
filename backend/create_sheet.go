package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

const SheetTypeMap = "map"
const SheetTypeList = "list"

func (cfg *apiConfig) createSheetHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := database.CreateSheetParams{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	sheetId := uuid.UUID{}
	fmt.Println(params.Type)
	if params.Type == SheetTypeMap {
		sheetId, err = cfg.createMapSheet(r.Context(), params.Name, params.BranchID)
		if err != nil {
			msg := fmt.Sprintf("Could not create map sheet: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		sheet, err := cfg.db.CreateSheet(r.Context(), params)
		if err != nil {
			msg := fmt.Sprintf("Could not create list sheet: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
		sheetId = sheet.ID
	}

	optionalSheetId := uuid.NullUUID{
		UUID:  sheetId,
		Valid: true,
	}
	sheetData, err := cfg.GetSheet(optionalSheetId, r.Context())
	if err != nil {
		msg := fmt.Sprintf("Could not get sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	setOpenedSheetParams := database.SetOpenedSheetParams{
		ID:          userId,
		OpenedSheet: optionalSheetId,
	}
	err = cfg.db.SetOpenedSheet(r.Context(), setOpenedSheetParams)
	if err != nil {
		msg := fmt.Sprintf("Could not set opened sheet: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	respondWithJSON(w, http.StatusCreated, sheetData)
}

func (cfg *apiConfig) createMapSheet(ctx context.Context, name string, branchID uuid.UUID) (uuid.UUID, error) {
	createMapSheetParams := database.CreateMapSheetParams{
		Name:     name,
		BranchID: branchID,
	}
	sheetId, err := cfg.db.CreateMapSheet(ctx, createMapSheetParams)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = cfg.db.CreateMapSheetColumns(ctx, sheetId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sheetId, nil
}
