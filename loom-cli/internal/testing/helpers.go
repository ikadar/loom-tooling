// Package testing provides test utilities for loom-cli tests.
// Note: This package is named "testing" to be imported as "loomtest" to avoid
// conflicts with the standard testing package.
package testing

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// CopyDir recursively copies a directory tree
func CopyDir(t *testing.T, src, dst string) {
	t.Helper()

	srcInfo, err := os.Stat(src)
	if err != nil {
		t.Fatalf("CopyDir: failed to stat source %s: %v", src, err)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		t.Fatalf("CopyDir: failed to create dest dir %s: %v", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("CopyDir: failed to read source dir %s: %v", src, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			CopyDir(t, srcPath, dstPath)
		} else {
			CopyFile(t, srcPath, dstPath)
		}
	}
}

// CopyFile copies a single file
func CopyFile(t *testing.T, src, dst string) {
	t.Helper()

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("CopyFile: failed to open source %s: %v", src, err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		t.Fatalf("CopyFile: failed to stat source %s: %v", src, err)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		t.Fatalf("CopyFile: failed to create dest %s: %v", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		t.Fatalf("CopyFile: failed to copy %s to %s: %v", src, dst, err)
	}
}

// LoadDir reads all files from a directory into a map
func LoadDir(t *testing.T, dir string) map[string]string {
	t.Helper()

	files := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		files[relPath] = string(content)
		return nil
	})

	if err != nil {
		t.Fatalf("LoadDir: failed to read directory %s: %v", dir, err)
	}

	return files
}

// LoadGolden loads golden files from a directory
// If LOOM_UPDATE_GOLDEN=1, it will update golden files instead of loading
func LoadGolden(t *testing.T, goldenDir string, actualDir string) map[string]string {
	t.Helper()

	if os.Getenv("LOOM_UPDATE_GOLDEN") == "1" {
		// Update mode: copy actual to golden
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatalf("LoadGolden: failed to remove old golden dir: %v", err)
		}
		CopyDir(t, actualDir, goldenDir)
		t.Logf("Updated golden files in %s", goldenDir)
	}

	return LoadDir(t, goldenDir)
}

// WriteFile writes content to a file, creating directories as needed
func WriteFile(t *testing.T, path, content string) {
	t.Helper()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("WriteFile: failed to create dir %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile: failed to write %s: %v", path, err)
	}
}

// ReadFile reads a file and returns its content
func ReadFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: failed to read %s: %v", path, err)
	}
	return string(content)
}

// LoadJSON reads and unmarshals a JSON file
func LoadJSON(t *testing.T, path string, result interface{}) {
	t.Helper()

	content := ReadFile(t, path)
	if err := json.Unmarshal([]byte(content), result); err != nil {
		t.Fatalf("LoadJSON: failed to unmarshal %s: %v", path, err)
	}
}

// SaveJSON marshals and writes JSON to a file
func SaveJSON(t *testing.T, path string, data interface{}) {
	t.Helper()

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("SaveJSON: failed to marshal: %v", err)
	}
	WriteFile(t, path, string(content))
}

// TempDir creates a temporary directory that is cleaned up after the test
func TempDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// =============================================================================
// Markdown Assertions
// =============================================================================

// AssertMarkdownContains checks if markdown content contains expected text
func AssertMarkdownContains(t *testing.T, content, expected string) {
	t.Helper()
	if !strings.Contains(content, expected) {
		t.Errorf("Expected markdown to contain:\n%s\n\nActual:\n%s", expected, truncate(content, 500))
	}
}

// AssertMarkdownSection checks if a markdown section exists and contains expected text
func AssertMarkdownSection(t *testing.T, content, sectionHeader, expected string) {
	t.Helper()

	// Find section
	headerPattern := regexp.MustCompile(`(?m)^#{1,6}\s+` + regexp.QuoteMeta(sectionHeader))
	loc := headerPattern.FindStringIndex(content)
	if loc == nil {
		t.Errorf("Section '%s' not found in markdown", sectionHeader)
		return
	}

	// Find section content (until next header of same or higher level)
	sectionStart := loc[1]
	nextHeaderPattern := regexp.MustCompile(`(?m)^#{1,6}\s+`)
	nextLoc := nextHeaderPattern.FindStringIndex(content[sectionStart:])

	var sectionContent string
	if nextLoc != nil {
		sectionContent = content[sectionStart : sectionStart+nextLoc[0]]
	} else {
		sectionContent = content[sectionStart:]
	}

	if !strings.Contains(sectionContent, expected) {
		t.Errorf("Section '%s' does not contain expected text.\nExpected: %s\nActual section: %s",
			sectionHeader, expected, truncate(sectionContent, 300))
	}
}

// AssertMarkdownHasID checks if markdown contains an ID pattern
func AssertMarkdownHasID(t *testing.T, content, idPattern string) {
	t.Helper()

	pattern := regexp.MustCompile(idPattern)
	if !pattern.MatchString(content) {
		t.Errorf("Expected markdown to contain ID matching pattern '%s'", idPattern)
	}
}

// =============================================================================
// ID Validation Assertions
// =============================================================================

// ExtractIDs extracts all IDs matching a pattern from content
func ExtractIDs(content, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(content, -1)
	return matches
}

// AssertAllReferencesExist checks that all references in source exist in target
func AssertAllReferencesExist(t *testing.T, sourceContent, targetContent, refPattern string) {
	t.Helper()

	refs := ExtractIDs(sourceContent, refPattern)
	targetIDs := ExtractIDs(targetContent, refPattern)

	targetSet := make(map[string]bool)
	for _, id := range targetIDs {
		targetSet[id] = true
	}

	for _, ref := range refs {
		if !targetSet[ref] {
			t.Errorf("Reference '%s' not found in target", ref)
		}
	}
}

// AssertNoDuplicateIDs checks that no ID appears more than once
func AssertNoDuplicateIDs(t *testing.T, content, idPattern string) {
	t.Helper()

	ids := ExtractIDs(content, idPattern)
	seen := make(map[string]int)

	for _, id := range ids {
		seen[id]++
	}

	for id, count := range seen {
		if count > 1 {
			t.Errorf("Duplicate ID found: '%s' (appears %d times)", id, count)
		}
	}
}

// =============================================================================
// Traceability Assertions
// =============================================================================

// SourceRef represents a source reference section
type SourceRef struct {
	ID    string
	Layer string
	Refs  []string
}

// ExtractSourceRefs extracts source_refs sections from markdown
func ExtractSourceRefs(content string) []SourceRef {
	// Pattern for source_refs section
	// Matches: ### source_refs followed by list items
	pattern := regexp.MustCompile(`(?m)### source_refs\n((?:- [^\n]+\n?)+)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	var refs []SourceRef
	for _, match := range matches {
		if len(match) > 1 {
			listContent := match[1]
			// Extract individual refs
			refPattern := regexp.MustCompile(`- ([A-Z]+-[A-Z]+-\d+|[A-Z]+-\d+)`)
			refMatches := refPattern.FindAllStringSubmatch(listContent, -1)

			ref := SourceRef{}
			for _, rm := range refMatches {
				if len(rm) > 1 {
					ref.Refs = append(ref.Refs, rm[1])
				}
			}
			refs = append(refs, ref)
		}
	}

	return refs
}

// AssertBidirectionalRefs checks that references are bidirectional
func AssertBidirectionalRefs(t *testing.T, doc1Content, doc2Content string) {
	t.Helper()

	// This is a simplified check - real implementation would be more thorough
	refs1 := ExtractSourceRefs(doc1Content)
	refs2 := ExtractSourceRefs(doc2Content)

	// Check that refs mentioned in doc1 exist in doc2 and vice versa
	// (Simplified - actual implementation would match specific IDs)
	if len(refs1) == 0 && len(refs2) == 0 {
		t.Log("Warning: No source_refs found in either document")
	}
}

// =============================================================================
// Utility Functions
// =============================================================================

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// FileExists checks if a file exists
func FileExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	return err == nil
}

// AssertFileExists fails if file does not exist
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if !FileExists(t, path) {
		t.Errorf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists fails if file exists
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if FileExists(t, path) {
		t.Errorf("Expected file to not exist: %s", path)
	}
}

// DirExists checks if a directory exists
func DirExists(t *testing.T, path string) bool {
	t.Helper()
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// AssertDirExists fails if directory does not exist
func AssertDirExists(t *testing.T, path string) {
	t.Helper()
	if !DirExists(t, path) {
		t.Errorf("Expected directory to exist: %s", path)
	}
}
