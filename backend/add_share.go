package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type ShareParams struct {
	Email   string `json:"email"`
	Perm    string `json:"perm"`
	TableId string `json:"table_id"`
}

func (cfg *apiConfig) addShareHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := ShareParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding column: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if _, err := mail.ParseAddress(params.Email); err != nil {
		msg := fmt.Sprintf("Invalid email address: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	tableId, err := uuid.Parse(params.TableId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the table id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !cfg.checkTablePermission(userId, tableId, "write", r.Context()) {
		respondWithError(w, http.StatusForbidden, "Insufficient write permissions")
		return
	}

	sharedTo, err := cfg.db.GetUserByMail(r.Context(), params.Email)
	if err != nil {
		msg := fmt.Sprintf("User with email %s not found: %s", params.Email, err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	if !cfg.canAssignPermision(userId, tableId, params.Perm, r.Context()) {
		msg := fmt.Sprintf("Could not parse the table id: %s", err)
		respondWithError(w, http.StatusForbidden, msg)
		return
	}

	createUserTable := database.CreateUserTableParams{
		UserID:     sharedTo.ID,
		TableID:    tableId,
		Permission: params.Perm,
	}

	getUserTableParams := database.GetUserTablesParams{
		UserID:  sharedTo.ID,
		TableID: tableId,
	}

	userTable, err := cfg.db.GetUserTables(r.Context(), getUserTableParams)
	if err != nil {
		_, err = cfg.db.CreateUserTable(r.Context(), createUserTable)
		if err != nil {
			msg := fmt.Sprintf("User table could not be updated: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	} else {
		updateUserTableParams := database.UpdateUserTableParams{
			UserID:     userTable.UserID,
			TableID:    userTable.TableID,
			Permission: params.Perm,
		}

		err = cfg.db.UpdateUserTable(r.Context(), updateUserTableParams)
		if err != nil {
			msg := fmt.Sprintf("User user table could not be updated: %s", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, "")
}
