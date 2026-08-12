package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"jakes-resume-compiler/server/config"
	"jakes-resume-compiler/server/middleware"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
}

func main() {
	app := config.NewApp()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", ping)

	var handler http.Handler = middleware.SetupMiddleware(mux, app)

	log.Println("listening on :8080")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}