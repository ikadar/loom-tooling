# Phase 3: CLI Commands

## Cél

Implementáld az összes CLI parancsot a `cmd/` package-ben.

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/interface-contracts.md | IC-ANL-001 | analyze command |
| l2/interface-contracts.md | IC-INT-001 | interview command |
| l2/interface-contracts.md | IC-DRV-001 | derive command |
| l2/interface-contracts.md | IC-DRV-002 | derive-l2 command |
| l2/interface-contracts.md | IC-DRV-003 | derive-l3 command |
| l2/interface-contracts.md | IC-VAL-001 | validate command |
| l2/interface-contracts.md | IC-SYN-001 | sync-links command |
| l2/interface-contracts.md | IC-CAS-001 | cascade command |
| l2/sequence-design.md | SEQ-* | Execution flows |
| l2/tech-specs.md | TS-ARCH-001 | Command router |

---

## KRITIKUS: Prompts Package Használata

**MINDEN parancsnak a `prompts` package-ből kell a promptokat használnia!**

Lásd: `l2/tech-specs.md`, `l2/build-test-guide.md`

### Import

```go
import (
    "loom-cli/internal/claude"
    "loom-cli/prompts"
)
```

### Prompt Használat

**HELYTELEN (inline prompt):**
```go
// ❌ NE ÍGY!
prompt := "You are a DDD expert. Generate a domain-model.md..."
response, _ := client.Call(prompt)
```

**HELYES (prompts package):**
```go
// ✅ ÍGY KELL!
prompt := claude.BuildPrompt(prompts.DeriveDomainModel, inputContent)
var result SomeStruct
err := client.CallJSON(prompt, &result)
```

### Prompt → Parancs Mapping

| Parancs | Prompts |
|---------|---------|
| analyze | `prompts.DomainDiscovery`, `prompts.EntityAnalysis`, `prompts.OperationAnalysis` |
| interview | `prompts.Interview` |
| derive | `prompts.DeriveDomainModel`, `prompts.DeriveBoundedContext`, `prompts.Derivation` |
| derive-l2 | `prompts.DeriveTechSpecs`, `prompts.DeriveInterfaceContracts`, `prompts.DeriveAggregateDesign`, `prompts.DeriveSequenceDesign`, `prompts.DeriveDataModel` |
| derive-l3 | `prompts.DeriveTestCases`, `prompts.DeriveL3API`, `prompts.DeriveL3Skeletons`, `prompts.DeriveFeatureTickets`, `prompts.DeriveServiceBoundaries`, `prompts.DeriveEventDesign`, `prompts.DeriveDependencyGraph` |

### Output Pattern

A promptok **JSON outputot** várnak, amit a formatter package markdown-ra alakít:

```go
// 1. Claude JSON-t ad vissza
var techSpecs []formatter.TechSpec
err := client.CallJSON(prompt, &techSpecs)

// 2. Formatter markdown-ra alakítja
markdown := formatter.FormatTechSpecs(techSpecs)

// 3. Fájlba írás
os.WriteFile("tech-specs.md", []byte(markdown), 0644)
```

---

## Implementálandó Fájlok

### 1. cmd/analyze.go

```
☐ Fájl: cmd/analyze.go
☐ Spec: IC-ANL-001, TS-ARCH-001a, SEQ-CAS-001
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-ANL-001
// See: l2/sequence-design.md SEQ-ANL-001
```

**Options:**
- `--input-file <path>` - Single L0 input file
- `--input-dir <path>` - Directory with L0 files
- `--output <path>` - Output file (default: stdout)
- `--decisions <path>` - Existing decisions.md

**Flow:**
1. Read input file(s)
2. Call Claude with domain-discovery prompt
3. Call Claude with entity-analysis prompt
4. Call Claude with operation-analysis prompt
5. Merge results into AnalyzeResult
6. Output JSON

---

### 2. cmd/interview.go

```
☐ Fájl: cmd/interview.go
☐ Spec: IC-INT-001, SEQ-INT-001, SEQ-INT-002
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-INT-001
// See: l2/sequence-design.md SEQ-INT-001, SEQ-INT-002
```

**Options:**
- `--init <path>` - Initialize from analysis JSON
- `--state <path>` - State file path
- `--answer <json>` - Answer JSON: `{"question_id":"...","answer":"...","source":"user"}`
- `--skip` - Skip current question with default

**Exit Codes:**
- 0: Interview complete
- 1: Error
- 100: More questions available (outputs question JSON)

**Flow (--init):**
1. Load analysis JSON
2. Create InterviewState with questions
3. Output first question, exit 100

**Flow (--answer):**
1. Load state
2. Record answer as Decision
3. Check skip conditions (DEC-L1-008)
4. If more questions: output next, exit 100
5. If complete: output final state, exit 0

---

### 3. cmd/derive_new.go

```
☐ Fájl: cmd/derive_new.go
☐ Spec: IC-DRV-001, TS-ARCH-001b
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-DRV-001
// See: l2/sequence-design.md SEQ-DRV-001
```

**Options:**
- `--output-dir <path>` - Output directory (required)
- `--analysis-file <path>` - Analysis JSON or interview state
- `--decisions <path>` - Existing decisions.md
- `--vocabulary <path>` - Optional vocabulary file
- `--nfr <path>` - Optional NFR file

**Output Files:**
- `domain-model.md`
- `bounded-context-map.md`
- `acceptance-criteria.md`
- `business-rules.md`
- `decisions.md`

---

### 4. cmd/derive_l2.go

```
☐ Fájl: cmd/derive_l2.go
☐ Spec: IC-DRV-002, SEQ-DRV-002
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-DRV-002
// See: l2/sequence-design.md SEQ-L2-001
```

**Options:**
- `--input-dir <path>` - L1 docs directory
- `--output-dir <path>` - Output directory (required)
- `--interactive, -i` - Interactive approval mode

**Output Files:**
- `tech-specs.md`
- `interface-contracts.md`
- `aggregate-design.md`
- `sequence-design.md`
- `initial-data-model.md`

---

### 5. cmd/derive_l2_convert.go

```
☐ Fájl: cmd/derive_l2_convert.go
☐ Spec: l2/package-structure.md
```

**Purpose:** JSON to Markdown conversion utilities for L2 documents.

**Functions:**
- `convertTechSpecsToMarkdown()`
- `convertContractsToMarkdown()`
- `convertAggregateToMarkdown()`
- `convertSequenceToMarkdown()`
- `convertDataModelToMarkdown()`

---

### 6. cmd/derive_l3.go

```
☐ Fájl: cmd/derive_l3.go
☐ Spec: IC-DRV-003
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-DRV-003
// See: l2/sequence-design.md SEQ-L3-001
```

**Options:**
- `--input-dir <path>` - L2 docs directory
- `--l1-dir <path>` - L1 docs directory (for AC refs)
- `--output-dir <path>` - Output directory (required)

**Output Files:**
- `test-cases.md`
- `openapi.json`
- `implementation-skeletons.md`
- `feature-tickets.md`
- `service-boundaries.md`
- `event-message-design.md`
- `dependency-graph.md`

---

### 7. cmd/validate.go

```
☐ Fájl: cmd/validate.go
☐ Spec: IC-VAL-001, SEQ-VAL-001, BR-VAL-001/002/003
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-VAL-001
// See: l2/sequence-design.md SEQ-VAL-001
// See: l1/business-rules.md BR-VAL-001/002/003
```

**Options:**
- `--input-dir <path>` - Directory to validate
- `--level <L1|L2|L3|ALL>` - Validation level (default: ALL)
- `--json` - Output as JSON

**Validation Rules:**
| Rule | Level | Description |
|------|-------|-------------|
| V001 | ALL | Documents have IDs |
| V002 | ALL | IDs follow patterns |
| V003 | ALL | References exist |
| V004 | ALL | Bidirectional links (deferred - DEC-L1-013) |
| V005 | L2+ | AC has test cases |
| V006 | L2+ | Entity has aggregate |
| V007 | L2+ | Service has contract |
| V008 | L3 | Negative test ratio >= 20% |
| V009 | L3 | Hallucination tests exist |
| V010 | ALL | No duplicate IDs |

---

### 8. cmd/sync_links.go

```
☐ Fájl: cmd/sync_links.go
☐ Spec: IC-SYN-001, SEQ-SYN-001
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-SYN-001
// See: l2/sequence-design.md SEQ-SYN-001
```

**Options:**
- `--input-dir <path>` - Directory to sync
- `--dry-run` - Show changes without modifying

---

### 9. cmd/cascade.go

```
☐ Fájl: cmd/cascade.go
☐ Spec: IC-CAS-001, SEQ-CAS-001, SEQ-CAS-002, AGG-CAS-001
```

**Traceability:**
```go
// Implements: l2/interface-contracts.md IC-CAS-001
// See: l2/sequence-design.md SEQ-CAS-001, SEQ-CAS-002
// See: l2/aggregate-design.md AGG-CAS-001
```

**Options:**
- `--input-file <path>` - L0 input file
- `--input-dir <path>` - L0 input directory
- `--output-dir <path>` - Output directory (required)
- `--skip-interview` - Use AI defaults
- `--decisions <path>` - Existing decisions
- `--interactive, -i` - Interactive approval
- `--resume` - Resume from saved state
- `--from <level>` - Resume from level (l1, l2, l3)

**State File:** `.cascade-state.json`

**Phases:**
1. Analyze → `.analysis.json`
2. Interview (if not --skip-interview) → `.interview-state.json`
3. Derive L1 → `l1/`
4. Derive L2 → `l2/`
5. Derive L3 → `l3/`

---

### 10. cmd/interactive.go

```
☐ Fájl: cmd/interactive.go
☐ Spec: l2/tech-specs.md TS-ARCH-004, DEC-L1-016, DEC-L1-017
```

**Traceability:**
```go
// Implements: l2/tech-specs.md TS-ARCH-004
// Implements: DEC-L1-016 (preview limits: 20 lines, 80 chars)
// Implements: DEC-L1-017 (editor priority: EDITOR→VISUAL→vim→nano→vi)
```

**Functions:**
- `showPreview(content string)` - Truncated preview (20 lines, 80 chars)
- `requestApproval(prompt string) (ApprovalAction, error)`
- `openInEditor(content string) (string, error)` - Edit in external editor
- `findEditor() string` - Find editor per DEC-L1-017

---

## Definition of Done

```
☐ MINDEN cmd/*.go IMPORTÁLJA: "loom-cli/prompts"
☐ MINDEN Claude hívás prompts.X változót használ (NEM inline string!)
☐ MINDEN derivation CallJSON()-t használ (NEM Call()-t!)
☐ cmd/analyze.go - reads input, calls Claude, outputs JSON
☐ cmd/interview.go - init, answer, skip flow works
☐ cmd/interview.go - exit code 100 when more questions
☐ cmd/derive_new.go - generates 5 L1 files
☐ cmd/derive_l2.go - generates 5 L2 files
☐ cmd/derive_l2_convert.go - JSON→Markdown utilities
☐ cmd/derive_l3.go - generates 7 L3 files
☐ cmd/validate.go - all 10 rules implemented
☐ cmd/sync_links.go - finds and fixes missing links
☐ cmd/cascade.go - orchestrates all phases
☐ cmd/cascade.go - state persistence works
☐ cmd/cascade.go - --resume works
☐ cmd/interactive.go - preview, approval, editor functions
☐ Minden fájl tartalmaz traceability kommentet
☐ `go build` HIBA NÉLKÜL fut
```

**Ellenőrzés:**
```bash
# Minden cmd file importálja a prompts package-et?
grep -L "loom-cli/prompts" cmd/*.go  # ÜRES kell legyen!

# Nincs inline prompt?
grep -r "You are a" cmd/*.go  # ÜRES kell legyen!
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
