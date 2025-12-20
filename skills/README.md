# Loom Skills

Claude Code skills for Loom document derivation.

## Overview

These skills enable AI-driven document derivation following the Loom (AI-DOP) methodology.

## Available Skills

| Skill | Level | Input | Output |
|-------|-------|-------|--------|
| `loom-derive.md` | L0 → L1 | user-stories.md | acceptance-criteria.md, business-rules.md |
| `loom-derive-l2.md` | L1 → L2 | AC + BR | interface-contracts.md, sequence-design.md |
| `loom-derive-l3.md` | L2 → L3 | contracts | test-cases.md |

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

### L0 → L1 Derivation

```
/loom-derive --input-file ai-dop/requirements/user-stories.md --output-dir ai-dop/requirements
```

**Input:** User stories in standard format
**Output:**
- `acceptance-criteria.md` - Given/When/Then ACs
- `business-rules.md` - BR-XXX rules with enforcement

### L1 → L2 Derivation

```
/loom-derive-l2 --input-ac ai-dop/requirements/acceptance-criteria.md --input-br ai-dop/requirements/business-rules.md --output-dir ai-dop/system-design
```

**Output:**
- `interface-contracts.md` - REST API specifications
- `sequence-design.md` - Mermaid sequence diagrams

### L2 → L3 Derivation

```
/loom-derive-l3 --contracts ai-dop/system-design/interface-contracts.md --output-dir ai-dop/test-plan
```

**Output:**
- `test-cases.md` - Comprehensive test cases (TDAI)

## PoC Results

| Metric | Target | Achieved |
|--------|--------|----------|
| Given/When/Then format | 100% | 100% |
| Negative test ratio | ≥20% | 33% |
| "Should NOT" tests | ≥5% | 13% |
| AC coverage | 100% | 100% |
| Content expansion | 10x+ | 26x |

## Customization

Skills can be customized per project by editing the copied `.md` files:
- Adjust ID conventions (US-XXX, AC-XXX-X)
- Modify output templates
- Add project-specific rules
