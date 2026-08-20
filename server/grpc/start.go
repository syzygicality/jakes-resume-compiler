package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"jakes-resume-compiler/server/grpc/compiler"
	"jakes-resume-compiler/server/grpc/health"
	"jakes-resume-compiler/server/grpc/platform/interceptor"
	"jakes-resume-compiler/server/shared/config"
)

func Start(app *config.App) {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(interceptor.SetupInterceptors(app))

	health.RegisterServer(srv)

	compiler.RegisterServer(srv)

	log.Println("listening on :8080")

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
