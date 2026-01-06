# Test Data Directory

This directory contains test fixtures, recorded responses, and golden files for loom-cli tests.

## Structure

```
testdata/
├── fixtures/           # Input files for tests
│   ├── l0-*/          # L0 input documents
│   ├── l1-*/          # L1 documents for derive-l2 tests
│   ├── l2-*/          # L2 documents for derive-l3 tests
│   └── broken-*/      # Invalid documents for validation tests
│
├── responses/          # Recorded LLM responses for integration tests
│   ├── analyze/       # Responses for analyze command
│   ├── derive-l2/     # Responses for derive-l2 command
│   └── derive-l3/     # Responses for derive-l3 command
│
└── golden/            # Expected output files for golden tests
    ├── derive-l2/     # Expected L2 outputs
    └── derive-l3/     # Expected L3 outputs
```

## Usage

### Fixtures

Fixtures are input files used to run tests. They should be complete, valid documents.

```go
func TestDeriveL2(t *testing.T) {
    tmpDir := t.TempDir()
    CopyDir(t, "testdata/fixtures/l1-ecommerce-order", tmpDir)
    // ... run test
}
```

### Responses

Recorded LLM responses for integration tests with mocked Claude client.

```go
func TestWithMockedLLM(t *testing.T) {
    mock := NewMockClient()
    mock.AddResponseFile("derive-tech-specs", "testdata/responses/derive-l2/ecommerce-order/tech-specs.json")
    // ... run test
}
```

### Golden Files

Expected outputs for snapshot testing. Update with:

```bash
LOOM_UPDATE_GOLDEN=1 go test ./cmd/... -run TestDeriveL2_Golden
```

## Recording New Responses

Use RecordingClient to capture LLM responses:

```go
func RecordResponses(t *testing.T) {
    client := claude.NewRecordingClient(claude.NewClient())
    // ... run commands
    client.SaveRecords("testdata/responses/my-test")
}
```

## Conventions

1. **Naming**: Use kebab-case for directories, match fixture names to test names
2. **Completeness**: Fixtures should be complete, valid documents
3. **Minimal**: Keep fixtures as small as possible while still being representative
4. **Versioned**: Commit all testdata to git
