package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"jakes-resume-compiler/server/config"
)

func loggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	st, _ := status.FromError(err)

	slog.Info("request",
		"request_id", config.GetRequestID(ctx),
		"method", info.FullMethod,
		"duration", time.Since(start),
		"code", st.Code(),
		"error", err,
	)

	return resp, err
}
