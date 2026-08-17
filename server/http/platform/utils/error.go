package utils

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func HTTPError(err error, w http.ResponseWriter, r *http.Request, logMsg string, status int) {
	if status >= 500 {
		slog.Error(logMsg,
			"error", err,
			"stack", string(debug.Stack()),
			"path", r.URL.Path,
			"method", r.Method,
		)
	} else {
		slog.Warn(logMsg,
			"error", err,
			"path", r.URL.Path,
			"method", r.Method,
		)
	}
	http.Error(w, http.StatusText(status), status)
}
