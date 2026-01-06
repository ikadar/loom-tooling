# Synthetic LLM Responses

This directory contains synthetic (mock) LLM responses for integration testing.

## Purpose

These files simulate Claude API responses without making real API calls.
They allow:
- Fast, deterministic tests
- CI/CD pipeline execution without API keys
- Testing error handling and edge cases

## File Format

Each JSON file represents a typical LLM response for a specific operation:

| File | Command | Description |
|------|---------|-------------|
| `domain-model-response.json` | `derive` | DomainModelDoc structure for domain-model.md |
| `bounded-context-response.json` | `derive` | BoundedContextMap for bounded-context-map.md |
| `derivation-response.json` | `derive` | AC and BR derivation results |
| `analyze-domain-response.json` | `analyze` | Domain model extraction from L0 |
| `analyze-ambiguities-response.json` | `analyze` | Ambiguity detection results |

## Usage with MockClient

```go
import (
    "testing"
    loomtest "github.com/ikadar/loom-cli/internal/testing"
    "github.com/ikadar/loom-cli/internal/claude"
)

func TestDeriveWithMock(t *testing.T) {
    // Load synthetic response
    response := loomtest.LoadJSON(t, "../testdata/responses/domain-model-response.json")

    // Create mock client
    mock := claude.NewMockClient()
    mock.AddContainsResponse("derive domain model", response)

    // Use mock in test...
}
```

## Extending

When adding new responses:
1. Follow the existing JSON schema from the actual code
2. Use realistic IDs following the ID pattern conventions (AC-XXX-NNN, BR-XXX-NNN, etc.)
3. Include both success and error cases if applicable
4. Document the response in this README

## Recording Real Responses

To record actual LLM responses for new scenarios:

```go
func TestRecordResponse(t *testing.T) {
    // Create recording client wrapping real client
    realClient := claude.NewClient()
    recorder := claude.NewRecordingClient(realClient)

    // Make the call
    result, err := someOperation(recorder)

    // Save recorded responses
    for prompt, response := range recorder.Responses() {
        hash := sha256.Sum256([]byte(prompt))
        filename := fmt.Sprintf("response-%x.json", hash[:8])
        loomtest.SaveJSON(t, filename, response)
    }
}
```
