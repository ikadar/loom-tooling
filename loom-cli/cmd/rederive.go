package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ikadar/loom-cli/internal/derivation"
)

// RederiveConfig holds configuration for the rederive command
type RederiveConfig struct {
	ProjectDir     string
	ArtifactIDs    []string // Specific artifacts to derive
	Layer          string   // Filter by layer
	All            bool     // Derive all stale artifacts
	DryRun         bool     // Preview without executing
	Verbose        bool     // Detailed output
	PreserveManual bool     // Keep manual sections
	Interactive    bool     // Confirm each derivation
}

func runRederive() error {
	rederiveFlags := flag.NewFlagSet("rederive", flag.ExitOnError)
	projectDir := rederiveFlags.String("project-dir", ".", "Project root directory")
	layer := rederiveFlags.String("layer", "", "Only derive artifacts in this layer (l1, l2, l3)")
	all := rederiveFlags.Bool("all", false, "Derive all stale artifacts")
	dryRun := rederiveFlags.Bool("dry-run", false, "Preview without making changes")
	verbose := rederiveFlags.Bool("verbose", false, "Show detailed output")
	preserveManual := rederiveFlags.Bool("preserve-manual", true, "Keep manual sections during re-derivation")
	interactive := rederiveFlags.Bool("interactive", false, "Confirm each derivation")

	if len(os.Args) > 2 {
		rederiveFlags.Parse(os.Args[2:])
	}

	// Remaining args are artifact IDs
	artifactIDs := rederiveFlags.Args()

	cfg := &RederiveConfig{
		ProjectDir:     *projectDir,
		ArtifactIDs:    artifactIDs,
		Layer:          *layer,
		All:            *all,
		DryRun:         *dryRun,
		Verbose:        *verbose,
		PreserveManual: *preserveManual,
		Interactive:    *interactive,
	}

	return executeRederive(cfg)
}

func executeRederive(cfg *RederiveConfig) error {
	// Validate config
	if !cfg.All && len(cfg.ArtifactIDs) == 0 && cfg.Layer == "" {
		return fmt.Errorf("specify artifact IDs, --layer, or --all")
	}

	// Load state
	sm := derivation.NewStateManager(cfg.ProjectDir)
	state, err := sm.Load()
	if err != nil {
		return fmt.Errorf("failed to load state (run 'loom-cli init' first): %w", err)
	}

	// Create executor
	executor := derivation.NewExecutor(state, cfg.ProjectDir)
	executor.DryRun = cfg.DryRun
	executor.Verbose = cfg.Verbose
	executor.PreserveManual = cfg.PreserveManual

	// Set up progress callback
	executor.ProgressCallback = func(event derivation.ProgressEvent) {
		switch event.Type {
		case derivation.ProgressStart:
			fmt.Printf("\n%s\n", event.Message)
			fmt.Println(strings.Repeat("─", 50))
		case derivation.ProgressStep:
			prefix := "  "
			if cfg.DryRun {
				prefix = "  [DRY-RUN] "
			}
			fmt.Printf("%s[%d/%d] %s\n", prefix, event.Current, event.Total, event.Message)
		case derivation.ProgressComplete:
			fmt.Println(strings.Repeat("─", 50))
			fmt.Printf("%s\n\n", event.Message)
		case derivation.ProgressError:
			fmt.Printf("  [ERROR] %s: %v\n", event.ArtifactID, event.Error)
		case derivation.ProgressSkip:
			fmt.Printf("  [SKIP] %s: %s\n", event.ArtifactID, event.Message)
		}
	}

	// Set up deriver function
	// This is a placeholder - real implementation would call AI
	executor.DeriverFunc = createDeriverFunc(cfg)

	// Determine what to derive
	var artifactIDs []string

	if cfg.All {
		// Detect all stale artifacts
		stale, err := executor.Tracker.DetectStaleArtifacts()
		if err != nil {
			return fmt.Errorf("failed to detect stale artifacts: %w", err)
		}
		for _, a := range stale {
			if cfg.Layer == "" || a.Layer == cfg.Layer {
				artifactIDs = append(artifactIDs, a.ID)
			}
		}
	} else if cfg.Layer != "" {
		// Get all stale artifacts in layer
		stale, err := executor.Tracker.DetectStaleArtifacts()
		if err != nil {
			return fmt.Errorf("failed to detect stale artifacts: %w", err)
		}
		for _, a := range stale {
			if a.Layer == cfg.Layer {
				artifactIDs = append(artifactIDs, a.ID)
			}
		}
	} else {
		artifactIDs = cfg.ArtifactIDs
	}

	if len(artifactIDs) == 0 {
		fmt.Println("No stale artifacts to derive.")
		return nil
	}

	// Preview mode
	if cfg.DryRun || cfg.Interactive {
		plan, impact, err := executor.PreviewExecution(artifactIDs)
		if err != nil {
			return fmt.Errorf("failed to create preview: %w", err)
		}

		printPreview(plan, impact)

		if cfg.DryRun {
			return nil
		}

		if cfg.Interactive {
			if !confirmExecution() {
				fmt.Println("Derivation cancelled.")
				return nil
			}
		}
	}

	// Execute derivation
	result, err := executor.Execute(artifactIDs)
	if err != nil {
		return fmt.Errorf("derivation failed: %w", err)
	}

	// Print results
	printResults(result)

	// Save updated state
	if !cfg.DryRun && len(result.Derived) > 0 {
		if err := sm.Save(state); err != nil {
			return fmt.Errorf("failed to save state: %w", err)
		}
		fmt.Println("State saved.")
	}

	// Return error if any derivations failed
	if len(result.Errors) > 0 {
		return fmt.Errorf("%d derivation(s) failed", len(result.Errors))
	}

	return nil
}

func createDeriverFunc(cfg *RederiveConfig) derivation.DeriverFunc {
	// This is a placeholder implementation
	// Real implementation would integrate with Claude API
	return func(artifact *derivation.Artifact, upstreamContent map[string]string, projectDir string) (string, error) {
		// For now, return a placeholder message
		// The actual AI-based derivation would go here
		return fmt.Sprintf(`<!-- LOOM:BEGIN generated id="%s" -->
# %s

**Status:** Re-derived from upstream changes

**Upstream artifacts:**
%s

<!-- LOOM:END generated -->
`,
			artifact.ID,
			artifact.ID,
			formatUpstreamList(upstreamContent),
		), nil
	}
}

func formatUpstreamList(upstream map[string]string) string {
	var sb strings.Builder
	for id := range upstream {
		sb.WriteString(fmt.Sprintf("- %s\n", id))
	}
	return sb.String()
}

func printPreview(plan *derivation.DerivationPlan, impact *derivation.ImpactReport) {
	fmt.Println("\n=== Derivation Preview ===")
	fmt.Println()

	fmt.Printf("Artifacts to derive: %d\n", plan.TotalCount)

	// By layer
	if len(plan.ByLayer) > 0 {
		fmt.Println("\nBy Layer:")
		for layer, count := range plan.ByLayer {
			fmt.Printf("  %s: %d\n", strings.ToUpper(layer), count)
		}
	}

	// Derivation order
	fmt.Println("\nDerivation Order:")
	for _, step := range plan.Artifacts {
		manual := ""
		if step.HasManual {
			manual = " [manual sections will be preserved]"
		}
		fmt.Printf("  %d. %s (%s)%s\n", step.Order, step.ArtifactID, step.Layer, manual)
		if len(step.Upstream) > 0 {
			fmt.Printf("     ← %s\n", strings.Join(step.Upstream, ", "))
		}
	}

	// Impact
	if len(impact.AffectedArtifacts) > 0 {
		fmt.Println("\nDownstream Impact:")
		fmt.Printf("  %d artifacts may be affected\n", len(impact.AffectedArtifacts))
		for layer, ids := range impact.AffectedByLayer {
			fmt.Printf("  %s: %d\n", strings.ToUpper(layer), len(ids))
		}
	}

	// Manual edit warnings
	if len(impact.ManualEditWarnings) > 0 {
		fmt.Println("\n⚠ Manual Edit Warnings:")
		for _, warning := range impact.ManualEditWarnings {
			fmt.Printf("  %s: has manual sections %v\n", warning.ArtifactID, warning.ManualSections)
		}
	}

	// Plan warnings
	if len(plan.Warnings) > 0 {
		fmt.Println("\n⚠ Warnings:")
		for _, warning := range plan.Warnings {
			fmt.Printf("  %s\n", warning)
		}
	}

	fmt.Println()
}

func printResults(result *derivation.ExecutionResult) {
	fmt.Println("\n=== Derivation Results ===")
	fmt.Println()

	fmt.Printf("Duration: %v\n", result.Duration.Round(100*1000000)) // Round to 100ms
	fmt.Printf("Derived:  %d\n", len(result.Derived))
	fmt.Printf("Skipped:  %d\n", len(result.Skipped))
	fmt.Printf("Errors:   %d\n", len(result.Errors))

	// Show derived artifacts
	if len(result.Derived) > 0 {
		fmt.Println("\nDerived:")
		for _, d := range result.Derived {
			fmt.Printf("  ✓ %s (%s) → %s\n", d.ArtifactID, d.Layer, d.OutputFile)
		}
	}

	// Show skipped
	if len(result.Skipped) > 0 {
		fmt.Println("\nSkipped:")
		for _, s := range result.Skipped {
			fmt.Printf("  - %s: %s\n", s.ArtifactID, s.Reason)
		}
	}

	// Show errors
	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, e := range result.Errors {
			fmt.Printf("  ✗ %s: %s\n", e.ArtifactID, e.Error)
		}
	}

	fmt.Println()
}

func confirmExecution() bool {
	fmt.Print("Proceed with derivation? [y/N] ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}
