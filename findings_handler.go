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

type Finding struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	Severity      string    `json:"severity"`
	SeverityScore int       `json:"severity_score"`
	File          string    `json:"file"`
	AtLine        int       `json:"at_line"`
	Description   string    `json:"description"`
	PentestId     uuid.UUID `json:"pentest_id"`
}

func (cfg *apiConfig) addFindingHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Status        string    `json:"status"`
		Severity      string    `json:"severity"`
		SeverityScore int       `json:"severity_score"`
		File          string    `json:"file"`
		AtLine        int       `json:"at_line"`
		Description   string    `json:"description"`
		PentestId     uuid.UUID `json:"pentest_id"`
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
	finding, err := cfg.db.AddFindingToPentest(r.Context(),
		database.AddFindingToPentestParams{
			Title:         params.Title,
			Status:        params.Status,
			Severity:      params.Severity,
			SeverityScore: int32(params.SeverityScore),
			File:          params.File,
			AtLine:        int32(params.AtLine),
			Description:   params.Description,
			PentestID:     params.PentestId,
		})
	if err != nil {
		respondWithError(w, 500, "unable to create finding")
		log.Printf("Error: %v", err)
		return
	}
	respondWithJSON(w, 201, finding)
}

func (cfg *apiConfig) getFindingsHandler(w http.ResponseWriter, r *http.Request) {
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
	findings, err := cfg.db.GetAllFindings(r.Context())
	if err != nil {
		respondWithError(w, 500, "unable to retrieve findings")
		return
	}
	formattedFindings := []Finding{}
	for _, finding := range findings {
		formattedFindings = append(formattedFindings, databaseFindingtoFinding(finding))
	}
	respondWithJSON(w, 200, formattedFindings)
}

func (cfg *apiConfig) wipeFindingsHandler(w http.ResponseWriter, r *http.Request) {
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
	err = cfg.db.WipeFindings(r.Context())
	if err != nil {
		respondWithError(w, 500, "unable to reset findings data")
		return
	}
	respondWithJSON(w, 204, "")
}

func (cfg *apiConfig) addFindingsUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Status        string    `json:"status"`
		Severity      string    `json:"severity"`
		SeverityScore int       `json:"severity_score"`
		File          string    `json:"file"`
		AtLine        int       `json:"at_line"`
		Description   string    `json:"description"`
		PentestId     uuid.UUID `json:"pentest_id"`
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
	test, err = cfg.db.CheckPentestAccess(r.Context(), database.CheckPentestAccessParams{
		ID:     params.PentestId,
		TeamID: team.ID,
	})
	if err != nil {
		respondWithError(w, 500, "unable to validate pentest to team mapping")
		log.Printf("Error: %v", err)
		return
	}
	if !test {
		respondWithError(w, 401, "pentest report is not owned by the provided team")
		log.Printf("Error: teamID %v, expected %v", team.ID, params.PentestId)
		return
	}
	finding, err := cfg.db.AddFindingToPentest(r.Context(),
		database.AddFindingToPentestParams{
			Title:         params.Title,
			Status:        params.Status,
			Severity:      params.Severity,
			SeverityScore: int32(params.SeverityScore),
			File:          params.File,
			AtLine:        int32(params.AtLine),
			Description:   params.Description,
			PentestID:     params.PentestId,
		})
	if err != nil {
		respondWithError(w, 500, "unable to create finding")
		log.Printf("Error: %v", err)
		return
	}
	respondWithJSON(w, 201, finding)
}

func (cfg *apiConfig) getPentestFindingsUserHandler(w http.ResponseWriter, r *http.Request) {
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
	pentestTitle := r.PathValue("PentestTitle")
	pentest, err := cfg.db.GetPentestByTitle(r.Context(), pentestTitle)
	if err != nil {
		respondWithError(w, 401, "unable to retrieve pentest")
		return
	}
	test, err = cfg.db.CheckPentestAccess(r.Context(), database.CheckPentestAccessParams{
		ID:     pentest.ID,
		TeamID: team.ID,
	})
	if err != nil {
		respondWithError(w, 500, "unable to validate pentest to team mapping")
		log.Printf("Error: %v", err)
		return
	}
	if !test {
		respondWithError(w, 401, "pentest report is not owned by the provided team")
		log.Printf("Error: teamID %v, expected %v", team.ID, pentest.TeamID)
		return
	}
	dbFindings, err := cfg.db.GetFindingsForPentest(r.Context(), pentest.ID)
	if err != nil {
		respondWithError(w, 500, "unable to retrieve findings")
		return
	}
	formattedFindings := []Finding{}
	for _, finding := range dbFindings {
		formattedFindings = append(formattedFindings, databaseFindingtoFinding(finding))
	}
	respondWithJSON(w, 200, formattedFindings)
}
