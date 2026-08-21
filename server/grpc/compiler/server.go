package compiler

import (
	"context"

	"google.golang.org/grpc"

	compilerpb "jakes-resume-compiler/proto"
	"jakes-resume-compiler/server/grpc/platform/utils"
	"jakes-resume-compiler/server/shared/services"
)

type compilerServer struct {
	compilerpb.UnimplementedCompilerServer
}

func RegisterServer(srv *grpc.Server) {
	compilerpb.RegisterCompilerServer(srv, &compilerServer{})
}

func (s *compilerServer) Compile(ctx context.Context, req *compilerpb.CompileRequest) (*compilerpb.CompileResponse, error) {
	method, _ := grpc.Method(ctx)

	pdfBytes, err := services.Compile(ctx, req.TexSource)
	if err != nil {
		return nil, utils.ServerError(err.Is, method, err.LogMsg, err.RPCCode, err.ReqID)
	}

	return &compilerpb.CompileResponse{Pdf: pdfBytes}, nil
}
