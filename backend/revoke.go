package main

import (
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/auth"
)

func (cfg *apiConfig) revoke_handler(w http.ResponseWriter, req *http.Request) {
	token_req, err := auth.GetRerfreshCookie(req)
	if err != nil {
		msg := fmt.Sprintf("Could not get the token from cookie: %v", err)
		respondWithError(w, 400, msg)
		return
	}
	err = cfg.db.RevokeRefreshToken(req.Context(), token_req)
	if err != nil {
		msg := fmt.Sprintf("Could not revoke the token: %v", err)
		respondWithError(w, 500, msg)
		return
	}
	respondWithJSON(w, 204, "")
}
