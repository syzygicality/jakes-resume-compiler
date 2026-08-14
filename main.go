package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"jakes-resume-compiler/server/compiler"
	"jakes-resume-compiler/server/health"
	"jakes-resume-compiler/server/platform/config"
	"jakes-resume-compiler/server/platform/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}
	app := config.NewApp()
	mux := http.NewServeMux()

	health.SetupHandlers(mux)

	compiler.SetupHandlers(mux)

	middleware.SetupLogger(app.Settings.Prod)

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
