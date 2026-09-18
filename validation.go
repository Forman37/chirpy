package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)

type param struct {
	Body string `json:"body"`
}


type validation struct {}

func (v validation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	naughtyWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

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
	
	words := strings.Fields(params.Body)

	for i, word := range words {
		fmt.Println("Comparing : " + word)
		if _, ok := naughtyWords[strings.ToLower(word)]; ok {
			fmt.Println("true")
			words[i] = "****"
		}
	}

	fmt.Println(words)
	responseBody := response {
		Valid: true,
		CleanedBody: strings.Join(words, " "),
	}
	respondWithJSON(w, 200, responseBody)
	return
} 
