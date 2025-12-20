"""
Loom RAG Engine - PoC

Uses the guidelines from 9300-guidelines/ as knowledge base
for enhanced document derivation.
"""

import os
from pathlib import Path

from langchain_community.document_loaders import DirectoryLoader, TextLoader
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_chroma import Chroma
from anthropic import Anthropic


class LoomRAG:
    """RAG engine for Loom document derivation."""

    def __init__(
        self,
        guidelines_dir: str,
        persist_dir: str = "./chroma_db",
        chunk_size: int = 500,
        chunk_overlap: int = 50,
    ):
        """
        Initialize the RAG engine.

        Args:
            guidelines_dir: Path to the guidelines directory
            persist_dir: Path to persist the vector database
            chunk_size: Size of each text chunk
            chunk_overlap: Overlap between chunks
        """
        self.guidelines_dir = Path(guidelines_dir)
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
        """Load guidelines and create vector database."""

        # Check if we already have a persisted database
        if Path(self.persist_dir).exists():
            print(f"Loading existing vector database from {self.persist_dir}")
            self.vectordb = Chroma(
                persist_directory=self.persist_dir,
                embedding_function=self.embeddings,
            )
            return

        print(f"Creating new vector database from {self.guidelines_dir}")

        # Load all markdown files from guidelines directory
        loader = DirectoryLoader(
            str(self.guidelines_dir),
            glob="**/*.md",
            loader_cls=TextLoader,
            loader_kwargs={"encoding": "utf-8"},
        )
        documents = loader.load()

        print(f"Loaded {len(documents)} documents")

        # Split documents into chunks
        splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size,
            chunk_overlap=self.chunk_overlap,
            separators=["\n## ", "\n### ", "\n\n", "\n", " "],
        )
        chunks = splitter.split_documents(documents)

        print(f"Created {len(chunks)} chunks")

        # Create vector database
        self.vectordb = Chroma.from_documents(
            documents=chunks,
            embedding=self.embeddings,
            persist_directory=self.persist_dir,
        )

        print(f"Vector database created and persisted to {self.persist_dir}")

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
            List of dicts with content and source
        """
        docs = self.vectordb.similarity_search(query, k=k)
        return [
            {
                "content": doc.page_content,
                "source": doc.metadata.get("source", "unknown"),
            }
            for doc in docs
        ]

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
        # Retrieve relevant guidelines
        context_chunks = self.retrieve(derivation_task, k=k)
        context = "\n\n---\n\n".join(context_chunks)

        # Build prompt with RAG context
        prompt = f"""You are a Loom documentation derivation expert. Use the following guidelines
from the knowledge base to inform your derivation.

## Guidelines (from knowledge base):

{context}

---

## Input Document:

{input_document}

---

## Derivation Task:

{derivation_task}

---

## Instructions:

1. Follow the guidelines strictly when deriving the output
2. Use the structure and format specified in the guidelines
3. Include all required sections mentioned in the guidelines
4. Provide rationale based on the guidelines where appropriate

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


if __name__ == "__main__":
    # Quick test
    guidelines_path = get_default_guidelines_path()
    print(f"Guidelines path: {guidelines_path}")

    rag = LoomRAG(guidelines_path)

    # Test retrieval
    query = "What must an aggregate definition contain?"
    results = rag.retrieve(query, k=3)

    print(f"\nQuery: {query}")
    print(f"\nRetrieved {len(results)} chunks:")
    for i, chunk in enumerate(results, 1):
        print(f"\n--- Chunk {i} ---")
        print(chunk[:300] + "..." if len(chunk) > 300 else chunk)
