package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type validation struct{}

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
		if _, ok := naughtyWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	responseBody := response{
		Valid:       true,
		CleanedBody: strings.Join(words, " "),
	}
	respondWithJSON(w, 200, responseBody)
	return
}
