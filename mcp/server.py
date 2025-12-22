"""
Loom RAG MCP Server

Provides RAG tools for Claude Code commands:
- rag_retrieve: Retrieve relevant context from knowledge base
- rag_index: Index a new or updated document
- rag_get_decisions: Get past SI decisions for a topic

Usage:
    python server.py

Configure in Claude Code with .mcp.json
"""

import sys
from pathlib import Path

# Add parent directory to path to import rag module
sys.path.insert(0, str(Path(__file__).parent.parent))

from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent
import json

from rag.rag_engine import LoomRAG, KnowledgeSource, create_self_learning_rag


# Global RAG instance (initialized on first use)
_rag_instance: LoomRAG | None = None


def get_rag() -> LoomRAG:
    """Get or create the RAG instance."""
    global _rag_instance

    if _rag_instance is None:
        # Default configuration
        # These paths will be configured via environment or first call
        _rag_instance = LoomRAG(
            knowledge_sources=[],
            persist_dir=str(Path(__file__).parent.parent / "data" / "chroma_db")
        )

    return _rag_instance


def initialize_rag(
    guidelines_dir: str | None = None,
    project_dir: str | None = None,
    persist_dir: str | None = None
) -> LoomRAG:
    """Initialize RAG with specific directories."""
    global _rag_instance

    sources = []

    if guidelines_dir:
        sources.append(
            KnowledgeSource(guidelines_dir, KnowledgeSource.GUIDELINES, priority=1)
        )

    if project_dir:
        sources.append(
            KnowledgeSource(project_dir, KnowledgeSource.PROJECT, priority=2)
        )

    _rag_instance = LoomRAG(
        knowledge_sources=sources,
        persist_dir=persist_dir or str(Path(__file__).parent.parent / "data" / "chroma_db")
    )

    return _rag_instance


# Create MCP server
server = Server("loom-rag")


@server.list_tools()
async def list_tools() -> list[Tool]:
    """List available RAG tools."""
    return [
        Tool(
            name="rag_retrieve",
            description="Retrieve relevant context from the Loom knowledge base. Use for getting guidelines, past decisions, and project context during derivation.",
            inputSchema={
                "type": "object",
                "properties": {
                    "query": {
                        "type": "string",
                        "description": "The search query (e.g., 'acceptance criteria format', 'entity deletion cascade')"
                    },
                    "sources": {
                        "type": "array",
                        "items": {"type": "string", "enum": ["guidelines", "project", "decisions"]},
                        "description": "Which sources to search (default: all)"
                    },
                    "limit": {
                        "type": "integer",
                        "description": "Number of results to return (default: 5)",
                        "default": 5
                    }
                },
                "required": ["query"]
            }
        ),
        Tool(
            name="rag_index",
            description="Index a new or updated document into the knowledge base. Call after creating or updating project documents.",
            inputSchema={
                "type": "object",
                "properties": {
                    "file_path": {
                        "type": "string",
                        "description": "Path to the file to index"
                    },
                    "source_type": {
                        "type": "string",
                        "enum": ["guidelines", "project", "decisions"],
                        "description": "Type of source (determines priority)"
                    }
                },
                "required": ["file_path", "source_type"]
            }
        ),
        Tool(
            name="rag_get_decisions",
            description="Get past SI decisions for a specific topic or entity. Use to check if a decision has already been made.",
            inputSchema={
                "type": "object",
                "properties": {
                    "topic": {
                        "type": "string",
                        "description": "The topic to search for (e.g., 'Station deletion', 'TimeSlot entity vs VO')"
                    },
                    "entity": {
                        "type": "string",
                        "description": "Optional: specific entity name to filter by"
                    }
                },
                "required": ["topic"]
            }
        ),
        Tool(
            name="rag_initialize",
            description="Initialize the RAG engine with specific directories. Call once at the start of a derivation session.",
            inputSchema={
                "type": "object",
                "properties": {
                    "guidelines_dir": {
                        "type": "string",
                        "description": "Path to guidelines directory"
                    },
                    "project_dir": {
                        "type": "string",
                        "description": "Path to project documents directory"
                    },
                    "persist_dir": {
                        "type": "string",
                        "description": "Path to persist vector database (optional)"
                    }
                },
                "required": []
            }
        )
    ]


@server.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    """Handle tool calls."""

    if name == "rag_retrieve":
        return await handle_retrieve(arguments)
    elif name == "rag_index":
        return await handle_index(arguments)
    elif name == "rag_get_decisions":
        return await handle_get_decisions(arguments)
    elif name == "rag_initialize":
        return await handle_initialize(arguments)
    else:
        return [TextContent(type="text", text=f"Unknown tool: {name}")]


async def handle_retrieve(arguments: dict) -> list[TextContent]:
    """Handle rag_retrieve tool call."""
    query = arguments["query"]
    limit = arguments.get("limit", 5)
    sources_filter = arguments.get("sources", None)

    rag = get_rag()

    # Retrieve with sources
    results = rag.retrieve_with_sources(query, k=limit)

    # Filter by source type if specified
    if sources_filter:
        results = [r for r in results if r["source_type"] in sources_filter]

    # Format results
    formatted = []
    for i, result in enumerate(results, 1):
        formatted.append({
            "rank": i,
            "content": result["content"],
            "source": result["source"],
            "source_type": result["source_type"],
            "priority": result["priority"]
        })

    return [TextContent(
        type="text",
        text=json.dumps({
            "query": query,
            "results_count": len(formatted),
            "results": formatted
        }, indent=2)
    )]


async def handle_index(arguments: dict) -> list[TextContent]:
    """Handle rag_index tool call."""
    file_path = arguments["file_path"]
    source_type = arguments["source_type"]

    # Validate file exists
    path = Path(file_path)
    if not path.exists():
        return [TextContent(
            type="text",
            text=json.dumps({"error": f"File not found: {file_path}"})
        )]

    rag = get_rag()

    # Map source type to priority
    priority_map = {
        "guidelines": 1,
        "project": 2,
        "decisions": 3  # Highest priority for explicit decisions
    }
    priority = priority_map.get(source_type, 1)

    # Add as knowledge source
    source = KnowledgeSource(
        str(path.parent),  # Index the directory containing the file
        source_type,
        priority=priority
    )

    try:
        rag.add_knowledge_source(source, rebuild=True)

        return [TextContent(
            type="text",
            text=json.dumps({
                "success": True,
                "indexed": file_path,
                "source_type": source_type,
                "priority": priority
            })
        )]
    except Exception as e:
        return [TextContent(
            type="text",
            text=json.dumps({"error": str(e)})
        )]


async def handle_get_decisions(arguments: dict) -> list[TextContent]:
    """Handle rag_get_decisions tool call."""
    topic = arguments["topic"]
    entity = arguments.get("entity", "")

    rag = get_rag()

    # Build query
    query = f"decision {topic}"
    if entity:
        query = f"{entity} {query}"

    # Search in project/decisions sources
    results = rag.retrieve_prioritized(query, k=5, prefer_project=True)

    # Filter for decision-like content
    decisions = []
    for result in results:
        content = result["content"].lower()
        # Look for decision indicators
        if any(indicator in content for indicator in ["decision", "answer:", "a:", "resolved", "decided"]):
            decisions.append({
                "content": result["content"],
                "source": result["source"],
                "source_type": result["source_type"]
            })

    return [TextContent(
        type="text",
        text=json.dumps({
            "topic": topic,
            "entity": entity if entity else None,
            "decisions_found": len(decisions),
            "decisions": decisions
        }, indent=2)
    )]


async def handle_initialize(arguments: dict) -> list[TextContent]:
    """Handle rag_initialize tool call."""
    guidelines_dir = arguments.get("guidelines_dir")
    project_dir = arguments.get("project_dir")
    persist_dir = arguments.get("persist_dir")

    try:
        rag = initialize_rag(
            guidelines_dir=guidelines_dir,
            project_dir=project_dir,
            persist_dir=persist_dir
        )

        return [TextContent(
            type="text",
            text=json.dumps({
                "success": True,
                "guidelines_dir": guidelines_dir,
                "project_dir": project_dir,
                "knowledge_sources": len(rag.knowledge_sources)
            })
        )]
    except Exception as e:
        return [TextContent(
            type="text",
            text=json.dumps({"error": str(e)})
        )]


async def main():
    """Run the MCP server."""
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            server.create_initialization_options()
        )


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
