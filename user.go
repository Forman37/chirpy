package main

import (
	"encoding/json"
	"net/http"
)

func (conf *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Set up JSON decoder for expected json response. Request in format
	// {
	//		"email":"user@example.com"
	// }
	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := params.Email
	newUser, err := conf.db.CreateUser(r.Context(), email)
	if err != nil {
		respondWithError(w, 500, "Issue creating user", err)
		return
	}

	responseBody := userResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     email,
	}

	respondWithJSON(w, 201, responseBody)
}

func (conf *apiConfig) resetUsers(w http.ResponseWriter, r *http.Request) {
	if conf.platform != "dev" {
		respondWithError(w, 403, "403 Forbidden", nil)
	}

	err := conf.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error deleting all users", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain charset:utf-8")
	responseBody := response{
		Message: "OK",
	}

	respondWithJSON(w, 200, responseBody)
}
