package middleware

import (
	"jakes-resume-compiler/server/shared/config"
	"net/http"
)

var publicRoutes = map[string]bool{
	"/health": true,
}

func authMiddleware(app *config.App) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestAPIKey := r.Header.Get(config.APIKeyHeader)

			if requestAPIKey != app.Settings.APIKey && !publicRoutes[r.URL.Path] {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
