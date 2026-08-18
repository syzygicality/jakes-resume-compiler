package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"jakes-resume-compiler/server/config"
	"jakes-resume-compiler/server/grpc/health"
)

func Start(app *config.App) {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	health.RegisterServer(srv)

	// TODO: register compiler service (generated from .proto)

	// TODO: interceptor chain (recovery, request ID, logging, auth), mirroring server/http/platform/middleware

	log.Println("listening on :8080")

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
