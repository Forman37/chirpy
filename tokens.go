package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Forman37/chirpy/internal/auth"
	"github.com/Forman37/chirpy/internal/database"
)

func (c *apiConfig) refreshRefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Refreshing Token")
	// Gather token from header
	token, err := getAccessTokenFromRequest(r)
	if err != nil {
		respondWithError(w, 404, "Unauthorized Access", err)
		return
	}

	log.Println("Checking refresh token")
	// Check validity of refresh token
	refreshToken, err := c.db.CheckRefreshToken(r.Context(), token)
	if err != nil || refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Unauthorized Access", err)
		return
	}

	log.Println("Making JWT")
	// Generate new token
	newToken, err := auth.MakeJWT(refreshToken.UserID, c.jwtSecret)

	log.Println("Updating token")
	// Update User in db
	refreshShape := database.UpdateRefreshTokenParams{
		Token:     refreshToken.Token,
		Token_2:   newToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err = c.db.UpdateRefreshToken(r.Context(), refreshShape)
	if err != nil {
		respondWithError(w, 500, "Error updating token", err)
		return
	}

	log.Println("Token Authenticated")
	responseBody := NewToken{Token: newToken}
	respondWithJSON(w, 200, responseBody)
}

func (c *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Println("Revoking Token")
	token, err := getAccessTokenFromRequest(r)
	if err != nil {
		respondWithError(w, 401, "Unauthorized Access", err)
		return
	}
	c.db.RevokeRefreshToken(r.Context(), token)

	respondWithJSON(w, 204, nil)
}
