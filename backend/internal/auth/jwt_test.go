package auth_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Dass33/administratum/backend/internal/auth"
	"github.com/google/uuid"
)

func TestTokenGeneration(t *testing.T) {
	for range 5 {
		new_uuid := uuid.New()
		secret := "top secret string"
		expires, _ := time.ParseDuration("10s")

		jwt_string, err := auth.MakeJWT(new_uuid, secret, expires)
		if err != nil {
			t.Fatal(err)
		}

		uuid, err := auth.ValidateJWT(jwt_string, secret)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(uuid)
	}
}
