package jwt

import (
	"fmt"
	"os"
	"time"

	"auth-service/internal/auth/models"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Username string
	ID       int32
	Role     models.Role
}

func GenerateJWT(user models.User, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   user.Username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenStr string) (*jwt.Token, *CustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if token.Valid {
		return token, claims, nil
	}

	return nil, nil, fmt.Errorf("invalid token")
}
