# Loom Skills

Claude Code skills for Loom document derivation and validation.

## Overview

These skills enable AI-driven document derivation following the Loom (AI-DOP) methodology.

## Quick Start

Use the unified `/loom` command for derivation and validation:

```bash
# Derive documents
/loom derive --level L1 --input user-stories.md --output-dir output/
/loom derive --level L2 --input ac.md,br.md --output-dir output/
/loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir output/

# Validate documents
/loom validate --dir output/
/loom validate --dir output/ --check coverage
```

## Architecture

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
    Specialized skills with focused prompts & SI catalogs
```

## Available Skills

| Skill | Purpose | Input | Output |
|-------|---------|-------|--------|
| `loom.md` | Dispatcher | any | routes to specialized skill |
| `loom-derive.md` | L0 → L1 | user-stories.md | acceptance-criteria.md, business-rules.md |
| `loom-derive-domain.md` | L0 → Domain | stories + vocabulary | domain-model.md |
| `loom-derive-l2.md` | L1 → L2 | AC + BR | interface-contracts.md, sequence-design.md |
| `loom-derive-l3.md` | L2 → L3 | contracts + AC + BR | test-cases.md |
| `loom-validate.md` | Validation | all docs | validation report |

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

## PoC Results

| Metric | Target | Achieved |
|--------|--------|----------|
| Given/When/Then format | 100% | 100% |
| Negative test ratio | ≥20% | 33% |
| "Should NOT" tests | ≥5% | 13% |
| AC coverage | 100% | 100% |
| Content expansion | 10x+ | 26x |
| SI decision points | - | 66 across all skills |

## Customization

Skills can be customized per project by editing the copied `.md` files:
- Adjust ID conventions (US-XXX, AC-XXX-X)
- Modify output templates
- Add project-specific rules
- Customize SI decision point catalogs
- Add project-specific validation rules
