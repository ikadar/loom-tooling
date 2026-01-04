// Package claude provides the Claude CLI integration layer.
//
// This file implements retry logic with exponential backoff.
// Implements: l2/tech-specs.md TS-RETRY-001, TS-RETRY-002
// Implements: DEC-L1-015 (retry configuration)
package claude

import (
	"fmt"
	"strings"
	"time"
)

// RetryConfig configures retry behavior.
// Implements: l2/tech-specs.md TS-RETRY-001
type RetryConfig struct {
	MaxAttempts int           // Default: 3
	BaseDelay   time.Duration // Default: 2 seconds
	MaxDelay    time.Duration // Default: 30 seconds
}

// DefaultRetryConfig returns the default retry configuration.
// Implements: DEC-L1-015 (3 attempts, 2s base, 30s max)
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// WithRetry executes the given function with retry logic.
// Implements: l2/tech-specs.md TS-RETRY-001
//
// Backoff Formula:
// delay = min(BaseDelay * 2^(attempt-1), MaxDelay)
//
// Attempt 1: 2s
// Attempt 2: 4s
// Attempt 3: 8s (capped at 30s if exceeded)
func WithRetry(cfg RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}

		// Don't sleep after last attempt
		if attempt < cfg.MaxAttempts {
			delay := calculateDelay(cfg, attempt)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// calculateDelay computes the delay for the given attempt.
// Formula: min(BaseDelay * 2^(attempt-1), MaxDelay)
func calculateDelay(cfg RetryConfig, attempt int) time.Duration {
	delay := cfg.BaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
	}
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	return delay
}

// isRetryableError determines if an error should trigger a retry.
// Implements: l2/tech-specs.md TS-RETRY-002 (Error Classification)
//
// Retryable: rate limit, timeout, 503, 502, 500, overloaded, connection refused,
//
//	connection reset, temporary failure
//
// Non-retryable: invalid api key, unauthorized, bad request, 404
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Retryable patterns
	retryable := []string{
		"rate limit",
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"503",
		"502",
		"500",
		"overloaded",
		"service unavailable",
		"gateway timeout",
		"bad gateway",
	}

	for _, pattern := range retryable {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	// Non-retryable patterns - return false immediately
	nonRetryable := []string{
		"invalid api key",
		"unauthorized",
		"bad request",
		"404",
		"not found",
		"forbidden",
		"401",
		"403",
	}

	for _, pattern := range nonRetryable {
		if strings.Contains(errStr, pattern) {
			return false
		}
	}

	// Default: retry unknown errors
	return true
}
