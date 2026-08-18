package utils

import (
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ServerError(err error, method string, logMsg string, code codes.Code) error {
	if code == codes.Internal || code == codes.Unknown {
		slog.Error(logMsg,
			"error", err,
			"stack", string(debug.Stack()),
			"method", method,
		)
	} else {
		slog.Warn(logMsg,
			"error", err,
			"method", method,
		)
	}
	return status.Error(code, code.String())
}
