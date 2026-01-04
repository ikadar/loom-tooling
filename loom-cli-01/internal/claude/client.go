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

// Client wraps the Claude CLI for AI-powered derivation.
// Implements: l2/tech-specs.md TS-ARCH-002
type Client struct {
	SessionID string
	Verbose   bool
}

// NewClient creates a new Claude client.
func NewClient() *Client {
	return &Client{}
}

// Call invokes Claude with the given prompt and returns the response.
// Implements: l2/tech-specs.md TS-ARCH-002 (CLI Invocation)
func (c *Client) Call(prompt string) (string, error) {
	args := []string{"-p", prompt}
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	// Set max output tokens for large derivations
	// Implements: DEC-L1-002 (100,000 max output tokens)
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[claude] Calling with prompt length: %d\n", len(prompt))
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude command failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("claude command failed: %w", err)
	}

	return string(output), nil
}

// CallWithSystemPrompt invokes Claude with system and user prompts.
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	fullPrompt := systemPrompt + "\n\n" + userPrompt
	return c.Call(fullPrompt)
}

// CallJSON invokes Claude and parses the response as JSON.
// Implements: l2/tech-specs.md SEQ-JSON-001 (Response Parsing)
func (c *Client) CallJSON(prompt string, result interface{}) error {
	response, err := c.Call(prompt)
	if err != nil {
		return err
	}

	jsonStr, err := extractJSON(response)
	if err != nil {
		// Conversion fallback: ask Claude to fix the JSON
		// Implements: l2/tech-specs.md SEQ-JSON-001 Step 4
		if c.Verbose {
			fmt.Fprintf(os.Stderr, "[claude] JSON extraction failed, attempting conversion\n")
		}
		conversionPrompt := fmt.Sprintf("Convert this to valid JSON only, no other text:\n\n%s", response)
		response2, err2 := c.Call(conversionPrompt)
		if err2 != nil {
			return fmt.Errorf("JSON extraction failed and conversion failed: %w", err)
		}
		jsonStr, err = extractJSON(response2)
		if err != nil {
			return fmt.Errorf("JSON extraction failed after conversion: %w", err)
		}
	}

	if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
		return fmt.Errorf("JSON parse error: %w", err)
	}

	return nil
}

// CallJSONWithRetry invokes Claude with retry logic for transient failures.
// Implements: l2/tech-specs.md TS-RETRY-001
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error {
	return WithRetry(cfg, func() error {
		return c.CallJSON(prompt, result)
	})
}

// extractJSON extracts JSON from Claude's response.
// Implements: l2/tech-specs.md SEQ-JSON-001 (JSON Extraction Strategy)
func extractJSON(response string) (string, error) {
	// Step 1: Markdown Code Block Extraction
	// Pattern: ```json\n{...}\n```
	if idx := strings.Index(response, "```json"); idx != -1 {
		start := idx + 7 // len("```json")
		// Skip any whitespace/newline after ```json
		for start < len(response) && (response[start] == '\n' || response[start] == '\r' || response[start] == ' ') {
			start++
		}
		end := strings.Index(response[start:], "```")
		if end != -1 {
			jsonStr := strings.TrimSpace(response[start : start+end])
			if json.Valid([]byte(jsonStr)) {
				return jsonStr, nil
			}
			// Try sanitization
			sanitized := sanitizeJSON(jsonStr)
			if json.Valid([]byte(sanitized)) {
				return sanitized, nil
			}
		}
	}

	// Step 2: Raw JSON Extraction
	// Find first '{' or '[' character
	startObj := strings.Index(response, "{")
	startArr := strings.Index(response, "[")

	var start int
	var endChar byte
	if startObj == -1 && startArr == -1 {
		return "", fmt.Errorf("no JSON found in response")
	} else if startObj == -1 {
		start = startArr
		endChar = ']'
	} else if startArr == -1 {
		start = startObj
		endChar = '}'
	} else if startObj < startArr {
		start = startObj
		endChar = '}'
	} else {
		start = startArr
		endChar = ']'
	}

	// Find matching end character
	end := strings.LastIndex(response, string(endChar))
	if end <= start {
		return "", fmt.Errorf("no matching JSON end found")
	}

	jsonStr := response[start : end+1]
	if json.Valid([]byte(jsonStr)) {
		return jsonStr, nil
	}

	// Step 3: JSON Sanitization
	sanitized := sanitizeJSON(jsonStr)
	if json.Valid([]byte(sanitized)) {
		return sanitized, nil
	}

	return "", fmt.Errorf("invalid JSON after sanitization")
}

// sanitizeJSON fixes common LLM JSON issues.
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

		// Remove control characters
		if c < 32 && c != '\t' && c != '\n' && c != '\r' {
			continue
		}

		result.WriteByte(c)
	}

	return result.String()
}

// BuildPrompt injects context into a prompt template.
// Implements: l2/tech-specs.md (Context Injection Pattern)
//
// All prompts use XML-style context markers:
// <context>
// </context>
//
// Context is placed at end of prompt following Anthropic best practices.
func BuildPrompt(template string, documents ...string) string {
	context := strings.Join(documents, "\n\n---\n\n")
	return strings.Replace(template, "</context>", context+"\n</context>", 1)
}
