package domain

import "errors"

// Common domain errors.
var (
	ErrInternal     = errors.New("internal error")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
)
