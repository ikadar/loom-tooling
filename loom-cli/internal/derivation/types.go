// Package derivation provides modular, incremental derivation capabilities
// for the Loom specification system.
package derivation

import (
	"time"
)

// =============================================================================
// Artifact Types
// =============================================================================

// ArtifactType represents the type of a derivation artifact
type ArtifactType string

const (
	// L0 Types - Input layer
	ArtifactUserStory  ArtifactType = "user_story"
	ArtifactNFR        ArtifactType = "nfr"
	ArtifactVocabulary ArtifactType = "vocabulary"

	// L1 Types - Domain Model layer
	ArtifactEntity         ArtifactType = "entity"
	ArtifactValueObject    ArtifactType = "value_object"
	ArtifactAggregate      ArtifactType = "aggregate"
	ArtifactRelationship   ArtifactType = "relationship"
	ArtifactBusinessRule   ArtifactType = "business_rule"
	ArtifactAcceptanceCrit ArtifactType = "acceptance_criteria"
	ArtifactBoundedContext ArtifactType = "bounded_context"

	// L2 Types - Technical Design layer
	ArtifactTechSpec        ArtifactType = "tech_spec"
	ArtifactInterfaceOp     ArtifactType = "interface_operation"
	ArtifactAggregateDesign ArtifactType = "aggregate_design"
	ArtifactSequence        ArtifactType = "sequence"
	ArtifactDataTable       ArtifactType = "data_table"
	ArtifactDataEnum        ArtifactType = "data_enum"

	// L3 Types - Implementation layer
	ArtifactTestCase     ArtifactType = "test_case"
	ArtifactAPIEndpoint  ArtifactType = "api_endpoint"
	ArtifactCodeSkeleton ArtifactType = "code_skeleton"
	ArtifactTicket       ArtifactType = "ticket"
	ArtifactEvent        ArtifactType = "event"
	ArtifactService      ArtifactType = "service"
)

// Layer returns the layer (l0, l1, l2, l3) for an artifact type
func (t ArtifactType) Layer() string {
	switch t {
	case ArtifactUserStory, ArtifactNFR, ArtifactVocabulary:
		return "l0"
	case ArtifactEntity, ArtifactValueObject, ArtifactAggregate,
		ArtifactRelationship, ArtifactBusinessRule, ArtifactAcceptanceCrit,
		ArtifactBoundedContext:
		return "l1"
	case ArtifactTechSpec, ArtifactInterfaceOp, ArtifactAggregateDesign,
		ArtifactSequence, ArtifactDataTable, ArtifactDataEnum:
		return "l2"
	case ArtifactTestCase, ArtifactAPIEndpoint, ArtifactCodeSkeleton,
		ArtifactTicket, ArtifactEvent, ArtifactService:
		return "l3"
	default:
		return "unknown"
	}
}

// =============================================================================
// Artifact Status
// =============================================================================

// ArtifactStatus represents the current state of an artifact
type ArtifactStatus string

const (
	// StatusCurrent means the artifact is up to date with its sources
	StatusCurrent ArtifactStatus = "current"

	// StatusStale means one or more upstream sources have changed
	StatusStale ArtifactStatus = "stale"

	// StatusAffected means the artifact depends on a stale artifact
	StatusAffected ArtifactStatus = "affected"

	// StatusModified means the artifact has manual edits since last derivation
	StatusModified ArtifactStatus = "modified"

	// StatusNew means the artifact has not been derived yet
	StatusNew ArtifactStatus = "new"

	// StatusOrphaned means the artifact's sources have been deleted
	StatusOrphaned ArtifactStatus = "orphaned"
)

// IsActionRequired returns true if the status requires user action
func (s ArtifactStatus) IsActionRequired() bool {
	switch s {
	case StatusStale, StatusAffected, StatusNew, StatusOrphaned:
		return true
	default:
		return false
	}
}

// =============================================================================
// Location
// =============================================================================

// ArtifactLocation specifies where an artifact is located in the file system
type ArtifactLocation struct {
	// File is the path to the file containing this artifact
	File string `json:"file"`

	// Anchor is the markdown anchor (e.g., "br-ord-001")
	Anchor string `json:"anchor,omitempty"`

	// LineStart is the starting line number (1-indexed)
	LineStart int `json:"line_start,omitempty"`

	// LineEnd is the ending line number (1-indexed, inclusive)
	LineEnd int `json:"line_end,omitempty"`
}

// =============================================================================
// Artifact
// =============================================================================

// Artifact represents a single derivable element in the specification
type Artifact struct {
	// ID is the unique identifier (e.g., "BR-ORD-001", "TS-BR-ORD-001")
	ID string `json:"id"`

	// Type is the artifact type
	Type ArtifactType `json:"type"`

	// Layer is the specification layer (l0, l1, l2, l3)
	Layer string `json:"layer"`

	// Location specifies where this artifact is in the file system
	Location ArtifactLocation `json:"location"`

	// ContentHash is the SHA256 hash of the artifact's content
	ContentHash string `json:"content_hash"`

	// Upstream maps upstream artifact IDs to their hashes at derivation time
	// This is used to detect when sources have changed
	Upstream map[string]string `json:"upstream,omitempty"`

	// Downstream lists artifact IDs that depend on this artifact
	Downstream []string `json:"downstream,omitempty"`

	// Decisions lists decision IDs that affected this artifact
	Decisions []string `json:"decisions,omitempty"`

	// DerivedAt is when this artifact was last derived
	DerivedAt time.Time `json:"derived_at,omitempty"`

	// DerivedFromHashes records the hashes of upstream artifacts at derivation time
	DerivedFromHashes map[string]string `json:"derived_from_hashes,omitempty"`

	// ManualSections lists the names of manually-edited sections
	ManualSections []string `json:"manual_sections,omitempty"`

	// Status is the current status of this artifact
	Status ArtifactStatus `json:"status"`
}

// IsStale returns true if any upstream artifact has changed since derivation
func (a *Artifact) IsStale(currentUpstreamHashes map[string]string) bool {
	if a.DerivedFromHashes == nil {
		return true // Never derived
	}

	for id, derivedHash := range a.DerivedFromHashes {
		currentHash, exists := currentUpstreamHashes[id]
		if !exists {
			return true // Upstream was deleted
		}
		if currentHash != derivedHash {
			return true // Upstream changed
		}
	}

	return false
}

// HasManualEdits returns true if the artifact has manual sections
func (a *Artifact) HasManualEdits() bool {
	return len(a.ManualSections) > 0
}

// =============================================================================
// Decision
// =============================================================================

// Decision represents a resolved ambiguity that affects derivation
type Decision struct {
	// ID is the unique identifier (e.g., "DEC-L1-015")
	ID string `json:"id"`

	// Layer is the specification layer this decision applies to
	Layer string `json:"layer"`

	// Question is the ambiguity question that was resolved
	Question string `json:"question"`

	// Answer is the chosen resolution
	Answer string `json:"answer"`

	// Source indicates how this decision was made
	// Values: "user", "default", "existing", "inferred"
	Source string `json:"source"`

	// DecidedAt is when this decision was made
	DecidedAt time.Time `json:"decided_at"`

	// Affects lists artifact IDs affected by this decision
	Affects []string `json:"affects,omitempty"`

	// Category is the ambiguity category
	Category string `json:"category,omitempty"`

	// Subject is what the decision is about (e.g., "Order", "Customer")
	Subject string `json:"subject,omitempty"`
}

// =============================================================================
// Dependency Edge
// =============================================================================

// EdgeType represents the type of dependency between artifacts
type EdgeType string

const (
	// EdgeDerives means the downstream artifact is derived from the upstream
	EdgeDerives EdgeType = "derives"

	// EdgeAffects means changes to upstream may affect downstream
	EdgeAffects EdgeType = "affects"

	// EdgeReferences means downstream references upstream but isn't derived from it
	EdgeReferences EdgeType = "references"
)

// DependencyEdge represents a dependency relationship between two artifacts
type DependencyEdge struct {
	// From is the upstream artifact ID
	From string `json:"from"`

	// To is the downstream artifact ID
	To string `json:"to"`

	// Type is the relationship type
	Type EdgeType `json:"type"`
}

// =============================================================================
// Status Report
// =============================================================================

// StatusReport provides an overview of the derivation state
type StatusReport struct {
	// Project is the project name
	Project string `json:"project"`

	// LastFullDerive is when a full derivation was last run
	LastFullDerive time.Time `json:"last_full_derive"`

	// LayerSummary maps layer names to their status counts
	LayerSummary map[string]*LayerStatus `json:"layer_summary"`

	// StaleArtifacts lists artifacts that need re-derivation
	StaleArtifacts []ArtifactSummary `json:"stale_artifacts,omitempty"`

	// AffectedArtifacts lists artifacts affected by stale artifacts
	AffectedArtifacts []ArtifactSummary `json:"affected_artifacts,omitempty"`

	// ModifiedArtifacts lists artifacts with manual edits
	ModifiedArtifacts []ArtifactSummary `json:"modified_artifacts,omitempty"`

	// TotalArtifacts is the total number of tracked artifacts
	TotalArtifacts int `json:"total_artifacts"`
}

// LayerStatus summarizes the status of artifacts in a layer
type LayerStatus struct {
	Current  int `json:"current"`
	Stale    int `json:"stale"`
	Affected int `json:"affected"`
	Modified int `json:"modified"`
	New      int `json:"new"`
	Orphaned int `json:"orphaned"`
}

// ArtifactSummary provides a brief summary of an artifact for reports
type ArtifactSummary struct {
	ID       string         `json:"id"`
	Type     ArtifactType   `json:"type"`
	Location string         `json:"location"`
	Status   ArtifactStatus `json:"status"`
	Message  string         `json:"message,omitempty"`
}

// =============================================================================
// Derive Options and Results
// =============================================================================

// DeriveOptions configures a derivation operation
type DeriveOptions struct {
	// Interactive enables interactive ambiguity resolution
	Interactive bool `json:"interactive"`

	// DryRun simulates derivation without making changes
	DryRun bool `json:"dry_run"`

	// StaleOnly only derives stale artifacts (not affected)
	StaleOnly bool `json:"stale_only"`

	// Layer limits derivation to a specific layer
	Layer string `json:"layer,omitempty"`

	// ArtifactIDs limits derivation to specific artifacts
	ArtifactIDs []string `json:"artifact_ids,omitempty"`

	// MergeStrategy specifies how to handle manual edits
	MergeStrategy MergeStrategy `json:"merge_strategy"`

	// MaxIterations limits ambiguity resolution iterations
	MaxIterations int `json:"max_iterations"`

	// Verbose enables detailed logging
	Verbose bool `json:"verbose"`
}

// MergeStrategy specifies how to merge generated content with existing content
type MergeStrategy string

const (
	// MergePreserveManual keeps manual sections, updates generated sections
	MergePreserveManual MergeStrategy = "preserve_manual"

	// MergeOverwriteAll replaces all content
	MergeOverwriteAll MergeStrategy = "overwrite_all"

	// MergeInteractive asks for each conflict
	MergeInteractive MergeStrategy = "interactive"

	// MergeReportOnly shows diff without making changes
	MergeReportOnly MergeStrategy = "report_only"
)

// DeriveResult contains the results of a derivation operation
type DeriveResult struct {
	// Success indicates if derivation completed without errors
	Success bool `json:"success"`

	// DerivedCount is the number of artifacts derived
	DerivedCount int `json:"derived_count"`

	// SkippedCount is the number of artifacts skipped (current or manual)
	SkippedCount int `json:"skipped_count"`

	// ErrorCount is the number of errors encountered
	ErrorCount int `json:"error_count"`

	// Derived lists the IDs of derived artifacts
	Derived []string `json:"derived,omitempty"`

	// Skipped lists the IDs of skipped artifacts
	Skipped []string `json:"skipped,omitempty"`

	// Errors lists any errors that occurred
	Errors []DeriveError `json:"errors,omitempty"`

	// NewDecisions lists decisions made during derivation
	NewDecisions []Decision `json:"new_decisions,omitempty"`

	// Duration is how long derivation took
	Duration time.Duration `json:"duration"`
}

// DeriveError represents an error during derivation
type DeriveError struct {
	ArtifactID string `json:"artifact_id"`
	Phase      string `json:"phase"`
	Message    string `json:"message"`
	Recoverable bool  `json:"recoverable"`
}
