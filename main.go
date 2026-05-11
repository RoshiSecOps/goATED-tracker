package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/RoshiSecOps/goATED-tracker/internal/database"
	_ "github.com/lib/pq"
)

func main() {
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
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}

type apiConfig struct {
	db *database.Queries
}
