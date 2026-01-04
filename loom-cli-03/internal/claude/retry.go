// Package claude provides the Claude CLI integration layer.
//
// This file implements retry logic with exponential backoff.
// Implements: l2/tech-specs.md TS-RETRY-001, TS-RETRY-002
// Implements: DEC-L1-015 (retry configuration: 3 attempts, 2s base, 30s max)
package claude

import (
	"strings"
	"time"
)

// RetryConfig defines retry behavior for Claude API calls.
//
// Implements: TS-RETRY-001
type RetryConfig struct {
	MaxAttempts int           // Number of attempts (not retries). Default: 3
	BaseDelay   time.Duration // Initial delay between attempts. Default: 2s
	MaxDelay    time.Duration // Maximum delay cap. Default: 30s
}

// DefaultRetryConfig returns the default retry configuration.
//
// Implements: DEC-L1-015 (3 attempts, 2s base, 30s max)
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// WithRetry executes a function with retry logic and exponential backoff.
//
// Implements: TS-RETRY-001
// Backoff formula: delay = min(BaseDelay * 2^(attempt-1), MaxDelay)
//
//	Attempt 1: 2s
//	Attempt 2: 4s
//	Attempt 3: 8s (capped at MaxDelay if exceeded)
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
			return err // Non-retryable error, fail immediately
		}

		// Don't sleep after last attempt
		if attempt < cfg.MaxAttempts {
			delay := calculateDelay(cfg, attempt)
			time.Sleep(delay)
		}
	}

	return lastErr
}

// calculateDelay computes the delay for a given attempt using exponential backoff.
//
// Implements: TS-RETRY-001
// Formula: delay = min(BaseDelay * 2^(attempt-1), MaxDelay)
func calculateDelay(cfg RetryConfig, attempt int) time.Duration {
	// 2^(attempt-1)
	multiplier := 1 << (attempt - 1)
	delay := time.Duration(multiplier) * cfg.BaseDelay

	if delay > cfg.MaxDelay {
		return cfg.MaxDelay
	}
	return delay
}

// isRetryableError determines if an error should trigger a retry.
//
// Implements: TS-RETRY-002
//
// Retryable errors (will retry):
//   - rate limit, timeout, connection refused, connection reset
//   - temporary failure, 503, 502, 500, overloaded
//   - service unavailable, gateway timeout
//
// Non-retryable errors (fail immediately):
//   - invalid api key, unauthorized, bad request, 404, not found
//   - forbidden, 401, 403
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Non-retryable patterns - check these first
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
	}

	for _, pattern := range retryable {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	// Default: assume retryable for unknown errors
	// This is safer as transient network issues may not match known patterns
	return true
}
