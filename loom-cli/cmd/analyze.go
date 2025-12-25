package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/internal/config"
	"github.com/ikadar/loom-cli/internal/domain"
	"github.com/ikadar/loom-cli/prompts"
)

// AnalyzeResult is the output of the analyze command
type AnalyzeResult struct {
	DomainModel   *domain.Domain      `json:"domain_model"`
	Ambiguities   []domain.Ambiguity  `json:"ambiguities"`
	Decisions     []domain.Decision   `json:"existing_decisions"`
	InputFiles    []string            `json:"input_files"`
	InputContent  string              `json:"input_content"`
}

func runAnalyze() error {
	// Parse arguments (skip "analyze" command itself)
	cfg, err := config.ParseArgsForAnalyze(os.Args[2:])
	if err != nil {
		return err
	}

	// Create Claude client
	client := claude.NewClient()

	// === PHASE 0: Read Input ===
	fmt.Fprintln(os.Stderr, "Phase 0: Reading input files...")

	inputContent, inputFiles, err := cfg.ReadInputFiles()
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "  Read %d file(s)\n", len(inputFiles))

	// Load existing decisions
	existingDecisions := loadDecisions(cfg.DecisionsFile)
	fmt.Fprintf(os.Stderr, "  Loaded %d existing decisions\n", len(existingDecisions))

	// === PHASE 1: Domain Discovery ===
	fmt.Fprintln(os.Stderr, "\nPhase 1: Discovering domain model...")

	domainModel, err := discoverDomain(client, inputContent)
	if err != nil {
		return fmt.Errorf("domain discovery failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Found: %d entities, %d operations, %d relationships\n",
		len(domainModel.Entities),
		len(domainModel.Operations),
		len(domainModel.Relationships))

	// === PHASE 2: Completeness Analysis ===
	fmt.Fprintln(os.Stderr, "\nPhase 2: Analyzing completeness...")

	entityAmbiguities, err := analyzeEntities(client, domainModel.Entities)
	if err != nil {
		return fmt.Errorf("entity analysis failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Entity ambiguities: %d\n", len(entityAmbiguities))

	operationAmbiguities, err := analyzeOperations(client, domainModel.Operations)
	if err != nil {
		return fmt.Errorf("operation analysis failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Operation ambiguities: %d\n", len(operationAmbiguities))

	// Combine all ambiguities
	allAmbiguities := append(entityAmbiguities, operationAmbiguities...)

	// === PHASE 3: Filter Already Resolved ===
	fmt.Fprintln(os.Stderr, "\nPhase 3: Filtering resolved ambiguities...")

	unresolvedAmbiguities := filterResolved(allAmbiguities, existingDecisions)
	fmt.Fprintf(os.Stderr, "  Unresolved: %d (of %d total)\n", len(unresolvedAmbiguities), len(allAmbiguities))

	// Output result as JSON to stdout
	result := AnalyzeResult{
		DomainModel:   domainModel,
		Ambiguities:   unresolvedAmbiguities,
		Decisions:     existingDecisions,
		InputFiles:    inputFiles,
		InputContent:  inputContent,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

// Phase 1: Domain Discovery
func discoverDomain(client *claude.Client, input string) (*domain.Domain, error) {
	prompt := prompts.DomainDiscovery + input

	var result domain.Domain
	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Phase 2: Analyze Entities
func analyzeEntities(client *claude.Client, entities []domain.Entity) ([]domain.Ambiguity, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	entitiesJSON, _ := json.MarshalIndent(entities, "", "  ")
	prompt := prompts.EntityAnalysis + string(entitiesJSON)

	var result struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}

	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return result.Ambiguities, nil
}

// Phase 2: Analyze Operations
func analyzeOperations(client *claude.Client, operations []domain.Operation) ([]domain.Ambiguity, error) {
	if len(operations) == 0 {
		return nil, nil
	}

	opsJSON, _ := json.MarshalIndent(operations, "", "  ")
	prompt := prompts.OperationAnalysis + string(opsJSON)

	var result struct {
		Ambiguities []domain.Ambiguity `json:"ambiguities"`
	}

	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return result.Ambiguities, nil
}

// Phase 3: Filter resolved ambiguities
func filterResolved(ambiguities []domain.Ambiguity, decisions []domain.Decision) []domain.Ambiguity {
	decisionMap := make(map[string]bool)
	for _, d := range decisions {
		// Match by question similarity (lowercase, trimmed)
		key := strings.ToLower(strings.TrimSpace(d.Question))
		decisionMap[key] = true
	}

	var unresolved []domain.Ambiguity
	for _, a := range ambiguities {
		key := strings.ToLower(strings.TrimSpace(a.Question))
		if !decisionMap[key] {
			unresolved = append(unresolved, a)
		}
	}

	return unresolved
}

// Load existing decisions from decisions.md
func loadDecisions(path string) []domain.Decision {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil // File doesn't exist, return empty
	}

	var decisions []domain.Decision
	lines := strings.Split(string(content), "\n")

	var currentDecision *domain.Decision
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "- **AMB-") {
			if currentDecision != nil && currentDecision.Question != "" {
				decisions = append(decisions, *currentDecision)
			}
			currentDecision = &domain.Decision{
				Source: "existing",
			}
			end := strings.Index(line, ":")
			if end > 6 {
				currentDecision.ID = line[4:end]
			}
		} else if currentDecision != nil {
			if strings.HasPrefix(line, "- Q:") {
				currentDecision.Question = strings.TrimSpace(line[4:])
			} else if strings.HasPrefix(line, "- A:") {
				currentDecision.Answer = strings.TrimSpace(line[4:])
			}
		}
	}

	if currentDecision != nil && currentDecision.Question != "" {
		decisions = append(decisions, *currentDecision)
	}

	return decisions
}
