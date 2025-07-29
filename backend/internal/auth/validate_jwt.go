package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	f := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, f)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Could not parse into token: %v", err)
	}

	id_str, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Could not get the claims subject: %v", err)
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != TokenTypeAccess {
		return uuid.Nil, errors.New("invalid issuer")
	}

	validated_uuid, err := uuid.Parse(id_str)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Could not parse to uuid: %v", err)
	}

	return validated_uuid, nil
}
