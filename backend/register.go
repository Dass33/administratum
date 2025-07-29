package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/Dass33/administratum/backend/internal/database"
)

func (cfg *apiConfig) create_user_handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	_, err = mail.ParseAddress(params.Email)
	if err != nil {
		msg := fmt.Sprintf("Invalid email address: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		msg := fmt.Sprintf("Error password hashing failed: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	user_par := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	user, err := cfg.db.CreateUser(req.Context(), user_par)
	if err != nil {
		msg := fmt.Sprintf("Error creating user: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwt_key, acc_expire_time)
	if err != nil {
		msg := fmt.Sprintf("Problem with creating access token: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	ref_token, _ := auth.MakeRefreshToken()
	ref_params := database.CreateRefreshTokenParams{
		Token:     ref_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(ref_expire_time),
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), ref_params)
	if err != nil {
		msg := fmt.Sprintf("Problem with creating refresh token: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	ret := Login{
		ID:           user.ID,
		Created_at:   user.CreatedAt,
		Updated_at:   user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: ref_token,
	}

	respondWithJSON(w, 201, ret)
}
