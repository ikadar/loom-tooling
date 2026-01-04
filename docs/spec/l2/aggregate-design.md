---
title: "Loom CLI Aggregate Design"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Aggregate Design

## Overview

This document defines the aggregate design for loom-cli, mapping domain entities to Go struct implementations. In this CLI context, "aggregates" are coherent groups of structs that are persisted and manipulated together.

**Traceability:** Derived from [domain-model.md](../l1/domain-model.md) entities.

---

## Aggregates

### AGG-CAS-001: Cascade Aggregate

**Aggregate Root:** CascadeState

**Purpose:** Tracks the complete state of a cascade derivation pipeline.

**File:** `cmd/cascade.go`

```go
// CascadeState is the aggregate root for cascade derivation
type CascadeState struct {
    Version    string                 `json:"version"`
    InputHash  string                 `json:"input_hash"`
    Phases     map[string]*PhaseState `json:"phases"`
    Config     CascadeStateConfig     `json:"config"`
    Timestamps struct {
        Started   time.Time `json:"started"`
        Completed time.Time `json:"completed,omitempty"`
    } `json:"timestamps"`
}

type PhaseState struct {
    Status    string    `json:"status"`    // pending, running, completed, failed
    Timestamp time.Time `json:"timestamp,omitempty"`
    Error     string    `json:"error,omitempty"`
}

type CascadeStateConfig struct {
    SkipInterview bool `json:"skip_interview"`
    Interactive   bool `json:"interactive"`
}
```

**Invariants:**
- Version must be "1.0"
- Phases map must contain all 5 phases
- Status must be one of: pending, running, completed, failed
- InputHash is computed from input file content

**Lifecycle:**
1. Created by `newCascadeState()` when cascade starts
2. Updated after each phase completes
3. Persisted to `.cascade-state.json`
4. Loaded by `loadCascadeState()` for resume

**Related Entity:** E-DRV-001 (DerivationPipeline), E-DRV-002 (Phase)

---

### AGG-INT-001: Interview Aggregate

**Aggregate Root:** InterviewState

**Purpose:** Manages interview session with questions and recorded decisions.

**File:** `cmd/interview.go`, `internal/domain/types.go`

```go
// InterviewState is the aggregate root for interview sessions
type InterviewState struct {
    SessionID       string       `json:"session_id"`
    DomainModel     *Domain      `json:"domain_model"`
    Questions       []Ambiguity  `json:"questions"`
    Decisions       []Decision   `json:"decisions"`
    CurrentIndex    int          `json:"current_index"`
    Skipped         []string     `json:"skipped"`
    InputContent    string       `json:"input_content"`
    Complete        bool         `json:"complete"`
}

// Decision represents a resolved ambiguity
type Decision struct {
    ID         string    `json:"id"`
    Question   string    `json:"question"`
    Answer     string    `json:"answer"`
    DecidedAt  time.Time `json:"decided_at"`
    Source     string    `json:"source"`    // user, default, existing, user_accepted_suggested
    Category   string    `json:"category"`
    Subject    string    `json:"subject"`
}

// QuestionGroup for grouped interview mode
type QuestionGroup struct {
    ID        string      `json:"id"`
    Subject   string      `json:"subject"`
    Category  string      `json:"category"`
    Questions []Ambiguity `json:"questions"`
}

// InterviewOutput is what the CLI outputs for each question
type InterviewOutput struct {
    Status          string         `json:"status"`           // question, group, complete, error
    Question        *Ambiguity     `json:"question,omitempty"`
    Group           *QuestionGroup `json:"group,omitempty"`
    Progress        string         `json:"progress,omitempty"`
    RemainingCount  int            `json:"remaining_count"`
    SkippedCount    int            `json:"skipped_count"`
    Message         string         `json:"message,omitempty"`
}
```

**Invariants:**
- CurrentIndex must be 0 <= index <= len(Questions)
- Each Decision must reference existing Question ID
- Source must be one of: user, default, existing, user_accepted_suggested
- Status must be one of: question, group, complete, error

**Lifecycle:**
1. Created from analysis output (`--init`)
2. Questions presented one at a time (or grouped)
3. Decisions recorded with source attribution
4. Questions can be auto-skipped based on DependsOn conditions
5. Persisted to state file after each answer
6. Complete flag set when all questions answered

**Related Entity:** E-INT-001 (Interview), E-INT-002 (Question), E-INT-003 (Decision)

**Related DEC:** DEC-L1-004, DEC-L1-008, DEC-L1-011, DEC-L1-012

#### Automatic Dependency Inference

**File:** `cmd/interview.go:addDependencies()`

**Purpose:** Automatically infers question dependencies based on semantic patterns to reduce redundant questions.

**Process:**

```
┌─────────────────────────────────────────────────────────────┐
│ Phase 1: Build Capability Question Map                       │
├─────────────────────────────────────────────────────────────┤
│ For each question:                                           │
│   IF contains "can" AND "deleted" → map[subject+"_delete"]  │
│   IF contains "can" AND "modified" → map[subject+"_modify"] │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 2: Attach Dependencies to Follow-up Questions         │
├─────────────────────────────────────────────────────────────┤
│ For each question:                                           │
│   Match patterns → Add DependsOn with SkipConditions        │
└─────────────────────────────────────────────────────────────┘
```

**Dependency Pattern Matching:**

| Follow-up Pattern | Depends On | Skip If Answer Contains |
|-------------------|------------|-------------------------|
| `after delet*`, `when delet*`, `deletion cascade`, `upon deletion` | `{subject}_delete` capability | `cannot be deleted`, `no deletion`, `not deletable`, `cannot delete`, `no, `, `soft delete only` |
| `after modif*`, `when modif*`, `modification trigger` | `{subject}_modify` capability | `cannot be modified`, `immutable`, `no modification`, `cannot modify` |
| `when expir*`, `after expir*`, `expiration notification` | Expiration capability question (same subject, contains "have"/"support" + "expir") | `no expiration`, `does not expire`, `never expires` |

**Implementation:**
```go
func addDependencies(questions []domain.Ambiguity) []domain.Ambiguity {
    result := make([]domain.Ambiguity, len(questions))
    copy(result, questions)

    // Phase 1: Build capability question map
    questionMap := make(map[string]string) // keyword -> question ID
    for i := range result {
        q := &result[i]
        qLower := strings.ToLower(q.Question)

        if strings.Contains(qLower, "can") && strings.Contains(qLower, "deleted") {
            questionMap[q.Subject+"_delete"] = q.ID
        }
        if strings.Contains(qLower, "can") && strings.Contains(qLower, "modified") {
            questionMap[q.Subject+"_modify"] = q.ID
        }
    }

    // Phase 2: Add dependencies
    for i := range result {
        q := &result[i]
        qLower := strings.ToLower(q.Question)

        // Deletion follow-up
        if strings.Contains(qLower, "after delet") || /* ... */ {
            if depID, ok := questionMap[q.Subject+"_delete"]; ok {
                q.DependsOn = append(q.DependsOn, domain.SkipCondition{
                    QuestionID:   depID,
                    SkipIfAnswer: []string{"cannot be deleted", "no deletion", /*...*/},
                })
            }
        }
        // Similar for modification and expiration...
    }
    return result
}
```

#### Skip Condition Evaluation

**File:** `internal/interview/grouping.go:shouldSkipQuestion()`

**Purpose:** Determines if a question should be skipped based on previously answered dependent questions.

**Algorithm:**
```go
func shouldSkipQuestion(q *domain.Ambiguity, decisions []domain.Decision) bool {
    if len(q.DependsOn) == 0 {
        return false
    }

    for _, dep := range q.DependsOn {
        for _, d := range decisions {
            if d.ID == dep.QuestionID {
                // Check if answer matches any skip condition
                for _, skipPhrase := range dep.SkipIfAnswer {
                    if containsIgnoreCase(d.Answer, skipPhrase) {
                        return true
                    }
                }
            }
        }
    }
    return false
}
```

**Matching Rules:**
- Case-insensitive substring matching
- Any matching skip phrase triggers skip
- Multiple DependsOn conditions are OR'ed (any match skips)
- Empty DependsOn means never auto-skip

**Example Flow:**
```
Q1: "Can orders be deleted?" → User answers: "No, orders are immutable"
Q2: "What happens after order deletion?" → DependsOn: Q1, SkipIf: ["immutable"]
    → containsIgnoreCase("No, orders are immutable", "immutable") = true
    → Q2 is SKIPPED
```

---

### AGG-ANL-001: Analysis Aggregate

**Aggregate Root:** AnalyzeResult

**Purpose:** Contains discovered domain model and identified ambiguities.

**File:** `cmd/analyze.go`, `internal/domain/types.go`

```go
// AnalyzeResult is the aggregate root for analysis output
type AnalyzeResult struct {
    DomainModel   *Domain          `json:"domain_model"`
    Ambiguities   []Ambiguity      `json:"ambiguities"`
    Decisions     []Decision       `json:"existing_decisions"`
    InputFiles    []string         `json:"input_files"`
    InputContent  string           `json:"input_content"`
}

type Domain struct {
    Entities      []Entity       `json:"entities"`
    Operations    []Operation    `json:"operations"`
    Relationships []Relationship `json:"relationships"`
    BusinessRules []string       `json:"business_rules"`
    UIMentions    []string       `json:"ui_mentions"`
}

// Entity captures what was mentioned in L0, not final design
type Entity struct {
    Name                string   `json:"name"`
    MentionedAttributes []string `json:"mentioned_attributes"`
    MentionedOperations []string `json:"mentioned_operations"`
    MentionedStates     []string `json:"mentioned_states"`
}

// Operation captures who does what, when, to what
type Operation struct {
    Name            string   `json:"name"`
    Actor           string   `json:"actor"`
    Trigger         string   `json:"trigger"`
    Target          string   `json:"target"`
    MentionedInputs []string `json:"mentioned_inputs"`
    MentionedRules  []string `json:"mentioned_rules"`
}

type Relationship struct {
    From        string `json:"from"`
    To          string `json:"to"`
    Type        string `json:"type"`          // has_many, belongs_to, has_one
    Cardinality string `json:"cardinality"`   // 1:1, 1:N, N:M
}

// Ambiguity represents an unresolved question with skip logic
type Ambiguity struct {
    ID              string          `json:"id"`
    Category        string          `json:"category"`      // entity, operation, ui
    Subject         string          `json:"subject"`       // entity/operation name
    Question        string          `json:"question"`
    Severity        string          `json:"severity"`      // critical, important, minor
    SuggestedAnswer string          `json:"suggested_answer,omitempty"`
    Options         []string        `json:"options,omitempty"`
    ChecklistItem   string          `json:"checklist_item"`
    DependsOn       []SkipCondition `json:"depends_on,omitempty"`
}

// SkipCondition defines when a question should be skipped
type SkipCondition struct {
    QuestionID   string   `json:"question_id"`
    SkipIfAnswer []string `json:"skip_if_answer"`
}
```

**Invariants:**
- Entities must have unique names
- Relationships must reference existing entities
- Ambiguities must have unique IDs
- Severity must be one of: critical, important, minor

**Lifecycle:**
1. Created by analyze phase from L0 input
2. Ambiguities become interview questions
3. Domain model feeds into L1 derivation

**Related Entity:** E-DOC-001 (Document)

**Related DEC:** DEC-L1-003, DEC-L1-005, DEC-L1-006, DEC-L1-007, DEC-L1-008

---

### AGG-VAL-001: Validation Aggregate

**Aggregate Root:** ValidationResult

**Purpose:** Contains validation results for a set of documents.

**File:** `cmd/validate.go`

```go
// ValidationResult is the aggregate root for validation output
type ValidationResult struct {
    Level    string             `json:"level"`
    Errors   []ValidationError  `json:"errors"`
    Warnings []ValidationWarning `json:"warnings"`
    Checks   []ValidationCheck  `json:"checks"`
    Summary  ValidationSummary  `json:"summary"`
}

type ValidationError struct {
    File    string `json:"file"`
    Line    int    `json:"line,omitempty"`
    Rule    string `json:"rule"`         // V001, V002, ...
    Message string `json:"message"`
    RefID   string `json:"ref_id,omitempty"`
}

type ValidationWarning struct {
    File    string `json:"file"`
    Line    int    `json:"line,omitempty"`
    Rule    string `json:"rule"`
    Message string `json:"message"`
}

type ValidationCheck struct {
    Rule    string `json:"rule"`
    Status  string `json:"status"`       // pass, fail, skip
    Message string `json:"message"`
    Count   int    `json:"count,omitempty"`
}

type ValidationSummary struct {
    TotalChecks int `json:"total_checks"`
    Passed      int `json:"passed"`
    Failed      int `json:"failed"`
    Warnings    int `json:"warnings"`
    ErrorCount  int `json:"error_count"`
}
```

**Invariants:**
- Rule must be V001-V010
- Status must be one of: pass, fail, skip
- Summary counts must match actual error/warning counts

**Lifecycle:**
1. Created by validate command
2. Populated as each rule is checked
3. Output as JSON (--json) or text

**Related Entity:** E-VAL-001 (ValidationRun), E-VAL-002 (ValidationResult)

---

### AGG-L4-001: L4 Derivation Aggregate

**Aggregate Root:** L4DerivationState

**Purpose:** Tracks the state and outputs of L4 (Implementation Design) derivation.

**File:** `cmd/derive_l4.go`, `internal/domain/types.go`

```go
// L4DerivationState is the aggregate root for L4 derivation
type L4DerivationState struct {
    Language     string           `json:"language"`     // go, typescript, python
    Methodology  string           `json:"methodology"`  // tdd, code-first
    Config       L4Config         `json:"config"`
    Outputs      L4Outputs        `json:"outputs"`
    SourceRefs   L4SourceRefs     `json:"source_refs"`
    Timestamp    time.Time        `json:"timestamp"`
}

// L4Config from loom.config.yaml
type L4Config struct {
    Language       string           `yaml:"language"`
    Framework      string           `yaml:"framework"`
    Methodology    string           `yaml:"methodology"`
    Architecture   L4ArchConfig     `yaml:"architecture"`
    Testing        L4TestConfig     `yaml:"testing"`
    LanguageConfig map[string]any   `yaml:"-"` // go, typescript, python specific
}

type L4ArchConfig struct {
    Pattern string   `yaml:"pattern"` // clean, hexagonal, layered
    Layers  []string `yaml:"layers"`
}

type L4TestConfig struct {
    CoverageMinimum int  `yaml:"minimum"`
    NegativeRatio   int  `yaml:"negative_ratio"`
    Hallucination   bool `yaml:"hallucination_tests"`
}

// L4Outputs tracks which files were generated
type L4Outputs struct {
    Architecture    *L4FileOutput `json:"architecture,omitempty"`
    Patterns        *L4FileOutput `json:"patterns,omitempty"`
    CodingStandards *L4FileOutput `json:"coding_standards,omitempty"`
    ProjectStructure *L4FileOutput `json:"project_structure,omitempty"`
    TestingStrategy *L4FileOutput `json:"testing_strategy,omitempty"`
}

type L4FileOutput struct {
    File      string    `json:"file"`
    Generated time.Time `json:"generated"`
    ItemCount int       `json:"item_count"`
}

// L4SourceRefs tracks what L2/L3 documents were used as input
type L4SourceRefs struct {
    L2Aggregates  []string `json:"l2_aggregates"`
    L2Sequences   []string `json:"l2_sequences"`
    L2TechSpecs   []string `json:"l2_tech_specs"`
    L3TestCases   []string `json:"l3_test_cases"`
    L3Services    []string `json:"l3_services"`
    L3Skeletons   []string `json:"l3_skeletons"`
}
```

**Invariants:**
- Language must be one of: go, typescript, python
- Methodology must be one of: tdd, code-first
- Architecture pattern must be one of: clean, hexagonal, layered
- All L4 outputs must trace back to L2/L3 source documents

**Lifecycle:**
1. Created by derive-l4 command
2. Config loaded from loom.config.yaml
3. L2/L3 documents read as input
4. 5 L4 documents generated
5. State persisted to `.l4-state.json`

**Related Entity:** E-L4-001 (L4DerivationState)

**Related Commands:** IC-DRV-004, IC-GEN-001

---

## Value Objects

### VO-001: CascadeConfig

**Purpose:** Immutable configuration for cascade command.

```go
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
```

---

### VO-002: PhaseResult

**Purpose:** Result of a single derivation phase for interactive approval.

```go
type PhaseResult struct {
    PhaseName   string
    FileName    string
    Content     string
    Summary     string
    ItemCount   int
    ItemType    string
}
```

---

### VO-003: ApprovalAction

**Purpose:** User action in interactive mode.

```go
type ApprovalAction int

const (
    ActionApprove ApprovalAction = iota
    ActionEdit
    ActionRegenerate
    ActionSkip
    ActionQuit
)
```

**Related:** DEC-031

---

## Aggregate Relationships

```
┌─────────────────────────────────────────────────────────┐
│                    CASCADE FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  CascadeState ─────────────────────────────────────┐    │
│       │                                             │    │
│       ├── PhaseState["analyze"]                    │    │
│       │        │                                    │    │
│       │        └──> AnalysisResult                 │    │
│       │                   │                         │    │
│       ├── PhaseState["interview"]                  │    │
│       │        │                                    │    │
│       │        └──> InterviewState                 │    │
│       │                   │                         │    │
│       │                   ├── Question[]           │    │
│       │                   └── Decision[]           │    │
│       │                                             │    │
│       ├── PhaseState["derive-l1"]                  │    │
│       ├── PhaseState["derive-l2"]                  │    │
│       ├── PhaseState["derive-l3"]                  │    │
│       └── PhaseState["derive-l4"]                  │    │
│                   │                                 │    │
│                   └──> L4DerivationState           │    │
│                                                     │    │
│  ValidationResult ◄─────────────────────────────────┘    │
│       │                                                  │
│       ├── ValidationCheck[] (V001-V010)                 │
│       ├── ValidationError[]                             │
│       └── ValidationWarning[]                           │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Persistence Strategy

| Aggregate | Persistence | File |
|-----------|-------------|------|
| CascadeState | JSON file | `.cascade-state.json` |
| InterviewState | JSON file | `.interview-state.json` |
| AnalysisResult | JSON file | `.analysis.json` |
| L4DerivationState | JSON file | `.l4-state.json` |
| ValidationResult | stdout (JSON or text) | - |

**File Location:** `{output-dir}/` for cascade, current directory otherwise.

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model (source) |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt Catalog |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria |
| L2 | [tech-specs.md](tech-specs.md) | Technical Specifications |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI Interface Contract |
| L2 | This document | Aggregate Design |
| L2 | [sequence-design.md](sequence-design.md) | Sequence Design |
| L2 | [initial-data-model.md](initial-data-model.md) | Data Model |
