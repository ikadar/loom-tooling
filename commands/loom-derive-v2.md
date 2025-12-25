---
name: loom-derive-v2
description: Derive L1 documents (AC + BR) from L0 specifications using the loom-cli wrapper
version: "2.0.0"
arguments:
  - name: input-file
    description: "Path to single L0 input file (use this OR --input-dir)"
    required: false
  - name: input-dir
    description: "Path to directory with L0 files (use this OR --input-file)"
    required: false
  - name: output-dir
    description: "Directory for generated L1 documents"
    required: true
---

# Loom Derive v2

This command uses the **loom-cli** binary for L0 → L1 derivation.

## Architecture

```
/loom-derive-v2 (this slash command)
    │
    └─→ loom-cli derive (Go binary with embedded prompts)
            │
            └─→ claude -p (headless mode, 6 phases)
                    │
                    └─→ L1 documents (AC + BR + decisions)
```

## Execution

### Step 1: Validate Arguments

Check that either `--input-file` OR `--input-dir` is provided (not both, not neither).
Check that `--output-dir` is provided.

### Step 2: Run loom-cli

Execute the loom-cli binary:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli derive \
  --input-file "$ARGUMENTS.input-file" \
  --output-dir "$ARGUMENTS.output-dir"
```

Or if using directory mode:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli derive \
  --input-dir "$ARGUMENTS.input-dir" \
  --output-dir "$ARGUMENTS.output-dir"
```

### Step 3: Display Results

After the CLI completes:
1. Show the summary output from the CLI
2. Read and display a preview of the generated files:
   - `{output-dir}/acceptance-criteria.md` (first 30 lines)
   - `{output-dir}/business-rules.md` (first 30 lines)
3. Mention the decisions.md file location

### Step 4: Offer Next Steps

After showing results, offer:
- "Would you like me to review the generated ACs and BRs?"
- "Would you like to run another derivation?"
- "Would you like to see the full output files?"

## Example Usage

```bash
# Single file
/loom-derive-v2 --input-file ./specs/l0/user-story.md --output-dir ./specs/l1

# Directory
/loom-derive-v2 --input-dir ./specs/l0 --output-dir ./specs/l1
```

## What This Demonstrates

| Component | Location | Visibility |
|-----------|----------|------------|
| This command | Plugin | Visible (orchestration only) |
| loom-cli binary | Wrapper | Contains secret prompts |
| claude -p | Headless | Called by wrapper |

The secret derivation prompts are embedded in the Go binary, NOT in this slash command.

---

## Now: Execute

1. Validate the arguments
2. Run the loom-cli command
3. Display results
