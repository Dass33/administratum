package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess string = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	curr_time := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    TokenTypeAccess,
		IssuedAt:  jwt.NewNumericDate(curr_time),
		ExpiresAt: jwt.NewNumericDate(curr_time.Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed_strng, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Creation of signed string failed: %v", err)
	}

	return signed_strng, nil
}
