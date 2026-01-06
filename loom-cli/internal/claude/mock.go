package claude

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// MockClient is a test double for Claude API calls.
// It returns pre-recorded responses based on prompt content.
type MockClient struct {
	// Responses maps prompt identifiers to responses
	// Key can be: exact prompt, prompt hash, or prompt prefix
	Responses map[string]string

	// ResponseFiles maps prompt identifiers to file paths containing responses
	ResponseFiles map[string]string

	// CallLog records all calls made to the mock
	CallLog []MockCall

	// DefaultResponse is returned when no matching response is found
	// If empty, an error is returned instead
	DefaultResponse string

	// StrictMode when true, returns error for unmatched prompts
	// When false, returns DefaultResponse or empty string
	StrictMode bool

	mu sync.Mutex
}

// MockCall records a single call to the mock client
type MockCall struct {
	Method       string // "Call", "CallWithSystemPrompt", "CallJSON"
	Prompt       string
	SystemPrompt string // Only for CallWithSystemPrompt
	PromptHash   string
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		Responses:     make(map[string]string),
		ResponseFiles: make(map[string]string),
		CallLog:       make([]MockCall, 0),
		StrictMode:    true,
	}
}

// AddResponse adds a response for a specific prompt
func (m *MockClient) AddResponse(prompt, response string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Responses[prompt] = response
	return m
}

// AddHashedResponse adds a response for a prompt hash
// Use this when prompts are too long to use as keys
func (m *MockClient) AddHashedResponse(promptHash, response string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Responses["hash:"+promptHash] = response
	return m
}

// AddResponseFile adds a response file for a prompt identifier
func (m *MockClient) AddResponseFile(promptID, filepath string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseFiles[promptID] = filepath
	return m
}

// AddPrefixResponse adds a response for prompts starting with a prefix
func (m *MockClient) AddPrefixResponse(prefix, response string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Responses["prefix:"+prefix] = response
	return m
}

// AddContainsResponse adds a response for prompts containing a substring
func (m *MockClient) AddContainsResponse(substring, response string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Responses["contains:"+substring] = response
	return m
}

// SetDefaultResponse sets the default response for unmatched prompts
func (m *MockClient) SetDefaultResponse(response string) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DefaultResponse = response
	m.StrictMode = false
	return m
}

// SetStrictMode enables or disables strict mode
func (m *MockClient) SetStrictMode(strict bool) *MockClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.StrictMode = strict
	return m
}

// Call implements ClaudeClient.Call
func (m *MockClient) Call(prompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	hash := hashPrompt(prompt)
	m.CallLog = append(m.CallLog, MockCall{
		Method:     "Call",
		Prompt:     prompt,
		PromptHash: hash,
	})

	return m.findResponse(prompt, hash)
}

// CallWithSystemPrompt implements ClaudeClient.CallWithSystemPrompt
func (m *MockClient) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	combinedPrompt := systemPrompt + "\n---\n" + userPrompt
	hash := hashPrompt(combinedPrompt)
	m.CallLog = append(m.CallLog, MockCall{
		Method:       "CallWithSystemPrompt",
		Prompt:       userPrompt,
		SystemPrompt: systemPrompt,
		PromptHash:   hash,
	})

	// Try combined first, then just user prompt
	response, err := m.findResponse(combinedPrompt, hash)
	if err != nil {
		return m.findResponse(userPrompt, hashPrompt(userPrompt))
	}
	return response, nil
}

// CallJSON implements ClaudeClient.CallJSON
func (m *MockClient) CallJSON(prompt string, result interface{}) error {
	m.mu.Lock()
	hash := hashPrompt(prompt)
	m.CallLog = append(m.CallLog, MockCall{
		Method:     "CallJSON",
		Prompt:     prompt,
		PromptHash: hash,
	})
	m.mu.Unlock()

	response, err := m.findResponse(prompt, hash)
	if err != nil {
		return err
	}

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(response), result); err != nil {
		// Try extracting JSON from markdown code block
		return extractJSON(response, result)
	}
	return nil
}

// findResponse looks up a response for the given prompt
func (m *MockClient) findResponse(prompt, hash string) (string, error) {
	// 1. Try exact match
	if response, ok := m.Responses[prompt]; ok {
		return response, nil
	}

	// 2. Try hash match
	if response, ok := m.Responses["hash:"+hash]; ok {
		return response, nil
	}

	// 3. Try prefix matches
	for key, response := range m.Responses {
		if strings.HasPrefix(key, "prefix:") {
			prefix := strings.TrimPrefix(key, "prefix:")
			if strings.HasPrefix(prompt, prefix) {
				return response, nil
			}
		}
	}

	// 4. Try contains matches
	for key, response := range m.Responses {
		if strings.HasPrefix(key, "contains:") {
			substring := strings.TrimPrefix(key, "contains:")
			if strings.Contains(prompt, substring) {
				return response, nil
			}
		}
	}

	// 5. Try response files
	for key, filepath := range m.ResponseFiles {
		if key == prompt || key == hash || key == "hash:"+hash {
			content, err := os.ReadFile(filepath)
			if err != nil {
				return "", fmt.Errorf("failed to read response file %s: %w", filepath, err)
			}
			return string(content), nil
		}
	}

	// 6. Use default or error
	if m.StrictMode {
		return "", fmt.Errorf("no mock response found for prompt (hash: %s, len: %d)\nPrompt preview: %s",
			hash, len(prompt), truncate(prompt, 200))
	}
	return m.DefaultResponse, nil
}

// GetCallCount returns the number of calls made
func (m *MockClient) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.CallLog)
}

// GetCalls returns a copy of the call log
func (m *MockClient) GetCalls() []MockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]MockCall, len(m.CallLog))
	copy(calls, m.CallLog)
	return calls
}

// Reset clears the call log
func (m *MockClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CallLog = make([]MockCall, 0)
}

// hashPrompt creates a short hash of a prompt for identification
func hashPrompt(prompt string) string {
	h := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(h[:8]) // First 8 bytes = 16 hex chars
}

// truncate shortens a string for display
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Ensure MockClient implements ClaudeClient
var _ ClaudeClient = (*MockClient)(nil)

// =============================================================================
// Recording Client - for capturing real responses
// =============================================================================

// RecordingClient wraps a real client and records all responses
type RecordingClient struct {
	Real    ClaudeClient
	Records []RecordedCall
	mu      sync.Mutex
}

// RecordedCall stores a prompt-response pair
type RecordedCall struct {
	Method       string `json:"method"`
	Prompt       string `json:"prompt"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	PromptHash   string `json:"prompt_hash"`
	Response     string `json:"response"`
	Error        string `json:"error,omitempty"`
}

// NewRecordingClient creates a client that records all calls
func NewRecordingClient(real ClaudeClient) *RecordingClient {
	return &RecordingClient{
		Real:    real,
		Records: make([]RecordedCall, 0),
	}
}

// Call implements ClaudeClient.Call with recording
func (r *RecordingClient) Call(prompt string) (string, error) {
	response, err := r.Real.Call(prompt)

	r.mu.Lock()
	defer r.mu.Unlock()

	record := RecordedCall{
		Method:     "Call",
		Prompt:     prompt,
		PromptHash: hashPrompt(prompt),
		Response:   response,
	}
	if err != nil {
		record.Error = err.Error()
	}
	r.Records = append(r.Records, record)

	return response, err
}

// CallWithSystemPrompt implements ClaudeClient.CallWithSystemPrompt with recording
func (r *RecordingClient) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error) {
	response, err := r.Real.CallWithSystemPrompt(systemPrompt, userPrompt)

	r.mu.Lock()
	defer r.mu.Unlock()

	combinedPrompt := systemPrompt + "\n---\n" + userPrompt
	record := RecordedCall{
		Method:       "CallWithSystemPrompt",
		Prompt:       userPrompt,
		SystemPrompt: systemPrompt,
		PromptHash:   hashPrompt(combinedPrompt),
		Response:     response,
	}
	if err != nil {
		record.Error = err.Error()
	}
	r.Records = append(r.Records, record)

	return response, err
}

// CallJSON implements ClaudeClient.CallJSON with recording
func (r *RecordingClient) CallJSON(prompt string, result interface{}) error {
	// Call the real client
	err := r.Real.CallJSON(prompt, result)

	r.mu.Lock()
	defer r.mu.Unlock()

	// Serialize the result for recording
	responseJSON, _ := json.MarshalIndent(result, "", "  ")

	record := RecordedCall{
		Method:     "CallJSON",
		Prompt:     prompt,
		PromptHash: hashPrompt(prompt),
		Response:   string(responseJSON),
	}
	if err != nil {
		record.Error = err.Error()
	}
	r.Records = append(r.Records, record)

	return err
}

// SaveRecords saves all recorded calls to a directory
func (r *RecordingClient) SaveRecords(dir string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for i, record := range r.Records {
		filename := filepath.Join(dir, fmt.Sprintf("%03d-%s-%s.json", i, record.Method, record.PromptHash[:8]))
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

// GetRecords returns a copy of all records
func (r *RecordingClient) GetRecords() []RecordedCall {
	r.mu.Lock()
	defer r.mu.Unlock()
	records := make([]RecordedCall, len(r.Records))
	copy(records, r.Records)
	return records
}

// Ensure RecordingClient implements ClaudeClient
var _ ClaudeClient = (*RecordingClient)(nil)
