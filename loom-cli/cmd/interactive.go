package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InteractiveConfig holds configuration for interactive mode
type InteractiveConfig struct {
	Enabled bool
}

// ApprovalAction represents the user's action choice
type ApprovalAction int

const (
	ActionApprove ApprovalAction = iota
	ActionEdit
	ActionRegenerate
	ActionSkip
	ActionQuit
)

// PhaseResult holds the generated content for a phase
type PhaseResult struct {
	PhaseName   string
	FileName    string
	Content     string
	Summary     string
	ItemCount   int
	ItemType    string
}

// Preview box characters
const (
	boxTopLeft     = "┌"
	boxTopRight    = "┐"
	boxBottomLeft  = "└"
	boxBottomRight = "┘"
	boxHorizontal  = "─"
	boxVertical    = "│"
)

// RenderPreview renders a preview box with the first N lines of content
func RenderPreview(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "...")
	}

	// Find max line length
	maxLen := 60
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	if maxLen > 80 {
		maxLen = 80
	}

	var sb strings.Builder

	// Top border
	sb.WriteString(boxTopLeft)
	sb.WriteString(strings.Repeat(boxHorizontal, maxLen+2))
	sb.WriteString(boxTopRight)
	sb.WriteString("\n")

	// Content lines
	for _, line := range lines {
		// Truncate long lines
		if len(line) > maxLen {
			line = line[:maxLen-3] + "..."
		}
		sb.WriteString(boxVertical)
		sb.WriteString(" ")
		sb.WriteString(line)
		sb.WriteString(strings.Repeat(" ", maxLen-len(line)))
		sb.WriteString(" ")
		sb.WriteString(boxVertical)
		sb.WriteString("\n")
	}

	// Bottom border
	sb.WriteString(boxBottomLeft)
	sb.WriteString(strings.Repeat(boxHorizontal, maxLen+2))
	sb.WriteString(boxBottomRight)

	return sb.String()
}

// AskApproval prompts the user for approval action
func AskApproval(result PhaseResult) (ApprovalAction, error) {
	fmt.Fprintf(os.Stderr, "\n%s\n", result.PhaseName)
	fmt.Fprintf(os.Stderr, "\nPreview (first 20 lines):\n")
	fmt.Fprintln(os.Stderr, RenderPreview(result.Content, 20))
	fmt.Fprintf(os.Stderr, "\nGenerated: %d %s\n", result.ItemCount, result.ItemType)
	if result.Summary != "" {
		fmt.Fprintf(os.Stderr, "%s\n", result.Summary)
	}
	fmt.Fprintf(os.Stderr, "\n[A]pprove / [E]dit / [R]egenerate / [S]kip / [Q]uit? ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ActionQuit, err
	}

	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "a", "approve", "":
		return ActionApprove, nil
	case "e", "edit":
		return ActionEdit, nil
	case "r", "regenerate":
		return ActionRegenerate, nil
	case "s", "skip":
		return ActionSkip, nil
	case "q", "quit":
		return ActionQuit, nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown action '%s', defaulting to Approve\n", input)
		return ActionApprove, nil
	}
}

// EditContent opens the content in the user's editor
func EditContent(content string, filename string) (string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "loom-edit-*.md")
	if err != nil {
		return content, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write content to temp file
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return content, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors
		for _, e := range []string{"vim", "nano", "vi"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}
	if editor == "" {
		return content, fmt.Errorf("no editor found (set $EDITOR)")
	}

	// Open editor
	fmt.Fprintf(os.Stderr, "Opening %s in %s...\n", filename, editor)
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return content, fmt.Errorf("editor failed: %w", err)
	}

	// Read edited content
	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return content, fmt.Errorf("failed to read edited file: %w", err)
	}

	return string(edited), nil
}

// ConfirmQuit asks the user to confirm quitting
func ConfirmQuit() bool {
	fmt.Fprintf(os.Stderr, "Are you sure you want to quit? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// PrintInteractiveHeader prints the interactive mode header
func PrintInteractiveHeader() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "   INTERACTIVE MODE")
	fmt.Fprintln(os.Stderr, "   Review each generated document")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Actions:")
	fmt.Fprintln(os.Stderr, "  [A]pprove    - Accept and write to file")
	fmt.Fprintln(os.Stderr, "  [E]dit       - Open in $EDITOR")
	fmt.Fprintln(os.Stderr, "  [R]egenerate - Generate again with AI")
	fmt.Fprintln(os.Stderr, "  [S]kip       - Skip this file")
	fmt.Fprintln(os.Stderr, "  [Q]uit       - Exit derivation")
	fmt.Fprintln(os.Stderr, "")
}
