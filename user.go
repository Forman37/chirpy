package main

import (
	"encoding/json"
	"net/http"

	"github.com/Forman37/chirpy/internal/auth"
	"github.com/Forman37/chirpy/internal/database"
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
	password := params.Password

	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password", err)
		return
	}

	newParams := database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPass,
	}

	newUser, err := conf.db.CreateUser(r.Context(), newParams)
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

func (c *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := params.Email
	password := params.Password

	user, err := c.db.FetchUser(r.Context(), email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password1", err)
		return
	}

	correctPass, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil || correctPass == false {
		respondWithError(w, 401, "Incorrect email or password2", err)
		return
	}

	responseBody := userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 200, responseBody)
}
