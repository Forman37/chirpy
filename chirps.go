package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Forman37/chirpy/internal/auth"
	"github.com/Forman37/chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("Posting Chirp")
	naughtyWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

	w.Header().Set("Content-Type", "application/json")

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Header", err)
		return
	}

	userID, err := auth.ValidateJWT(bearer, c.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := param{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong decoding params", err)
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
		UserID: userID,
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

func (c *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting chirps")
	allChirps, err := c.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong getting chirps", err)
		return
	}

	responseBody := []chirpResponse{}

	for _, chirp := range allChirps {
		newResponse := chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		responseBody = append(responseBody, newResponse)
	}

	respondWithJSON(w, 200, responseBody)
}

func (c *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Chirp")
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, "Error with chirp ID: ", err)
		return
	}

	newChirp, err := c.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error fetching chirp", err)
		return
	}

	responseBody := chirpResponse{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	respondWithJSON(w, 200, responseBody)
}

func (c *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting Chirp")
	token, err := getAccessTokenFromRequest(r)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	user, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Token", err)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 401, "Invalid ID", err)
		return
	}

	// Verify user is owner of chirp
	chirp, err := c.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Error fetching chirp", err)
		return
	}
	if chirp.UserID != user {
		respondWithError(w, 403, "Unauthorized to Delete chirp", nil)
		return
	}
	// END : Verify user is owner of chirp

	chirpParams := database.DeleteChirpParams{
		ID:     chirpID,
		UserID: user,
	}

	err = c.db.DeleteChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 403, "Unauthorized to delete this Chirp", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
