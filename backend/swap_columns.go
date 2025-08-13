package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type SwapColumnsParams struct {
	ColumnID1 uuid.UUID `json:"column_id1"`
	ColumnID2 uuid.UUID `json:"column_id2"`
}

func (cfg *apiConfig) swapColumnsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := SwapColumnsParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding columns IDs: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if params.ColumnID1 == uuid.Nil || params.ColumnID2 == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Both column IDs must be provided")
		return
	}

	if params.ColumnID1 == params.ColumnID2 {
		respondWithError(w, http.StatusBadRequest, "Cannot swap a column with itself")
		return
	}

	columns, err := cfg.db.GetColumnOrderIndexes(r.Context(), database.GetColumnOrderIndexesParams{
		ID:     params.ColumnID1,
		ID_2:   params.ColumnID2,
		UserID: id,
	})
	if err != nil {
		msg := fmt.Sprintf("Could not get columns: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	if len(columns) != 2 {
		respondWithError(w, http.StatusBadRequest, "Could not find both columns or insufficient permissions")
		return
	}

	var col1Order, col2Order int64
	for _, col := range columns {
		switch col.ID {
		case params.ColumnID1:
			col1Order = col.OrderIndex
		case params.ColumnID2:
			col2Order = col.OrderIndex
		}
	}

	_, err = cfg.db.SwapColumnsWithPermissionCheck(r.Context(), database.SwapColumnsWithPermissionCheckParams{
		ID:           params.ColumnID1,
		ID_2:         params.ColumnID2,
		OrderIndex:   col2Order,
		OrderIndex_2: col1Order,
		UserID:       id,
	})
	if err != nil {
		msg := fmt.Sprintf("columns could not be swapped: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	respondWithJSON(w, http.StatusOK, "")
}
