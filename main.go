package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Forman37/chirpy/internal/database"
)

func main() {
	// START: DB initiation and opening
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	log.Printf("Successfully connected to database!")

	dbQueries := database.New(db)
	log.Printf("%s\n", dbQueries)
	// END : DB initiation and opening

	// Load PLATFORM into variable
	platform := os.Getenv("PLATFORM")

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/", fileserver),
		),
	)

	mux.Handle("GET /api/healthz", healthHandler{})
	mux.HandleFunc("GET /admin/metrics", apiCfg.checkViews)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetViews)
	mux.Handle("POST /api/validate_chirp", validation{})
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetUsers)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// Blocks main until server is closed
	log.Fatal(server.ListenAndServe())
}
