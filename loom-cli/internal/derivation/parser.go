package derivation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// =============================================================================
// LOOM Marker Constants
// =============================================================================

const (
	// MarkerBegin marks the start of a generated section
	MarkerBegin = "<!-- LOOM:BEGIN"

	// MarkerEnd marks the end of a generated section
	MarkerEnd = "<!-- LOOM:END"

	// MarkerManual marks a manual section
	MarkerManual = "<!-- LOOM:MANUAL"

	// MarkerMeta contains artifact metadata
	MarkerMeta = "<!-- LOOM:META"

	// MarkerRef marks a reference to another artifact
	MarkerRef = "<!-- LOOM:REF"
)

// =============================================================================
// Parser Types
// =============================================================================

// Parser extracts artifacts and references from LOOM-marked documents
type Parser struct {
	// StrictMode requires all sections to have proper LOOM markers
	StrictMode bool

	// ExtractRefs enables extraction of references from content
	ExtractRefs bool

	// IDPatterns are regex patterns for recognizing artifact IDs
	IDPatterns map[string]*regexp.Regexp
}

// ParsedDocument represents a fully parsed specification document
type ParsedDocument struct {
	// Path is the file path of the document
	Path string `json:"path"`

	// Layer is the detected layer (l0, l1, l2, l3)
	Layer string `json:"layer"`

	// Sections contains all parsed sections
	Sections []ParsedSection `json:"sections"`

	// Artifacts contains extracted artifact definitions
	Artifacts []*Artifact `json:"artifacts"`

	// References contains all cross-references found
	References []Reference `json:"references"`

	// Metadata contains document-level metadata
	Metadata map[string]string `json:"metadata,omitempty"`

	// Errors contains any parsing errors
	Errors []ParseError `json:"errors,omitempty"`
}

// ParsedSection represents a section within a document
type ParsedSection struct {
	// ID is the section identifier (artifact ID or section name)
	ID string `json:"id"`

	// Type is the section type (generated, manual, mixed)
	Type string `json:"type"`

	// StartLine is the starting line number (1-indexed)
	StartLine int `json:"start_line"`

	// EndLine is the ending line number (1-indexed)
	EndLine int `json:"end_line"`

	// Content is the raw section content
	Content string `json:"content"`

	// ArtifactType is the detected artifact type
	ArtifactType ArtifactType `json:"artifact_type,omitempty"`

	// ManualSections lists manual subsection names within this section
	ManualSections []string `json:"manual_sections,omitempty"`
}

// Reference represents a cross-reference between artifacts
type Reference struct {
	// FromID is the source artifact ID
	FromID string `json:"from_id"`

	// ToID is the target artifact ID
	ToID string `json:"to_id"`

	// Type is the reference type (derives, references, implements, tests)
	Type string `json:"type"`

	// Line is the line number where the reference was found
	Line int `json:"line"`
}

// ParseError represents a parsing error
type ParseError struct {
	// Line is the line number where the error occurred
	Line int `json:"line"`

	// Message is the error description
	Message string `json:"message"`

	// Severity is the error severity (error, warning)
	Severity string `json:"severity"`
}

// =============================================================================
// Parser Constructor
// =============================================================================

// NewParser creates a new parser with default settings
func NewParser() *Parser {
	return &Parser{
		StrictMode:  false,
		ExtractRefs: true,
		IDPatterns:  defaultIDPatterns(),
	}
}

// defaultIDPatterns returns the standard artifact ID patterns
func defaultIDPatterns() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		// L0 patterns
		"US": regexp.MustCompile(`US-[A-Z]+-\d{3}`),

		// L1 patterns
		"AC":  regexp.MustCompile(`AC-[A-Z]+-\d{3}`),
		"BR":  regexp.MustCompile(`BR-[A-Z]+-\d{3}`),
		"ENT": regexp.MustCompile(`ENT-[A-Z]+`),
		"VO":  regexp.MustCompile(`VO-[A-Z]+`),
		"BC":  regexp.MustCompile(`BC-[A-Z]+`),

		// L2 patterns
		"TS":  regexp.MustCompile(`TS-[A-Z]+-\d{3}`),
		"IC":  regexp.MustCompile(`IC-[A-Z]+-\d{3}`),
		"AGG": regexp.MustCompile(`AGG-[A-Z]+-\d{3}`),
		"SEQ": regexp.MustCompile(`SEQ-[A-Z]+-\d{3}`),
		"DT":  regexp.MustCompile(`DT-[A-Z]+-\d{3}`),

		// L3 patterns
		"TC":  regexp.MustCompile(`TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2}`),
		"API": regexp.MustCompile(`API-[A-Z]+-\d{3}`),
		"EVT": regexp.MustCompile(`EVT-[A-Z]+-\d{3}`),
		"CMD": regexp.MustCompile(`CMD-[A-Z]+-\d{3}`),
		"TKT": regexp.MustCompile(`TKT-[A-Z]+-\d{3}`),
	}
}

// =============================================================================
// Document Parsing
// =============================================================================

// ParseFile parses a file and extracts all LOOM-marked content
func (p *Parser) ParseFile(path string) (*ParsedDocument, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	doc := p.ParseContent(string(content), path)
	return doc, nil
}

// ParseContent parses string content and extracts LOOM-marked sections
func (p *Parser) ParseContent(content, path string) *ParsedDocument {
	doc := &ParsedDocument{
		Path:       path,
		Layer:      detectLayerFromPath(path),
		Sections:   make([]ParsedSection, 0),
		Artifacts:  make([]*Artifact, 0),
		References: make([]Reference, 0),
		Metadata:   make(map[string]string),
		Errors:     make([]ParseError, 0),
	}

	lines := strings.Split(content, "\n")

	// Parse LOOM markers
	p.parseMarkers(lines, doc)

	// Extract artifacts from sections
	p.extractArtifacts(doc)

	// Extract references if enabled
	if p.ExtractRefs {
		p.extractReferences(lines, doc)
	}

	return doc
}

// =============================================================================
// Marker Parsing
// =============================================================================

// parseMarkers extracts all LOOM-marked sections from lines
func (p *Parser) parseMarkers(lines []string, doc *ParsedDocument) {
	// Regex patterns for LOOM markers
	beginPattern := regexp.MustCompile(`<!--\s*LOOM:BEGIN\s+(\w+)(?:\s+id="([^"]+)")?(?:\s+type="([^"]+)")?\s*-->`)
	endPattern := regexp.MustCompile(`<!--\s*LOOM:END\s+(\w+)\s*-->`)
	manualPattern := regexp.MustCompile(`<!--\s*LOOM:MANUAL\s+section="([^"]+)"\s*-->`)
	metaPattern := regexp.MustCompile(`<!--\s*LOOM:META\s+([^>]+)\s*-->`)

	// Track section stack (for nested sections)
	type sectionState struct {
		section   *ParsedSection
		content   strings.Builder
		startLine int
	}
	var stack []sectionState

	for i, line := range lines {
		lineNum := i + 1

		// Check for LOOM:META
		if matches := metaPattern.FindStringSubmatch(line); matches != nil {
			p.parseMetadata(matches[1], doc)
			continue
		}

		// Check for LOOM:BEGIN
		if matches := beginPattern.FindStringSubmatch(line); matches != nil {
			sectionType := matches[1]
			sectionID := ""
			artifactType := ""

			if len(matches) > 2 && matches[2] != "" {
				sectionID = matches[2]
			}
			if len(matches) > 3 && matches[3] != "" {
				artifactType = matches[3]
			}

			section := &ParsedSection{
				ID:           sectionID,
				Type:         sectionType,
				StartLine:    lineNum,
				ArtifactType: ArtifactType(artifactType),
			}

			stack = append(stack, sectionState{
				section:   section,
				startLine: lineNum,
			})
			continue
		}

		// Check for LOOM:END
		if matches := endPattern.FindStringSubmatch(line); matches != nil {
			if len(stack) > 0 {
				state := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				state.section.EndLine = lineNum
				state.section.Content = state.content.String()
				doc.Sections = append(doc.Sections, *state.section)
			} else {
				doc.Errors = append(doc.Errors, ParseError{
					Line:     lineNum,
					Message:  "LOOM:END without matching LOOM:BEGIN",
					Severity: "error",
				})
			}
			continue
		}

		// Check for LOOM:MANUAL
		if matches := manualPattern.FindStringSubmatch(line); matches != nil {
			manualSection := ParsedSection{
				ID:        matches[1],
				Type:      "manual",
				StartLine: lineNum,
				EndLine:   lineNum,
				Content:   line,
			}
			doc.Sections = append(doc.Sections, manualSection)

			// Also track in parent section if exists
			if len(stack) > 0 {
				parent := stack[len(stack)-1].section
				parent.ManualSections = append(parent.ManualSections, matches[1])
			}
			continue
		}

		// Accumulate content for current section
		if len(stack) > 0 {
			if stack[len(stack)-1].content.Len() > 0 {
				stack[len(stack)-1].content.WriteString("\n")
			}
			stack[len(stack)-1].content.WriteString(line)
		}
	}

	// Check for unclosed sections
	for _, state := range stack {
		doc.Errors = append(doc.Errors, ParseError{
			Line:     state.startLine,
			Message:  fmt.Sprintf("Unclosed section starting at line %d", state.startLine),
			Severity: "error",
		})
	}
}

// parseMetadata extracts key="value" pairs from metadata
func (p *Parser) parseMetadata(metaStr string, doc *ParsedDocument) {
	// Pattern for key="value" pairs
	pattern := regexp.MustCompile(`(\w+)="([^"]*)"`)
	matches := pattern.FindAllStringSubmatch(metaStr, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			doc.Metadata[match[1]] = match[2]
		}
	}
}

// =============================================================================
// Artifact Extraction
// =============================================================================

// extractArtifacts creates Artifact objects from parsed sections
func (p *Parser) extractArtifacts(doc *ParsedDocument) {
	for _, section := range doc.Sections {
		if section.Type == "manual" {
			continue
		}

		// Try to find artifact ID in section
		id := section.ID
		if id == "" {
			// Try to extract ID from content
			id = p.extractIDFromContent(section.Content)
		}

		if id == "" {
			continue
		}

		// Determine artifact type
		artType := section.ArtifactType
		if artType == "" {
			artType = p.detectArtifactType(id)
		}

		artifact := &Artifact{
			ID:    id,
			Type:  artType,
			Layer: doc.Layer,
			Location: ArtifactLocation{
				File:      doc.Path,
				LineStart: section.StartLine,
				LineEnd:   section.EndLine,
			},
			Status:         StatusCurrent,
			ManualSections: section.ManualSections,
			Upstream:       make(map[string]string),
		}

		doc.Artifacts = append(doc.Artifacts, artifact)
	}
}

// extractIDFromContent tries to extract an artifact ID from content
func (p *Parser) extractIDFromContent(content string) string {
	// Try each ID pattern
	for _, pattern := range p.IDPatterns {
		if match := pattern.FindString(content); match != "" {
			return match
		}
	}
	return ""
}

// detectArtifactType determines artifact type from ID prefix
func (p *Parser) detectArtifactType(id string) ArtifactType {
	// Map of ID prefixes to artifact types
	prefixMap := map[string]ArtifactType{
		"US":  ArtifactUserStory,
		"AC":  ArtifactAcceptanceCrit,
		"BR":  ArtifactBusinessRule,
		"ENT": ArtifactEntity,
		"VO":  ArtifactValueObject,
		"BC":  ArtifactBoundedContext,
		"TS":  ArtifactTechSpec,
		"IC":  ArtifactInterfaceOp,
		"AGG": ArtifactAggregateDesign,
		"SEQ": ArtifactSequence,
		"DT":  ArtifactDataTable,
		"TC":  ArtifactTestCase,
		"API": ArtifactAPIEndpoint,
		"EVT": ArtifactType("event"),
		"CMD": ArtifactType("command"),
		"TKT": ArtifactTicket,
	}

	// Find matching prefix
	for prefix, artType := range prefixMap {
		if strings.HasPrefix(id, prefix+"-") || strings.HasPrefix(id, prefix) {
			return artType
		}
	}

	return ArtifactType("unknown")
}

// =============================================================================
// Reference Extraction
// =============================================================================

// extractReferences finds all cross-references in the document
func (p *Parser) extractReferences(lines []string, doc *ParsedDocument) {
	// Build map of artifact IDs in this document
	localIDs := make(map[string]bool)
	for _, artifact := range doc.Artifacts {
		localIDs[artifact.ID] = true
	}

	// Pattern to match references
	// Looks for: "See AC-ORD-001", "References: BR-ORD-002", "Implements AC-ORD-001"
	refPatterns := []struct {
		pattern *regexp.Regexp
		refType string
	}{
		{regexp.MustCompile(`(?i)(?:see|ref|references?)\s*:?\s*`), "references"},
		{regexp.MustCompile(`(?i)(?:implements?|realizes?)\s*:?\s*`), "implements"},
		{regexp.MustCompile(`(?i)(?:tests?|verifies?)\s*:?\s*`), "tests"},
		{regexp.MustCompile(`(?i)(?:derives?\s+from|derived\s+from)\s*:?\s*`), "derives"},
	}

	// Find all artifact IDs mentioned in each line
	for i, line := range lines {
		lineNum := i + 1

		// Find all artifact IDs in this line
		for prefix, pattern := range p.IDPatterns {
			matches := pattern.FindAllString(line, -1)
			for _, targetID := range matches {
				// Skip if it's a local ID definition (not a reference)
				if p.isDefinition(line, targetID) {
					continue
				}

				// Determine reference type
				refType := "references"
				for _, rp := range refPatterns {
					if rp.pattern.MatchString(line) {
						refType = rp.refType
						break
					}
				}

				// Find source artifact (the one making the reference)
				sourceID := p.findSourceArtifact(lineNum, doc)

				if sourceID != "" && sourceID != targetID {
					doc.References = append(doc.References, Reference{
						FromID: sourceID,
						ToID:   targetID,
						Type:   refType,
						Line:   lineNum,
					})
				}

				// Suppress unused variable warning
				_ = prefix
			}
		}
	}
}

// isDefinition checks if the ID appears as a definition (heading or marker)
func (p *Parser) isDefinition(line, id string) bool {
	// Check if line is a heading with the ID
	if strings.HasPrefix(strings.TrimSpace(line), "#") && strings.Contains(line, id) {
		return true
	}

	// Check if it's in a LOOM marker
	if strings.Contains(line, "LOOM:") && strings.Contains(line, id) {
		return true
	}

	return false
}

// findSourceArtifact finds which artifact contains the given line
func (p *Parser) findSourceArtifact(lineNum int, doc *ParsedDocument) string {
	for _, section := range doc.Sections {
		if section.Type != "manual" && lineNum >= section.StartLine && lineNum <= section.EndLine {
			return section.ID
		}
	}
	return ""
}

// =============================================================================
// Directory Parsing
// =============================================================================

// ParseDirectory parses all markdown files in a directory
func (p *Parser) ParseDirectory(dir string) ([]*ParsedDocument, error) {
	var docs []*ParsedDocument

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		doc, err := p.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		docs = append(docs, doc)
		return nil
	})

	return docs, err
}

// =============================================================================
// Utility Functions
// =============================================================================

// GetAllArtifacts extracts all artifacts from multiple documents
func GetAllArtifacts(docs []*ParsedDocument) []*Artifact {
	var artifacts []*Artifact
	for _, doc := range docs {
		artifacts = append(artifacts, doc.Artifacts...)
	}
	return artifacts
}

// GetAllReferences extracts all references from multiple documents
func GetAllReferences(docs []*ParsedDocument) []Reference {
	var refs []Reference
	for _, doc := range docs {
		refs = append(refs, doc.References...)
	}
	return refs
}

// BuildReferenceMap creates a map of artifact ID to its references
func BuildReferenceMap(refs []Reference) map[string][]Reference {
	refMap := make(map[string][]Reference)
	for _, ref := range refs {
		refMap[ref.FromID] = append(refMap[ref.FromID], ref)
	}
	return refMap
}

// GetArtifactIDs extracts all unique artifact IDs from content
func (p *Parser) GetArtifactIDs(content string) []string {
	ids := make(map[string]bool)

	for _, pattern := range p.IDPatterns {
		matches := pattern.FindAllString(content, -1)
		for _, id := range matches {
			ids[id] = true
		}
	}

	// Convert to sorted slice
	result := make([]string, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	sort.Strings(result)

	return result
}

// ValidateDocument checks a parsed document for errors
func (p *Parser) ValidateDocument(doc *ParsedDocument) []ParseError {
	var errors []ParseError

	// Check for unclosed sections (already done in parsing)
	errors = append(errors, doc.Errors...)

	// Check for duplicate artifact IDs
	idCounts := make(map[string]int)
	for _, artifact := range doc.Artifacts {
		idCounts[artifact.ID]++
	}
	for id, count := range idCounts {
		if count > 1 {
			errors = append(errors, ParseError{
				Line:     0,
				Message:  fmt.Sprintf("Duplicate artifact ID: %s (appears %d times)", id, count),
				Severity: "error",
			})
		}
	}

	// Check for invalid references
	localIDs := make(map[string]bool)
	for _, artifact := range doc.Artifacts {
		localIDs[artifact.ID] = true
	}

	return errors
}
