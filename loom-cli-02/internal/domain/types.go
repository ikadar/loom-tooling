// Package domain provides shared domain types for loom-cli.
//
// Implements: l2/package-structure.md PKG-010
// Implements: l2/internal-api.md internal/domain
// Implements: l2/aggregate-design.md AGG-*
package domain

import "time"

// Exit codes as specified in l2/interface-contracts.md.
//
// Implements: l2/tech-specs.md TS-ERR-001
const (
	ExitCodeSuccess  = 0   // Success
	ExitCodeError    = 1   // General error
	ExitCodeQuestion = 100 // Interview: question available
)

// Constants as specified in l2/tech-specs.md.
//
// Implements: DEC-L1-011, DEC-L1-016
const (
	MaxGroupSize    = 5  // Maximum questions per group (DEC-L1-011)
	MaxPreviewLines = 20 // Preview rendering limit (DEC-L1-016)
	MaxLineWidth    = 80 // Terminal width (DEC-L1-016)
)

// Severity levels for ambiguities.
//
// Implements: DEC-L1-003
type Severity string

const (
	SeverityCritical  Severity = "critical"
	SeverityImportant Severity = "important"
	SeverityMinor     Severity = "minor"
)

// =============================================================================
// Analysis Types (AGG-ANL-001)
// =============================================================================

// Entity captures what was mentioned in L0, not final design.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
// Implements: DEC-L1-005
type Entity struct {
	Name                string   `json:"name"`
	MentionedAttributes []string `json:"mentioned_attributes"`
	MentionedOperations []string `json:"mentioned_operations"`
	MentionedStates     []string `json:"mentioned_states"`
}

// Operation captures who does what, when, to what.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
// Implements: DEC-L1-006
type Operation struct {
	Name            string   `json:"name"`
	Actor           string   `json:"actor"`
	Trigger         string   `json:"trigger"`
	Target          string   `json:"target"`
	MentionedInputs []string `json:"mentioned_inputs"`
	MentionedRules  []string `json:"mentioned_rules"`
}

// Relationship represents a relationship between entities.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
type Relationship struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`        // has_many, belongs_to, has_one, contains, references
	Cardinality string `json:"cardinality"` // 1:1, 1:N, N:1, N:M
}

// Domain is the complete L0 analysis result.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
// Implements: l2/internal-api.md internal/domain
type Domain struct {
	Entities      []Entity       `json:"entities"`
	Operations    []Operation    `json:"operations"`
	Relationships []Relationship `json:"relationships"`
	BusinessRules []string       `json:"business_rules"`
	UIMentions    []string       `json:"ui_mentions"`
}

// SkipCondition defines when a question should be skipped.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-008
type SkipCondition struct {
	QuestionID   string   `json:"question_id"`
	SkipIfAnswer []string `json:"skip_if_answer"`
}

// Ambiguity represents an unresolved question with skip logic.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
// Implements: DEC-L1-007
type Ambiguity struct {
	ID              string          `json:"id"`
	Category        string          `json:"category"` // entity, operation, ui
	Subject         string          `json:"subject"`  // entity/operation name
	Question        string          `json:"question"`
	Severity        Severity        `json:"severity"`
	SuggestedAnswer string          `json:"suggested_answer,omitempty"`
	Options         []string        `json:"options,omitempty"`
	ChecklistItem   string          `json:"checklist_item"`
	DependsOn       []SkipCondition `json:"depends_on,omitempty"`
}

// AnalyzeResult is the aggregate root for analysis output.
//
// Implements: l2/aggregate-design.md AGG-ANL-001
type AnalyzeResult struct {
	DomainModel  *Domain     `json:"domain_model"`
	Ambiguities  []Ambiguity `json:"ambiguities"`
	Decisions    []Decision  `json:"existing_decisions"`
	InputFiles   []string    `json:"input_files"`
	InputContent string      `json:"input_content"`
}

// =============================================================================
// Interview Types (AGG-INT-001)
// =============================================================================

// Decision represents a resolved ambiguity.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-004
type Decision struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	DecidedAt time.Time `json:"decided_at"`
	Source    string    `json:"source"` // user, default, existing, user_accepted_suggested
	Category  string    `json:"category"`
	Subject   string    `json:"subject"`
}

// InterviewState is the aggregate root for interview sessions.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-012
type InterviewState struct {
	SessionID    string      `json:"session_id"`
	DomainModel  *Domain     `json:"domain_model"`
	Questions    []Ambiguity `json:"questions"`
	Decisions    []Decision  `json:"decisions"`
	CurrentIndex int         `json:"current_index"`
	Skipped      []string    `json:"skipped"`
	InputContent string      `json:"input_content"`
	Complete     bool        `json:"complete"`
}

// QuestionGroup for grouped interview mode.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-011
type QuestionGroup struct {
	ID        string      `json:"id"`
	Subject   string      `json:"subject"`
	Category  string      `json:"category"`
	Questions []Ambiguity `json:"questions"`
}

// InterviewOutput is what the CLI outputs for each question.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-011
type InterviewOutput struct {
	Status         string         `json:"status"` // question, group, complete, error
	Question       *Ambiguity     `json:"question,omitempty"`
	Group          *QuestionGroup `json:"group,omitempty"`
	Progress       string         `json:"progress,omitempty"`
	RemainingCount int            `json:"remaining_count"`
	SkippedCount   int            `json:"skipped_count"`
	Message        string         `json:"message,omitempty"`
}

// AnswerInput is the JSON input for answering a question.
//
// Implements: l2/interface-contracts.md TBL-OUT-003
type AnswerInput struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
	Source     string `json:"source"` // user, default, user_accepted_suggested
}

// =============================================================================
// Cascade Types (AGG-CAS-001)
// =============================================================================

// PhaseState tracks the state of a single phase.
//
// Implements: l2/aggregate-design.md AGG-CAS-001
type PhaseState struct {
	Status    string    `json:"status"` // pending, running, completed, failed
	Timestamp time.Time `json:"timestamp,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// CascadeStateConfig holds cascade configuration within state.
//
// Implements: l2/aggregate-design.md AGG-CAS-001
type CascadeStateConfig struct {
	SkipInterview bool `json:"skip_interview"`
	Interactive   bool `json:"interactive"`
}

// CascadeState is the aggregate root for cascade derivation.
//
// Implements: l2/aggregate-design.md AGG-CAS-001
// Implements: DEC-L1-009
type CascadeState struct {
	Version    string                 `json:"version"` // "1.0"
	InputHash  string                 `json:"input_hash"`
	Phases     map[string]*PhaseState `json:"phases"`
	Config     CascadeStateConfig     `json:"config"`
	Timestamps struct {
		Started   time.Time `json:"started"`
		Completed time.Time `json:"completed,omitempty"`
	} `json:"timestamps"`
}

// CascadeConfig is the immutable configuration for cascade command.
//
// Implements: l2/aggregate-design.md VO-001
type CascadeConfig struct {
	InputFile     string
	InputDir      string
	OutputDir     string
	SkipInterview bool
	DecisionsFile string
	Interactive   bool
	Resume        bool
	FromLevel     string
}

// =============================================================================
// Validation Types (AGG-VAL-001)
// =============================================================================

// ValidationError represents a validation error.
//
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationError struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"` // V001, V002, ...
	Message string `json:"message"`
	RefID   string `json:"ref_id,omitempty"`
}

// ValidationWarning represents a validation warning.
//
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationWarning struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// ValidationCheck represents a single check result.
//
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationCheck struct {
	Rule    string `json:"rule"`
	Status  string `json:"status"` // pass, fail, skip
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// ValidationSummary contains validation statistics.
//
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationSummary struct {
	TotalChecks int `json:"total_checks"`
	Passed      int `json:"passed"`
	Failed      int `json:"failed"`
	Warnings    int `json:"warnings"`
	ErrorCount  int `json:"error_count"`
}

// ValidationResult is the aggregate root for validation output.
//
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationResult struct {
	Level    string              `json:"level"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	Checks   []ValidationCheck   `json:"checks"`
	Summary  ValidationSummary   `json:"summary"`
}

// =============================================================================
// Interactive Mode Types (VO-002, VO-003)
// =============================================================================

// PhaseResult is the result of a single derivation phase for interactive approval.
//
// Implements: l2/aggregate-design.md VO-002
// Implements: DEC-031
type PhaseResult struct {
	PhaseName string
	FileName  string
	Content   string
	Summary   string
	ItemCount int
	ItemType  string
}

// ApprovalAction represents user action in interactive mode.
//
// Implements: l2/aggregate-design.md VO-003
// Implements: DEC-031
type ApprovalAction int

const (
	ActionApprove ApprovalAction = iota
	ActionEdit
	ActionRegenerate
	ActionSkip
	ActionQuit
)

// =============================================================================
// Derivation Types (L1 Output)
// =============================================================================

// AcceptanceCriteria represents a derived acceptance criterion.
//
// Implements: l2/internal-api.md Derivation Types
type AcceptanceCriteria struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Given        string   `json:"given"`
	When         string   `json:"when"`
	Then         string   `json:"then"`
	ErrorCases   []string `json:"error_cases,omitempty"`
	SourceRefs   []string `json:"source_refs"`
	DecisionRefs []string `json:"decision_refs,omitempty"`
}

// BusinessRule represents a derived business rule.
//
// Implements: l2/internal-api.md Derivation Types
type BusinessRule struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Rule         string   `json:"rule"`
	Invariant    string   `json:"invariant"`
	Enforcement  string   `json:"enforcement"`
	ErrorCode    string   `json:"error_code,omitempty"`
	SourceRefs   []string `json:"source_refs"`
	DecisionRefs []string `json:"decision_refs,omitempty"`
}

// DerivationResult is the final L1 output aggregate.
//
// Implements: l2/internal-api.md Derivation Types
type DerivationResult struct {
	AcceptanceCriteria []AcceptanceCriteria `json:"acceptance_criteria"`
	BusinessRules      []BusinessRule       `json:"business_rules"`
	Decisions          []Decision           `json:"decisions"`
	Stats              DerivationStats      `json:"stats"`
}

// DerivationStats holds statistics about the derivation process.
//
// Implements: l2/internal-api.md Derivation Types
type DerivationStats struct {
	InputFiles          int `json:"input_files"`
	InputLines          int `json:"input_lines"`
	EntitiesAnalyzed    int `json:"entities_analyzed"`
	OperationsAnalyzed  int `json:"operations_analyzed"`
	AmbiguitiesFound    int `json:"ambiguities_found"`
	AmbiguitiesResolved int `json:"ambiguities_resolved"`
	ExistingDecisions   int `json:"existing_decisions"`
	NewDecisions        int `json:"new_decisions"`
	ACsGenerated        int `json:"acs_generated"`
	BRsGenerated        int `json:"brs_generated"`
}
