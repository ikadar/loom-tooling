// Package claude provides the Claude CLI integration layer.
//
// Implements: l2/tech-specs.md TS-ARCH-002 (Anti-Corruption Layer via Claude Code CLI)
// See: l1/bounded-context-map.md (ACL pattern)
// See: DEC-L1-001 (Claude CLI integration), DEC-L1-002 (100,000 max output tokens)
package claude

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps the Claude Code CLI for AI-powered derivation.
//
// Implements: l2/internal-api.md, TS-ARCH-002
type Client struct {
	SessionID string // For multi-turn conversations (--resume)
	Verbose   bool   // Debug output to stderr
}

// NewClient creates a new Claude client.
//
// Implements: l2/internal-api.md
func NewClient() *Client {
	return &Client{}
}

// Call sends a prompt to Claude and returns the raw text response.
//
// Implements: TS-ARCH-002
// Uses: claude -p "<prompt>" [--resume <session_id>]
// Sets: CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000 (DEC-L1-002)
func (c *Client) Call(prompt string) (string, error) {
	args := []string{"-p", prompt, "--output-format", "text"}
	if c.SessionID != "" {
		args = append(args, "--resume", c.SessionID)
	}

	cmd := exec.Command("claude", args...)

	// Implements: DEC-L1-002 (100,000 max output tokens)
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[claude] Calling: claude %s\n", strings.Join(args, " "))
	}

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("claude failed: %s: %s", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("claude failed: %w", err)
	}

	return string(output), nil
}

// CallWithSystemPrompt calls Claude with a system prompt and user prompt.
//
// Implements: l2/internal-api.md
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	// Combine system and user prompts - Claude Code CLI doesn't have separate system prompt
	// so we prefix the user prompt with system instructions
	combinedPrompt := fmt.Sprintf("<system>\n%s\n</system>\n\n%s", systemPrompt, userPrompt)
	return c.Call(combinedPrompt)
}

// CallJSON sends a prompt and parses the response as JSON into result.
//
// Implements: SEQ-JSON-001
// Extraction strategy: markdown block → raw JSON → sanitize → conversion fallback
func (c *Client) CallJSON(prompt string, result interface{}) error {
	response, err := c.Call(prompt)
	if err != nil {
		return err
	}

	// Try to extract JSON
	if err := extractJSON(response, result); err != nil {
		// Conversion fallback: ask Claude to fix the JSON
		if c.Verbose {
			fmt.Fprintf(os.Stderr, "[claude] JSON extraction failed, trying conversion fallback\n")
		}

		fixPrompt := fmt.Sprintf("Convert this to valid JSON. Output ONLY the JSON, no explanation:\n\n%s", response)
		fixedResponse, fixErr := c.Call(fixPrompt)
		if fixErr != nil {
			return fmt.Errorf("JSON extraction failed and conversion fallback failed: %w", err)
		}

		if err := extractJSON(fixedResponse, result); err != nil {
			return fmt.Errorf("JSON extraction failed after conversion: %w", err)
		}
	}

	return nil
}

// CallJSONWithRetry sends a prompt with retry logic for transient failures.
//
// Implements: TS-RETRY-001
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error {
	return WithRetry(cfg, func() error {
		return c.CallJSON(prompt, result)
	})
}

// extractJSON attempts to parse JSON from Claude's response.
//
// Implements: SEQ-JSON-001
// Strategy:
//  1. Markdown code block extraction: ```json ... ```
//  2. Raw JSON extraction: find first { or [
//  3. JSON sanitization: fix unescaped newlines
func extractJSON(response string, result interface{}) error {
	// Step 1: Try markdown code block extraction
	jsonStr := extractMarkdownJSON(response)

	// Step 2: If no markdown block, try raw JSON extraction
	if jsonStr == "" {
		jsonStr = extractRawJSON(response)
	}

	if jsonStr == "" {
		return errors.New("no JSON found in response")
	}

	// Step 3: Try parsing as-is first
	if err := json.Unmarshal([]byte(jsonStr), result); err == nil {
		return nil
	}

	// Step 4: Sanitize and retry
	sanitized := sanitizeJSON(jsonStr)
	if err := json.Unmarshal([]byte(sanitized), result); err != nil {
		return fmt.Errorf("JSON parse error: %w", err)
	}

	return nil
}

// extractMarkdownJSON extracts JSON from ```json ... ``` blocks.
//
// Implements: SEQ-JSON-001 Step 1
func extractMarkdownJSON(response string) string {
	const jsonMarker = "```json"
	const endMarker = "```"

	start := strings.Index(response, jsonMarker)
	if start == -1 {
		return ""
	}

	start += len(jsonMarker)
	// Skip any whitespace/newline after marker
	for start < len(response) && (response[start] == '\n' || response[start] == '\r' || response[start] == ' ') {
		start++
	}

	end := strings.Index(response[start:], endMarker)
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(response[start : start+end])
}

// extractRawJSON finds JSON object or array in response.
//
// Implements: SEQ-JSON-001 Step 2
func extractRawJSON(response string) string {
	// Find first { or [
	objStart := strings.Index(response, "{")
	arrStart := strings.Index(response, "[")

	var start int
	var openChar, closeChar byte

	switch {
	case objStart == -1 && arrStart == -1:
		return ""
	case objStart == -1:
		start = arrStart
		openChar, closeChar = '[', ']'
	case arrStart == -1:
		start = objStart
		openChar, closeChar = '{', '}'
	case objStart < arrStart:
		start = objStart
		openChar, closeChar = '{', '}'
	default:
		start = arrStart
		openChar, closeChar = '[', ']'
	}

	// Find matching close bracket, accounting for nesting
	depth := 0
	inString := false
	escaped := false

	for i := start; i < len(response); i++ {
		c := response[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' && inString {
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

		if c == openChar {
			depth++
		} else if c == closeChar {
			depth--
			if depth == 0 {
				return response[start : i+1]
			}
		}
	}

	return ""
}

// sanitizeJSON fixes common LLM JSON output issues.
//
// Implements: SEQ-JSON-001 Step 3
// Fixes: unescaped newlines inside strings, control characters
func sanitizeJSON(jsonStr string) string {
	var result strings.Builder
	result.Grow(len(jsonStr))

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

		// Replace literal newlines inside strings with escaped newlines
		if inString && (c == '\n' || c == '\r') {
			result.WriteString("\\n")
			if c == '\r' && i+1 < len(jsonStr) && jsonStr[i+1] == '\n' {
				i++ // Skip \n after \r
			}
			continue
		}

		// Skip other control characters inside strings
		if inString && c < 32 && c != '\t' {
			continue
		}

		result.WriteByte(c)
	}

	return result.String()
}

// BuildPrompt injects context into a prompt template.
//
// Implements: l2/internal-api.md, TS-ARCH-002
// All prompts use XML-style context markers: <context></context>
// Context is injected just before the closing </context> tag.
func BuildPrompt(template string, documents ...string) string {
	if len(documents) == 0 {
		return template
	}

	context := strings.Join(documents, "\n\n---\n\n")
	return strings.Replace(template, "</context>", context+"\n</context>", 1)
}
