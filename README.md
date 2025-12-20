# Loom Tooling

Loom tooling for AI-driven development: skills, RAG engine, and MCP server.

## Overview

This repository contains the tooling components for the Loom (AI-DOP) framework:

```
loom-tooling/
├── skills/           # Claude Code skills for derivation
├── rag/              # RAG engine for knowledge-enhanced derivation
├── templates/        # Project templates
├── docs/             # Documentation
└── mcp-server/       # MCP Server (planned)
```

## Components

### Skills (`skills/`)

Claude Code skills for document derivation:

| Skill | Purpose |
|-------|---------|
| `loom-derive.md` | L0 → L1 derivation (user stories → AC, BR) |
| `loom-derive-l2.md` | L1 → L2 derivation (AC, BR → API contracts, sequences) |
| `loom-derive-l3.md` | L2 → L3 derivation (contracts → test cases) |

**Usage in project:**
```bash
# Copy skills to your project
cp loom-tooling/skills/*.md my-project/.claude/skills/

# Use in Claude Code
/loom-derive --input-file ai-dop/requirements/user-stories.md
```

### RAG Engine (`rag/`)

Knowledge-enhanced derivation using RAG (Retrieval-Augmented Generation).

**Features:**
- Vector DB: Chroma (local, free)
- Embeddings: HuggingFace all-MiniLM-L6-v2 (local, free)
- Knowledge base: Loom guidelines

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

### Option 1: Copy Skills (Simple)

```bash
# Copy skills to your project
mkdir -p my-project/.claude/skills
cp loom-tooling/skills/*.md my-project/.claude/skills/
```

### Option 2: Symlink (Development)

```bash
# Symlink for development (skills update automatically)
ln -s /path/to/loom-tooling/skills my-project/.claude/skills
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

- v0.1.0 - Initial release (skills + RAG PoC)

## License

MIT
