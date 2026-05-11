package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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
