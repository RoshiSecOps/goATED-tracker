package main

import (
	"encoding/json"
	"net/http"

	"github.com/RoshiSecOps/goATED-tracker/internal/database"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func databaseTeamtoTeam(dbTeam database.Team) Team {
	return Team{
		ID:        dbTeam.ID,
		CreatedAt: dbTeam.CreatedAt,
		UpdatedAt: dbTeam.UpdatedAt,
		Name:      dbTeam.Teamname,
	}
}

func databaseTeamMembersToTeamMembers(dbTeamMembers database.TeamMember) TeamMember {
	return TeamMember{
		ID:        dbTeamMembers.ID,
		CreatedAt: dbTeamMembers.CreatedAt,
		UpdatedAt: dbTeamMembers.UpdatedAt,
		UserID:    dbTeamMembers.UserID,
		TeamID:    dbTeamMembers.TeamID,
	}
}

func databasePentesttoPentest(dbPentest database.Pentest) Pentest {
	return Pentest{
		ID:        dbPentest.ID,
		CreatedAt: dbPentest.CreatedAt,
		UpdatedAt: dbPentest.CreatedAt,
		Title:     dbPentest.Title,
		TeamID:    dbPentest.TeamID,
	}
}

func databaseFindingtoFinding(dbFinding database.Finding) Finding {
	return Finding{
		ID:            dbFinding.ID,
		CreatedAt:     dbFinding.CreatedAt,
		UpdatedAt:     dbFinding.UpdatedAt,
		Title:         dbFinding.Title,
		Status:        dbFinding.Status,
		Severity:      dbFinding.Severity,
		SeverityScore: int(dbFinding.SeverityScore),
		File:          dbFinding.File,
		AtLine:        int(dbFinding.AtLine),
		Description:   dbFinding.Description,
		PentestId:     dbFinding.PentestID,
	}
}
