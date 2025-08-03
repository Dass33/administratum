package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteColumnHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := ColParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	sheet_id, err := uuid.Parse(params.Sheet_id)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	deleteColumnParams := database.DeleteColumnParams{
		Name:    params.Col.Name,
		SheetID: sheet_id,
	}
	err = cfg.db.DeleteColumn(r.Context(), deleteColumnParams)
	if err != nil {
		msg := fmt.Sprintf("Column could not be deleted: %s", err)
		respondWithError(w, 500, msg)
	}
	respondWithJSON(w, 200, "")
}
