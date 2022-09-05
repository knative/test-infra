package modules

import "errors"

var (
	// ErrBug is returned when the error is probably a bug.
	ErrBug = errors.New("probably a bug")

	// ErrInvalidGowork is returned when the go.work is invalid.
	ErrInvalidGowork = errors.New("invalid go.work")

	// ErrInvalidGomod is returned when the go.mod is invalid.
	ErrInvalidGomod = errors.New("invalid go.mod")
)
