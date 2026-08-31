package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	grpcserver "jakes-resume-compiler/server/grpc"
	httpserver "jakes-resume-compiler/server/http"
	"jakes-resume-compiler/server/shared/config"
)

// router dispatches each request to the gRPC server or the REST handler.
//
// gRPC always speaks HTTP/2 and always tags its bodies application/grpc (plus a
// codec suffix such as +proto), so that pair identifies a gRPC call with no
// ambiguity against the REST routes, which are JSON over either protocol
// version. Routing happens before the HTTP middleware chain, so gRPC calls run
// under their own interceptors instead of both stacks.
func router(handler http.Handler, grpcSrv *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcSrv.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func Start(app *config.App) {
	grpcSrv := grpcserver.Server(app)

	handler := httpserver.Handler(app)

	log.Println("listening on :8080 (HTTP and gRPC)")

	// Unencrypted HTTP/2 is required for gRPC over this listener; HTTP/1.1
	// stays enabled so REST clients that do not negotiate h2c still work.
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router(handler, grpcSrv),
		Protocols:    protocols,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
