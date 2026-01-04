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
//
// Implements: l2/internal-api.md
type Progress struct {
	Label   string
	Total   int
	Current int
}

// NewProgress creates a progress tracker.
//
// Implements: l2/internal-api.md
func NewProgress(label string, total int) *Progress {
	p := &Progress{
		Label:   label,
		Total:   total,
		Current: 0,
	}
	p.display("")
	return p
}

// Increment advances progress by 1.
//
// Implements: l2/internal-api.md
func (p *Progress) Increment() {
	p.Current++
	p.display("")
}

// Update sets current progress with message.
//
// Implements: l2/internal-api.md
func (p *Progress) Update(current int, message string) {
	p.Current = current
	p.display(message)
}

// Done marks progress complete.
//
// Implements: l2/internal-api.md
func (p *Progress) Done() {
	p.Current = p.Total
	fmt.Fprintf(os.Stderr, "\r[%s] %d/%d ✓ Complete\n", p.Label, p.Current, p.Total)
}

// display renders the progress bar.
// Format: [label] current/total message
func (p *Progress) display(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, "\r[%s] %d/%d %s...", p.Label, p.Current, p.Total, message)
	} else {
		fmt.Fprintf(os.Stderr, "\r[%s] %d/%d", p.Label, p.Current, p.Total)
	}
}
