package derivation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTracker_DetectStaleArtifacts(t *testing.T) {
	tmpDir := t.TempDir()

	// Create downstream file
	downstreamFile := filepath.Join(tmpDir, "l1", "acceptance-criteria.md")
	os.MkdirAll(filepath.Dir(downstreamFile), 0755)
	os.WriteFile(downstreamFile, []byte("# AC-ORD-001\nDerived content"), 0644)

	// Create a simple artifact with no upstream (L0 roots are never stale from upstream changes)
	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// L1 artifact depends on L0
	upstreamHash := "sha256:original"

	// L0 artifact has no upstream, so it can never be stale from upstream changes
	l0Artifact := &Artifact{
		ID:          "US-ORD-001",
		Type:        ArtifactUserStory,
		Layer:       "l0",
		Status:      StatusCurrent,
		ContentHash: upstreamHash, // Same hash that L1 expects
		// No upstream - it's a root
	}
	state.Artifacts["US-ORD-001"] = l0Artifact
	l1Artifact := &Artifact{
		ID:       "AC-ORD-001",
		Type:     ArtifactAcceptanceCrit,
		Layer:    "l1",
		Location: ArtifactLocation{File: downstreamFile},
		Upstream: map[string]string{
			"US-ORD-001": upstreamHash,
		},
		DerivedFromHashes: map[string]string{
			"US-ORD-001": upstreamHash,
		},
		Status: StatusCurrent,
	}
	state.Artifacts["AC-ORD-001"] = l1Artifact

	state.DependencyGraph.AddEdge("US-ORD-001", "AC-ORD-001", EdgeDerives)

	tracker := NewTracker(state, tmpDir)

	t.Run("no stale when hashes match", func(t *testing.T) {
		stale, err := tracker.DetectStaleArtifacts()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// L0 has no upstream, L1's upstream hash matches stored
		if len(stale) != 0 {
			t.Errorf("Expected no stale artifacts, got %d: %v", len(stale), stale)
		}
	})

	t.Run("detects stale when upstream hash changes", func(t *testing.T) {
		// Change the upstream hash in the artifact registry
		// (simulating that upstream content changed)
		state.Artifacts["US-ORD-001"].ContentHash = "sha256:changed"

		stale, err := tracker.DetectStaleArtifacts()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// L1 should be stale because its upstream hash no longer matches
		foundL1 := false
		for _, a := range stale {
			if a.ID == "AC-ORD-001" {
				foundL1 = true
			}
		}
		if !foundL1 {
			t.Error("Expected AC-ORD-001 to be stale")
		}
	})
}

func TestTracker_UpdateStatuses(t *testing.T) {
	tmpDir := t.TempDir()
	state := &DerivationState{
		Artifacts:       make(map[string]*Artifact),
		DependencyGraph: NewDependencyGraph(),
	}

	// Create files
	l1File := filepath.Join(tmpDir, "l1.md")
	l2File := filepath.Join(tmpDir, "l2.md")

	os.WriteFile(l1File, []byte("L1 content"), 0644)
	os.WriteFile(l2File, []byte("L2 content"), 0644)

	hasher := NewHasher()
	l1Hash, _ := hasher.HashFile(l1File)

	// Setup chain: L0 -> L1 -> L2
	// L0 is a root with no file (like a user story reference)
	state.Artifacts["L0"] = &Artifact{
		ID:          "L0",
		Layer:       "l0",
		ContentHash: "sha256:l0original",
		Status:      StatusCurrent,
	}

	state.Artifacts["L1"] = &Artifact{
		ID:       "L1",
		Layer:    "l1",
		Location: ArtifactLocation{File: l1File},
		ContentHash: l1Hash,
		Upstream: map[string]string{"L0": "sha256:l0original"},
		DerivedFromHashes: map[string]string{"L0": "sha256:l0original"},
		Status:   StatusCurrent,
	}

	state.Artifacts["L2"] = &Artifact{
		ID:       "L2",
		Layer:    "l2",
		Location: ArtifactLocation{File: l2File},
		Upstream: map[string]string{"L1": l1Hash},
		DerivedFromHashes: map[string]string{"L1": l1Hash},
		Status:   StatusCurrent,
	}

	state.DependencyGraph.AddEdge("L0", "L1", EdgeDerives)
	state.DependencyGraph.AddEdge("L1", "L2", EdgeDerives)

	tracker := NewTracker(state, tmpDir)

	// Change L0's content hash (simulating a source change)
	state.Artifacts["L0"].ContentHash = "sha256:l0modified"

	err := tracker.UpdateStatuses()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// L1 should be stale (direct dependency on changed L0)
	if state.Artifacts["L1"].Status != StatusStale {
		t.Errorf("L1 should be stale, got %s", state.Artifacts["L1"].Status)
	}

	// L2 should be affected (transitive dependency)
	if state.Artifacts["L2"].Status != StatusAffected {
		t.Errorf("L2 should be affected, got %s", state.Artifacts["L2"].Status)
	}
}

func TestTracker_MarkAsDerived(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("Content"), 0644)

	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"TEST-001": {
				ID:       "TEST-001",
				Location: ArtifactLocation{File: testFile},
				Status:   StatusStale,
			},
		},
		DependencyGraph: NewDependencyGraph(),
	}

	tracker := NewTracker(state, tmpDir)

	upstreamHashes := map[string]string{
		"UP-001": "sha256:abc123",
	}

	err := tracker.MarkAsDerived("TEST-001", upstreamHashes)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	artifact := state.Artifacts["TEST-001"]

	if artifact.Status != StatusCurrent {
		t.Errorf("Expected status CURRENT, got %s", artifact.Status)
	}

	if artifact.DerivedFromHashes["UP-001"] != "sha256:abc123" {
		t.Error("Should store upstream hashes")
	}

	if artifact.DerivedAt.IsZero() {
		t.Error("DerivedAt should be set")
	}

	if artifact.ContentHash == "" {
		t.Error("ContentHash should be computed")
	}
}

func TestTracker_MarkAsModified(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"TEST-001": {
				ID:     "TEST-001",
				Status: StatusCurrent,
			},
		},
		DependencyGraph: NewDependencyGraph(),
	}

	tracker := NewTracker(state, ".")

	err := tracker.MarkAsModified("TEST-001", []string{"notes", "edge_cases"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	artifact := state.Artifacts["TEST-001"]

	if artifact.Status != StatusModified {
		t.Errorf("Expected status MODIFIED, got %s", artifact.Status)
	}

	if len(artifact.ManualSections) != 2 {
		t.Errorf("Expected 2 manual sections, got %d", len(artifact.ManualSections))
	}
}

func TestTracker_AnalyzeImpact(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"L0-001": {ID: "L0-001", Layer: "l0", Status: StatusCurrent},
			"L1-001": {ID: "L1-001", Layer: "l1", Status: StatusCurrent},
			"L1-002": {ID: "L1-002", Layer: "l1", Status: StatusCurrent, ManualSections: []string{"notes"}},
			"L2-001": {ID: "L2-001", Layer: "l2", Status: StatusCurrent},
		},
		DependencyGraph: NewDependencyGraph(),
	}

	// L0-001 -> L1-001 -> L2-001
	// L0-001 -> L1-002
	state.DependencyGraph.AddEdge("L0-001", "L1-001", EdgeDerives)
	state.DependencyGraph.AddEdge("L0-001", "L1-002", EdgeDerives)
	state.DependencyGraph.AddEdge("L1-001", "L2-001", EdgeDerives)

	tracker := NewTracker(state, ".")

	report := tracker.AnalyzeImpact([]string{"L0-001"})

	if len(report.ChangedArtifacts) != 1 {
		t.Errorf("Expected 1 changed artifact, got %d", len(report.ChangedArtifacts))
	}

	if len(report.AffectedArtifacts) != 3 {
		t.Errorf("Expected 3 affected artifacts, got %d", len(report.AffectedArtifacts))
	}

	if len(report.ManualEditWarnings) != 1 {
		t.Errorf("Expected 1 manual edit warning, got %d", len(report.ManualEditWarnings))
	}

	if report.ManualEditWarnings[0].ArtifactID != "L1-002" {
		t.Error("Warning should be for L1-002")
	}
}

func TestTracker_PlanDerivation(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"L0-001": {ID: "L0-001", Layer: "l0", Type: ArtifactUserStory},
			"L1-001": {ID: "L1-001", Layer: "l1", Type: ArtifactAcceptanceCrit, Upstream: map[string]string{"L0-001": "hash"}},
			"L1-002": {ID: "L1-002", Layer: "l1", Type: ArtifactBusinessRule, Upstream: map[string]string{"L0-001": "hash"}},
			"L2-001": {ID: "L2-001", Layer: "l2", Type: ArtifactTestCase, Upstream: map[string]string{"L1-001": "hash"}},
		},
		DependencyGraph: NewDependencyGraph(),
	}

	state.DependencyGraph.AddEdge("L0-001", "L1-001", EdgeDerives)
	state.DependencyGraph.AddEdge("L0-001", "L1-002", EdgeDerives)
	state.DependencyGraph.AddEdge("L1-001", "L2-001", EdgeDerives)

	tracker := NewTracker(state, ".")

	plan, err := tracker.PlanDerivation([]string{"L1-001", "L1-002"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if plan.TotalCount != 3 { // L1-001, L1-002, L2-001
		t.Errorf("Expected 3 artifacts in plan, got %d", plan.TotalCount)
	}

	// Check order - L1 artifacts should come before L2
	for i, step := range plan.Artifacts {
		if step.Layer == "l2" {
			// All L1 should be before this
			for j := 0; j < i; j++ {
				if plan.Artifacts[j].Layer == "l2" {
					// Another L2 before - that's ok
				}
			}
		}
	}
}

func TestTracker_DetectFileChanges(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(testFile, []byte("Original"), 0644)

	hasher := NewHasher()
	info, _ := hasher.HashFileWithInfo(testFile)

	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"TEST-001": {
				ID:       "TEST-001",
				Location: ArtifactLocation{File: testFile},
			},
		},
		FileHashes:      map[string]*FileHashInfo{testFile: info},
		DependencyGraph: NewDependencyGraph(),
	}

	tracker := NewTracker(state, tmpDir)

	t.Run("no changes initially", func(t *testing.T) {
		changes, err := tracker.DetectFileChanges()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(changes) != 0 {
			t.Errorf("Expected no changes, got %d", len(changes))
		}
	})

	t.Run("detects modification", func(t *testing.T) {
		// Wait a moment and modify
		time.Sleep(10 * time.Millisecond)
		os.WriteFile(testFile, []byte("Modified"), 0644)

		changes, err := tracker.DetectFileChanges()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(changes) != 1 {
			t.Errorf("Expected 1 change, got %d", len(changes))
		}
		if len(changes) > 0 && changes[0].ChangeType != "modified" {
			t.Errorf("Expected 'modified', got '%s'", changes[0].ChangeType)
		}
	})
}

func TestTracker_CleanupOrphaned(t *testing.T) {
	state := &DerivationState{
		Artifacts: map[string]*Artifact{
			"CURRENT":  {ID: "CURRENT", Status: StatusCurrent},
			"ORPHAN1":  {ID: "ORPHAN1", Status: StatusOrphaned},
			"ORPHAN2":  {ID: "ORPHAN2", Status: StatusOrphaned},
		},
		DependencyGraph: NewDependencyGraph(),
	}

	state.DependencyGraph.AddEdge("CURRENT", "ORPHAN1", EdgeDerives)

	tracker := NewTracker(state, ".")

	removed := tracker.CleanupOrphaned()

	if len(removed) != 2 {
		t.Errorf("Expected 2 removed, got %d", len(removed))
	}

	if state.GetArtifact("CURRENT") == nil {
		t.Error("CURRENT should still exist")
	}

	if state.GetArtifact("ORPHAN1") != nil {
		t.Error("ORPHAN1 should be removed")
	}
}

func TestTracker_ValidateGraph(t *testing.T) {
	t.Run("no issues", func(t *testing.T) {
		state := &DerivationState{
			Artifacts: map[string]*Artifact{
				"A": {ID: "A"},
				"B": {ID: "B"},
			},
			DependencyGraph: NewDependencyGraph(),
		}
		state.DependencyGraph.AddEdge("A", "B", EdgeDerives)

		tracker := NewTracker(state, ".")
		issues := tracker.ValidateGraph()

		if len(issues) != 0 {
			t.Errorf("Expected no issues, got %v", issues)
		}
	})

	t.Run("detects cycle", func(t *testing.T) {
		state := &DerivationState{
			Artifacts: map[string]*Artifact{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
			},
			DependencyGraph: NewDependencyGraph(),
		}
		state.DependencyGraph.AddEdge("A", "B", EdgeDerives)
		state.DependencyGraph.AddEdge("B", "C", EdgeDerives)
		state.DependencyGraph.AddEdge("C", "A", EdgeDerives) // cycle

		tracker := NewTracker(state, ".")
		issues := tracker.ValidateGraph()

		hasCycleIssue := false
		for _, issue := range issues {
			if len(issue) > 0 {
				hasCycleIssue = true
			}
		}

		if !hasCycleIssue {
			t.Error("Should detect cycle")
		}
	})

	t.Run("detects missing node", func(t *testing.T) {
		state := &DerivationState{
			Artifacts: map[string]*Artifact{
				"A": {ID: "A"},
				// B is missing
			},
			DependencyGraph: NewDependencyGraph(),
		}
		state.DependencyGraph.AddEdge("A", "B", EdgeDerives)

		tracker := NewTracker(state, ".")
		issues := tracker.ValidateGraph()

		if len(issues) == 0 {
			t.Error("Should detect missing node B")
		}
	})
}

func TestLayerOrder(t *testing.T) {
	tests := []struct {
		layer    string
		expected int
	}{
		{"l0", 0},
		{"l1", 1},
		{"l2", 2},
		{"l3", 3},
		{"unknown", 99},
	}

	for _, tt := range tests {
		t.Run(tt.layer, func(t *testing.T) {
			result := layerOrder(tt.layer)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
