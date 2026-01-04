// Package cmd provides CLI commands for loom-cli.
//
// This file implements the analyze command.
// Implements: l2/interface-contracts.md IC-ANL-001
// Implements: l2/tech-specs.md TS-ARCH-001a
// See: l2/sequence-design.md SEQ-CAS-001
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
	"loom-cli/prompts"
)

// runAnalyze implements the analyze command.
// Implements: IC-ANL-001
//
// Flags:
//
//	--input-file    Single L0 file to analyze
//	--input-dir     Directory containing L0 files
//	--decisions     Existing decisions file to filter resolved ambiguities
//	--output        Output file (default: stdout)
func runAnalyze(args []string) int {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	inputFile := fs.String("input-file", "", "Single L0 file to analyze")
	inputDir := fs.String("input-dir", "", "Directory containing L0 files")
	decisionsFile := fs.String("decisions", "", "Existing decisions file")
	outputFile := fs.String("output", "", "Output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate input
	if *inputFile == "" && *inputDir == "" {
		fmt.Fprintf(os.Stderr, "Error: either --input-file or --input-dir is required\n")
		return domain.ExitCodeError
	}

	// Read input content
	inputContent, inputFiles, err := readInputContent(*inputFile, *inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return domain.ExitCodeError
	}

	// Load existing decisions if provided
	var existingDecisions []domain.Decision
	if *decisionsFile != "" {
		existingDecisions, err = loadDecisions(*decisionsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading decisions: %v\n", err)
			return domain.ExitCodeError
		}
	}

	// Run analysis
	result, err := analyze(inputContent, inputFiles, existingDecisions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during analysis: %v\n", err)
		return domain.ExitCodeError
	}

	// Output result
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling result: %v\n", err)
		return domain.ExitCodeError
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, output, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			return domain.ExitCodeError
		}
	} else {
		fmt.Println(string(output))
	}

	return domain.ExitCodeSuccess
}

// analyze performs the analysis pipeline.
// Implements: l2/tech-specs.md TS-ARCH-001a
//
// Pipeline:
// 1. Domain discovery - extract entities, operations, relationships
// 2. Entity analysis - enhance entity details, find ambiguities
// 3. Operation analysis - enhance operation details, find ambiguities
// 4. Filter already resolved ambiguities
func analyze(inputContent string, inputFiles []string, existingDecisions []domain.Decision) (*domain.AnalyzeResult, error) {
	client := claude.NewClient()
	client.Verbose = Verbose()

	// Step 1: Domain Discovery
	// Implements: l2/prompt-catalog.md PRM-ANL-001
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[analyze] Step 1: Domain discovery...")
	}

	discoveryPrompt := claude.BuildPrompt(prompts.DomainDiscovery, inputContent)

	var discovery struct {
		Entities      []domain.Entity       `json:"entities"`
		Operations    []domain.Operation    `json:"operations"`
		Relationships []domain.Relationship `json:"relationships"`
		Aggregates    []domain.Aggregate    `json:"aggregates"`
		BusinessRules []string              `json:"business_rules"`
	}

	if err := client.CallJSONWithRetry(discoveryPrompt, &discovery, claude.DefaultRetryConfig()); err != nil {
		return nil, fmt.Errorf("domain discovery failed: %w", err)
	}

	// Step 2: Entity Analysis
	// Implements: l2/prompt-catalog.md PRM-ANL-002
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[analyze] Step 2: Entity analysis...")
	}

	entitiesJSON, _ := json.Marshal(discovery.Entities)
	entityPrompt := claude.BuildPrompt(prompts.EntityAnalysis, string(entitiesJSON))

	var entityAnalysis struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}

	if err := client.CallJSONWithRetry(entityPrompt, &entityAnalysis, claude.DefaultRetryConfig()); err != nil {
		return nil, fmt.Errorf("entity analysis failed: %w", err)
	}

	// Step 3: Operation Analysis
	// Implements: l2/prompt-catalog.md PRM-ANL-003
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[analyze] Step 3: Operation analysis...")
	}

	operationsJSON, _ := json.Marshal(discovery.Operations)
	operationPrompt := claude.BuildPrompt(prompts.OperationAnalysis, string(operationsJSON))

	var operationAnalysis struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}

	if err := client.CallJSONWithRetry(operationPrompt, &operationAnalysis, claude.DefaultRetryConfig()); err != nil {
		return nil, fmt.Errorf("operation analysis failed: %w", err)
	}

	// Combine ambiguities
	allAmbiguities := append(entityAnalysis.Ambiguities, operationAnalysis.Ambiguities...)

	// Step 4: Filter resolved ambiguities
	// Implements: DEC-L1-008 (automatic dependency inference)
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[analyze] Step 4: Filtering resolved ambiguities...")
	}

	resolvedIDs := make(map[string]bool)
	for _, d := range existingDecisions {
		resolvedIDs[d.ID] = true
	}

	var unresolvedAmbiguities []domain.Ambiguity
	for _, amb := range allAmbiguities {
		if !resolvedIDs[amb.ID] {
			unresolvedAmbiguities = append(unresolvedAmbiguities, amb)
		}
	}

	// Build result
	result := &domain.AnalyzeResult{
		DomainModel: &domain.Domain{
			Entities:      discovery.Entities,
			Operations:    discovery.Operations,
			Relationships: discovery.Relationships,
			Aggregates:    discovery.Aggregates,
			BusinessRules: discovery.BusinessRules,
		},
		Ambiguities:  unresolvedAmbiguities,
		Decisions:    existingDecisions,
		InputFiles:   inputFiles,
		InputContent: inputContent,
	}

	if Verbose() {
		fmt.Fprintf(os.Stderr, "[analyze] Done. Found %d entities, %d operations, %d ambiguities\n",
			len(discovery.Entities), len(discovery.Operations), len(unresolvedAmbiguities))
	}

	return result, nil
}

// readInputContent reads L0 content from file or directory.
func readInputContent(inputFile, inputDir string) (string, []string, error) {
	var content strings.Builder
	var files []string

	if inputFile != "" {
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return "", nil, err
		}
		content.Write(data)
		files = append(files, inputFile)
	}

	if inputDir != "" {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			return "", nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			path := filepath.Join(inputDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				return "", nil, err
			}

			if content.Len() > 0 {
				content.WriteString("\n\n---\n\n")
			}
			content.Write(data)
			files = append(files, path)
		}
	}

	return content.String(), files, nil
}

// loadDecisions loads existing decisions from a JSON file.
func loadDecisions(path string) ([]domain.Decision, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var decisions []domain.Decision
	if err := json.Unmarshal(data, &decisions); err != nil {
		// Try loading as interview state
		var state domain.InterviewState
		if err2 := json.Unmarshal(data, &state); err2 != nil {
			return nil, err
		}
		return state.Decisions, nil
	}

	return decisions, nil
}
