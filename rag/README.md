# Loom RAG Engine

Knowledge-enhanced derivation using RAG (Retrieval-Augmented Generation).

## Overview

The RAG engine retrieves relevant guidelines from a knowledge base to improve derivation quality.

**PoC Results (validated):**

| Aspect | Without RAG | With RAG |
|--------|-------------|----------|
| Structure | 3 ad-hoc sections | 7 guideline-compliant sections |
| Entity/VO Decision | Implicit | Explicit rationale |
| Invariants | Missing | 4 identified |
| Aggregate Boundaries | Implicit | Documented |

## Setup

### Prerequisites

- Python 3.13+ (3.14 not yet supported)
- `uv` package manager (recommended)

### Installation

```bash
cd rag

# Install uv if not present
curl -LsSf https://astral.sh/uv/install.sh | sh

# Create virtual environment
uv venv --python 3.13
source .venv/bin/activate

# Install dependencies
uv pip install -r requirements.txt
```

## Usage

### Initialize Vector Database

```python
from rag_engine import LoomRAG

# Initialize with guidelines directory
rag = LoomRAG(
    guidelines_dir="/path/to/9300-guidelines",
    persist_dir="chroma_db"
)

print(f"Indexed {len(rag.vectordb._collection.count())} chunks")
```

### Retrieve Context

```bash
# CLI usage
python rag_retrieve.py "aggregate design entities value objects"

# Output: Retrieved guidelines chunks in markdown format
```

### Use with Claude Code

1. Run retrieval:
```bash
python rag_retrieve.py "domain model aggregate design" > context.md
```

2. Provide to Claude Code:
```
Here is relevant guidelines context:
[paste context.md]

Now derive a domain model for the Quote entity.
```

## Files

| File | Purpose |
|------|---------|
| `rag_engine.py` | Core RAG engine class |
| `rag_retrieve.py` | CLI retrieval script |
| `requirements.txt` | Python dependencies |
| `test_derivation.py` | Comparison test script |

## Configuration

### Tuning Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `chunk_size` | 500 | Tokens per chunk |
| `chunk_overlap` | 50 | Overlap between chunks |
| `retrieval_k` | 5 | Number of chunks to retrieve |

### Embedding Model

Using HuggingFace `all-MiniLM-L6-v2`:
- Free (no API key required)
- Local (no network calls)
- Fast (~100ms per query)

## Key Insight

RAG doesn't give "correct answers" - it provides **informed, reasoned decisions**:
- Entity vs Value Object: both can be valid depending on context
- RAG ensures decisions are **conscious and documented**
- Guidelines compliance is automatic, not manual
