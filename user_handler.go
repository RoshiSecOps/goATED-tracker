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

func (cfg *apiConfig) wipeUsersHandler(w http.ResponseWriter, r *http.Request) {
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
	err = cfg.db.WipeUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not delete users")
		log.Printf("error wiping users: %v", err)
		return
	}
	respondWithJSON(w, 200, "users successfully reset")
}

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := cfg.db.GetUserByName(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, 400, "unable to get user")
	}
	match, err := auth.CheckPasswordHash(params.Password, user.Passwordhash)
	if err != nil {
		respondWithError(w, 400, "wrong credentials")
		return
	}
	if match != true {
		respondWithError(w, 400, "wrong credentials")
		return
	}
	secret := os.Getenv("JWT_SECRET")
	expiresIn := time.Hour

	jwtToken, err := auth.MakeJWT(
		user.ID,
		secret,
		expiresIn,
	)
	if err != nil {
		respondWithError(w, 500, "unable to create token")
		return
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

func (cfg *apiConfig) getUserTeamsHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "unable to get auth token")
		return
	}
	userId, err := auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 400, "unable to validate jwt")
		return
	}
	teams, err := cfg.db.GetTeamsForUser(r.Context(), userId)
	formattedTeams := []Team{}
	for _, team := range teams {
		formattedTeams = append(formattedTeams, databaseTeamtoTeam(team))
	}
	respondWithJSON(w, 200, formattedTeams)
}

func (cfg *apiConfig) getTeamPentestsHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "unable to get auth token")
		return
	}
	userId, err := auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 400, "unable to validate jwt")
		return
	}
	teamName := r.PathValue("TeamName")
	team, err := cfg.db.GetTeamByName(r.Context(), teamName)
	if err != nil {
		respondWithError(w, 500, "unable to retrieve team")
		log.Printf("Error: %v", err)
		return
	}
	test, err := cfg.db.CheckMembership(r.Context(), database.CheckMembershipParams{
		UserID: userId,
		TeamID: team.ID,
	})
	if err != nil {
		respondWithError(w, 500, "unable to check membership")
		log.Printf("Error: %v", err)
		return
	}
	if !test {
		respondWithError(w, 401, "not member of the team")
		return
	}
	pentests, err := cfg.db.GetPentestsForTeam(r.Context(), team.ID)
	if err != nil {
		respondWithError(w, 500, "unable to fetch pentests")
		log.Printf("Error: %v", err)
		return
	}
	formattedTests := []Pentest{}
	for _, pentest := range pentests {
		formattedTests = append(formattedTests, databasePentesttoPentest(pentest))
	}
	respondWithJSON(w, 200, formattedTests)
}

func (cfg *apiConfig) addTeamPentestsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title  string    `json:"title"`
		TeamID uuid.UUID `json:"team_id"`
	}
	secret := os.Getenv("JWT_SECRET")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 500, "unable to get auth token")
		return
	}
	userId, err := auth.ValidateJWT(token, secret)
	if err != nil {
		respondWithError(w, 400, "unable to validate jwt")
		return
	}
	teamName := r.PathValue("TeamName")
	team, err := cfg.db.GetTeamByName(r.Context(), teamName)
	if err != nil {
		respondWithError(w, 500, "unable to retrieve team")
		log.Printf("Error: %v", err)
		return
	}
	test, err := cfg.db.CheckMembership(r.Context(), database.CheckMembershipParams{
		UserID: userId,
		TeamID: team.ID,
	})
	if err != nil {
		respondWithError(w, 500, "unable to check membership")
		log.Printf("Error: %v", err)
		return
	}
	if !test {
		respondWithError(w, 401, "not member of the team")
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
	pentest, err := cfg.db.CreatePentest(r.Context(), database.CreatePentestParams{
		Title:  params.Title,
		TeamID: params.TeamID,
	})
	if err != nil {
		respondWithError(w, 500, "unable to create pentest")
		return
	}
	respondWithJSON(w, 201, Pentest{
		ID:        pentest.ID,
		CreatedAt: pentest.CreatedAt,
		UpdatedAt: pentest.UpdatedAt,
		Title:     pentest.Title,
		TeamID:    pentest.TeamID,
	})
}
