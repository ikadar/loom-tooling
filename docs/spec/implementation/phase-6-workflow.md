# Phase 6: Workflow & Infrastructure Packages

## Cél

Implementáld a maradék infrastructure package-eket:
- `internal/workflow/` - Interactive workflow utilities
- `internal/config/` - Configuration management
- `internal/checkpoint/` - Resume functionality
- `internal/interview/` - Interview question grouping

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/package-structure.md | PKG-003 | checkpoint package |
| l2/package-structure.md | PKG-004 | config package |
| l2/package-structure.md | PKG-008 | interview package |
| l2/package-structure.md | PKG-009 | workflow package |
| l2/internal-api.md | internal/workflow | Progress, Approval |

---

## Implementálandó Fájlok

### 1. internal/workflow/progress.go

```
☐ Fájl: internal/workflow/progress.go
☐ Spec: l2/package-structure.md PKG-009, l2/internal-api.md
```

**Traceability:**
```go
// Package workflow provides interactive workflow utilities.
//
// Implements: l2/package-structure.md PKG-009
// See: l2/internal-api.md
package workflow
```

**Types and Functions:**

```go
// Progress tracks and displays progress.
type Progress struct {
    Label   string
    Total   int
    Current int
}

// NewProgress creates a progress tracker.
func NewProgress(label string, total int) *Progress

// Increment advances progress by 1.
func (p *Progress) Increment()

// Update sets current progress with message.
func (p *Progress) Update(current int, message string)

// Done marks progress complete.
func (p *Progress) Done()
```

**Display Format:**
```
[derive-l2] 3/5 Generating tech-specs...
```

---

### 2. internal/workflow/approval.go

```
☐ Fájl: internal/workflow/approval.go
☐ Spec: l2/package-structure.md PKG-009, l2/internal-api.md
```

**Functions:**

```go
// RequestApproval prompts user for yes/no.
//
// Implements: l2/internal-api.md
func RequestApproval(prompt string) (bool, error)

// ShowDiff displays a diff between before and after content.
func ShowDiff(before, after string) error

// ShowPreview displays truncated content preview.
//
// Implements: DEC-L1-016 (20 lines, 80 chars)
func ShowPreview(content string, maxLines, maxWidth int) string
```

---

### 3. internal/config/config.go

```
☐ Fájl: internal/config/config.go
☐ Spec: l2/package-structure.md PKG-004
```

**Traceability:**
```go
// Package config provides configuration management.
//
// Implements: l2/package-structure.md PKG-004
package config
```

**Types and Functions:**

```go
// Config holds CLI configuration.
type Config struct {
    InputFile      string
    InputDir       string
    OutputDir      string
    DecisionsFile  string
    AnalysisFile   string
    VocabularyFile string
    NFRFile        string
    Format         string  // "text" or "json"
    BatchMode      bool
    Verbose        bool
}

// ParseArgsForAnalyze parses analyze command arguments.
func ParseArgsForAnalyze(args []string) (*Config, error)

// ParseArgsForDerive parses derive command arguments.
func ParseArgsForDerive(args []string) (*Config, error)

// ReadInputFiles reads all L0 input markdown files.
// Returns: combined content, list of file paths, error
func (cfg *Config) ReadInputFiles() (string, []string, error)

// ReadVocabulary reads optional domain vocabulary file.
func (cfg *Config) ReadVocabulary() (string, error)

// ReadNFR reads optional non-functional requirements file.
func (cfg *Config) ReadNFR() (string, error)
```

---

### 4. internal/checkpoint/checkpoint.go

```
☐ Fájl: internal/checkpoint/checkpoint.go
☐ Spec: l2/package-structure.md PKG-003
```

**Traceability:**
```go
// Package checkpoint provides resume functionality for long operations.
//
// Implements: l2/package-structure.md PKG-003
package checkpoint
```

**Types and Functions:**

```go
// Checkpoint holds state for resumable operations.
type Checkpoint struct {
    Phase     string    `json:"phase"`
    Timestamp time.Time `json:"timestamp"`
    Data      any       `json:"data"`
}

// Save persists checkpoint to file.
func Save(path string, cp *Checkpoint) error

// Load reads checkpoint from file.
func Load(path string) (*Checkpoint, error)

// Exists checks if checkpoint file exists.
func Exists(path string) bool

// Clear removes checkpoint file.
func Clear(path string) error
```

---

### 5. internal/interview/grouping.go

```
☐ Fájl: internal/interview/grouping.go
☐ Spec: l2/package-structure.md PKG-008, DEC-L1-011
```

**Traceability:**
```go
// Package interview provides interview question grouping and processing.
//
// Implements: l2/package-structure.md PKG-008
// Implements: DEC-L1-011 (question grouping)
package interview
```

**Types and Functions:**

```go
// GroupQuestions groups ambiguities by subject/category.
//
// Implements: DEC-L1-011
// Max group size: 5 questions (domain.MaxGroupSize)
func GroupQuestions(ambiguities []domain.Ambiguity) []domain.QuestionGroup

// ShouldSkip checks if a question should be skipped based on previous answers.
//
// Implements: DEC-L1-008 (skip conditions)
func ShouldSkip(question domain.Ambiguity, decisions []domain.Decision) bool

// ProcessAnswer records an answer and returns updated state.
func ProcessAnswer(state *domain.InterviewState, answer domain.AnswerInput) error
```

**Grouping Strategy:**
1. Group by Subject (entity/operation name)
2. Within subject, group by Category
3. Max 5 questions per group (MaxGroupSize)
4. If more, split into multiple groups

---

## Definition of Done

```
☐ internal/workflow/progress.go - Progress struct, NewProgress, Increment, Update, Done
☐ internal/workflow/approval.go - RequestApproval, ShowDiff, ShowPreview
☐ internal/config/config.go - Config struct, ParseArgs*, ReadInputFiles, ReadVocabulary, ReadNFR
☐ internal/checkpoint/checkpoint.go - Checkpoint struct, Save, Load, Exists, Clear
☐ internal/interview/grouping.go - GroupQuestions (max 5), ShouldSkip, ProcessAnswer
☐ Minden fájl tartalmaz traceability kommentet
☐ `go build` HIBA NÉLKÜL fut
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
