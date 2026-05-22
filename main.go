package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/RoshiSecOps/goATED-tracker/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to init database")
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{db: dbQueries}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("DELETE /api/users", apiCfg.wipeUsersHandler)
	mux.HandleFunc("POST /api/login", apiCfg.userLoginHandler)
	mux.HandleFunc("POST /api/admin/teams", apiCfg.createTeamHandler)
	mux.HandleFunc("GET /api/admin/teams", apiCfg.getTeamsHandler)
	mux.HandleFunc("DELETE /api/admin/teams", apiCfg.wipeTeamsHandler)
	mux.HandleFunc("GET /api/admin/teams/{TeamName}", apiCfg.getTeamHandler)
	mux.HandleFunc("POST /api/admin/teams/members", apiCfg.addMemberHandler)
	mux.HandleFunc("GET /api/admin/teams/members", apiCfg.getTeamsMembersHandler)
	mux.HandleFunc("DELETE /api/admin/teams/members", apiCfg.wipeTeamsMembersHandler)
	mux.HandleFunc("POST /api/admin/teams/pentests", apiCfg.addPentestHandler)
	mux.HandleFunc("GET /api/admin/teams/pentests", apiCfg.getPentestsHandler)
	mux.HandleFunc("DELETE /api/admin/teams/pentests", apiCfg.wipePentestsHandler)
	mux.HandleFunc("POST /api/admin/findings", apiCfg.addFindingHandler)
	mux.HandleFunc("GET /api/admin/findings", apiCfg.getFindingsHandler)
	mux.HandleFunc("DELETE /api/admin/findings", apiCfg.wipeFindingsHandler)
	mux.HandleFunc("GET /api/users/teams", apiCfg.getUserTeamsHandler) // Get Team membership as a user
	mux.HandleFunc("GET /api/v1/{TeamName}/pentests", apiCfg.getTeamPentestsHandler)
	mux.HandleFunc("POST /api/v1/{TeamName}/pentests", apiCfg.addTeamPentestsHandler)
	mux.HandleFunc("POST /api/v1/{TeamName}/pentests/findings", apiCfg.addFindingsUserHandler)
	mux.HandleFunc("GET /api/v1/{TeamName}/{PentestTitle}/findings", apiCfg.getPentestFindingsUserHandler)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
