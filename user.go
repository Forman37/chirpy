package main

import (
	"encoding/json"
	"net/http"
)

func (conf *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Set up JSON decoder for expected json response. Request in format
	// {
	//		"email": "user@example.com"
	// }
	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := params.Body
	newUser, err := conf.db.CreateUser(r.Context(), email)
	if err != nil {
		respondWithError(w, 500, "Issue creating user", err)
		return
	}

	responseBody := userResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	respondWithJSON(w, 201, responseBody)
	return
}

func (conf *apiConfig) resetUsers(w http.ResponseWriter, r *http.Request) {
	return
}
