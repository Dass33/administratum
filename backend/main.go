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

	h := http.HandlerFunc(apiCfg.request_count_handler)
	mux.Handle("GET /admin", h)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Listening failed: %v\n", err)
		os.Exit(1)
	}
}

func (cfg *apiConfig) request_count_handler(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "text/json")
	wr.Header().Set("Access-Control-Allow-Origin", "*")
	wr.WriteHeader(200)
	content := `{"text": "hello from go"}`
	wr.Write([]byte(content))
}
