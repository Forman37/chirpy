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
	// Load env vars
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("JWTSECRET")
	pKey := os.Getenv("POLKA_KEY")

	// START: DB initiation and opening
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
	// END : DB initiation and opening

	const filepathRoot = "."
	const port = "8080"

	// Config to pass to functions
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      secret,
		polkaKey:       pKey,
	}

	// Server MUX and fileserver initialization
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))

	// START : Handlers and Endpoints
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/", fileserver),
		),
	)

	mux.Handle("GET /api/healthz", healthHandler{})
	//mux.HandleFunc("GET /admin/metrics", apiCfg.checkViews)
	mux.Handle("POST /api/validate_chirp", validation{})
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/login", apiCfg.login)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.updateUserToChirpyRed)
	// END : Handlers and Endpoints

	// Setup server with mux as the handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// Blocks main until server is closed
	log.Fatal(server.ListenAndServe())
}
