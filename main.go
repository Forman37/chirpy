package main

import (
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main(){
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{}, 
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

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	// Blocks main until server is closed
	log.Fatal(server.ListenAndServe())
}

