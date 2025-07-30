package main

import (
	"fmt"
	"net/http"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, uuid.UUID)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			msg := fmt.Sprintf("Couldn't find bearer token: %s", err)
			respondWithError(w, http.StatusUnauthorized, msg)
			return
		}

		user_id, err := auth.ValidateJWT(token, cfg.jwt_key)
		if err != nil {
			msg := fmt.Sprintf("Couldn't get user id: %s", err)
			respondWithError(w, http.StatusNotFound, msg)
			return
		}

		handler(w, r, user_id)
	}
}
