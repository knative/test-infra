package gowork

import (
	"errors"
	"fmt"
)

var (
	// ErrBug is returned when the error is probably a bug.
	ErrBug = errors.New("probably a bug")

	// ErrInvalidGowork is returned when the go.work is invalid.
	ErrInvalidGowork = errors.New("invalid go.work")

	// ErrInvalidGomod is returned when the go.mod is invalid.
	ErrInvalidGomod = errors.New("invalid go.mod")
)

// errWrap wraps an error as parent error.
func errWrap(err error, as error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, as) {
		return err
	}
	return fmt.Errorf("%w: %v", as, err)
}
