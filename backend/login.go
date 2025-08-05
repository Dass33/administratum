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

type LoginData struct {
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	OpenedTable  TableData `json:"opened_table"`
	OpenedSheet  Sheet     `json:"opened_sheet"`
	TableIdNames []IdName  `json:"table_names"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
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
		respondWithError(w, 400, msg)
		return
	}

	if !auth.IsValidEmail(params.Email) {
		msg := fmt.Sprintf("Provided email is not valid: %s", err)
		respondWithError(w, 400, msg)
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

	table := TableData{}

	sheet, err := cfg.GetSheet(user.OpenedSheet, ctx)
	if err == nil {
		DbTable, err := cfg.db.GetTableFromSheet(ctx, sheet.ID)
		if err != nil {
			msg := fmt.Sprintf("Problem with getting table from sheet id", err)
			respondWithError(w, 500, msg)
			return
		}
		optional_table_id := uuid.NullUUID{
			UUID:  DbTable.ID,
			Valid: true,
		}
		table, err = cfg.GetTable(user.ID, optional_table_id, ctx)
		if err != nil {
			msg := fmt.Sprintf("Problem with getting table data: %s", err)
			respondWithError(w, 500, msg)
			return
		}
	}

	tables, _ := cfg.db.GetTablesFromUser(ctx, user.ID)
	tableIdNames := make([]IdName, 0, len(tables))

	for i := range tables {
		item := IdName{
			ID:   tables[i].ID,
			Name: tables[i].Name,
		}
		tableIdNames = append(tableIdNames, item)
	}

	ret := LoginData{
		Email:        user.Email,
		Token:        token,
		OpenedTable:  table,
		OpenedSheet:  sheet,
		TableIdNames: tableIdNames,
	}
	respondWithJSON(w, code, ret)
}
