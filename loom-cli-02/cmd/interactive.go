// Implements: l2/tech-specs.md TS-ARCH-004
// Implements: DEC-L1-016 (preview limits: 20 lines, 80 chars)
// Implements: DEC-L1-017 (editor priority: EDITOR→VISUAL→vim→nano→vi)
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"loom-cli/internal/domain"
)

// Box drawing characters (Unicode)
const (
	boxTopLeft     = "┌"
	boxTopRight    = "┐"
	boxBottomLeft  = "└"
	boxBottomRight = "┘"
	boxHorizontal  = "─"
	boxVertical    = "│"
	boxTeeRight    = "├"
	boxTeeLeft     = "┤"
)

// showPreview renders a preview of content in a box.
//
// Implements: DEC-L1-016 (max 20 lines, 80 chars)
func showPreview(filename, content string) {
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	// Truncate if needed
	if len(lines) > domain.MaxPreviewLines {
		lines = lines[:domain.MaxPreviewLines]
	}

	// Print box
	printBoxTop(filename)
	for _, line := range lines {
		printBoxLine(truncateLine(line, domain.MaxLineWidth-4))
	}
	if totalLines > domain.MaxPreviewLines {
		printBoxLine(fmt.Sprintf("(truncated - showing %d of %d lines)",
			domain.MaxPreviewLines, totalLines))
	}
	printBoxBottom()
}

func printBoxTop(title string) {
	width := domain.MaxLineWidth - 2
	titleLen := len(title) + 2 // space on each side
	leftPad := (width - titleLen) / 2
	rightPad := width - titleLen - leftPad

	fmt.Fprint(os.Stderr, boxTopLeft)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, leftPad))
	fmt.Fprintf(os.Stderr, " %s ", title)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, rightPad))
	fmt.Fprintln(os.Stderr, boxTopRight)
}

func printBoxLine(content string) {
	width := domain.MaxLineWidth - 4 // Account for box chars and padding
	if len(content) > width {
		content = content[:width-3] + "..."
	}
	padding := width - len(content)
	fmt.Fprintf(os.Stderr, "%s %s%s %s\n",
		boxVertical, content, strings.Repeat(" ", padding), boxVertical)
}

func printBoxBottom() {
	width := domain.MaxLineWidth - 2
	fmt.Fprint(os.Stderr, boxBottomLeft)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, width))
	fmt.Fprintln(os.Stderr, boxBottomRight)
}

func truncateLine(line string, maxLen int) string {
	if len(line) <= maxLen {
		return line
	}
	return line[:maxLen-3] + "..."
}

// requestApproval prompts user for approval action.
//
// Implements: TS-ARCH-004
// Returns: ActionApprove, ActionEdit, ActionRegenerate, ActionSkip, ActionQuit
func requestApproval() domain.ApprovalAction {
	for {
		fmt.Fprint(os.Stderr, "[A]pprove  [E]dit  [R]egenerate  [S]kip  [Q]uit  > ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))

		// Empty input defaults to approve
		if input == "" || input == "a" {
			return domain.ActionApprove
		}

		switch input {
		case "e":
			return domain.ActionEdit
		case "r":
			return domain.ActionRegenerate
		case "s":
			return domain.ActionSkip
		case "q":
			return domain.ActionQuit
		default:
			fmt.Fprintln(os.Stderr, "Invalid option. Please enter A, E, R, S, or Q.")
		}
	}
}

// confirmQuit asks for confirmation before quitting.
func confirmQuit() bool {
	fmt.Fprint(os.Stderr, "Are you sure you want to quit? Approved files will be saved. [y/N] > ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// openInEditor opens content in an external editor.
//
// Implements: DEC-L1-017 (editor priority)
func openInEditor(content string) (string, error) {
	editor := findEditor()
	if editor == "" {
		return "", fmt.Errorf("no editor found (set $EDITOR)")
	}

	// Create temp file with .md extension for syntax highlighting
	tmpFile, err := os.CreateTemp("", "loom-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write current content
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Launch editor
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}

	// Read modified content
	modified, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return string(modified), nil
}

// findEditor finds an available text editor.
//
// Implements: DEC-L1-017
// Priority: $EDITOR → $VISUAL → vim → nano → vi
func findEditor() string {
	// Check environment variables first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// Fallback chain
	editors := []string{"vim", "nano", "vi"}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	return "vi" // Ultimate fallback
}
