package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"

	"jakes-resume-compiler/server/platform/config"
	"jakes-resume-compiler/server/platform/middleware"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}
	app := config.NewApp()
	mux := http.NewServeMux()

	middleware.SetupLogger(app.Settings.Prod)

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
