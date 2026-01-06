package derivation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// =============================================================================
// Content Hasher
// =============================================================================

// Hasher computes content hashes for artifacts and files
type Hasher struct {
	// Algorithm is the hash algorithm to use (default: sha256)
	Algorithm string

	// IgnoreWhitespace controls whether to normalize whitespace before hashing
	IgnoreWhitespace bool

	// IgnoreComments controls whether to strip comments before hashing
	IgnoreComments bool

	// SectionHasher enables per-section hashing for LOOM-marked content
	SectionHasher bool
}

// NewHasher creates a new hasher with default settings
func NewHasher() *Hasher {
	return &Hasher{
		Algorithm:        "sha256",
		IgnoreWhitespace: false,
		IgnoreComments:   false,
		SectionHasher:    true,
	}
}

// =============================================================================
// File Hashing
// =============================================================================

// HashFile computes the hash of a file's contents
func (h *Hasher) HashFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return h.HashContent(string(content)), nil
}

// HashContent computes the hash of a string content
func (h *Hasher) HashContent(content string) string {
	// Normalize content if configured
	normalized := content
	if h.IgnoreWhitespace {
		normalized = normalizeWhitespace(normalized)
	}
	if h.IgnoreComments {
		normalized = stripComments(normalized)
	}

	// Compute hash
	hash := sha256.Sum256([]byte(normalized))
	return "sha256:" + hex.EncodeToString(hash[:])
}

// HashFileWithInfo computes hash and returns file info for caching
func (h *Hasher) HashFileWithInfo(path string) (*FileHashInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	hash, err := h.HashFile(path)
	if err != nil {
		return nil, err
	}

	return &FileHashInfo{
		Path:    path,
		Hash:    hash,
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}, nil
}

// =============================================================================
// Incremental Hashing
// =============================================================================

// NeedsRehash checks if a file needs to be rehashed based on cached info
func (h *Hasher) NeedsRehash(path string, cached *FileHashInfo) bool {
	if cached == nil {
		return true
	}

	info, err := os.Stat(path)
	if err != nil {
		return true // File doesn't exist or can't be read
	}

	// Check if file has changed
	if info.ModTime().After(cached.ModTime) || info.Size() != cached.Size {
		return true
	}

	return false
}

// UpdateHashCache updates the hash cache for changed files
func (h *Hasher) UpdateHashCache(state *DerivationState, paths []string) error {
	for _, path := range paths {
		cached := state.GetFileHash(path)
		if h.NeedsRehash(path, cached) {
			info, err := h.HashFileWithInfo(path)
			if err != nil {
				return fmt.Errorf("failed to hash %s: %w", path, err)
			}
			state.SetFileHash(info)
		}
	}
	return nil
}

// =============================================================================
// Section Hashing
// =============================================================================

// SectionHash represents the hash of a specific section within a document
type SectionHash struct {
	// SectionID is the identifier of the section (e.g., "AC-ORD-001")
	SectionID string `json:"section_id"`

	// SectionType is the type of section (e.g., "generated", "manual")
	SectionType string `json:"section_type"`

	// Hash is the content hash of the section
	Hash string `json:"hash"`

	// StartLine is the starting line number of the section
	StartLine int `json:"start_line"`

	// EndLine is the ending line number of the section
	EndLine int `json:"end_line"`
}

// HashSections extracts and hashes individual sections from LOOM-marked content
func (h *Hasher) HashSections(content string) []SectionHash {
	var sections []SectionHash

	lines := strings.Split(content, "\n")

	// Track current section
	var currentSection *SectionHash
	var sectionContent strings.Builder

	// Regex patterns for LOOM markers
	beginPattern := regexp.MustCompile(`<!--\s*LOOM:BEGIN\s+(\w+)(?:\s+id="([^"]+)")?\s*-->`)
	endPattern := regexp.MustCompile(`<!--\s*LOOM:END\s+(\w+)\s*-->`)
	manualPattern := regexp.MustCompile(`<!--\s*LOOM:MANUAL\s+section="([^"]+)"\s*-->`)

	for i, line := range lines {
		lineNum := i + 1

		// Check for LOOM:BEGIN
		if matches := beginPattern.FindStringSubmatch(line); matches != nil {
			if currentSection != nil {
				// Close previous section
				currentSection.Hash = h.HashContent(sectionContent.String())
				sections = append(sections, *currentSection)
			}

			sectionType := matches[1]
			sectionID := ""
			if len(matches) > 2 {
				sectionID = matches[2]
			}

			currentSection = &SectionHash{
				SectionID:   sectionID,
				SectionType: sectionType,
				StartLine:   lineNum,
			}
			sectionContent.Reset()
			continue
		}

		// Check for LOOM:END
		if matches := endPattern.FindStringSubmatch(line); matches != nil {
			if currentSection != nil {
				currentSection.EndLine = lineNum
				currentSection.Hash = h.HashContent(sectionContent.String())
				sections = append(sections, *currentSection)
				currentSection = nil
				sectionContent.Reset()
			}
			continue
		}

		// Check for LOOM:MANUAL (single-line marker for manual sections)
		if matches := manualPattern.FindStringSubmatch(line); matches != nil {
			sections = append(sections, SectionHash{
				SectionID:   matches[1],
				SectionType: "manual",
				StartLine:   lineNum,
				EndLine:     lineNum,
				Hash:        h.HashContent(line),
			})
			continue
		}

		// Accumulate content for current section
		if currentSection != nil {
			if sectionContent.Len() > 0 {
				sectionContent.WriteString("\n")
			}
			sectionContent.WriteString(line)
		}
	}

	// Handle unclosed section at end of file
	if currentSection != nil {
		currentSection.EndLine = len(lines)
		currentSection.Hash = h.HashContent(sectionContent.String())
		sections = append(sections, *currentSection)
	}

	return sections
}

// =============================================================================
// Artifact Hashing
// =============================================================================

// HashArtifact computes the hash of an artifact's content
func (h *Hasher) HashArtifact(artifact *Artifact, projectDir string) (string, error) {
	if artifact.Location.File == "" {
		return "", fmt.Errorf("artifact %s has no file location", artifact.ID)
	}

	// Resolve file path
	filePath := artifact.Location.File
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(projectDir, filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact file: %w", err)
	}

	// If we have line range, extract only that portion
	if artifact.Location.LineStart > 0 {
		lines := strings.Split(string(content), "\n")
		startIdx := artifact.Location.LineStart - 1
		endIdx := len(lines)
		if artifact.Location.LineEnd > 0 && artifact.Location.LineEnd < len(lines) {
			endIdx = artifact.Location.LineEnd
		}

		if startIdx < len(lines) {
			content = []byte(strings.Join(lines[startIdx:endIdx], "\n"))
		}
	}

	return h.HashContent(string(content)), nil
}

// HashArtifacts computes hashes for multiple artifacts
func (h *Hasher) HashArtifacts(artifacts []*Artifact, projectDir string) (map[string]string, error) {
	hashes := make(map[string]string)

	for _, artifact := range artifacts {
		hash, err := h.HashArtifact(artifact, projectDir)
		if err != nil {
			return nil, fmt.Errorf("failed to hash artifact %s: %w", artifact.ID, err)
		}
		hashes[artifact.ID] = hash
	}

	return hashes, nil
}

// =============================================================================
// Upstream Hash Collection
// =============================================================================

// CollectUpstreamHashes gathers current hashes for all upstream dependencies
func (h *Hasher) CollectUpstreamHashes(artifact *Artifact, state *DerivationState, projectDir string) (map[string]string, error) {
	hashes := make(map[string]string)

	for upstreamID := range artifact.Upstream {
		upstream := state.GetArtifact(upstreamID)
		if upstream == nil {
			// Upstream artifact not found - could be external reference
			continue
		}

		// If artifact has a file, hash it; otherwise use stored ContentHash
		if upstream.Location.File != "" {
			hash, err := h.HashArtifact(upstream, projectDir)
			if err != nil {
				return nil, fmt.Errorf("failed to hash upstream %s: %w", upstreamID, err)
			}
			hashes[upstreamID] = hash
		} else if upstream.ContentHash != "" {
			// Use stored content hash for artifacts without files
			hashes[upstreamID] = upstream.ContentHash
		}
	}

	return hashes, nil
}

// =============================================================================
// Directory Hashing
// =============================================================================

// HashDirectory computes hashes for all markdown files in a directory
func (h *Hasher) HashDirectory(dir string) (map[string]string, error) {
	hashes := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		hash, err := h.HashFile(path)
		if err != nil {
			return fmt.Errorf("failed to hash %s: %w", path, err)
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}

		hashes[relPath] = hash
		return nil
	})

	return hashes, err
}

// =============================================================================
// Change Detection
// =============================================================================

// DetectChanges compares current hashes with stored hashes and returns changed files
func (h *Hasher) DetectChanges(currentHashes, storedHashes map[string]string) []string {
	var changed []string

	// Check for changed or new files
	for path, currentHash := range currentHashes {
		storedHash, exists := storedHashes[path]
		if !exists || storedHash != currentHash {
			changed = append(changed, path)
		}
	}

	return changed
}

// DetectDeleted finds files that were in stored hashes but not in current
func (h *Hasher) DetectDeleted(currentHashes, storedHashes map[string]string) []string {
	var deleted []string

	for path := range storedHashes {
		if _, exists := currentHashes[path]; !exists {
			deleted = append(deleted, path)
		}
	}

	return deleted
}

// =============================================================================
// Helper Functions
// =============================================================================

// normalizeWhitespace collapses multiple whitespace into single space
func normalizeWhitespace(content string) string {
	// Replace multiple spaces/tabs with single space
	re := regexp.MustCompile(`[ \t]+`)
	content = re.ReplaceAllString(content, " ")

	// Replace multiple newlines with single newline
	re = regexp.MustCompile(`\n{2,}`)
	content = re.ReplaceAllString(content, "\n")

	return strings.TrimSpace(content)
}

// stripComments removes HTML comments from content (but preserves LOOM markers)
func stripComments(content string) string {
	// Find all HTML comments
	re := regexp.MustCompile(`<!--[^>]*-->`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// Preserve LOOM markers
		if strings.Contains(match, "LOOM:") {
			return match
		}
		return ""
	})
}

// StreamHash computes hash for large files using streaming
func (h *Hasher) StreamHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return "sha256:" + hex.EncodeToString(hasher.Sum(nil)), nil
}

// HashWithTimestamp includes modification time in hash for cache invalidation
func (h *Hasher) HashWithTimestamp(path string) (string, time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", time.Time{}, err
	}

	hash, err := h.HashFile(path)
	if err != nil {
		return "", time.Time{}, err
	}

	return hash, info.ModTime(), nil
}
