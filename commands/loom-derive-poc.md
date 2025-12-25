---
name: loom-derive-poc
description: POC - Derive L1 from L0 using CLI wrapper (demonstrates Plugin + Headless architecture)
version: "0.1.0"
arguments:
  - name: input-file
    description: "Path to L0 input file (user story markdown)"
    required: true
  - name: output-dir
    description: "Directory for generated L1 documents (optional, displays if not provided)"
    required: false
---

# Loom Derive POC

This is a **proof of concept** demonstrating the Plugin + Headless CLI Wrapper architecture.

## Architecture

```
This slash command (Plugin - orchestration)
    │
    └─→ loom-cli-poc derive (Wrapper - contains secret prompt)
            │
            └─→ claude -p "secret prompt + input" (Headless mode)
```

## Execution

### Step 1: Call the CLI wrapper

Run the loom-cli-poc binary to perform derivation:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli-poc/loom-cli-poc derive "$ARGUMENTS.input-file" --format text
```

### Step 2: Display results

Show the derivation output to the user.

### Step 3: Optionally save to file

If `output-dir` is provided:
1. Create the directory if it doesn't exist
2. Save the output to `{output-dir}/derived-l1.md`
3. Confirm the file was written

If `output-dir` is NOT provided:
- Just display the results (don't save)

## Example Usage

```bash
# Display only
/loom-derive-poc --input-file ./user-story.md

# Save to directory
/loom-derive-poc --input-file ./user-story.md --output-dir ./output
```

## What This Demonstrates

1. **Plugin**: This markdown file (visible, orchestration only)
2. **CLI Wrapper**: The Go binary (contains secret derivation prompt)
3. **Headless Mode**: Claude -p called by the binary
4. **Separation**: User sees this command, NOT the secret prompt

---

## Now: Execute Step 1

Call the CLI wrapper with the input file.
