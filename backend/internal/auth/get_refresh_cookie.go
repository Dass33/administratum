package auth

import (
	"errors"
	"net/http"
)

const RefreshTokenName = "refresh_token"

func GetRerfreshCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("No refresh token found")
		}
		return "", errors.New("Error reading cookie")
	}

	refreshToken := cookie.Value
	return refreshToken, nil
}
