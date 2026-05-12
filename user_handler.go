package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
	}
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "could not read request body")
		return
	}
	params := parameters{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 500, "unable to unmarshal data")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 500, "unable to create user")
		log.Printf("Error: %v", err)
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
	})
}

func (cfg *apiConfig) wipeUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.WipeUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not delete users")
		log.Printf("error wiping users: %v", err)
		return
	}
	respondWithJSON(w, 200, "users successfully reset")
}
