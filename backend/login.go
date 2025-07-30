package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

const acc_expire_time = time.Hour
const ref_expire_time = time.Hour * 24 * 60

type Login struct {
	ID          uuid.UUID `json:"id"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	OpenedTable Table     `json:"opened_table"`
}

func (cfg *apiConfig) login_handler(w http.ResponseWriter, req *http.Request) {
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

	user, err := cfg.db.GetUserByMail(req.Context(), params.Email)
	if err != nil {
		msg := fmt.Sprintf("User with email %s not found: %s", params.Email, err)
		respondWithError(w, 500, msg)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		msg := fmt.Sprintf("Password check failed: %s", err)
		respondWithError(w, 401, msg)
		return
	}

	cfg.ReturnLoginData(w, user, req.Context(), 200)
}

func (cfg *apiConfig) ReturnLoginData(w http.ResponseWriter, user database.User, ctx context.Context, code int) {
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
	_, err = cfg.db.CreateRefreshToken(ctx, ref_params)
	if err != nil {
		msg := fmt.Sprintf("Problem with creating refresh token: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.RefreshTokenName,
		Value:    ref_token,
		Path:     "/",
		Expires:  time.Now().Add(ref_expire_time),
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   cfg.platform != PlatformDev,
		SameSite: http.SameSiteStrictMode,
	})

	table, err := cfg.GetTable(user.ID, user.OpenedTable, ctx)
	if err != nil {
		msg := fmt.Sprintf("With getting an opened table: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	ret := Login{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		Token:      token,
	}
	respondWithJSON(w, code, ret)
}
