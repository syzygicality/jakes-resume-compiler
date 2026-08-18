package middleware

import (
	"jakes-resume-compiler/server/config"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request",
			"request_id", config.GetRequestID(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

func SetupLogger(prod bool) {
	var logger slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if prod {
		logger = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(logger))
}
