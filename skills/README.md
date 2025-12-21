# Loom Skills

Claude Code skills for Loom document derivation and validation.

## Overview

These skills enable AI-driven document derivation following the Loom (AI-DOP) methodology. Includes separate skill chains for **backend** (`/loom`) and **UI/UX** (`/loom-ui`) derivation.

## Quick Start

### Backend Derivation (`/loom`)

```bash
# Derive documents
/loom derive --level L1 --input user-stories.md --output-dir output/
/loom derive --level L2 --input ac.md,br.md --output-dir output/
/loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir output/

# Validate documents
/loom validate --dir output/
/loom validate --dir output/ --check coverage
```

### UI/UX Derivation (`/loom-ui`)

```bash
# Generate cross-cutting patterns (run once per project)
/loom-ui patterns --output-dir ui/

# Derive UI documents
/loom-ui derive --level L1 --input stories.md,mockups/,br.md --output-dir ui/
/loom-ui derive --level L2 --input ui-stories.md,ui-ac.md --output-dir ui/
/loom-ui derive --level L3 --input component-specs.md,state-machines.md --output-dir ui/

# Validate UI documents
/loom-ui validate --dir ui/
/loom-ui validate --dir ui/ --check traceability
```

## Architecture

### Backend Skills

```
┌─────────────────────────────────────────────────────────────────┐
│                      /loom (dispatcher)                         │
│           Unified interface for derivation & validation         │
└───────────────────────────┬─────────────────────────────────────┘
                            │
          ┌─────────────────┼─────────────────┐
          │                 │                 │
          ▼                 ▼                 ▼
    ┌───────────┐     ┌───────────┐     ┌─────────────┐
    │  derive   │     │  derive   │     │  validate   │
    │  L1/L2/L3 │     │  domain   │     │             │
    └───────────┘     └───────────┘     └─────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
    ┌───────────┐     ┌───────────┐     ┌─────────────┐
    │loom-derive│     │loom-domain│     │loom-validate│
    │loom-l2/l3 │     │           │     │             │
    └───────────┘     └───────────┘     └─────────────┘
```

### UI/UX Skills

```
┌─────────────────────────────────────────────────────────────────┐
│                     /loom-ui (dispatcher)                        │
│          Unified interface for UI derivation & validation        │
└───────────────────────────┬─────────────────────────────────────┘
                            │
     ┌──────────┬───────────┼───────────┬──────────┐
     │          │           │           │          │
     ▼          ▼           ▼           ▼          ▼
┌─────────┐┌─────────┐┌─────────┐┌─────────┐┌──────────┐
│patterns ││derive L1││derive L2││derive L3││ validate │
│(once)   ││         ││         ││         ││          │
└─────────┘└─────────┘└─────────┘└─────────┘└──────────┘
     │          │           │           │          │
     ▼          ▼           ▼           ▼          ▼
┌─────────┐┌─────────┐┌─────────┐┌─────────┐┌──────────┐
│loom-ui- ││loom-ui- ││loom-ui- ││loom-ui- ││loom-ui-  │
│patterns ││derive-l1││derive-l2││derive-l3││validate  │
└─────────┘└─────────┘└─────────┘└─────────┘└──────────┘
```

### Traceability Bridge (Backend ↔ UI)

```
Backend                              UI/UX
───────                              ─────
user-stories.md ──┐
                  │
                  ▼
            business-rules.md ◄──── ui-interaction-stories.md
            (BR-* IDs)              (US-UI-* references BR-*)
                  │                        │
                  ▼                        ▼
            acceptance-criteria.md  ui-acceptance-criteria.md
                  │                        │
                  ▼                        ▼
            interface-contracts.md  component-specs.md
                  │                 state-machines.md
                  ▼                        │
            test-cases.md           e2e-tests.md
                                    visual-tests.md
```

## Available Skills

### Backend Skills

| Skill | Purpose | Input | Output |
|-------|---------|-------|--------|
| `loom.md` | Dispatcher | any | routes to specialized skill |
| `loom-derive.md` | L0 → L1 | user-stories.md | acceptance-criteria.md, business-rules.md |
| `loom-derive-domain.md` | L0 → Domain | stories + vocabulary | domain-model.md |
| `loom-derive-l2.md` | L1 → L2 | AC + BR | interface-contracts.md, sequence-design.md |
| `loom-derive-l3.md` | L2 → L3 | contracts + AC + BR | test-cases.md |
| `loom-validate.md` | Validation | all docs | validation report |

### UI/UX Skills

| Skill | Purpose | Input | Output |
|-------|---------|-------|--------|
| `loom-ui.md` | UI Dispatcher | any | routes to specialized UI skill |
| `loom-ui-patterns.md` | Cross-cutting patterns | design system | ui-patterns.md |
| `loom-ui-derive-l1.md` | L0 → L1-UI | stories + mockups + BR | ui-interaction-stories.md, ui-acceptance-criteria.md |
| `loom-ui-derive-l2.md` | L1-UI → L2-UI | UI stories + AC | component-specs.md, state-machines.md, interaction-patterns/*.md |
| `loom-ui-derive-l3.md` | L2-UI → L3-UI | components + state machines | e2e-tests.md, visual-tests.md, manual-qa.md, accessibility-audit.md |
| `loom-ui-validate.md` | UI Validation | all UI docs | validation report |

## Installation

### Copy to Project

```bash
mkdir -p my-project/.claude/skills
cp loom-tooling/skills/*.md my-project/.claude/skills/
```

### Symlink (Development)

```bash
ln -s /path/to/loom-tooling/skills my-project/.claude/skills
```

## Usage

### Unified Command (Recommended)

```bash
# Derivation
/loom derive --level <L1|domain|L2|L3> --input <file(s)> --output-dir <dir>

# Validation
/loom validate --dir <dir> [--check <check>]
```

**Derivation Examples:**

```bash
# Derive AC and BR from user stories
/loom derive --level L1 --input requirements/user-stories.md --output-dir requirements/

# Derive domain model
/loom derive --level domain --input stories.md,vocabulary.md --output-dir design/

# Derive API contracts and sequences
/loom derive --level L2 --input ac.md,br.md --output-dir design/

# Derive test cases
/loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir tests/
```

**Validation Examples:**

```bash
# Run all checks
/loom validate --dir output/

# Check only traceability
/loom validate --dir output/ --check traceability

# Check test coverage
/loom validate --dir output/ --check coverage
```

### Direct Skill Invocation (Advanced)

```bash
# Derivation skills
/loom-derive --input-file user-stories.md --output-dir output/
/loom-derive-domain --user-stories-file stories.md --output-dir output/
/loom-derive-l2 --ac-file ac.md --br-file br.md --output-dir output/
/loom-derive-l3 --contracts-file contracts.md --ac-file ac.md --br-file br.md --output-dir output/

# Validation skill
/loom-validate --dir output/ --check all
```

## Validation Checks

| Check | Description |
|-------|-------------|
| `traceability` | Verify all IDs and cross-references are valid |
| `format` | Check YAML frontmatter, ID conventions, document structure |
| `coverage` | Ensure every requirement has tests, every error code is tested |
| `consistency` | Detect duplicate IDs, contradicting SI decisions |
| `all` | Run all checks (default) |

### Validation Output

```
LOOM VALIDATION REPORT
======================

Directory: poc/booking-system/

Traceability .......... ✅ PASS (0 issues)
Format ................ ⚠️  WARN (1 issue)
Coverage .............. ✅ PASS (100%)
Consistency ........... ✅ PASS (0 issues)

─────────────────────────────────────────
Overall: 1 warning, 0 errors
```

## Recommended Workflow

### Backend Workflow

```bash
# 1. Derive L1 from user stories
/loom derive --level L1 --input user-stories.md --output-dir output/

# 2. Derive domain model
/loom derive --level domain --input stories.md,vocabulary.md --output-dir output/

# 3. Derive L2
/loom derive --level L2 --input ac.md,br.md --output-dir output/

# 4. Derive L3
/loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir output/

# 5. Validate before commit
/loom validate --dir output/

# 6. Fix any issues, commit
```

### UI/UX Workflow

```bash
# 1. Generate cross-cutting patterns (once per project)
/loom-ui patterns --output-dir ui/

# 2. Derive L1-UI from stories + mockups
/loom-ui derive --level L1 --input stories.md,mockups/,br.md --output-dir ui/

# 3. Validate L1-UI traceability
/loom-ui validate --dir ui/ --check traceability

# 4. Derive L2-UI (components, state machines)
/loom-ui derive --level L2 --input ui-stories.md,ui-ac.md --output-dir ui/

# 5. Derive L3-UI (tests)
/loom-ui derive --level L3 --input component-specs.md,state-machines.md --output-dir ui/

# 6. Full validation before commit
/loom-ui validate --dir ui/ --check all
```

### Full-Stack Workflow

For projects with both backend and UI:

```bash
# 1. Backend L1 (produces BR that UI will reference)
/loom derive --level L1 --input user-stories.md --output-dir backend/

# 2. UI patterns (once)
/loom-ui patterns --output-dir ui/

# 3. UI L1 (references BR from backend)
/loom-ui derive --level L1 --input stories.md,mockups/,backend/business-rules.md --output-dir ui/

# 4. Continue with L2/L3 for both chains...
```

## Structured Interview

All derivation skills use the **Structured Interview** pattern:

1. **Identify decision points** that need resolution
2. **Check if input provides answers** to those decision points
3. **Ask targeted questions** for any gaps
4. **Only then derive** with full explicit context

This ensures AI never makes implicit decisions about:
- Error handling strategies
- Authorization models
- Entity vs Value Object classification
- API design patterns
- Test strategies

### UI-Specific SI Questions

The UI skill chain has 21 dedicated SI questions:

| Category | Count | Examples |
|----------|-------|----------|
| Cross-cutting (CC-*) | 4 | Loading strategy, error display, validation timing, empty states |
| L1-UI (UI-COMP-*, etc.) | 7 | Component granularity, state management, styling, design system, accessibility, navigation, forms |
| L2-UI (UI-LOAD-*, etc.) | 5 | Loading state primary, error display primary, empty state style, validation timing, transitions |
| L3-UI (UI-E2E-*, etc.) | 5 | E2E framework, visual testing tool, browser matrix, a11y tools, QA process |

## PoC Results

### Backend Derivation

| Metric | Target | Achieved |
|--------|--------|----------|
| Given/When/Then format | 100% | 100% |
| Negative test ratio | ≥20% | 33% |
| "Should NOT" tests | ≥5% | 13% |
| AC coverage | 100% | 100% |
| Content expansion | 10x+ | 26x |
| SI decision points | - | 66 |

### UI/UX Derivation

| Metric | Target | Status |
|--------|--------|--------|
| UI Stories format | 100% | Pending test |
| E2E test coverage | ≥90% | Pending test |
| Visual test coverage | 100% states | Pending test |
| Accessibility coverage | WCAG 2.1 AA | Pending test |
| SI decision points | - | 21 |

### Total SI Decision Points

| Skill Chain | Count |
|-------------|-------|
| Backend | 66 |
| UI/UX | 21 |
| **Total** | **87** |

## Customization

Skills can be customized per project by editing the copied `.md` files:
- Adjust ID conventions (US-XXX, AC-XXX-X)
- Modify output templates
- Add project-specific rules
- Customize SI decision point catalogs
- Add project-specific validation rules
