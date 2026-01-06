package derivation

import (
	"testing"
	"time"
)

func TestArtifactType_Layer(t *testing.T) {
	tests := []struct {
		artifactType ArtifactType
		expected     string
	}{
		// L0
		{ArtifactUserStory, "l0"},
		{ArtifactNFR, "l0"},
		{ArtifactVocabulary, "l0"},
		// L1
		{ArtifactEntity, "l1"},
		{ArtifactValueObject, "l1"},
		{ArtifactBusinessRule, "l1"},
		{ArtifactAcceptanceCrit, "l1"},
		{ArtifactBoundedContext, "l1"},
		// L2
		{ArtifactTechSpec, "l2"},
		{ArtifactInterfaceOp, "l2"},
		{ArtifactAggregateDesign, "l2"},
		{ArtifactSequence, "l2"},
		{ArtifactDataTable, "l2"},
		// L3
		{ArtifactTestCase, "l3"},
		{ArtifactAPIEndpoint, "l3"},
		{ArtifactCodeSkeleton, "l3"},
		{ArtifactTicket, "l3"},
		// Unknown
		{ArtifactType("invalid"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.artifactType), func(t *testing.T) {
			if got := tt.artifactType.Layer(); got != tt.expected {
				t.Errorf("Layer() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestArtifactStatus_IsActionRequired(t *testing.T) {
	tests := []struct {
		status   ArtifactStatus
		expected bool
	}{
		{StatusCurrent, false},
		{StatusModified, false},
		{StatusStale, true},
		{StatusAffected, true},
		{StatusNew, true},
		{StatusOrphaned, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsActionRequired(); got != tt.expected {
				t.Errorf("IsActionRequired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestArtifact_IsStale(t *testing.T) {
	t.Run("never derived", func(t *testing.T) {
		a := &Artifact{
			ID:                "BR-001",
			DerivedFromHashes: nil,
		}

		if !a.IsStale(map[string]string{"US-001": "hash1"}) {
			t.Error("Expected never-derived artifact to be stale")
		}
	})

	t.Run("upstream unchanged", func(t *testing.T) {
		a := &Artifact{
			ID: "BR-001",
			DerivedFromHashes: map[string]string{
				"US-001": "hash1",
				"US-002": "hash2",
			},
		}

		currentHashes := map[string]string{
			"US-001": "hash1",
			"US-002": "hash2",
		}

		if a.IsStale(currentHashes) {
			t.Error("Expected artifact to not be stale when upstream unchanged")
		}
	})

	t.Run("upstream changed", func(t *testing.T) {
		a := &Artifact{
			ID: "BR-001",
			DerivedFromHashes: map[string]string{
				"US-001": "hash1",
			},
		}

		currentHashes := map[string]string{
			"US-001": "hash1-modified",
		}

		if !a.IsStale(currentHashes) {
			t.Error("Expected artifact to be stale when upstream changed")
		}
	})

	t.Run("upstream deleted", func(t *testing.T) {
		a := &Artifact{
			ID: "BR-001",
			DerivedFromHashes: map[string]string{
				"US-001": "hash1",
			},
		}

		currentHashes := map[string]string{} // US-001 deleted

		if !a.IsStale(currentHashes) {
			t.Error("Expected artifact to be stale when upstream deleted")
		}
	})
}

func TestArtifact_HasManualEdits(t *testing.T) {
	t.Run("no manual sections", func(t *testing.T) {
		a := &Artifact{
			ID:             "BR-001",
			ManualSections: nil,
		}
		if a.HasManualEdits() {
			t.Error("Expected HasManualEdits to be false for nil ManualSections")
		}
	})

	t.Run("empty manual sections", func(t *testing.T) {
		a := &Artifact{
			ID:             "BR-001",
			ManualSections: []string{},
		}
		if a.HasManualEdits() {
			t.Error("Expected HasManualEdits to be false for empty ManualSections")
		}
	})

	t.Run("has manual sections", func(t *testing.T) {
		a := &Artifact{
			ID:             "BR-001",
			ManualSections: []string{"notes", "edge_cases"},
		}
		if !a.HasManualEdits() {
			t.Error("Expected HasManualEdits to be true")
		}
	})
}

func TestDecision(t *testing.T) {
	d := Decision{
		ID:        "DEC-L1-001",
		Layer:     "l1",
		Question:  "What is the minimum order amount?",
		Answer:    "$10.00",
		Source:    "user",
		DecidedAt: time.Now(),
		Affects:   []string{"BR-ORD-001", "AC-ORD-001"},
		Category:  "missing_definition",
		Subject:   "Order",
	}

	if d.ID != "DEC-L1-001" {
		t.Errorf("Expected ID 'DEC-L1-001', got '%s'", d.ID)
	}

	if len(d.Affects) != 2 {
		t.Errorf("Expected 2 affected artifacts, got %d", len(d.Affects))
	}
}

func TestDeriveOptions_Defaults(t *testing.T) {
	opts := DeriveOptions{}

	if opts.Interactive {
		t.Error("Interactive should default to false")
	}
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
	if opts.MergeStrategy != "" {
		t.Error("MergeStrategy should default to empty")
	}
}

func TestDeriveResult(t *testing.T) {
	result := DeriveResult{
		Success:      true,
		DerivedCount: 5,
		SkippedCount: 2,
		ErrorCount:   0,
		Derived:      []string{"BR-001", "BR-002", "AC-001", "AC-002", "AC-003"},
		Skipped:      []string{"BR-003", "AC-004"},
		Duration:     time.Second * 30,
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}
	if result.DerivedCount != 5 {
		t.Errorf("Expected DerivedCount 5, got %d", result.DerivedCount)
	}
	if result.SkippedCount != 2 {
		t.Errorf("Expected SkippedCount 2, got %d", result.SkippedCount)
	}
}

func TestLayerStatus(t *testing.T) {
	ls := LayerStatus{
		Current:  10,
		Stale:    2,
		Affected: 5,
		Modified: 1,
		New:      0,
		Orphaned: 0,
	}

	total := ls.Current + ls.Stale + ls.Affected + ls.Modified + ls.New + ls.Orphaned
	if total != 18 {
		t.Errorf("Expected total 18, got %d", total)
	}
}
