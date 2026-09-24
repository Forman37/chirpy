package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Forman37/chirpy/internal/auth"
	"github.com/Forman37/chirpy/internal/database"
)

func (conf *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating user")
	w.Header().Set("Content-Type", "application/json")
	// Set up JSON decoder for expected json response. Request in format
	// {
	//		"email":"user@example.com"
	// }
	params := param{}
	err := parseRequestBody(r, &params)
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
	log.Println("Resetting User")
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
	log.Println("Logging in")
	w.Header().Set("Content-Type", "application/json")
	params := param{}
	err := parseRequestBody(r, &params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	email := params.Email
	password := params.Password
	expiresIn := params.ExpiresIn

	if expiresIn == 0 {
		expiresIn = 3600
	}
	if expiresIn > 3600 {
		expiresIn = 3600
	}

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

	token, err := auth.MakeJWT(user.ID, c.jwtSecret)
	if err != nil {
		respondWithError(w, 500, "Error making token", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()

	newRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	newToken, err := c.db.CreateRefreshToken(r.Context(), newRefreshTokenParams)
	if err != nil {
		respondWithError(w, 500, "Error creating new refresh token", err)
		return
	}

	responseBody := userResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: newToken.Token,
	}

	respondWithJSON(w, 200, responseBody)
}

func (c *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	token, err := getAccessTokenFromRequest(r)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	user, err := auth.ValidateJWT(token, c.jwtSecret)

}
