package domain

import (
	"fmt"
)

type ValidationError struct {
	Field  string
	Reason error
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{Field: field, Reason: err}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("value of field '%s' is invalid: %s", e.Field, e.Reason)
}

func (e *ValidationError) Unwrap() error {
	return e.Reason
}
