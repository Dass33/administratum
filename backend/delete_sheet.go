package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type deleteSheetParams struct {
	SheetId string `json:"SheetId"`
}

func (cfg *apiConfig) deleteSheetHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := deleteSheetParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	sheetId, err := uuid.Parse(params.SheetId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	err = cfg.db.DeleteSheet(r.Context(), sheetId)
	if err != nil {
		msg := fmt.Sprintf("Sheet could not be deleted: %s", err)
		respondWithError(w, 500, msg)
		return
	}
	respondWithJSON(w, 204, "")
}
