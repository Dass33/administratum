package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getJsonHandler(w http.ResponseWriter, r *http.Request) {
	branchIdStr := chi.URLParam(r, "branch_id")
	branchId, err := uuid.Parse(branchIdStr)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the sheet id from url: %s", err)
		respondWithError(w, 400, msg)
		return
	}

	sheetsDb, err := cfg.db.GetSheetsFromBranch(r.Context(), branchId)
	if err != nil {
		msg := fmt.Sprintf("Could not get sheets from branch id: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	fmt.Println(sheetsDb)

	//todo
	respondWithJSON(w, 200, "")
}
