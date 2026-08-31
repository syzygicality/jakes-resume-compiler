package http

import (
	"net/http"

	"jakes-resume-compiler/server/http/compiler"
	"jakes-resume-compiler/server/http/health"
	"jakes-resume-compiler/server/http/platform/middleware"
	"jakes-resume-compiler/server/shared/config"
)

// Handler builds the REST mux and wraps it in the HTTP middleware chain. The
// caller owns the listener, since the gRPC server shares it.
func Handler(app *config.App) http.Handler {
	mux := http.NewServeMux()

	health.SetupHandlers(mux)

	compiler.SetupHandlers(mux)

	return middleware.SetupMiddleware(mux, app)
}
