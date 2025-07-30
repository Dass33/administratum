package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

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

	cfg.ReturnLoginData(w, user, req.Context(), 201)
}
