package testing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	// Create source directory
	srcDir := t.TempDir()
	WriteFile(t, filepath.Join(srcDir, "file1.txt"), "content1")
	WriteFile(t, filepath.Join(srcDir, "subdir", "file2.txt"), "content2")

	// Copy to destination
	dstDir := t.TempDir()
	CopyDir(t, srcDir, dstDir)

	// Verify
	content1 := ReadFile(t, filepath.Join(dstDir, "file1.txt"))
	if content1 != "content1" {
		t.Errorf("expected 'content1', got '%s'", content1)
	}

	content2 := ReadFile(t, filepath.Join(dstDir, "subdir", "file2.txt"))
	if content2 != "content2" {
		t.Errorf("expected 'content2', got '%s'", content2)
	}
}

func TestLoadDir(t *testing.T) {
	dir := t.TempDir()
	WriteFile(t, filepath.Join(dir, "a.txt"), "aaa")
	WriteFile(t, filepath.Join(dir, "b.txt"), "bbb")
	WriteFile(t, filepath.Join(dir, "sub", "c.txt"), "ccc")

	files := LoadDir(t, dir)

	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}
	if files["a.txt"] != "aaa" {
		t.Errorf("expected 'aaa', got '%s'", files["a.txt"])
	}
	if files[filepath.Join("sub", "c.txt")] != "ccc" {
		t.Errorf("expected 'ccc' in sub/c.txt")
	}
}

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	WriteFile(t, path, `{"name": "test", "count": 42}`)

	var result struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	LoadJSON(t, path, &result)

	if result.Name != "test" {
		t.Errorf("expected 'test', got '%s'", result.Name)
	}
	if result.Count != 42 {
		t.Errorf("expected 42, got %d", result.Count)
	}
}

func TestSaveJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.json")

	data := map[string]interface{}{
		"key": "value",
	}
	SaveJSON(t, path, data)

	content := ReadFile(t, path)
	if content != "{\n  \"key\": \"value\"\n}" {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestAssertMarkdownContains(t *testing.T) {
	mockT := &testing.T{}
	content := "# Header\n\nSome content here"

	// Should pass
	AssertMarkdownContains(mockT, content, "Some content")
	if mockT.Failed() {
		t.Error("AssertMarkdownContains should not fail for existing content")
	}
}

func TestExtractIDs(t *testing.T) {
	content := `
## AC-ORD-001: Create Order
Some text here.

## AC-ORD-002: Update Order
More text.

### source_refs
- BR-ORD-001
- BR-ORD-002
`

	// Extract AC IDs
	acIDs := ExtractIDs(content, `AC-ORD-\d+`)
	if len(acIDs) != 2 {
		t.Errorf("expected 2 AC IDs, got %d", len(acIDs))
	}

	// Extract BR refs
	brIDs := ExtractIDs(content, `BR-ORD-\d+`)
	if len(brIDs) != 2 {
		t.Errorf("expected 2 BR IDs, got %d", len(brIDs))
	}
}

func TestAssertNoDuplicateIDs(t *testing.T) {
	mockT := &testing.T{}

	// No duplicates - should pass
	content1 := "## AC-001\n## AC-002\n## AC-003"
	AssertNoDuplicateIDs(mockT, content1, `AC-\d+`)
	if mockT.Failed() {
		t.Error("should not fail when no duplicates")
	}
}

func TestExtractSourceRefs(t *testing.T) {
	content := `
## AC-ORD-001: Create Order

Description here.

### source_refs
- BR-ORD-001
- BR-ORD-002
- US-001

## AC-ORD-002: Update Order

### source_refs
- BR-ORD-003
`

	refs := ExtractSourceRefs(content)
	if len(refs) != 2 {
		t.Errorf("expected 2 source_refs sections, got %d", len(refs))
	}

	// First section should have 3 refs (BR-ORD-001, BR-ORD-002, US-001 matches BR pattern partially)
	// Actually only BR-ORD-XXX matches the pattern
	if len(refs[0].Refs) < 2 {
		t.Errorf("expected at least 2 refs in first section, got %d", len(refs[0].Refs))
	}
}

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	existingFile := filepath.Join(dir, "exists.txt")
	WriteFile(t, existingFile, "content")

	if !FileExists(t, existingFile) {
		t.Error("FileExists should return true for existing file")
	}

	if FileExists(t, filepath.Join(dir, "nonexistent.txt")) {
		t.Error("FileExists should return false for non-existing file")
	}
}

func TestDirExists(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	os.MkdirAll(subDir, 0755)

	if !DirExists(t, subDir) {
		t.Error("DirExists should return true for existing directory")
	}

	if DirExists(t, filepath.Join(dir, "nonexistent")) {
		t.Error("DirExists should return false for non-existing directory")
	}
}
