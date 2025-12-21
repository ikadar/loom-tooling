# Loom Tooling

Loom tooling for AI-driven development: commands, RAG engine, and MCP server.

## Overview

This repository contains the tooling components for the Loom (AI-DOP) framework:

```
loom-tooling/
├── commands/         # Claude Code slash commands for derivation
├── rag/              # RAG engine for knowledge-enhanced derivation
├── templates/        # Project templates
├── docs/             # Documentation
└── mcp-server/       # MCP Server (planned)
```

## Components

### Commands (`commands/`)

Claude Code slash commands for document derivation with **Structured Interview** pattern:

| Command | Purpose |
|---------|---------|
| `/loom` | Backend dispatcher (routes to derive/validate) |
| `/loom-derive` | L0 → L1 derivation (user stories → AC, BR) |
| `/loom-derive-l2` | L1 → L2 derivation (AC, BR → API contracts) |
| `/loom-derive-l3` | L2 → L3 derivation (contracts → test cases) |
| `/loom-validate` | Validate all documents |
| `/loom-ui` | UI/UX dispatcher (routes to UI derive/validate) |
| `/loom-ui-patterns` | Generate cross-cutting UI patterns |
| `/loom-ui-derive-l1` | UI stories + AC from mockups |
| `/loom-ui-derive-l2` | Component specs + state machines |
| `/loom-ui-derive-l3` | E2E + visual + manual QA tests |
| `/loom-ui-validate` | Validate UI documents |

**All commands use Structured Interview (SI):** AI asks targeted questions before making decisions, preventing implicit/wrong choices.

**Usage in project:**
```bash
# Symlink commands to your project
ln -s /path/to/loom-tooling/commands my-project/.claude/commands

# Use in Claude Code
/loom derive --level L1 --input user-stories.md --output-dir output/
/loom-ui validate --dir ui/
```

### RAG Engine (`rag/`)

Knowledge-enhanced derivation using RAG (Retrieval-Augmented Generation) with **Self-Learning System**.

**Features:**
- Vector DB: Chroma (local, free)
- Embeddings: HuggingFace all-MiniLM-L6-v2 (local, free)
- Knowledge sources: Guidelines + Project docs (learns from derived documents)
- Cross-chain SI reuse: Backend and UI decisions in single knowledge base

**Setup:**
```bash
cd rag
uv venv --python 3.13
source .venv/bin/activate
uv pip install -r requirements.txt
```

**Usage:**
```bash
# Initialize vector DB with guidelines
python rag_engine.py --init --guidelines-dir /path/to/guidelines

# Retrieve context for derivation
python rag_retrieve.py "aggregate design entities value objects"
```

### MCP Server (`mcp-server/`) [Planned]

Model Context Protocol server for Loom tools integration.

**Planned tools:**
- `loom_validate` - Validate traceability and documentation
- `loom_derive` - Derive documents with RAG context
- `loom_trace` - Generate traceability maps

## Installation in Project

### Option 1: Symlink (Recommended)

```bash
# Symlink commands (updates automatically)
mkdir -p my-project/.claude
ln -s /path/to/loom-tooling/commands my-project/.claude/commands
```

### Option 2: Copy

```bash
# Copy commands to your project
mkdir -p my-project/.claude/commands
cp loom-tooling/commands/*.md my-project/.claude/commands/
```

### Option 3: MCP Server (Future)

```json
// my-project/.mcp.json
{
  "mcpServers": {
    "loom": {
      "command": "npx",
      "args": ["-y", "@loom/mcp-server"]
    }
  }
}
```

## Related Repositories

- **loom-spec** - Loom (AI-DOP) specification
- **loom-project** - Loom development project (thinking docs, evaluations)

## Version

- v0.3.0 - Commands restructure + UI/UX skill chain (2025-12-21)
- v0.2.0 - Structured Interview pattern added to all skills (2025-12-21)
- v0.1.0 - Initial release (skills + RAG PoC)

## License

MIT
