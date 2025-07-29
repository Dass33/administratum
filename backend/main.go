package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/Dass33/administratum/backend/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	apiCfg := apiConfig{}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiCfg.DB = dbQueries
	log.Println("Connected to database!")

	router := chi.NewRouter()

	// todo remove http in prod
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// router.Post("/users", apiCfg.handlerUsersCreate)
	// router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
	router.Get("/columns", apiCfg.testColumnHandler)
	router.Get("/sheets", apiCfg.testSheetsHandler)
	router.Get("/save", apiCfg.testSaveHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 0,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
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

	log.Println("Received body:", string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Saved"}`))
}
