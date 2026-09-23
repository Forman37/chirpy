package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}

	if !match {
		t.Errorf("expected password to match hash")
	}
}

func TestWrongPassword(t *testing.T) {
	password := "super-secret-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}

	if match {
		t.Errorf("expected wrong password to not match hash")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, 5*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if returnedID != userID {
		t.Errorf("expected user ID %v, got %v", userID, returnedID)
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	token, err := MakeJWT(userID, "correct-secret", 5*time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")

	if err == nil {
		t.Errorf("expected validation to fail with wrong secret")
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, -1*time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	_, err = ValidateJWT(token, secret)

	if err == nil {
		t.Errorf("expected expired token validation to fail")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer blackjack")

	resp, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}

	if resp != "blackjack" {
		t.Errorf("Token did not match expected value")
	}
}

func TestGetWrongBearerToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer blackjack")

	resp, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}

	if resp == "wrong_token" {
		t.Errorf("Token did not match expected value")
	}
}

func TestGetBearerTokenNoHeader(t *testing.T) {
	header := http.Header{}

	resp, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("Expected to have an error")
	}

	if resp != "" {
		t.Errorf("Should have returned an error and no response")
	}
}
