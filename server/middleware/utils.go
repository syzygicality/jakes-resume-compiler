package middleware

import (
	"jakes-resume-compiler/server/config"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func SetupMiddleware(mux http.Handler, app *config.App) http.Handler {
	return chain(mux,
		authMiddleware(app),
	)
}
