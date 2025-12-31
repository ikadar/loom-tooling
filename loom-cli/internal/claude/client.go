package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Response represents the JSON output from claude -p
type Response struct {
	Result    string  `json:"result"`
	SessionID string  `json:"session_id"`
	CostUSD   float64 `json:"cost_usd"`
}

// Client wraps Claude CLI calls
type Client struct {
	SessionID string
	Verbose   bool
}

// NewClient creates a new Claude client
func NewClient() *Client {
	return &Client{}
}

// Call sends a prompt to Claude and returns the response
func (c *Client) Call(prompt string) (string, error) {
	// Don't use --output-format json as it returns empty result for multi-turn responses
	args := []string{"-p", prompt}

	// Resume session if we have one
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)
	// Set high output token limit for large generations
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CallWithSystemPrompt calls Claude with an additional system prompt
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	// Don't use --output-format json as it returns empty result for multi-turn responses
	args := []string{
		"-p", userPrompt,
		"--append-system-prompt", systemPrompt,
	}

	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)
	// Set high output token limit for large generations
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// sanitizeJSON attempts to fix common JSON issues from LLM output
func sanitizeJSON(jsonStr string) string {
	// Fix line breaks inside strings by reconstructing the JSON
	// This is a simple approach - replace literal newlines within quoted strings
	var result strings.Builder
	inString := false
	escaped := false

	for i := 0; i < len(jsonStr); i++ {
		c := jsonStr[i]

		if escaped {
			result.WriteByte(c)
			escaped = false
			continue
		}

		if c == '\\' {
			result.WriteByte(c)
			escaped = true
			continue
		}

		if c == '"' {
			inString = !inString
			result.WriteByte(c)
			continue
		}

		if inString && (c == '\n' || c == '\r') {
			// Replace literal newline with escaped version
			result.WriteString("\\n")
			continue
		}

		result.WriteByte(c)
	}

	return result.String()
}

// CallJSON sends a prompt expecting JSON response
func (c *Client) CallJSON(prompt string, result interface{}) error {
	response, err := c.Call(prompt)
	if err != nil {
		return err
	}

	// Try to extract JSON from markdown code block first
	codeBlockStart := strings.Index(response, "```json")
	if codeBlockStart != -1 {
		codeBlockStart += 7 // skip ```json
		codeBlockEnd := strings.Index(response[codeBlockStart:], "```")
		if codeBlockEnd != -1 {
			jsonStr := strings.TrimSpace(response[codeBlockStart : codeBlockStart+codeBlockEnd])
			// Sanitize JSON - fix line breaks inside strings
			jsonStr = sanitizeJSON(jsonStr)
			if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
				// Show debug info on parse error
				preview := jsonStr
				if len(preview) > 500 {
					preview = preview[:250] + "\n...[truncated]...\n" + preview[len(preview)-250:]
				}
				return fmt.Errorf("JSON parse error from code block: %w\nJSON preview:\n%s", err, preview)
			}
			return nil
		}
	}

	// Fallback: Try to extract raw JSON from response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		// Try array
		jsonStart = strings.Index(response, "[")
		jsonEnd = strings.LastIndex(response, "]")
	}

	if jsonStart == -1 || jsonEnd == -1 {
		preview := response
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		return fmt.Errorf("no JSON found in response: %s", preview)
	}

	jsonStr := response[jsonStart : jsonEnd+1]
	// Sanitize JSON - fix line breaks inside strings
	jsonStr = sanitizeJSON(jsonStr)
	if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
		// Show debug info on parse error
		preview := jsonStr
		if len(preview) > 500 {
			preview = preview[:250] + "\n...[truncated]...\n" + preview[len(preview)-250:]
		}
		return fmt.Errorf("JSON parse error: %w\nJSON preview:\n%s", err, preview)
	}
	return nil
}
