package jwt

import (
	"errors"
	"testing"
	"time"

	"auth-service/internal/auth/models"

	"github.com/golang-jwt/jwt/v5"
)

// Checking if token ie generated correctly and is Valid with signature
func TestGenerateAndValidate(t *testing.T) {
	subject := models.User{
		Username: "John Wick",
		Password: "ilovedogs123",
		Email:    "johnwick43@pancilover.com",
		ID:       12345,
		Role:     1,
	}

	token, err := GenerateJWT(subject, time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate token")
	}

	parse, _, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate token")
	}

	if !parse.Valid {
		t.Fatalf("Token no valid")
	}
}

// Checking if token is expaired aftet ttl ran out
func TestValidateExpiredToken(t *testing.T) {
	subject := models.User{
		Username: "John Wick",
		Password: "ilovedogs123",
		Email:    "johnwick43@pancilover.com",
		ID:       12345,
		Role:     1,
	}

	// Using -time.Minutes to make an expired token
	token, err := GenerateJWT(subject, -time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate token for expiration test: %v", err)
	}

	_, _, err = ValidateJWT(token)
	if err == nil {
		t.Fatalf("Expected an error for expired token, but got nil")
	}

	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("Expected token expired error, but got a different error: %v", err)
	}
}
