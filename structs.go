package main

import (
	"sync/atomic"
	"time"

	"github.com/Forman37/chirpy/internal/database"
	"github.com/google/uuid"
)

type param struct {
	Body string `json:"body"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type response struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
