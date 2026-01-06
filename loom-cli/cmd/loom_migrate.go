package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ikadar/loom-cli/internal/derivation"
)

// MigrateConfig holds configuration for the migrate command
type MigrateConfig struct {
	ProjectDir string
	DryRun     bool
	BackupDir  string
	Verbose    bool
	Force      bool // Force migration even if already has markers
}

func runMigrate() error {
	migrateFlags := flag.NewFlagSet("migrate", flag.ExitOnError)
	projectDir := migrateFlags.String("project-dir", ".", "Project root directory")
	dryRun := migrateFlags.Bool("dry-run", false, "Preview without making changes")
	backupDir := migrateFlags.String("backup-dir", "", "Directory for backups (default: .loom/backups)")
	verbose := migrateFlags.Bool("verbose", false, "Show detailed output")
	force := migrateFlags.Bool("force", false, "Force migration even if markers exist")

	if len(os.Args) > 2 {
		migrateFlags.Parse(os.Args[2:])
	}

	cfg := &MigrateConfig{
		ProjectDir: *projectDir,
		DryRun:     *dryRun,
		BackupDir:  *backupDir,
		Verbose:    *verbose,
		Force:      *force,
	}

	return executeMigrate(cfg)
}

func executeMigrate(cfg *MigrateConfig) error {
	fmt.Println("=== LOOM Migration ===")
	fmt.Println()

	if cfg.DryRun {
		fmt.Println("DRY RUN MODE - no changes will be made")
		fmt.Println()
	}

	// Create migrator
	migrator := derivation.NewMigrator()
	migrator.DryRun = cfg.DryRun
	migrator.Verbose = cfg.Verbose

	if cfg.BackupDir != "" {
		migrator.BackupDir = cfg.BackupDir
	}

	// Run migration
	result, err := migrator.MigrateProject(cfg.ProjectDir)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Print migration report
	printMigrationReport(result, cfg.Verbose)

	// Validate migration
	issues := migrator.ValidateMigration(result)
	if len(issues) > 0 {
		fmt.Println("\n⚠ Validation Issues:")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
	}

	// Initialize state if not dry run
	if !cfg.DryRun && len(result.DiscoveredArtifacts) > 0 {
		fmt.Println("\nInitializing derivation state...")

		sm := derivation.NewStateManager(cfg.ProjectDir)
		state := sm.NewState()

		// Add discovered artifacts to state
		for _, artifact := range result.DiscoveredArtifacts {
			state.SetArtifact(artifact)
		}

		// Build dependency graph
		state.DependencyGraph.BuildFromArtifacts(state.Artifacts)

		// Save state
		if err := sm.Save(state); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}

		fmt.Printf("State saved with %d artifacts\n", len(state.Artifacts))
	}

	return nil
}

func printMigrationReport(result *derivation.MigrationResult, verbose bool) {
	stats := result.Statistics

	fmt.Println("Migration Statistics:")
	fmt.Printf("  Files scanned:   %d\n", stats.FilesScanned)
	fmt.Printf("  Files migrated:  %d\n", stats.FilesMigrated)
	fmt.Printf("  Files skipped:   %d\n", stats.FilesSkipped)
	fmt.Printf("  Artifacts found: %d\n", stats.ArtifactsFound)
	fmt.Printf("  Markers added:   %d\n", stats.MarkersAdded)

	// Artifacts by layer
	if len(result.DiscoveredArtifacts) > 0 {
		layerCounts := make(map[string]int)
		typeCounts := make(map[derivation.ArtifactType]int)

		for _, a := range result.DiscoveredArtifacts {
			layerCounts[a.Layer]++
			typeCounts[a.Type]++
		}

		fmt.Println("\nBy Layer:")
		for _, layer := range []string{"l0", "l1", "l2", "l3"} {
			if count := layerCounts[layer]; count > 0 {
				fmt.Printf("  %s: %d\n", strings.ToUpper(layer), count)
			}
		}

		if verbose {
			fmt.Println("\nBy Type:")
			for artType, count := range typeCounts {
				fmt.Printf("  %s: %d\n", artType, count)
			}
		}
	}

	// Migrated files detail
	if verbose && len(result.MigratedFiles) > 0 {
		fmt.Println("\nMigrated Files:")
		for _, f := range result.MigratedFiles {
			if f.Skipped {
				fmt.Printf("  [SKIP] %s: %s\n", f.Path, f.SkipReason)
			} else {
				fmt.Printf("  [OK] %s: %d artifacts, %d markers\n", f.Path, f.ArtifactCount, f.MarkersAdded)
			}
		}
	}

	// Warnings
	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warning)
		}
	}

	// Errors
	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, e := range result.Errors {
			fmt.Printf("  ✗ %s: %s\n", e.File, e.Message)
		}
	}
}
