// Package workflow provides interactive workflow utilities.
//
// Implements: l2/package-structure.md PKG-009
// See: l2/internal-api.md
package workflow

import (
	"fmt"
	"os"
)

// Progress tracks and displays progress.
type Progress struct {
	Label   string
	Total   int
	Current int
}

// NewProgress creates a progress tracker.
func NewProgress(label string, total int) *Progress {
	return &Progress{
		Label:   label,
		Total:   total,
		Current: 0,
	}
}

// Increment advances progress by 1.
func (p *Progress) Increment() {
	p.Current++
	p.display("")
}

// Update sets current progress with message.
func (p *Progress) Update(current int, message string) {
	p.Current = current
	p.display(message)
}

// Done marks progress complete.
func (p *Progress) Done() {
	p.Current = p.Total
	p.display("Complete")
}

// display outputs the current progress.
// Format: [label] current/total message...
func (p *Progress) display(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, "[%s] %d/%d %s\n", p.Label, p.Current, p.Total, message)
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %d/%d\n", p.Label, p.Current, p.Total)
	}
}
