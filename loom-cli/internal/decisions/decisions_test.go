package decisions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadFromFile_NonExistent(t *testing.T) {
	ds, err := LoadFromFile("/nonexistent/path/decisions.md")
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}
	if len(ds.Decisions) != 0 {
		t.Errorf("Expected empty decision set, got %d decisions", len(ds.Decisions))
	}
}

func TestLoadFromFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "decisions.md")

	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	ds, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(ds.Decisions) != 0 {
		t.Errorf("Expected empty decision set, got %d decisions", len(ds.Decisions))
	}
}

func TestLoadFromFile_SingleDecision(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "decisions.md")

	content := `# Decisions

### AMB-DEF-001

**Question:** What is the maximum order amount?

**Decision:** $10,000 per order

**Source:** user
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ds, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(ds.Decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(ds.Decisions))
	}

	d := ds.Decisions[0]
	if d.AmbiguityID != "AMB-DEF-001" {
		t.Errorf("Expected ID 'AMB-DEF-001', got '%s'", d.AmbiguityID)
	}
	if d.Question != "What is the maximum order amount?" {
		t.Errorf("Expected question about max order, got '%s'", d.Question)
	}
	if d.Answer != "$10,000 per order" {
		t.Errorf("Expected answer '$10,000 per order', got '%s'", d.Answer)
	}
	if d.Source != "user" {
		t.Errorf("Expected source 'user', got '%s'", d.Source)
	}
}

func TestLoadFromFile_MultipleDecisions(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "decisions.md")

	content := `# Decisions

## Missing Definitions

### AMB-DEF-001

**Question:** What is the maximum order amount?

**Decision:** $10,000 per order

**Source:** user

---

### AMB-DEF-002

**Question:** Should orders include tax?

**Decision:** Tax calculated separately

**Source:** default

---

## Business Rules

### AMB-BRG-001

**Question:** Can orders be cancelled?

**Decision:** Yes, within 24 hours

**Source:** existing
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ds, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(ds.Decisions) != 3 {
		t.Fatalf("Expected 3 decisions, got %d", len(ds.Decisions))
	}

	// Check IDs
	ids := make(map[string]bool)
	for _, d := range ds.Decisions {
		ids[d.AmbiguityID] = true
	}

	expectedIDs := []string{"AMB-DEF-001", "AMB-DEF-002", "AMB-BRG-001"}
	for _, id := range expectedIDs {
		if !ids[id] {
			t.Errorf("Expected to find decision with ID '%s'", id)
		}
	}
}

func TestDecisionSet_HasDecision(t *testing.T) {
	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{AmbiguityID: "AMB-001", Answer: "Yes"},
			{AmbiguityID: "AMB-002", Answer: "No"},
		},
	}

	if !ds.HasDecision("AMB-001") {
		t.Error("Expected HasDecision to return true for AMB-001")
	}

	if !ds.HasDecision("AMB-002") {
		t.Error("Expected HasDecision to return true for AMB-002")
	}

	if ds.HasDecision("AMB-003") {
		t.Error("Expected HasDecision to return false for AMB-003")
	}

	if ds.HasDecision("") {
		t.Error("Expected HasDecision to return false for empty string")
	}
}

func TestDecisionSet_GetDecision(t *testing.T) {
	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{AmbiguityID: "AMB-001", Answer: "Yes", Question: "Question 1"},
			{AmbiguityID: "AMB-002", Answer: "No", Question: "Question 2"},
		},
	}

	// Existing decision
	d := ds.GetDecision("AMB-001")
	if d == nil {
		t.Fatal("Expected to get decision for AMB-001")
	}
	if d.Answer != "Yes" {
		t.Errorf("Expected answer 'Yes', got '%s'", d.Answer)
	}
	if d.Question != "Question 1" {
		t.Errorf("Expected 'Question 1', got '%s'", d.Question)
	}

	// Non-existent decision
	d = ds.GetDecision("AMB-999")
	if d != nil {
		t.Error("Expected nil for non-existent decision")
	}
}

func TestDecisionSet_AddDecision_New(t *testing.T) {
	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{},
	}

	ds.AddDecision(AmbiguityDecision{
		AmbiguityID: "AMB-001",
		Question:    "Test question?",
		Answer:      "Test answer",
		Source:      "user",
	})

	if len(ds.Decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(ds.Decisions))
	}

	if ds.Decisions[0].AmbiguityID != "AMB-001" {
		t.Errorf("Expected ID 'AMB-001', got '%s'", ds.Decisions[0].AmbiguityID)
	}
}

func TestDecisionSet_AddDecision_Update(t *testing.T) {
	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{AmbiguityID: "AMB-001", Answer: "Old answer", Source: "default"},
		},
	}

	ds.AddDecision(AmbiguityDecision{
		AmbiguityID: "AMB-001",
		Answer:      "New answer",
		Source:      "user",
	})

	if len(ds.Decisions) != 1 {
		t.Fatalf("Expected 1 decision after update, got %d", len(ds.Decisions))
	}

	if ds.Decisions[0].Answer != "New answer" {
		t.Errorf("Expected updated answer 'New answer', got '%s'", ds.Decisions[0].Answer)
	}

	if ds.Decisions[0].Source != "user" {
		t.Errorf("Expected updated source 'user', got '%s'", ds.Decisions[0].Source)
	}
}

func TestDecisionSet_WriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "output", "decisions.md")

	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{
				AmbiguityID: "AMB-DEF-001",
				Question:    "What is the max amount?",
				Answer:      "$10,000",
				Source:      "user",
				Category:    "missing_definition",
				Severity:    "critical",
				DecidedAt:   time.Now(),
			},
			{
				AmbiguityID: "AMB-BRG-001",
				Question:    "Can orders be cancelled?",
				Answer:      "Yes, within 24h",
				Source:      "default",
				Category:    "business_rule_gap",
				Severity:    "important",
				DecidedAt:   time.Now(),
			},
		},
	}

	err := ds.WriteToFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("Expected file to be created")
	}

	// Read and verify content
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Check header
	if !strings.Contains(contentStr, "# Ambiguity Decisions") {
		t.Error("Expected document title")
	}

	// Check frontmatter
	if !strings.Contains(contentStr, "---") {
		t.Error("Expected YAML frontmatter delimiters")
	}
	if !strings.Contains(contentStr, "title: \"Ambiguity Decisions\"") {
		t.Error("Expected title in frontmatter")
	}

	// Check category headers
	if !strings.Contains(contentStr, "## Missing Definitions") {
		t.Error("Expected formatted category header for missing_definition")
	}
	if !strings.Contains(contentStr, "## Business Rule Gaps") {
		t.Error("Expected formatted category header for business_rule_gap")
	}

	// Check decision content
	if !strings.Contains(contentStr, "### AMB-DEF-001") {
		t.Error("Expected AMB-DEF-001 ID")
	}
	if !strings.Contains(contentStr, "**Question:** What is the max amount?") {
		t.Error("Expected question")
	}
	if !strings.Contains(contentStr, "**Decision:** $10,000") {
		t.Error("Expected decision/answer")
	}
	if !strings.Contains(contentStr, "**Source:** user") {
		t.Error("Expected source")
	}

	// Check summary table
	if !strings.Contains(contentStr, "## Summary") {
		t.Error("Expected Summary section")
	}
	if !strings.Contains(contentStr, "| ID | Category | Severity | Source |") {
		t.Error("Expected summary table header")
	}

	// Check statistics
	if !strings.Contains(contentStr, "Total decisions: 2") {
		t.Error("Expected total decisions count")
	}
	if !strings.Contains(contentStr, "User: 1") {
		t.Error("Expected user count")
	}
	if !strings.Contains(contentStr, "Default: 1") {
		t.Error("Expected default count")
	}
}

func TestWriteToFile_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// Nested path that doesn't exist
	path := filepath.Join(tmpDir, "a", "b", "c", "decisions.md")

	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{AmbiguityID: "AMB-001", Answer: "Test"},
		},
	}

	err := ds.WriteToFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Expected file to be created in nested directory")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
		{"", 10, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "..."},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestFormatCategory(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"missing_definition", "Missing Definitions"},
		{"unclear_relationship", "Unclear Relationships"},
		{"synonym_resolution", "Synonym Resolution"},
		{"boundary_ambiguity", "Boundary Ambiguities"},
		{"business_rule_gap", "Business Rule Gaps"},
		{"state_lifecycle", "State & Lifecycle"},
		{"unknown_category", "Unknown Category"},
		{"some_other_thing", "Some Other Thing"},
	}

	for _, tt := range tests {
		result := formatCategory(tt.input)
		if result != tt.expected {
			t.Errorf("formatCategory(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestLoadFromFile_MalformedContent(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "decisions.md")

	// Content with missing fields
	content := `# Decisions

### AMB-DEF-001

**Question:** Question without answer

### AMB-DEF-002

**Decision:** Answer without question

**Source:** user
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ds, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only have AMB-DEF-002 since AMB-DEF-001 has no answer
	// and AMB-DEF-002 has an answer
	if len(ds.Decisions) != 1 {
		t.Errorf("Expected 1 valid decision, got %d", len(ds.Decisions))
	}

	if len(ds.Decisions) > 0 && ds.Decisions[0].AmbiguityID != "AMB-DEF-002" {
		t.Errorf("Expected AMB-DEF-002, got %s", ds.Decisions[0].AmbiguityID)
	}
}

func TestLoadFromFile_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "decisions.md")

	// Create original decision set
	original := &DecisionSet{
		Decisions: []AmbiguityDecision{
			{
				AmbiguityID: "AMB-DEF-001",
				Question:    "What is the maximum?",
				Answer:      "100 items",
				Source:      "user",
				Category:    "missing_definition",
				Severity:    "critical",
			},
			{
				AmbiguityID: "AMB-REL-001",
				Question:    "How are orders related to customers?",
				Answer:      "One customer can have many orders",
				Source:      "default",
				Category:    "unclear_relationship",
				Severity:    "important",
			},
		},
	}

	// Write to file
	if err := original.WriteToFile(path); err != nil {
		t.Fatal(err)
	}

	// Read back
	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// Verify count
	if len(loaded.Decisions) != len(original.Decisions) {
		t.Fatalf("Expected %d decisions, got %d", len(original.Decisions), len(loaded.Decisions))
	}

	// Verify each decision is present
	for _, orig := range original.Decisions {
		found := loaded.GetDecision(orig.AmbiguityID)
		if found == nil {
			t.Errorf("Expected to find decision %s after round-trip", orig.AmbiguityID)
			continue
		}

		if found.Question != orig.Question {
			t.Errorf("Question mismatch for %s: got %q, want %q", orig.AmbiguityID, found.Question, orig.Question)
		}

		if found.Answer != orig.Answer {
			t.Errorf("Answer mismatch for %s: got %q, want %q", orig.AmbiguityID, found.Answer, orig.Answer)
		}
	}
}
