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
	fmt.Fprintf(os.Stderr, "%s [y/N] > ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}

// ShowDiff displays a diff between before and after content.
//
// Implements: l2/internal-api.md
func ShowDiff(before, after string) error {
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")

	fmt.Fprintln(os.Stderr, "--- Before")
	fmt.Fprintln(os.Stderr, "+++ After")
	fmt.Fprintln(os.Stderr, "")

	// Simple line-by-line diff
	maxLines := len(beforeLines)
	if len(afterLines) > maxLines {
		maxLines = len(afterLines)
	}

	for i := 0; i < maxLines; i++ {
		var bLine, aLine string
		if i < len(beforeLines) {
			bLine = beforeLines[i]
		}
		if i < len(afterLines) {
			aLine = afterLines[i]
		}

		if bLine != aLine {
			if bLine != "" {
				fmt.Fprintf(os.Stderr, "- %s\n", bLine)
			}
			if aLine != "" {
				fmt.Fprintf(os.Stderr, "+ %s\n", aLine)
			}
		} else if bLine != "" {
			fmt.Fprintf(os.Stderr, "  %s\n", bLine)
		}
	}

	return nil
}

// ShowPreview displays truncated content preview.
//
// Implements: DEC-L1-016 (20 lines, 80 chars)
func ShowPreview(content string, maxLines, maxWidth int) string {
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	// Truncate lines
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	var result strings.Builder

	for _, line := range lines {
		// Truncate width
		if len(line) > maxWidth {
			line = line[:maxWidth-3] + "..."
		}
		result.WriteString(line)
		result.WriteString("\n")
	}

	if totalLines > maxLines {
		result.WriteString(fmt.Sprintf("... (%d more lines)\n", totalLines-maxLines))
	}

	return result.String()
}
