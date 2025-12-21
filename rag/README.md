# Loom RAG Engine - Self-Learning System

Knowledge-enhanced derivation using RAG (Retrieval-Augmented Generation) with support for the **Self-Learning System** pattern.

## Overview

The RAG engine retrieves relevant context from multiple knowledge sources:

1. **Guidelines** (global) - Format, structure, best practices
2. **Project docs** (project-specific) - Previous derivations with SI decisions
3. **Code patterns** (optional) - Implementation patterns from source code

```
┌────────────────────────────────────────────────────────────────────┐
│                     SELF-LEARNING SYSTEM                           │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌──────────────┐                                                 │
│  │  Guidelines  │ ─────────────┐                                  │
│  │  (global)    │              │                                  │
│  └──────────────┘              │                                  │
│                                ▼                                  │
│  ┌──────────────┐       ┌──────────────┐       ┌──────────────┐  │
│  │   Project    │ ─────►│  RAG Engine  │◄──────│   Derive     │  │
│  │   Docs       │       │  (ChromaDB)  │       │   Request    │  │
│  └──────────────┘       └──────┬───────┘       └──────────────┘  │
│         ▲                      │                                  │
│         │                      │ retrieve context                 │
│         │                      ▼                                  │
│         │               ┌──────────────┐                         │
│         │               │   Derived    │                         │
│         └───────────────│   Document   │ (includes SI decisions) │
│           re-index      │   (output)   │                         │
│                         └──────────────┘                         │
│                                                                    │
│  The system learns from its own output!                           │
└────────────────────────────────────────────────────────────────────┘
```

## Key Concept: SI Decision Reuse

Previous Structured Interview decisions are stored in document frontmatter:

```yaml
# acceptance-criteria.md
---
structured-interview:
  decisions:
    EH-1: blocking-error
    AU-1: customer-or-provider
---
```

When re-deriving or deriving related documents, the RAG engine:
1. Finds previous decisions in indexed project docs
2. Prioritizes project docs over guidelines
3. Reuses decisions automatically (or asks to override)

**No separate cache file needed!**

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

### Basic Usage (Backward Compatible)

```python
from rag_engine import LoomRAG

# Initialize with guidelines only
rag = LoomRAG(guidelines_dir="/path/to/9300-guidelines")

# Retrieve context
results = rag.retrieve("aggregate design entities", k=5)
```

### Self-Learning System Setup

```python
from rag_engine import create_self_learning_rag, KnowledgeSource

# Create with multiple sources
rag = create_self_learning_rag(
    guidelines_dir="/path/to/9300-guidelines",
    project_dir="/path/to/project/docs",
    code_dir="/path/to/src",  # Optional
    persist_dir="./chroma_db"
)

# Or manually configure sources
from rag_engine import LoomRAG, KnowledgeSource

rag = LoomRAG(knowledge_sources=[
    KnowledgeSource("/path/to/guidelines", "guidelines", priority=1),
    KnowledgeSource("/path/to/project", "project", priority=2),
])
```

### Retrieve Previous SI Decisions

```python
# Check if a decision was already made
decision = rag.retrieve_si_decision("EH-1", domain="booking-system")

if decision:
    print(f"Found: {decision['decision_id']} in {decision['source']}")
    # Reuse the decision
else:
    # Ask the user
    pass
```

### Prioritized Retrieval

```python
# Get context, prioritizing project docs over guidelines
results = rag.retrieve_prioritized(
    query="error handling race condition",
    k=5,
    prefer_project=True  # Project docs first
)

for result in results:
    print(f"[{result['source_type']}] {result['content'][:100]}...")
```

### Refresh After Derivation

```python
# After deriving new documents, refresh the index
rag.refresh_project_knowledge("/path/to/project/docs")
```

## Knowledge Source Types

| Type | Priority | Description |
|------|----------|-------------|
| `guidelines` | 1 (low) | Global format and structure guidance |
| `project` | 2 (high) | Project-specific docs with SI decisions |
| `code` | 1 (low) | Implementation patterns (reference) |

Higher priority sources are preferred when conflicts exist.

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

## PoC Results

| Aspect | Without RAG | With RAG | With Self-Learning |
|--------|-------------|----------|-------------------|
| Structure | 3 ad-hoc sections | 7 guideline-compliant | 7 sections |
| Entity/VO Decision | Implicit | Explicit rationale | Reused from previous |
| SI Questions | All asked | All asked | Only new ones |
| Consistency | Variable | Good | Excellent |

## Files

| File | Purpose |
|------|---------|
| `rag_engine.py` | Core RAG engine with Self-Learning support |
| `rag_retrieve.py` | CLI retrieval script |
| `requirements.txt` | Python dependencies |
| `test_derivation.py` | Comparison test script |

## Integration with Commands

The derivation commands can use the RAG engine:

```python
# In command execution
rag = create_self_learning_rag(
    guidelines_dir="specs-for-ai/ai-pds-specification/9300-guidelines",
    project_dir="project/docs"
)

# Check for cached SI decisions before asking
for decision_point in required_decisions:
    cached = rag.retrieve_si_decision(decision_point.id)
    if cached:
        # Use cached decision
        decision_point.answer = extract_answer(cached)
    else:
        # Add to questions list
        questions.append(decision_point)
```

## Full-Stack Self-Learning

The Self-Learning System works across **both backend and UI/UX** document chains. All derived documents feed into the same knowledge base:

```
┌─────────────────────────────────────────────────────────────────┐
│                 FULL-STACK SELF-LEARNING                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Backend Docs                    UI/UX Docs                     │
│  ┌─────────────┐                ┌─────────────┐                │
│  │ AC, BR, L2  │                │ UI-AC, L2-UI│                │
│  │ SI: EH-1,   │                │ SI: CC-1,   │                │
│  │     AU-1    │                │     UI-DS-1 │                │
│  └──────┬──────┘                └──────┬──────┘                │
│         │                              │                        │
│         └──────────┬───────────────────┘                        │
│                    ▼                                            │
│           ┌───────────────┐                                    │
│           │  RAG Engine   │ ◄── Single knowledge base          │
│           │  (ChromaDB)   │                                    │
│           └───────────────┘                                    │
│                    │                                            │
│         ┌─────────┴─────────┐                                  │
│         ▼                   ▼                                  │
│    Backend derive      UI derive                               │
│    (reuses EH-1,       (reuses CC-1,                          │
│     finds UI-DS-1)      finds AU-1)                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Indexed Document Types

| Chain | Documents | SI Prefixes |
|-------|-----------|-------------|
| Backend | acceptance-criteria.md, business-rules.md, interface-contracts.md, test-cases.md | EH-*, AU-*, VAL-*, API-* |
| UI/UX | ui-patterns.md, ui-interaction-stories.md, ui-acceptance-criteria.md, component-specs.md, state-machines.md | CC-*, UI-COMP-*, UI-STATE-*, UI-E2E-* |

### Cross-Chain Benefits

1. **Shared decisions** - Authorization model (AU-1) used by both backend AC and UI components
2. **Consistency** - Same error handling strategy reflected in API contracts and UI error displays
3. **Traceability bridge** - Business Rules (BR-*) link backend to UI; both chains reference them

### Example: Cross-Chain Reuse

```python
# UI derivation finds backend decision
rag.retrieve_si_decision("AU-1")
# → Found in acceptance-criteria.md: "customer-or-provider"

# Backend derivation finds UI decision
rag.retrieve_si_decision("UI-DS-1")
# → Found in ui-interaction-stories.md: "shadcn"
# → Informs: API should return data compatible with shadcn patterns
```

## Key Insight

The Self-Learning System creates a **virtuous cycle**:

1. Derive documents with SI decisions
2. Index derived documents into RAG
3. Future derivations find and reuse decisions
4. Consistency improves over time
5. Human effort decreases per derivation

**Single Source of Truth**: Decisions live in the documents, not in a separate cache.
