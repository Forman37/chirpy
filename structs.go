package main

import (
	"sync/atomic"
	"time"

	"github.com/Forman37/chirpy/internal/database"
	"github.com/google/uuid"
)

type param struct {
	Body      string    `json:"body"`
	Email     string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
	Password  string    `json:"password"`
	ExpiresIn int       `json:"expires_in_seconds"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

type response struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
	Message     string `json:"message"`
}

type userResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChripyRed  bool      `json:"is_chirpy_red"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type NewToken struct {
	Token string `json:"token"`
}
