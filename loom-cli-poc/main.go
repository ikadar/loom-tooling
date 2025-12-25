package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ClaudeResponse represents the JSON output from claude -p
type ClaudeResponse struct {
	Result    string  `json:"result"`
	SessionID string  `json:"session_id"`
	CostUSD   float64 `json:"cost_usd"`
}

// DerivationResult is returned to the plugin
type DerivationResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// SECRET PROMPT - This would be encrypted/obfuscated in production
var derivationPrompt = `You are an AI documentation derivation expert.

Given the following L0 user story, derive:
1. Acceptance Criteria (AC) in Given/When/Then format
2. Business Rules (BR) with enforcement mechanisms

Format your output EXACTLY as follows:

## Acceptance Criteria

### AC-001 – [Title]
**Given** [precondition]
**When** [action]
**Then** [outcome]

### AC-002 – [Title]
...

## Business Rules

### BR-001 – [Title]
**Rule:** [Statement]
**Enforcement:** [How it's enforced]

### BR-002 – [Title]
...

---

Now derive from this L0 input:

`

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "derive":
		handleDerive()
	case "prompt":
		handlePrompt()
	case "version":
		fmt.Println("loom-cli-poc v0.1.0")
	default:
		// Legacy: treat as direct prompt
		prompt := strings.Join(os.Args[1:], " ")
		result, err := callClaude(prompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	}
}

func printUsage() {
	fmt.Println(`loom-cli-poc - AI-DOP CLI Proof of Concept

Usage:
  loom-cli-poc derive <input-file> [--format json|text]
  loom-cli-poc prompt <your prompt>
  loom-cli-poc version

Commands:
  derive    Derive L1 (AC + BR) from L0 input file
  prompt    Send a prompt directly to Claude
  version   Show version

Examples:
  loom-cli-poc derive ./user-story.md
  loom-cli-poc derive ./user-story.md --format json
  loom-cli-poc prompt "What is 2+2?"`)
}

func handleDerive() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: input file required")
		fmt.Fprintln(os.Stderr, "Usage: loom-cli-poc derive <input-file> [--format json|text]")
		os.Exit(1)
	}

	inputFile := os.Args[2]
	format := "text"

	// Check for --format flag
	for i, arg := range os.Args {
		if arg == "--format" && i+1 < len(os.Args) {
			format = os.Args[i+1]
		}
	}

	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		outputError(format, fmt.Sprintf("Failed to read input file: %v", err))
		os.Exit(1)
	}

	// Combine secret prompt with user input
	fullPrompt := derivationPrompt + string(content)

	// Call Claude
	result, err := callClaude(fullPrompt)
	if err != nil {
		outputError(format, fmt.Sprintf("Claude error: %v", err))
		os.Exit(1)
	}

	// Output result
	if format == "json" {
		output := DerivationResult{
			Success: true,
			Output:  result,
		}
		jsonBytes, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Println(result)
	}
}

func handlePrompt() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: prompt required")
		os.Exit(1)
	}

	prompt := strings.Join(os.Args[2:], " ")
	result, err := callClaude(prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func outputError(format string, message string) {
	if format == "json" {
		output := DerivationResult{
			Success: false,
			Error:   message,
		}
		jsonBytes, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	}
}

func callClaude(prompt string) (string, error) {
	cmd := exec.Command("claude", "-p", prompt, "--output-format", "json")

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude: %w", err)
	}

	var response ClaudeResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return string(output), nil
	}

	return response.Result, nil
}
