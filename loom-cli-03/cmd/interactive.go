// Package cmd provides CLI commands for loom-cli.
//
// This file implements interactive mode utilities.
//
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
	"loom-cli/prompts"
)

// Ensure prompts package is imported
var _ = prompts.Derivation

// Constants for preview rendering
//
// Implements: DEC-L1-016
const (
	maxPreviewLines = 20
	maxLineWidth    = 80
	boxWidth        = 79
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

// showPreview renders a truncated preview of content in a box.
//
// Implements: TS-ARCH-004
func showPreview(filename, content string) {
	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	// Truncate if needed
	truncated := false
	if len(lines) > maxPreviewLines {
		lines = lines[:maxPreviewLines]
		truncated = true
	}

	// Print top border
	title := fmt.Sprintf(" Preview: %s ", filename)
	titleLen := len(title)
	leftPad := (boxWidth - titleLen - 2) / 2
	rightPad := boxWidth - titleLen - 2 - leftPad

	fmt.Fprint(os.Stderr, boxTopLeft)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, leftPad))
	fmt.Fprint(os.Stderr, title)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, rightPad))
	fmt.Fprintln(os.Stderr, boxTopRight)

	// Print separator
	fmt.Fprint(os.Stderr, boxTeeRight)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, boxWidth-2))
	fmt.Fprintln(os.Stderr, boxTeeLeft)

	// Print content lines
	for _, line := range lines {
		printBoxLine(line)
	}

	// Print truncation notice if applicable
	if truncated {
		notice := fmt.Sprintf("(truncated - showing %d of %d lines)", maxPreviewLines, totalLines)
		printBoxLine(notice)
	}

	// Print bottom border
	fmt.Fprint(os.Stderr, boxBottomLeft)
	fmt.Fprint(os.Stderr, strings.Repeat(boxHorizontal, boxWidth-2))
	fmt.Fprintln(os.Stderr, boxBottomRight)
}

// printBoxLine prints a line inside the box, truncating if necessary.
func printBoxLine(line string) {
	// Truncate line if too long
	maxContent := maxLineWidth - 4 // Account for "│ " and " │"
	if len(line) > maxContent {
		line = line[:maxContent-3] + "..."
	}

	// Pad to fill the box
	padding := maxContent - len(line)
	if padding < 0 {
		padding = 0
	}

	fmt.Fprintf(os.Stderr, "%s %s%s %s\n",
		boxVertical,
		line,
		strings.Repeat(" ", padding),
		boxVertical)
}

// requestApproval prompts the user for an approval action.
//
// Implements: TS-ARCH-004
func requestApproval() domain.ApprovalAction {
	fmt.Fprint(os.Stderr, "\n[A]pprove  [E]dit  [R]egenerate  [S]kip  [Q]uit  > ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return domain.ActionApprove
	}

	input = strings.TrimSpace(strings.ToLower(input))

	// Empty input defaults to approve
	if input == "" {
		return domain.ActionApprove
	}

	switch input[0] {
	case 'a':
		return domain.ActionApprove
	case 'e':
		return domain.ActionEdit
	case 'r':
		return domain.ActionRegenerate
	case 's':
		return domain.ActionSkip
	case 'q':
		return domain.ActionQuit
	default:
		fmt.Fprintln(os.Stderr, "Invalid input. Please enter a, e, r, s, or q.")
		return requestApproval()
	}
}

// openInEditor opens content in an external editor and returns the modified content.
//
// Implements: TS-ARCH-004, DEC-L1-017
func openInEditor(content string) (string, error) {
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

	// Find editor
	editor := findEditor()
	if editor == "" {
		return "", fmt.Errorf("no editor found (set $EDITOR)")
	}

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
		return "", fmt.Errorf("failed to read edited file: %w", err)
	}

	return string(modified), nil
}

// findEditor finds an available editor following priority order.
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
	fallbacks := []string{"vim", "nano", "vi"}
	for _, editor := range fallbacks {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	// Ultimate fallback
	return "vi"
}

// confirmQuit asks the user to confirm quitting.
//
// Implements: TS-ARCH-004
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
