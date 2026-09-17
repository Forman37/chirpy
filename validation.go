package main

import (
	"encoding/json"
	"net/http"
)

type param struct {
	Body string `json:"body"`
}


type validation struct {}

func (v validation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	
	responseBody := response {
		Valid: true,
	}
	respondWithJSON(w, 200, responseBody)
	return
} 
