package validate

import (
	"github.com/go-playground/validator/v10"
)

type Validate = validator.Validate

// NewValidator provides a wrapper for the validator package.
func NewValidator() *Validate {
	return validator.New()
}
