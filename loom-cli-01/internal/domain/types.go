// Package domain contains all domain types for loom-cli.
//
// Implements: l2/aggregate-design.md (AGG-ANL-001, AGG-INT-001, AGG-CAS-001, AGG-VAL-001)
// Implements: l2/internal-api.md
// See: l2/initial-data-model.md for JSON schemas
package domain

import "time"

// =============================================================================
// Analysis Types (AGG-ANL-001)
// =============================================================================

// Entity captures what was mentioned in L0, not final design.
// Implements: l2/aggregate-design.md AGG-ANL-001
type Entity struct {
	Name                string   `json:"name"`
	Classification      string   `json:"classification,omitempty"`       // entity, value_object, unknown
	Confidence          string   `json:"confidence,omitempty"`           // high, medium, low
	NeedsInterview      bool     `json:"needs_interview,omitempty"`
	InterviewQuestions  []string `json:"interview_questions,omitempty"`
	MentionedAttributes []string `json:"mentioned_attributes,omitempty"`
	MentionedOperations []string `json:"mentioned_operations,omitempty"`
	MentionedStates     []string `json:"mentioned_states,omitempty"`
}

// Operation captures who does what, when, to what.
// Implements: l2/aggregate-design.md AGG-ANL-001
type Operation struct {
	Name            string   `json:"name"`
	Actor           string   `json:"actor,omitempty"`
	Trigger         string   `json:"trigger,omitempty"`
	Target          string   `json:"target,omitempty"`
	MentionedInputs []string `json:"mentioned_inputs,omitempty"`
	MentionedRules  []string `json:"mentioned_rules,omitempty"`
}

// Relationship represents a relationship between entities.
// Implements: l2/aggregate-design.md AGG-ANL-001
type Relationship struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`                  // has_many, belongs_to, has_one, contains, references, many_to_many
	Cardinality string `json:"cardinality,omitempty"` // 1:1, 1:N, N:M
	Confidence  string `json:"confidence,omitempty"`  // high, medium, low
}

// Aggregate represents a DDD aggregate boundary.
type Aggregate struct {
	Name       string   `json:"name"`
	Root       string   `json:"root"`
	Contains   []string `json:"contains"`
	Confidence string   `json:"confidence,omitempty"` // high, medium, low
}

// Domain represents the discovered domain model.
// Implements: l2/aggregate-design.md AGG-ANL-001
type Domain struct {
	Entities      []Entity       `json:"entities"`
	Operations    []Operation    `json:"operations"`
	Relationships []Relationship `json:"relationships"`
	Aggregates    []Aggregate    `json:"aggregates,omitempty"`
	BusinessRules []string       `json:"business_rules,omitempty"`
	UIMentions    []string       `json:"ui_mentions,omitempty"`
}

// SkipCondition defines when a question should be skipped.
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-008 (automatic dependency inference)
type SkipCondition struct {
	QuestionID   string   `json:"question_id"`
	SkipIfAnswer []string `json:"skip_if_answer"`
}

// Ambiguity represents an unresolved question with skip logic.
// Implements: l2/aggregate-design.md AGG-ANL-001
type Ambiguity struct {
	ID              string          `json:"id"`
	Category        string          `json:"category"`                   // entity, operation, ui
	Subject         string          `json:"subject"`                    // entity/operation name
	Question        string          `json:"question"`
	Severity        string          `json:"severity"`                   // critical, important, minor
	SuggestedAnswer string          `json:"suggested_answer,omitempty"`
	Options         []string        `json:"options,omitempty"`
	ChecklistItem   string          `json:"checklist_item,omitempty"`
	DependsOn       []SkipCondition `json:"depends_on,omitempty"`
}

// AnalyzeResult is the aggregate root for analysis output.
// Implements: l2/aggregate-design.md AGG-ANL-001
type AnalyzeResult struct {
	DomainModel  *Domain     `json:"domain_model"`
	Ambiguities  []Ambiguity `json:"ambiguities"`
	Decisions    []Decision  `json:"existing_decisions,omitempty"`
	InputFiles   []string    `json:"input_files,omitempty"`
	InputContent string      `json:"input_content,omitempty"`
}

// =============================================================================
// Interview Types (AGG-INT-001)
// =============================================================================

// Decision represents a resolved ambiguity.
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-004 (decision source attribution)
type Decision struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	DecidedAt time.Time `json:"decided_at,omitempty"`
	Source    string    `json:"source"`             // user, default, existing, user_accepted_suggested
	Category  string    `json:"category,omitempty"`
	Subject   string    `json:"subject,omitempty"`
}

// InterviewState is the aggregate root for interview sessions.
// Implements: l2/aggregate-design.md AGG-INT-001
type InterviewState struct {
	SessionID    string      `json:"session_id"`
	DomainModel  *Domain     `json:"domain_model,omitempty"`
	Questions    []Ambiguity `json:"questions"`
	Decisions    []Decision  `json:"decisions"`
	CurrentIndex int         `json:"current_index"`
	Skipped      []string    `json:"skipped,omitempty"`
	InputContent string      `json:"input_content,omitempty"`
	Complete     bool        `json:"complete"`
}

// QuestionGroup for grouped interview mode.
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: DEC-L1-011 (question grouping)
type QuestionGroup struct {
	ID        string      `json:"id"`
	Subject   string      `json:"subject"`
	Category  string      `json:"category"`
	Questions []Ambiguity `json:"questions"`
}

// InterviewOutput is what the CLI outputs for each question.
// Implements: l2/aggregate-design.md AGG-INT-001
type InterviewOutput struct {
	Status         string         `json:"status"` // question, group, complete, error
	Question       *Ambiguity     `json:"question,omitempty"`
	Group          *QuestionGroup `json:"group,omitempty"`
	Progress       string         `json:"progress,omitempty"`
	RemainingCount int            `json:"remaining_count"`
	SkippedCount   int            `json:"skipped_count"`
	Message        string         `json:"message,omitempty"`
}

// AnswerInput is the JSON structure for --answer flag.
// Implements: l2/initial-data-model.md TBL-OUT-003
type AnswerInput struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
	Source     string `json:"source"` // user, default, existing, user_accepted_suggested
}

// =============================================================================
// Cascade Types (AGG-CAS-001)
// =============================================================================

// PhaseState represents the state of a single phase.
// Implements: l2/aggregate-design.md AGG-CAS-001
type PhaseState struct {
	Status    string    `json:"status"` // pending, running, completed, failed
	Timestamp time.Time `json:"timestamp,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// CascadeStateConfig holds cascade configuration.
// Implements: l2/aggregate-design.md AGG-CAS-001
type CascadeStateConfig struct {
	SkipInterview bool `json:"skip_interview"`
	Interactive   bool `json:"interactive"`
}

// CascadeState is the aggregate root for cascade derivation.
// Implements: l2/aggregate-design.md AGG-CAS-001
type CascadeState struct {
	Version    string                 `json:"version"` // Must be "1.0"
	InputHash  string                 `json:"input_hash"`
	Phases     map[string]*PhaseState `json:"phases"`
	Config     CascadeStateConfig     `json:"config"`
	Timestamps struct {
		Started   time.Time `json:"started"`
		Completed time.Time `json:"completed,omitempty"`
	} `json:"timestamps"`
}

// CascadeConfig is the command configuration (immutable value object).
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
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationError struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"` // V001, V002, ...
	Message string `json:"message"`
	RefID   string `json:"ref_id,omitempty"`
}

// ValidationWarning represents a validation warning.
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationWarning struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

// ValidationCheck represents a single validation check result.
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationCheck struct {
	Rule    string `json:"rule"`
	Status  string `json:"status"` // pass, fail, skip
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// ValidationSummary contains validation statistics.
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationSummary struct {
	TotalChecks int `json:"total_checks"`
	Passed      int `json:"passed"`
	Failed      int `json:"failed"`
	Warnings    int `json:"warnings"`
	ErrorCount  int `json:"error_count"`
}

// ValidationResult is the aggregate root for validation output.
// Implements: l2/aggregate-design.md AGG-VAL-001
type ValidationResult struct {
	Level    string              `json:"level"` // L1, L2, L3, ALL
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	Checks   []ValidationCheck   `json:"checks"`
	Summary  ValidationSummary   `json:"summary"`
}

// =============================================================================
// Interactive Mode Types
// =============================================================================

// PhaseResult is the result of a single derivation phase for interactive approval.
// Implements: l2/aggregate-design.md VO-002
// Implements: DEC-031 (interactive approval)
type PhaseResult struct {
	PhaseName string
	FileName  string
	Content   string
	Summary   string
	ItemCount int
	ItemType  string
}

// ApprovalAction is user action in interactive mode.
// Implements: l2/aggregate-design.md VO-003
// Implements: DEC-031 (interactive approval)
type ApprovalAction int

const (
	ActionApprove ApprovalAction = iota
	ActionEdit
	ActionRegenerate
	ActionSkip
	ActionQuit
)

// =============================================================================
// Exit Codes
// =============================================================================

// Exit codes as specified in l2/tech-specs.md TS-ERR-001.
const (
	ExitCodeSuccess  = 0
	ExitCodeError    = 1
	ExitCodeQuestion = 100 // Interview has more questions
)

// =============================================================================
// Constants
// =============================================================================

// Implements: DEC-L1-011, l2/tech-specs.md
const (
	MaxGroupSize    = 5  // Maximum questions per group
	MaxPreviewLines = 20 // Preview rendering limit (DEC-L1-016)
	MaxLineWidth    = 80 // Terminal width (DEC-L1-016)
)
