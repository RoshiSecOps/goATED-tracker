package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}
	return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	parts := strings.Fields(token)
	if len(parts) != 2 {
		return "", fmt.Errorf("malformed authorization header")
	}
	if !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("authorization scheme must be Bearer")
	}
	return parts[1], nil
}
