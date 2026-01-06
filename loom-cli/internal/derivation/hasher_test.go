package derivation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHasher_HashContent(t *testing.T) {
	h := NewHasher()

	t.Run("consistent hashing", func(t *testing.T) {
		content := "Hello, World!"
		hash1 := h.HashContent(content)
		hash2 := h.HashContent(content)

		if hash1 != hash2 {
			t.Errorf("Same content should produce same hash")
		}
	})

	t.Run("different content different hash", func(t *testing.T) {
		hash1 := h.HashContent("Hello")
		hash2 := h.HashContent("World")

		if hash1 == hash2 {
			t.Errorf("Different content should produce different hash")
		}
	})

	t.Run("hash format", func(t *testing.T) {
		hash := h.HashContent("test")

		if !strings.HasPrefix(hash, "sha256:") {
			t.Errorf("Hash should have sha256: prefix, got %s", hash)
		}
	})
}

func TestHasher_HashFile(t *testing.T) {
	h := NewHasher()
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.md")
	content := "# Test\n\nThis is a test file."
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("hash existing file", func(t *testing.T) {
		hash, err := h.HashFile(testFile)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !strings.HasPrefix(hash, "sha256:") {
			t.Errorf("Hash should have sha256: prefix")
		}

		// Hash should match content hash
		expectedHash := h.HashContent(content)
		if hash != expectedHash {
			t.Errorf("File hash should match content hash")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := h.HashFile("/nonexistent/file.md")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestHasher_HashFileWithInfo(t *testing.T) {
	h := NewHasher()
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	info, err := h.HashFileWithInfo(testFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if info.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, info.Path)
	}

	if info.Hash == "" {
		t.Error("Hash should not be empty")
	}

	if info.Size != 12 { // "test content" is 12 bytes
		t.Errorf("Expected size 12, got %d", info.Size)
	}

	if info.ModTime.IsZero() {
		t.Error("ModTime should not be zero")
	}
}

func TestHasher_NeedsRehash(t *testing.T) {
	h := NewHasher()
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("original"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("nil cached", func(t *testing.T) {
		if !h.NeedsRehash(testFile, nil) {
			t.Error("Should need rehash when cached is nil")
		}
	})

	t.Run("unchanged file", func(t *testing.T) {
		info, _ := h.HashFileWithInfo(testFile)
		if h.NeedsRehash(testFile, info) {
			t.Error("Should not need rehash for unchanged file")
		}
	})

	t.Run("modified file", func(t *testing.T) {
		info, _ := h.HashFileWithInfo(testFile)

		// Modify file
		if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
			t.Fatalf("Failed to modify file: %v", err)
		}

		if !h.NeedsRehash(testFile, info) {
			t.Error("Should need rehash for modified file")
		}
	})
}

func TestHasher_HashSections(t *testing.T) {
	h := NewHasher()

	content := `# Document

<!-- LOOM:BEGIN generated id="AC-ORD-001" -->
## AC-ORD-001

Acceptance criteria content here.

<!-- LOOM:MANUAL section="notes" -->

<!-- LOOM:END generated -->

<!-- LOOM:BEGIN generated id="BR-ORD-001" -->
## BR-ORD-001

Business rule content.

<!-- LOOM:END generated -->
`

	sections := h.HashSections(content)

	if len(sections) < 2 {
		t.Errorf("Expected at least 2 sections, got %d", len(sections))
	}

	// Check first section
	found := false
	for _, s := range sections {
		if s.SectionID == "AC-ORD-001" {
			found = true
			if s.SectionType != "generated" {
				t.Errorf("Expected type 'generated', got '%s'", s.SectionType)
			}
			if s.Hash == "" {
				t.Error("Section hash should not be empty")
			}
		}
	}

	if !found {
		t.Error("Should find AC-ORD-001 section")
	}
}

func TestHasher_HashDirectory(t *testing.T) {
	h := NewHasher()
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"doc1.md":         "# Document 1",
		"doc2.md":         "# Document 2",
		"subdir/doc3.md":  "# Document 3",
		"other.txt":       "Not a markdown file",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		os.MkdirAll(filepath.Dir(path), 0755)
		os.WriteFile(path, []byte(content), 0644)
	}

	hashes, err := h.HashDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have 3 markdown files
	if len(hashes) != 3 {
		t.Errorf("Expected 3 hashes, got %d", len(hashes))
	}

	// Should not include .txt file
	for path := range hashes {
		if strings.HasSuffix(path, ".txt") {
			t.Error("Should not hash .txt files")
		}
	}
}

func TestHasher_DetectChanges(t *testing.T) {
	h := NewHasher()

	current := map[string]string{
		"file1.md": "sha256:aaa",
		"file2.md": "sha256:bbb",
		"file3.md": "sha256:ccc",
	}

	stored := map[string]string{
		"file1.md": "sha256:aaa", // unchanged
		"file2.md": "sha256:xxx", // changed
		// file3.md is new
	}

	changed := h.DetectChanges(current, stored)

	if len(changed) != 2 {
		t.Errorf("Expected 2 changed files, got %d", len(changed))
	}

	// Should include file2.md (changed) and file3.md (new)
	hasFile2 := false
	hasFile3 := false
	for _, f := range changed {
		if f == "file2.md" {
			hasFile2 = true
		}
		if f == "file3.md" {
			hasFile3 = true
		}
	}

	if !hasFile2 {
		t.Error("Should detect file2.md as changed")
	}
	if !hasFile3 {
		t.Error("Should detect file3.md as new")
	}
}

func TestHasher_DetectDeleted(t *testing.T) {
	h := NewHasher()

	current := map[string]string{
		"file1.md": "sha256:aaa",
	}

	stored := map[string]string{
		"file1.md": "sha256:aaa",
		"file2.md": "sha256:bbb",
		"file3.md": "sha256:ccc",
	}

	deleted := h.DetectDeleted(current, stored)

	if len(deleted) != 2 {
		t.Errorf("Expected 2 deleted files, got %d", len(deleted))
	}
}

func TestHasher_NormalizeWhitespace(t *testing.T) {
	h := NewHasher()
	h.IgnoreWhitespace = true

	content1 := "Hello   World"
	content2 := "Hello World"

	hash1 := h.HashContent(content1)
	hash2 := h.HashContent(content2)

	if hash1 != hash2 {
		t.Error("Hashes should match when ignoring whitespace")
	}
}

func TestHasher_StripComments(t *testing.T) {
	h := NewHasher()
	h.IgnoreComments = true

	content1 := "Hello <!-- comment --> World"
	content2 := "Hello  World"

	hash1 := h.HashContent(content1)
	hash2 := h.HashContent(content2)

	if hash1 != hash2 {
		t.Error("Hashes should match when ignoring comments")
	}

	// LOOM markers should not be stripped
	contentWithLoom := "<!-- LOOM:BEGIN generated -->"
	hash := h.HashContent(contentWithLoom)
	if hash == h.HashContent("") {
		t.Error("LOOM markers should not be stripped")
	}
}
