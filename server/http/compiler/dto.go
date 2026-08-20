package compiler

import (
	"errors"
	"jakes-resume-compiler/server/config"
	"strings"
)

type latexRequest struct {
	Source string `json:"source" validate:"required,min=1,max=1000000"`
}

// trimToDocument drops everything ahead of \begin{document}. The precompiled
// format already holds the preamble, so a submitted one re-runs \documentclass
// and errors with "Two \documentclass or \documentstyle commands".
func trimToDocument(source string) (string, error) {
	i := strings.Index(source, config.DocumentMarker)
	if i < 0 {
		return "", errors.New("source is missing " + config.DocumentMarker)
	}
	return source[i:], nil
}
