package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", fmt.Errorf("Password has to be under 72 bytes")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("There was a problem with hashing the password")
	}
	return string(hash), nil
}
