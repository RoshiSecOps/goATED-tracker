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
	mux.HandleFunc("DELETE /api/users", apiCfg.wipeUsers)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
