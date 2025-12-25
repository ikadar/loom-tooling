---
name: loom-derive-v4
description: Derive L1 documents with CLI-controlled iterative interview
version: "4.0.0"
arguments:
  - name: input-file
    description: "Path to single L0 input file"
    required: false
  - name: input-dir
    description: "Path to directory with L0 files"
    required: false
  - name: output-dir
    description: "Directory for generated L1 documents"
    required: true
  - name: batch
    description: "Skip interview, use AI-suggested defaults"
    required: false
---

# Loom Derive v4 - CLI-Controlled Iterative Interview

This command uses the **CLI to control the interview flow**, allowing:
- One question at a time
- Dependent questions automatically skipped based on answers
- Full state persistence (can resume later)

## Architecture

```
/loom-derive-v4 (this slash command)
    │
    ├─→ loom-cli analyze → analysis.json
    │
    ├─→ loom-cli interview --init → first question (exit 100)
    │       │
    │       └─→ LOOP until exit 0:
    │             ├─→ Claude presents question (from JSON output)
    │             ├─→ User answers
    │             ├─→ loom-cli interview --answer → next question
    │             └─→ (CLI may skip dependent questions automatically)
    │
    └─→ loom-cli derive → AC + BR files
```

## Execution Steps

### Step 1: Validate Arguments

Check that either `--input-file` OR `--input-dir` is provided.
Check that `--output-dir` is provided.

### Step 2: Run Analysis

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli analyze \
  --input-file "$ARGUMENTS.input-file" \
  > /tmp/loom-analysis-$$.json
```

### Step 3: Initialize Interview

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli interview \
  --init /tmp/loom-analysis-$$.json \
  --state /tmp/loom-interview-$$.json
```

This outputs JSON with the first question and exits with code 100.

### Step 4: Interview Loop

If `--batch` flag is set, accept all suggested answers automatically.

Otherwise, for each question:

1. **Parse the JSON output** from the CLI:
   ```json
   {
     "status": "question",
     "question": {
       "id": "AMB-ENT-001",
       "question": "Are there additional quote statuses?",
       "suggested_answer": "Add Draft, Expired, Cancelled",
       "options": ["Option A", "Option B", "Option C"]
     },
     "progress": "5/71",
     "skipped_count": 2
   }
   ```

2. **Present the question** to the user:
   ```
   ═══════════════════════════════════════════
   Question 5/71 (2 skipped)
   Category: entity | Subject: Quote
   ═══════════════════════════════════════════

   Are there additional quote statuses beyond 'Pending' and 'Accepted'?

   Suggested: Add Draft, Expired, Cancelled statuses

   Options:
   1. Option A
   2. Option B
   3. Option C
   4. [Accept suggested]
   5. [Custom answer]
   ```

3. **Get user response** - wait for user input

4. **Call CLI with answer**:
   ```bash
   loom-cli interview --state /tmp/loom-interview-$$.json \
     --answer '{"question_id":"AMB-ENT-001","answer":"User answer here","source":"user"}'
   ```

5. **Check exit code**:
   - `100` → More questions, repeat from step 1
   - `0` → Interview complete, proceed to derive
   - `1` → Error, show message

6. **Handle skips**: The CLI automatically skips questions based on dependencies.
   If a question is skipped, the output shows updated `skipped_count`.

### Step 5: Run Derivation

When interview exits with code 0:

```bash
/Users/istvan/Code/loom/loom-tooling/loom-cli/loom-cli derive \
  --output-dir "$ARGUMENTS.output-dir" \
  --analysis-file /tmp/loom-interview-$$.json
```

### Step 6: Display Results

Show summary and file previews.

## Batch Mode

With `--batch` flag, automatically accept all suggested answers:

```bash
# For each question until complete
while exit_code == 100:
  question = parse_json_output()
  answer = {"question_id": question.id, "answer": question.suggested_answer, "source": "auto"}
  run: loom-cli interview --state ... --answer '$answer'
```

## Skip Logic Examples

**Example 1: Deletion cascade**
```
Q1: Can quotes be deleted?
A1: "Quotes cannot be deleted"

Q2: What happens to line items when a quote is deleted? ← SKIPPED
    (depends_on Q1, skip_if "cannot be deleted")
```

**Example 2: Expiration handling**
```
Q1: Do quotes have expiration?
A1: "Yes, default 30 days"

Q2: What notification is sent before expiration? ← ASKED (not skipped)
    (depends_on Q1, skip_if "no expiration")
```

## State File Structure

The state file (`/tmp/loom-interview-$$.json`) contains:
- `domain_model`: Extracted entities, operations, relationships
- `questions`: All ambiguities with dependencies
- `decisions`: Answered questions so far
- `current_index`: Position in question list
- `skipped`: IDs of skipped questions
- `complete`: Whether interview is done

This allows:
- Resuming interrupted interviews
- Reviewing all decisions before derivation
- Modifying answers if needed

## Example Session

```
═══════════════════════════════════════════
Question 1/71 (0 skipped)
Category: entity | Subject: Quote
═══════════════════════════════════════════

Are there additional quote statuses beyond 'Pending' and 'Accepted'?

Suggested: Add Draft, Expired, Cancelled statuses

Your answer (Enter to accept, or type custom):
> yes

Decision recorded! Next question...

═══════════════════════════════════════════
Question 2/71 (0 skipped)
Category: entity | Subject: Quote
═══════════════════════════════════════════

What happens when a quote expires?
...
```

---

## Now: Execute

1. Validate arguments
2. Run loom-cli analyze
3. Run loom-cli interview --init
4. Loop: present question → get answer → loom-cli interview --answer
5. When complete (exit 0): run loom-cli derive
6. Display results
