# Phase 1: Core Infrastructure

## Cél

Implementáld a loom-cli core infrastruktúráját: entry point, command router, domain types.

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/package-structure.md | PKG-001, PKG-002, PKG-010 | main.go, cmd/root.go, domain/types.go |
| l2/internal-api.md | internal/domain | Domain type definitions |
| l2/aggregate-design.md | AGG-* | Aggregate structs |
| l2/interface-contracts.md | Exit codes | Exit code constants |
| l1/decisions.md | DEC-L1-014 | Go 1.21+ requirement |

---

## Implementálandó Fájlok

### 1. go.mod

```
☐ Fájl: go.mod
☐ Spec: l2/package-structure.md, DEC-L1-014
```

**Követelmények:**
- Module name: `loom-cli`
- Go version: 1.21 (DEC-L1-014)
- No external dependencies (DEC-L1-001)

**Traceability komment:**
```go
// Implements: DEC-L1-014 (Go 1.21+)
// Implements: DEC-L1-001 (stdlib only, no external dependencies)
```

---

### 2. main.go

```
☐ Fájl: main.go
☐ Spec: l2/package-structure.md PKG-001
```

**Követelmények:**
- Entry point
- Calls cmd.Execute()
- Returns exit code

**Minta:**
```go
// Package main is the entry point for loom-cli.
//
// Implements: l2/package-structure.md PKG-001
// See: l2/interface-contracts.md, l2/tech-specs.md TS-ARCH-001
package main

import (
    "os"
    "loom-cli/cmd"
)

func main() {
    os.Exit(cmd.Execute())
}
```

---

### 3. cmd/root.go

```
☐ Fájl: cmd/root.go
☐ Spec: l2/package-structure.md PKG-002, l2/interface-contracts.md
```

**Követelmények:**
- Execute() returns int (exit code)
- Routes to subcommands: analyze, interview, derive, derive-l2, derive-l3, validate, sync-links, cascade
- Uses flag.FlagSet for argument parsing
- Implements help and version commands

**Traceability komment:**
```go
// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md (CLI interface)
// Implements: l2/tech-specs.md TS-ARCH-001 (Command router architecture)
// See: l0/decisions.md DEC-L1-001 (stdlib only, no cobra)
package cmd
```

**Exit Codes (l2/interface-contracts.md):**
- 0: Success
- 1: Error
- 100: Interview has more questions

---

### 4. internal/domain/types.go

```
☐ Fájl: internal/domain/types.go
☐ Spec: l2/internal-api.md, l2/aggregate-design.md
```

**Implementálandó típusok:**

| Type | Spec Reference | Leírás |
|------|----------------|--------|
| Entity | AGG-ANL-001 | L0 analysis entity |
| Operation | AGG-ANL-001 | L0 analysis operation |
| Relationship | AGG-ANL-001 | Entity relationships |
| Aggregate | AGG-ANL-001 | DDD aggregate |
| Domain | AGG-ANL-001 | Complete domain model |
| SkipCondition | AGG-INT-001, DEC-L1-008 | Question skip logic |
| Ambiguity | AGG-ANL-001 | Unresolved question |
| AnalyzeResult | AGG-ANL-001 | Analysis aggregate root |
| Decision | AGG-INT-001, DEC-L1-004 | Resolved ambiguity |
| InterviewState | AGG-INT-001 | Interview aggregate root |
| QuestionGroup | AGG-INT-001, DEC-L1-011 | Grouped questions |
| InterviewOutput | AGG-INT-001 | CLI output for question |
| AnswerInput | TBL-OUT-003 | Answer JSON input |
| **AcceptanceCriteria** | l2/internal-api.md | L1 derivation output |
| **BusinessRule** | l2/internal-api.md | L1 derivation output |
| **DerivationResult** | l2/internal-api.md | L1 derivation aggregate |
| **DerivationStats** | l2/internal-api.md | Derivation statistics |
| PhaseState | AGG-CAS-001 | Cascade phase state |
| CascadeStateConfig | AGG-CAS-001 | Cascade config |
| CascadeState | AGG-CAS-001 | Cascade aggregate root |
| CascadeConfig | VO-001 | Cascade command config |
| ValidationError | AGG-VAL-001 | Validation error |
| ValidationWarning | AGG-VAL-001 | Validation warning |
| ValidationCheck | AGG-VAL-001 | Single check result |
| ValidationSummary | AGG-VAL-001 | Validation stats |
| ValidationResult | AGG-VAL-001 | Validation aggregate root |
| PhaseResult | VO-002, DEC-031 | Interactive phase result |
| ApprovalAction | VO-003, DEC-031 | Interactive approval action |

**Constants (l2/tech-specs.md):**
```go
// Exit codes as specified in l2/tech-specs.md TS-ERR-001.
const (
    ExitCodeSuccess  = 0
    ExitCodeError    = 1
    ExitCodeQuestion = 100
)

// Implements: DEC-L1-011, DEC-L1-016, l2/tech-specs.md
const (
    MaxGroupSize    = 5  // Maximum questions per group
    MaxPreviewLines = 20 // Preview rendering limit (DEC-L1-016)
    MaxLineWidth    = 80 // Terminal width (DEC-L1-016)
)
```

---

## Pattern Követelmények

**FACTORY PATTERN (DP-FAC-001):**

Minden aggregate/entity struct-hoz **kötelező** `New{Type}()` constructor:

```go
// ✅ HELYES - Factory pattern
func NewAnalyzeResult(domain Domain, entities []Entity) *AnalyzeResult {
    return &AnalyzeResult{
        Domain:   domain,
        Entities: entities,
        // ... validation, defaults
    }
}

// ❌ HIBÁS - Közvetlen struct literal
result := &AnalyzeResult{Domain: d, Entities: e}
```

**Ellenőrzés:** Nincs közvetlen struct literal exportált aggregate-re a kódban.

---

## Definition of Done

```
☐ go.mod létezik, Go 1.21, no dependencies
☐ main.go létezik, traceability kommenttel
☐ cmd/root.go létezik, Execute() returns int
☐ cmd/root.go routes all 8 commands (analyze, interview, derive, derive-l2, derive-l3, validate, sync-links, cascade)
☐ internal/domain/types.go tartalmazza MIND a 28 típust
☐ Exit code constants definiálva
☐ MaxGroupSize, MaxPreviewLines, MaxLineWidth constants definiálva
☐ `go build` HIBA NÉLKÜL fut (stub functions OK)
☐ Minden fájl tartalmaz traceability kommentet
☐ **PATTERN:** Minden aggregate-hez van New{Type}() constructor
☐ **PATTERN:** Nincs közvetlen struct literal exportált aggregate-re
```

---

## Stub Pattern

A többi phase-ig használj stub függvényeket:

```go
func runAnalyze(args []string) int {
    // TODO: Implement in Phase 3
    fmt.Fprintln(os.Stderr, "analyze: not yet implemented")
    return domain.ExitCodeError
}
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
