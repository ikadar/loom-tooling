package workflow

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ApprovalAction represents the user's action choice
type ApprovalAction int

const (
	ActionApprove ApprovalAction = iota
	ActionEdit
	ActionRegenerate
	ActionSkip
	ActionQuit
)

// WriteConfig configures a file write with optional approval
type WriteConfig struct {
	Path      string
	Content   string
	PhaseName string
	ItemCount int
	ItemType  string
	Summary   string
}

// WriteResult indicates what happened during the write
type WriteResult struct {
	Written bool
	Skipped bool
	Edited  bool
	Path    string
}

// WriteWithApproval writes a file with optional interactive approval
// Returns: result, needsRegenerate, error
func WriteWithApproval(cfg WriteConfig, interactive bool) (*WriteResult, bool, error) {
	// Write the file first
	if err := os.WriteFile(cfg.Path, []byte(cfg.Content), 0644); err != nil {
		return nil, false, fmt.Errorf("failed to write %s: %w", cfg.Path, err)
	}

	result := &WriteResult{
		Path: cfg.Path,
	}

	if !interactive {
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", cfg.Path)
		return result, false, nil
	}

	// Interactive mode: ask for approval
	action, err := askApproval(cfg)
	if err != nil {
		return nil, false, err
	}

	switch action {
	case ActionApprove:
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", cfg.Path)

	case ActionEdit:
		edited, err := editContent(cfg.Content, cfg.PhaseName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Edit error: %v\n", err)
			result.Written = true
		} else {
			if err := os.WriteFile(cfg.Path, []byte(edited), 0644); err != nil {
				return nil, false, fmt.Errorf("failed to write edited content: %w", err)
			}
			result.Written = true
			result.Edited = true
			fmt.Fprintf(os.Stderr, "  Written (edited): %s\n", cfg.Path)
		}

	case ActionSkip:
		os.Remove(cfg.Path)
		result.Skipped = true
		fmt.Fprintf(os.Stderr, "  Skipped: %s\n", cfg.Path)

	case ActionRegenerate:
		os.Remove(cfg.Path)
		return result, true, nil

	case ActionQuit:
		if confirmQuit() {
			return nil, false, fmt.Errorf("user quit")
		}
		// User cancelled quit, treat as approve
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", cfg.Path)
	}

	return result, false, nil
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

// renderPreview renders a preview box with the first N lines of content
func renderPreview(content string, maxLines int) string {
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

// askApproval prompts the user for approval action
func askApproval(cfg WriteConfig) (ApprovalAction, error) {
	fmt.Fprintf(os.Stderr, "\n%s\n", cfg.PhaseName)
	fmt.Fprintf(os.Stderr, "\nPreview (first 20 lines):\n")
	fmt.Fprintln(os.Stderr, renderPreview(cfg.Content, 20))
	fmt.Fprintf(os.Stderr, "\nGenerated: %d %s\n", cfg.ItemCount, cfg.ItemType)
	if cfg.Summary != "" {
		fmt.Fprintf(os.Stderr, "%s\n", cfg.Summary)
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

// editContent opens the content in the user's editor
func editContent(content string, filename string) (string, error) {
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

// confirmQuit asks the user to confirm quitting
func confirmQuit() bool {
	fmt.Fprintf(os.Stderr, "Are you sure you want to quit? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// HandleFileApproval handles approval for an already written file
// Returns: result, needsRegenerate, error
func HandleFileApproval(path, phaseName string, itemCount int, itemType, summary string, interactive bool) (*WriteResult, bool, error) {
	result := &WriteResult{
		Path: path,
	}

	if !interactive {
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", path)
		return result, false, nil
	}

	// Read file content for preview
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read %s: %w", path, err)
	}

	cfg := WriteConfig{
		Path:      path,
		Content:   string(content),
		PhaseName: phaseName,
		ItemCount: itemCount,
		ItemType:  itemType,
		Summary:   summary,
	}

	action, err := askApproval(cfg)
	if err != nil {
		return nil, false, err
	}

	switch action {
	case ActionApprove:
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", path)

	case ActionEdit:
		edited, err := editContent(string(content), phaseName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Edit error: %v\n", err)
			result.Written = true
		} else {
			if err := os.WriteFile(path, []byte(edited), 0644); err != nil {
				return nil, false, fmt.Errorf("failed to write edited content: %w", err)
			}
			result.Written = true
			result.Edited = true
			fmt.Fprintf(os.Stderr, "  Written (edited): %s\n", path)
		}

	case ActionSkip:
		os.Remove(path)
		result.Skipped = true
		fmt.Fprintf(os.Stderr, "  Skipped: %s\n", path)

	case ActionRegenerate:
		os.Remove(path)
		return result, true, nil

	case ActionQuit:
		if confirmQuit() {
			return nil, false, fmt.Errorf("user quit")
		}
		// User cancelled quit, treat as approve
		result.Written = true
		fmt.Fprintf(os.Stderr, "  Written: %s\n", path)
	}

	return result, false, nil
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
