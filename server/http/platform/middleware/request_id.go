package middleware

import (
	"context"
	"net/http"

	"jakes-resume-compiler/server/config"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(config.RequestIDHeader)

		if reqID == "" {
			reqID = config.NewRequestID()
		}

		ctx := context.WithValue(r.Context(), config.RequestIDKey, reqID)

		w.Header().Set(config.RequestIDHeader, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
