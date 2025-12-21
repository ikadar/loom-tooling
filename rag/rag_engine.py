"""
Loom RAG Engine - Self-Learning System

Uses multiple knowledge sources for enhanced document derivation:
1. Guidelines (global) - ai-pds-specification/9300-guidelines
2. Project docs (project-specific) - derived documents with SI decisions
3. Code patterns (optional) - src/ directory for implementation patterns

The system learns from its own output: derived documents become part of
the knowledge base, enabling decision reuse across derivations.
"""

import os
from pathlib import Path
from typing import Optional

from langchain_community.document_loaders import DirectoryLoader, TextLoader
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_chroma import Chroma
from anthropic import Anthropic


class KnowledgeSource:
    """Represents a knowledge source with type and priority."""

    GUIDELINES = "guidelines"  # Global guidelines (lowest priority for conflicts)
    PROJECT = "project"        # Project documentation (high priority)
    CODE = "code"              # Source code patterns (reference only)

    def __init__(self, path: str, source_type: str, priority: int = 1):
        self.path = Path(path)
        self.source_type = source_type
        self.priority = priority  # Higher = more important in conflicts


class LoomRAG:
    """
    RAG engine for Loom document derivation.

    Supports the Self-Learning System pattern where:
    - Guidelines provide base knowledge
    - Project documents (with SI decisions) are indexed
    - Previous decisions inform future derivations
    """

    def __init__(
        self,
        knowledge_sources: Optional[list[KnowledgeSource]] = None,
        guidelines_dir: Optional[str] = None,  # Backward compatible
        persist_dir: str = "./chroma_db",
        chunk_size: int = 500,
        chunk_overlap: int = 50,
    ):
        """
        Initialize the RAG engine.

        Args:
            knowledge_sources: List of KnowledgeSource objects
            guidelines_dir: (Deprecated) Path to guidelines directory
            persist_dir: Path to persist the vector database
            chunk_size: Size of each text chunk
            chunk_overlap: Overlap between chunks
        """
        # Backward compatibility: convert guidelines_dir to KnowledgeSource
        if knowledge_sources is None and guidelines_dir is not None:
            knowledge_sources = [
                KnowledgeSource(guidelines_dir, KnowledgeSource.GUIDELINES, priority=1)
            ]

        self.knowledge_sources = knowledge_sources or []
        self.persist_dir = persist_dir
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

        self.vectordb = None
        self.claude = Anthropic()

        # Use local HuggingFace embeddings (free, no API key needed)
        self.embeddings = HuggingFaceEmbeddings(
            model_name="all-MiniLM-L6-v2",
            model_kwargs={"device": "cpu"},
        )

        self._initialize_vectordb()

    def _initialize_vectordb(self):
        """Load all knowledge sources and create vector database."""

        # Check if we already have a persisted database
        if Path(self.persist_dir).exists():
            print(f"Loading existing vector database from {self.persist_dir}")
            self.vectordb = Chroma(
                persist_directory=self.persist_dir,
                embedding_function=self.embeddings,
            )
            return

        all_chunks = []

        for source in self.knowledge_sources:
            if not source.path.exists():
                print(f"Warning: Knowledge source not found: {source.path}")
                continue

            print(f"Loading {source.source_type} from {source.path}")

            # Load all markdown files from this source
            loader = DirectoryLoader(
                str(source.path),
                glob="**/*.md",
                loader_cls=TextLoader,
                loader_kwargs={"encoding": "utf-8"},
            )
            documents = loader.load()

            # Add source metadata to each document
            for doc in documents:
                doc.metadata["source_type"] = source.source_type
                doc.metadata["priority"] = source.priority

            print(f"  Loaded {len(documents)} documents")

            # Split documents into chunks
            splitter = RecursiveCharacterTextSplitter(
                chunk_size=self.chunk_size,
                chunk_overlap=self.chunk_overlap,
                separators=["\n## ", "\n### ", "\n\n", "\n", " "],
            )
            chunks = splitter.split_documents(documents)

            # Preserve metadata in chunks
            for chunk in chunks:
                chunk.metadata["source_type"] = source.source_type
                chunk.metadata["priority"] = source.priority

            all_chunks.extend(chunks)
            print(f"  Created {len(chunks)} chunks")

        if not all_chunks:
            print("Warning: No documents found in any knowledge source")
            self.vectordb = Chroma(
                embedding_function=self.embeddings,
                persist_directory=self.persist_dir,
            )
            return

        print(f"Total: {len(all_chunks)} chunks from {len(self.knowledge_sources)} sources")

        # Create vector database
        self.vectordb = Chroma.from_documents(
            documents=all_chunks,
            embedding=self.embeddings,
            persist_directory=self.persist_dir,
        )

        print(f"Vector database created and persisted to {self.persist_dir}")

    def add_knowledge_source(self, source: KnowledgeSource, rebuild: bool = True):
        """
        Add a new knowledge source to the RAG engine.

        Args:
            source: The KnowledgeSource to add
            rebuild: If True, rebuild the vector database
        """
        self.knowledge_sources.append(source)

        if rebuild:
            # Remove existing database and rebuild
            import shutil
            if Path(self.persist_dir).exists():
                shutil.rmtree(self.persist_dir)
            self._initialize_vectordb()

    def refresh_project_knowledge(self, project_dir: str):
        """
        Refresh the knowledge base with updated project documents.

        Call this after derivation to include new documents in the knowledge base.

        Args:
            project_dir: Path to the project directory
        """
        # Remove old project source if exists
        self.knowledge_sources = [
            s for s in self.knowledge_sources
            if s.source_type != KnowledgeSource.PROJECT
        ]

        # Add updated project source
        self.add_knowledge_source(
            KnowledgeSource(project_dir, KnowledgeSource.PROJECT, priority=2),
            rebuild=True
        )

    def retrieve(self, query: str, k: int = 5) -> list[str]:
        """
        Retrieve relevant chunks for a query.

        Args:
            query: The search query
            k: Number of chunks to retrieve

        Returns:
            List of relevant text chunks
        """
        docs = self.vectordb.similarity_search(query, k=k)
        return [doc.page_content for doc in docs]

    def retrieve_with_sources(self, query: str, k: int = 5) -> list[dict]:
        """
        Retrieve relevant chunks with source information.

        Args:
            query: The search query
            k: Number of chunks to retrieve

        Returns:
            List of dicts with content, source, and source_type
        """
        docs = self.vectordb.similarity_search(query, k=k)
        return [
            {
                "content": doc.page_content,
                "source": doc.metadata.get("source", "unknown"),
                "source_type": doc.metadata.get("source_type", "unknown"),
                "priority": doc.metadata.get("priority", 1),
            }
            for doc in docs
        ]

    def retrieve_si_decision(self, decision_id: str, domain: str = "") -> Optional[dict]:
        """
        Retrieve a previous SI decision from the knowledge base.

        Args:
            decision_id: The decision point ID (e.g., "EH-1", "API-1")
            domain: Optional domain filter (e.g., "booking-system")

        Returns:
            Dict with decision info if found, None otherwise
        """
        query = f"structured-interview decision {decision_id} {domain}".strip()
        results = self.retrieve_with_sources(query, k=3)

        # Look for results that contain the decision ID
        for result in results:
            if decision_id in result["content"]:
                return {
                    "decision_id": decision_id,
                    "content": result["content"],
                    "source": result["source"],
                    "source_type": result["source_type"],
                }

        return None

    def retrieve_prioritized(
        self,
        query: str,
        k: int = 5,
        prefer_project: bool = True
    ) -> list[dict]:
        """
        Retrieve chunks with priority-based ordering.

        Project documents are prioritized over guidelines when prefer_project=True.

        Args:
            query: The search query
            k: Number of chunks to retrieve
            prefer_project: If True, prioritize project docs over guidelines

        Returns:
            List of dicts sorted by priority
        """
        # Retrieve more than k to allow for filtering
        docs = self.vectordb.similarity_search(query, k=k * 2)

        results = [
            {
                "content": doc.page_content,
                "source": doc.metadata.get("source", "unknown"),
                "source_type": doc.metadata.get("source_type", "unknown"),
                "priority": doc.metadata.get("priority", 1),
            }
            for doc in docs
        ]

        if prefer_project:
            # Sort by priority (higher first), then by original order
            results.sort(key=lambda x: -x["priority"])

        return results[:k]

    def derive_with_rag(
        self,
        input_document: str,
        derivation_task: str,
        k: int = 5,
    ) -> str:
        """
        Derive a document using RAG-enhanced context.

        Args:
            input_document: The source document to derive from
            derivation_task: Description of what to derive
            k: Number of context chunks to retrieve

        Returns:
            The derived document
        """
        # Retrieve relevant context, prioritizing project docs
        context_results = self.retrieve_prioritized(derivation_task, k=k)

        # Format context with source info
        context_parts = []
        for result in context_results:
            source_label = f"[{result['source_type'].upper()}]"
            context_parts.append(f"{source_label}\n{result['content']}")

        context = "\n\n---\n\n".join(context_parts)

        # Build prompt with RAG context
        prompt = f"""You are a Loom documentation derivation expert. Use the following context
from the knowledge base to inform your derivation.

Note: [PROJECT] sources contain previous decisions that should be reused.
      [GUIDELINES] sources provide format and structure guidance.

## Knowledge Base Context:

{context}

---

## Input Document:

{input_document}

---

## Derivation Task:

{derivation_task}

---

## Instructions:

1. Check [PROJECT] context for previous SI decisions - reuse them
2. Follow [GUIDELINES] for format and structure
3. Include all required sections
4. Document any new decisions in YAML frontmatter

Generate the derived document:
"""

        # Call Claude API
        response = self.claude.messages.create(
            model="claude-sonnet-4-20250514",
            max_tokens=4096,
            messages=[{"role": "user", "content": prompt}],
        )

        return response.content[0].text

    def derive_without_rag(
        self,
        input_document: str,
        derivation_task: str,
    ) -> str:
        """
        Derive a document WITHOUT RAG (for comparison).

        Args:
            input_document: The source document to derive from
            derivation_task: Description of what to derive

        Returns:
            The derived document
        """
        prompt = f"""You are a documentation derivation expert.

## Input Document:

{input_document}

---

## Derivation Task:

{derivation_task}

---

Generate the derived document:
"""

        # Call Claude API
        response = self.claude.messages.create(
            model="claude-sonnet-4-20250514",
            max_tokens=4096,
            messages=[{"role": "user", "content": prompt}],
        )

        return response.content[0].text


def get_default_guidelines_path() -> str:
    """Get the default path to the guidelines directory."""
    # Relative to this file's location
    current_dir = Path(__file__).parent
    # Go up to specs-for-ai root, then to guidelines
    guidelines_path = current_dir.parent.parent / "ai-pds-specification" / "9000-appendix" / "9300-guidelines"
    return str(guidelines_path)


def create_self_learning_rag(
    guidelines_dir: str,
    project_dir: Optional[str] = None,
    code_dir: Optional[str] = None,
    persist_dir: str = "./chroma_db",
) -> LoomRAG:
    """
    Create a RAG engine configured for the Self-Learning System.

    Args:
        guidelines_dir: Path to the guidelines directory
        project_dir: Optional path to project documentation
        code_dir: Optional path to source code
        persist_dir: Path to persist the vector database

    Returns:
        Configured LoomRAG instance
    """
    sources = [
        KnowledgeSource(guidelines_dir, KnowledgeSource.GUIDELINES, priority=1)
    ]

    if project_dir:
        sources.append(
            KnowledgeSource(project_dir, KnowledgeSource.PROJECT, priority=2)
        )

    if code_dir:
        sources.append(
            KnowledgeSource(code_dir, KnowledgeSource.CODE, priority=1)
        )

    return LoomRAG(knowledge_sources=sources, persist_dir=persist_dir)


if __name__ == "__main__":
    # Quick test with self-learning system
    guidelines_path = get_default_guidelines_path()
    print(f"Guidelines path: {guidelines_path}")

    # Create with just guidelines (backward compatible)
    rag = LoomRAG(guidelines_dir=guidelines_path)

    # Test retrieval
    query = "What must an aggregate definition contain?"
    results = rag.retrieve(query, k=3)

    print(f"\nQuery: {query}")
    print(f"\nRetrieved {len(results)} chunks:")
    for i, chunk in enumerate(results, 1):
        print(f"\n--- Chunk {i} ---")
        print(chunk[:300] + "..." if len(chunk) > 300 else chunk)

    # Test SI decision retrieval
    print("\n\n=== SI Decision Retrieval Test ===")
    decision = rag.retrieve_si_decision("EH-1", "booking")
    if decision:
        print(f"Found decision: {decision['decision_id']}")
        print(f"Source: {decision['source']}")
    else:
        print("No cached decision found (expected if project not indexed)")
