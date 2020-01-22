package ensure

import "github.com/pkg/errors"

// NoError checks is there is no error given
func NoError(err error) {
	if err != nil {
		panic(errors.WithMessage(err, "expected to be no error, but that was"))
	}
}
