// Implements: l2/interface-contracts.md IC-DRV-002
// See: l2/sequence-design.md SEQ-DRV-001, SEQ-L2-001
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
)

// runDeriveL2 implements the derive-l2 command.
//
// Implements: IC-DRV-002
// Output files:
//   - tech-specs.md
//   - interface-contracts.md
//   - aggregate-design.md
//   - sequence-design.md
//   - initial-data-model.md
func runDeriveL2(args []string) int {
	fs := flag.NewFlagSet("derive-l2", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "L1 docs directory (required)")
	outputDir := fs.String("output-dir", "", "Output directory (required)")
	verbose := fs.Bool("verbose", false, "Verbose output")
	interactive := fs.Bool("interactive", false, "Interactive approval mode")
	fs.Bool("i", false, "Alias for --interactive")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate required options
	if *inputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --input-dir is required (directory containing L1 documents)")
		return domain.ExitCodeError
	}
	if *outputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --output-dir is required")
		return domain.ExitCodeError
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create output directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Read L1 documents
	domainModel, err := os.ReadFile(filepath.Join(*inputDir, "domain-model.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read domain-model.md: %v\n", err)
		return domain.ExitCodeError
	}

	acceptanceCriteria, err := os.ReadFile(filepath.Join(*inputDir, "acceptance-criteria.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read acceptance-criteria.md: %v\n", err)
		return domain.ExitCodeError
	}

	businessRules, err := os.ReadFile(filepath.Join(*inputDir, "business-rules.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read business-rules.md: %v\n", err)
		return domain.ExitCodeError
	}

	boundedContext, _ := os.ReadFile(filepath.Join(*inputDir, "bounded-context-map.md"))

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	l1Context := fmt.Sprintf(`Domain Model:
%s

Acceptance Criteria:
%s

Business Rules:
%s

Bounded Context Map:
%s`, string(domainModel), string(acceptanceCriteria), string(businessRules), string(boundedContext))

	// Generate L2 documents (could be parallelized in future)
	type l2Doc struct {
		name    string
		content string
	}
	docs := make([]l2Doc, 0, 5)

	// Tech Specs
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Generating tech-specs.md...")
	}
	techSpecs, err := deriveTechSpecs(client, l1Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate tech-specs.md: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		techSpecs, _ = interactiveApproval("tech-specs.md", techSpecs, client, func() (string, error) {
			return deriveTechSpecs(client, l1Context)
		})
	}
	docs = append(docs, l2Doc{"tech-specs.md", techSpecs})

	// Interface Contracts
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Generating interface-contracts.md...")
	}
	contracts, err := deriveInterfaceContracts(client, l1Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate interface-contracts.md: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		contracts, _ = interactiveApproval("interface-contracts.md", contracts, client, func() (string, error) {
			return deriveInterfaceContracts(client, l1Context)
		})
	}
	docs = append(docs, l2Doc{"interface-contracts.md", contracts})

	// Aggregate Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Generating aggregate-design.md...")
	}
	aggregates, err := deriveAggregateDesign(client, l1Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate aggregate-design.md: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		aggregates, _ = interactiveApproval("aggregate-design.md", aggregates, client, func() (string, error) {
			return deriveAggregateDesign(client, l1Context)
		})
	}
	docs = append(docs, l2Doc{"aggregate-design.md", aggregates})

	// Sequence Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Generating sequence-design.md...")
	}
	sequences, err := deriveSequenceDesign(client, l1Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate sequence-design.md: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		sequences, _ = interactiveApproval("sequence-design.md", sequences, client, func() (string, error) {
			return deriveSequenceDesign(client, l1Context)
		})
	}
	docs = append(docs, l2Doc{"sequence-design.md", sequences})

	// Data Model
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Generating initial-data-model.md...")
	}
	dataModel, err := deriveDataModel(client, l1Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate initial-data-model.md: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		dataModel, _ = interactiveApproval("initial-data-model.md", dataModel, client, func() (string, error) {
			return deriveDataModel(client, l1Context)
		})
	}
	docs = append(docs, l2Doc{"initial-data-model.md", dataModel})

	// Write all files
	for _, doc := range docs {
		path := filepath.Join(*outputDir, doc.name)
		if err := os.WriteFile(path, []byte(doc.content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write %s: %v\n", doc.name, err)
			return domain.ExitCodeError
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "[derive-l2] Written: %s\n", path)
		}
	}

	fmt.Fprintf(os.Stderr, "L2 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

func deriveTechSpecs(client *claude.Client, l1Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a technical architect. Generate a tech-specs.md document from L1 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L2)
2. Technology Stack section
3. Technical Specifications with ID pattern: TS-{CTX}-{NNN}
4. Each spec should have: implementation details, validation points, error handling
5. Traceability to L1 documents (AC, BR references)

Use proper markdown formatting.`, l1Context)

	return client.Call(prompt)
}

func deriveInterfaceContracts(client *claude.Client, l1Context string) (string, error) {
	prompt := fmt.Sprintf(`You are an API architect. Generate an interface-contracts.md document from L1 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L2)
2. API Overview section
3. Interface Contracts with ID pattern: IC-{CTX}-{NNN}
4. Each contract should have: endpoint, method, request/response schemas, error codes
5. Traceability to L1 documents (AC references)

Use proper markdown formatting.`, l1Context)

	return client.Call(prompt)
}

func deriveAggregateDesign(client *claude.Client, l1Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a DDD expert. Generate an aggregate-design.md document from L1 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L2)
2. Aggregate Overview section
3. Aggregates with ID pattern: AGG-{CTX}-{NNN}
4. Each aggregate should have: root entity, invariants, entities, value objects, behaviors, events
5. Traceability to L1 domain model

Use proper markdown formatting.`, l1Context)

	return client.Call(prompt)
}

func deriveSequenceDesign(client *claude.Client, l1Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a software architect. Generate a sequence-design.md document from L1 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L2)
2. Sequence Overview section
3. Sequences with ID pattern: SEQ-{CTX}-{NNN}
4. Each sequence should have: trigger, participants, steps, outcome
5. Mermaid sequence diagrams for each flow
6. Traceability to L1 acceptance criteria

Use proper markdown formatting.`, l1Context)

	return client.Call(prompt)
}

func deriveDataModel(client *claude.Client, l1Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a data architect. Generate an initial-data-model.md document from L1 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L2)
2. Data Model Overview section
3. Tables/Collections derived from aggregates
4. Each table should have: fields, types, constraints, indexes, foreign keys
5. Enums section
6. Traceability to L1 domain model and aggregates

Use proper markdown formatting.`, l1Context)

	return client.Call(prompt)
}
