# Loom RAG MCP Server

MCP server that provides RAG (Retrieval-Augmented Generation) tools for Loom commands.

## Tools

### `rag_retrieve`

Retrieve relevant context from the knowledge base.

```json
{
  "query": "acceptance criteria format",
  "sources": ["guidelines", "project"],
  "limit": 5
}
```

Returns chunks with source information and priority.

### `rag_index`

Index a new or updated document.

```json
{
  "file_path": "/path/to/decisions.md",
  "source_type": "decisions"
}
```

Source types:
- `guidelines` (priority 1) - Format and structure guidance
- `project` (priority 2) - Domain context and vocabulary
- `decisions` (priority 3) - Past SI answers (highest priority)

### `rag_get_decisions`

Get past SI decisions for a topic.

```json
{
  "topic": "Station deletion",
  "entity": "Station"
}
```

Returns matching decisions from the knowledge base.

### `rag_initialize`

Initialize RAG with specific directories.

```json
{
  "guidelines_dir": "/path/to/guidelines",
  "project_dir": "/path/to/project/docs"
}
```

## Setup

1. Install dependencies:
   ```bash
   cd /Users/istvan/Code/loom-tooling
   pip install -r mcp/requirements.txt
   ```

2. Copy MCP config to Claude Code:
   ```bash
   # For project-specific config
   cp mcp.json /path/to/your/project/.mcp.json

   # Or for global config
   cp mcp.json ~/.claude/mcp.json
   ```

3. Restart Claude Code to load the MCP server.

## Usage in Commands

Commands can now use RAG tools:

```markdown
## Phase 0: Initialize RAG

Call `rag_initialize` with:
- guidelines_dir: path to ai-pds-specification/9000-appendix/9300-guidelines
- project_dir: path to input directory

## Phase 1: Retrieve Context

Before analyzing entities, call `rag_retrieve`:
- query: "entity completeness checklist"
- sources: ["guidelines"]

## Phase 3: Check Past Decisions

Before asking a question, call `rag_get_decisions`:
- topic: "{entity_name} {aspect}"

If decision found → use it, don't ask again.

## Phase 6: Index New Decisions

After interview, call `rag_index`:
- file_path: "{input-dir}/decisions.md"
- source_type: "decisions"
```

## Architecture

```
Claude Code
    │
    ├── /loom-derive command
    │       │
    │       └── MCP Tool Calls ────────────────────┐
    │                                               │
    │                                               ▼
    │                                     ┌─────────────────┐
    │                                     │  loom-rag MCP   │
    │                                     │                 │
    │                                     │  rag_retrieve   │
    │                                     │  rag_index      │
    │                                     │  rag_get_dec... │
    │                                     └────────┬────────┘
    │                                              │
    │                                              ▼
    │                                     ┌─────────────────┐
    │                                     │   LoomRAG       │
    │                                     │   (rag_engine)  │
    │                                     └────────┬────────┘
    │                                              │
    │                                              ▼
    │                                     ┌─────────────────┐
    │                                     │   ChromaDB      │
    │                                     │   (vector DB)   │
    │                                     └─────────────────┘
    │
    └── Returns context to command
```
