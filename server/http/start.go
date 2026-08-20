package http

import (
	"log"
	"net/http"
	"time"

	"jakes-resume-compiler/server/http/compiler"
	"jakes-resume-compiler/server/http/health"
	"jakes-resume-compiler/server/http/platform/middleware"
	"jakes-resume-compiler/server/shared/config"
)

func Start(app *config.App) {
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
