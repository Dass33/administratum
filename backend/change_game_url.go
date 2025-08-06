package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Dass33/administratum/backend/internal/database"
	"github.com/google/uuid"
)

type changeGameUrlParams struct {
	GameUrl sql.NullString `json:"game_url"`
	TableId string         `json:"table_id"`
}

func (cfg *apiConfig) changeGameUrlHandler(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	decoder := json.NewDecoder(r.Body)
	params := changeGameUrlParams{}

	err := decoder.Decode(&params)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	if !isValidURL(params.GameUrl.String) && params.GameUrl.Valid {
		msg := "Given url is not valid"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	projectId, err := uuid.Parse(params.TableId)
	if err != nil {
		msg := fmt.Sprintf("Could not parse the project id: %s", err)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	changeGameUrlParams := database.ChangeGameUrlParams{
		GameUrl: params.GameUrl,
		ID:      projectId,
	}
	err = cfg.db.ChangeGameUrl(r.Context(), changeGameUrlParams)
	if err != nil {
		msg := fmt.Sprintf("Game url could not be changed: %s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
