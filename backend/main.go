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
	rawDB    *sql.DB
	platform string
	jwt_key  string
}

type IdName struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
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
		rawDB:    db,
		platform: os.Getenv("PLATFORM"),
		jwt_key:  os.Getenv("JWT_KEY"),
	}

	log.Println("Connected to database!")

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://dass33.github.io",
			"http://localhost:5173",
		},
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
	router.Post("/create_sheet", apiCfg.middlewareAuth(apiCfg.createSheetHandler))
	router.Get("/json/{branch_id}", apiCfg.getJsonHandler)
	router.Put("/rename_sheet", apiCfg.middlewareAuth(apiCfg.renameSheetHandler))
	router.Delete("/delete_sheet", apiCfg.middlewareAuth(apiCfg.deleteSheetHandler))
	router.Put("/rename_project", apiCfg.middlewareAuth(apiCfg.renemeProjectHandler))
	router.Delete("/delete_project", apiCfg.middlewareAuth(apiCfg.deleteProjectHandler))
	router.Post("/add_share", apiCfg.middlewareAuth(apiCfg.addShareHandler))
	router.Delete("/delete_row", apiCfg.middlewareAuth(apiCfg.deleteRowHandler))
	router.Put("/change_game_url", apiCfg.middlewareAuth(apiCfg.changeGameUrlHandler))
	router.Post("/create_branch", apiCfg.middlewareAuth(apiCfg.createBranchHandler))
	router.Get("/get_branch/{branch_id}", apiCfg.middlewareAuth(apiCfg.getBranchHandler))
	router.Delete("/delete_branch", apiCfg.middlewareAuth(apiCfg.deleteBranchHandler))
	router.Put("/update_branch", apiCfg.middlewareAuth(apiCfg.updateBranchHandler))
	router.Post("/merge_preview", apiCfg.middlewareAuth(apiCfg.mergePreviewHandler))
	router.Post("/merge_execute", apiCfg.middlewareAuth(apiCfg.mergeExecuteHandler))
	router.Get("/merge_targets", apiCfg.middlewareAuth(apiCfg.getMergeTargetsHandler))

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 0,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
