// Package cmd provides CLI commands for loom-cli.
//
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
	"loom-cli/prompts"
)

// runAnalyze handles the analyze command.
//
// Implements: IC-ANL-001
// Flow:
//  1. Read input file(s)
//  2. Call Claude with domain-discovery prompt
//  3. Call Claude with entity-analysis prompt
//  4. Call Claude with operation-analysis prompt
//  5. Merge results into AnalyzeResult
//  6. Output JSON to stdout
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
	inputContent, inputFiles, err := readInputFiles(*inputFile, *inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read input: %v\n", err)
		return domain.ExitCodeError
	}

	// Load existing decisions if provided
	var existingDecisions []domain.Decision
	if *decisionsFile != "" {
		existingDecisions, err = loadDecisions(*decisionsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load decisions: %v\n", err)
		}
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Domain Discovery
	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 1: Domain Discovery")
	}
	domainModel, err := discoverDomain(client, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: domain discovery failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 2: Entity Analysis
	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 2: Entity Analysis")
	}
	entityAmbiguities, err := analyzeEntities(client, domainModel, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: entity analysis failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Phase 3: Operation Analysis
	if *verbose {
		fmt.Fprintln(os.Stderr, "[analyze] Phase 3: Operation Analysis")
	}
	operationAmbiguities, err := analyzeOperations(client, domainModel, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: operation analysis failed: %v\n", err)
		return domain.ExitCodeError
	}

	// Merge ambiguities and filter already resolved
	allAmbiguities := append(entityAmbiguities, operationAmbiguities...)
	allAmbiguities = addDependencies(allAmbiguities)
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
			fmt.Fprintf(os.Stderr, "Error: failed to write output: %v\n", err)
			return domain.ExitCodeError
		}
	} else {
		fmt.Println(string(output))
	}

	return domain.ExitCodeSuccess
}

// readInputFiles reads markdown files from file or directory.
func readInputFiles(inputFile, inputDir string) (string, []string, error) {
	var files []string
	var contents []string

	if inputFile != "" {
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read %s: %w", inputFile, err)
		}
		files = []string{inputFile}
		contents = []string{string(content)}
	} else {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read directory %s: %w", inputDir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".md") {
				continue
			}
			path := filepath.Join(inputDir, name)
			content, err := os.ReadFile(path)
			if err != nil {
				return "", nil, fmt.Errorf("failed to read %s: %w", path, err)
			}
			files = append(files, path)
			contents = append(contents, string(content))
		}
	}

	if len(files) == 0 {
		return "", nil, fmt.Errorf("no markdown files found")
	}

	return strings.Join(contents, "\n\n---\n\n"), files, nil
}

// loadDecisions loads existing decisions from markdown file.
func loadDecisions(path string) ([]domain.Decision, error) {
	// TODO: Parse decisions.md format
	// For now, return empty slice
	return nil, nil
}

// discoverDomain calls Claude to discover domain model from L0 input.
//
// Implements: TS-ARCH-001a Phase 1
func discoverDomain(client *claude.Client, inputContent string) (*domain.Domain, error) {
	prompt := claude.BuildPrompt(prompts.DomainDiscovery, inputContent)

	var result domain.Domain
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return &result, nil
}

// analyzeEntities calls Claude to analyze entities for completeness.
//
// Implements: TS-ARCH-001a Phase 2
func analyzeEntities(client *claude.Client, domainModel *domain.Domain, inputContent string) ([]domain.Ambiguity, error) {
	// Build context with domain model
	domainJSON, _ := json.MarshalIndent(domainModel, "", "  ")
	context := fmt.Sprintf("## Original L0 Input\n\n%s\n\n## Discovered Domain Model\n\n```json\n%s\n```",
		inputContent, string(domainJSON))

	prompt := claude.BuildPrompt(prompts.EntityAnalysis, context)

	var result []domain.Ambiguity
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return result, nil
}

// analyzeOperations calls Claude to analyze operations for completeness.
//
// Implements: TS-ARCH-001a Phase 2
func analyzeOperations(client *claude.Client, domainModel *domain.Domain, inputContent string) ([]domain.Ambiguity, error) {
	// Build context with domain model
	domainJSON, _ := json.MarshalIndent(domainModel, "", "  ")
	context := fmt.Sprintf("## Original L0 Input\n\n%s\n\n## Discovered Domain Model\n\n```json\n%s\n```",
		inputContent, string(domainJSON))

	prompt := claude.BuildPrompt(prompts.OperationAnalysis, context)

	var result []domain.Ambiguity
	if err := client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return result, nil
}

// addDependencies adds automatic dependency inference to questions.
//
// Implements: AGG-INT-001 (Automatic Dependency Inference)
// See: DEC-L1-008
func addDependencies(questions []domain.Ambiguity) []domain.Ambiguity {
	result := make([]domain.Ambiguity, len(questions))
	copy(result, questions)

	// Phase 1: Build capability question map
	questionMap := make(map[string]string) // keyword -> question ID
	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		if strings.Contains(qLower, "can") && strings.Contains(qLower, "deleted") {
			questionMap[q.Subject+"_delete"] = q.ID
		}
		if strings.Contains(qLower, "can") && strings.Contains(qLower, "modified") {
			questionMap[q.Subject+"_modify"] = q.ID
		}
		if (strings.Contains(qLower, "have") || strings.Contains(qLower, "support")) &&
			strings.Contains(qLower, "expir") {
			questionMap[q.Subject+"_expire"] = q.ID
		}
	}

	// Phase 2: Add dependencies
	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		// Deletion follow-up questions
		if strings.Contains(qLower, "after delet") ||
			strings.Contains(qLower, "when delet") ||
			strings.Contains(qLower, "deletion cascade") ||
			strings.Contains(qLower, "upon deletion") {
			if depID, ok := questionMap[q.Subject+"_delete"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"cannot be deleted",
						"no deletion",
						"not deletable",
						"cannot delete",
						"no, ",
						"soft delete only",
					},
				})
			}
		}

		// Modification follow-up questions
		if strings.Contains(qLower, "after modif") ||
			strings.Contains(qLower, "when modif") ||
			strings.Contains(qLower, "modification trigger") {
			if depID, ok := questionMap[q.Subject+"_modify"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"cannot be modified",
						"immutable",
						"no modification",
						"cannot modify",
					},
				})
			}
		}

		// Expiration follow-up questions
		if strings.Contains(qLower, "when expir") ||
			strings.Contains(qLower, "after expir") ||
			strings.Contains(qLower, "expiration notification") {
			if depID, ok := questionMap[q.Subject+"_expire"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"no expiration",
						"does not expire",
						"never expires",
					},
				})
			}
		}
	}

	return result
}

// filterResolvedAmbiguities removes questions that have already been answered.
func filterResolvedAmbiguities(ambiguities []domain.Ambiguity, decisions []domain.Decision) []domain.Ambiguity {
	decisionSet := make(map[string]bool)
	for _, d := range decisions {
		decisionSet[d.ID] = true
	}

	var filtered []domain.Ambiguity
	for _, a := range ambiguities {
		if !decisionSet[a.ID] {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
