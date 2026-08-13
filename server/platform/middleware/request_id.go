package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "requestID"
const requestIDHeader = "X-Request-ID"

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDHeader)

		if reqID == "" {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err != nil {
				reqID = "buy-a-lottery-ticket-bro"
			} else {
				reqID = hex.EncodeToString(bytes)
			}
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)

		w.Header().Set(requestIDHeader, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
