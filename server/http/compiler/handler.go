package compiler

import (
	"log/slog"
	"net/http"

	"jakes-resume-compiler/server/http/platform/utils"
	"jakes-resume-compiler/server/shared/services"
)

type latexRequest struct {
	Source string `json:"source" validate:"required,min=1,max=1000000"`
}

func SetupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /compile", compileHandler)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	lr, valErr := utils.DecodeAndValidate[latexRequest](r)
	reqID := services.GetRequestID(r.Context())

	if valErr != nil {
		utils.HTTPError(valErr, w, r, "dto failure", http.StatusBadRequest, reqID)
		return
	}

	pdfBytes, err := services.Compile(r.Context(), lr.Source)
	if err != nil {
		utils.HTTPError(err.Is, w, r, err.LogMsg, err.HTTPCode, err.ReqID)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(pdfBytes); err != nil {
		slog.Error("pdf response write failure", "error", err, "path", r.URL.Path)
	}
}
