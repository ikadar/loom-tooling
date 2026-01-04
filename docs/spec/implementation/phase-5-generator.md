# Phase 5: Generator Package

## Cél

Implementáld az `internal/generator/` package-et: document generation orchestration.

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/package-structure.md | PKG-006 | Package structure |
| l2/internal-api.md | internal/generator | ChunkedTestCaseGenerator, TestCaseResult |
| l2/tech-specs.md | Chunking | AC batching strategy |

---

## Implementálandó Fájlok

### 1. internal/generator/testcases.go

```
☐ Fájl: internal/generator/testcases.go
☐ Spec: l2/package-structure.md PKG-006, l2/internal-api.md
```

**Traceability:**
```go
// Package generator provides document generation orchestration.
//
// Implements: l2/package-structure.md PKG-006
// See: l2/internal-api.md
package generator
```

**Típusok és függvények:**

```go
// ChunkedTestCaseGenerator generates test cases in batches.
//
// Implements: l2/internal-api.md
type ChunkedTestCaseGenerator struct {
    Client    *claude.Client
    ChunkSize int  // Default: 5 ACs per chunk
}

// NewChunkedTestCaseGenerator creates generator with default settings.
func NewChunkedTestCaseGenerator(client *claude.Client) *ChunkedTestCaseGenerator

// Generate generates test cases from AC markdown content.
//
// Strategy:
// 1. Parse ACs from markdown
// 2. Chunk into groups of ChunkSize
// 3. Generate test cases per chunk (parallel if possible)
// 4. Merge results
func (g *ChunkedTestCaseGenerator) Generate(acContent string) (*TestCaseResult, error)

// TestCaseResult is the result of test case generation.
type TestCaseResult struct {
    TestSuites []TestSuite `json:"test_suites"`
    Summary    TDAISummary `json:"summary"`
}

type TestSuite struct {
    ACRef   string     `json:"ac_ref"`
    ACTitle string     `json:"ac_title"`
    Tests   []TestCase `json:"tests"`
}

// FlattenTestCases extracts all test cases from suites.
func FlattenTestCases(suites []TestSuite) []TestCase

// ChunkACs splits ACs into groups.
func ChunkACs(acs []AC, chunkSize int) [][]AC
```

**Chunk Strategy:**
- Default chunk size: 5 ACs per batch (l2/internal-api.md)
- Prevents context overflow
- Enables parallel processing

---

### 2. internal/generator/parallel.go

```
☐ Fájl: internal/generator/parallel.go
☐ Spec: l2/package-structure.md PKG-006
```

**Traceability:**
```go
// Implements: l2/package-structure.md PKG-006
package generator
```

**Functions:**

```go
// ProcessInParallel processes items in parallel with worker limit.
//
// Generic parallel processor for any derivation task.
func ProcessInParallel[T, R any](items []T, fn func(T) (R, error), workers int) ([]R, error)

// Example usage:
// results, err := ProcessInParallel(chunks, func(chunk []AC) ([]TestCase, error) {
//     return generateTestCasesForChunk(chunk)
// }, 3)
```

**Implementation Notes:**
- Use `sync.WaitGroup` for coordination
- Use buffered channel for results
- Collect errors, don't fail on first error
- Default workers: 3 (balance between speed and rate limits)

---

## Definition of Done

```
☐ internal/generator/testcases.go létezik
☐ ChunkedTestCaseGenerator struct implemented
☐ NewChunkedTestCaseGenerator() returns generator with ChunkSize=5
☐ Generate() parses ACs, chunks, generates, merges
☐ TestCaseResult, TestSuite types defined
☐ FlattenTestCases() implemented
☐ ChunkACs() implemented
☐ internal/generator/parallel.go létezik
☐ ProcessInParallel() generic function implemented
☐ Minden fájl tartalmaz traceability kommentet
☐ `go build` HIBA NÉLKÜL fut
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
