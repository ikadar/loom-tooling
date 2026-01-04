// Package claude provides the Claude CLI integration layer.
//
// Implements: l2/tech-specs.md TS-ARCH-002 (Anti-Corruption Layer via Claude Code CLI)
// See: l1/bounded-context-map.md (ACL pattern)
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps the Claude Code CLI for AI-powered operations.
//
// Implements: l2/internal-api.md internal/claude
type Client struct {
	SessionID string // For multi-turn conversations
	Verbose   bool   // Debug output
}

// NewClient creates a new Claude client.
//
// Implements: l2/internal-api.md
func NewClient() *Client {
	return &Client{}
}

// Call sends a prompt to Claude and returns the raw text response.
//
// Implements: l2/tech-specs.md TS-ARCH-002
func (c *Client) Call(prompt string) (string, error) {
	args := []string{"-p", prompt, "--output-format", "text"}
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	// Implements: DEC-L1-002 (100,000 max output tokens)
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[claude] Calling with prompt length: %d\n", len(prompt))
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude call failed: %s: %s", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("claude call failed: %w", err)
	}

	return string(output), nil
}

// CallWithSystemPrompt calls Claude with a system prompt and user prompt.
//
// Implements: l2/internal-api.md
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	args := []string{"-p", userPrompt, "--output-format", "text"}
	if systemPrompt != "" {
		args = append(args, "--system-prompt", systemPrompt)
	}
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	// Implements: DEC-L1-002 (100,000 max output tokens)
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[claude] Calling with system prompt length: %d, user prompt length: %d\n",
			len(systemPrompt), len(userPrompt))
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude call failed: %s: %s", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("claude call failed: %w", err)
	}

	return string(output), nil
}

// CallJSON sends a prompt and parses the response as JSON.
//
// Implements: l2/tech-specs.md SEQ-JSON-001
func (c *Client) CallJSON(prompt string, result interface{}) error {
	response, err := c.Call(prompt)
	if err != nil {
		return err
	}

	if err := extractJSON(response, result); err != nil {
		// Step 4: Conversion Fallback - ask Claude to fix
		if c.Verbose {
			fmt.Fprintf(os.Stderr, "[claude] JSON extraction failed, trying conversion fallback\n")
		}

		fixPrompt := fmt.Sprintf("Convert the following to valid JSON only, no explanation:\n\n%s", response)
		fixedResponse, fixErr := c.Call(fixPrompt)
		if fixErr != nil {
			return fmt.Errorf("JSON extraction failed: %w (fallback also failed: %v)", err, fixErr)
		}

		if err := extractJSON(fixedResponse, result); err != nil {
			return fmt.Errorf("JSON extraction failed after fallback: %w", err)
		}
	}

	return nil
}

// CallJSONWithRetry calls Claude with retry logic and parses JSON response.
//
// Implements: l2/tech-specs.md TS-RETRY-001
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error {
	return WithRetry(cfg, func() error {
		return c.CallJSON(prompt, result)
	})
}

// extractJSON attempts to parse JSON from a Claude response.
//
// Implements: l2/tech-specs.md SEQ-JSON-001
// Strategy:
// 1. Markdown Code Block Extraction
// 2. Raw JSON Extraction
// 3. JSON Sanitization
func extractJSON(response string, result interface{}) error {
	// Step 1: Markdown Code Block Extraction
	if idx := strings.Index(response, "```json"); idx != -1 {
		start := idx + len("```json")
		// Skip any whitespace/newline after ```json
		for start < len(response) && (response[start] == '\n' || response[start] == '\r' || response[start] == ' ') {
			start++
		}
		end := strings.Index(response[start:], "```")
		if end != -1 {
			jsonStr := strings.TrimSpace(response[start : start+end])
			if err := json.Unmarshal([]byte(jsonStr), result); err == nil {
				return nil
			}
			// Try with sanitization
			sanitized := sanitizeJSON(jsonStr)
			if err := json.Unmarshal([]byte(sanitized), result); err == nil {
				return nil
			}
		}
	}

	// Also try generic ``` block (some responses use ```\n{...}\n```)
	if idx := strings.Index(response, "```\n"); idx != -1 {
		start := idx + len("```\n")
		end := strings.Index(response[start:], "```")
		if end != -1 {
			jsonStr := strings.TrimSpace(response[start : start+end])
			if len(jsonStr) > 0 && (jsonStr[0] == '{' || jsonStr[0] == '[') {
				if err := json.Unmarshal([]byte(jsonStr), result); err == nil {
					return nil
				}
				sanitized := sanitizeJSON(jsonStr)
				if err := json.Unmarshal([]byte(sanitized), result); err == nil {
					return nil
				}
			}
		}
	}

	// Step 2: Raw JSON Extraction
	jsonStr := extractRawJSON(response)
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), result); err == nil {
			return nil
		}

		// Step 3: JSON Sanitization
		sanitized := sanitizeJSON(jsonStr)
		if err := json.Unmarshal([]byte(sanitized), result); err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not extract valid JSON from response")
}

// extractRawJSON finds JSON object or array in response text.
func extractRawJSON(response string) string {
	// Find first { or [
	startObj := strings.Index(response, "{")
	startArr := strings.Index(response, "[")

	var start int
	var isObject bool

	switch {
	case startObj == -1 && startArr == -1:
		return ""
	case startObj == -1:
		start = startArr
		isObject = false
	case startArr == -1:
		start = startObj
		isObject = true
	case startObj < startArr:
		start = startObj
		isObject = true
	default:
		start = startArr
		isObject = false
	}

	// Find matching closing bracket
	var end int
	if isObject {
		end = findMatchingBracket(response[start:], '{', '}')
	} else {
		end = findMatchingBracket(response[start:], '[', ']')
	}

	if end == -1 {
		return ""
	}

	return response[start : start+end+1]
}

// findMatchingBracket finds the position of the matching closing bracket.
func findMatchingBracket(s string, open, close byte) int {
	depth := 0
	inString := false
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if c == open {
			depth++
		} else if c == close {
			depth--
			if depth == 0 {
				return i
			}
		}
	}

	return -1
}

// sanitizeJSON fixes common LLM JSON output issues.
//
// Implements: l2/tech-specs.md SEQ-JSON-001 Step 3
func sanitizeJSON(jsonStr string) string {
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

		// Replace literal newlines inside strings
		if inString && (c == '\n' || c == '\r') {
			result.WriteString("\\n")
			continue
		}

		// Remove other control characters inside strings
		if inString && c < 32 && c != '\t' {
			continue
		}

		result.WriteByte(c)
	}

	return result.String()
}

// BuildPrompt injects context documents into a prompt template.
//
// Implements: l2/internal-api.md prompts.BuildPrompt
// All prompts use XML-style context markers:
// <context>
// </context>
func BuildPrompt(template string, documents ...string) string {
	if len(documents) == 0 {
		return template
	}
	context := strings.Join(documents, "\n\n---\n\n")
	return strings.Replace(template, "</context>", context+"\n</context>", 1)
}
