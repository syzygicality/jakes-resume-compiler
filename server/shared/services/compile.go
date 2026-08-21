package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
)

type CompileError struct {
	Is       error
	LogMsg   string
	RPCCode  codes.Code
	HTTPCode int
	ReqID    string
}

// Is is nil when the failure has no underlying error: pdflatex can exit 0 on a
// zero-byte PDF, which fails the compile with nothing to report but LogMsg.
func (e *CompileError) Error() string {
	if e.Is == nil {
		return e.LogMsg
	}
	return e.Is.Error()
}

const PreambleFormat = "preamble"

const DocumentMarker = `\begin{document}`

const CompileTimeout = 5 * time.Second

func trimToDocument(source string) (string, error) {
	i := strings.Index(source, DocumentMarker)
	if i < 0 {
		return "", errors.New("source is missing " + DocumentMarker)
	}
	return source[i:], nil
}

func Compile(ctx context.Context, reqSource string) ([]byte, *CompileError) {
	reqID := GetRequestID(ctx)

	source, err := trimToDocument(reqSource)
	if err != nil {
		return nil, &CompileError{
			Is:       err,
			LogMsg:   "source validation failure",
			RPCCode:  codes.InvalidArgument,
			HTTPCode: http.StatusBadRequest,
			ReqID:    reqID,
		}
	}

	dir, err := os.MkdirTemp("", "resume-dir")
	if err != nil {
		return nil, &CompileError{
			Is:       err,
			LogMsg:   "temp dir creation failure",
			RPCCode:  codes.Internal,
			HTTPCode: http.StatusInternalServerError,
			ReqID:    reqID,
		}
	}
	defer os.RemoveAll(dir)

	tex, err := os.CreateTemp(dir, fmt.Sprintf("resume_%s-*.tex", reqID))
	if err != nil {
		return nil, &CompileError{
			Is:       err,
			LogMsg:   ".tex file creation failure",
			RPCCode:  codes.Internal,
			HTTPCode: http.StatusInternalServerError,
			ReqID:    reqID,
		}
	}
	texPath := tex.Name()

	if _, err := tex.Write([]byte(source)); err != nil {
		tex.Close()
		return nil, &CompileError{
			Is:       err,
			LogMsg:   ".tex file write failure",
			RPCCode:  codes.Internal,
			HTTPCode: http.StatusInternalServerError,
			ReqID:    reqID,
		}
	}

	if err := tex.Close(); err != nil {
		return nil, &CompileError{
			Is:       err,
			LogMsg:   ".tex file close failure",
			RPCCode:  codes.Internal,
			HTTPCode: http.StatusInternalServerError,
			ReqID:    reqID,
		}
	}

	pdfPath := strings.TrimSuffix(texPath, ".tex") + ".pdf"

	ctx, cancel := context.WithTimeout(ctx, CompileTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"pdflatex",
		"-fmt="+PreambleFormat,
		"-interaction=nonstopmode",
		"-no-shell-escape",
		"-output-directory", dir,
		texPath,
	)

	stdout, runErr := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		return nil, &CompileError{
			Is:       ctx.Err(),
			LogMsg:   "pdflatex compile timeout",
			RPCCode:  codes.DeadlineExceeded,
			HTTPCode: http.StatusGatewayTimeout,
			ReqID:    reqID,
		}
	}

	// A source yielding no pages exits 0 with a zero-byte PDF, so existence alone
	// would serve it as a success.
	info, statErr := os.Stat(pdfPath)
	if statErr != nil || info.Size() == 0 {
		return nil, &CompileError{
			Is:       runErr,
			LogMsg:   "pdflatex compile failure: " + string(stdout),
			RPCCode:  codes.InvalidArgument,
			HTTPCode: http.StatusUnprocessableEntity,
			ReqID:    reqID,
		}
	}

	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, &CompileError{
			Is:       err,
			LogMsg:   "pdf read failure",
			RPCCode:  codes.Internal,
			HTTPCode: http.StatusInternalServerError,
			ReqID:    reqID,
		}
	}

	return pdfBytes, nil
}
