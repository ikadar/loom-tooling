// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md IC-DRV-001
// See: l2/sequence-design.md SEQ-DRV-001
// See: l2/tech-specs.md TS-ARCH-001b
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
	"loom-cli/prompts"
)

// runDerive handles the derive command (L1 derivation).
//
// Implements: IC-DRV-001
// Output Files:
//   - domain-model.md
//   - bounded-context-map.md
//   - acceptance-criteria.md
//   - business-rules.md
//   - decisions.md
func runDerive(args []string) int {
	fs := flag.NewFlagSet("derive", flag.ContinueOnError)
	outputDir := fs.String("output-dir", "", "Output directory for L1 documents (required)")
	analysisFile := fs.String("analysis-file", "", "Path to analysis JSON or interview state file")
	decisionsFile := fs.String("decisions", "", "Path to existing decisions.md")
	vocabularyFile := fs.String("vocabulary", "", "Domain vocabulary file for enhanced accuracy")
	nfrFile := fs.String("nfr", "", "Non-functional requirements file")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate required options
	if *outputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --output-dir is required")
		return domain.ExitCodeError
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Load input
	input, err := loadDeriveInput(*analysisFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Load optional files
	vocabulary := ""
	if *vocabularyFile != "" {
		data, err := os.ReadFile(*vocabularyFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read vocabulary: %v\n", err)
		} else {
			vocabulary = string(data)
		}
	}

	nfr := ""
	if *nfrFile != "" {
		data, err := os.ReadFile(*nfrFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read NFR: %v\n", err)
		} else {
			nfr = string(data)
		}
	}

	existingDecisions := ""
	if *decisionsFile != "" {
		data, err := os.ReadFile(*decisionsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read decisions: %v\n", err)
		} else {
			existingDecisions = string(data)
		}
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Generate Domain Model
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Phase 1: Generating domain model...")
	}
	domainModelMD, err := generateDomainModel(client, input, vocabulary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: domain model derivation failed: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "domain-model.md", domainModelMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 2: Generate Bounded Context Map
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Phase 2: Generating bounded context map...")
	}
	bcmMD, err := generateBoundedContextMap(client, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: bounded context map derivation failed: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeFile(*outputDir, "bounded-context-map.md", bcmMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 3: Generate Acceptance Criteria and Business Rules
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Phase 3: Generating acceptance criteria and business rules...")
	}
	derivationResult, err := generateACsAndBRs(client, input, nfr, existingDecisions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: derivation failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Format and write acceptance criteria
	acMD := formatter.FormatAcceptanceCriteria(derivationResult.AcceptanceCriteria)
	if err := writeFile(*outputDir, "acceptance-criteria.md", acMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Format and write business rules
	brMD := formatter.FormatBusinessRules(derivationResult.BusinessRules)
	if err := writeFile(*outputDir, "business-rules.md", brMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 4: Write Decisions
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Phase 4: Writing decisions...")
	}
	decisionsMD := formatDecisions(input.Decisions, derivationResult.Decisions)
	if err := writeFile(*outputDir, "decisions.md", decisionsMD); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive] Complete. Output written to %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// deriveInput holds input data for derivation.
type deriveInput struct {
	DomainModel  *domain.Domain
	Decisions    []domain.Decision
	InputContent string
}

// loadDeriveInput loads input from analysis JSON or interview state.
func loadDeriveInput(path string) (*deriveInput, error) {
	if path == "" {
		return nil, fmt.Errorf("--analysis-file is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis file: %w", err)
	}

	// Try interview state first
	var state domain.InterviewState
	if err := json.Unmarshal(data, &state); err == nil && state.DomainModel != nil {
		return &deriveInput{
			DomainModel:  state.DomainModel,
			Decisions:    state.Decisions,
			InputContent: state.InputContent,
		}, nil
	}

	// Try analysis result
	var analysis domain.AnalyzeResult
	if err := json.Unmarshal(data, &analysis); err == nil && analysis.DomainModel != nil {
		return &deriveInput{
			DomainModel:  analysis.DomainModel,
			Decisions:    analysis.Decisions,
			InputContent: analysis.InputContent,
		}, nil
	}

	return nil, fmt.Errorf("domain_model is required in input")
}

// generateDomainModel generates L1 domain model markdown.
func generateDomainModel(client *claude.Client, input *deriveInput, vocabulary string) (string, error) {
	// Build context
	domainJSON, _ := json.MarshalIndent(input.DomainModel, "", "  ")
	context := fmt.Sprintf("## Discovered Domain Model\n\n```json\n%s\n```", string(domainJSON))

	if vocabulary != "" {
		context += fmt.Sprintf("\n\n## Domain Vocabulary\n\n%s", vocabulary)
	}

	prompt := claude.BuildPrompt(prompts.DeriveDomainModel, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateBoundedContextMap generates L1 bounded context map markdown.
func generateBoundedContextMap(client *claude.Client, input *deriveInput) (string, error) {
	domainJSON, _ := json.MarshalIndent(input.DomainModel, "", "  ")
	context := fmt.Sprintf("## Discovered Domain Model\n\n```json\n%s\n```", string(domainJSON))

	prompt := claude.BuildPrompt(prompts.DeriveBoundedContext, context)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// generateACsAndBRs generates acceptance criteria and business rules.
func generateACsAndBRs(client *claude.Client, input *deriveInput, nfr, existingDecisions string) (*domain.DerivationResult, error) {
	domainJSON, _ := json.MarshalIndent(input.DomainModel, "", "  ")
	decisionsJSON, _ := json.MarshalIndent(input.Decisions, "", "  ")

	context := fmt.Sprintf("## Domain Model\n\n```json\n%s\n```\n\n## Decisions\n\n```json\n%s\n```",
		string(domainJSON), string(decisionsJSON))

	if nfr != "" {
		context += fmt.Sprintf("\n\n## Non-Functional Requirements\n\n%s", nfr)
	}

	if existingDecisions != "" {
		context += fmt.Sprintf("\n\n## Existing Decisions\n\n%s", existingDecisions)
	}

	if input.InputContent != "" {
		context += fmt.Sprintf("\n\n## Original L0 Input\n\n%s", input.InputContent)
	}

	prompt := claude.BuildPrompt(prompts.Derivation, context)

	var result domain.DerivationResult
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return &result, nil
}

// formatDecisions formats decisions as markdown.
func formatDecisions(existingDecisions, newDecisions []domain.Decision) string {
	all := append(existingDecisions, newDecisions...)
	return formatter.FormatDecisions(all)
}

// writeFile writes content to a file in the output directory.
func writeFile(dir, filename, content string) error {
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}
