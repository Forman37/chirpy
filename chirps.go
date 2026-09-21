package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Forman37/chirpy/internal/database"
)

func (c *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpBody := database.CreateChirpParams{
		UserID: params.UserID,
		Body:   strings.Join(words, " "),
	}

	newChirp, err := c.db.CreateChirp(r.Context(), chirpBody)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	responseBody := chirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	respondWithJSON(w, 201, responseBody)
}
