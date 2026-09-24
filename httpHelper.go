package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func getAccessTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization Header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("No Bearer in Authorization header")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	return token, nil
}

func parseRequestBody(r *http.Request, b any) error {
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(b)
	if err != nil {
		return err
	}
	return nil
}
