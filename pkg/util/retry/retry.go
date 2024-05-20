// Package retry provides simple retry wrappers for functions that return an error
package retry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	// DefaultTimeout is a default timeout for retry operations
	DefaultTimeout = 2 * time.Minute
	// Interval is the time to wait between retry attempts
	Interval = 5 * time.Second
	// ErrAbort should be returned when an error occurs on which retrying should be aborted
	ErrAbort = errors.New("retrying aborted")
)

// Context is a retry wrapper that will retry the given function until it succeeds or the context is cancelled
func Context(ctx context.Context, f func(ctx context.Context) error) error {
	var lastErr error

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Execute the function immediately for the first try
	lastErr = f(ctx)
	if lastErr == nil || errors.Is(lastErr, ErrAbort) {
		return lastErr
	}

	ticker := time.NewTicker(Interval)
	defer ticker.Stop()

	attempt := 0

	for {
		select {
		case <-ctx.Done():
			slog.Debug(fmt.Sprintf("retry.Context: context cancelled after %d attempts", attempt))
			return errors.Join(ctx.Err(), lastErr)
		case <-ticker.C:
			attempt++
			if lastErr != nil {
				slog.DebugContext(ctx, fmt.Sprintf("retrying, attempt %d - last error: %v", attempt, lastErr))
			}
			lastErr = f(ctx)

			if errors.Is(lastErr, ErrAbort) {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Context: aborted after %d attempts", attempt))
				return lastErr
			}

			if lastErr == nil {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Context: succeeded after %d attempts", attempt))
				return nil
			} else {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Context: attempt %d failed: %s", attempt, lastErr))
			}
		}
	}
}

// Timeout is a retry wrapper that will retry the given function until it succeeds, the context
// is cancelled, or the timeout is reached
func Timeout(ctx context.Context, timeout time.Duration, f func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return Context(ctx, f)
}

// Times is a retry wrapper that will retry the given function until it succeeds or the given number of
// attempts have been made
func Times(ctx context.Context, times int, f func(context.Context) error) error {
	var lastErr error

	// Execute the function immediately for the first try
	lastErr = f(ctx)
	if lastErr == nil || errors.Is(lastErr, ErrAbort) {
		return lastErr
	}

	i := 1

	ticker := time.NewTicker(Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Debug(fmt.Sprintf("retry.Times: context cancelled after %d attempts", i))
			return errors.Join(ctx.Err(), lastErr)
		case <-ticker.C:
			if lastErr != nil {
				slog.DebugContext(ctx, fmt.Sprintf("retrying: attempt %d of %d (previous error: %v)", i+1, times, lastErr))
			}

			lastErr = f(ctx)

			if errors.Is(lastErr, ErrAbort) {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Times: aborted after %d attempts", i))
				return lastErr
			}

			if lastErr == nil {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Times: succeeded on attempt %d", i))
				return nil
			}

			i++

			if i >= times {
				slog.DebugContext(ctx, fmt.Sprintf("retry.Times: exceeded %d attempts", times))
				return fmt.Errorf("retry limit exceeded after %d attempts: %w", times, lastErr)
			}
		}
	}
}
