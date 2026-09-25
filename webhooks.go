package main

import (
	"net/http"

	"github.com/Forman37/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) updateUserToChirpyRed(w http.ResponseWriter, r *http.Request) {

	// Structs for parsing header
	type userIDStruct struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type userInputBody struct {
		Event string       `json:"event"`
		Data  userIDStruct `json:"data"`
	}

	// Gather and Verify API Key
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != c.polkaKey {
		respondWithError(w, 401, "Unauthorized Access", err)
		return
	}

	userInput := userInputBody{}

	err = parseRequestBody(r, &userInput)
	if err != nil {
		respondWithError(w, 500, "Error parsing body", err)
		return
	}

	if userInput.Event != "user.upgraded" {
		respondWithError(w, 204, "Inapproprate event", nil)
		return
	}

	err = c.db.UpdateToChirpyRed(r.Context(), userInput.Data.UserID)
	if err != nil {
		respondWithError(w, 401, "Error with updating to chirpy red", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
