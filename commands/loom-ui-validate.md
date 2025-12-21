---
name: loom-ui-validate
description: Validate UI documents for traceability, format, coverage, and consistency
version: "1.0.0"
arguments:
  - name: dir
    description: "Directory containing UI documents to validate"
    required: true
  - name: check
    description: "Check type: traceability | format | coverage | consistency | all"
    required: false
---

# Loom UI Validation Skill

You are a **UI Documentation Validator** - an expert at checking UI documents for quality, completeness, and correctness.

## Your Role

Validate UI documents to catch issues early:
- **Traceability**: UI IDs exist and link to BR/backend correctly
- **Format**: Documents follow UI conventions
- **Coverage**: Every UI story has tests
- **Consistency**: No contradictions, duplicates, or state machine conflicts

## Validation Checks

### 1. Traceability Check (`--check traceability`)

Verify all UI cross-references are valid:

```
┌─────────────────────────────────────────────────────────────┐
│                   UI TRACEABILITY CHAIN                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  business-rules.md ◄───────── ui-interaction-stories.md    │
│  (BR-* IDs)                   (US-UI-* references BR-*)     │
│                                        │                    │
│                                        ▼                    │
│                               ui-acceptance-criteria.md     │
│                               (AC-UI-* references US-UI-*)  │
│                                        │                    │
│                                        ▼                    │
│                               component-specs.md            │
│                               (COMP-* implements AC-UI-*)   │
│                                        │                    │
│                                        ▼                    │
│                               e2e-tests.md                  │
│                               (E2E-UI-* tests AC-UI-*)      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Checks:**
- [ ] Every `US-UI-*` references an existing `BR-*`
- [ ] Every `AC-UI-*` references an existing `US-UI-*`
- [ ] Every `COMP-*` implements existing `AC-UI-*` or `US-UI-*`
- [ ] Every `E2E-UI-*` tests existing `AC-UI-*`
- [ ] Every `SM-*` is referenced by at least one component
- [ ] No dangling references (IDs that don't exist)

**Output format:**
```
UI TRACEABILITY CHECK
=====================

✅ US-UI → BR references: 15/15 valid
✅ AC-UI → US-UI references: 42/42 valid
✅ COMP → AC-UI references: 12/12 valid
⚠️  E2E-UI → AC-UI references: 38/42 valid
   - E2E-UI-DRAG-005 references AC-UI-DRAG-099 (not found)
   - E2E-UI-FORM-003 references AC-UI-FORM-015 (not found)
✅ SM referenced: 5/5 state machines used

Result: 2 issues found
```

### 2. Format Check (`--check format`)

Verify documents follow UI conventions:

**YAML Frontmatter:**
- [ ] All UI `.md` files have valid YAML frontmatter
- [ ] Required fields: `status`, `derived-from`, `derived-at`
- [ ] SI decisions recorded if applicable

**ID Conventions:**
- [ ] UI Stories: `US-UI-{CATEGORY}-{NNN}` (e.g., US-UI-DRAG-001)
- [ ] UI AC: `AC-UI-{CATEGORY}-{NNN}-{N}` (e.g., AC-UI-DRAG-001-1)
- [ ] Components: `COMP-{NAME}` (e.g., COMP-TILE)
- [ ] State machines: `SM-{NAME}` (e.g., SM-DRAG)
- [ ] E2E tests: `E2E-UI-{CATEGORY}-{NNN}` (e.g., E2E-UI-DRAG-001)
- [ ] Visual tests: `VIS-UI-{CATEGORY}-{NNN}` (e.g., VIS-UI-TILE-001)

**Structure:**
- [ ] UI Stories have "As a... I want... so that..." format
- [ ] AC-UI uses Given/When/Then format
- [ ] Component specs have Props Interface section
- [ ] State machines have States, Transitions, Diagram sections

**Output format:**
```
UI FORMAT CHECK
===============

ui-interaction-stories.md
  ✅ YAML frontmatter valid
  ✅ ID convention: 15/15 correct
  ✅ Story format: 15/15 correct

ui-acceptance-criteria.md
  ✅ YAML frontmatter valid
  ✅ ID convention: 42/42 correct
  ✅ Given/When/Then format: 42/42

component-specs.md
  ⚠️  COMP-SIDEBAR missing Props Interface section
  ✅ ID convention: 12/12 correct

state-machines.md
  ✅ All state machines have required sections

Result: 1 issue found
```

### 3. Coverage Check (`--check coverage`)

Verify completeness of UI documentation:

**Story Coverage:**
- [ ] Every UI story category has at least one story
- [ ] Every story has at least one AC

**AC Coverage:**
- [ ] Every AC-UI has at least one E2E test
- [ ] Every AC-UI is implemented by at least one component

**Component Coverage:**
- [ ] Every component has visual tests for all states
- [ ] Every component references cross-cutting patterns

**State Machine Coverage:**
- [ ] Every state has at least one incoming transition
- [ ] Every state has at least one outgoing transition (except terminal)
- [ ] No unreachable states

**Output format:**
```
UI COVERAGE CHECK
=================

UI Story → AC Coverage
  ✅ US-UI-DRAG-001: 4 ACs
  ✅ US-UI-DRAG-002: 3 ACs
  ⚠️  US-UI-FORM-003: 0 ACs (missing coverage)

AC → E2E Test Coverage
  ✅ 38/42 AC-UIs have E2E tests (90%)
  Missing:
   - AC-UI-DRAG-004-2
   - AC-UI-FORM-001-3
   - AC-UI-FORM-002-1
   - AC-UI-NAV-001-2

Component → Visual Test Coverage
  ✅ COMP-TILE: 5/5 states covered
  ⚠️  COMP-STATION: 2/4 states covered
     Missing: Collapsed, QuickPlacement

State Machine Coverage
  ✅ SM-DRAG: All states reachable
  ✅ SM-TILE: All states reachable

Result: 90% coverage (target: 100%)
```

### 4. Consistency Check (`--check consistency`)

Verify no contradictions:

**Duplicate Detection:**
- [ ] No duplicate UI IDs within a document
- [ ] No duplicate IDs across documents
- [ ] No conflicting component names

**State Machine Consistency:**
- [ ] States referenced in transitions exist
- [ ] No contradicting state transitions
- [ ] Component states match state machine states

**SI Decision Consistency:**
- [ ] Same SI decision point has same answer across docs
- [ ] No contradicting style/pattern decisions

**Cross-Cutting Consistency:**
- [ ] Components reference existing patterns
- [ ] Pattern IDs in ui-patterns.md exist

**Output format:**
```
UI CONSISTENCY CHECK
====================

Duplicate IDs
  ✅ No duplicates found

State Machine Consistency
  ✅ SM-DRAG transitions consistent
  ⚠️  SM-TILE: COMP-TILE uses "focused" state not in SM-TILE
     Consider adding FOCUSED state to state machine

SI Decisions
  ✅ UI-STATE-1: zustand (consistent across 3 docs)
  ✅ UI-STYLE-1: tailwind (consistent across 3 docs)

Cross-Cutting References
  ✅ All pattern references valid

Result: 1 warning found
```

### 5. All Checks (`--check all`)

Run all checks and provide summary:

```
LOOM UI VALIDATION REPORT
=========================

Directory: ui/

Traceability .......... ⚠️  WARN (2 issues)
Format ................ ✅ PASS (0 issues)
Coverage .............. ⚠️  WARN (90% - target 100%)
Consistency ........... ⚠️  WARN (1 issue)

─────────────────────────────────────────
Overall: 3 warnings, 0 errors

Details:
  [TRACEABILITY] E2E-UI-DRAG-005 references AC-UI-DRAG-099 (not found)
  [TRACEABILITY] E2E-UI-FORM-003 references AC-UI-FORM-015 (not found)
  [COVERAGE] US-UI-FORM-003 has no acceptance criteria
  [COVERAGE] 4 AC-UIs missing E2E tests
  [CONSISTENCY] COMP-TILE uses "focused" state not in SM-TILE
```

## Execution Steps

### Step 1: Discover UI Documents

Find all UI documents in the specified directory:

```bash
{dir}/
├── ui-patterns.md
├── ui-interaction-stories.md
├── ui-acceptance-criteria.md
├── component-specs.md
├── state-machines.md
├── interaction-patterns/
│   ├── drag-drop.md
│   └── form-handling.md
└── tests/
    ├── e2e-tests.md
    ├── visual-tests.md
    ├── manual-qa.md
    └── accessibility-audit.md
```

### Step 2: Parse Documents

For each document:
1. Extract YAML frontmatter
2. Parse markdown structure
3. Extract all IDs and references
4. Build reference graph

### Step 3: Run Requested Checks

Based on `--check` argument:
- `traceability`: Run traceability check only
- `format`: Run format check only
- `coverage`: Run coverage check only
- `consistency`: Run consistency check only
- `all` or omitted: Run all checks

### Step 4: Report Results

Output validation results with:
- Clear pass/fail/warn status per check
- Specific issues with file and line references
- Actionable suggestions for fixes

## Integration with /loom-ui

Recommended workflow:

```bash
# 1. Generate cross-cutting patterns
/loom-ui patterns --output-dir ui/

# 2. Derive L1-UI
/loom-ui derive --level L1 --input stories.md,mockups/,br.md --output-dir ui/

# 3. Validate L1-UI
/loom-ui validate --dir ui/ --check traceability

# 4. Derive L2-UI
/loom-ui derive --level L2 --input ui-stories.md,ui-ac.md --output-dir ui/

# 5. Derive L3-UI
/loom-ui derive --level L3 --input component-specs.md,state-machines.md --output-dir ui/

# 6. Full validation before commit
/loom-ui validate --dir ui/ --check all
```

## Exit Codes (for CI)

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | Warnings found (non-blocking) |
| 2 | Errors found (blocking) |
