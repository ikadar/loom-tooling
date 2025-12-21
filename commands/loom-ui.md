---
name: loom-ui
description: Unified UI/UX derivation skill - routes to specialized UI skills
version: "1.0.0"
arguments:
  - name: command
    description: "Command: derive | validate | patterns (default: derive)"
    required: false
  - name: level
    description: "Derivation level: L1 | L2 | L3 (for derive command)"
    required: false
  - name: input
    description: "Input file(s). Comma-separated for multiple inputs."
    required: false
  - name: output-dir
    description: "Directory for generated documents"
    required: false
  - name: check
    description: "Validation check: traceability | format | coverage | consistency | all"
    required: false
  - name: dir
    description: "Alias for output-dir, used with validate command"
    required: false
---

# Loom UI Unified Skill (Dispatcher)

You are the **Loom UI Dispatcher** - a routing skill that invokes the appropriate specialized UI skill based on the command.

## Your Role

You do NOT perform UI derivation or validation directly. Instead, you:

1. Parse the command (`derive`, `validate`, or `patterns`)
2. Route to the appropriate specialized UI skill
3. Map arguments to the target skill's format

## Commands

| Command | Description | Target Skills |
|---------|-------------|---------------|
| `derive` | Generate UI documents from source | loom-ui-derive-l1, loom-ui-derive-l2, loom-ui-derive-l3 |
| `validate` | Check UI documents for issues | loom-ui-validate |
| `patterns` | Generate cross-cutting patterns | loom-ui-patterns |

## Derive Routing Table

| Level | Target Skill | Input Mapping |
|-------|--------------|---------------|
| `L1` | `loom-ui-derive-l1` | `--input` → `--stories-file`, `--mockups-dir`, `--br-file` |
| `L2` | `loom-ui-derive-l2` | `--input` → `--ui-stories-file`, `--ui-ac-file` |
| `L3` | `loom-ui-derive-l3` | `--input` → `--component-specs-file`, `--state-machines-file` |

## Execution Flow

### Derive Flow
```
┌─────────────────────────────────────────────────────────────────┐
│  /loom-ui derive --level L1 --input stories.md,mockups/,br.md   │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
                ┌───────────────────────┐
                │  Parse --level = L1   │
                └───────────┬───────────┘
                              │
                              ▼
                ┌───────────────────────┐
                │  Invoke Skill tool:   │
                │  skill: "loom-ui-     │
                │         derive-l1"    │
                └───────────────────────┘
```

### Patterns Flow
```
┌─────────────────────────────────────────────────────────────────┐
│  /loom-ui patterns --output-dir ui/                             │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
                ┌───────────────────────┐
                │  Invoke Skill tool:   │
                │  skill: "loom-ui-     │
                │         patterns"     │
                └───────────────────────┘
```

### Validate Flow
```
┌─────────────────────────────────────────────────────────────────┐
│  /loom-ui validate --dir ui/ --check all                        │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
                ┌───────────────────────┐
                │  Invoke Skill tool:   │
                │  skill: "loom-ui-     │
                │         validate"     │
                └───────────────────────┘
```

## Execution Steps

### Step 0: Determine Command

Check the first argument or command:
- If `derive` or `--level` is present → **Derive flow**
- If `patterns` is present → **Patterns flow**
- If `validate` or `--check` is present → **Validate flow**
- If neither → Show help

### Step 1a: Derive Flow - Check Level

For derive command, check that `--level` is one of: `L1`, `L2`, `L3`

If invalid, respond with:
```
Error: Invalid level "{level}". Valid levels are: L1, L2, L3

Usage:
  /loom-ui derive --level L1 --input stories.md,mockups/,br.md --output-dir ui/
  /loom-ui derive --level L2 --input ui-stories.md,ui-ac.md --output-dir ui/
  /loom-ui derive --level L3 --input component-specs.md,state-machines.md --output-dir ui/
```

### Step 1b: Patterns Flow

For patterns command, ensure `--output-dir` is provided.

### Step 1c: Validate Flow

For validate command, ensure `--dir` is provided.

### Step 2: Parse Input Files

Split `--input` by comma to get individual file paths:
- `--input file1.md` → `["file1.md"]`
- `--input file1.md,mockups/,file2.md` → `["file1.md", "mockups/", "file2.md"]`

### Step 3: Build Target Skill Arguments

**For L1:**
```
skill: loom-ui-derive-l1
args: --stories-file {input[0]} --mockups-dir {input[1]} --br-file {input[2]} --output-dir {output-dir}
```

**For L2:**
```
skill: loom-ui-derive-l2
args: --ui-stories-file {input[0]} --ui-ac-file {input[1]} --output-dir {output-dir}
```

**For L3:**
```
skill: loom-ui-derive-l3
args: --component-specs-file {input[0]} --state-machines-file {input[1]} --output-dir {output-dir}
```

**For patterns:**
```
skill: loom-ui-patterns
args: --output-dir {output-dir}
```

**For validate:**
```
skill: loom-ui-validate
args: --dir {dir} --check {check}
```

### Step 4: Invoke Target Skill

Use the **Skill tool** to invoke the target skill with the mapped arguments.

## Usage Examples

### Cross-Cutting Patterns (once per project)
```bash
/loom-ui patterns --output-dir ui/
```

### L1-UI: User Stories + Mockups → UI Interaction Stories
```bash
/loom-ui derive --level L1 --input user-stories.md,mockups/,business-rules.md --output-dir ui/
```

### L2-UI: UI Stories → Component Specs + State Machines
```bash
/loom-ui derive --level L2 --input ui/ui-interaction-stories.md,ui/ui-acceptance-criteria.md --output-dir ui/
```

### L3-UI: Component Specs → Tests
```bash
/loom-ui derive --level L3 --input ui/component-specs.md,ui/state-machines.md --output-dir ui/tests/
```

### Validate UI Documents
```bash
/loom-ui validate --dir ui/ --check all
```

## Help Command

If user invokes `/loom-ui` without arguments or with `--help`, display:

```
Loom UI - Frontend/UX Document Derivation & Validation

Usage:
  /loom-ui derive --level <level> --input <file(s)> --output-dir <dir>
  /loom-ui validate --dir <dir> [--check <check>]
  /loom-ui patterns --output-dir <dir>

Commands:
  derive    Generate UI documents from source (default)
  validate  Check UI documents for issues
  patterns  Generate cross-cutting UI patterns (once per project)

Derive Levels:
  L1      UI Interaction Stories + AC from user stories and mockups
  L2      Component Specs + State Machines from UI stories
  L3      E2E Tests + Visual Tests + Manual QA from component specs

Validation Checks:
  traceability  Verify all UI IDs and references are valid
  format        Check document structure and conventions
  coverage      Ensure all UI requirements have tests
  consistency   Detect contradictions and duplicates
  all           Run all checks (default)

Examples:
  # Generate cross-cutting patterns first
  /loom-ui patterns --output-dir ui/

  # Derive L1-UI from stories and mockups
  /loom-ui derive --level L1 --input stories.md,mockups/,br.md --output-dir ui/

  # Derive L2-UI from UI stories
  /loom-ui derive --level L2 --input ui-stories.md,ui-ac.md --output-dir ui/

  # Validate UI documents
  /loom-ui validate --dir ui/

For more info:
  /loom-ui-patterns --help     # Cross-cutting patterns
  /loom-ui-derive-l1 --help    # L1-UI derivation
  /loom-ui-validate --help     # Validation details
```

## Important Notes

1. **Do NOT perform derivation yourself** - always dispatch to the specialized skill
2. **Patterns first** - Run `/loom-ui patterns` before first L1 derivation
3. **Business Rules are the bridge** - L1-UI requires BR file for traceability
4. **Mockups required for L1** - UI derivation needs visual reference
5. **Validate after derive** - recommended workflow
