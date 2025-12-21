# Loom Commands

Claude Code slash commands for Loom document derivation and validation.

## Overview

These commands enable AI-driven document derivation following the Loom (AI-DOP) methodology. Use explicit `/command` invocation for predictable workflows.

## Quick Start

### Backend Derivation (`/loom`)

```bash
# Derive documents
/loom derive --level L1 --input user-stories.md --output-dir output/
/loom derive --level L2 --input ac.md,br.md --output-dir output/
/loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir output/

# Validate documents
/loom validate --dir output/
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
```

## Available Commands

### Backend Commands

| Command | Purpose | Input | Output |
|---------|---------|-------|--------|
| `/loom` | Dispatcher | any | routes to specialized command |
| `/loom-derive` | L0 → L1 | user-stories.md | acceptance-criteria.md, business-rules.md |
| `/loom-derive-domain` | L0 → Domain | stories + vocabulary | domain-model.md |
| `/loom-derive-l2` | L1 → L2 | AC + BR | interface-contracts.md, sequence-design.md |
| `/loom-derive-l3` | L2 → L3 | contracts + AC + BR | test-cases.md |
| `/loom-validate` | Validation | all docs | validation report |

### UI/UX Commands

| Command | Purpose | Input | Output |
|---------|---------|-------|--------|
| `/loom-ui` | UI Dispatcher | any | routes to specialized UI command |
| `/loom-ui-patterns` | Cross-cutting patterns | design system | ui-patterns.md |
| `/loom-ui-derive-l1` | L0 → L1-UI | stories + mockups + BR | ui-interaction-stories.md, ui-acceptance-criteria.md |
| `/loom-ui-derive-l2` | L1-UI → L2-UI | UI stories + AC | component-specs.md, state-machines.md |
| `/loom-ui-derive-l3` | L2-UI → L3-UI | components + state machines | e2e-tests.md, visual-tests.md, manual-qa.md |
| `/loom-ui-validate` | UI Validation | all UI docs | validation report |

## Installation

### Symlink (Recommended)

```bash
mkdir -p my-project/.claude
ln -s /path/to/loom-tooling/commands my-project/.claude/commands
```

### Copy

```bash
mkdir -p my-project/.claude/commands
cp loom-tooling/commands/*.md my-project/.claude/commands/
```

## Structured Interview (SI)

All derivation commands use the **Structured Interview** pattern - asking explicit questions before making decisions. SI decisions are recorded in YAML frontmatter for reuse.

| Chain | SI Questions |
|-------|--------------|
| Backend | 66 |
| UI/UX | 21 |
| **Total** | **87** |

## Traceability

Backend and UI chains share Business Rules (BR-*) as the traceability bridge:

```
Backend: user-stories → BR → AC → contracts → tests
                         ↑
UI/UX:   mockups ────────┴──→ UI-stories → components → UI-tests
```
