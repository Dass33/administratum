package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type RowParams struct {
	SheetId string `json:"sheet_id"`
	RowIdx  int64  `json:"row_idx"`
}

func (cfg *apiConfig) deleteRowHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := RowParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	sheet_id, err := uuid.Parse(params.SheetId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !cfg.checkSheetPermission(id, sheet_id, "write", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient write permissions")
		return
	}

	deleteRowParams := database.DeleteRowParams{
		SheetID: sheet_id,
		Idx:     params.RowIdx,
	}
	err = cfg.db.DeleteRow(r.Context(), deleteRowParams)
	if err != nil {
		msg := fmt.Sprintf("Row could not be deleted: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}
