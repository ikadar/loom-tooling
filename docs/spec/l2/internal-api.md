---
title: "Loom CLI Internal API"
generated: 2025-01-03T15:45:00Z
status: draft
level: L2
---

# Loom CLI Internal API

## Overview

This document specifies the internal Go APIs between loom-cli packages.

**Traceability:** Implements [package-structure.md](package-structure.md) interfaces.

---

## internal/claude

### Client

**Purpose:** Claude API client via Claude Code CLI wrapper

```go
type Client struct {
    SessionID string  // For multi-turn conversations
    Verbose   bool    // Debug output
}

// NewClient creates a new Claude client
func NewClient() *Client

// Call sends a prompt and returns raw text response
func (c *Client) Call(prompt string) (string, error)

// CallWithSystemPrompt calls with additional system prompt
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error)

// CallJSON expects JSON response, auto-retries on parse failure
func (c *Client) CallJSON(prompt string, result interface{}) error

// CallJSONWithRetry calls with retry configuration
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error
```

### Response

```go
type Response struct {
    Result    string  `json:"result"`
    SessionID string  `json:"session_id"`
    CostUSD   float64 `json:"cost_usd"`
}
```

### RetryConfig

```go
type RetryConfig struct {
    MaxAttempts int           // Number of attempts (not retries)
    BaseDelay   time.Duration // Initial delay between attempts
    MaxDelay    time.Duration // Maximum delay cap
}

func DefaultRetryConfig() RetryConfig  // 3 attempts, 2s base, 30s max
```

### Helper Functions

```go
// extractJSON tries to parse JSON from response (handles markdown code blocks)
func extractJSON(response string, result interface{}) error

// sanitizeJSON fixes common LLM JSON output issues (unescaped newlines)
func sanitizeJSON(jsonStr string) string
```

---

## internal/config

### Config

**Purpose:** CLI configuration and argument parsing

```go
type Config struct {
    InputFile      string  // Single input file
    InputDir       string  // Input directory
    OutputDir      string  // Output directory
    DecisionsFile  string  // decisions.md path
    AnalysisFile   string  // analysis.json path
    VocabularyFile string  // Optional vocabulary
    NFRFile        string  // Optional NFR file
    Format         string  // "text" or "json"
    BatchMode      bool    // Non-interactive
    Verbose        bool    // Debug output
}
```

### Argument Parsers

```go
// ParseArgsForAnalyze parses analyze command arguments
// Requires: --input-file OR --input-dir
// Optional: --decisions, --verbose
func ParseArgsForAnalyze(args []string) (*Config, error)

// ParseArgsForDerive parses derive command arguments
// Requires: --output-dir
// Optional: --decisions, --analysis-file, --vocabulary, --nfr, --verbose
func ParseArgsForDerive(args []string) (*Config, error)

// ParseArgs is legacy alias for ParseArgsForDerive
func ParseArgs(args []string) (*Config, error)
```

### File Operations

```go
// ReadInputFiles reads all L0 input markdown files
// Returns: combined content, list of file paths, error
func (cfg *Config) ReadInputFiles() (string, []string, error)

// ReadVocabulary reads optional domain vocabulary file
func (cfg *Config) ReadVocabulary() (string, error)

// ReadNFR reads optional non-functional requirements file
func (cfg *Config) ReadNFR() (string, error)
```

---

## internal/domain

### Analysis Types

```go
// Entity from L0 analysis
type Entity struct {
    Name                string   `json:"name"`
    MentionedAttributes []string `json:"mentioned_attributes"`
    MentionedOperations []string `json:"mentioned_operations"`
    MentionedStates     []string `json:"mentioned_states"`
}

// Operation from L0 analysis
type Operation struct {
    Name            string   `json:"name"`
    Actor           string   `json:"actor"`
    Trigger         string   `json:"trigger"`
    Target          string   `json:"target"`
    MentionedInputs []string `json:"mentioned_inputs"`
    MentionedRules  []string `json:"mentioned_rules"`
}

// Relationship between entities
type Relationship struct {
    From        string `json:"from"`
    To          string `json:"to"`
    Type        string `json:"type"`        // contains, references, belongs_to, many_to_many
    Cardinality string `json:"cardinality"` // 1:1, 1:N, N:1, N:M
}

// Domain is the complete L0 analysis result
type Domain struct {
    Entities      []Entity       `json:"entities"`
    Operations    []Operation    `json:"operations"`
    Relationships []Relationship `json:"relationships"`
    BusinessRules []string       `json:"business_rules"`
    UIMentions    []string       `json:"ui_mentions"`
}
```

### Interview Types

```go
type Severity string
const (
    SeverityCritical  Severity = "critical"
    SeverityImportant Severity = "important"
    SeverityMinor     Severity = "minor"
)

// SkipCondition for conditional question skipping
type SkipCondition struct {
    QuestionID   string   `json:"question_id"`
    SkipIfAnswer []string `json:"skip_if_answer"`
}

// Ambiguity is an unresolved question
type Ambiguity struct {
    ID              string          `json:"id"`
    Category        string          `json:"category"`  // entity, operation, ui
    Subject         string          `json:"subject"`
    Question        string          `json:"question"`
    Severity        Severity        `json:"severity"`
    SuggestedAnswer string          `json:"suggested_answer,omitempty"`
    Options         []string        `json:"options,omitempty"`
    ChecklistItem   string          `json:"checklist_item"`
    DependsOn       []SkipCondition `json:"depends_on,omitempty"`
}

// Decision is a resolved ambiguity
type Decision struct {
    ID        string    `json:"id"`
    Question  string    `json:"question"`
    Answer    string    `json:"answer"`
    DecidedAt time.Time `json:"decided_at"`
    Source    string    `json:"source"`   // user, default, existing
    Category  string    `json:"category"`
    Subject   string    `json:"subject"`
}

// InterviewState holds ongoing interview state
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

// QuestionGroup for grouped questions (R6)
type QuestionGroup struct {
    ID        string      `json:"id"`
    Subject   string      `json:"subject"`
    Category  string      `json:"category"`
    Questions []Ambiguity `json:"questions"`
}

// InterviewOutput is CLI output for each question
type InterviewOutput struct {
    Status         string         `json:"status"` // question, group, complete, error
    Question       *Ambiguity     `json:"question,omitempty"`
    Group          *QuestionGroup `json:"group,omitempty"`
    Progress       string         `json:"progress,omitempty"`
    RemainingCount int            `json:"remaining_count"`
    SkippedCount   int            `json:"skipped_count"`
    Message        string         `json:"message,omitempty"`
}
```

### Derivation Types

```go
// AcceptanceCriteria (L1 output)
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

// BusinessRule (L1 output)
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

// DerivationResult is the final L1 output
type DerivationResult struct {
    AcceptanceCriteria []AcceptanceCriteria `json:"acceptance_criteria"`
    BusinessRules      []BusinessRule       `json:"business_rules"`
    Decisions          []Decision           `json:"decisions"`
    Stats              DerivationStats      `json:"stats"`
}

// DerivationStats for reporting
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
```

---

## internal/generator

### ChunkedTestCaseGenerator

**Purpose:** Generate test cases in batches to handle large AC sets

```go
type ChunkedTestCaseGenerator struct {
    Client    *claude.Client
    ChunkSize int  // Default: 5 ACs per chunk
}

// NewChunkedTestCaseGenerator creates generator with default settings
func NewChunkedTestCaseGenerator(client *claude.Client) *ChunkedTestCaseGenerator

// Generate generates test cases from AC markdown content
func (g *ChunkedTestCaseGenerator) Generate(acContent string) (*TestCaseResult, error)
```

### TestCaseResult

```go
type TestCaseResult struct {
    TestSuites []TestSuite `json:"test_suites"`
    Summary    TDAISummary `json:"summary"`
}

type TestSuite struct {
    ACRef   string     `json:"ac_ref"`
    ACTitle string     `json:"ac_title"`
    Tests   []TestCase `json:"tests"`
}

// FlattenTestCases extracts all test cases from suites
func FlattenTestCases(suites []TestSuite) []TestCase
```

### Helper Functions

```go
// buildPrompt injects content into <context> tag
func buildPrompt(promptTemplate, document string) string
```

---

## internal/formatter

### Test Case Formatting

```go
type TestCase struct {
    ID              string     `json:"id"`
    Name            string     `json:"name"`
    Category        string     `json:"category"` // positive, negative, boundary, hallucination
    ACRef           string     `json:"ac_ref"`
    BRRefs          []string   `json:"br_refs"`
    Preconditions   []string   `json:"preconditions"`
    TestData        []TestData `json:"test_data"`
    Steps           []string   `json:"steps"`
    ExpectedResults []string   `json:"expected_results"`
    ShouldNot       string     `json:"should_not,omitempty"`
}

type TestData struct {
    Field string      `json:"field"`
    Value interface{} `json:"value"`
    Notes string      `json:"notes"`
}

type TDAISummary struct {
    Total      int
    ByCategory struct {
        Positive, Negative, Boundary, Hallucination int
    }
    Coverage struct {
        ACsCovered            int
        PositiveRatio         float64
        NegativeRatio         float64
        HasHallucinationTests bool
    }
}

// FormatTestCases formats test cases as markdown
func FormatTestCases(suites []TestSuite, summary TDAISummary) string
```

### Tech Spec Formatting

```go
type TechSpec struct {
    ID               string          `json:"id"`
    Name             string          `json:"name"`
    BRRef            string          `json:"br_ref"`
    Rule             string          `json:"rule"`
    Implementation   string          `json:"implementation"`
    ValidationPoints []string        `json:"validation_points"`
    DataRequirements []DataReq       `json:"data_requirements"`
    ErrorHandling    []ErrorHandling `json:"error_handling"`
    RelatedACs       []string        `json:"related_acs"`
}

// FormatTechSpecs formats tech specs as markdown
func FormatTechSpecs(specs []TechSpec) string
```

### Interface Contract Formatting

```go
type InterfaceContract struct {
    ID                   string               `json:"id"`
    ServiceName          string               `json:"serviceName"`
    Purpose              string               `json:"purpose"`
    BaseURL              string               `json:"baseUrl"`
    Operations           []ContractOperation  `json:"operations"`
    Events               []ContractEvent      `json:"events"`
    SecurityRequirements SecurityRequirements `json:"securityRequirements"`
}

// FormatInterfaceContracts formats contracts as markdown
func FormatInterfaceContracts(contracts []InterfaceContract, sharedTypes []SharedType) string
```

### Aggregate Design Formatting

```go
type AggregateDesign struct {
    ID                 string           `json:"id"`
    Name               string           `json:"name"`
    Purpose            string           `json:"purpose"`
    Invariants         []AggInvariant   `json:"invariants"`
    Root               AggRoot          `json:"root"`
    Entities           []AggEntity      `json:"entities"`
    ValueObjects       []string         `json:"valueObjects"`
    Behaviors          []AggBehavior    `json:"behaviors"`
    Events             []AggEvent       `json:"events"`
    Repository         AggRepository    `json:"repository"`
    ExternalReferences []AggExternalRef `json:"externalReferences"`
}

// FormatAggregateDesign formats aggregate design as markdown
func FormatAggregateDesign(aggregates []AggregateDesign) string
```

### Sequence Design Formatting

```go
type SequenceDesign struct {
    ID           string              `json:"id"`
    Name         string              `json:"name"`
    Description  string              `json:"description"`
    Trigger      SequenceTrigger     `json:"trigger"`
    Participants []SeqParticipant    `json:"participants"`
    Steps        []SequenceStep      `json:"steps"`
    Outcome      SequenceOutcome     `json:"outcome"`
    Exceptions   []SequenceException `json:"exceptions"`
    RelatedACs   []string            `json:"relatedACs"`
    RelatedBRs   []string            `json:"relatedBRs"`
}

// FormatSequenceDesign formats sequences as markdown with Mermaid diagrams
func FormatSequenceDesign(sequences []SequenceDesign) string
```

### Data Model Formatting

```go
type DataTable struct {
    ID               string           `json:"id"`
    Name             string           `json:"name"`
    Aggregate        string           `json:"aggregate"`
    Purpose          string           `json:"purpose"`
    Fields           []DataField      `json:"fields"`
    PrimaryKey       DataPrimaryKey   `json:"primaryKey"`
    Indexes          []DataIndex      `json:"indexes"`
    ForeignKeys      []DataForeignKey `json:"foreignKeys"`
    CheckConstraints []DataConstraint `json:"checkConstraints"`
}

type DataEnum struct {
    Name   string   `json:"name"`
    Values []string `json:"values"`
}

// FormatDataModel formats data model as markdown
func FormatDataModel(tables []DataTable, enums []DataEnum) string
```

### Frontmatter

```go
// FormatFrontmatter generates YAML frontmatter
func FormatFrontmatter(title string, level string) string
```

### L4 Architecture Formatting

```go
type L4Architecture struct {
    Pattern            string           `json:"pattern"` // clean, hexagonal, layered
    Language           string           `json:"language"`
    Rationale          string           `json:"rationale"`
    Layers             []L4Layer        `json:"layers"`
    DependencyInjection L4DIConfig      `json:"dependency_injection"`
    ErrorHandling      L4ErrorConfig    `json:"error_handling"`
}

type L4Layer struct {
    Name        string   `json:"name"`
    Purpose     string   `json:"purpose"`
    AllowedDeps []string `json:"allowed_deps"`
    Packages    []string `json:"packages"`
}

type L4DIConfig struct {
    Pattern               string `json:"pattern"` // constructor, framework
    InterfacesLocation    string `json:"interfaces_location"`
    ImplementationsLocation string `json:"implementations_location"`
    Example               string `json:"example"`
}

type L4ErrorConfig struct {
    Pattern     string `json:"pattern"`
    ErrorType   string `json:"error_type"`
    Propagation string `json:"propagation"`
    CodePrefix  string `json:"code_prefix"`
}

// FormatL4Architecture formats architecture design as markdown
func FormatL4Architecture(arch L4Architecture, decisions []L4Decision, aggregateMapping []L4AggregateMapping) string
```

### L4 Patterns Formatting

```go
type L4Pattern struct {
    ID                  string   `json:"id"`
    Name                string   `json:"name"`
    Problem             string   `json:"problem"`
    SourceRefs          []string `json:"source_refs"`
    AppliesTo           []string `json:"applies_to"`
    Interface           L4Code   `json:"interface"`
    ImplementationNotes []string `json:"implementation_notes"`
    ExampleUsage        string   `json:"example_usage"`
}

type L4Code struct {
    Name     string `json:"name"`
    Language string `json:"language"`
    Code     string `json:"code"`
}

// FormatL4Patterns formats patterns design as markdown
func FormatL4Patterns(patterns []L4Pattern, languagePatterns []L4LanguagePattern) string
```

### L4 Coding Standards Formatting

```go
type L4CodingStandard struct {
    ID           string   `json:"id"`
    Category     string   `json:"category"` // naming, error_handling, testing, documentation
    Title        string   `json:"title"`
    Rule         string   `json:"rule"`
    GoodExamples []string `json:"good_examples"`
    BadExamples  []string `json:"bad_examples"`
    Rationale    string   `json:"rationale"`
}

type L4ErrorHandlingConfig struct {
    ErrorType      L4Code            `json:"error_type"`
    ErrorCodes     L4ErrorCodeConfig `json:"error_codes"`
    WrappingPattern string           `json:"wrapping_pattern"`
}

// FormatL4CodingStandards formats coding standards as markdown
func FormatL4CodingStandards(standards []L4CodingStandard, errorConfig L4ErrorHandlingConfig, tooling L4Tooling) string
```

### L4 Project Structure Formatting

```go
type L4ProjectStructure struct {
    Language       string              `json:"language"`
    RootStructure  L4RootStructure     `json:"root_structure"`
    PackageStructure []L4PackageStructure `json:"package_structure"`
    ServiceMapping []L4ServiceMapping  `json:"service_mapping"`
    TestStructure  L4TestStructure     `json:"test_structure"`
}

type L4RootStructure struct {
    Directories []L4Directory `json:"directories"`
    Files       []L4File      `json:"files"`
}

type L4Directory struct {
    Path           string   `json:"path"`
    Purpose        string   `json:"purpose"`
    Contains       []string `json:"contains"`
    Subdirectories []string `json:"subdirectories,omitempty"`
}

// FormatL4ProjectStructure formats project structure as markdown
func FormatL4ProjectStructure(structure L4ProjectStructure) string
```

### L4 Testing Strategy Formatting

```go
type L4TestingStrategy struct {
    Language       string                `json:"language"`
    Methodology    string                `json:"methodology"` // tdd, code-first
    TestFramework  L4TestFramework       `json:"test_framework"`
    TestPyramid    L4TestPyramid         `json:"test_pyramid"`
    TestCaseMapping []L4TestCaseMapping  `json:"test_case_mapping"`
    MockingStrategy L4MockingStrategy    `json:"mocking_strategy"`
    Fixtures       L4Fixtures            `json:"fixtures"`
    TDDWorkflow    *L4TDDWorkflow        `json:"tdd_workflow,omitempty"`
    Coverage       L4Coverage            `json:"coverage"`
}

type L4TestCaseMapping struct {
    L3TestCase  string `json:"l3_test_case"`
    ACRef       string `json:"ac_ref"`
    Category    string `json:"category"`
    TestType    string `json:"test_type"` // unit, integration, e2e
    File        string `json:"file"`
    Function    string `json:"function"`
    Description string `json:"description"`
}

type L4TDDWorkflow struct {
    Enabled       bool     `json:"enabled"`
    Steps         []string `json:"steps"`
    CIEnforcement L4CIConfig `json:"ci_enforcement"`
}

// FormatL4TestingStrategy formats testing strategy as markdown
func FormatL4TestingStrategy(strategy L4TestingStrategy) string
```

---

## internal/workflow

### Progress

```go
type Progress struct {
    Label   string
    Total   int
    Current int
}

// NewProgress creates a progress tracker
func NewProgress(label string, total int) *Progress

// Increment advances progress by 1
func (p *Progress) Increment()

// Done marks progress complete
func (p *Progress) Done()
```

### Approval

```go
// RequestApproval prompts user for yes/no
func RequestApproval(prompt string) (bool, error)
```

---

## prompts

### Embedded Prompts

```go
// All prompts embedded via go:embed directive
var (
    DomainDiscovery          string  // domain-discovery.md
    EntityAnalysis           string  // entity-analysis.md
    OperationAnalysis        string  // operation-analysis.md
    DerivationPrompt         string  // derivation.md
    InterviewPrompt          string  // interview.md
    DeriveL2                 string  // derive-l2.md
    DeriveTestCases          string  // derive-test-cases.md
    DeriveTechSpecs          string  // derive-tech-specs.md
    DeriveL3                 string  // derive-l3.md
    DeriveL3API              string  // derive-l3-api.md
    DeriveL3Skeletons        string  // derive-l3-skeletons.md
    DeriveDomainModel        string  // derive-domain-model.md
    DeriveBoundedContext     string  // derive-bounded-context.md
    DeriveInterfaceContracts string  // derive-interface-contracts.md
    DeriveSequenceDesign     string  // derive-sequence-design.md
    DeriveAggregateDesign    string  // derive-aggregate-design.md
    DeriveDataModel          string  // derive-data-model.md
    DeriveFeatureTickets     string  // derive-feature-tickets.md
    DeriveServiceBoundaries  string  // derive-service-boundaries.md
    DeriveEventDesign        string  // derive-event-design.md
    DeriveDependencyGraph    string  // derive-dependency-graph.md

    // L4 prompts
    DeriveL4Architecture     string  // derive-l4-architecture.md
    DeriveL4Patterns         string  // derive-l4-patterns.md
    DeriveL4CodingStandards  string  // derive-l4-coding-standards.md
    DeriveL4ProjectStructure string  // derive-l4-project-structure.md
    DeriveL4TestingStrategy  string  // derive-l4-testing-strategy.md
)
```

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L2 | [package-structure.md](package-structure.md) | Package structure |
| L2 | [tech-specs.md](tech-specs.md) | Technical specifications |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt specifications |
| L2 | This document | Internal API |
