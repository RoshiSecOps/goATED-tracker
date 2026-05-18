package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RoshiSecOps/goATED-tracker/internal/auth"
	"github.com/RoshiSecOps/goATED-tracker/internal/database"
	"github.com/google/uuid"
)

type TeamMember struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	TeamID    uuid.UUID `json:"team_id"`
}

func (cfg *apiConfig) addMember(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Teamname string `json:"teamname"`
		Username string `json:"username"`
	}
	secret := os.Getenv("ADMIN_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	if token != secret {
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
	team, err := cfg.db.GetTeamByName(r.Context(), params.Teamname)
	if err != nil {
		respondWithError(w, 404, "Team not found")
		return
	}
	user, err := cfg.db.GetUserByName(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}
	teamMember, err := cfg.db.CreateTeamMember(r.Context(), database.CreateTeamMemberParams{
		UserID: user.ID,
		TeamID: team.ID,
	})
	if err != nil {
		respondWithError(w, 400, "Could not add member")
		return
	}
	respondWithJSON(w, 201, TeamMember{
		ID:        teamMember.ID,
		CreatedAt: teamMember.CreatedAt,
		UpdatedAt: teamMember.UpdatedAt,
		UserID:    teamMember.UserID,
		TeamID:    teamMember.TeamID,
	})
}

func (cfg *apiConfig) getTeamsMembersHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("ADMIN_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	if token != secret {
		respondWithError(w, 401, "unable to validate jwt")
		return
	}
	teamsMembers, err := cfg.db.GetAllTeamMember(r.Context())
	if err != nil {
		respondWithError(w, 404, "cannot get all teams and members")
		return
	}
	listTeamsMembers := []TeamMember{}
	for _, teamMem := range teamsMembers {
		listTeamsMembers = append(listTeamsMembers, databaseTeamMembersToTeamMembers(teamMem))
	}
	respondWithJSON(w, 200, listTeamsMembers)
}

func (cfg *apiConfig) wipeTeamsMembersHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("ADMIN_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "wrong/malformed auth header")
		return
	}
	if token != secret {
		respondWithError(w, 401, "unable to validate jwt")
		return
	}
	err = cfg.db.WipeTeamMember(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not wipe teams members table")
		return
	}
	respondWithJSON(w, 204, "")
}
