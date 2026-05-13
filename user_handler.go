package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RoshiSecOps/goATED-tracker/internal/auth"
	"github.com/RoshiSecOps/goATED-tracker/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

type LogedUser struct {
	User
	Token string `json:"token"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
	if len(params.Password) == 0 {
		respondWithError(w, 400, "password cannot be empty")
		return
	}
	passHash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "unable to hash password")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username: params.Username, Passwordhash: passHash,
	})
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

func (cfg *apiConfig) userLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
	if len(params.Password) == 0 {
		respondWithError(w, 400, "password cannot be empty")
		return
	}
	passHash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "unable to hash password")
		return
	}
	user, err := cfg.db.GetUserByNameAndPass(r.Context(), database.GetUserByNameAndPassParams{
		Username:     params.Username,
		Passwordhash: passHash,
	})

	secret := os.Getenv("JWT_SECRET")
	expiresIn := time.Hour

	jwtToken, err := auth.MakeJWT(
		user.ID,
		secret,
		expiresIn,
	)
	if err != nil {
		respondWithError(w, 500, "unable to create token")
	}
	respondWithJSON(w, 200, LogedUser{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
		},
		Token: jwtToken,
	})
}
