package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"jakes-resume-compiler/server/grpc/platform/utils"
)

func recoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			recErr, ok := rec.(error)
			if !ok {
				recErr = status.Errorf(codes.Internal, "%v", rec)
			}
			err = utils.ServerError(recErr, info.FullMethod, "panic recovered", codes.Internal)
		}
	}()
	return handler(ctx, req)
}
