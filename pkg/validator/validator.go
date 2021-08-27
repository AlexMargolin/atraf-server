package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator is an alias for validator.Validate
type Validator = validator.Validate

// NewValidator returns a new validator instance
func NewValidator() *Validator {
	return validator.New()
}
