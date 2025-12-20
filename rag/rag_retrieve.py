"""
RAG Retrieve - Simple retrieval script for Claude Code integration.

Usage:
    python rag_retrieve.py "your query here"
    python rag_retrieve.py "derive domain model with aggregates" --k 5
"""

import sys
from pathlib import Path
from rag_engine import LoomRAG, get_default_guidelines_path


def main():
    if len(sys.argv) < 2:
        print("Usage: python rag_retrieve.py 'your query' [--k N]")
        sys.exit(1)

    query = sys.argv[1]
    k = 5  # default

    if "--k" in sys.argv:
        k_idx = sys.argv.index("--k")
        k = int(sys.argv[k_idx + 1])

    # Initialize RAG
    guidelines_path = get_default_guidelines_path()
    rag = LoomRAG(
        guidelines_dir=guidelines_path,
        persist_dir=str(Path(__file__).parent / "chroma_db"),
    )

    # Retrieve
    results = rag.retrieve_with_sources(query, k=k)

    # Output as markdown
    print(f"# RAG Context for: {query}\n")
    print(f"Retrieved {len(results)} relevant chunks from guidelines.\n")
    print("---\n")

    for i, item in enumerate(results, 1):
        source = Path(item["source"]).name
        print(f"## [{i}] {source}\n")
        print(item["content"])
        print("\n---\n")


if __name__ == "__main__":
    main()
