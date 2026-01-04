// Package generator provides document generation utilities.
//
// This file implements parallel processing utilities.
//
// Implements: l2/tech-specs.md TS-ARCH-003
package generator

import (
	"sync"
)

// ProcessInParallel processes items in parallel with bounded concurrency.
//
// Implements: TS-ARCH-003
func ProcessInParallel[T any, R any](items []T, fn func(T) (R, error), workers int) ([]R, []error) {
	if workers <= 0 {
		workers = 3 // Default max concurrent
	}

	results := make([]R, len(items))
	errors := make([]error, len(items))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers)

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it T) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := fn(it)
			results[idx] = result
			errors[idx] = err
		}(i, item)
	}

	wg.Wait()

	// Filter out nil errors
	var actualErrors []error
	for _, err := range errors {
		if err != nil {
			actualErrors = append(actualErrors, err)
		}
	}

	return results, actualErrors
}

// ParallelExecutor manages parallel task execution.
//
// Implements: TS-ARCH-003
type ParallelExecutor struct {
	MaxConcurrent int
	semaphore     chan struct{}
}

// NewParallelExecutor creates a new executor with specified concurrency.
func NewParallelExecutor(maxConcurrent int) *ParallelExecutor {
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}
	return &ParallelExecutor{
		MaxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

// Task represents a unit of work.
type Task func() (interface{}, error)

// Result holds the result of a task execution.
type Result struct {
	Value interface{}
	Error error
	Index int
}

// Execute runs tasks in parallel with bounded concurrency.
func (pe *ParallelExecutor) Execute(tasks []Task) []Result {
	results := make([]Result, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t Task) {
			defer wg.Done()

			// Acquire semaphore
			pe.semaphore <- struct{}{}
			defer func() { <-pe.semaphore }()

			value, err := t()
			results[idx] = Result{
				Value: value,
				Error: err,
				Index: idx,
			}
		}(i, task)
	}

	wg.Wait()
	return results
}
