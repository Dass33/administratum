package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth_heder := headers.Get("Authorization")
	if auth_heder == "" {
		return "", errors.New("Authorization header not found")
	}
	token_string := strings.TrimPrefix(auth_heder, "Bearer")
	return strings.TrimSpace(token_string), nil
}
