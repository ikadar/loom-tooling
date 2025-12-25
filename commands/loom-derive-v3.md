---
name: loom-derive-v3
description: Derive L1 documents with interactive structured interview
version: "3.0.0"
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
  - name: batch
    description: "Skip interview, use AI-suggested defaults"
    required: false
---

# Loom Derive v3 - Interactive Interview

This command derives L1 documents (AC + BR) from L0 specifications with an **interactive structured interview** for ambiguity resolution.

## Architecture

```
/loom-derive-v3 (this slash command)
    │
    ├─→ loom-cli analyze (phases 0-3)
    │       └─→ Returns: domain model + ambiguities JSON
    │
    ├─→ [INTERACTIVE] Present questions in Claude Code session
    │       └─→ User answers questions
    │       └─→ Decisions collected
    │
    └─→ loom-cli derive (phases 5-6)
            └─→ Generates: AC + BR files
```

## Execution Steps

### Step 1: Validate Arguments

Check that either `--input-file` OR `--input-dir` is provided.
Check that `--output-dir` is provided.

### Step 2: Run Analysis Phase

Execute loom-cli analyze to discover domain model and find ambiguities:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli analyze \
  --input-file "$ARGUMENTS.input-file" \
  > /tmp/loom-analysis-$$.json
```

Or for directory mode:
```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli analyze \
  --input-dir "$ARGUMENTS.input-dir" \
  > /tmp/loom-analysis-$$.json
```

Parse the JSON output to get:
- `domain_model`: entities, operations, relationships
- `ambiguities`: unresolved questions needing answers
- `existing_decisions`: previously resolved decisions

### Step 3: Structured Interview (Interactive)

If `--batch` flag is set, skip to Step 4 with AI-suggested answers.

Otherwise, for each ambiguity in the JSON:

1. **Present the question** to the user in a clear format:
   ```
   ## Question [1/N]: {ambiguity.category} - {ambiguity.subject}

   {ambiguity.question}

   Suggested answer: {ambiguity.suggested_answer}
   Options (if any): {ambiguity.options}
   ```

2. **Wait for user response** - the user can:
   - Accept the suggested answer (press Enter or type "yes")
   - Provide their own answer
   - Skip this question for now
   - Request more context

3. **Record the decision** with:
   - The question
   - The user's answer
   - Timestamp
   - Source: "user" or "ai_suggested"

4. Continue until all ambiguities are addressed.

**Important interview guidelines:**
- Group related questions together (by entity or operation)
- Allow batch acceptance of suggested answers ("accept all remaining")
- Show progress indicator (Question 5/23)
- Allow going back to previous questions
- Save progress periodically so user can resume

### Step 4: Prepare Derive Input

Create the derive input JSON by merging:
- `domain_model` from analysis
- `decisions`: existing_decisions + new decisions from interview
- `input_content` from analysis

Save to: `/tmp/loom-derive-input-$$.json`

### Step 5: Run Derivation Phase

Execute loom-cli derive:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli derive \
  --output-dir "$ARGUMENTS.output-dir" \
  --analysis-file /tmp/loom-derive-input-$$.json
```

### Step 6: Display Results

After derivation completes:
1. Show summary of generated documents
2. Display preview of acceptance-criteria.md (first 30 lines)
3. Display preview of business-rules.md (first 30 lines)
4. Show count of new decisions recorded

### Step 7: Cleanup & Next Steps

Clean up temporary files.

Offer next steps:
- "Would you like to review the generated ACs in detail?"
- "Would you like to run another derivation?"
- "Should I save the interview decisions to your decisions.md?"

## Example Usage

```bash
# Interactive mode (recommended)
/loom-derive-v3 --input-file ./specs/l0/user-story.md --output-dir ./specs/l1

# Batch mode (use AI-suggested defaults)
/loom-derive-v3 --input-file ./specs/l0/user-story.md --output-dir ./specs/l1 --batch

# Directory of stories
/loom-derive-v3 --input-dir ./specs/l0 --output-dir ./specs/l1
```

## Interview Example Session

```
========================================
     STRUCTURED INTERVIEW (23 questions)
========================================

## Question [1/23]: Entity - Quote

What additional states can a Quote have beyond 'Pending' and 'Accepted'?
(e.g., Draft, Expired, Rejected, Cancelled, Converted to Order)

Suggested: Draft, Pending, Accepted, Rejected, Expired, Cancelled

Your answer (Enter to accept suggested, or type your answer):
> _

[User types: "yes" or provides custom answer]

Decision recorded!

## Question [2/23]: Entity - Quote
...
```

---

## Now: Execute

1. Validate arguments
2. Run loom-cli analyze
3. Parse ambiguities from JSON
4. Conduct interactive interview (or use batch mode)
5. Prepare derive input with decisions
6. Run loom-cli derive
7. Display results
