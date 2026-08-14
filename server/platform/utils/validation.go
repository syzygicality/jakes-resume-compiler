package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate[T any](r *http.Request) (*T, error) {
	var body T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&body); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if err := validate.Struct(body); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &body, nil
}
