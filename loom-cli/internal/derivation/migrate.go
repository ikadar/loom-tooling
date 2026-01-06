package derivation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// =============================================================================
// Migration Types
// =============================================================================

// Migrator handles migration of existing projects to LOOM-marked format
type Migrator struct {
	// Parser is used for content parsing
	Parser *Parser

	// Hasher is used for content hashing
	Hasher *Hasher

	// DryRun if true, doesn't modify files
	DryRun bool

	// Verbose enables detailed logging
	Verbose bool

	// BackupDir is where to store backups (empty = no backup)
	BackupDir string
}

// MigrationResult holds the results of a migration operation
type MigrationResult struct {
	// ProjectDir is the migrated project directory
	ProjectDir string `json:"project_dir"`

	// MigratedFiles lists all migrated files
	MigratedFiles []MigratedFile `json:"migrated_files"`

	// DiscoveredArtifacts lists all discovered artifacts
	DiscoveredArtifacts []*Artifact `json:"discovered_artifacts"`

	// CreatedState is the newly created state (if not dry run)
	CreatedState *DerivationState `json:"-"`

	// Errors lists any migration errors
	Errors []MigrationError `json:"errors,omitempty"`

	// Warnings lists migration warnings
	Warnings []string `json:"warnings,omitempty"`

	// Statistics contains migration statistics
	Statistics MigrationStats `json:"statistics"`
}

// MigratedFile describes a file that was migrated
type MigratedFile struct {
	// Path is the file path
	Path string `json:"path"`

	// ArtifactCount is the number of artifacts found
	ArtifactCount int `json:"artifact_count"`

	// MarkersAdded is the number of LOOM markers added
	MarkersAdded int `json:"markers_added"`

	// BackupPath is the backup file path (if backed up)
	BackupPath string `json:"backup_path,omitempty"`

	// Skipped indicates if the file was skipped
	Skipped bool `json:"skipped,omitempty"`

	// SkipReason explains why the file was skipped
	SkipReason string `json:"skip_reason,omitempty"`
}

// MigrationError describes an error during migration
type MigrationError struct {
	// File is the file where the error occurred
	File string `json:"file"`

	// Message is the error message
	Message string `json:"message"`

	// Recoverable indicates if the error is recoverable
	Recoverable bool `json:"recoverable"`
}

// MigrationStats holds migration statistics
type MigrationStats struct {
	// FilesScanned is the total number of files scanned
	FilesScanned int `json:"files_scanned"`

	// FilesMigrated is the number of files that were modified
	FilesMigrated int `json:"files_migrated"`

	// FilesSkipped is the number of files skipped
	FilesSkipped int `json:"files_skipped"`

	// ArtifactsFound is the total number of artifacts found
	ArtifactsFound int `json:"artifacts_found"`

	// MarkersAdded is the total number of markers added
	MarkersAdded int `json:"markers_added"`

	// Duration is the total migration duration
	Duration time.Duration `json:"duration"`
}

// =============================================================================
// Constructor
// =============================================================================

// NewMigrator creates a new migrator with default settings
func NewMigrator() *Migrator {
	return &Migrator{
		Parser:    NewParser(),
		Hasher:    NewHasher(),
		DryRun:    false,
		Verbose:   false,
		BackupDir: "",
	}
}

// =============================================================================
// Migration Entry Points
// =============================================================================

// MigrateProject migrates an entire project to LOOM-marked format
func (m *Migrator) MigrateProject(projectDir string) (*MigrationResult, error) {
	startTime := time.Now()

	result := &MigrationResult{
		ProjectDir:          projectDir,
		MigratedFiles:       make([]MigratedFile, 0),
		DiscoveredArtifacts: make([]*Artifact, 0),
		Errors:              make([]MigrationError, 0),
		Warnings:            make([]string, 0),
	}

	// Find all spec directories
	specDirs := m.findSpecDirectories(projectDir)
	if len(specDirs) == 0 {
		result.Warnings = append(result.Warnings,
			"No specification directories found (l0/, l1/, l2/, l3/, specs/)")
	}

	// Process each spec directory
	for _, dir := range specDirs {
		if err := m.migrateDirectory(dir, result); err != nil {
			result.Errors = append(result.Errors, MigrationError{
				File:        dir,
				Message:     err.Error(),
				Recoverable: true,
			})
		}
	}

	// Create state if not dry run
	if !m.DryRun && len(result.DiscoveredArtifacts) > 0 {
		state, err := m.createState(projectDir, result.DiscoveredArtifacts)
		if err != nil {
			result.Errors = append(result.Errors, MigrationError{
				File:        ".loom/derivation-state.json",
				Message:     err.Error(),
				Recoverable: false,
			})
		} else {
			result.CreatedState = state
		}
	}

	// Calculate statistics
	result.Statistics.Duration = time.Since(startTime)
	for _, f := range result.MigratedFiles {
		result.Statistics.FilesScanned++
		if f.Skipped {
			result.Statistics.FilesSkipped++
		} else if f.MarkersAdded > 0 {
			result.Statistics.FilesMigrated++
		}
		result.Statistics.ArtifactsFound += f.ArtifactCount
		result.Statistics.MarkersAdded += f.MarkersAdded
	}

	return result, nil
}

// MigrateFile migrates a single file to LOOM-marked format
func (m *Migrator) MigrateFile(filePath string) (*MigratedFile, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Check if already has LOOM markers
	if strings.Contains(string(content), "LOOM:BEGIN") {
		return &MigratedFile{
			Path:       filePath,
			Skipped:    true,
			SkipReason: "Already has LOOM markers",
		}, nil
	}

	// Parse content to find artifacts
	doc := m.Parser.ParseContent(string(content), filePath)

	// Find artifact sections and add markers
	migratedContent, markersAdded := m.addMarkers(string(content), filePath)

	result := &MigratedFile{
		Path:          filePath,
		ArtifactCount: len(doc.Artifacts),
		MarkersAdded:  markersAdded,
	}

	// Write migrated content
	if !m.DryRun && markersAdded > 0 {
		// Backup if configured
		if m.BackupDir != "" {
			backupPath, err := m.backupFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to backup: %w", err)
			}
			result.BackupPath = backupPath
		}

		// Write modified content
		if err := os.WriteFile(filePath, []byte(migratedContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
	}

	return result, nil
}

// =============================================================================
// Internal Methods
// =============================================================================

// findSpecDirectories finds specification directories in a project
func (m *Migrator) findSpecDirectories(projectDir string) []string {
	var dirs []string

	candidates := []string{
		"l0", "l1", "l2", "l3",
		"specs/l0", "specs/l1", "specs/l2", "specs/l3",
		"spec/l0", "spec/l1", "spec/l2", "spec/l3",
		"docs/l0", "docs/l1", "docs/l2", "docs/l3",
	}

	for _, candidate := range candidates {
		path := filepath.Join(projectDir, candidate)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			dirs = append(dirs, path)
		}
	}

	return dirs
}

// migrateDirectory migrates all markdown files in a directory
func (m *Migrator) migrateDirectory(dir string, result *MigrationResult) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		migrated, err := m.MigrateFile(path)
		if err != nil {
			result.Errors = append(result.Errors, MigrationError{
				File:        path,
				Message:     err.Error(),
				Recoverable: true,
			})
			return nil // Continue with other files
		}

		result.MigratedFiles = append(result.MigratedFiles, *migrated)

		// Extract artifacts from migrated file
		if !migrated.Skipped {
			doc, err := m.Parser.ParseFile(path)
			if err == nil {
				result.DiscoveredArtifacts = append(result.DiscoveredArtifacts, doc.Artifacts...)
			}
		}

		return nil
	})
}

// addMarkers adds LOOM markers to content
func (m *Migrator) addMarkers(content, filePath string) (string, int) {
	lines := strings.Split(content, "\n")
	var result []string
	markersAdded := 0

	// Detect artifact sections based on headings
	layer := detectLayerFromPath(filePath)
	sections := m.detectSections(lines, layer)

	// Track which lines have markers
	markerLines := make(map[int]string)
	endMarkerLines := make(map[int]string)

	for _, section := range sections {
		if section.ID != "" {
			// Add begin marker before section start
			markerLines[section.StartLine] = fmt.Sprintf(
				"<!-- LOOM:BEGIN generated id=\"%s\" type=\"%s\" -->",
				section.ID, section.Type)
			// Add end marker after section end
			endMarkerLines[section.EndLine] = fmt.Sprintf(
				"<!-- LOOM:END generated -->")
			markersAdded += 2
		}
	}

	// Build result with markers
	for i, line := range lines {
		lineNum := i + 1

		// Add begin marker if needed
		if marker, ok := markerLines[lineNum]; ok {
			result = append(result, marker)
		}

		result = append(result, line)

		// Add end marker if needed
		if marker, ok := endMarkerLines[lineNum]; ok {
			result = append(result, marker)
		}
	}

	return strings.Join(result, "\n"), markersAdded
}

// detectSections finds artifact sections in lines
func (m *Migrator) detectSections(lines []string, layer string) []detectedSection {
	var sections []detectedSection

	// Heading patterns for different artifact types
	headingPatterns := []struct {
		pattern      *regexp.Regexp
		artifactType string
	}{
		{regexp.MustCompile(`^##\s+(AC-[A-Z]+-\d{3})`), "acceptance_criteria"},
		{regexp.MustCompile(`^##\s+(BR-[A-Z]+-\d{3})`), "business_rule"},
		{regexp.MustCompile(`^##\s+(ENT-[A-Z]+)`), "entity"},
		{regexp.MustCompile(`^##\s+(VO-[A-Z]+)`), "value_object"},
		{regexp.MustCompile(`^##\s+(BC-[A-Z]+)`), "bounded_context"},
		{regexp.MustCompile(`^##\s+(TS-[A-Z]+-\d{3})`), "tech_spec"},
		{regexp.MustCompile(`^##\s+(IC-[A-Z]+-\d{3})`), "interface_operation"},
		{regexp.MustCompile(`^##\s+(AGG-[A-Z]+-\d{3})`), "aggregate_design"},
		{regexp.MustCompile(`^##\s+(SEQ-[A-Z]+-\d{3})`), "sequence"},
		{regexp.MustCompile(`^##\s+(TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2})`), "test_case"},
		{regexp.MustCompile(`^##\s+(API-[A-Z]+-\d{3})`), "api_endpoint"},
		{regexp.MustCompile(`^##\s+(TKT-[A-Z]+-\d{3})`), "ticket"},
	}

	var currentSection *detectedSection

	for i, line := range lines {
		lineNum := i + 1

		// Check for artifact heading
		for _, hp := range headingPatterns {
			if matches := hp.pattern.FindStringSubmatch(line); matches != nil {
				// Close previous section
				if currentSection != nil {
					currentSection.EndLine = lineNum - 1
					sections = append(sections, *currentSection)
				}

				// Start new section
				currentSection = &detectedSection{
					ID:        matches[1],
					Type:      hp.artifactType,
					StartLine: lineNum,
				}
				break
			}
		}

		// Check for next heading (closes current section)
		if currentSection != nil && strings.HasPrefix(line, "## ") {
			isArtifactHeading := false
			for _, hp := range headingPatterns {
				if hp.pattern.MatchString(line) {
					isArtifactHeading = true
					break
				}
			}
			if !isArtifactHeading {
				// Non-artifact heading - close current section
				currentSection.EndLine = lineNum - 1
				sections = append(sections, *currentSection)
				currentSection = nil
			}
		}
	}

	// Close final section
	if currentSection != nil {
		currentSection.EndLine = len(lines)
		sections = append(sections, *currentSection)
	}

	return sections
}

type detectedSection struct {
	ID        string
	Type      string
	StartLine int
	EndLine   int
}

// backupFile creates a backup of a file
func (m *Migrator) backupFile(filePath string) (string, error) {
	// Create backup directory if needed
	if err := os.MkdirAll(m.BackupDir, 0755); err != nil {
		return "", err
	}

	// Generate backup filename
	base := filepath.Base(filePath)
	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(m.BackupDir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	// Copy file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", err
	}

	return backupPath, nil
}

// createState creates initial derivation state from discovered artifacts
func (m *Migrator) createState(projectDir string, artifacts []*Artifact) (*DerivationState, error) {
	sm := NewStateManager(projectDir)
	state := sm.NewState()

	// Add artifacts to state
	for _, artifact := range artifacts {
		// Compute content hash
		hash, err := m.Hasher.HashArtifact(artifact, projectDir)
		if err == nil {
			artifact.ContentHash = hash
		}

		// Set initial status
		artifact.Status = StatusCurrent
		artifact.DerivedAt = time.Now()

		state.SetArtifact(artifact)
	}

	// Build dependency graph from references
	m.buildDependencies(state, artifacts)

	// Save state
	if err := sm.Save(state); err != nil {
		return nil, err
	}

	return state, nil
}

// buildDependencies builds dependency graph from artifact references
func (m *Migrator) buildDependencies(state *DerivationState, artifacts []*Artifact) {
	// Group artifacts by ID for lookup
	artifactMap := make(map[string]*Artifact)
	for _, artifact := range artifacts {
		artifactMap[artifact.ID] = artifact
	}

	// Parse each artifact file to find references
	for _, artifact := range artifacts {
		if artifact.Location.File == "" {
			continue
		}

		content, err := os.ReadFile(artifact.Location.File)
		if err != nil {
			continue
		}

		// Find all artifact IDs referenced in the content
		referencedIDs := m.Parser.GetArtifactIDs(string(content))

		for _, refID := range referencedIDs {
			if refID == artifact.ID {
				continue // Skip self-references
			}

			// Check if referenced artifact exists
			if ref, ok := artifactMap[refID]; ok {
				// Determine direction based on layers
				refLayer := layerOrder(ref.Layer)
				artLayer := layerOrder(artifact.Layer)

				if refLayer < artLayer {
					// ref is upstream of artifact
					if artifact.Upstream == nil {
						artifact.Upstream = make(map[string]string)
					}
					artifact.Upstream[refID] = ref.ContentHash
					state.DependencyGraph.AddEdge(refID, artifact.ID, EdgeDerives)
				}
			}
		}
	}
}

// =============================================================================
// Migration Validation
// =============================================================================

// ValidateMigration checks if a migration was successful
func (m *Migrator) ValidateMigration(result *MigrationResult) []string {
	var issues []string

	// Check for errors
	for _, err := range result.Errors {
		if !err.Recoverable {
			issues = append(issues, fmt.Sprintf("Critical error: %s - %s", err.File, err.Message))
		}
	}

	// Check artifact count
	if result.Statistics.ArtifactsFound == 0 {
		issues = append(issues, "No artifacts were discovered during migration")
	}

	// Check for orphaned artifacts (no upstream)
	rootCount := 0
	for _, artifact := range result.DiscoveredArtifacts {
		if len(artifact.Upstream) == 0 && artifact.Layer != "l0" {
			rootCount++
		}
	}
	if rootCount > 0 && result.Statistics.ArtifactsFound > rootCount {
		issues = append(issues, fmt.Sprintf("%d artifacts have no upstream dependencies (may need manual linking)", rootCount))
	}

	return issues
}

// =============================================================================
// Migration Report
// =============================================================================

// GenerateMigrationReport creates a human-readable migration report
func (m *Migrator) GenerateMigrationReport(result *MigrationResult) string {
	var sb strings.Builder

	sb.WriteString("# Migration Report\n\n")

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Project**: %s\n", result.ProjectDir))
	sb.WriteString(fmt.Sprintf("- **Duration**: %v\n", result.Statistics.Duration))
	sb.WriteString(fmt.Sprintf("- **Files Scanned**: %d\n", result.Statistics.FilesScanned))
	sb.WriteString(fmt.Sprintf("- **Files Migrated**: %d\n", result.Statistics.FilesMigrated))
	sb.WriteString(fmt.Sprintf("- **Files Skipped**: %d\n", result.Statistics.FilesSkipped))
	sb.WriteString(fmt.Sprintf("- **Artifacts Found**: %d\n", result.Statistics.ArtifactsFound))
	sb.WriteString(fmt.Sprintf("- **Markers Added**: %d\n", result.Statistics.MarkersAdded))
	sb.WriteString("\n")

	// Artifacts by layer
	sb.WriteString("## Artifacts by Layer\n\n")
	layerCounts := make(map[string]int)
	for _, artifact := range result.DiscoveredArtifacts {
		layerCounts[artifact.Layer]++
	}
	layers := []string{"l0", "l1", "l2", "l3"}
	for _, layer := range layers {
		if count, ok := layerCounts[layer]; ok {
			sb.WriteString(fmt.Sprintf("- **%s**: %d artifacts\n", strings.ToUpper(layer), count))
		}
	}
	sb.WriteString("\n")

	// Migrated files
	if len(result.MigratedFiles) > 0 {
		sb.WriteString("## Migrated Files\n\n")
		for _, f := range result.MigratedFiles {
			if f.Skipped {
				sb.WriteString(fmt.Sprintf("- ~~%s~~ (skipped: %s)\n", f.Path, f.SkipReason))
			} else if f.MarkersAdded > 0 {
				sb.WriteString(fmt.Sprintf("- %s (%d artifacts, %d markers)\n",
					f.Path, f.ArtifactCount, f.MarkersAdded))
			}
		}
		sb.WriteString("\n")
	}

	// Errors
	if len(result.Errors) > 0 {
		sb.WriteString("## Errors\n\n")
		for _, err := range result.Errors {
			severity := "Error"
			if err.Recoverable {
				severity = "Warning"
			}
			sb.WriteString(fmt.Sprintf("- **%s**: %s - %s\n", severity, err.File, err.Message))
		}
		sb.WriteString("\n")
	}

	// Warnings
	if len(result.Warnings) > 0 {
		sb.WriteString("## Warnings\n\n")
		for _, warn := range result.Warnings {
			sb.WriteString(fmt.Sprintf("- %s\n", warn))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// =============================================================================
// Artifact Discovery (without markers)
// =============================================================================

// DiscoverArtifacts finds artifacts in content without requiring LOOM markers
func (m *Migrator) DiscoverArtifacts(projectDir string) ([]*Artifact, error) {
	var artifacts []*Artifact

	specDirs := m.findSpecDirectories(projectDir)

	for _, dir := range specDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			// Extract artifact IDs from content
			ids := m.Parser.GetArtifactIDs(string(content))
			layer := detectLayerFromPath(path)

			for _, id := range ids {
				artifact := &Artifact{
					ID:    id,
					Type:  m.Parser.detectArtifactType(id),
					Layer: layer,
					Location: ArtifactLocation{
						File: path,
					},
					Status:   StatusNew,
					Upstream: make(map[string]string),
				}
				artifacts = append(artifacts, artifact)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	// Sort by ID for consistent ordering
	sort.Slice(artifacts, func(i, j int) bool {
		return artifacts[i].ID < artifacts[j].ID
	})

	return artifacts, nil
}
