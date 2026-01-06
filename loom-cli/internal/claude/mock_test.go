package claude

import (
	"testing"
)

func TestMockClient_ExactMatch(t *testing.T) {
	mock := NewMockClient()
	mock.AddResponse("hello", "world")

	response, err := mock.Call("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != "world" {
		t.Errorf("expected 'world', got '%s'", response)
	}
}

func TestMockClient_PrefixMatch(t *testing.T) {
	mock := NewMockClient()
	mock.AddPrefixResponse("You are an assistant", `{"result": "ok"}`)

	response, err := mock.Call("You are an assistant that helps with testing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != `{"result": "ok"}` {
		t.Errorf("unexpected response: %s", response)
	}
}

func TestMockClient_ContainsMatch(t *testing.T) {
	mock := NewMockClient()
	mock.AddContainsResponse("domain model", `{"entities": []}`)

	response, err := mock.Call("Please analyze the domain model for this application")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != `{"entities": []}` {
		t.Errorf("unexpected response: %s", response)
	}
}

func TestMockClient_HashMatch(t *testing.T) {
	mock := NewMockClient()

	// Get hash of a prompt
	prompt := "This is a very long prompt that would be cumbersome to use as a key"
	hash := hashPrompt(prompt)

	mock.AddHashedResponse(hash, "hashed response")

	response, err := mock.Call(prompt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != "hashed response" {
		t.Errorf("unexpected response: %s", response)
	}
}

func TestMockClient_StrictMode_NoMatch(t *testing.T) {
	mock := NewMockClient()
	mock.SetStrictMode(true)

	_, err := mock.Call("unmatched prompt")
	if err == nil {
		t.Error("expected error for unmatched prompt in strict mode")
	}
}

func TestMockClient_NonStrictMode_DefaultResponse(t *testing.T) {
	mock := NewMockClient()
	mock.SetDefaultResponse("default")
	mock.SetStrictMode(false)

	response, err := mock.Call("unmatched prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != "default" {
		t.Errorf("expected 'default', got '%s'", response)
	}
}

func TestMockClient_CallJSON(t *testing.T) {
	mock := NewMockClient()
	mock.AddResponse("get entities", `{"entities": ["Order", "Customer"]}`)

	var result struct {
		Entities []string `json:"entities"`
	}

	err := mock.CallJSON("get entities", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entities) != 2 {
		t.Errorf("expected 2 entities, got %d", len(result.Entities))
	}
	if result.Entities[0] != "Order" {
		t.Errorf("expected 'Order', got '%s'", result.Entities[0])
	}
}

func TestMockClient_CallJSON_FromCodeBlock(t *testing.T) {
	mock := NewMockClient()
	mock.AddResponse("get entities", "Here's the response:\n\n```json\n{\"entities\": [\"Order\"]}\n```\n\nDone!")

	var result struct {
		Entities []string `json:"entities"`
	}

	err := mock.CallJSON("get entities", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entities) != 1 {
		t.Errorf("expected 1 entity, got %d", len(result.Entities))
	}
}

func TestMockClient_CallLog(t *testing.T) {
	mock := NewMockClient()
	mock.SetStrictMode(false)

	mock.Call("first call")
	mock.Call("second call")
	mock.CallWithSystemPrompt("system", "user")

	if mock.GetCallCount() != 3 {
		t.Errorf("expected 3 calls, got %d", mock.GetCallCount())
	}

	calls := mock.GetCalls()
	if calls[0].Prompt != "first call" {
		t.Errorf("expected 'first call', got '%s'", calls[0].Prompt)
	}
	if calls[2].Method != "CallWithSystemPrompt" {
		t.Errorf("expected 'CallWithSystemPrompt', got '%s'", calls[2].Method)
	}
	if calls[2].SystemPrompt != "system" {
		t.Errorf("expected 'system', got '%s'", calls[2].SystemPrompt)
	}
}

func TestMockClient_Reset(t *testing.T) {
	mock := NewMockClient()
	mock.SetStrictMode(false)

	mock.Call("test")
	if mock.GetCallCount() != 1 {
		t.Errorf("expected 1 call before reset")
	}

	mock.Reset()
	if mock.GetCallCount() != 0 {
		t.Errorf("expected 0 calls after reset")
	}
}

func TestMockClient_MatchPriority(t *testing.T) {
	mock := NewMockClient()
	// Add multiple matching responses
	mock.AddResponse("exact match", "exact response")
	mock.AddPrefixResponse("exact", "prefix response")
	mock.AddContainsResponse("match", "contains response")

	// Exact match should take priority
	response, err := mock.Call("exact match")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != "exact response" {
		t.Errorf("expected exact match to take priority, got '%s'", response)
	}
}

func TestRecordingClient(t *testing.T) {
	// Create a mock as the "real" client
	innerMock := NewMockClient()
	innerMock.AddResponse("test prompt", "test response")

	// Wrap it in a recording client
	recorder := NewRecordingClient(innerMock)

	response, err := recorder.Call("test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != "test response" {
		t.Errorf("unexpected response: %s", response)
	}

	records := recorder.GetRecords()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Prompt != "test prompt" {
		t.Errorf("expected 'test prompt', got '%s'", records[0].Prompt)
	}
	if records[0].Response != "test response" {
		t.Errorf("expected 'test response', got '%s'", records[0].Response)
	}
}
