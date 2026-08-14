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

	"jakes-resume-compiler/server/platform/utils"
)

const compileTimeout = 30 * time.Second

func SetupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /compile", compileHandler)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	lr, err := utils.DecodeAndValidate[latexRequest](r)
	if err != nil {
		utils.HTTPError(err, w, r, "dto failure", http.StatusBadRequest)
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

	if _, err := tex.Write([]byte(lr.Source)); err != nil {
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

	if _, statErr := os.Stat(pdfPath); statErr != nil {
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
