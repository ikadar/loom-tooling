package derivation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	executor := NewExecutor(state, "/test/project")

	if executor.State != state {
		t.Error("State should be set")
	}

	if executor.ProjectDir != "/test/project" {
		t.Errorf("Expected project dir /test/project, got %s", executor.ProjectDir)
	}

	if executor.Tracker == nil {
		t.Error("Tracker should be created")
	}

	if executor.Hasher == nil {
		t.Error("Hasher should be created")
	}

	if !executor.PreserveManual {
		t.Error("PreserveManual should default to true")
	}
}

func TestExecutor_Execute_NoDeriverFunc(t *testing.T) {
	tmpDir := t.TempDir()

	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Add upstream artifact
	upstream := &Artifact{
		ID:          "US-ORD-001",
		Type:        ArtifactUserStory,
		Layer:       "l0",
		Status:      StatusCurrent,
		ContentHash: "sha256:upstream",
	}
	state.Artifacts[upstream.ID] = upstream

	// Add a downstream artifact
	artifact := &Artifact{
		ID:     "AC-ORD-001",
		Type:   ArtifactAcceptanceCrit,
		Layer:  "l1",
		Status: StatusStale,
		Location: ArtifactLocation{
			File: filepath.Join(tmpDir, "l1", "ac.md"),
		},
		Upstream: map[string]string{
			"US-ORD-001": "sha256:upstream",
		},
	}
	state.Artifacts[artifact.ID] = artifact
	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	executor := NewExecutor(state, tmpDir)
	// DeriverFunc is nil

	result, err := executor.Execute([]string{"AC-ORD-001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should skip because no deriver function
	if len(result.Skipped) != 1 {
		t.Errorf("Expected 1 skipped, got %d", len(result.Skipped))
	}

	if len(result.Derived) != 0 {
		t.Errorf("Expected 0 derived, got %d", len(result.Derived))
	}
}

func TestExecutor_Execute_WithDeriverFunc(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	l1Dir := filepath.Join(tmpDir, "l1")
	os.MkdirAll(l1Dir, 0755)

	// Create upstream file
	l0Dir := filepath.Join(tmpDir, "l0")
	os.MkdirAll(l0Dir, 0755)
	upstreamFile := filepath.Join(l0Dir, "user-stories.md")
	os.WriteFile(upstreamFile, []byte("# User Story\nUS-ORD-001: As a user..."), 0644)

	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Add upstream artifact
	upstream := &Artifact{
		ID:          "US-ORD-001",
		Type:        ArtifactUserStory,
		Layer:       "l0",
		Status:      StatusCurrent,
		ContentHash: "sha256:upstream",
		Location: ArtifactLocation{
			File: upstreamFile,
		},
	}
	state.Artifacts[upstream.ID] = upstream

	// Add downstream artifact
	acFile := filepath.Join(l1Dir, "ac.md")
	os.WriteFile(acFile, []byte("# Old Content"), 0644)

	artifact := &Artifact{
		ID:     "AC-ORD-001",
		Type:   ArtifactAcceptanceCrit,
		Layer:  "l1",
		Status: StatusStale,
		Location: ArtifactLocation{
			File: acFile,
		},
		Upstream: map[string]string{
			"US-ORD-001": "sha256:old-hash", // Different from current
		},
	}
	state.Artifacts[artifact.ID] = artifact
	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	executor := NewExecutor(state, tmpDir)

	// Set a simple deriver function
	derivedContent := ""
	executor.DeriverFunc = func(art *Artifact, upstream map[string]string, projectDir string) (string, error) {
		derivedContent = "# Derived Content\n## AC-ORD-001\nNew content"
		return derivedContent, nil
	}

	result, err := executor.Execute([]string{"AC-ORD-001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should derive successfully
	if len(result.Derived) != 1 {
		t.Errorf("Expected 1 derived, got %d", len(result.Derived))
	}

	if len(result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d: %v", len(result.Errors), result.Errors)
	}

	// Check file was written
	content, err := os.ReadFile(acFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(content) != derivedContent {
		t.Errorf("File content mismatch.\nExpected: %s\nGot: %s", derivedContent, string(content))
	}
}

func TestExecutor_Execute_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	l1Dir := filepath.Join(tmpDir, "l1")
	os.MkdirAll(l1Dir, 0755)

	// Create upstream file
	l0Dir := filepath.Join(tmpDir, "l0")
	os.MkdirAll(l0Dir, 0755)
	upstreamFile := filepath.Join(l0Dir, "user-stories.md")
	os.WriteFile(upstreamFile, []byte("# User Story"), 0644)

	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Add upstream
	upstream := &Artifact{
		ID:          "US-ORD-001",
		Type:        ArtifactUserStory,
		Layer:       "l0",
		Status:      StatusCurrent,
		ContentHash: "sha256:upstream",
		Location: ArtifactLocation{
			File: upstreamFile,
		},
	}
	state.Artifacts[upstream.ID] = upstream

	// Add downstream
	acFile := filepath.Join(l1Dir, "ac.md")
	originalContent := "# Original Content"
	os.WriteFile(acFile, []byte(originalContent), 0644)

	artifact := &Artifact{
		ID:          "AC-ORD-001",
		Type:        ArtifactAcceptanceCrit,
		Layer:       "l1",
		Status:      StatusStale,
		ContentHash: "sha256:original",
		Location: ArtifactLocation{
			File: acFile,
		},
		Upstream: map[string]string{
			"US-ORD-001": "sha256:old",
		},
	}
	state.Artifacts[artifact.ID] = artifact
	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	executor := NewExecutor(state, tmpDir)
	executor.DryRun = true

	executor.DeriverFunc = func(art *Artifact, upstream map[string]string, projectDir string) (string, error) {
		return "# New Content", nil
	}

	result, err := executor.Execute([]string{"AC-ORD-001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should still report as derived
	if len(result.Derived) != 1 {
		t.Errorf("Expected 1 derived, got %d", len(result.Derived))
	}

	// But file should NOT be modified
	content, err := os.ReadFile(acFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != originalContent {
		t.Error("File should not be modified in dry run mode")
	}

	// State should not be modified
	if artifact.ContentHash != "sha256:original" {
		t.Error("Artifact hash should not be modified in dry run mode")
	}
}

func TestExecutor_ExecuteAll(t *testing.T) {
	tmpDir := t.TempDir()

	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Add a current artifact (should not be derived)
	current := &Artifact{
		ID:     "AC-ORD-001",
		Type:   ArtifactAcceptanceCrit,
		Layer:  "l1",
		Status: StatusCurrent,
	}
	state.Artifacts[current.ID] = current

	executor := NewExecutor(state, tmpDir)

	result, err := executor.ExecuteAll()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// No stale artifacts, so nothing to derive
	if result.Plan.TotalCount != 0 {
		t.Errorf("Expected 0 in plan, got %d", result.Plan.TotalCount)
	}
}

func TestExecutor_ExecuteLayer(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	l1Dir := filepath.Join(tmpDir, "l1")
	l2Dir := filepath.Join(tmpDir, "l2")
	os.MkdirAll(l1Dir, 0755)
	os.MkdirAll(l2Dir, 0755)

	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// L0 upstream
	l0 := &Artifact{
		ID:          "US-ORD-001",
		Type:        ArtifactUserStory,
		Layer:       "l0",
		Status:      StatusCurrent,
		ContentHash: "sha256:new",
	}
	state.Artifacts[l0.ID] = l0

	// L1 artifact - stale
	l1File := filepath.Join(l1Dir, "ac.md")
	os.WriteFile(l1File, []byte("# AC"), 0644)
	l1 := &Artifact{
		ID:     "AC-ORD-001",
		Type:   ArtifactAcceptanceCrit,
		Layer:  "l1",
		Status: StatusCurrent, // Will be detected as stale
		Location: ArtifactLocation{
			File: l1File,
		},
		Upstream: map[string]string{
			"US-ORD-001": "sha256:old", // Different
		},
		DerivedFromHashes: map[string]string{
			"US-ORD-001": "sha256:old",
		},
	}
	state.Artifacts[l1.ID] = l1
	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	// L2 artifact - current
	l2File := filepath.Join(l2Dir, "ts.md")
	os.WriteFile(l2File, []byte("# TS"), 0644)
	l2 := &Artifact{
		ID:     "TS-ORD-001",
		Type:   ArtifactTechSpec,
		Layer:  "l2",
		Status: StatusCurrent,
		Location: ArtifactLocation{
			File: l2File,
		},
	}
	state.Artifacts[l2.ID] = l2

	executor := NewExecutor(state, tmpDir)
	executor.DryRun = true
	executor.DeriverFunc = func(art *Artifact, upstream map[string]string, projectDir string) (string, error) {
		return "# Derived", nil
	}

	// Execute only L1
	result, err := executor.ExecuteLayer("l1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include L1 stale artifact
	if result.Plan.TotalCount != 1 {
		t.Errorf("Expected 1 in plan, got %d", result.Plan.TotalCount)
	}

	if len(result.Derived) != 1 {
		t.Errorf("Expected 1 derived, got %d", len(result.Derived))
	}

	if result.Derived[0].Layer != "l1" {
		t.Errorf("Expected layer l1, got %s", result.Derived[0].Layer)
	}
}

func TestExecutor_PreviewExecution(t *testing.T) {
	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Add artifacts
	l0 := &Artifact{
		ID:     "US-ORD-001",
		Type:   ArtifactUserStory,
		Layer:  "l0",
		Status: StatusCurrent,
	}
	state.Artifacts[l0.ID] = l0

	l1 := &Artifact{
		ID:     "AC-ORD-001",
		Type:   ArtifactAcceptanceCrit,
		Layer:  "l1",
		Status: StatusStale,
		Upstream: map[string]string{
			"US-ORD-001": "sha256:hash",
		},
	}
	state.Artifacts[l1.ID] = l1
	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	l2 := &Artifact{
		ID:     "TS-ORD-001",
		Type:   ArtifactTechSpec,
		Layer:  "l2",
		Status: StatusCurrent,
		Upstream: map[string]string{
			"AC-ORD-001": "sha256:hash",
		},
	}
	state.Artifacts[l2.ID] = l2
	state.DependencyGraph.AddEdge("AC-ORD-001", "TS-ORD-001", EdgeDerives)

	executor := NewExecutor(state, "/test")

	plan, impact, err := executor.PreviewExecution([]string{"AC-ORD-001"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check plan - includes AC-ORD-001 and its downstream TS-ORD-001
	if plan.TotalCount != 2 {
		t.Errorf("Expected 2 in plan (AC and downstream TS), got %d", plan.TotalCount)
	}

	// Check impact - TS-ORD-001 should be affected
	if len(impact.AffectedArtifacts) != 1 {
		t.Errorf("Expected 1 affected, got %d", len(impact.AffectedArtifacts))
	}

	if len(impact.AffectedArtifacts) > 0 && impact.AffectedArtifacts[0] != "TS-ORD-001" {
		t.Errorf("Expected TS-ORD-001 affected, got %s", impact.AffectedArtifacts[0])
	}
}

func TestExecutionResult_Duration(t *testing.T) {
	result := &ExecutionResult{
		StartTime: time.Now(),
	}

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if result.Duration < 10*time.Millisecond {
		t.Error("Duration should be at least 10ms")
	}
}

func TestExecutor_ProgressCallback(t *testing.T) {
	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	artifact := &Artifact{
		ID:       "AC-ORD-001",
		Type:     ArtifactAcceptanceCrit,
		Layer:    "l1",
		Status:   StatusStale,
		Upstream: make(map[string]string),
	}
	state.Artifacts[artifact.ID] = artifact

	executor := NewExecutor(state, "/test")

	var events []ProgressEvent
	executor.ProgressCallback = func(event ProgressEvent) {
		events = append(events, event)
	}

	executor.Execute([]string{"AC-ORD-001"})

	// Should have received events
	if len(events) < 2 {
		t.Errorf("Expected at least 2 events, got %d", len(events))
	}

	// First should be start
	if events[0].Type != ProgressStart {
		t.Errorf("Expected first event to be start, got %s", events[0].Type)
	}

	// Last should be complete
	if events[len(events)-1].Type != ProgressComplete {
		t.Errorf("Expected last event to be complete, got %s", events[len(events)-1].Type)
	}
}
