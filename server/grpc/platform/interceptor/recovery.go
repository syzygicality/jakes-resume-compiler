package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"jakes-resume-compiler/server/grpc/platform/utils"
)

func recoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			recErr, ok := rec.(error)
			if !ok {
				recErr = fmt.Errorf("%v", rec)
			}
			err = utils.ServerError(recErr, info.FullMethod, "panic recovered", codes.Internal, "not available")
		}
	}()
	return handler(ctx, req)
}
