package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/Dass33/administratum/backend/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const PlatformDev = "dev"
const PlatformProd = "production"

type apiConfig struct {
	db       *database.Queries
	platform string
	jwt_key  string
}

type IdName struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

const OwnerPermission string = "owner"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		db:       dbQueries,
		platform: os.Getenv("PLATFORM"),
		jwt_key:  os.Getenv("JWT_KEY"),
	}

	log.Println("Connected to database!")

	router := chi.NewRouter()

	// todo change origin in prod
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/login", apiCfg.loginHandler)
	router.Post("/register", apiCfg.createUserHandler)
	router.Post("/refresh", apiCfg.refreshHandler)
	router.Post("/logout", apiCfg.revokeHandler)
	router.Put("/update_column", apiCfg.middlewareAuth(apiCfg.updateColumnHandler))
	router.Post("/add_column", apiCfg.middlewareAuth(apiCfg.addColumnHandler))
	router.Put("/update_column_data", apiCfg.middlewareAuth(apiCfg.updateColumnDataHandler))
	router.Post("/add_column_data", apiCfg.middlewareAuth(apiCfg.addColumnDataHandler))
	router.Delete("/delete_column", apiCfg.middlewareAuth(apiCfg.deleteColumnHandler))
	router.Get("/get_sheet/{sheet_id}", apiCfg.middlewareAuth(apiCfg.getSheetHandler))
	router.Get("/get_project/{table_id}", apiCfg.middlewareAuth(apiCfg.getProjectHandler))
	router.Post("/create_project", apiCfg.middlewareAuth(apiCfg.createProjectHandler))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 0,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
