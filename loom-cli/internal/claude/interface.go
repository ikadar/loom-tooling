package claude

// ClaudeClient defines the interface for Claude API calls.
// This allows mocking in tests.
type ClaudeClient interface {
	// Call sends a prompt to Claude and returns the response
	Call(prompt string) (string, error)

	// CallWithSystemPrompt calls Claude with an additional system prompt
	CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error)

	// CallJSON sends a prompt expecting JSON response and unmarshals it
	CallJSON(prompt string, result interface{}) error
}

// Ensure Client implements ClaudeClient
var _ ClaudeClient = (*Client)(nil)
