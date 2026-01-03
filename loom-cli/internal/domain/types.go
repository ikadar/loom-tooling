package domain

import "time"

// Entity represents a domain entity extracted from L0
type Entity struct {
	Name                string   `json:"name"`
	MentionedAttributes []string `json:"mentioned_attributes"`
	MentionedOperations []string `json:"mentioned_operations"`
	MentionedStates     []string `json:"mentioned_states"`
}

// Operation represents a domain operation extracted from L0
type Operation struct {
	Name            string   `json:"name"`
	Actor           string   `json:"actor"`
	Trigger         string   `json:"trigger"`
	Target          string   `json:"target"`
	MentionedInputs []string `json:"mentioned_inputs"`
	MentionedRules  []string `json:"mentioned_rules"`
}

// Relationship represents a relationship between entities
type Relationship struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`
	Cardinality string `json:"cardinality"`
}

// Domain represents the extracted domain model from L0
type Domain struct {
	Entities      []Entity       `json:"entities"`
	Operations    []Operation    `json:"operations"`
	Relationships []Relationship `json:"relationships"`
	BusinessRules []string       `json:"business_rules"`
	UIMentions    []string       `json:"ui_mentions"`
}

// Severity levels for ambiguities
type Severity string

const (
	SeverityCritical  Severity = "critical"
	SeverityImportant Severity = "important"
	SeverityMinor     Severity = "minor"
)

// SkipCondition defines when a question should be skipped
type SkipCondition struct {
	QuestionID   string   `json:"question_id"`    // ID of the question this depends on
	SkipIfAnswer []string `json:"skip_if_answer"` // Skip if answer contains any of these
}

// Ambiguity represents an unresolved question
type Ambiguity struct {
	ID              string          `json:"id"`
	Category        string          `json:"category"` // entity, operation, ui
	Subject         string          `json:"subject"`  // entity/operation name
	Question        string          `json:"question"`
	Severity        Severity        `json:"severity"`
	SuggestedAnswer string          `json:"suggested_answer,omitempty"`
	Options         []string        `json:"options,omitempty"`
	ChecklistItem   string          `json:"checklist_item"`
	DependsOn       []SkipCondition `json:"depends_on,omitempty"` // Skip conditions
}

// InterviewState holds the state of an ongoing interview
type InterviewState struct {
	SessionID       string       `json:"session_id"`
	DomainModel     *Domain      `json:"domain_model"`
	Questions       []Ambiguity  `json:"questions"`
	Decisions       []Decision   `json:"decisions"`        // answered so far
	CurrentIndex    int          `json:"current_index"`    // where we are
	Skipped         []string     `json:"skipped"`          // skipped question IDs
	InputContent    string       `json:"input_content"`    // original L0 content
	Complete        bool         `json:"complete"`         // interview done?
}

// QuestionGroup represents a group of related questions
type QuestionGroup struct {
	ID        string      `json:"id"`
	Subject   string      `json:"subject"`   // Common subject (e.g., "Order", "Customer")
	Category  string      `json:"category"`  // Common category (e.g., "entity", "operation")
	Questions []Ambiguity `json:"questions"` // Questions in this group
}

// InterviewOutput is what the CLI outputs for each question
type InterviewOutput struct {
	Status          string         `json:"status"`                    // "question", "group", "complete", "error"
	Question        *Ambiguity     `json:"question,omitempty"`        // Single question (legacy)
	Group           *QuestionGroup `json:"group,omitempty"`           // Grouped questions (R6)
	Progress        string         `json:"progress,omitempty"`        // e.g., "5/23"
	RemainingCount  int            `json:"remaining_count"`
	SkippedCount    int            `json:"skipped_count"`
	Message         string         `json:"message,omitempty"`
}

// Decision represents a resolved ambiguity
type Decision struct {
	ID         string    `json:"id"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	DecidedAt  time.Time `json:"decided_at"`
	Source     string    `json:"source"` // "user", "default", "existing", "user_accepted_suggested"
	Category   string    `json:"category"`
	Subject    string    `json:"subject"`
}

// AcceptanceCriteria represents a derived AC
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

// BusinessRule represents a derived BR
type BusinessRule struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Rule        string   `json:"rule"`
	Invariant   string   `json:"invariant"`
	Enforcement string   `json:"enforcement"`
	ErrorCode   string   `json:"error_code,omitempty"`
	SourceRefs  []string `json:"source_refs"`
	DecisionRefs []string `json:"decision_refs,omitempty"`
}

// DerivationResult is the final output
type DerivationResult struct {
	AcceptanceCriteria []AcceptanceCriteria `json:"acceptance_criteria"`
	BusinessRules      []BusinessRule       `json:"business_rules"`
	Decisions          []Decision           `json:"decisions"`
	Stats              DerivationStats      `json:"stats"`
}

// DerivationStats holds statistics about the derivation
type DerivationStats struct {
	InputFiles         int `json:"input_files"`
	InputLines         int `json:"input_lines"`
	EntitiesAnalyzed   int `json:"entities_analyzed"`
	OperationsAnalyzed int `json:"operations_analyzed"`
	AmbiguitiesFound   int `json:"ambiguities_found"`
	AmbiguitiesResolved int `json:"ambiguities_resolved"`
	ExistingDecisions  int `json:"existing_decisions"`
	NewDecisions       int `json:"new_decisions"`
	ACsGenerated       int `json:"acs_generated"`
	BRsGenerated       int `json:"brs_generated"`
}
