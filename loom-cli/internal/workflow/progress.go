package workflow

import (
	"fmt"
	"os"
	"strings"
)

// Progress tracks and displays progress for long-running operations
type Progress struct {
	Label     string
	Total     int
	Completed int
	Width     int // Progress bar width in characters
}

// NewProgress creates a new progress tracker
func NewProgress(label string, total int) *Progress {
	return &Progress{
		Label: label,
		Total: total,
		Width: 30,
	}
}

// Update updates the progress and redraws the progress bar
func (p *Progress) Update(completed int) {
	p.Completed = completed
	p.draw()
}

// Increment increases completed by 1 and redraws
func (p *Progress) Increment() {
	p.Completed++
	p.draw()
}

// draw renders the progress bar to stderr
func (p *Progress) draw() {
	if p.Total == 0 {
		return
	}

	pct := float64(p.Completed) / float64(p.Total)
	filled := int(pct * float64(p.Width))
	if filled > p.Width {
		filled = p.Width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.Width-filled)
	fmt.Fprintf(os.Stderr, "\r  %s: [%s] %d/%d (%.0f%%)", p.Label, bar, p.Completed, p.Total, pct*100)
}

// Done finishes the progress bar with a newline
func (p *Progress) Done() {
	fmt.Fprintln(os.Stderr)
}

// DoneWithMessage finishes the progress bar with a custom message
func (p *Progress) DoneWithMessage(msg string) {
	fmt.Fprintf(os.Stderr, "\r  %s: %s\n", p.Label, msg)
}

// Error shows an error state
func (p *Progress) Error(err error) {
	fmt.Fprintf(os.Stderr, "\r  %s: ERROR - %v\n", p.Label, err)
}
