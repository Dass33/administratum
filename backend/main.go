package main

import (
	"fmt"
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
	mux.Handle("GET /admin", h)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Listening failed: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *apiConfig) testColumnHandler(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "text/json")
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
