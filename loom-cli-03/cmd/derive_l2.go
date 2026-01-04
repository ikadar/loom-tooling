// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md IC-DRV-002
// See: l2/sequence-design.md SEQ-L2-001
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
	"loom-cli/internal/formatter"
	"loom-cli/prompts"
)

// runDeriveL2 handles the derive-l2 command.
//
// Implements: IC-DRV-002
// Output Files:
//   - tech-specs.md
//   - interface-contracts.md
//   - aggregate-design.md
//   - sequence-design.md
//   - initial-data-model.md
func runDeriveL2(args []string) int {
	fs := flag.NewFlagSet("derive-l2", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "Directory containing L1 documents (required)")
	outputDir := fs.String("output-dir", "", "Output directory for L2 documents (required)")
	interactive := fs.Bool("interactive", false, "Interactive approval mode")
	fs.BoolVar(interactive, "i", false, "Interactive approval mode (shorthand)")
	verbose := fs.Bool("verbose", false, "Verbose output")

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
		fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Read L1 documents
	l1Input, err := readL1Documents(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Generate Tech Specs
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase 1: Generating tech specs...")
	}
	techSpecsMD, err := generateTechSpecs(client, l1Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate tech-specs: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		techSpecsMD, err = interactiveApproval("tech-specs.md", techSpecsMD, client, l1Input, generateTechSpecs)
		if err != nil {
			return domain.ExitCodeError
		}
	}
	if err := writeFile(*outputDir, "tech-specs.md", techSpecsMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 2: Generate Interface Contracts
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase 2: Generating interface contracts...")
	}
	contractsMD, err := generateInterfaceContracts(client, l1Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate interface-contracts: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		contractsMD, err = interactiveApproval("interface-contracts.md", contractsMD, client, l1Input, generateInterfaceContracts)
		if err != nil {
			return domain.ExitCodeError
		}
	}
	if err := writeFile(*outputDir, "interface-contracts.md", contractsMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 3: Generate Aggregate Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase 3: Generating aggregate design...")
	}
	aggregateMD, err := generateAggregateDesign(client, l1Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate aggregate-design: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		aggregateMD, err = interactiveApproval("aggregate-design.md", aggregateMD, client, l1Input, generateAggregateDesign)
		if err != nil {
			return domain.ExitCodeError
		}
	}
	if err := writeFile(*outputDir, "aggregate-design.md", aggregateMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 4: Generate Sequence Design
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase 4: Generating sequence design...")
	}
	sequenceMD, err := generateSequenceDesign(client, l1Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate sequence-design: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		sequenceMD, err = interactiveApproval("sequence-design.md", sequenceMD, client, l1Input, generateSequenceDesign)
		if err != nil {
			return domain.ExitCodeError
		}
	}
	if err := writeFile(*outputDir, "sequence-design.md", sequenceMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 5: Generate Data Model
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase 5: Generating initial data model...")
	}
	dataModelMD, err := generateDataModel(client, l1Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate initial-data-model: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		dataModelMD, err = interactiveApproval("initial-data-model.md", dataModelMD, client, l1Input, generateDataModel)
		if err != nil {
			return domain.ExitCodeError
		}
	}
	if err := writeFile(*outputDir, "initial-data-model.md", dataModelMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive-l2] Complete. Output written to %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// l1Input holds L1 documents content.
type l1Input struct {
	DomainModel        string
	BoundedContextMap  string
	AcceptanceCriteria string
	BusinessRules      string
	Decisions          string
}

// readL1Documents reads L1 documents from directory.
func readL1Documents(dir string) (*l1Input, error) {
	input := &l1Input{}

	// Read domain-model.md (optional but recommended)
	if data, err := os.ReadFile(filepath.Join(dir, "domain-model.md")); err == nil {
		input.DomainModel = string(data)
	}

	// Read bounded-context-map.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "bounded-context-map.md")); err == nil {
		input.BoundedContextMap = string(data)
	}

	// Read acceptance-criteria.md (required)
	data, err := os.ReadFile(filepath.Join(dir, "acceptance-criteria.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to read acceptance-criteria.md: %w", err)
	}
	input.AcceptanceCriteria = string(data)

	// Read business-rules.md (required)
	data, err = os.ReadFile(filepath.Join(dir, "business-rules.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to read business-rules.md: %w", err)
	}
	input.BusinessRules = string(data)

	// Read decisions.md (optional)
	if data, err := os.ReadFile(filepath.Join(dir, "decisions.md")); err == nil {
		input.Decisions = string(data)
	}

	return input, nil
}

// generateTechSpecs generates tech specs from L1 input.
func generateTechSpecs(client *claude.Client, input *l1Input) (string, error) {
	context := buildL1Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveTechSpecs, context)

	var specs []formatter.TechSpec
	if err := client.CallJSONWithRetry(prompt, &specs, claude.DefaultRetryConfig()); err != nil {
		return "", err
	}

	return formatter.FormatTechSpecs(specs), nil
}

// generateInterfaceContracts generates interface contracts from L1 input.
func generateInterfaceContracts(client *claude.Client, input *l1Input) (string, error) {
	context := buildL1Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveInterfaceContracts, context)

	var contracts []formatter.InterfaceContract
	if err := client.CallJSONWithRetry(prompt, &contracts, claude.DefaultRetryConfig()); err != nil {
		return "", err
	}

	return formatter.FormatInterfaceContracts(contracts, nil), nil
}

// generateAggregateDesign generates aggregate design from L1 input.
func generateAggregateDesign(client *claude.Client, input *l1Input) (string, error) {
	context := buildL1Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveAggregateDesign, context)

	var aggregates []formatter.AggregateDesign
	if err := client.CallJSONWithRetry(prompt, &aggregates, claude.DefaultRetryConfig()); err != nil {
		return "", err
	}

	return formatter.FormatAggregateDesign(aggregates), nil
}

// generateSequenceDesign generates sequence design from L1 input.
func generateSequenceDesign(client *claude.Client, input *l1Input) (string, error) {
	context := buildL1Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveSequenceDesign, context)

	var sequences []formatter.SequenceDesign
	if err := client.CallJSONWithRetry(prompt, &sequences, claude.DefaultRetryConfig()); err != nil {
		return "", err
	}

	return formatter.FormatSequenceDesign(sequences), nil
}

// generateDataModel generates initial data model from L1 input.
func generateDataModel(client *claude.Client, input *l1Input) (string, error) {
	context := buildL1Context(input)
	prompt := claude.BuildPrompt(prompts.DeriveDataModel, context)

	var result struct {
		Tables []formatter.DataTable `json:"tables"`
		Enums  []formatter.DataEnum  `json:"enums"`
	}
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return "", err
	}

	return formatter.FormatDataModel(result.Tables, result.Enums), nil
}

// buildL1Context builds context string from L1 input.
func buildL1Context(input *l1Input) string {
	var context string

	if input.DomainModel != "" {
		context += fmt.Sprintf("## Domain Model\n\n%s\n\n", input.DomainModel)
	}
	if input.BoundedContextMap != "" {
		context += fmt.Sprintf("## Bounded Context Map\n\n%s\n\n", input.BoundedContextMap)
	}
	if input.AcceptanceCriteria != "" {
		context += fmt.Sprintf("## Acceptance Criteria\n\n%s\n\n", input.AcceptanceCriteria)
	}
	if input.BusinessRules != "" {
		context += fmt.Sprintf("## Business Rules\n\n%s\n\n", input.BusinessRules)
	}
	if input.Decisions != "" {
		context += fmt.Sprintf("## Decisions\n\n%s\n\n", input.Decisions)
	}

	return context
}

// generateFunc is a function type for regeneration in interactive mode.
type generateFunc func(*claude.Client, *l1Input) (string, error)

// interactiveApproval handles interactive approval flow for a document.
func interactiveApproval(filename, content string, client *claude.Client, input *l1Input, regenerate generateFunc) (string, error) {
	for {
		showPreview(filename, content)
		action := requestApproval()

		switch action {
		case domain.ActionApprove:
			return content, nil

		case domain.ActionEdit:
			edited, err := openInEditor(content)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			content = edited

		case domain.ActionRegenerate:
			fmt.Fprintln(os.Stderr, "Regenerating...")
			newContent, err := regenerate(client, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: regeneration failed: %v\n", err)
				continue
			}
			content = newContent

		case domain.ActionSkip:
			return content, nil

		case domain.ActionQuit:
			if confirmQuit() {
				os.Exit(0)
			}
		}
	}
}
