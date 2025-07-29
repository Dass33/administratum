package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type Users struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

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

	ret := Users{
		ID:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	respondWithJSON(w, 201, ret)
}
