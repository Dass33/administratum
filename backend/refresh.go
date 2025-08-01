package main

import (
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/auth"
)

func (cfg *apiConfig) refresh_handler(w http.ResponseWriter, req *http.Request) {
	token_req, err := auth.GetRerfreshCookie(req)
	if err != nil {
		msg := fmt.Sprintf("Could not get the token from cookie: %v", err)
		respondWithError(w, 400, msg)
		return
	}
	token_db, err := cfg.db.GetRefreshToken(req.Context(), token_req)
	if err != nil {
		msg := fmt.Sprintf("Could not find the token in database: %v", err)
		respondWithError(w, 401, msg)
		return
	}
	if token_db.RevokedAt.Valid {
		msg := fmt.Sprintf("Refresh token had been revoked: %v", err)
		respondWithError(w, 401, msg)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), token_db.UserID)
	if err != nil {
		msg := fmt.Sprintf("User with id not found: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	cfg.ReturnLoginData(w, user, req.Context(), 200)
}
