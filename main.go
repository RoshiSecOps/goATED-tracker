package main

import "net/http"

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
