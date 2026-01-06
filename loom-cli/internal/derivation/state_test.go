package derivation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStateManager_NewState(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewStateManager(tmpDir)

	state := sm.NewState()

	if state.Version != StateVersion {
		t.Errorf("Expected version %s, got %s", StateVersion, state.Version)
	}
	if state.Project != filepath.Base(tmpDir) {
		t.Errorf("Expected project name %s, got %s", filepath.Base(tmpDir), state.Project)
	}
	if state.Artifacts == nil {
		t.Error("Artifacts map should be initialized")
	}
	if state.Decisions == nil {
		t.Error("Decisions map should be initialized")
	}
	if state.DependencyGraph == nil {
		t.Error("DependencyGraph should be initialized")
	}
}

func TestStateManager_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewStateManager(tmpDir)

	// Create state with data
	state := sm.NewState()
	state.Project = "test-project"
	state.LastFullDerive = time.Now().Truncate(time.Second)
	state.Artifacts["BR-001"] = &Artifact{
		ID:          "BR-001",
		Type:        ArtifactBusinessRule,
		Layer:       "l1",
		ContentHash: "sha256:abc123",
		Status:      StatusCurrent,
	}
	state.Decisions["DEC-001"] = &Decision{
		ID:       "DEC-001",
		Question: "Test question?",
		Answer:   "Test answer",
	}

	// Save state
	if err := sm.Save(state); err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(sm.StatePath); os.IsNotExist(err) {
		t.Fatal("State file was not created")
	}

	// Load state
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// Verify loaded data
	if loaded.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", loaded.Project)
	}
	if len(loaded.Artifacts) != 1 {
		t.Errorf("Expected 1 artifact, got %d", len(loaded.Artifacts))
	}
	if loaded.Artifacts["BR-001"] == nil {
		t.Error("Expected BR-001 artifact")
	}
	if loaded.Artifacts["BR-001"].ContentHash != "sha256:abc123" {
		t.Error("Artifact content hash mismatch")
	}
	if len(loaded.Decisions) != 1 {
		t.Errorf("Expected 1 decision, got %d", len(loaded.Decisions))
	}
}

func TestStateManager_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewStateManager(tmpDir)

	// Load from non-existent file should return new state
	state, err := sm.Load()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if state == nil {
		t.Fatal("Expected non-nil state")
	}
	if state.Version != StateVersion {
		t.Errorf("Expected version %s, got %s", StateVersion, state.Version)
	}
}

func TestStateManager_Lock(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewStateManager(tmpDir)

	// Lock should succeed
	if err := sm.Lock(); err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Should be locked
	if !sm.IsLocked() {
		t.Error("Expected IsLocked to be true")
	}

	// Lock file should exist
	if _, err := os.Stat(sm.LockPath); os.IsNotExist(err) {
		t.Error("Lock file was not created")
	}

	// Re-locking same manager should succeed (no-op)
	if err := sm.Lock(); err != nil {
		t.Fatalf("Re-lock failed: %v", err)
	}

	// Unlock
	if err := sm.Unlock(); err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	// Should not be locked
	if sm.IsLocked() {
		t.Error("Expected IsLocked to be false after unlock")
	}

	// Lock file should be removed
	if _, err := os.Stat(sm.LockPath); !os.IsNotExist(err) {
		t.Error("Lock file should be removed after unlock")
	}
}

func TestDerivationState_ArtifactOperations(t *testing.T) {
	state := &DerivationState{
		Artifacts: make(map[string]*Artifact),
	}

	// Set artifact
	artifact := &Artifact{
		ID:     "BR-001",
		Type:   ArtifactBusinessRule,
		Layer:  "l1",
		Status: StatusCurrent,
	}
	state.SetArtifact(artifact)

	// Get artifact
	got := state.GetArtifact("BR-001")
	if got == nil {
		t.Fatal("Expected to get artifact")
	}
	if got.ID != "BR-001" {
		t.Errorf("Expected ID 'BR-001', got '%s'", got.ID)
	}

	// Get non-existent
	if state.GetArtifact("BR-999") != nil {
		t.Error("Expected nil for non-existent artifact")
	}

	// Remove artifact
	state.RemoveArtifact("BR-001")
	if state.GetArtifact("BR-001") != nil {
		t.Error("Expected nil after removal")
	}
}

func TestDerivationState_GetArtifactsByLayer(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"BR-001": {ID: "BR-001", Layer: "l1"},
			"BR-002": {ID: "BR-002", Layer: "l1"},
			"TS-001": {ID: "TS-001", Layer: "l2"},
			"TC-001": {ID: "TC-001", Layer: "l3"},
		},
	}

	l1 := state.GetArtifactsByLayer("l1")
	if len(l1) != 2 {
		t.Errorf("Expected 2 L1 artifacts, got %d", len(l1))
	}

	l2 := state.GetArtifactsByLayer("l2")
	if len(l2) != 1 {
		t.Errorf("Expected 1 L2 artifact, got %d", len(l2))
	}

	l0 := state.GetArtifactsByLayer("l0")
	if len(l0) != 0 {
		t.Errorf("Expected 0 L0 artifacts, got %d", len(l0))
	}
}

func TestDerivationState_GetArtifactsByStatus(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"BR-001": {ID: "BR-001", Status: StatusCurrent},
			"BR-002": {ID: "BR-002", Status: StatusStale},
			"BR-003": {ID: "BR-003", Status: StatusStale},
			"TS-001": {ID: "TS-001", Status: StatusAffected},
		},
	}

	stale := state.GetArtifactsByStatus(StatusStale)
	if len(stale) != 2 {
		t.Errorf("Expected 2 stale artifacts, got %d", len(stale))
	}

	current := state.GetArtifactsByStatus(StatusCurrent)
	if len(current) != 1 {
		t.Errorf("Expected 1 current artifact, got %d", len(current))
	}
}

func TestDerivationState_GetStaleArtifacts(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"BR-001": {ID: "BR-001", Status: StatusCurrent},
			"BR-002": {ID: "BR-002", Status: StatusStale},
			"TS-001": {ID: "TS-001", Status: StatusAffected},
			"TC-001": {ID: "TC-001", Status: StatusModified},
		},
	}

	stale := state.GetStaleArtifacts()
	if len(stale) != 2 {
		t.Errorf("Expected 2 stale/affected artifacts, got %d", len(stale))
	}
}

func TestDerivationState_GetStatusReport(t *testing.T) {
	state := &DerivationState{
		Project:        "test-project",
		LastFullDerive: time.Now(),
		Artifacts: map[string]*Artifact{
			"US-001": {ID: "US-001", Layer: "l0", Status: StatusCurrent},
			"BR-001": {ID: "BR-001", Layer: "l1", Status: StatusCurrent},
			"BR-002": {ID: "BR-002", Layer: "l1", Status: StatusStale, Location: ArtifactLocation{File: "l1/br.md"}},
			"TS-001": {ID: "TS-001", Layer: "l2", Status: StatusAffected, Location: ArtifactLocation{File: "l2/ts.md"}},
			"TC-001": {ID: "TC-001", Layer: "l3", Status: StatusModified, Location: ArtifactLocation{File: "l3/tc.md"}, ManualSections: []string{"notes"}},
		},
	}

	report := state.GetStatusReport()

	if report.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", report.Project)
	}
	if report.TotalArtifacts != 5 {
		t.Errorf("Expected 5 total artifacts, got %d", report.TotalArtifacts)
	}

	// Check layer summaries
	if report.LayerSummary["l0"].Current != 1 {
		t.Error("Expected 1 current L0 artifact")
	}
	if report.LayerSummary["l1"].Current != 1 {
		t.Error("Expected 1 current L1 artifact")
	}
	if report.LayerSummary["l1"].Stale != 1 {
		t.Error("Expected 1 stale L1 artifact")
	}
	if report.LayerSummary["l2"].Affected != 1 {
		t.Error("Expected 1 affected L2 artifact")
	}
	if report.LayerSummary["l3"].Modified != 1 {
		t.Error("Expected 1 modified L3 artifact")
	}

	// Check lists
	if len(report.StaleArtifacts) != 1 {
		t.Errorf("Expected 1 stale artifact in list, got %d", len(report.StaleArtifacts))
	}
	if len(report.AffectedArtifacts) != 1 {
		t.Errorf("Expected 1 affected artifact in list, got %d", len(report.AffectedArtifacts))
	}
	if len(report.ModifiedArtifacts) != 1 {
		t.Errorf("Expected 1 modified artifact in list, got %d", len(report.ModifiedArtifacts))
	}
}

func TestDerivationState_DecisionOperations(t *testing.T) {
	state := &DerivationState{
		Decisions: make(map[string]*Decision),
	}

	// Set decision
	decision := &Decision{
		ID:       "DEC-001",
		Question: "Test?",
		Answer:   "Yes",
	}
	state.SetDecision(decision)

	// Get decision
	got := state.GetDecision("DEC-001")
	if got == nil {
		t.Fatal("Expected to get decision")
	}
	if got.Answer != "Yes" {
		t.Errorf("Expected answer 'Yes', got '%s'", got.Answer)
	}

	// Get non-existent
	if state.GetDecision("DEC-999") != nil {
		t.Error("Expected nil for non-existent decision")
	}
}

func TestStateManager_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewStateManager(tmpDir)

	state := sm.NewState()
	state.Project = "atomic-test"

	// Save should not leave temp files
	if err := sm.Save(state); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Check no .tmp file exists
	tmpPath := sm.StatePath + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("Temp file should not exist after successful save")
	}
}
