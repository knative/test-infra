package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrTimeoutExceeded = errors.New("timeout exceeded")

func WithDeadline(ctx context.Context, fn func() error) error {
	deadline, ok := ctx.Deadline()
	if ok {
		if time.Now().After(deadline) {
			return fmt.Errorf("%w: %v", ErrTimeoutExceeded, deadline)
		}
	}

	return fn()
}
