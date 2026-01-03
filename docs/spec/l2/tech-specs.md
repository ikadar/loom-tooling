---
title: "Loom CLI Technical Specifications"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Technical Specifications

## Overview

This document defines the technical specifications for loom-cli implementation. It provides detailed technical requirements derived from L1 business rules and acceptance criteria.

**Traceability:** Derived from L1 documents (domain-model.md, business-rules.md, acceptance-criteria.md).

---

## Technology Stack

### TS-CLI-001: Core Technology

**Language:** Go 1.21+

**Rationale:**
- Single binary distribution (no runtime dependencies)
- Fast startup time (CLI responsiveness)
- Strong typing for reliability
- Excellent concurrency support (parallel phase execution)

**Related BR:** BR-DRV-002 (parallel execution)

**Related DEC:** DEC-L1-014

---

### TS-CLI-002: External Dependencies

| Dependency | Purpose | Version |
|------------|---------|---------|
| Claude Code CLI | AI-powered derivation | `claude` command |
| Standard library | File I/O, JSON | Go stdlib |

**Rationale:** Minimal dependencies for maintainability and security. Uses Claude Code CLI for AI integration instead of direct API calls, enabling session management and tool use.

**Related NFR:** NFR-MNT-002 (minimal dependencies)

**Related DEC:** DEC-L1-001

---

## Architecture Specifications

### TS-ARCH-001: Command Structure

**Pattern:** Command Router with Function Handlers

```
main() → parseArgs() → router → commandHandler()
```

**Commands:**
| Command | Handler Function | File |
|---------|------------------|------|
| analyze | runAnalyze() | analyze.go |
| interview | runInterview() | interview.go |
| derive | runDeriveNew() | derive_new.go |
| derive-l2 | runDeriveL2() | derive_l2.go |
| derive-l3 | runDeriveL3() | derive_l3.go |
| validate | runValidate() | validate.go |
| sync-links | runSyncLinks() | sync_links.go |
| cascade | runCascade() | cascade.go |
| help | printHelp() | root.go |
| version | printVersion() | root.go |

**Related:** interface-contracts.md

---

### TS-ARCH-001a: Analyze Command Flow

**File:** `cmd/analyze.go`

**Purpose:** Discovers domain model and identifies ambiguities from L0 input.

**Execution Phases:**

```
┌─────────────────────────────────────────────────────────────┐
│ Phase 0: Read Input                                          │
├─────────────────────────────────────────────────────────────┤
│ - Read markdown files (--input-file or --input-dir)         │
│ - Load existing decisions (if --decisions provided)          │
│ - Output: inputContent string, existingDecisions []Decision  │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 1: Domain Discovery                                    │
├─────────────────────────────────────────────────────────────┤
│ Function: discoverDomain(client, inputContent)               │
│ Prompt: prompts.DomainDiscovery                              │
│ Output: Domain{Entities, Operations, Relationships, Rules}   │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 2: Completeness Analysis                               │
├─────────────────────────────────────────────────────────────┤
│ Function: analyzeEntities(client, entities)                  │
│ Prompt: prompts.EntityAnalysis                               │
│ Output: entityAmbiguities []Ambiguity                        │
│                                                              │
│ Function: analyzeOperations(client, operations)              │
│ Prompt: prompts.OperationAnalysis                            │
│ Output: operationAmbiguities []Ambiguity                     │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 3: Filter Already Resolved                             │
├─────────────────────────────────────────────────────────────┤
│ - Compare ambiguities against existingDecisions              │
│ - Remove questions already answered                          │
│ - Output: filteredAmbiguities []Ambiguity                    │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Output: JSON to stdout                                       │
├─────────────────────────────────────────────────────────────┤
│ {                                                            │
│   "domain_model": {...},                                     │
│   "ambiguities": [...],                                      │
│   "existing_decisions": [...],                               │
│   "input_files": [...],                                      │
│   "input_content": "..."                                     │
│ }                                                            │
└─────────────────────────────────────────────────────────────┘
```

**API Calls:** 3 (domain discovery, entity analysis, operation analysis)

**Related:** AC-ANL-001, AC-ANL-002, prompts/domain-discovery.md, prompts/entity-analysis.md, prompts/operation-analysis.md

---

### TS-ARCH-001b: Derive Command Flow

**File:** `cmd/derive_new.go`

**Purpose:** Generates L1 Strategic Design documents from analysis.

**Execution Phases:**

```
┌─────────────────────────────────────────────────────────────┐
│ Phase 0: Load Input                                          │
├─────────────────────────────────────────────────────────────┤
│ - Load analysis JSON (--analysis-file)                       │
│ - Load optional vocabulary (--vocabulary)                    │
│ - Load optional NFR (--nfr)                                  │
│ - Create output directory                                    │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 1: Generate Domain Model                               │
├─────────────────────────────────────────────────────────────┤
│ Prompt: prompts.DeriveDomainModel                            │
│ Context: domain_model + vocabulary                           │
│ Output: domain-model.md                                      │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 2: Generate Bounded Context Map                        │
├─────────────────────────────────────────────────────────────┤
│ Prompt: prompts.DeriveBoundedContext                         │
│ Context: domain_model + entities                             │
│ Output: bounded-context-map.md                               │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 3: Generate Acceptance Criteria & Business Rules       │
├─────────────────────────────────────────────────────────────┤
│ Prompt: prompts.DerivationPrompt                             │
│ Context: domain_model + decisions + nfr                      │
│ Output: acceptance-criteria.md, business-rules.md            │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Phase 4: Write Decisions                                     │
├─────────────────────────────────────────────────────────────┤
│ - Merge existing decisions with new decisions                │
│ - Format as markdown                                         │
│ Output: decisions.md                                         │
└─────────────────────────────────────────────────────────────┘
```

**API Calls:** 3-4 (domain model, bounded context, AC+BR derivation)

**Output Files:**
- `{output-dir}/domain-model.md`
- `{output-dir}/bounded-context-map.md`
- `{output-dir}/acceptance-criteria.md`
- `{output-dir}/business-rules.md`
- `{output-dir}/decisions.md`

**Related:** AC-DRV-001, prompts/derive-domain-model.md, prompts/derive-bounded-context.md, prompts/derivation.md

---

### TS-ARCH-001c: Interview Command Flow

**File:** `cmd/interview.go`, `internal/interview/grouping.go`

**Purpose:** Conducts structured interview to resolve ambiguities.

**Modes:**

| Mode | Flag | Description |
|------|------|-------------|
| Initialize | `--init` | Create new interview from analysis JSON |
| Single Question | (default) | One question at a time |
| Grouped | `--grouped` | All questions grouped by subject |
| Batch Answer | `--answers` | Answer multiple questions at once |

**Exit Codes:**
| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitCodeComplete` | Interview complete, no more questions |
| 1 | `ExitCodeError` | Error occurred |
| 100 | `ExitCodeQuestion` | Question available (output contains question JSON) |

**Question Skip Logic:**

```
┌─────────────────────────────────────────────────────────────┐
│ SkipCondition                                                │
├─────────────────────────────────────────────────────────────┤
│ type SkipCondition struct {                                  │
│     QuestionID   string   // ID of question this depends on  │
│     SkipIfAnswer []string // Skip if answer contains these   │
│ }                                                            │
└─────────────────────────────────────────────────────────────┘
```

**Automatic Dependency Detection:**

The `addDependencies()` function detects common patterns:

| Pattern in Question | Depends On | Skip If Answer Contains |
|---------------------|------------|-------------------------|
| "after deletion", "when deleted" | "{Subject}_delete" question | "cannot be deleted", "no deletion", "not deletable" |
| "after modification", "when modified" | "{Subject}_modify" question | "cannot be modified", "immutable" |
| "when expired", "after expiration" | "{Subject} have expiration" question | "no expiration", "does not expire" |

**shouldSkip() Algorithm:**
```go
func shouldSkip(q *Ambiguity, decisions []Decision) bool {
    if len(q.DependsOn) == 0 {
        return false
    }
    for _, dep := range q.DependsOn {
        for _, d := range decisions {
            if d.ID == dep.QuestionID {
                // Check if answer matches any skip phrase (case-insensitive)
                for _, skipPhrase := range dep.SkipIfAnswer {
                    if strings.Contains(lower(d.Answer), lower(skipPhrase)) {
                        return true  // Skip this question
                    }
                }
            }
        }
    }
    return false
}
```

**Question Grouping (--grouped mode):**

```
┌─────────────────────────────────────────────────────────────┐
│ GroupQuestions(questions) → []QuestionGroup                  │
├─────────────────────────────────────────────────────────────┤
│ 1. Collect questions by Subject                              │
│ 2. Maintain original order (first occurrence of subject)     │
│ 3. Split into chunks of max 5 questions per group            │
│ 4. Assign group IDs: GRP-001, GRP-002, ...                   │
│ 5. Determine common category (or "mixed" if different)       │
└─────────────────────────────────────────────────────────────┘
```

**Constants:**
| Constant | Value | Purpose |
|----------|-------|---------|
| `MaxGroupSize` | 5 | Maximum questions per group |

**Related:** AC-INT-001 through AC-INT-006, BR-INT-001 through BR-INT-004

---

### TS-ARCH-002: Claude CLI Integration

**Pattern:** Anti-Corruption Layer (ACL) via Claude Code CLI

**File:** `internal/claude/client.go`, `internal/claude/retry.go`

#### Client Interface

```go
// internal/claude/client.go
type Client struct {
    SessionID string
    Verbose   bool
}

func NewClient() *Client
func (c *Client) Call(prompt string) (string, error)
func (c *Client) CallWithSystemPrompt(systemPrompt, userPrompt string) (string, error)
func (c *Client) CallJSON(prompt string, result interface{}) error
func (c *Client) CallJSONWithRetry(prompt string, result interface{}, cfg RetryConfig) error
```

#### CLI Invocation

**Command:** `claude -p "<prompt>" [--resume <session_id>]`

**Environment Variables:**
| Variable | Value | Purpose |
|----------|-------|---------|
| `CLAUDE_CODE_MAX_OUTPUT_TOKENS` | 100000 | Large output support for derivation |

**Implementation:**
```go
cmd := exec.Command("claude", "-p", prompt)
if c.SessionID != "" {
    cmd.Args = append(cmd.Args, "--resume", c.SessionID)
}
cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")
output, err := cmd.Output()
```

#### JSON Extraction Strategy

**File:** `internal/claude/client.go:extractJSON()`

**SEQ-JSON-001: Response Parsing**

LLM responses may contain JSON in various formats. The extraction follows this priority:

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Markdown Code Block Extraction                       │
├─────────────────────────────────────────────────────────────┤
│ Pattern: ```json\n{...}\n```                                 │
│ Regex: strings.Index(response, "```json")                   │
│ Extract content between ```json and next ```                 │
│ If found and valid JSON → RETURN SUCCESS                    │
└─────────────────────────────────────────────────────────────┘
                         │ FAIL
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 2: Raw JSON Extraction                                  │
├─────────────────────────────────────────────────────────────┤
│ Find first '{' or '[' character                             │
│ Find last matching '}' or ']' character                     │
│ Extract substring between markers                            │
│ If valid JSON → RETURN SUCCESS                              │
└─────────────────────────────────────────────────────────────┘
                         │ FAIL
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 3: JSON Sanitization                                    │
├─────────────────────────────────────────────────────────────┤
│ Fix common LLM JSON issues:                                  │
│ - Replace literal \n inside strings with \\n                │
│ - Remove control characters                                  │
│ - Handle escaped quotes                                      │
└─────────────────────────────────────────────────────────────┘
                         │ FAIL
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 4: Conversion Fallback (CallJSON only)                  │
├─────────────────────────────────────────────────────────────┤
│ Make second Claude call with prompt:                         │
│ "Convert this to valid JSON: {response}"                    │
│ Re-apply Steps 1-3 on new response                          │
└─────────────────────────────────────────────────────────────┘
```

**Sanitization Algorithm:**
```go
func sanitizeJSON(jsonStr string) string {
    var result strings.Builder
    inString := false
    escaped := false

    for i := 0; i < len(jsonStr); i++ {
        c := jsonStr[i]

        if escaped {
            result.WriteByte(c)
            escaped = false
            continue
        }

        if c == '\\' {
            result.WriteByte(c)
            escaped = true
            continue
        }

        if c == '"' {
            inString = !inString
            result.WriteByte(c)
            continue
        }

        // Replace literal newlines inside strings
        if inString && (c == '\n' || c == '\r') {
            result.WriteString("\\n")
            continue
        }

        result.WriteByte(c)
    }

    return result.String()
}
```

#### Retry Strategy

**File:** `internal/claude/retry.go`

**TS-RETRY-001: Exponential Backoff Configuration**

```go
type RetryConfig struct {
    MaxAttempts int           // Default: 3
    BaseDelay   time.Duration // Default: 2 seconds
    MaxDelay    time.Duration // Default: 30 seconds
}

func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   2 * time.Second,
        MaxDelay:    30 * time.Second,
    }
}
```

**Related DEC:** DEC-L1-015

**Backoff Formula:**
```
delay = min(BaseDelay * 2^(attempt-1), MaxDelay)

Attempt 1: 2s
Attempt 2: 4s
Attempt 3: 8s (capped at 30s if exceeded)
```

**TS-RETRY-002: Error Classification**

| Category | Patterns | Action |
|----------|----------|--------|
| Retryable | `rate limit`, `timeout`, `503`, `502`, `500`, `overloaded`, `connection refused`, `connection reset`, `temporary failure` | Retry with backoff |
| Non-retryable | `invalid api key`, `unauthorized`, `bad request`, `404` | Return error immediately |

**Implementation:**
```go
func isRetryableError(err error) bool {
    if err == nil {
        return false
    }
    errStr := strings.ToLower(err.Error())

    retryable := []string{
        "rate limit", "timeout", "connection refused",
        "connection reset", "temporary failure",
        "503", "502", "500", "overloaded",
    }

    for _, r := range retryable {
        if strings.Contains(errStr, r) {
            return true
        }
    }
    return false
}
```

#### Prompt Injection Pattern

**Context Injection:**

All prompts use XML-style context markers:
```markdown
<context>
</context>
```

**Injection Method:**
```go
func buildPrompt(template string, documents ...string) string {
    context := strings.Join(documents, "\n\n---\n\n")
    return strings.Replace(template, "</context>", context + "\n</context>", 1)
}
```

**Rationale:** Places context at end of prompt following Anthropic best practices for optimal attention.

**Related:** [prompt-catalog.md](prompt-catalog.md) for full prompt specifications

**Related BR:** BR-ENV-001, DEC-003, DEC-L1-001, DEC-L1-002

---

### TS-ARCH-003: Parallel Execution

**Pattern:** Bounded Parallelism with Semaphore

**Implementation:**
```go
type ParallelExecutor struct {
    maxConcurrent int  // Default: 3
    semaphore     chan struct{}
}

func (pe *ParallelExecutor) Execute(tasks []Task) []Result
```

**Constraints:**
- Maximum 3 concurrent API calls (BR-DRV-002)
- Error in one task does not stop others
- Results collected in order

**Related BR:** BR-DRV-002, NFR-PERF-002

---

### TS-ARCH-004: Interactive Mode

**Pattern:** Preview-Approve Loop

**Flow:**
```
1. Generate content
2. Render preview (first 20 lines, box border)
3. Prompt: [A]pprove / [E]dit / [R]egenerate / [S]kip / [Q]uit
4. Handle action
5. Repeat for next file
```

**Terminal Requirements:**
- Reads from stdin
- Writes prompts to stderr
- Supports $EDITOR for edit action

#### Preview Rendering

**Box Drawing Characters (Unicode):**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Preview: filename.md                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ # Document Title                                                             │
│                                                                              │
│ Content lines...                                                             │
│ ...                                                                          │
│ (truncated - showing 20 of 150 lines)                                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Rendering Constraints:**
| Constraint | Value | Purpose |
|------------|-------|---------|
| Max preview lines | 20 | Fit in terminal without scrolling |
| Max line width | 80 | Standard terminal width |
| Box width | 79 | Accommodate border characters |

**Related DEC:** DEC-L1-016

**Implementation:**
```go
const (
    maxPreviewLines = 20
    maxLineWidth    = 80
    boxWidth        = 79
)

// Box drawing characters
const (
    boxTopLeft     = "┌"
    boxTopRight    = "┐"
    boxBottomLeft  = "└"
    boxBottomRight = "┘"
    boxHorizontal  = "─"
    boxVertical    = "│"
    boxTeeRight    = "├"
    boxTeeLeft     = "┤"
)

func renderPreview(filename, content string) {
    lines := strings.Split(content, "\n")
    totalLines := len(lines)

    // Truncate if needed
    if len(lines) > maxPreviewLines {
        lines = lines[:maxPreviewLines]
    }

    // Draw box with content
    printBoxTop(filename)
    for _, line := range lines {
        printBoxLine(truncateLine(line, maxLineWidth-4))
    }
    if totalLines > maxPreviewLines {
        printBoxLine(fmt.Sprintf("(truncated - showing %d of %d lines)",
            maxPreviewLines, totalLines))
    }
    printBoxBottom()
}
```

#### User Input Handling

**Prompt Format:**
```
[A]pprove  [E]dit  [R]egenerate  [S]kip  [Q]uit  >
```

**Input Parsing Rules:**
| Input | Action | Notes |
|-------|--------|-------|
| `a`, `A`, `<Enter>` | Approve | Empty input = default approve |
| `e`, `E` | Edit | Opens in external editor |
| `r`, `R` | Regenerate | Re-calls Claude API |
| `s`, `S` | Skip | Continue to next file |
| `q`, `Q` | Quit | Requires confirmation |
| Other | Re-prompt | Invalid input message |

**Implementation:**
```go
func askApproval() string {
    fmt.Fprint(os.Stderr, "[A]pprove  [E]dit  [R]egenerate  [S]kip  [Q]uit  > ")

    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(strings.ToLower(input))

    // Empty input defaults to approve
    if input == "" {
        return "a"
    }
    return input
}
```

#### External Editor Integration

**Editor Selection Priority:**
```
1. $EDITOR environment variable
2. $VISUAL environment variable
3. vim (if available)
4. nano (if available)
5. vi (fallback)
```

**Related DEC:** DEC-L1-017

**Implementation:**
```go
func getEditor() string {
    if editor := os.Getenv("EDITOR"); editor != "" {
        return editor
    }
    if editor := os.Getenv("VISUAL"); editor != "" {
        return editor
    }

    // Fallback chain
    for _, editor := range []string{"vim", "nano", "vi"} {
        if _, err := exec.LookPath(editor); err == nil {
            return editor
        }
    }
    return "vi" // Ultimate fallback
}
```

**Edit Flow:**
```
┌─────────────────────────────────────────────────────────────────┐
│ 1. Write content to temp file (.md extension)                   │
│ 2. Launch editor: exec.Command(editor, tempFile)                │
│ 3. Wait for editor to exit                                      │
│ 4. Read modified content from temp file                         │
│ 5. Delete temp file                                             │
│ 6. Return to approval prompt with modified content              │
└─────────────────────────────────────────────────────────────────┘
```

**Temp File Handling:**
```go
func editContent(content string) (string, error) {
    // Create temp file with .md extension for syntax highlighting
    tmpFile, err := os.CreateTemp("", "loom-*.md")
    if err != nil {
        return "", err
    }
    tmpPath := tmpFile.Name()
    defer os.Remove(tmpPath)

    // Write current content
    tmpFile.WriteString(content)
    tmpFile.Close()

    // Launch editor
    editor := getEditor()
    cmd := exec.Command(editor, tmpPath)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("editor failed: %w", err)
    }

    // Read modified content
    modified, err := os.ReadFile(tmpPath)
    if err != nil {
        return "", err
    }

    return string(modified), nil
}
```

#### Quit Confirmation

**Flow:**
```
User: q
System: Are you sure you want to quit? Approved files will be saved. [y/N] >
User: y → Exit with code 0
User: n/other → Return to approval prompt
```

**Implementation:**
```go
func confirmQuit() bool {
    fmt.Fprint(os.Stderr,
        "Are you sure you want to quit? Approved files will be saved. [y/N] > ")

    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(strings.ToLower(input))

    return input == "y" || input == "yes"
}
```

#### Output Streams

| Stream | Content |
|--------|---------|
| stdout | Final approved content (for piping) |
| stderr | Preview boxes, prompts, status messages |

**Rationale:** Allows piping output while maintaining interactive UI.

**Related BR:** BR-DRV-004, DEC-031

---

## File Format Specifications

### TS-FMT-001: YAML Frontmatter

**Required Fields:**
```yaml
---
title: Document Title
generated: 2024-01-15T10:30:00Z  # ISO 8601
status: draft                     # draft | review | approved
level: L1                         # L0 | L1 | L2 | L3
loom-cli-version: 0.3.0          # Optional
---
```

**Related BR:** BR-DOC-001, DEC-015, DEC-029

---

### TS-FMT-002: Document ID Patterns

**Validation Regex (Go):**
```go
var idPatterns = map[string]*regexp.Regexp{
    // L1 patterns
    "AC":  regexp.MustCompile(`AC-[A-Z]+-\d{3}`),   // Acceptance Criteria
    "BR":  regexp.MustCompile(`BR-[A-Z]+-\d{3}`),   // Business Rules
    "ENT": regexp.MustCompile(`ENT-[A-Z]+`),        // Entities
    "BC":  regexp.MustCompile(`BC-[A-Z]+`),         // Bounded Contexts
    // L2 patterns
    "TC":   regexp.MustCompile(`TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2}`), // Test Cases
    "TS":   regexp.MustCompile(`TS-[A-Z]+-\d{3}`),  // Tech Specs
    "IC":   regexp.MustCompile(`IC-[A-Z]+-\d{3}`),  // Interface Contracts
    "AGG":  regexp.MustCompile(`AGG-[A-Z]+-\d{3}`), // Aggregates
    "SEQ":  regexp.MustCompile(`SEQ-[A-Z]+-\d{3}`), // Sequences
    // L3 patterns
    "EVT":  regexp.MustCompile(`EVT-[A-Z]+-\d{3}`),  // Domain Events
    "CMD":  regexp.MustCompile(`CMD-[A-Z]+-\d{3}`),  // Commands
    "INT":  regexp.MustCompile(`INT-[A-Z]+-\d{3}`),  // Integration Events
    "SVC":  regexp.MustCompile(`SVC-[A-Z]+`),        // Services
    "FDT":  regexp.MustCompile(`FDT-\d{3}`),         // Feature Definition Tickets
    "SKEL": regexp.MustCompile(`SKEL-[A-Z]+-\d{3}`), // Implementation Skeletons
    "DEP":  regexp.MustCompile(`DEP-[A-Z]+-\d{3}`),  // Dependency Graph entries
}
```

**ID Pattern Summary:**

| Level | Prefix | Pattern | Example | Document |
|-------|--------|---------|---------|----------|
| L1 | AC | `AC-{CTX}-{NNN}` | AC-ORD-001 | acceptance-criteria.md |
| L1 | BR | `BR-{CTX}-{NNN}` | BR-VAL-001 | business-rules.md |
| L1 | ENT | `ENT-{NAME}` | ENT-ORDER | domain-model.md |
| L1 | BC | `BC-{NAME}` | BC-ORDERING | bounded-context-map.md |
| L2 | TS | `TS-{CTX}-{NNN}` | TS-CLI-001 | tech-specs.md |
| L2 | IC | `IC-{CTX}-{NNN}` | IC-ANL-001 | interface-contracts.md |
| L2 | AGG | `AGG-{CTX}-{NNN}` | AGG-ORD-001 | aggregate-design.md |
| L2 | SEQ | `SEQ-{CTX}-{NNN}` | SEQ-ORD-001 | sequence-design.md |
| L3 | TC | `TC-AC-{CTX}-{NNN}-{T}{NN}` | TC-AC-ORD-001-P01 | test-cases.md |
| L3 | EVT | `EVT-{CTX}-{NNN}` | EVT-DOM-001 | event-message-design.md |
| L3 | CMD | `CMD-{CTX}-{NNN}` | CMD-ENT-001 | event-message-design.md |
| L3 | INT | `INT-{CTX}-{NNN}` | INT-EVT-001 | event-message-design.md |
| L3 | SVC | `SVC-{NAME}` | SVC-ORDR | service-boundaries.md |
| L3 | FDT | `FDT-{NNN}` | FDT-001 | feature-tickets.md |
| L3 | SKEL | `SKEL-{CTX}-{NNN}` | SKEL-ORDR-001 | implementation-skeletons.md |
| L3 | DEP | `DEP-{CTX}-{NNN}` | DEP-GRP-001 | dependency-graph.md |

**Legend:**
- `{CTX}` - Context/domain code (3-4 uppercase letters): ORD, CUST, VAL, CLI, etc.
- `{NNN}` - Sequential number (3 digits, zero-padded): 001, 002, ...
- `{NAME}` - Descriptive name (uppercase): ORDER, CUSTOMER, ORDERING
- `{T}` - Test type: P (positive), N (negative), B (boundary), H (hallucination)

**Related BR:** BR-DOC-002, DEC-016-022

---

### TS-FMT-003: JSON State Files

**Cascade State:** `.cascade-state.json`
```json
{
  "version": "1.0",
  "input_hash": "sha256...",
  "phases": {
    "analyze": { "status": "completed", "timestamp": "..." },
    "interview": { "status": "completed", "timestamp": "..." },
    "derive-l1": { "status": "running", "timestamp": "..." },
    "derive-l2": { "status": "pending" },
    "derive-l3": { "status": "pending" }
  },
  "config": {
    "skip_interview": false,
    "interactive": false
  },
  "timestamps": {
    "started": "2024-01-15T10:00:00Z",
    "completed": null
  }
}
```

**Interview State:** `.interview-state.json`
```json
{
  "questions": [...],
  "current_index": 2,
  "decisions": [
    { "question_id": "Q1", "answer": "...", "source": "user" }
  ]
}
```

**Related:** DEC-010, initial-data-model.md

---

## Error Handling Specifications

### TS-ERR-001: Exit Codes

| Code | Meaning | When |
|------|---------|------|
| 0 | Success | Command completed successfully |
| 1 | General error | Any unrecoverable error |
| 100 | Interview pending | Interview has more questions |

**Related BR:** BR-INT-003, DEC-008

---

### TS-ERR-002: Error Messages

**Format:** `Error: <message>` or `<command> failed: <cause>`

**Output:** stderr

#### Error Catalog by Command

**Global Errors:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `unknown command: {cmd}` | Invalid command name | Check `loom-cli help` for valid commands |

**analyze:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--input-file or --input-dir required` | Neither input specified | Provide one input option |
| `cannot specify both --input-file and --input-dir` | Both inputs specified | Use only one |
| `failed to read input file: {err}` | File not found or unreadable | Check file path and permissions |
| `domain discovery failed: {err}` | Claude API error | Check API key, retry |
| `entity analysis failed: {err}` | Claude API error | Check API key, retry |
| `operation analysis failed: {err}` | Claude API error | Check API key, retry |

**interview:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--state is required` | State file not specified | Provide `--state` path |
| `failed to read analysis file: {err}` | Analysis file not found | Run `analyze` first |
| `failed to parse analysis file: {err}` | Invalid JSON | Check analysis output |
| `failed to read state file: {err}` | State file not found | Use `--init` first |
| `failed to parse state file: {err}` | Corrupted state | Re-initialize interview |
| `failed to parse answer: {err}` | Invalid answer JSON | Check JSON format |
| `failed to parse answers: {err}` | Invalid batch answers | Check JSON array format |

**derive:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--output-dir is required` | Missing output directory | Specify output path |
| `failed to read analysis file: {err}` | Analysis not found | Run `analyze` first |
| `domain_model is required in input` | Missing domain model | Check analysis output |
| `domain model derivation failed: {err}` | Claude API error | Retry |
| `bounded context map derivation failed: {err}` | Claude API error | Retry |
| `derivation failed: {err}` | Claude API error | Retry |

**derive-l2:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--input-dir is required (directory containing L1 documents)` | Missing L1 input | Specify L1 directory |
| `--output-dir is required` | Missing output | Specify output path |
| `failed to read acceptance-criteria.md: {err}` | Missing L1 file | Run `derive` first |
| `failed to read business-rules.md: {err}` | Missing L1 file | Run `derive` first |
| `failed to read domain-model.md: {err}` | Missing L1 file | Run `derive` first |
| `no checkpoint found in {dir}, cannot resume` | No checkpoint for resume | Remove `--resume` or complete initial run |
| `failed to generate {phase}: {err}` | Phase generation failed | Check Claude API, retry |

**derive-l3:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--input-dir is required (directory containing L2 documents)` | Missing L2 input | Specify L2 directory |
| `--output-dir is required` | Missing output | Specify output path |
| `failed to read tech-specs.md: {err}` | Missing L2 file | Run `derive-l2` first |
| `failed to read acceptance-criteria.md: {err}` | Missing L1 file | Need AC for test cases |
| `failed to generate test cases: {err}` | Claude API error | Retry |
| `failed to generate API spec: {err}` | Claude API error | Retry |
| `failed to generate implementation skeletons: {err}` | Claude API error | Retry |

**cascade:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `either --input-file or --input-dir is required` | Missing L0 input | Provide input |
| `--output-dir is required` | Missing output | Specify output path |
| `failed to create directory {dir}: {err}` | Permission denied | Check directory permissions |
| `failed to load state: {err}` | Corrupted state file | Delete `.cascade-state.json` |
| `analyze failed: {err}` | Phase failed | Check Claude API |
| `derive L1 failed: {err}` | Phase failed | Check L0 input |
| `derive L2 failed: {err}` | Phase failed | Check L1 output |
| `derive L3 failed: {err}` | Phase failed | Check L2 output |

**validate:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--input-dir is required` | Missing input | Specify directory |
| `validation failed with {n} errors` | Validation errors found | Fix reported issues |

**sync-links:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `--input-dir is required` | Missing input | Specify directory |

**interactive mode:**
| Error Message | Cause | Resolution |
|--------------|-------|------------|
| `no editor found (set $EDITOR)` | No editor available | Set `$EDITOR` environment variable |
| `editor failed: {err}` | Editor error | Check editor installation |
| `failed to create temp file: {err}` | Temp directory issue | Check `/tmp` permissions |

---

## Performance Specifications

### TS-PERF-001: API Call Optimization

**Strategies:**
1. Batch related derivations in single prompt where possible
2. Use structured JSON output to reduce parsing errors
3. Retry transient failures (max 3 retries, exponential backoff)

**Targets:**
| Operation | Target Time | Max API Calls |
|-----------|-------------|---------------|
| analyze | < 30s | 3 |
| interview (per Q) | < 10s | 1 |
| derive (L1) | < 2 min | 5 |
| derive-l2 | < 3 min | 5 (parallel) |
| derive-l3 | < 5 min | 7 (parallel) |

**Related NFR:** NFR-PERF-001, NFR-PERF-002

---

### TS-PERF-002: File I/O Optimization

**Strategies:**
1. Read files once, cache in memory
2. Write files atomically (temp file + rename)
3. Validate before write (fail fast)

---

## Security Specifications

### TS-SEC-001: API Key Handling

**Requirements:**
- Read from environment variable only (never from file or command line)
- Never log or echo API key
- Clear from memory after use (best effort in Go)

**Related NFR:** NFR-SEC-001

---

### TS-SEC-002: Input Validation

**Requirements:**
- Validate file paths (no directory traversal)
- Validate JSON structure before processing
- Sanitize markdown output (no script injection)

---

## Testing Specifications

### TS-TEST-001: Test Categories

| Category | Coverage Target | Tools |
|----------|-----------------|-------|
| Unit tests | 70% | go test |
| Integration tests | Key flows | go test + fixtures |
| E2E tests | Happy paths | shell scripts |

**Related NFR:** NFR-TST-001

---

### TS-TEST-002: Test Fixtures

**Location:** `test/fixtures/`

**Structure:**
```
test/
├── fixtures/
│   ├── l0-input/
│   │   └── user-stories.md
│   ├── expected-l1/
│   │   ├── domain-model.md
│   │   └── ...
│   └── expected-l2/
│       └── ...
└── benchmark/
    └── 01-ecommerce-order/
```

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria |
| L2 | This document | Technical Specifications |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt Catalog |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI Interface Contract |
| L2 | [aggregate-design.md](aggregate-design.md) | Aggregate Design |
| L2 | [sequence-design.md](sequence-design.md) | Sequence Design |
| L2 | [initial-data-model.md](initial-data-model.md) | Data Model |
