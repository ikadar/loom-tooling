package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ikadar/loom-cli/internal/derivation"
)

// InitConfig holds configuration for the init command
type InitConfig struct {
	ProjectDir string
	Force      bool
	ScanDocs   bool
}

func runInit() error {
	// Parse init-specific flags
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	projectDir := initFlags.String("project-dir", ".", "Project root directory")
	force := initFlags.Bool("force", false, "Overwrite existing state file")
	scanDocs := initFlags.Bool("scan", false, "Scan existing documents and build initial state")

	// Parse arguments
	if len(os.Args) > 2 {
		initFlags.Parse(os.Args[2:])
	}

	cfg := &InitConfig{
		ProjectDir: *projectDir,
		Force:      *force,
		ScanDocs:   *scanDocs,
	}

	return executeInit(cfg)
}

func executeInit(cfg *InitConfig) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(cfg.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if directory exists
	if info, err := os.Stat(absPath); err != nil || !info.IsDir() {
		return fmt.Errorf("project directory does not exist: %s", absPath)
	}

	// Create state manager
	sm := derivation.NewStateManager(absPath)

	// Check if state already exists
	if _, err := os.Stat(sm.StatePath); err == nil {
		if !cfg.Force {
			return fmt.Errorf("state file already exists at %s (use --force to overwrite)", sm.StatePath)
		}
		fmt.Printf("Overwriting existing state file...\n")
	}

	// Create .loom directory
	if err := os.MkdirAll(sm.LoomDir, 0755); err != nil {
		return fmt.Errorf("failed to create .loom directory: %w", err)
	}

	// Create new state
	state := sm.NewState()

	// If scan mode, scan existing documents
	if cfg.ScanDocs {
		if err := scanExistingDocs(absPath, state); err != nil {
			fmt.Printf("Warning: failed to scan existing documents: %v\n", err)
		}
	}

	// Save state
	if err := sm.Save(state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("Initialized loom project at %s\n", absPath)
	fmt.Printf("State file: %s\n", sm.StatePath)

	if len(state.Artifacts) > 0 {
		fmt.Printf("Discovered %d artifacts\n", len(state.Artifacts))
	}

	return nil
}

// scanExistingDocs scans the project for existing specification documents
// and builds initial artifact state from them
func scanExistingDocs(projectDir string, state *derivation.DerivationState) error {
	// Look for specification directories
	specDirs := []string{"l1", "l2", "l3", "specs/l1", "specs/l2", "specs/l3"}

	for _, dir := range specDirs {
		fullPath := filepath.Join(projectDir, dir)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			if err := scanDirectory(fullPath, state); err != nil {
				return err
			}
		}
	}

	// Rebuild dependency graph from discovered artifacts
	if len(state.Artifacts) > 0 {
		state.DependencyGraph.BuildFromArtifacts(state.Artifacts)
	}

	return nil
}

// scanDirectory scans a directory for markdown files and extracts artifacts
func scanDirectory(dir string, state *derivation.DerivationState) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Scan file for artifacts
		artifacts, err := scanFileForArtifacts(path)
		if err != nil {
			return err
		}

		// Add discovered artifacts to state
		for _, artifact := range artifacts {
			state.SetArtifact(artifact)
		}

		return nil
	})
}

// scanFileForArtifacts scans a markdown file and extracts artifact definitions
// This is a basic implementation that will be enhanced in Phase 1
func scanFileForArtifacts(path string) ([]*derivation.Artifact, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Determine layer from path
	layer := detectLayerFromPath(path)

	// Extract artifact IDs from content
	artifacts := extractArtifactsFromContent(string(content), path, layer)

	return artifacts, nil
}

// detectLayerFromPath determines the layer (l0, l1, l2, l3) from file path
func detectLayerFromPath(path string) string {
	pathNorm := filepath.ToSlash(path)

	if strings.Contains(pathNorm, "/l3/") || strings.Contains(pathNorm, "/l3-") {
		return "l3"
	}
	if strings.Contains(pathNorm, "/l2/") || strings.Contains(pathNorm, "/l2-") {
		return "l2"
	}
	if strings.Contains(pathNorm, "/l1/") || strings.Contains(pathNorm, "/l1-") {
		return "l1"
	}
	if strings.Contains(pathNorm, "/l0/") || strings.Contains(pathNorm, "/l0-") {
		return "l0"
	}

	// Default to l1 for spec files without clear layer
	return "l1"
}

// extractArtifactsFromContent extracts artifact definitions from markdown content
// This is a basic implementation - Phase 1 will add proper parsing
func extractArtifactsFromContent(content, filePath, layer string) []*derivation.Artifact {
	var artifacts []*derivation.Artifact

	// Use the same ID patterns as validation
	idPatterns := map[string]derivation.ArtifactType{
		`AC-[A-Z]+-\d{3}`:                derivation.ArtifactAcceptanceCrit,
		`BR-[A-Z]+-\d{3}`:                derivation.ArtifactBusinessRule,
		`ENT-[A-Z]+`:                     derivation.ArtifactEntity,
		`BC-[A-Z]+`:                      derivation.ArtifactBoundedContext,
		`TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2}`: derivation.ArtifactTestCase,
		`TS-[A-Z]+-\d{3}`:                derivation.ArtifactTechSpec,
		`IC-[A-Z]+-\d{3}`:                derivation.ArtifactInterfaceOp,
		`AGG-[A-Z]+-\d{3}`:               derivation.ArtifactAggregateDesign,
		`SEQ-[A-Z]+-\d{3}`:               derivation.ArtifactSequence,
	}

	// Track found IDs to avoid duplicates
	found := make(map[string]bool)

	for pattern, artType := range idPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		matches := re.FindAllString(content, -1)
		for _, id := range matches {
			if found[id] {
				continue
			}
			found[id] = true

			artifact := &derivation.Artifact{
				ID:    id,
				Type:  artType,
				Layer: layer,
				Location: derivation.ArtifactLocation{
					File: filePath,
				},
				Status:   derivation.StatusCurrent,
				Upstream: make(map[string]string),
			}

			artifacts = append(artifacts, artifact)
		}
	}

	return artifacts
}
