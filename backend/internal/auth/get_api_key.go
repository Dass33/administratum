package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth_heder := headers.Get("Authorization")
	if auth_heder == "" {
		return "", errors.New("Authorization header not found")
	}
	token_string := strings.TrimPrefix(auth_heder, "ApiKey")
	return strings.TrimSpace(token_string), nil
}
