"""
RAG PoC Test - Compare RAG vs no-RAG document derivation.

This script tests the quality improvement from using RAG
with the guidelines as knowledge base.
"""

from pathlib import Path
from rag_engine import LoomRAG, get_default_guidelines_path


# Test input: Domain Vocabulary for Quote
TEST_INPUT = """
## Domain Vocabulary - Quote Management

### Core Concepts

- **Quote**: A formal price offer presented to a customer, containing one or more line items with pricing
- **Quote Line Item**: An individual product or service entry within a quote, including quantity and unit price
- **Quote Validity Period**: The time window during which the quote's pricing and terms remain valid
- **Quote Status**: The current state of the quote in its lifecycle (draft, sent, accepted, rejected, expired)
- **Quote Total**: The calculated sum of all line item prices, potentially including taxes and discounts
- **Customer Reference**: A link to the customer entity who requested or will receive the quote

### Business Context

Quotes are created by sales representatives in response to customer inquiries.
A quote may go through several revisions before being sent to the customer.
Once sent, the customer can accept or reject the quote within the validity period.
Accepted quotes typically convert into orders.
"""

# Derivation task
DERIVATION_TASK = """
Derive a Domain Model from this vocabulary. The domain model should include:
1. Aggregate definition with root entity
2. Classification of each concept as Entity or Value Object with rationale
3. Invariants that must be maintained
4. State transitions and behaviors
5. Aggregate boundaries and external interactions
"""


def run_comparison():
    """Run RAG vs no-RAG comparison test."""

    print("=" * 60)
    print("RAG PoC - Document Derivation Comparison")
    print("=" * 60)

    # Initialize RAG engine
    guidelines_path = get_default_guidelines_path()
    print(f"\nGuidelines path: {guidelines_path}")

    rag = LoomRAG(
        guidelines_dir=guidelines_path,
        persist_dir=str(Path(__file__).parent / "chroma_db"),
    )

    # Show what RAG retrieves for context
    print("\n" + "=" * 60)
    print("RAG Context Preview")
    print("=" * 60)

    retrieved = rag.retrieve_with_sources(DERIVATION_TASK, k=3)
    for i, item in enumerate(retrieved, 1):
        source = Path(item["source"]).name
        content_preview = item["content"][:200] + "..."
        print(f"\n[{i}] Source: {source}")
        print(f"    Content: {content_preview}")

    # Derive WITHOUT RAG
    print("\n" + "=" * 60)
    print("Derivation WITHOUT RAG")
    print("=" * 60)

    result_no_rag = rag.derive_without_rag(TEST_INPUT, DERIVATION_TASK)
    print(result_no_rag)

    # Save to file
    output_dir = Path(__file__).parent / "output"
    output_dir.mkdir(exist_ok=True)

    (output_dir / "derivation_no_rag.md").write_text(result_no_rag)
    print(f"\nSaved to: output/derivation_no_rag.md")

    # Derive WITH RAG
    print("\n" + "=" * 60)
    print("Derivation WITH RAG")
    print("=" * 60)

    result_with_rag = rag.derive_with_rag(TEST_INPUT, DERIVATION_TASK, k=5)
    print(result_with_rag)

    # Save to file
    (output_dir / "derivation_with_rag.md").write_text(result_with_rag)
    print(f"\nSaved to: output/derivation_with_rag.md")

    # Summary
    print("\n" + "=" * 60)
    print("Summary")
    print("=" * 60)
    print(f"\nOutput files created in: {output_dir}")
    print("- derivation_no_rag.md  (baseline)")
    print("- derivation_with_rag.md (RAG-enhanced)")
    print("\nCompare these files to evaluate RAG quality improvement.")


def test_retrieval_only():
    """Test just the retrieval part."""

    print("=" * 60)
    print("RAG Retrieval Test")
    print("=" * 60)

    guidelines_path = get_default_guidelines_path()
    rag = LoomRAG(
        guidelines_dir=guidelines_path,
        persist_dir=str(Path(__file__).parent / "chroma_db"),
    )

    queries = [
        "What must an aggregate definition contain?",
        "How to identify entity vs value object?",
        "What are aggregate invariants?",
        "Service boundary design rules",
    ]

    for query in queries:
        print(f"\n{'─' * 40}")
        print(f"Query: {query}")
        print("─" * 40)

        results = rag.retrieve_with_sources(query, k=2)
        for i, item in enumerate(results, 1):
            source = Path(item["source"]).name
            print(f"\n[{i}] {source}")
            print(f"    {item['content'][:150]}...")


if __name__ == "__main__":
    import sys

    if len(sys.argv) > 1 and sys.argv[1] == "--retrieval-only":
        test_retrieval_only()
    else:
        run_comparison()
