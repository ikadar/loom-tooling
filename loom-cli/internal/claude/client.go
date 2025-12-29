package claude

import (
	"encoding/json"
	"fmt"
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
	args := []string{"-p", prompt, "--output-format", "json"}

	// Resume session if we have one
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude: %w", err)
	}

	var response Response
	if err := json.Unmarshal(output, &response); err != nil {
		// If JSON parsing fails, return raw output
		return strings.TrimSpace(string(output)), nil
	}

	// Store session ID for continuation
	if response.SessionID != "" {
		c.SessionID = response.SessionID
	}

	return response.Result, nil
}

// CallWithSystemPrompt calls Claude with an additional system prompt
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	args := []string{
		"-p", userPrompt,
		"--append-system-prompt", systemPrompt,
		"--output-format", "json",
	}

	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude: %w", err)
	}

	var response Response
	if err := json.Unmarshal(output, &response); err != nil {
		return strings.TrimSpace(string(output)), nil
	}

	if response.SessionID != "" {
		c.SessionID = response.SessionID
	}

	return response.Result, nil
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
			return json.Unmarshal([]byte(jsonStr), result)
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
		return fmt.Errorf("no JSON found in response: %s", response)
	}

	jsonStr := response[jsonStart : jsonEnd+1]
	return json.Unmarshal([]byte(jsonStr), result)
}
