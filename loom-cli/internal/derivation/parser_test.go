package derivation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParser_ParseContent_LOOMMarkers(t *testing.T) {
	p := NewParser()

	content := `# Acceptance Criteria

<!-- LOOM:BEGIN generated id="AC-ORD-001" type="acceptance_criteria" -->
## AC-ORD-001

Given a valid order
When the customer submits
Then the order is created

<!-- LOOM:END generated -->

<!-- LOOM:BEGIN generated id="AC-ORD-002" type="acceptance_criteria" -->
## AC-ORD-002

Given an invalid order
When submitted
Then an error is shown

<!-- LOOM:END generated -->
`

	doc := p.ParseContent(content, "l1/acceptance-criteria.md")

	t.Run("parses sections", func(t *testing.T) {
		if len(doc.Sections) != 2 {
			t.Errorf("Expected 2 sections, got %d", len(doc.Sections))
		}
	})

	t.Run("extracts section IDs", func(t *testing.T) {
		found := make(map[string]bool)
		for _, s := range doc.Sections {
			found[s.ID] = true
		}

		if !found["AC-ORD-001"] {
			t.Error("Should find AC-ORD-001")
		}
		if !found["AC-ORD-002"] {
			t.Error("Should find AC-ORD-002")
		}
	})

	t.Run("extracts artifacts", func(t *testing.T) {
		if len(doc.Artifacts) != 2 {
			t.Errorf("Expected 2 artifacts, got %d", len(doc.Artifacts))
		}
	})

	t.Run("detects layer", func(t *testing.T) {
		if doc.Layer != "l1" {
			t.Errorf("Expected layer 'l1', got '%s'", doc.Layer)
		}
	})

	t.Run("no errors", func(t *testing.T) {
		if len(doc.Errors) != 0 {
			t.Errorf("Expected no errors, got %v", doc.Errors)
		}
	})
}

func TestParser_ParseContent_ManualSections(t *testing.T) {
	p := NewParser()

	content := `<!-- LOOM:BEGIN generated id="BR-ORD-001" -->
## BR-ORD-001

Business rule content.

<!-- LOOM:MANUAL section="notes" -->
Manual notes here.

<!-- LOOM:MANUAL section="edge_cases" -->
Edge case handling.

<!-- LOOM:END generated -->
`

	doc := p.ParseContent(content, "l1/business-rules.md")

	if len(doc.Sections) < 1 {
		t.Fatalf("Expected at least 1 section, got %d", len(doc.Sections))
	}

	// Find generated section
	var generatedSection *ParsedSection
	for i := range doc.Sections {
		if doc.Sections[i].ID == "BR-ORD-001" {
			generatedSection = &doc.Sections[i]
			break
		}
	}

	if generatedSection == nil {
		t.Fatal("Should find BR-ORD-001 section")
	}

	if len(generatedSection.ManualSections) != 2 {
		t.Errorf("Expected 2 manual sections, got %d", len(generatedSection.ManualSections))
	}
}

func TestParser_ParseContent_Metadata(t *testing.T) {
	p := NewParser()

	content := `<!-- LOOM:META version="1.0" layer="l1" author="test" -->
# Document
Content here.
`

	doc := p.ParseContent(content, "test.md")

	if doc.Metadata["version"] != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", doc.Metadata["version"])
	}

	if doc.Metadata["layer"] != "l1" {
		t.Errorf("Expected layer 'l1', got '%s'", doc.Metadata["layer"])
	}

	if doc.Metadata["author"] != "test" {
		t.Errorf("Expected author 'test', got '%s'", doc.Metadata["author"])
	}
}

func TestParser_ParseContent_UnclosedSection(t *testing.T) {
	p := NewParser()
	p.StrictMode = true

	content := `<!-- LOOM:BEGIN generated id="AC-ORD-001" -->
## AC-ORD-001
Content without end marker.
`

	doc := p.ParseContent(content, "test.md")

	if len(doc.Errors) == 0 {
		t.Error("Expected error for unclosed section")
	}
}

func TestParser_ParseFile(t *testing.T) {
	p := NewParser()
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "l1", "test.md")
	os.MkdirAll(filepath.Dir(testFile), 0755)

	content := `<!-- LOOM:BEGIN generated id="AC-TEST-001" -->
## AC-TEST-001
Test content
<!-- LOOM:END generated -->
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	doc, err := p.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if doc.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, doc.Path)
	}

	if len(doc.Artifacts) != 1 {
		t.Errorf("Expected 1 artifact, got %d", len(doc.Artifacts))
	}
}

func TestParser_DetectArtifactType(t *testing.T) {
	p := NewParser()

	tests := []struct {
		id       string
		expected ArtifactType
	}{
		{"AC-ORD-001", ArtifactAcceptanceCrit},
		{"BR-ORD-001", ArtifactBusinessRule},
		{"ENT-ORDER", ArtifactEntity},
		{"VO-MONEY", ArtifactValueObject},
		{"BC-ORDERING", ArtifactBoundedContext},
		{"TS-ORD-001", ArtifactTechSpec},
		{"IC-ORD-001", ArtifactInterfaceOp},
		{"AGG-ORD-001", ArtifactAggregateDesign},
		{"SEQ-ORD-001", ArtifactSequence},
		{"TC-AC-ORD-001-P01", ArtifactTestCase},
		{"TKT-ORD-001", ArtifactTicket},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			result := p.detectArtifactType(tt.id)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParser_GetArtifactIDs(t *testing.T) {
	p := NewParser()

	content := `# Document

See AC-ORD-001 for details.
This implements BR-ORD-002 and BR-ORD-003.

References: ENT-ORDER, VO-MONEY

Test case TC-AC-ORD-001-P01 verifies this.
`

	ids := p.GetArtifactIDs(content)

	expectedIDs := map[string]bool{
		"AC-ORD-001":        true,
		"BR-ORD-002":        true,
		"BR-ORD-003":        true,
		"ENT-ORDER":         true,
		"VO-MONEY":          true,
		"TC-AC-ORD-001-P01": true,
	}

	if len(ids) != len(expectedIDs) {
		t.Errorf("Expected %d IDs, got %d", len(expectedIDs), len(ids))
	}

	for _, id := range ids {
		if !expectedIDs[id] {
			t.Errorf("Unexpected ID: %s", id)
		}
	}
}

func TestParser_ExtractReferences(t *testing.T) {
	p := NewParser()
	p.ExtractRefs = true

	content := `<!-- LOOM:BEGIN generated id="BR-ORD-001" -->
## BR-ORD-001

This rule implements AC-ORD-001.
See also: AC-ORD-002

References: ENT-ORDER

<!-- LOOM:END generated -->
`

	doc := p.ParseContent(content, "l1/business-rules.md")

	if len(doc.References) == 0 {
		t.Error("Expected to find references")
	}

	// Check for specific reference
	foundACRef := false
	for _, ref := range doc.References {
		if ref.FromID == "BR-ORD-001" && ref.ToID == "AC-ORD-001" {
			foundACRef = true
			if ref.Type != "implements" {
				t.Errorf("Expected type 'implements', got '%s'", ref.Type)
			}
		}
	}

	if !foundACRef {
		t.Error("Should find reference from BR-ORD-001 to AC-ORD-001")
	}
}

func TestParser_ParseDirectory(t *testing.T) {
	p := NewParser()
	tmpDir := t.TempDir()

	// Create test files
	l1Dir := filepath.Join(tmpDir, "l1")
	os.MkdirAll(l1Dir, 0755)

	file1 := filepath.Join(l1Dir, "acceptance-criteria.md")
	content1 := `<!-- LOOM:BEGIN generated id="AC-ORD-001" -->
## AC-ORD-001
Content
<!-- LOOM:END generated -->
`
	os.WriteFile(file1, []byte(content1), 0644)

	file2 := filepath.Join(l1Dir, "business-rules.md")
	content2 := `<!-- LOOM:BEGIN generated id="BR-ORD-001" -->
## BR-ORD-001
Content
<!-- LOOM:END generated -->
`
	os.WriteFile(file2, []byte(content2), 0644)

	docs, err := p.ParseDirectory(l1Dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(docs) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(docs))
	}

	// Check GetAllArtifacts helper
	artifacts := GetAllArtifacts(docs)
	if len(artifacts) != 2 {
		t.Errorf("Expected 2 artifacts, got %d", len(artifacts))
	}
}

func TestParser_ValidateDocument(t *testing.T) {
	p := NewParser()

	t.Run("duplicate IDs", func(t *testing.T) {
		content := `<!-- LOOM:BEGIN generated id="AC-ORD-001" -->
## AC-ORD-001
First
<!-- LOOM:END generated -->

<!-- LOOM:BEGIN generated id="AC-ORD-001" -->
## AC-ORD-001
Duplicate
<!-- LOOM:END generated -->
`

		doc := p.ParseContent(content, "test.md")
		errors := p.ValidateDocument(doc)

		hasDuplicateError := false
		for _, err := range errors {
			if err.Severity == "error" && err.Message != "" {
				hasDuplicateError = true
			}
		}

		if !hasDuplicateError {
			t.Error("Should detect duplicate IDs")
		}
	})
}

func TestParser_IsDefinition(t *testing.T) {
	p := NewParser()

	tests := []struct {
		line     string
		id       string
		expected bool
	}{
		{"## AC-ORD-001 - Title", "AC-ORD-001", true},
		{"See AC-ORD-001 for details", "AC-ORD-001", false},
		{"<!-- LOOM:BEGIN generated id=\"AC-ORD-001\" -->", "AC-ORD-001", true},
		{"References: AC-ORD-001", "AC-ORD-001", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := p.isDefinition(tt.line, tt.id)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBuildReferenceMap(t *testing.T) {
	refs := []Reference{
		{FromID: "BR-001", ToID: "AC-001", Type: "implements"},
		{FromID: "BR-001", ToID: "AC-002", Type: "implements"},
		{FromID: "BR-002", ToID: "AC-001", Type: "implements"},
	}

	refMap := BuildReferenceMap(refs)

	if len(refMap["BR-001"]) != 2 {
		t.Errorf("Expected 2 refs from BR-001, got %d", len(refMap["BR-001"]))
	}

	if len(refMap["BR-002"]) != 1 {
		t.Errorf("Expected 1 ref from BR-002, got %d", len(refMap["BR-002"]))
	}
}
