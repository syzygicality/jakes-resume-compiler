package compiler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	compilerpb "jakes-resume-compiler/proto"
	"jakes-resume-compiler/server/shared/config"
	"jakes-resume-compiler/server/shared/services"
)

type compilerServer struct {
	compilerpb.UnimplementedCompilerServer
}

func RegisterServer(srv *grpc.Server) {
	compilerpb.RegisterCompilerServer(srv, &compilerServer{})
}

func trimToDocument(source string) (string, error) {
	i := strings.Index(source, config.DocumentMarker)
	if i < 0 {
		return "", errors.New("source is missing " + config.DocumentMarker)
	}
	return source[i:], nil
}

func (s *compilerServer) Compile(ctx context.Context, req *compilerpb.CompileRequest) (*compilerpb.CompileResponse, error) {
	source, err := trimToDocument(req.TexSource)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "source validation failure: "+err.Error())
	}

	reqID := services.GetRequestID(ctx)

	dir, err := os.MkdirTemp("", "resume-dir")
	if err != nil {
		slog.Error("temp dir creation failure", "error", err, "request_id", reqID)
		return nil, status.Error(codes.Internal, "temp dir creation failure")
	}
	defer os.RemoveAll(dir)

	tex, err := os.CreateTemp(dir, fmt.Sprintf("resume_%s-*.tex", reqID))
	if err != nil {
		slog.Error(".tex file creation failure", "error", err, "request_id", reqID)
		return nil, status.Error(codes.Internal, ".tex file creation failure")
	}
	texPath := tex.Name()

	if _, err := tex.Write([]byte(source)); err != nil {
		tex.Close()
		slog.Error(".tex file write failure", "error", err, "request_id", reqID)
		return nil, status.Error(codes.Internal, ".tex file write failure")
	}
	if err := tex.Close(); err != nil {
		slog.Error(".tex file close failure", "error", err, "request_id", reqID)
		return nil, status.Error(codes.Internal, ".tex file close failure")
	}

	pdfPath := strings.TrimSuffix(texPath, ".tex") + ".pdf"

	ctx, cancel := context.WithTimeout(ctx, config.CompileTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"pdflatex",
		"-fmt="+config.PreambleFormat,
		"-interaction=nonstopmode",
		"-no-shell-escape",
		"-output-directory", dir,
		texPath,
	)

	stdout, runErr := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "pdflatex compile timeout")
	}

	// A source yielding no pages exits 0 with a zero-byte PDF, so existence alone
	// would serve it as a success.
	info, statErr := os.Stat(pdfPath)
	if statErr != nil || info.Size() == 0 {
		slog.Warn("pdflatex compile failure", "error", runErr, "stdout", string(stdout), "request_id", reqID)
		return nil, status.Error(codes.InvalidArgument, "pdflatex compile failure: "+string(stdout))
	}

	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		slog.Error("pdf read failure", "error", err, "request_id", reqID)
		return nil, status.Error(codes.Internal, "pdf read failure")
	}

	return &compilerpb.CompileResponse{Pdf: pdfBytes}, nil
}
