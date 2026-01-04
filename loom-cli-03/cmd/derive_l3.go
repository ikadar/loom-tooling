// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md IC-DRV-003
// See: l2/sequence-design.md SEQ-L3-001
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
	"loom-cli/internal/formatter"
	"loom-cli/internal/generator"
	"loom-cli/prompts"
)

// runDeriveL3 handles the derive-l3 command.
//
// Implements: IC-DRV-003
// Output Files:
//   - test-cases.md
//   - openapi.json
//   - implementation-skeletons.md
//   - feature-tickets.md
//   - service-boundaries.md
//   - event-message-design.md
//   - dependency-graph.md
func runDeriveL3(args []string) int {
	fs := flag.NewFlagSet("derive-l3", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "Directory containing L2 documents (required)")
	l1Dir := fs.String("l1-dir", "", "Directory containing L1 documents (for AC refs)")
	outputDir := fs.String("output-dir", "", "Output directory for L3 documents (required)")
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
		fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Read L2 documents
	l2Input, err := readL2Documents(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Read L1 acceptance criteria if available
	acContent := ""
	if *l1Dir != "" {
		if data, err := os.ReadFile(filepath.Join(*l1Dir, "acceptance-criteria.md")); err == nil {
			acContent = string(data)
		}
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Generate Test Cases
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 1: Generating test cases...")
	}
	testCasesMD, err := generateTestCases(client, l2Input, acContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate test cases: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "test-cases.md", testCasesMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 2: Generate OpenAPI Spec
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 2: Generating OpenAPI spec...")
	}
	openAPIMD, err := generateOpenAPI(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate API spec: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "openapi.json", openAPIMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 3: Generate Implementation Skeletons
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 3: Generating implementation skeletons...")
	}
	skeletonsMD, err := generateSkeletons(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate implementation skeletons: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "implementation-skeletons.md", skeletonsMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 4: Generate Feature Tickets
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 4: Generating feature tickets...")
	}
	ticketsMD, err := generateFeatureTickets(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate feature tickets: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "feature-tickets.md", ticketsMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 5: Generate Service Boundaries
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 5: Generating service boundaries...")
	}
	serviceMD, err := generateServiceBoundaries(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate service boundaries: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "service-boundaries.md", serviceMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 6: Generate Event/Message Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 6: Generating event/message design...")
	}
	eventMD, err := generateEventDesign(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate event design: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "event-message-design.md", eventMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 7: Generate Dependency Graph
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase 7: Generating dependency graph...")
	}
	depMD, err := generateDependencyGraph(client, l2Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate dependency graph: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "dependency-graph.md", depMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Write combined L3 output as JSON
	l3Output := map[string]interface{}{
		"test_cases":              testCasesMD,
		"openapi":                 openAPIMD,
		"implementation_skeletons": skeletonsMD,
		"feature_tickets":          ticketsMD,
		"service_boundaries":       serviceMD,
		"event_design":            eventMD,
		"dependency_graph":        depMD,
	}
	l3JSON, _ := json.MarshalIndent(l3Output, "", "  ")
	if err := writeFile(*outputDir, "l3-output.json", string(l3JSON)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive-l3] Complete. Output written to %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// l2Input holds L2 documents content.
type l2Input struct {
	TechSpecs          string
	InterfaceContracts string
	AggregateDesign    string
	SequenceDesign     string
	DataModel          string
}

// readL2Documents reads L2 documents from directory.
func readL2Documents(dir string) (*l2Input, error) {
	input := &l2Input{}

	// Read tech-specs.md (required)
	data, err := os.ReadFile(filepath.Join(dir, "tech-specs.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to read tech-specs.md: %w", err)
	}
	input.TechSpecs = string(data)

	// Read interface-contracts.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "interface-contracts.md")); err == nil {
		input.InterfaceContracts = string(data)
	}

	// Read aggregate-design.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "aggregate-design.md")); err == nil {
		input.AggregateDesign = string(data)
	}

	// Read sequence-design.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "sequence-design.md")); err == nil {
		input.SequenceDesign = string(data)
	}

	// Read initial-data-model.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "initial-data-model.md")); err == nil {
		input.DataModel = string(data)
	}

	return input, nil
}

// buildL2Context builds context string from L2 input.
func buildL2Context(input *l2Input) string {
	var context string

	if input.TechSpecs != "" {
		context += fmt.Sprintf("## Technical Specifications\n\n%s\n\n", input.TechSpecs)
	}
	if input.InterfaceContracts != "" {
		context += fmt.Sprintf("## Interface Contracts\n\n%s\n\n", input.InterfaceContracts)
	}
	if input.AggregateDesign != "" {
		context += fmt.Sprintf("## Aggregate Design\n\n%s\n\n", input.AggregateDesign)
	}
	if input.SequenceDesign != "" {
		context += fmt.Sprintf("## Sequence Design\n\n%s\n\n", input.SequenceDesign)
	}
	if input.DataModel != "" {
		context += fmt.Sprintf("## Data Model\n\n%s\n\n", input.DataModel)
	}

	return context
}

// generateTestCases generates test cases using chunked generation.
func generateTestCases(client *claude.Client, input *l2Input, acContent string) (string, error) {
	context := buildL2Context(input)
	if acContent != "" {
		context = fmt.Sprintf("## Acceptance Criteria\n\n%s\n\n", acContent) + context
	}

	// Use chunked generator for large AC sets
	gen := generator.NewChunkedTestCaseGenerator(client)
	result, err := gen.Generate(context)
	if err != nil {
		return "", err
	}

	return formatter.FormatTestCases(result.TestSuites, result.Summary), nil
}

// generateOpenAPI generates OpenAPI specification.
func generateOpenAPI(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveL3API, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateSkeletons generates implementation skeletons.
func generateSkeletons(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveL3Skeletons, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateFeatureTickets generates feature tickets.
func generateFeatureTickets(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveFeatureTickets, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateServiceBoundaries generates service boundaries.
func generateServiceBoundaries(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveServiceBoundaries, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateEventDesign generates event/message design.
func generateEventDesign(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveEventDesign, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateDependencyGraph generates dependency graph.
func generateDependencyGraph(client *claude.Client, input *l2Input) (string, error) {
	context := buildL2Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveDependencyGraph, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
