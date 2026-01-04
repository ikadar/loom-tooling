// Package claude provides the Claude CLI integration layer.
//
// This file implements retry logic with exponential backoff.
//
// Implements: l2/tech-specs.md TS-RETRY-001, TS-RETRY-002
// Implements: DEC-L1-015 (retry configuration)
package claude

import (
	"strings"
	"time"
)

// RetryConfig configures retry behavior for Claude API calls.
//
// Implements: l2/tech-specs.md TS-RETRY-001
type RetryConfig struct {
	MaxAttempts int           // Number of attempts (not retries)
	BaseDelay   time.Duration // Initial delay between attempts
	MaxDelay    time.Duration // Maximum delay cap
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
// Implements: l2/tech-specs.md TS-RETRY-001
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

	return lastErr
}

// calculateDelay computes the delay for a given attempt using exponential backoff.
//
// Implements: l2/tech-specs.md TS-RETRY-001
// Formula: delay = min(BaseDelay * 2^(attempt-1), MaxDelay)
func calculateDelay(cfg RetryConfig, attempt int) time.Duration {
	// 2^(attempt-1)
	multiplier := 1 << (attempt - 1) // bit shift for power of 2

	delay := cfg.BaseDelay * time.Duration(multiplier)

	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}

	return delay
}

// isRetryableError determines if an error should trigger a retry.
//
// Implements: l2/tech-specs.md TS-RETRY-002
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Retryable patterns
	retryable := []string{
		"rate limit",
		"rate_limit",
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
		"gateway timeout",
		"overloaded",
		"503",
		"502",
		"500",
		"429", // Too Many Requests
	}

	for _, pattern := range retryable {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	// Non-retryable patterns - fail immediately
	nonRetryable := []string{
		"invalid api key",
		"invalid_api_key",
		"unauthorized",
		"bad request",
		"bad_request",
		"not found",
		"not_found",
		"forbidden",
		"401",
		"403",
		"404",
	}

	for _, pattern := range nonRetryable {
		if strings.Contains(errStr, pattern) {
			return false
		}
	}

	// Default: retry unknown errors
	return true
}
