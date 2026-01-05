# Loom CLI Implementation Master Checklist

## Overview

Ez a master checklist követi az összes implementációs fázis állapotát.

**Használat:**
1. Fázisonként haladj (ne ugorj előre!)
2. Minden fázis végén jelöld meg a checklistet
3. Csak akkor lépj tovább, ha MINDEN elem kész
4. Ha egy session megszakad, a következő innen folytathatja

---

## Phase Status

| Phase | Leírás | Státusz |
|-------|--------|---------|
| Phase 1 | Core Infrastructure | ☐ Not Started |
| Phase 2 | Claude Client | ☐ Not Started |
| Phase 3 | CLI Commands | ☐ Not Started |
| Phase 4 | Formatter Package | ☐ Not Started |
| Phase 5 | Generator Package | ☐ Not Started |
| Phase 6 | Workflow & Infrastructure | ☐ Not Started |
| Phase 7 | Prompt Content | ☐ Not Started |

---

## Phase 1: Core Infrastructure

**Prompt:** `implementation/phase-1-core.md`

**Files:**
```
☐ go.mod (Go 1.21, no deps)
☐ main.go
☐ cmd/root.go
☐ internal/domain/types.go
```

**Verification:**
```
☐ go build succeeds
☐ All 28 domain types defined
☐ Exit codes defined
☐ All files have traceability comments
```

**Pattern Check:**
```
☐ FACTORY: Minden domain type-hoz van New{Type}() constructor (ha van invariáns)
☐ FACTORY: Nincs közvetlen struct literal exportált aggregate-re
```

**Phase Complete:** ☐

---

## Phase 2: Claude Client

**Prompt:** `implementation/phase-2-claude.md`

**Files:**
```
☐ internal/claude/client.go
☐ internal/claude/retry.go
```

**Verification:**
```
☐ Client.Call() uses exec.Command("claude", ...)
☐ CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000 set
☐ extractJSON() handles markdown blocks
☐ DefaultRetryConfig() returns 3/2s/30s
☐ isRetryableError() classifies correctly
☐ go build succeeds
```

**Pattern Check:**
```
☐ ADAPTER (ACL): Claude CLI = external service, internal/claude/ = adapter layer
☐ ADAPTER: Domain nem függ közvetlenül a claude CLI-től
☐ FACTORY: NewClient() constructor használata
```

**Phase Complete:** ☐

---

## Phase 3: CLI Commands

**Prompt:** `implementation/phase-3-commands.md`

**Files:**
```
☐ cmd/analyze.go
☐ cmd/interview.go
☐ cmd/derive_new.go
☐ cmd/derive_l2.go
☐ cmd/derive_l2_convert.go
☐ cmd/derive_l3.go
☐ cmd/validate.go
☐ cmd/sync_links.go
☐ cmd/cascade.go
☐ cmd/interactive.go
```

**Verification:**
```
☐ MINDEN cmd/*.go IMPORTÁLJA: "loom-cli/prompts"
☐ MINDEN Claude hívás prompts.X változót használ (NEM inline string!)
☐ analyze: reads input, calls Claude, outputs JSON
☐ interview: init/answer flow, exit code 100
☐ derive: generates 5 L1 files
☐ derive-l2: generates 5 L2 files
☐ derive-l3: generates 7 L3 files
☐ validate: all 10 rules (V001-V010)
☐ cascade: orchestrates phases, state persistence
☐ interactive: preview (20 lines), editor selection
☐ go build succeeds
```

**Pattern Check:**
```
☐ FACADE: Cascade command = Facade a teljes derivációs workflow-hoz
☐ CHAIN OF RESP: Cascade phases egymás után futnak, bármely megállíthatja
☐ FACTORY: Minden result struct New*() constructorral jön létre
```

**Phase Complete:** ☐

---

## Phase 4: Formatter Package

**Prompt:** `implementation/phase-4-formatter.md`

**Files (MUST BE 9!):**
```
☐ internal/formatter/types.go
☐ internal/formatter/frontmatter.go
☐ internal/formatter/anchor.go
☐ internal/formatter/techspecs.go
☐ internal/formatter/testcases.go
☐ internal/formatter/contracts.go
☐ internal/formatter/aggregates.go
☐ internal/formatter/sequences.go
☐ internal/formatter/datamodel.go
```

**Verification:**
```
☐ All types from l2/internal-api.md defined
☐ FormatFrontmatter() generates YAML
☐ FormatSequenceDesign() generates Mermaid
☐ go build succeeds
☐ File count = 9 (not 8!)
```

**Phase Complete:** ☐

---

## Phase 5: Generator Package

**Prompt:** `implementation/phase-5-generator.md`

**Files:**
```
☐ internal/generator/testcases.go
☐ internal/generator/parallel.go
```

**Verification:**
```
☐ ChunkedTestCaseGenerator with ChunkSize=5
☐ ProcessInParallel generic function works
☐ go build succeeds
```

**Phase Complete:** ☐

---

## Phase 6: Workflow & Infrastructure

**Prompt:** `implementation/phase-6-workflow.md`

**Files:**
```
☐ internal/workflow/progress.go
☐ internal/workflow/approval.go
☐ internal/config/config.go
☐ internal/checkpoint/checkpoint.go
☐ internal/interview/grouping.go
```

**Verification:**
```
☐ Progress tracking works
☐ Approval prompts work
☐ Config parsing works
☐ Checkpoint save/load works
☐ Question grouping (max 5) works
☐ go build succeeds
```

**Phase Complete:** ☐

---

## Phase 7: Prompt Content

**Prompt:** `implementation/phase-7-prompts.md`

**Files (MUST BE 21!):**
```
Analyze (3):
☐ prompts/domain-discovery.md
☐ prompts/entity-analysis.md
☐ prompts/operation-analysis.md

Interview (1):
☐ prompts/interview.md

L1 Derivation (3):
☐ prompts/derivation.md
☐ prompts/derive-domain-model.md
☐ prompts/derive-bounded-context.md

L2 Derivation (6):
☐ prompts/derive-l2.md
☐ prompts/derive-tech-specs.md
☐ prompts/derive-interface-contracts.md
☐ prompts/derive-aggregate-design.md
☐ prompts/derive-sequence-design.md
☐ prompts/derive-data-model.md

L3 Derivation (8):
☐ prompts/derive-test-cases.md
☐ prompts/derive-l3.md
☐ prompts/derive-l3-api.md
☐ prompts/derive-l3-skeletons.md
☐ prompts/derive-feature-tickets.md
☐ prompts/derive-service-boundaries.md
☐ prompts/derive-event-design.md
☐ prompts/derive-dependency-graph.md
```

**Verification:**
```
☐ Every prompt has <role> section
☐ Every prompt has <task> section
☐ Every prompt has <thinking_process> section
☐ Every prompt has <instructions> section
☐ Every prompt has <output_format> section
☐ Every prompt has <examples> section (2+ examples)
☐ Every prompt has <self_review> section
☐ Every prompt has <critical_output_format> section
☐ Every prompt has <context></context>
☐ prompts.go loads all 21 prompts
☐ go build succeeds
☐ Prompt count = 21
```

**Phase Complete:** ☐

---

## Final Verification

```
☐ All 7 phases complete
☐ go build -o loom-cli . succeeds
☐ ./loom-cli --help shows all commands
☐ ./loom-cli cascade --help shows options
```

---

## Test Run

```bash
# Build
cd loom-cli && go build -o loom-cli .

# Full cascade test
./loom-cli cascade \
  --input-file ../test/benchmark/01-ecommerce-order/input-l0.md \
  --output-dir /tmp/test-output \
  --skip-interview

# Validation
./loom-cli validate --input-dir /tmp/test-output --level ALL

# Expected: Exit 0, no validation errors
```

**Test Passed:** ☐

---

## Implementation Complete

**All phases done:** ☐
**All tests passed:** ☐
**Ready for review:** ☐
