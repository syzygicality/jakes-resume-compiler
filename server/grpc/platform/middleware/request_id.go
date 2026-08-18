package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"jakes-resume-compiler/server/config"
)

func requestIDInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	reqID := ""

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(config.RequestIDHeader); len(vals) > 0 {
			reqID = vals[0]
		}
	}

	if reqID == "" {
		reqID = config.NewRequestID()
	}

	ctx = context.WithValue(ctx, config.RequestIDKey, reqID)

	grpc.SetHeader(ctx, metadata.Pairs(config.RequestIDHeader, reqID))

	return handler(ctx, req)
}
