// Implements: l2/interface-contracts.md IC-DRV-003
// See: l2/sequence-design.md SEQ-L3-001
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
)

// runDeriveL3 implements the derive-l3 command.
//
// Implements: IC-DRV-003
// Output files:
//   - test-cases.md
//   - openapi.json
//   - implementation-skeletons.md
//   - feature-tickets.md
//   - service-boundaries.md
//   - event-message-design.md
//   - dependency-graph.md
func runDeriveL3(args []string) int {
	fs := flag.NewFlagSet("derive-l3", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "L2 docs directory (required)")
	l1Dir := fs.String("l1-dir", "", "L1 docs directory (for AC refs)")
	outputDir := fs.String("output-dir", "", "Output directory (required)")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate required options
	if *inputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --input-dir is required (directory containing L2 documents)")
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

	// Read L2 documents
	techSpecs, err := os.ReadFile(filepath.Join(*inputDir, "tech-specs.md"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read tech-specs.md: %v\n", err)
		return domain.ExitCodeError
	}

	interfaceContracts, _ := os.ReadFile(filepath.Join(*inputDir, "interface-contracts.md"))
	aggregateDesign, _ := os.ReadFile(filepath.Join(*inputDir, "aggregate-design.md"))
	sequenceDesign, _ := os.ReadFile(filepath.Join(*inputDir, "sequence-design.md"))
	dataModel, _ := os.ReadFile(filepath.Join(*inputDir, "initial-data-model.md"))

	// Read L1 documents if provided
	var acceptanceCriteria, businessRules []byte
	if *l1Dir != "" {
		acceptanceCriteria, _ = os.ReadFile(filepath.Join(*l1Dir, "acceptance-criteria.md"))
		businessRules, _ = os.ReadFile(filepath.Join(*l1Dir, "business-rules.md"))
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	l2Context := fmt.Sprintf(`Tech Specs:
%s

Interface Contracts:
%s

Aggregate Design:
%s

Sequence Design:
%s

Data Model:
%s

Acceptance Criteria:
%s

Business Rules:
%s`, string(techSpecs), string(interfaceContracts), string(aggregateDesign),
		string(sequenceDesign), string(dataModel), string(acceptanceCriteria), string(businessRules))

	// Generate L3 documents
	type l3Doc struct {
		name    string
		content string
	}
	docs := make([]l3Doc, 0, 7)

	// Test Cases
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating test-cases.md...")
	}
	testCases, err := deriveTestCases(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate test-cases.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"test-cases.md", testCases})

	// OpenAPI
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating openapi.json...")
	}
	openAPI, err := deriveOpenAPI(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate openapi.json: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"openapi.json", openAPI})

	// Implementation Skeletons
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating implementation-skeletons.md...")
	}
	skeletons, err := deriveSkeletons(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate implementation-skeletons.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"implementation-skeletons.md", skeletons})

	// Feature Tickets
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating feature-tickets.md...")
	}
	tickets, err := deriveFeatureTickets(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate feature-tickets.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"feature-tickets.md", tickets})

	// Service Boundaries
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating service-boundaries.md...")
	}
	services, err := deriveServiceBoundaries(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate service-boundaries.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"service-boundaries.md", services})

	// Event Message Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating event-message-design.md...")
	}
	events, err := deriveEventDesign(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate event-message-design.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"event-message-design.md", events})

	// Dependency Graph
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Generating dependency-graph.md...")
	}
	deps, err := deriveDependencyGraph(client, l2Context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate dependency-graph.md: %v\n", err)
		return domain.ExitCodeError
	}
	docs = append(docs, l3Doc{"dependency-graph.md", deps})

	// Write all files
	for _, doc := range docs {
		path := filepath.Join(*outputDir, doc.name)
		if err := os.WriteFile(path, []byte(doc.content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write %s: %v\n", doc.name, err)
			return domain.ExitCodeError
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "[derive-l3] Written: %s\n", path)
		}
	}

	fmt.Fprintf(os.Stderr, "L3 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

func deriveTestCases(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a QA architect. Generate a test-cases.md document from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Test Case Overview with TDAI summary
3. Test Cases with ID pattern: TC-AC-{CTX}-{NNN}-{T}{NN} where T is P(positive), N(negative), B(boundary), H(hallucination)
4. Each test should have: preconditions, test data, steps, expected results
5. Ensure at least 20%% negative tests
6. Include hallucination prevention tests (what system should NOT do)
7. Traceability to AC and BR

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}

func deriveOpenAPI(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are an API architect. Generate an OpenAPI 3.0 specification from L2 interface contracts.

<context>
%s
</context>

Generate a valid OpenAPI 3.0 JSON specification with:
1. info section with title, version, description
2. servers section
3. paths derived from interface contracts
4. components/schemas derived from data model
5. security schemes if applicable
6. Error response schemas

Return ONLY valid JSON, no markdown.`, l2Context)

	return client.Call(prompt)
}

func deriveSkeletons(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a software architect. Generate implementation-skeletons.md from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Implementation Skeletons with ID pattern: SKEL-{CTX}-{NNN}
3. Each skeleton should have: purpose, module structure, key interfaces, implementation notes
4. Code structure outlines (not full implementations)
5. Traceability to aggregates and tech specs

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}

func deriveFeatureTickets(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a product manager. Generate feature-tickets.md from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Feature Tickets with ID pattern: FDT-{NNN}
3. Each ticket should have: title, description, acceptance criteria, story points estimate, dependencies
4. Prioritized order
5. Traceability to AC

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}

func deriveServiceBoundaries(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a microservices architect. Generate service-boundaries.md from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Services with ID pattern: SVC-{NAME}
3. Each service should have: responsibilities, API surface, data ownership, dependencies
4. Service interaction diagram (Mermaid)
5. Traceability to bounded contexts

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}

func deriveEventDesign(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are an event-driven architect. Generate event-message-design.md from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Domain Events with ID pattern: EVT-{CTX}-{NNN}
3. Commands with ID pattern: CMD-{CTX}-{NNN}
4. Integration Events with ID pattern: INT-{CTX}-{NNN}
5. Each event/command should have: payload schema, producer, consumers
6. Event flow diagrams (Mermaid)
7. Traceability to sequences and aggregates

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}

func deriveDependencyGraph(client *claude.Client, l2Context string) (string, error) {
	prompt := fmt.Sprintf(`You are a software architect. Generate dependency-graph.md from L2 documents.

<context>
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L3)
2. Dependency entries with ID pattern: DEP-{CTX}-{NNN}
3. Module/service dependencies
4. Data dependencies
5. External service dependencies
6. Mermaid dependency graph diagram
7. Build/deployment order

Use proper markdown formatting.`, l2Context)

	return client.Call(prompt)
}
