// Implements: l2/interface-contracts.md IC-ANL-001
// See: l2/sequence-design.md SEQ-ANL-001
// See: l2/tech-specs.md TS-ARCH-001a
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
)

// runAnalyze implements the analyze command.
//
// Implements: IC-ANL-001
// Options:
//   - --input-file <path>: Single L0 input file
//   - --input-dir <path>: Directory with L0 files
//   - --output <path>: Output file (default: stdout)
//   - --decisions <path>: Existing decisions.md
func runAnalyze(args []string) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	inputFile := fs.String("input-file", "", "Single L0 input file (markdown)")
	inputDir := fs.String("input-dir", "", "Directory containing L0 input files")
	outputFile := fs.String("output", "", "Output file (default: stdout)")
	decisionsFile := fs.String("decisions", "", "Path to existing decisions.md")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate input options
	if *inputFile == "" && *inputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --input-file or --input-dir required")
		return domain.ExitCodeError
	}
	if *inputFile != "" && *inputDir != "" {
		fmt.Fprintln(os.Stderr, "Error: cannot specify both --input-file and --input-dir")
		return domain.ExitCodeError
	}

	// Read input content
	var inputContent string
	var inputFiles []string
	var err error

	if *inputFile != "" {
		content, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read input file: %v\n", err)
			return domain.ExitCodeError
		}
		inputContent = string(content)
		inputFiles = []string{*inputFile}
	} else {
		inputContent, inputFiles, err = readInputDir(*inputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read input directory: %v\n", err)
			return domain.ExitCodeError
		}
	}

	// Load existing decisions if provided
	var existingDecisions []domain.Decision
	if *decisionsFile != "" {
		existingDecisions, err = parseDecisionsFile(*decisionsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse decisions file: %v\n", err)
		}
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Domain Discovery
	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 1: Domain discovery...")
	}
	domainModel, err := discoverDomain(client, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: domain discovery failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 2: Completeness Analysis
	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 2: Entity analysis...")
	}
	entityAmbiguities, err := analyzeEntities(client, domainModel.Entities, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: entity analysis failed: %v\n", err)
		return domain.ExitCodeError
	}

	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 2: Operation analysis...")
	}
	operationAmbiguities, err := analyzeOperations(client, domainModel.Operations, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: operation analysis failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Merge ambiguities
	allAmbiguities := append(entityAmbiguities, operationAmbiguities...)

	// Phase 3: Filter already resolved
	filteredAmbiguities := filterResolvedAmbiguities(allAmbiguities, existingDecisions)

	// Build result
	result := domain.AnalyzeResult{
		DomainModel:  domainModel,
		Ambiguities:  filteredAmbiguities,
		Decisions:    existingDecisions,
		InputFiles:   inputFiles,
		InputContent: inputContent,
	}

	// Output JSON
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to marshal result: %v\n", err)
		return domain.ExitCodeError
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, output, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write output file: %v\n", err)
			return domain.ExitCodeError
		}
	} else {
		fmt.Println(string(output))
	}

	return domain.ExitCodeSuccess
}

// readInputDir reads all markdown files from a directory.
func readInputDir(dir string) (string, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", nil, err
	}

	var contents []string
	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
		contents = append(contents, string(content))
		files = append(files, path)
	}

	return strings.Join(contents, "\n\n---\n\n"), files, nil
}

// discoverDomain calls Claude to discover the domain model from L0 input.
func discoverDomain(client *claude.Client, inputContent string) (*domain.Domain, error) {
	prompt := buildDomainDiscoveryPrompt(inputContent)

	var result domain.Domain
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return &result, nil
}

// analyzeEntities calls Claude to analyze entities for completeness.
func analyzeEntities(client *claude.Client, entities []domain.Entity, inputContent string) ([]domain.Ambiguity, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	prompt := buildEntityAnalysisPrompt(entities, inputContent)

	var result struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return result.Ambiguities, nil
}

// analyzeOperations calls Claude to analyze operations for completeness.
func analyzeOperations(client *claude.Client, operations []domain.Operation, inputContent string) ([]domain.Ambiguity, error) {
	if len(operations) == 0 {
		return nil, nil
	}

	prompt := buildOperationAnalysisPrompt(operations, inputContent)

	var result struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return result.Ambiguities, nil
}

// filterResolvedAmbiguities removes ambiguities that have already been resolved.
func filterResolvedAmbiguities(ambiguities []domain.Ambiguity, decisions []domain.Decision) []domain.Ambiguity {
	if len(decisions) == 0 {
		return ambiguities
	}

	resolved := make(map[string]bool)
	for _, d := range decisions {
		resolved[d.ID] = true
	}

	var filtered []domain.Ambiguity
	for _, a := range ambiguities {
		if !resolved[a.ID] {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// parseDecisionsFile parses an existing decisions.md file.
func parseDecisionsFile(path string) ([]domain.Decision, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Simple parsing: extract decision blocks
	// This is a simplified implementation - full parsing would be more complex
	_ = content // TODO: Implement proper decisions.md parsing
	return nil, nil
}

// buildDomainDiscoveryPrompt builds the prompt for domain discovery.
func buildDomainDiscoveryPrompt(inputContent string) string {
	return fmt.Sprintf(`You are a domain-driven design expert. Analyze the following user story/requirements document and extract the domain model.

<context>
%s
</context>

Extract and return a JSON object with this structure:
{
  "entities": [
    {
      "name": "EntityName",
      "mentioned_attributes": ["attr1", "attr2"],
      "mentioned_operations": ["op1", "op2"],
      "mentioned_states": ["state1", "state2"]
    }
  ],
  "operations": [
    {
      "name": "OperationName",
      "actor": "User/System",
      "trigger": "When/What triggers this",
      "target": "Target entity",
      "mentioned_inputs": ["input1"],
      "mentioned_rules": ["rule1"]
    }
  ],
  "relationships": [
    {
      "from": "Entity1",
      "to": "Entity2",
      "type": "has_many|belongs_to|has_one|contains|references",
      "cardinality": "1:1|1:N|N:1|N:M"
    }
  ],
  "business_rules": ["rule1", "rule2"],
  "ui_mentions": ["UI element or screen mentioned"]
}

Return ONLY the JSON object, no explanation.`, inputContent)
}

// buildEntityAnalysisPrompt builds the prompt for entity completeness analysis.
func buildEntityAnalysisPrompt(entities []domain.Entity, inputContent string) string {
	entitiesJSON, _ := json.MarshalIndent(entities, "", "  ")

	return fmt.Sprintf(`You are a requirements analyst. Analyze the following entities and identify ambiguities or missing information.

<context>
Original Document:
%s

Discovered Entities:
%s
</context>

For each entity, identify:
1. Missing attributes that should be specified
2. Unclear state transitions
3. Missing validation rules
4. Ambiguous relationships

Return a JSON object with this structure:
{
  "ambiguities": [
    {
      "id": "Q-ENT-001",
      "category": "entity",
      "subject": "EntityName",
      "question": "What is the maximum length of X?",
      "severity": "critical|important|minor",
      "suggested_answer": "Suggested default if any",
      "options": ["Option 1", "Option 2"],
      "checklist_item": "Define maximum length for X"
    }
  ]
}

Return ONLY the JSON object, no explanation.`, inputContent, string(entitiesJSON))
}

// buildOperationAnalysisPrompt builds the prompt for operation completeness analysis.
func buildOperationAnalysisPrompt(operations []domain.Operation, inputContent string) string {
	operationsJSON, _ := json.MarshalIndent(operations, "", "  ")

	return fmt.Sprintf(`You are a requirements analyst. Analyze the following operations and identify ambiguities or missing information.

<context>
Original Document:
%s

Discovered Operations:
%s
</context>

For each operation, identify:
1. Missing error handling specifications
2. Unclear preconditions
3. Missing authorization rules
4. Ambiguous success/failure criteria

Return a JSON object with this structure:
{
  "ambiguities": [
    {
      "id": "Q-OP-001",
      "category": "operation",
      "subject": "OperationName",
      "question": "What happens if X fails?",
      "severity": "critical|important|minor",
      "suggested_answer": "Suggested default if any",
      "options": ["Option 1", "Option 2"],
      "checklist_item": "Define error handling for X"
    }
  ]
}

Return ONLY the JSON object, no explanation.`, inputContent, string(operationsJSON))
}
