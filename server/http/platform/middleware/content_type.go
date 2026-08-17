package middleware

import (
	"net/http"
	"strings"
)

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodDelete, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		ct := r.Header.Get("Content-Type")

		if !strings.HasPrefix(ct, "application/json") {
			http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
