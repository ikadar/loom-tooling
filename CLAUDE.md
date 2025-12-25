# loom-tooling

Kód és eszközök a Loom platformhoz.

## Struktúra

```
loom-tooling/
├── loom-cli/       # Go CLI v0.3.0
├── commands/       # Claude Code slash commands
├── mcp/            # MCP server
└── rag/            # RAG engine
```

## loom-cli

Go-based CLI a `loom-cli/` mappában.

```bash
loom-cli analyze --input-file story.md
loom-cli interview --init analysis.json
loom-cli interview --answer '{"...}'
loom-cli derive --output-dir ./out
```

Exit codes: 0=complete, 100=question available, 1=error
