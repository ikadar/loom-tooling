// Package workflow provides interactive workflow utilities.
//
// Implements: l2/package-structure.md PKG-009
// See: l2/internal-api.md
package workflow

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RequestApproval prompts user for yes/no.
//
// Implements: l2/internal-api.md
func RequestApproval(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "%s [y/n]: ", prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}

// ShowDiff displays a diff between before and after content.
func ShowDiff(before, after string) error {
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")

	fmt.Fprintln(os.Stderr, "--- before")
	fmt.Fprintln(os.Stderr, "+++ after")

	// Simple line-by-line diff
	maxLines := len(beforeLines)
	if len(afterLines) > maxLines {
		maxLines = len(afterLines)
	}

	for i := 0; i < maxLines; i++ {
		var beforeLine, afterLine string
		if i < len(beforeLines) {
			beforeLine = beforeLines[i]
		}
		if i < len(afterLines) {
			afterLine = afterLines[i]
		}

		if beforeLine != afterLine {
			if beforeLine != "" {
				fmt.Fprintf(os.Stderr, "-%s\n", beforeLine)
			}
			if afterLine != "" {
				fmt.Fprintf(os.Stderr, "+%s\n", afterLine)
			}
		}
	}

	return nil
}

// ShowPreview displays truncated content preview.
//
// Implements: DEC-L1-016 (20 lines, 80 chars)
func ShowPreview(content string, maxLines, maxWidth int) string {
	if maxLines <= 0 {
		maxLines = 20
	}
	if maxWidth <= 0 {
		maxWidth = 80
	}

	lines := strings.Split(content, "\n")
	var result []string

	for i, line := range lines {
		if i >= maxLines {
			result = append(result, fmt.Sprintf("... (%d more lines)", len(lines)-maxLines))
			break
		}

		if len(line) > maxWidth {
			line = line[:maxWidth-3] + "..."
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
