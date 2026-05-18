package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RoshiSecOps/goATED-tracker/internal/auth"
	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"teamname"`
}

func (cfg *apiConfig) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Teamname string `json:"teamname"`
	}
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	_, err = auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 401, "unable to validate jwt")
		return
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
	team, err := cfg.db.CreateTeam(r.Context(), params.Teamname)
	if err != nil {
		respondWithError(w, 401, "unable to create team")
		return
	}
	respondWithJSON(w, 201, Team{
		ID:        team.ID,
		CreatedAt: team.CreatedAt,
		UpdatedAt: team.UpdatedAt,
		Name:      team.Teamname,
	})

}

func (cfg *apiConfig) getTeamsHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	_, err = auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 401, "unable to validate jwt")
		return
	}
	teams, err := cfg.db.GetAllTeams(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not fetch teams")
	}
	formattedTeams := []Team{}
	for _, team := range teams {
		formattedTeams = append(formattedTeams, databaseTeamtoTeam(team))
	}
	respondWithJSON(w, 200, formattedTeams)
}

func (cfg *apiConfig) wipeTeamsHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	_, err = auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 401, "unable to validate jwt")
		return
	}
	err = cfg.db.WipeTeams(r.Context())
	if err != nil {
		respondWithError(w, 401, "teams not reset")
	}
	respondWithJSON(w, 204, "")
}

func (cfg *apiConfig) getTeamHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	_, err = auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 401, "unable to validate jwt")
		return
	}

	teamName := r.PathValue("teamName")

	team, err := cfg.db.GetTeamByName(r.Context(), teamName)
	if err != nil {
		respondWithError(w, 404, "Team not found")
		return
	}
	respondWithJSON(w, 200, databaseTeamtoTeam(team))
}
