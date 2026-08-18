package compiler

import (
	"errors"
	"strings"
)

type latexRequest struct {
	Source string `json:"source" validate:"required,min=1,max=1000000"`
}

const documentMarker = `\begin{document}`

// trimToDocument drops everything ahead of \begin{document}. The precompiled
// format already holds the preamble, so a submitted one re-runs \documentclass
// and errors with "Two \documentclass or \documentstyle commands".
func trimToDocument(source string) (string, error) {
	i := strings.Index(source, documentMarker)
	if i < 0 {
		return "", errors.New("source is missing " + documentMarker)
	}
	return source[i:], nil
}
