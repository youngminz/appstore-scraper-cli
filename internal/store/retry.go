package store

import (
	"context"
	"errors"
	"strings"
	"time"
)

type transientError struct {
	err error
}

func (e transientError) Error() string {
	return e.err.Error()
}

func (e transientError) Unwrap() error {
	return e.err
}

func retryTransient(ctx context.Context, fn func() error) error {
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		if err := fn(); err != nil {
			last = err
			var transient transientError
			if !errors.As(err, &transient) || ctx.Err() != nil {
				return err
			}
			timer := time.NewTimer(time.Duration(attempt+1) * 250 * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				continue
			}
		}
		return nil
	}
	return last
}

func isTransientStatus(status int) bool {
	return status == 408 || status == 429 || status >= 500
}

func isTransientText(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "timeout") ||
		strings.Contains(text, "temporary") ||
		strings.Contains(text, "connection reset") ||
		strings.Contains(text, "status code: 408") ||
		strings.Contains(text, "status code: 429") ||
		strings.Contains(text, "status code: 5")
}
