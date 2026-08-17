package compiler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"jakes-resume-compiler/server/http/platform/utils"
)

// Must stay below main.go's WriteTimeout, or the write deadline fires first and
// the 504 never reaches the client.
const compileTimeout = 5 * time.Second

// Preamble precompiled by the Dockerfile's fmt-builder stage. Named, not pathed:
// TEXFORMATS resolves it.
const baseFormat = "base"

func SetupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /compile", compileHandler)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	lr, err := utils.DecodeAndValidate[latexRequest](r)
	if err != nil {
		utils.HTTPError(err, w, r, "dto failure", http.StatusBadRequest)
		return
	}

	source, err := trimToDocument(lr.Source)
	if err != nil {
		utils.HTTPError(err, w, r, "source validation failure", http.StatusBadRequest)
		return
	}

	reqID := utils.GetRequestID(r.Context())

	dir, err := os.MkdirTemp("", "resume-dir")
	if err != nil {
		utils.HTTPError(err, w, r, "temp dir creation failure", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	tex, err := os.CreateTemp(dir, fmt.Sprintf("resume_%s-*.tex", reqID))
	if err != nil {
		utils.HTTPError(err, w, r, ".tex file creation failure", http.StatusInternalServerError)
		return
	}
	texPath := tex.Name()

	if _, err := tex.Write([]byte(source)); err != nil {
		tex.Close()
		utils.HTTPError(err, w, r, ".tex file write failure", http.StatusInternalServerError)
		return
	}
	if err := tex.Close(); err != nil {
		utils.HTTPError(err, w, r, ".tex file close failure", http.StatusInternalServerError)
		return
	}

	pdfPath := strings.TrimSuffix(texPath, ".tex") + ".pdf"

	ctx, cancel := context.WithTimeout(r.Context(), compileTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"pdflatex",
		"-fmt="+baseFormat,
		"-interaction=nonstopmode",
		"-no-shell-escape",
		"-output-directory", dir,
		texPath,
	)

	stdout, runErr := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		utils.HTTPError(ctx.Err(), w, r, "pdflatex compile timeout", http.StatusGatewayTimeout)
		return
	}

	// A source yielding no pages exits 0 with a zero-byte PDF, so existence alone
	// would serve it as a 200.
	info, statErr := os.Stat(pdfPath)
	if statErr != nil || info.Size() == 0 {
		utils.HTTPError(runErr, w, r, "pdflatex compile failure: "+string(stdout), http.StatusUnprocessableEntity)
		return
	}

	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		utils.HTTPError(err, w, r, "pdf read failure", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(pdfBytes); err != nil {
		slog.Error("pdf response write failure", "error", err, "path", r.URL.Path)
	}
}
