package middleware

import (
	"net/http"
	"jakes-resume-compiler/server/config"
)

var publicRoutes = map[string]bool{
	"/health": true,
}

func authMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestAPIKey := r.Header.Get("X-API-Key")

			if requestAPIKey != app.Settings.APIKey && !publicRoutes[r.URL.Path] {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}