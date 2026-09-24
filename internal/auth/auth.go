package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("Returning from hashing password with error")
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Println("Returning from pass compare with error")
		return false, err
	}

	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expiresIn := time.Hour * 1

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	byteSecret := []byte(tokenSecret)
	ss, err := token.SignedString(byteSecret)
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	log.Println("Validating JWT")
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

	if err != nil {
		log.Println("Returning Here at Parse of claims")
		return uuid.Nil, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Returning Here at GetSubject")
		return uuid.Nil, err
	}

	returnID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Returning Here at Parse of UUID")
		return uuid.Nil, err
	}
	log.Println("Returning Complete")
	return returnID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	//log.Printf("Headers for getting bearer token : %v\n", headers)
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No authorization header found")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("No bearer token found in header")
	}

	token := strings.TrimPrefix(authHeader, prefix)

	return token, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	randStr := hex.EncodeToString(key)

	return randStr
}
