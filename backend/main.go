package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type apiConfig struct {
}

func main() {
	mux := new(http.ServeMux)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	apiCfg := apiConfig{}

	h := http.HandlerFunc(apiCfg.testColumnHandler)
	mux.Handle("GET /columns", h)

	h = http.HandlerFunc(apiCfg.testSheetsHandler)
	mux.Handle("GET /sheets", h)

	h = http.HandlerFunc(apiCfg.testSaveHandler)
	mux.Handle("POST /save", h)
	mux.Handle("OPTIONS /save", h)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Listening failed: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *apiConfig) testColumnHandler(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	wr.Header().Set("Access-Control-Allow-Origin", "*")
	wr.WriteHeader(200)
	content := `[
{"name": "name", "columnType": "text", "required": true },
{"name": "age", "columnType": "number", "required": false },
{"name": "city", "columnType": "text", "required": false },
{"name": "active", "columnType": "bool", "required": false },
{"name": "salary", "columnType": "number", "required": false },
{"name": "questions", "columnType": "text", "required": true }
]`
	wr.Write([]byte(content))
}

func (cfg *apiConfig) testSheetsHandler(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	wr.Header().Set("Access-Control-Allow-Origin", "*")
	wr.WriteHeader(200)
	content := `["config", "questions"]`
	wr.Write([]byte(content))
}

func (cfg *apiConfig) testSaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println("Received body:", string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Saved"}`))
}
