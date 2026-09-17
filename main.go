package main

import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type healthHandler struct {}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte("OK")
	w.Write(body)
}

func (cfg *apiConfig) checkViews(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	val := cfg.fileserverHits.Load()

	fmt.Fprintf(w, "Hits: %d", val)
}

func (cfg *apiConfig) resetViews(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	val := cfg.fileserverHits.Load()

	fmt.Fprintf(w, "Hits: %d", val)
}

func main(){
	var apiCfg apiConfig
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("."))

	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app/", fileserver),
		),
	)

	mux.Handle("GET /healthz", healthHandler{})
	mux.HandleFunc("GET /metrics", apiCfg.checkViews)
	mux.HandleFunc("POST /reset", apiCfg.resetViews)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fmt.Println("Starting server on http://localhost:8080 ...")

	// Blocks main until server is closed
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

}

