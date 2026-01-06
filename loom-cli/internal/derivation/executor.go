package derivation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// =============================================================================
// Derivation Executor
// =============================================================================

// Executor handles the execution of artifact derivations
type Executor struct {
	// State is the current derivation state
	State *DerivationState

	// Tracker for status updates
	Tracker *Tracker

	// Hasher for content hashing
	Hasher *Hasher

	// ProjectDir is the project root directory
	ProjectDir string

	// DryRun prevents actual file modifications
	DryRun bool

	// Verbose enables detailed logging
	Verbose bool

	// PreserveManual keeps manual sections during re-derivation
	PreserveManual bool

	// DeriverFunc is the function that performs actual derivation
	// It receives the artifact and its upstream content, returns new content
	DeriverFunc DeriverFunc

	// ProgressCallback is called during execution to report progress
	ProgressCallback ProgressCallback
}

// DeriverFunc is the signature for derivation functions
// It receives the artifact to derive, upstream artifacts, and returns new content
type DeriverFunc func(artifact *Artifact, upstreamContent map[string]string, projectDir string) (string, error)

// ProgressCallback reports execution progress
type ProgressCallback func(event ProgressEvent)

// ProgressEvent describes an execution progress event
type ProgressEvent struct {
	Type       ProgressType
	ArtifactID string
	Message    string
	Current    int
	Total      int
	Error      error
}

// ProgressType is the type of progress event
type ProgressType string

const (
	ProgressStart    ProgressType = "start"
	ProgressStep     ProgressType = "step"
	ProgressComplete ProgressType = "complete"
	ProgressError    ProgressType = "error"
	ProgressSkip     ProgressType = "skip"
)

// NewExecutor creates a new derivation executor
func NewExecutor(state *DerivationState, projectDir string) *Executor {
	return &Executor{
		State:          state,
		Tracker:        NewTracker(state, projectDir),
		Hasher:         NewHasher(),
		ProjectDir:     projectDir,
		PreserveManual: true,
	}
}

// =============================================================================
// Execution
// =============================================================================

// ExecutionResult holds the result of a derivation execution
type ExecutionResult struct {
	// Plan is the derivation plan that was executed
	Plan *DerivationPlan

	// Derived lists successfully derived artifacts
	Derived []DerivedArtifact

	// Skipped lists artifacts that were skipped
	Skipped []SkippedArtifact

	// Errors lists artifacts that failed to derive
	Errors []DerivationError

	// Duration is the total execution time
	Duration time.Duration

	// StartTime when execution started
	StartTime time.Time

	// EndTime when execution completed
	EndTime time.Time
}

// DerivedArtifact describes a successfully derived artifact
type DerivedArtifact struct {
	ArtifactID  string            `json:"artifact_id"`
	Layer       string            `json:"layer"`
	OutputFile  string            `json:"output_file"`
	OldHash     string            `json:"old_hash"`
	NewHash     string            `json:"new_hash"`
	FromHashes  map[string]string `json:"from_hashes"`
	Duration    time.Duration     `json:"duration"`
	HasManual   bool              `json:"has_manual"`
	ManualKept  []string          `json:"manual_kept,omitempty"`
}

// SkippedArtifact describes an artifact that was skipped
type SkippedArtifact struct {
	ArtifactID string `json:"artifact_id"`
	Reason     string `json:"reason"`
}

// DerivationError describes a failed derivation
type DerivationError struct {
	ArtifactID string `json:"artifact_id"`
	Error      string `json:"error"`
	Recoverable bool  `json:"recoverable"`
}

// Execute runs derivation for specified artifacts
func (e *Executor) Execute(artifactIDs []string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Derived:   make([]DerivedArtifact, 0),
		Skipped:   make([]SkippedArtifact, 0),
		Errors:    make([]DerivationError, 0),
		StartTime: time.Now(),
	}

	// Create derivation plan
	plan, err := e.Tracker.PlanDerivation(artifactIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create derivation plan: %w", err)
	}
	result.Plan = plan

	// Report start
	e.reportProgress(ProgressEvent{
		Type:    ProgressStart,
		Message: fmt.Sprintf("Starting derivation of %d artifacts", plan.TotalCount),
		Total:   plan.TotalCount,
	})

	// Execute each step in order
	for i, step := range plan.Artifacts {
		stepResult := e.executeStep(step, i+1, plan.TotalCount)

		switch stepResult.status {
		case "derived":
			result.Derived = append(result.Derived, stepResult.derived)
		case "skipped":
			result.Skipped = append(result.Skipped, stepResult.skipped)
		case "error":
			result.Errors = append(result.Errors, stepResult.err)
			// Continue with other artifacts unless it's a critical error
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Report completion
	e.reportProgress(ProgressEvent{
		Type:    ProgressComplete,
		Message: fmt.Sprintf("Derivation complete: %d derived, %d skipped, %d errors",
			len(result.Derived), len(result.Skipped), len(result.Errors)),
	})

	return result, nil
}

// ExecuteAll detects and derives all stale artifacts
func (e *Executor) ExecuteAll() (*ExecutionResult, error) {
	// Detect stale artifacts
	stale, err := e.Tracker.DetectStaleArtifacts()
	if err != nil {
		return nil, fmt.Errorf("failed to detect stale artifacts: %w", err)
	}

	if len(stale) == 0 {
		return &ExecutionResult{
			Plan:      &DerivationPlan{},
			Derived:   make([]DerivedArtifact, 0),
			Skipped:   make([]SkippedArtifact, 0),
			Errors:    make([]DerivationError, 0),
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}, nil
	}

	// Collect IDs
	ids := make([]string, len(stale))
	for i, a := range stale {
		ids[i] = a.ID
	}

	return e.Execute(ids)
}

// stepResult holds the result of a single derivation step
type stepResult struct {
	status  string // "derived", "skipped", "error"
	derived DerivedArtifact
	skipped SkippedArtifact
	err     DerivationError
}

func (e *Executor) executeStep(step *DerivationStep, current, total int) stepResult {
	artifact := e.State.GetArtifact(step.ArtifactID)
	if artifact == nil {
		return stepResult{
			status: "error",
			err: DerivationError{
				ArtifactID:  step.ArtifactID,
				Error:       "artifact not found in state",
				Recoverable: false,
			},
		}
	}

	// Report step start
	e.reportProgress(ProgressEvent{
		Type:       ProgressStep,
		ArtifactID: step.ArtifactID,
		Message:    fmt.Sprintf("Deriving %s (%s)", step.ArtifactID, step.Layer),
		Current:    current,
		Total:      total,
	})

	startTime := time.Now()

	// Check if deriver function is set
	if e.DeriverFunc == nil {
		return stepResult{
			status: "skipped",
			skipped: SkippedArtifact{
				ArtifactID: step.ArtifactID,
				Reason:     "no deriver function configured",
			},
		}
	}

	// Collect upstream content
	upstreamContent, err := e.collectUpstreamContent(artifact)
	if err != nil {
		return stepResult{
			status: "error",
			err: DerivationError{
				ArtifactID:  step.ArtifactID,
				Error:       fmt.Sprintf("failed to collect upstream content: %v", err),
				Recoverable: true,
			},
		}
	}

	// Preserve manual sections if configured
	var manualContent map[string]string
	if e.PreserveManual && step.HasManual {
		manualContent, err = e.extractManualSections(artifact)
		if err != nil {
			e.log("Warning: failed to extract manual sections: %v", err)
		}
	}

	// Execute derivation
	oldHash := artifact.ContentHash
	newContent, err := e.DeriverFunc(artifact, upstreamContent, e.ProjectDir)
	if err != nil {
		return stepResult{
			status: "error",
			err: DerivationError{
				ArtifactID:  step.ArtifactID,
				Error:       fmt.Sprintf("derivation failed: %v", err),
				Recoverable: true,
			},
		}
	}

	// Restore manual sections
	if manualContent != nil && len(manualContent) > 0 {
		newContent = e.restoreManualSections(newContent, manualContent)
	}

	// Write output (unless dry run)
	outputFile := artifact.Location.File
	if !e.DryRun {
		if err := e.writeOutput(artifact, newContent); err != nil {
			return stepResult{
				status: "error",
				err: DerivationError{
					ArtifactID:  step.ArtifactID,
					Error:       fmt.Sprintf("failed to write output: %v", err),
					Recoverable: true,
				},
			}
		}
	}

	// Update artifact state
	upstreamHashes, _ := e.Hasher.CollectUpstreamHashes(artifact, e.State, e.ProjectDir)
	newHash := e.Hasher.HashContent(newContent)

	if !e.DryRun {
		artifact.ContentHash = newHash
		artifact.DerivedFromHashes = upstreamHashes
		artifact.DerivedAt = time.Now()
		artifact.Status = StatusCurrent
	}

	return stepResult{
		status: "derived",
		derived: DerivedArtifact{
			ArtifactID: step.ArtifactID,
			Layer:      step.Layer,
			OutputFile: outputFile,
			OldHash:    oldHash,
			NewHash:    newHash,
			FromHashes: upstreamHashes,
			Duration:   time.Since(startTime),
			HasManual:  step.HasManual,
			ManualKept: mapKeys(manualContent),
		},
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

func (e *Executor) collectUpstreamContent(artifact *Artifact) (map[string]string, error) {
	content := make(map[string]string)

	for upstreamID := range artifact.Upstream {
		upstream := e.State.GetArtifact(upstreamID)
		if upstream == nil {
			continue
		}

		// Read upstream content from file
		if upstream.Location.File != "" {
			filePath := upstream.Location.File
			if !filepath.IsAbs(filePath) {
				filePath = filepath.Join(e.ProjectDir, filePath)
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read upstream %s: %w", upstreamID, err)
			}

			// Extract relevant section if line range specified
			fileContent := string(data)
			if upstream.Location.LineStart > 0 {
				lines := strings.Split(fileContent, "\n")
				start := upstream.Location.LineStart - 1
				end := len(lines)
				if upstream.Location.LineEnd > 0 && upstream.Location.LineEnd < len(lines) {
					end = upstream.Location.LineEnd
				}
				if start < len(lines) {
					fileContent = strings.Join(lines[start:end], "\n")
				}
			}

			content[upstreamID] = fileContent
		}
	}

	return content, nil
}

func (e *Executor) extractManualSections(artifact *Artifact) (map[string]string, error) {
	if artifact.Location.File == "" {
		return nil, nil
	}

	filePath := artifact.Location.File
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(e.ProjectDir, filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Use parser to extract manual sections
	parser := NewParser()
	doc := parser.ParseContent(string(content), filePath)

	manualSections := make(map[string]string)
	for _, section := range doc.Sections {
		if section.Type == "manual" {
			manualSections[section.ID] = section.Content
		}
	}

	return manualSections, nil
}

func (e *Executor) restoreManualSections(content string, manualSections map[string]string) string {
	// For each manual section, find the placeholder and replace
	// This assumes generated content has markers like <!-- LOOM:MANUAL section="name" -->
	for sectionID, sectionContent := range manualSections {
		placeholder := fmt.Sprintf("<!-- LOOM:MANUAL section=\"%s\" -->", sectionID)
		if strings.Contains(content, placeholder) {
			// Replace placeholder with actual content
			replacement := fmt.Sprintf("%s\n%s", placeholder, sectionContent)
			content = strings.Replace(content, placeholder, replacement, 1)
		}
	}
	return content
}

func (e *Executor) writeOutput(artifact *Artifact, content string) error {
	filePath := artifact.Location.File
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(e.ProjectDir, filePath)
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write content
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (e *Executor) reportProgress(event ProgressEvent) {
	if e.ProgressCallback != nil {
		e.ProgressCallback(event)
	}

	if e.Verbose {
		switch event.Type {
		case ProgressStart:
			fmt.Printf("[START] %s\n", event.Message)
		case ProgressStep:
			fmt.Printf("[%d/%d] %s\n", event.Current, event.Total, event.Message)
		case ProgressComplete:
			fmt.Printf("[DONE] %s\n", event.Message)
		case ProgressError:
			fmt.Printf("[ERROR] %s: %v\n", event.ArtifactID, event.Error)
		case ProgressSkip:
			fmt.Printf("[SKIP] %s: %s\n", event.ArtifactID, event.Message)
		}
	}
}

func (e *Executor) log(format string, args ...interface{}) {
	if e.Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func mapKeys(m map[string]string) []string {
	if m == nil {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// =============================================================================
// Impact Preview
// =============================================================================

// PreviewExecution shows what would be derived without making changes
func (e *Executor) PreviewExecution(artifactIDs []string) (*DerivationPlan, *ImpactReport, error) {
	// Create plan
	plan, err := e.Tracker.PlanDerivation(artifactIDs)
	if err != nil {
		return nil, nil, err
	}

	// Analyze impact
	impact := e.Tracker.AnalyzeImpact(artifactIDs)

	return plan, impact, nil
}

// =============================================================================
// Selective Derivation
// =============================================================================

// ExecuteLayer derives all stale artifacts in a specific layer
func (e *Executor) ExecuteLayer(layer string) (*ExecutionResult, error) {
	stale, err := e.Tracker.DetectStaleArtifacts()
	if err != nil {
		return nil, err
	}

	// Filter by layer
	var layerStale []string
	for _, a := range stale {
		if a.Layer == layer {
			layerStale = append(layerStale, a.ID)
		}
	}

	if len(layerStale) == 0 {
		return &ExecutionResult{
			Plan:      &DerivationPlan{},
			Derived:   make([]DerivedArtifact, 0),
			Skipped:   make([]SkippedArtifact, 0),
			Errors:    make([]DerivationError, 0),
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}, nil
	}

	return e.Execute(layerStale)
}

// ExecuteType derives all stale artifacts of a specific type
func (e *Executor) ExecuteType(artType ArtifactType) (*ExecutionResult, error) {
	stale, err := e.Tracker.DetectStaleArtifacts()
	if err != nil {
		return nil, err
	}

	// Filter by type
	var typeStale []string
	for _, a := range stale {
		if a.Type == artType {
			typeStale = append(typeStale, a.ID)
		}
	}

	if len(typeStale) == 0 {
		return &ExecutionResult{
			Plan:      &DerivationPlan{},
			Derived:   make([]DerivedArtifact, 0),
			Skipped:   make([]SkippedArtifact, 0),
			Errors:    make([]DerivationError, 0),
			StartTime: time.Now(),
			EndTime:   time.Now(),
		}, nil
	}

	return e.Execute(typeStale)
}
