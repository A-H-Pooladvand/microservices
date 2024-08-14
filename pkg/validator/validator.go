package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate validates provided `i`. It is usually called after `Context#Bind()`.
// Validator must be registered using `Echo#Validator`.
func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}
