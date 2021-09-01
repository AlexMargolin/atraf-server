package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator = validator.Validate

// NewValidator provides a wrapper for the validator package.
func NewValidator() *Validator {
	return validator.New()
}
