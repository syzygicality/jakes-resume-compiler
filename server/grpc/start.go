package grpc

import (
	"google.golang.org/grpc"

	"jakes-resume-compiler/server/grpc/compiler"
	"jakes-resume-compiler/server/grpc/health"
	"jakes-resume-compiler/server/grpc/platform/interceptor"
	"jakes-resume-compiler/server/shared/config"
)

// Server builds the gRPC server with its services registered. The caller owns
// the listener: this is served through ServeHTTP off the shared HTTP/2 port
// rather than Serve(lis), so no listener is bound here.
func Server(app *config.App) *grpc.Server {
	srv := grpc.NewServer(interceptor.SetupInterceptors(app))

	health.RegisterServer(srv)

	compiler.RegisterServer(srv)

	return srv
}
