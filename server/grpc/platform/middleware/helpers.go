package middleware

import (
	"google.golang.org/grpc"

	"jakes-resume-compiler/server/config"
)

func SetupInterceptors(app *config.App) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		recoveryInterceptor,
		requestIDInterceptor,
		loggingInterceptor,
		authInterceptor(app),
	)
}
