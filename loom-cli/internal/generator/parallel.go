package generator

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// PhaseResult holds the result of a parallel phase execution
type PhaseResult struct {
	Name     string
	Duration time.Duration
	Error    error
	Data     interface{}
}

// Phase represents a single generation phase that can run in parallel
type Phase struct {
	Name    string
	Execute func() (interface{}, error)
}

// ParallelExecutor runs multiple phases concurrently
type ParallelExecutor struct {
	MaxConcurrent int
}

// NewParallelExecutor creates a new parallel executor
// maxConcurrent controls how many phases can run simultaneously (0 = unlimited)
func NewParallelExecutor(maxConcurrent int) *ParallelExecutor {
	return &ParallelExecutor{
		MaxConcurrent: maxConcurrent,
	}
}

// Execute runs all phases in parallel and returns results
func (e *ParallelExecutor) Execute(phases []Phase) []PhaseResult {
	results := make([]PhaseResult, len(phases))
	var wg sync.WaitGroup

	// Create semaphore for limiting concurrency
	var sem chan struct{}
	if e.MaxConcurrent > 0 {
		sem = make(chan struct{}, e.MaxConcurrent)
	}

	// Print start message
	fmt.Fprintf(os.Stderr, "\n[Parallel] Starting %d phases...\n", len(phases))

	for i, phase := range phases {
		wg.Add(1)
		go func(idx int, p Phase) {
			defer wg.Done()

			// Acquire semaphore if limited
			if sem != nil {
				sem <- struct{}{}
				defer func() { <-sem }()
			}

			start := time.Now()
			fmt.Fprintf(os.Stderr, "  [%s] Starting...\n", p.Name)

			data, err := p.Execute()

			results[idx] = PhaseResult{
				Name:     p.Name,
				Duration: time.Since(start),
				Error:    err,
				Data:     data,
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "  [%s] Failed: %v (%.1fs)\n", p.Name, err, results[idx].Duration.Seconds())
			} else {
				fmt.Fprintf(os.Stderr, "  [%s] Done (%.1fs)\n", p.Name, results[idx].Duration.Seconds())
			}
		}(i, phase)
	}

	wg.Wait()

	// Print summary
	var totalDuration time.Duration
	var maxDuration time.Duration
	errors := 0
	for _, r := range results {
		totalDuration += r.Duration
		if r.Duration > maxDuration {
			maxDuration = r.Duration
		}
		if r.Error != nil {
			errors++
		}
	}

	fmt.Fprintf(os.Stderr, "[Parallel] Complete: %d/%d succeeded (wall: %.1fs, total: %.1fs, saved: %.1fs)\n",
		len(phases)-errors, len(phases),
		maxDuration.Seconds(),
		totalDuration.Seconds(),
		(totalDuration - maxDuration).Seconds())

	return results
}

// GetResult extracts a typed result from PhaseResult
func GetResult[T any](result PhaseResult) (T, error) {
	var zero T
	if result.Error != nil {
		return zero, result.Error
	}
	if data, ok := result.Data.(T); ok {
		return data, nil
	}
	return zero, fmt.Errorf("type assertion failed for phase %s", result.Name)
}
