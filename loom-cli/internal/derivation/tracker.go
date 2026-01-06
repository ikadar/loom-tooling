package derivation

import (
	"fmt"
	"sort"
	"time"
)

// =============================================================================
// Status Tracker
// =============================================================================

// Tracker monitors artifact status and detects changes
type Tracker struct {
	// State is the derivation state being tracked
	State *DerivationState

	// Hasher is used for content hashing
	Hasher *Hasher

	// ProjectDir is the project root directory
	ProjectDir string

	// Verbose enables detailed logging
	Verbose bool
}

// NewTracker creates a new status tracker
func NewTracker(state *DerivationState, projectDir string) *Tracker {
	return &Tracker{
		State:      state,
		Hasher:     NewHasher(),
		ProjectDir: projectDir,
		Verbose:    false,
	}
}

// =============================================================================
// Status Detection
// =============================================================================

// DetectStaleArtifacts finds all artifacts that need re-derivation
func (t *Tracker) DetectStaleArtifacts() ([]*Artifact, error) {
	var stale []*Artifact

	for _, artifact := range t.State.Artifacts {
		// Skip artifacts with no upstream (root artifacts like L0)
		if len(artifact.Upstream) == 0 {
			continue
		}

		if artifact.Status == StatusCurrent || artifact.Status == StatusModified {
			isStale, err := t.isArtifactStale(artifact)
			if err != nil {
				return nil, fmt.Errorf("failed to check staleness of %s: %w", artifact.ID, err)
			}
			if isStale {
				stale = append(stale, artifact)
			}
		}
	}

	// Sort by layer order (L0 first, then L1, L2, L3)
	sort.Slice(stale, func(i, j int) bool {
		return layerOrder(stale[i].Layer) < layerOrder(stale[j].Layer)
	})

	return stale, nil
}

// isArtifactStale checks if an artifact's upstream dependencies have changed
func (t *Tracker) isArtifactStale(artifact *Artifact) (bool, error) {
	// Get current upstream hashes
	currentHashes, err := t.Hasher.CollectUpstreamHashes(artifact, t.State, t.ProjectDir)
	if err != nil {
		return false, err
	}

	// Compare with stored hashes
	return artifact.IsStale(currentHashes), nil
}

// layerOrder returns numeric order for layers (for sorting)
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

// =============================================================================
// Status Updates
// =============================================================================

// UpdateStatuses recalculates status for all artifacts
func (t *Tracker) UpdateStatuses() error {
	// Phase 1: Mark directly stale artifacts
	staleIDs := make(map[string]bool)
	for _, artifact := range t.State.Artifacts {
		isStale, err := t.isArtifactStale(artifact)
		if err != nil {
			return err
		}
		if isStale {
			artifact.Status = StatusStale
			staleIDs[artifact.ID] = true
		}
	}

	// Phase 2: Mark downstream artifacts as affected
	for id := range staleIDs {
		affected := t.State.DependencyGraph.GetAllDownstream(id)
		for _, affectedID := range affected {
			artifact := t.State.GetArtifact(affectedID)
			if artifact != nil && artifact.Status == StatusCurrent {
				artifact.Status = StatusAffected
			}
		}
	}

	return nil
}

// MarkAsDerived updates an artifact after successful derivation
func (t *Tracker) MarkAsDerived(artifactID string, upstreamHashes map[string]string) error {
	artifact := t.State.GetArtifact(artifactID)
	if artifact == nil {
		return fmt.Errorf("artifact not found: %s", artifactID)
	}

	// Update artifact state
	artifact.DerivedFromHashes = upstreamHashes
	artifact.DerivedAt = time.Now()
	artifact.Status = StatusCurrent

	// Compute new content hash
	hash, err := t.Hasher.HashArtifact(artifact, t.ProjectDir)
	if err != nil {
		return fmt.Errorf("failed to hash artifact: %w", err)
	}
	artifact.ContentHash = hash

	return nil
}

// MarkAsModified marks an artifact as having manual modifications
func (t *Tracker) MarkAsModified(artifactID string, manualSections []string) error {
	artifact := t.State.GetArtifact(artifactID)
	if artifact == nil {
		return fmt.Errorf("artifact not found: %s", artifactID)
	}

	artifact.Status = StatusModified
	artifact.ManualSections = manualSections

	return nil
}

// MarkAsOrphaned marks an artifact as orphaned (no longer has upstream)
func (t *Tracker) MarkAsOrphaned(artifactID string) error {
	artifact := t.State.GetArtifact(artifactID)
	if artifact == nil {
		return fmt.Errorf("artifact not found: %s", artifactID)
	}

	artifact.Status = StatusOrphaned
	return nil
}

// =============================================================================
// Change Impact Analysis
// =============================================================================

// AnalyzeImpact determines what would be affected if given artifacts are re-derived
func (t *Tracker) AnalyzeImpact(artifactIDs []string) *ImpactReport {
	report := &ImpactReport{
		ChangedArtifacts:   artifactIDs,
		AffectedArtifacts:  make([]string, 0),
		AffectedByLayer:    make(map[string][]string),
		DerivationOrder:    make([]string, 0),
		ManualEditWarnings: make([]ManualEditWarning, 0),
	}

	// Collect all affected artifacts
	affectedSet := make(map[string]bool)
	for _, id := range artifactIDs {
		downstream := t.State.DependencyGraph.GetAllDownstream(id)
		for _, downID := range downstream {
			if !affectedSet[downID] {
				affectedSet[downID] = true
				report.AffectedArtifacts = append(report.AffectedArtifacts, downID)
			}
		}
	}

	// Group by layer
	for _, id := range report.AffectedArtifacts {
		artifact := t.State.GetArtifact(id)
		if artifact != nil {
			report.AffectedByLayer[artifact.Layer] = append(
				report.AffectedByLayer[artifact.Layer], id)
		}
	}

	// Get derivation order
	allAffected := append(artifactIDs, report.AffectedArtifacts...)
	order, err := t.State.DependencyGraph.GetDerivationOrder(allAffected)
	if err == nil {
		report.DerivationOrder = order
	}

	// Check for manual edit warnings
	for _, id := range report.AffectedArtifacts {
		artifact := t.State.GetArtifact(id)
		if artifact != nil && artifact.HasManualEdits() {
			report.ManualEditWarnings = append(report.ManualEditWarnings, ManualEditWarning{
				ArtifactID:     id,
				ManualSections: artifact.ManualSections,
				Message:        fmt.Sprintf("Artifact %s has manual edits in sections: %v", id, artifact.ManualSections),
			})
		}
	}

	return report
}

// ImpactReport describes the impact of re-deriving artifacts
type ImpactReport struct {
	// ChangedArtifacts are the directly changed artifacts
	ChangedArtifacts []string `json:"changed_artifacts"`

	// AffectedArtifacts are downstream artifacts that would need update
	AffectedArtifacts []string `json:"affected_artifacts"`

	// AffectedByLayer groups affected artifacts by layer
	AffectedByLayer map[string][]string `json:"affected_by_layer"`

	// DerivationOrder is the order in which to re-derive
	DerivationOrder []string `json:"derivation_order"`

	// ManualEditWarnings lists artifacts with manual edits that would be affected
	ManualEditWarnings []ManualEditWarning `json:"manual_edit_warnings,omitempty"`
}

// ManualEditWarning warns about manual edits that may be affected
type ManualEditWarning struct {
	ArtifactID     string   `json:"artifact_id"`
	ManualSections []string `json:"manual_sections"`
	Message        string   `json:"message"`
}

// =============================================================================
// Derivation Planning
// =============================================================================

// PlanDerivation creates an ordered list of artifacts to derive
func (t *Tracker) PlanDerivation(staleIDs []string) (*DerivationPlan, error) {
	plan := &DerivationPlan{
		Artifacts:   make([]*DerivationStep, 0),
		TotalCount:  0,
		ByLayer:     make(map[string]int),
		Warnings:    make([]string, 0),
	}

	// Get derivation order
	order, err := t.State.DependencyGraph.GetDerivationOrder(staleIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to compute derivation order: %w", err)
	}

	for i, id := range order {
		artifact := t.State.GetArtifact(id)
		if artifact == nil {
			plan.Warnings = append(plan.Warnings,
				fmt.Sprintf("Artifact %s in derivation order but not in state", id))
			continue
		}

		step := &DerivationStep{
			Order:      i + 1,
			ArtifactID: id,
			Type:       artifact.Type,
			Layer:      artifact.Layer,
			Upstream:   make([]string, 0, len(artifact.Upstream)),
			HasManual:  artifact.HasManualEdits(),
		}

		// Collect upstream IDs
		for upID := range artifact.Upstream {
			step.Upstream = append(step.Upstream, upID)
		}
		sort.Strings(step.Upstream)

		plan.Artifacts = append(plan.Artifacts, step)
		plan.TotalCount++
		plan.ByLayer[artifact.Layer]++
	}

	return plan, nil
}

// DerivationPlan describes the plan for re-deriving artifacts
type DerivationPlan struct {
	// Artifacts is the ordered list of derivation steps
	Artifacts []*DerivationStep `json:"artifacts"`

	// TotalCount is the total number of artifacts to derive
	TotalCount int `json:"total_count"`

	// ByLayer counts artifacts by layer
	ByLayer map[string]int `json:"by_layer"`

	// Warnings lists any planning warnings
	Warnings []string `json:"warnings,omitempty"`
}

// DerivationStep describes one artifact to derive
type DerivationStep struct {
	// Order is the derivation order (1-indexed)
	Order int `json:"order"`

	// ArtifactID is the artifact to derive
	ArtifactID string `json:"artifact_id"`

	// Type is the artifact type
	Type ArtifactType `json:"type"`

	// Layer is the artifact layer
	Layer string `json:"layer"`

	// Upstream lists the upstream dependencies
	Upstream []string `json:"upstream"`

	// HasManual indicates if the artifact has manual sections
	HasManual bool `json:"has_manual"`
}

// =============================================================================
// File Watching
// =============================================================================

// DetectFileChanges compares current file hashes with cached hashes
func (t *Tracker) DetectFileChanges() ([]FileChange, error) {
	var changes []FileChange

	// Get all artifact files
	files := make(map[string]bool)
	for _, artifact := range t.State.Artifacts {
		if artifact.Location.File != "" {
			files[artifact.Location.File] = true
		}
	}

	// Check each file
	for file := range files {
		cached := t.State.GetFileHash(file)
		if t.Hasher.NeedsRehash(file, cached) {
			newInfo, err := t.Hasher.HashFileWithInfo(file)
			if err != nil {
				changes = append(changes, FileChange{
					Path:       file,
					ChangeType: "error",
					Error:      err.Error(),
				})
				continue
			}

			changeType := "modified"
			if cached == nil {
				changeType = "new"
			}

			changes = append(changes, FileChange{
				Path:       file,
				ChangeType: changeType,
				OldHash:    cachedHash(cached),
				NewHash:    newInfo.Hash,
			})
		}
	}

	return changes, nil
}

// FileChange describes a detected file change
type FileChange struct {
	// Path is the file path
	Path string `json:"path"`

	// ChangeType is the type of change (new, modified, deleted, error)
	ChangeType string `json:"change_type"`

	// OldHash is the previous hash (empty for new files)
	OldHash string `json:"old_hash,omitempty"`

	// NewHash is the new hash (empty for deleted files)
	NewHash string `json:"new_hash,omitempty"`

	// Error contains any error message
	Error string `json:"error,omitempty"`
}

// cachedHash safely gets hash from cached info
func cachedHash(info *FileHashInfo) string {
	if info == nil {
		return ""
	}
	return info.Hash
}

// =============================================================================
// Synchronization
// =============================================================================

// SyncFromFiles updates state based on current file contents
func (t *Tracker) SyncFromFiles(parser *Parser) error {
	// Track seen artifacts
	seen := make(map[string]bool)

	// Parse all artifact files
	for _, artifact := range t.State.Artifacts {
		if artifact.Location.File == "" {
			continue
		}

		doc, err := parser.ParseFile(artifact.Location.File)
		if err != nil {
			// File might have been deleted
			artifact.Status = StatusOrphaned
			continue
		}

		// Update artifact from parsed content
		for _, parsed := range doc.Artifacts {
			if parsed.ID == artifact.ID {
				seen[artifact.ID] = true

				// Update location
				artifact.Location = parsed.Location

				// Update manual sections
				artifact.ManualSections = parsed.ManualSections

				// Compute new hash
				hash, err := t.Hasher.HashArtifact(artifact, t.ProjectDir)
				if err == nil {
					if artifact.ContentHash != hash {
						// Content changed
						if artifact.Status == StatusCurrent {
							artifact.Status = StatusModified
						}
						artifact.ContentHash = hash
					}
				}
			}
		}
	}

	// Mark unseen artifacts as orphaned
	for id, artifact := range t.State.Artifacts {
		if !seen[id] && artifact.Location.File != "" {
			artifact.Status = StatusOrphaned
		}
	}

	return nil
}

// =============================================================================
// Cleanup
// =============================================================================

// CleanupOrphaned removes orphaned artifacts from state
func (t *Tracker) CleanupOrphaned() []string {
	var removed []string

	for id, artifact := range t.State.Artifacts {
		if artifact.Status == StatusOrphaned {
			t.State.RemoveArtifact(id)
			t.State.DependencyGraph.RemoveNode(id)
			removed = append(removed, id)
		}
	}

	return removed
}

// ValidateGraph checks the dependency graph for issues
func (t *Tracker) ValidateGraph() []string {
	var issues []string

	// Check for cycles
	if cycle := t.State.DependencyGraph.DetectCycle(); cycle != nil {
		issues = append(issues, fmt.Sprintf("Dependency cycle detected: %v", cycle))
	}

	// Check for missing nodes
	for _, edge := range t.State.DependencyGraph.Edges {
		if t.State.GetArtifact(edge.From) == nil {
			issues = append(issues, fmt.Sprintf("Edge references missing artifact: %s", edge.From))
		}
		if t.State.GetArtifact(edge.To) == nil {
			issues = append(issues, fmt.Sprintf("Edge references missing artifact: %s", edge.To))
		}
	}

	return issues
}
