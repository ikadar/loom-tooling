package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ikadar/loom-cli/internal/derivation"
)

// StatusConfig holds configuration for the status command
type StatusConfig struct {
	ProjectDir string
	Verbose    bool
	Layer      string // Filter by layer
	Format     string // output format: text, json
	ShowPlan   bool   // Show derivation plan for stale artifacts
}

func runStatus() error {
	statusFlags := flag.NewFlagSet("status", flag.ExitOnError)
	projectDir := statusFlags.String("project-dir", ".", "Project root directory")
	verbose := statusFlags.Bool("verbose", false, "Show detailed status information")
	layer := statusFlags.String("layer", "", "Filter by layer (l0, l1, l2, l3)")
	format := statusFlags.String("format", "text", "Output format (text, json)")
	showPlan := statusFlags.Bool("plan", false, "Show derivation plan for stale artifacts")

	if len(os.Args) > 2 {
		statusFlags.Parse(os.Args[2:])
	}

	cfg := &StatusConfig{
		ProjectDir: *projectDir,
		Verbose:    *verbose,
		Layer:      *layer,
		Format:     *format,
		ShowPlan:   *showPlan,
	}

	return executeStatus(cfg)
}

func executeStatus(cfg *StatusConfig) error {
	// Create state manager and load state
	sm := derivation.NewStateManager(cfg.ProjectDir)
	state, err := sm.Load()
	if err != nil {
		return fmt.Errorf("failed to load state (run 'loom-cli init' first): %w", err)
	}

	// Create tracker for status detection
	tracker := derivation.NewTracker(state, cfg.ProjectDir)
	tracker.Verbose = cfg.Verbose

	// Detect stale artifacts
	staleArtifacts, err := tracker.DetectStaleArtifacts()
	if err != nil {
		return fmt.Errorf("failed to detect stale artifacts: %w", err)
	}

	// Build status summary
	summary := buildStatusSummary(state, staleArtifacts, cfg.Layer)

	// Output based on format
	if cfg.Format == "json" {
		return outputStatusJSON(summary, cfg.ShowPlan, tracker)
	}

	return outputStatusText(summary, cfg, tracker)
}

// StatusSummary holds aggregated status information
type StatusSummary struct {
	TotalArtifacts int
	ByStatus       map[derivation.ArtifactStatus]int
	ByLayer        map[string]map[derivation.ArtifactStatus]int
	StaleArtifacts []*derivation.Artifact
	Artifacts      []*derivation.Artifact // Filtered list
}

func buildStatusSummary(state *derivation.DerivationState, stale []*derivation.Artifact, layerFilter string) *StatusSummary {
	summary := &StatusSummary{
		ByStatus:       make(map[derivation.ArtifactStatus]int),
		ByLayer:        make(map[string]map[derivation.ArtifactStatus]int),
		StaleArtifacts: stale,
		Artifacts:      make([]*derivation.Artifact, 0),
	}

	// Mark stale artifacts in state
	staleIDs := make(map[string]bool)
	for _, a := range stale {
		staleIDs[a.ID] = true
	}

	for _, artifact := range state.Artifacts {
		// Apply layer filter
		if layerFilter != "" && artifact.Layer != layerFilter {
			continue
		}

		// Update status if stale
		status := artifact.Status
		if staleIDs[artifact.ID] {
			status = derivation.StatusStale
		}

		summary.TotalArtifacts++
		summary.ByStatus[status]++

		// By layer
		if summary.ByLayer[artifact.Layer] == nil {
			summary.ByLayer[artifact.Layer] = make(map[derivation.ArtifactStatus]int)
		}
		summary.ByLayer[artifact.Layer][status]++

		summary.Artifacts = append(summary.Artifacts, artifact)
	}

	// Sort artifacts by layer, then ID
	sort.Slice(summary.Artifacts, func(i, j int) bool {
		if summary.Artifacts[i].Layer != summary.Artifacts[j].Layer {
			return layerOrder(summary.Artifacts[i].Layer) < layerOrder(summary.Artifacts[j].Layer)
		}
		return summary.Artifacts[i].ID < summary.Artifacts[j].ID
	})

	return summary
}

func layerOrder(layer string) int {
	switch layer {
	case "l0":
		return 0
	case "l1":
		return 1
	case "l2":
		return 2
	case "l3":
		return 3
	default:
		return 99
	}
}

func outputStatusText(summary *StatusSummary, cfg *StatusConfig, tracker *derivation.Tracker) error {
	fmt.Println("=== Derivation Status ===")
	fmt.Println()

	// Overall summary
	fmt.Printf("Total artifacts: %d\n", summary.TotalArtifacts)
	fmt.Println()

	// Status breakdown
	fmt.Println("By Status:")
	statusOrder := []derivation.ArtifactStatus{
		derivation.StatusCurrent,
		derivation.StatusStale,
		derivation.StatusAffected,
		derivation.StatusModified,
		derivation.StatusOrphaned,
	}
	for _, status := range statusOrder {
		if count := summary.ByStatus[status]; count > 0 {
			fmt.Printf("  %s: %d\n", statusIcon(status), count)
		}
	}
	fmt.Println()

	// Layer breakdown
	fmt.Println("By Layer:")
	layers := []string{"l0", "l1", "l2", "l3"}
	for _, layer := range layers {
		if layerStats, ok := summary.ByLayer[layer]; ok {
			total := 0
			for _, count := range layerStats {
				total += count
			}
			if total > 0 {
				fmt.Printf("  %s: %d artifacts", strings.ToUpper(layer), total)
				// Show stale count if any
				if stale := layerStats[derivation.StatusStale]; stale > 0 {
					fmt.Printf(" (%d stale)", stale)
				}
				fmt.Println()
			}
		}
	}
	fmt.Println()

	// Stale artifacts detail
	if len(summary.StaleArtifacts) > 0 {
		fmt.Println("Stale Artifacts (need re-derivation):")
		for _, artifact := range summary.StaleArtifacts {
			fmt.Printf("  ⚠ %s (%s)\n", artifact.ID, artifact.Layer)
			if cfg.Verbose && len(artifact.Upstream) > 0 {
				fmt.Printf("    Upstream: %s\n", formatUpstream(artifact.Upstream))
			}
		}
		fmt.Println()
	}

	// Show derivation plan if requested
	if cfg.ShowPlan && len(summary.StaleArtifacts) > 0 {
		staleIDs := make([]string, 0, len(summary.StaleArtifacts))
		for _, a := range summary.StaleArtifacts {
			staleIDs = append(staleIDs, a.ID)
		}

		plan, err := tracker.PlanDerivation(staleIDs)
		if err != nil {
			fmt.Printf("Warning: failed to create derivation plan: %v\n", err)
		} else {
			fmt.Println("Derivation Plan:")
			for _, step := range plan.Artifacts {
				manual := ""
				if step.HasManual {
					manual = " [has manual sections]"
				}
				fmt.Printf("  %d. %s (%s)%s\n", step.Order, step.ArtifactID, step.Layer, manual)
			}
			fmt.Println()
		}
	}

	// Verbose artifact list
	if cfg.Verbose {
		fmt.Println("All Artifacts:")
		currentLayer := ""
		for _, artifact := range summary.Artifacts {
			if artifact.Layer != currentLayer {
				currentLayer = artifact.Layer
				fmt.Printf("\n  [%s]\n", strings.ToUpper(currentLayer))
			}
			icon := statusIcon(artifact.Status)
			fmt.Printf("    %s %s (%s)\n", icon, artifact.ID, artifact.Type)
		}
		fmt.Println()
	}

	return nil
}

func statusIcon(status derivation.ArtifactStatus) string {
	switch status {
	case derivation.StatusCurrent:
		return "✓ CURRENT"
	case derivation.StatusStale:
		return "⚠ STALE"
	case derivation.StatusAffected:
		return "↯ AFFECTED"
	case derivation.StatusModified:
		return "✎ MODIFIED"
	case derivation.StatusOrphaned:
		return "✗ ORPHANED"
	default:
		return "? UNKNOWN"
	}
}

func formatUpstream(upstream map[string]string) string {
	ids := make([]string, 0, len(upstream))
	for id := range upstream {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return strings.Join(ids, ", ")
}

func outputStatusJSON(summary *StatusSummary, showPlan bool, tracker *derivation.Tracker) error {
	// Build JSON structure
	output := struct {
		TotalArtifacts int                                             `json:"total_artifacts"`
		ByStatus       map[derivation.ArtifactStatus]int               `json:"by_status"`
		ByLayer        map[string]map[derivation.ArtifactStatus]int    `json:"by_layer"`
		StaleArtifacts []string                                        `json:"stale_artifacts"`
		Plan           *derivation.DerivationPlan                      `json:"derivation_plan,omitempty"`
	}{
		TotalArtifacts: summary.TotalArtifacts,
		ByStatus:       summary.ByStatus,
		ByLayer:        summary.ByLayer,
		StaleArtifacts: make([]string, 0, len(summary.StaleArtifacts)),
	}

	for _, a := range summary.StaleArtifacts {
		output.StaleArtifacts = append(output.StaleArtifacts, a.ID)
	}

	if showPlan && len(summary.StaleArtifacts) > 0 {
		staleIDs := make([]string, 0, len(summary.StaleArtifacts))
		for _, a := range summary.StaleArtifacts {
			staleIDs = append(staleIDs, a.ID)
		}
		plan, err := tracker.PlanDerivation(staleIDs)
		if err == nil {
			output.Plan = plan
		}
	}

	// Simple JSON output (avoid importing encoding/json for now)
	fmt.Printf("{\n")
	fmt.Printf("  \"total_artifacts\": %d,\n", output.TotalArtifacts)
	fmt.Printf("  \"stale_count\": %d,\n", len(output.StaleArtifacts))
	if len(output.StaleArtifacts) > 0 {
		fmt.Printf("  \"stale_artifacts\": [%s]\n", formatJSONArray(output.StaleArtifacts))
	}
	fmt.Printf("}\n")

	return nil
}

func formatJSONArray(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("\"%s\"", item)
	}
	return strings.Join(quoted, ", ")
}
