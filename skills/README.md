# Loom Skills

Claude Code skills for Loom document derivation.

## Overview

These skills enable AI-driven document derivation following the Loom (AI-DOP) methodology.

## Quick Start

Use the unified `/loom` command for all derivations:

```bash
# L0 → L1: User stories to acceptance criteria + business rules
/loom --level L1 --input user-stories.md --output-dir output/

# Domain modeling
/loom --level domain --input stories.md,vocabulary.md --output-dir output/

# L1 → L2: AC + BR to interface contracts + sequences
/loom --level L2 --input ac.md,br.md --output-dir output/

# L2 → L3: Contracts to test cases
/loom --level L3 --input contracts.md,ac.md,br.md --output-dir output/
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    /loom (dispatcher)                   │
│              Unified interface for all levels           │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┬─────────────┐
        ▼             ▼             ▼             ▼
┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐
│loom-derive│  │loom-domain│  │ loom-l2   │  │ loom-l3   │
│   (L1)    │  │           │  │           │  │           │
└───────────┘  └───────────┘  └───────────┘  └───────────┘
   Specialized skills with focused prompts & SI catalogs
```

The `/loom` dispatcher routes to specialized skills, preserving their focused intelligence while providing a unified UX.

## Available Skills

| Skill | Level | Input | Output |
|-------|-------|-------|--------|
| `loom.md` | Dispatcher | any | routes to specialized skill |
| `loom-derive.md` | L0 → L1 | user-stories.md | acceptance-criteria.md, business-rules.md |
| `loom-derive-domain.md` | L0 → Domain | stories + vocabulary | domain-model.md |
| `loom-derive-l2.md` | L1 → L2 | AC + BR | interface-contracts.md, sequence-design.md |
| `loom-derive-l3.md` | L2 → L3 | contracts + AC + BR | test-cases.md |

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
/loom --level <L1|domain|L2|L3> --input <file(s)> --output-dir <dir>
```

**Examples:**

```bash
# Derive AC and BR from user stories
/loom --level L1 --input requirements/user-stories.md --output-dir requirements/

# Derive domain model
/loom --level domain --input requirements/user-stories.md,requirements/vocabulary.md --output-dir design/

# Derive API contracts and sequences
/loom --level L2 --input requirements/acceptance-criteria.md,requirements/business-rules.md --output-dir design/

# Derive test cases
/loom --level L3 --input design/interface-contracts.md,requirements/acceptance-criteria.md,requirements/business-rules.md --output-dir tests/
```

### Direct Skill Invocation (Advanced)

You can also invoke specialized skills directly for more control:

```bash
# L0 → L1
/loom-derive --input-file user-stories.md --output-dir output/

# Domain
/loom-derive-domain --user-stories-file stories.md --vocabulary-file vocab.md --output-dir output/

# L1 → L2
/loom-derive-l2 --ac-file ac.md --br-file br.md --output-dir output/

# L2 → L3
/loom-derive-l3 --contracts-file contracts.md --ac-file ac.md --br-file br.md --output-dir output/
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
