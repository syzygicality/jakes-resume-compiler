package compiler

type latexRequest struct {
	Source string `json:"source" validate:"required,min=1,max=1000000"`
}
