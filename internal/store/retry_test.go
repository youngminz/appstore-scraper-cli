package store

import (
	"context"
	"errors"
	"testing"
)

func TestRetryTransientRetriesConservativeAttempts(t *testing.T) {
	attempts := 0
	err := retryTransient(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return transientError{err: errors.New("temporary outage")}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("retryTransient() error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestRetryTransientDoesNotRetryPermanentErrors(t *testing.T) {
	attempts := 0
	err := retryTransient(context.Background(), func() error {
		attempts++
		return errors.New("bad request")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}
