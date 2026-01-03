---
title: "Loom CLI Build & Test Guide"
generated: 2025-01-03T16:00:00Z
status: draft
level: L2
---

# Loom CLI Build & Test Guide

## Overview

This document specifies build, test, and development procedures for loom-cli.

**Traceability:** Implements [package-structure.md](package-structure.md) build requirements.

---

## Prerequisites

### Required

| Requirement | Version | Purpose |
|-------------|---------|---------|
| Go | 1.21+ | Build toolchain |
| Claude Code CLI | Latest | AI backend (`claude` command) |
| ANTHROPIC_API_KEY | - | API authentication |

### Verify Installation

```bash
# Check Go version
go version

# Check Claude CLI
claude --version

# Check API key
echo $ANTHROPIC_API_KEY | head -c 10
```

---

## Build

### Development Build

```bash
cd loom-cli
go build -o loom-cli .
```

### Release Build

```bash
cd loom-cli
go build -ldflags="-s -w" -o loom-cli .
```

### Cross-Compilation

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o loom-cli-linux .

# Windows
GOOS=windows GOARCH=amd64 go build -o loom-cli.exe .
```

---

## Test Data

### Benchmark Dataset

**Location:** `test/benchmark/01-ecommerce-order/`

**Purpose:** E-commerce order domain for end-to-end testing

**Structure:**
```
01-ecommerce-order/
├── input-l0.md              # L0 user story
├── expected-entities.json   # Expected entities
├── expected-ambiguities.json # Expected ambiguities
├── expected-severity.json   # Expected severity levels
├── l1-*/                    # L1 output samples
├── l2-*/                    # L2 output samples
└── l3-*/                    # L3 output samples
```

### Test Fixtures

**Location:** `loom-cli/test/`

```
test/
├── user-story.md    # Sample L0 input
├── decisions.md     # Sample decisions
└── output/          # Sample outputs
```

---

## Manual Testing

### Full Pipeline Test

```bash
# 1. Analyze
./loom-cli analyze \
  --input-file ../test/benchmark/01-ecommerce-order/input-l0.md \
  > /tmp/analysis.json

# 2. Interview (batch mode with defaults)
./loom-cli interview \
  --input-file ../test/benchmark/01-ecommerce-order/input-l0.md \
  --state-file /tmp/state.json \
  --mode batch

# 3. Derive L1
./loom-cli derive \
  --output-dir /tmp/l1 \
  --analysis-file /tmp/analysis.json

# 4. Derive L2
./loom-cli derive-l2 \
  --input-dir /tmp/l1 \
  --output-dir /tmp/l2

# 5. Derive L3
./loom-cli derive-l3 \
  --input-dir /tmp/l2 \
  --output-dir /tmp/l3

# 6. Validate
./loom-cli validate \
  --input-dir /tmp/l1 \
  --level l1

./loom-cli validate \
  --input-dir /tmp/l2 \
  --level l2

./loom-cli validate \
  --input-dir /tmp/l3 \
  --level l3
```

### Cascade Test

```bash
./loom-cli cascade \
  --input-file ../test/benchmark/01-ecommerce-order/input-l0.md \
  --output-dir /tmp/cascade-test \
  --skip-interview
```

### Command-by-Command Tests

#### Analyze

```bash
# Single file
./loom-cli analyze --input-file test/user-story.md

# Directory
./loom-cli analyze --input-dir ../test/benchmark/01-ecommerce-order/
```

#### Interview

```bash
# Interactive mode (default)
./loom-cli interview \
  --input-file test/user-story.md \
  --state-file /tmp/state.json

# Batch mode (use defaults)
./loom-cli interview \
  --input-file test/user-story.md \
  --state-file /tmp/state.json \
  --mode batch
```

#### Validate

```bash
# Validate L1
./loom-cli validate --input-dir /tmp/l1 --level l1

# Validate L2
./loom-cli validate --input-dir /tmp/l2 --level l2

# Validate L3
./loom-cli validate --input-dir /tmp/l3 --level l3

# Validate all levels
./loom-cli validate --input-dir /tmp/all --level all
```

---

## Validation Rules

### L1 Validation (V001-V003)

| Rule | Description | Check |
|------|-------------|-------|
| V001 | Documents have IDs | AC-XXX-NNN, BR-XXX-NNN present |
| V002 | ID patterns valid | Regex match `^(AC\|BR)-[A-Z]+-\d{3}$` |
| V003 | References valid | All referenced IDs exist |

### L2 Validation (V004-V007)

| Rule | Description | Check |
|------|-------------|-------|
| V004 | Tech specs exist | tech-specs.md present |
| V005 | Every AC has tests | Test case references all ACs |
| V006 | Aggregates defined | aggregate-design.md has aggregates |
| V007 | Sequences exist | sequence-design.md has flows |

### L3 Validation (V008-V010)

| Rule | Description | Check |
|------|-------------|-------|
| V008 | Negative ratio >= 20% | Test case statistics |
| V009 | Hallucination tests | At least 1 H-test per AC |
| V010 | OpenAPI valid | openapi.json parseable |

---

## Development Workflow

### Adding a New Command

1. Create `cmd/{command}.go`
2. Add route in `cmd/root.go`
3. Add argument parser in `internal/config/config.go`
4. Update CLI help text
5. Test manually

### Adding a New Prompt

1. Create `prompts/{prompt-name}.md`
2. Add embed directive in `prompts/prompts.go`
3. Use in command via `prompts.{PromptName}`
4. Copy to `docs/spec/l2/prompts/`

### Adding a New Formatter

1. Add types in `internal/formatter/types.go`
2. Create `internal/formatter/{type}.go`
3. Implement `Format{Type}(data) string`
4. Use in derive command

---

## Debugging

### Verbose Mode

```bash
./loom-cli analyze --input-file test.md --verbose
```

### Claude Response Inspection

Set environment variable for raw output:

```bash
export LOOM_DEBUG=1
./loom-cli analyze --input-file test.md
```

### Common Issues

| Issue | Cause | Fix |
|-------|-------|-----|
| "claude: command not found" | Claude CLI not installed | Install Claude Code |
| "ANTHROPIC_API_KEY not set" | Missing API key | Export API key |
| "no JSON found in response" | LLM didn't return JSON | Check prompt, retry |
| "rate limit exceeded" | Too many API calls | Wait and retry |

---

## CI/CD Integration

### Build Step

```yaml
- name: Build
  run: |
    cd loom-cli
    go build -o loom-cli .
```

### Validation Step

```yaml
- name: Validate Specs
  run: |
    cd loom-cli
    ./loom-cli validate --input-dir ../docs/spec/l1 --level l1
    ./loom-cli validate --input-dir ../docs/spec/l2 --level l2
    ./loom-cli validate --input-dir ../docs/spec/l3 --level l3
```

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L2 | [package-structure.md](package-structure.md) | Package structure |
| L2 | [internal-api.md](internal-api.md) | Internal API |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI interface |
| L2 | This document | Build & test guide |
