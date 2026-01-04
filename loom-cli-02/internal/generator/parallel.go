// Implements: l2/package-structure.md PKG-006
// See: l2/internal-api.md
package generator

import (
	"sync"
)

// DefaultWorkers is the default number of parallel workers.
// Balance between speed and rate limits.
const DefaultWorkers = 3

// ProcessInParallel processes items in parallel with worker limit.
//
// Generic parallel processor for any derivation task.
// Uses sync.WaitGroup for coordination and buffered channels for results.
// Collects errors, doesn't fail on first error.
//
// Implements: l2/package-structure.md PKG-006
func ProcessInParallel[T, R any](items []T, fn func(T) (R, error), workers int) ([]R, error) {
	if workers <= 0 {
		workers = DefaultWorkers
	}

	if len(items) == 0 {
		return []R{}, nil
	}

	// For small item counts, don't use more workers than items
	if workers > len(items) {
		workers = len(items)
	}

	type indexedResult struct {
		index  int
		result R
		err    error
	}

	// Channels for work distribution and results
	jobs := make(chan struct {
		index int
		item  T
	}, len(items))
	results := make(chan indexedResult, len(items))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				result, err := fn(job.item)
				results <- indexedResult{
					index:  job.index,
					result: result,
					err:    err,
				}
			}
		}()
	}

	// Send jobs
	for i, item := range items {
		jobs <- struct {
			index int
			item  T
		}{index: i, item: item}
	}
	close(jobs)

	// Wait for completion in separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	orderedResults := make([]R, len(items))
	var errors []error

	for r := range results {
		if r.err != nil {
			errors = append(errors, r.err)
		} else {
			orderedResults[r.index] = r.result
		}
	}

	// Return first error if any (we collected all, but return first for simplicity)
	if len(errors) > 0 {
		return orderedResults, errors[0]
	}

	return orderedResults, nil
}

// ProcessSequentially processes items one by one.
// Useful for debugging or when rate limiting is strict.
func ProcessSequentially[T, R any](items []T, fn func(T) (R, error)) ([]R, error) {
	results := make([]R, len(items))

	for i, item := range items {
		result, err := fn(item)
		if err != nil {
			return results, err
		}
		results[i] = result
	}

	return results, nil
}
