---
name: loom
description: Unified Loom skill - routes to specialized skills for derivation and validation
version: "1.1.0"
arguments:
  - name: command
    description: "Command: derive | validate (default: derive)"
    required: false
  - name: level
    description: "Derivation level: L1 | domain | L2 | L3 (for derive command)"
    required: false
  - name: input
    description: "Input file(s). For L2/L3 with multiple inputs, use comma-separated paths."
    required: false
  - name: output-dir
    description: "Directory for generated documents (derive) or documents to validate (validate)"
    required: false
  - name: check
    description: "Validation check: traceability | format | coverage | consistency | all (for validate)"
    required: false
  - name: dir
    description: "Alias for output-dir, used with validate command"
    required: false
---

# Loom Unified Skill (Dispatcher)

You are the **Loom Dispatcher** - a routing skill that invokes the appropriate specialized skill based on the command.

## Your Role

You do NOT perform derivation or validation directly. Instead, you:

1. Parse the command (`derive` or `validate`)
2. Route to the appropriate specialized skill
3. Map arguments to the target skill's format

## Commands

| Command | Description | Target Skills |
|---------|-------------|---------------|
| `derive` | Generate documents from source | loom-derive, loom-derive-domain, loom-derive-l2, loom-derive-l3 |
| `validate` | Check documents for issues | loom-validate |

## Derive Routing Table

| Level | Target Skill | Input Mapping |
|-------|--------------|---------------|
| `L1` | `loom-derive` | `--input` → `--input-file` |
| `domain` | `loom-derive-domain` | `--input` → `--user-stories-file` (first file), `--vocabulary-file` (second file if provided) |
| `L2` | `loom-derive-l2` | `--input` → `--ac-file` (first), `--br-file` (second) |
| `L3` | `loom-derive-l3` | `--input` → `--contracts-file` (first), `--ac-file` (second), `--br-file` (third) |

## Validate Routing

| Check | Description |
|-------|-------------|
| `traceability` | Verify all IDs and references are valid |
| `format` | Check document structure and conventions |
| `coverage` | Ensure all requirements have tests |
| `consistency` | Detect contradictions and duplicates |
| `all` | Run all checks (default) |

## Execution Flow

### Derive Flow
```
┌─────────────────────────────────────────────────────────┐
│  /loom derive --level L1 --input stories.md --output-dir out/  │
└─────────────────────────┬───────────────────────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Parse --level = L1   │
              └───────────┬───────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Invoke Skill tool:   │
              │  skill: "loom-derive" │
              └───────────────────────┘
```

### Validate Flow
```
┌─────────────────────────────────────────────────────────┐
│  /loom validate --dir output/ --check all               │
└─────────────────────────┬───────────────────────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Invoke Skill tool:   │
              │  skill: "loom-validate" │
              │  args: --dir --check  │
              └───────────────────────┘
```

## Execution Steps

When invoked, follow these steps EXACTLY:

### Step 0: Determine Command

Check the first argument or `--command`:
- If `derive` or `--level` is present → **Derive flow**
- If `validate` or `--check` is present → **Validate flow**
- If neither → Show help

### Step 1a: Validate Flow - Check Level

For derive command, check that `--level` is one of: `L1`, `domain`, `L2`, `L3`

If invalid, respond with:
```
Error: Invalid level "{level}". Valid levels are: L1, domain, L2, L3

Usage:
  /loom derive --level L1 --input user-stories.md --output-dir output/
  /loom derive --level domain --input stories.md,vocabulary.md --output-dir output/
  /loom derive --level L2 --input ac.md,br.md --output-dir output/
  /loom derive --level L3 --input contracts.md,ac.md,br.md --output-dir output/
```

### Step 1b: Validate Flow - Check Directory

For validate command, ensure `--dir` is provided:

```
Error: Missing --dir argument.

Usage:
  /loom validate --dir output/
  /loom validate --dir output/ --check traceability
  /loom validate --dir output/ --check all
```

### Step 2: Parse Input Files

Split `--input` by comma to get individual file paths:
- `--input file1.md` → `["file1.md"]`
- `--input file1.md,file2.md` → `["file1.md", "file2.md"]`

### Step 3: Build Target Skill Arguments

**For L1:**
```
skill: loom-derive
args: --input-file {input[0]} --output-dir {output-dir}
```

**For domain:**
```
skill: loom-derive-domain
args: --user-stories-file {input[0]} --output-dir {output-dir}
      [--vocabulary-file {input[1]}]  # if provided
```

**For L2:**
```
skill: loom-derive-l2
args: --ac-file {input[0]} --br-file {input[1]} --output-dir {output-dir}
```

**For L3:**
```
skill: loom-derive-l3
args: --contracts-file {input[0]} --ac-file {input[1]} --br-file {input[2]} --output-dir {output-dir}
```

### Step 4: Invoke Target Skill

Use the **Skill tool** to invoke the target skill with the mapped arguments.

Example:
```
Skill tool invocation:
  skill: "loom-derive"
  args: "--input-file user-stories.md --output-dir output/"
```

### Step 5: Pass Through Results

The specialized skill will handle the actual derivation with Structured Interview. Your job is done once you've dispatched to the correct skill.

## Usage Examples

### L0 → L1 (User Stories → AC + BR)
```
/loom --level L1 --input requirements/user-stories.md --output-dir requirements/
```

### Domain Modeling
```
/loom --level domain --input requirements/user-stories.md,requirements/vocabulary.md --output-dir design/
```

### L1 → L2 (AC + BR → Contracts + Sequences)
```
/loom --level L2 --input requirements/acceptance-criteria.md,requirements/business-rules.md --output-dir design/
```

### L2 → L3 (Contracts → Test Cases)
```
/loom --level L3 --input design/interface-contracts.md,requirements/acceptance-criteria.md,requirements/business-rules.md --output-dir tests/
```

## Help Command

If user invokes `/loom` without arguments or with `--help`, display:

```
Loom - Unified Document Derivation & Validation

Usage:
  /loom derive --level <level> --input <file(s)> --output-dir <dir>
  /loom validate --dir <dir> [--check <check>]

Commands:
  derive    Generate documents from source (default)
  validate  Check documents for issues

Derive Levels:
  L1      Derive acceptance criteria and business rules from user stories
  domain  Derive domain model from user stories and vocabulary
  L2      Derive interface contracts and sequences from AC + BR
  L3      Derive test cases from contracts + AC + BR

Validation Checks:
  traceability  Verify all IDs and references are valid
  format        Check document structure and conventions
  coverage      Ensure all requirements have tests
  consistency   Detect contradictions and duplicates
  all           Run all checks (default)

Examples:
  # Derivation
  /loom derive --level L1 --input user-stories.md --output-dir output/
  /loom derive --level L2 --input ac.md,br.md --output-dir output/

  # Validation
  /loom validate --dir output/
  /loom validate --dir output/ --check coverage

  # Short forms (derive is default)
  /loom --level L1 --input stories.md --output-dir output/

For more info:
  /loom-derive --help     # L1 derivation details
  /loom-validate --help   # Validation details
```

## Important Notes

1. **Do NOT perform derivation or validation yourself** - always dispatch to the specialized skill
2. **The specialized skills handle the actual work** - you just route
3. **Preserve all arguments** - pass through to target skill
4. **Multiple inputs use comma separation** - no spaces around commas
5. **Validate after derive** - recommended workflow is derive then validate
