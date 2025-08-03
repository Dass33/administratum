package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectData struct {
	Table TableData `table:"table"`
	Sheet Sheet     `sheet:"sheet"`
}

func (cfg *apiConfig) getProjectHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	tableIdStr := chi.URLParam(r, "table_id")
	tableId, err := uuid.Parse(tableIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the table id from url: %s", err)
		respondWithError(w, 400, msg)
	}

	cfg.switchProject(w, r, tableId, userId, 200)
}
