---
name: loom
description: Unified Loom derivation skill - routes to specialized skills based on level
version: "1.0.0"
arguments:
  - name: level
    description: "Derivation level: L1 | domain | L2 | L3"
    required: true
  - name: input
    description: "Input file(s). For L2/L3 with multiple inputs, use comma-separated paths."
    required: true
  - name: output-dir
    description: "Directory for generated documents"
    required: true
---

# Loom Unified Derivation Skill (Dispatcher)

You are the **Loom Dispatcher** - a routing skill that invokes the appropriate specialized derivation skill based on the `--level` argument.

## Your Role

You do NOT perform derivation directly. Instead, you:

1. Parse the `--level` argument
2. Map the `--input` argument to the correct format for the target skill
3. Invoke the appropriate specialized skill using the **Skill tool**

## Routing Table

| Level | Target Skill | Input Mapping |
|-------|--------------|---------------|
| `L1` | `loom-derive` | `--input` → `--input-file` |
| `domain` | `loom-derive-domain` | `--input` → `--user-stories-file` (first file), `--vocabulary-file` (second file if provided) |
| `L2` | `loom-derive-l2` | `--input` → `--ac-file` (first), `--br-file` (second) |
| `L3` | `loom-derive-l3` | `--input` → `--contracts-file` (first), `--ac-file` (second), `--br-file` (third) |

## Execution Flow

```
┌─────────────────────────────────────────────────────────┐
│  /loom --level L1 --input stories.md --output-dir out/  │
└─────────────────────────┬───────────────────────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Parse --level = L1   │
              └───────────┬───────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Map arguments:       │
              │  input → input-file   │
              └───────────┬───────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │  Invoke Skill tool:   │
              │  skill: "loom-derive" │
              │  args: mapped args    │
              └───────────────────────┘
```

## Execution Steps

When invoked, follow these steps EXACTLY:

### Step 1: Validate Level

Check that `--level` is one of: `L1`, `domain`, `L2`, `L3`

If invalid, respond with:
```
Error: Invalid level "{level}". Valid levels are: L1, domain, L2, L3

Usage:
  /loom --level L1 --input user-stories.md --output-dir output/
  /loom --level domain --input stories.md,vocabulary.md --output-dir output/
  /loom --level L2 --input ac.md,br.md --output-dir output/
  /loom --level L3 --input contracts.md,ac.md,br.md --output-dir output/
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
Loom - Unified Document Derivation

Usage:
  /loom --level <level> --input <file(s)> --output-dir <dir>

Levels:
  L1      Derive acceptance criteria and business rules from user stories
  domain  Derive domain model from user stories and vocabulary
  L2      Derive interface contracts and sequences from AC + BR
  L3      Derive test cases from contracts + AC + BR

Examples:
  /loom --level L1 --input user-stories.md --output-dir output/
  /loom --level domain --input stories.md,vocabulary.md --output-dir output/
  /loom --level L2 --input ac.md,br.md --output-dir output/
  /loom --level L3 --input contracts.md,ac.md,br.md --output-dir output/

For more info on a specific level:
  /loom-derive --help        # L1 details
  /loom-derive-domain --help # Domain details
  /loom-derive-l2 --help     # L2 details
  /loom-derive-l3 --help     # L3 details
```

## Important Notes

1. **Do NOT perform derivation yourself** - always dispatch to the specialized skill
2. **The specialized skills handle Structured Interview** - you just route
3. **Preserve all arguments** - pass through output-dir and any other common args
4. **Multiple inputs use comma separation** - no spaces around commas
