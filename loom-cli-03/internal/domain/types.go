// Package domain provides core domain types for loom-cli.
//
// Implements: l2/internal-api.md, l2/aggregate-design.md
// See: l1/domain-model.md
package domain

import "time"

// =============================================================================
// Exit Codes
// Implements: l2/tech-specs.md TS-ERR-001, l2/interface-contracts.md
// =============================================================================

const (
	ExitCodeSuccess  = 0   // Command completed successfully
	ExitCodeError    = 1   // General error
	ExitCodeQuestion = 100 // Interview: question available
)

// =============================================================================
// Configuration Constants
// Implements: DEC-L1-011, DEC-L1-016, l2/tech-specs.md
// =============================================================================

const (
	MaxGroupSize    = 5  // Maximum questions per group (DEC-L1-011)
	MaxPreviewLines = 20 // Preview rendering limit (DEC-L1-016)
	MaxLineWidth    = 80 // Terminal width (DEC-L1-016)
)

// =============================================================================
// Severity Type
// Implements: DEC-L1-003
// =============================================================================

type Severity string

const (
	SeverityCritical  Severity = "critical"  // Blocks derivation, must be resolved
	SeverityImportant Severity = "important" // Significant impact, should be resolved
	SeverityMinor     Severity = "minor"     // Nice to have, can use AI suggestion
)

// =============================================================================
// Analysis Types (AGG-ANL-001)
// Implements: l2/aggregate-design.md AGG-ANL-001, l2/internal-api.md
// =============================================================================

// Entity captures what was mentioned in L0, not final design.
//
// Implements: AGG-ANL-001, DEC-L1-005
type Entity struct {
	Name                string   `json:"name"`
	MentionedAttributes []string `json:"mentioned_attributes"`
	MentionedOperations []string `json:"mentioned_operations"`
	MentionedStates     []string `json:"mentioned_states"`
}

// Operation captures who does what, when, to what.
//
// Implements: AGG-ANL-001, DEC-L1-006
type Operation struct {
	Name            string   `json:"name"`
	Actor           string   `json:"actor"`
	Trigger         string   `json:"trigger"`
	Target          string   `json:"target"`
	MentionedInputs []string `json:"mentioned_inputs"`
	MentionedRules  []string `json:"mentioned_rules"`
}

// Relationship between entities.
//
// Implements: AGG-ANL-001
type Relationship struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`        // has_many, belongs_to, has_one, contains, references, many_to_many
	Cardinality string `json:"cardinality"` // 1:1, 1:N, N:1, N:M
}

// Aggregate represents a DDD aggregate from analysis.
//
// Implements: AGG-ANL-001
type Aggregate struct {
	Name     string   `json:"name"`
	Root     string   `json:"root"`
	Entities []string `json:"entities"`
}

// Domain is the complete L0 analysis result.
//
// Implements: AGG-ANL-001
type Domain struct {
	Entities      []Entity       `json:"entities"`
	Operations    []Operation    `json:"operations"`
	Relationships []Relationship `json:"relationships"`
	BusinessRules []string       `json:"business_rules"`
	UIMentions    []string       `json:"ui_mentions"`
}

// SkipCondition defines when a question should be skipped.
//
// Implements: AGG-INT-001, DEC-L1-008
type SkipCondition struct {
	QuestionID   string   `json:"question_id"`
	SkipIfAnswer []string `json:"skip_if_answer"`
}

// Ambiguity represents an unresolved question with skip logic.
//
// Implements: AGG-ANL-001, DEC-L1-007
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
// Implements: AGG-ANL-001
type AnalyzeResult struct {
	DomainModel  *Domain     `json:"domain_model"`
	Ambiguities  []Ambiguity `json:"ambiguities"`
	Decisions    []Decision  `json:"existing_decisions"`
	InputFiles   []string    `json:"input_files"`
	InputContent string      `json:"input_content"`
}

// =============================================================================
// Interview Types (AGG-INT-001)
// Implements: l2/aggregate-design.md AGG-INT-001, l2/internal-api.md
// =============================================================================

// Decision represents a resolved ambiguity.
//
// Implements: AGG-INT-001, DEC-L1-004
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
// Implements: AGG-INT-001, DEC-L1-012
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
// Implements: AGG-INT-001, DEC-L1-011
type QuestionGroup struct {
	ID        string      `json:"id"`
	Subject   string      `json:"subject"`
	Category  string      `json:"category"`
	Questions []Ambiguity `json:"questions"`
}

// InterviewOutput is what the CLI outputs for each question.
//
// Implements: AGG-INT-001, DEC-L1-011
type InterviewOutput struct {
	Status         string         `json:"status"` // question, group, complete, error
	Question       *Ambiguity     `json:"question,omitempty"`
	Group          *QuestionGroup `json:"group,omitempty"`
	Progress       string         `json:"progress,omitempty"`
	RemainingCount int            `json:"remaining_count"`
	SkippedCount   int            `json:"skipped_count"`
	Message        string         `json:"message,omitempty"`
}

// AnswerInput is the JSON input format for answering questions.
//
// Implements: TBL-OUT-003
type AnswerInput struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
	Source     string `json:"source"` // user, default, user_accepted_suggested
}

// =============================================================================
// Derivation Types
// Implements: l2/internal-api.md
// =============================================================================

// AcceptanceCriteria (L1 output).
//
// Implements: l2/internal-api.md
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

// BusinessRule (L1 output).
//
// Implements: l2/internal-api.md
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

// DerivationResult is the final L1 output.
//
// Implements: l2/internal-api.md
type DerivationResult struct {
	AcceptanceCriteria []AcceptanceCriteria `json:"acceptance_criteria"`
	BusinessRules      []BusinessRule       `json:"business_rules"`
	Decisions          []Decision           `json:"decisions"`
	Stats              DerivationStats      `json:"stats"`
}

// DerivationStats for reporting.
//
// Implements: l2/internal-api.md
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

// =============================================================================
// Cascade Types (AGG-CAS-001)
// Implements: l2/aggregate-design.md AGG-CAS-001
// =============================================================================

// PhaseState tracks the status of a single derivation phase.
//
// Implements: AGG-CAS-001
type PhaseState struct {
	Status    string    `json:"status"` // pending, running, completed, failed
	Timestamp time.Time `json:"timestamp,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// CascadeStateConfig holds immutable cascade configuration.
//
// Implements: AGG-CAS-001
type CascadeStateConfig struct {
	SkipInterview bool `json:"skip_interview"`
	Interactive   bool `json:"interactive"`
}

// CascadeState is the aggregate root for cascade derivation.
//
// Implements: AGG-CAS-001
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

// CascadeConfig is a value object for cascade command configuration.
//
// Implements: VO-001
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
// Implements: l2/aggregate-design.md AGG-VAL-001
// =============================================================================

// ValidationError represents a validation failure.
//
// Implements: AGG-VAL-001
type ValidationError struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"` // V001, V002, ...
	Message string `json:"message"`
	RefID   string `json:"ref_id,omitempty"`
}

// ValidationWarning represents a non-fatal validation issue.
//
// Implements: AGG-VAL-001
type ValidationWarning struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// ValidationCheck represents a single validation rule check result.
//
// Implements: AGG-VAL-001
type ValidationCheck struct {
	Rule    string `json:"rule"`
	Status  string `json:"status"` // pass, fail, skip
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// ValidationSummary provides aggregate validation statistics.
//
// Implements: AGG-VAL-001
type ValidationSummary struct {
	TotalChecks int `json:"total_checks"`
	Passed      int `json:"passed"`
	Failed      int `json:"failed"`
	Warnings    int `json:"warnings"`
	ErrorCount  int `json:"error_count"`
}

// ValidationResult is the aggregate root for validation output.
//
// Implements: AGG-VAL-001
type ValidationResult struct {
	Level    string              `json:"level"` // L1, L2, L3, ALL
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	Checks   []ValidationCheck   `json:"checks"`
	Summary  ValidationSummary   `json:"summary"`
}

// =============================================================================
// Interactive Mode Types
// Implements: VO-002, VO-003, DEC-031
// =============================================================================

// PhaseResult is the result of a single derivation phase for interactive approval.
//
// Implements: VO-002, DEC-031
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
// Implements: VO-003, DEC-031
type ApprovalAction int

const (
	ActionApprove     ApprovalAction = iota // Accept the generated content
	ActionEdit                              // Open in external editor
	ActionRegenerate                        // Regenerate with Claude
	ActionSkip                              // Skip this phase
	ActionQuit                              // Abort the cascade
)
