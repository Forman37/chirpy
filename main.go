package main

import (
	"fmt"
	"net/http"
	"log"
)

type welcomeScreen struct{}

func (welcomeScreen) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fmt.Println("Starting server on http://localhost:8080 ...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

}

