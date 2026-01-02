package claude

import (
	"fmt"
	"strings"
	"time"
)

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns sensible defaults for retrying
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// isRetryableError checks if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()

	// Retryable conditions
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
	}

	for _, r := range retryable {
		if strings.Contains(strings.ToLower(errStr), r) {
			return true
		}
	}

	return false
}

// CallJSONWithRetry calls Claude with retry logic
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err := c.CallJSON(prompt, result)
		if err == nil {
			return nil
		}
		lastErr = err

		// Don't retry non-retryable errors
		if !isRetryableError(err) {
			return err
		}

		// Don't sleep after last attempt
		if attempt < cfg.MaxAttempts {
			delay := cfg.BaseDelay * time.Duration(1<<(attempt-1)) // Exponential backoff
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
			fmt.Printf("  Retry %d/%d after %v (error: %v)\n", attempt, cfg.MaxAttempts, delay, summarizeError(err))
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", cfg.MaxAttempts, lastErr)
}

// summarizeError returns a short version of the error for logging
func summarizeError(err error) string {
	s := err.Error()
	if len(s) > 100 {
		return s[:100] + "..."
	}
	return s
}
