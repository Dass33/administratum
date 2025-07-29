package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type failed struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	err_msg := failed{
		Error: msg,
	}
	dat, _ := json.Marshal(err_msg)
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		msg := fmt.Sprintf("Error encoding parameters: %s", err)
		respondWithError(w, 500, msg)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
