package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/internal/config"
	"github.com/ikadar/loom-cli/internal/domain"
	"github.com/ikadar/loom-cli/prompts"
)

// DeriveInput is the expected input for the derive command
type DeriveInput struct {
	DomainModel  *domain.Domain     `json:"domain_model"`
	Decisions    []domain.Decision  `json:"decisions"`
	InputContent string             `json:"input_content"`
}

func runDeriveNew() error {
	// Parse arguments
	cfg, err := config.ParseArgsForDerive(os.Args[2:])
	if err != nil {
		return err
	}

	// Read input JSON from stdin or --analysis-file
	var input DeriveInput

	if cfg.AnalysisFile != "" {
		// Read from file
		content, err := os.ReadFile(cfg.AnalysisFile)
		if err != nil {
			return fmt.Errorf("failed to read analysis file: %w", err)
		}
		if err := json.Unmarshal(content, &input); err != nil {
			return fmt.Errorf("failed to parse analysis file: %w", err)
		}
	} else {
		// Read from stdin
		decoder := json.NewDecoder(os.Stdin)
		if err := decoder.Decode(&input); err != nil {
			return fmt.Errorf("failed to read input from stdin: %w", err)
		}
	}

	// Validate input
	if input.DomainModel == nil {
		return fmt.Errorf("domain_model is required in input")
	}

	// Create Claude client
	client := claude.NewClient()

	// === PHASE 5: Derivation ===
	fmt.Fprintln(os.Stderr, "Phase 5: Deriving L1 documents...")

	result, err := deriveDocuments(client, input.DomainModel, input.Decisions, input.InputContent)
	if err != nil {
		return fmt.Errorf("derivation failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Generated: %d ACs, %d BRs\n",
		len(result.AcceptanceCriteria),
		len(result.BusinessRules))

	// === PHASE 6: Write Output ===
	fmt.Fprintln(os.Stderr, "\nPhase 6: Writing output...")

	// Find new decisions (those with source != "existing")
	var newDecisions []domain.Decision
	for _, d := range input.Decisions {
		if d.Source != "existing" {
			newDecisions = append(newDecisions, d)
		}
	}

	if err := writeOutputFiles(cfg, result, newDecisions); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Print summary
	printDeriveSummary(cfg, input.DomainModel, input.Decisions, result)

	return nil
}

// Phase 5: Derive AC and BR
func deriveDocuments(client *claude.Client, dm *domain.Domain, decisions []domain.Decision, input string) (*domain.DerivationResult, error) {
	domainJSON, _ := json.MarshalIndent(dm, "", "  ")
	decisionsJSON, _ := json.MarshalIndent(decisions, "", "  ")

	prompt := prompts.DerivationPrompt + "\n\nDOMAIN MODEL:\n" + string(domainJSON) + "\n\nDECISIONS:\n" + string(decisionsJSON)

	var result struct {
		AcceptanceCriteria []domain.AcceptanceCriteria `json:"acceptance_criteria"`
		BusinessRules      []domain.BusinessRule       `json:"business_rules"`
	}

	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return &domain.DerivationResult{
		AcceptanceCriteria: result.AcceptanceCriteria,
		BusinessRules:      result.BusinessRules,
		Decisions:          decisions,
	}, nil
}

// Phase 6: Write output files
func writeOutputFiles(cfg *config.Config, result *domain.DerivationResult, newDecisions []domain.Decision) error {
	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return err
	}

	// Write acceptance-criteria.md
	acPath := cfg.OutputDir + "/acceptance-criteria.md"
	acContent := formatAC(result.AcceptanceCriteria)
	if err := os.WriteFile(acPath, []byte(acContent), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", acPath)

	// Write business-rules.md
	brPath := cfg.OutputDir + "/business-rules.md"
	brContent := formatBR(result.BusinessRules)
	if err := os.WriteFile(brPath, []byte(brContent), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", brPath)

	// Append new decisions to decisions.md
	if len(newDecisions) > 0 {
		decContent := formatDecisions(newDecisions)
		f, err := os.OpenFile(cfg.DecisionsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(decContent); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "  Updated: %s (+%d decisions)\n", cfg.DecisionsFile, len(newDecisions))
	}

	return nil
}

func formatAC(acs []domain.AcceptanceCriteria) string {
	var sb strings.Builder

	sb.WriteString("# Acceptance Criteria\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("---\n\n")

	for _, ac := range acs {
		sb.WriteString(fmt.Sprintf("## %s – %s\n\n", ac.ID, ac.Title))
		sb.WriteString(fmt.Sprintf("**Given** %s\n", ac.Given))
		sb.WriteString(fmt.Sprintf("**When** %s\n", ac.When))
		sb.WriteString(fmt.Sprintf("**Then** %s\n\n", ac.Then))

		if len(ac.ErrorCases) > 0 {
			sb.WriteString("**Error Cases:**\n")
			for _, ec := range ac.ErrorCases {
				sb.WriteString(fmt.Sprintf("- %s\n", ec))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("**Traceability:**\n")
		for _, ref := range ac.SourceRefs {
			sb.WriteString(fmt.Sprintf("- Source: %s\n", ref))
		}
		for _, ref := range ac.DecisionRefs {
			sb.WriteString(fmt.Sprintf("- Decision: %s\n", ref))
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatBR(brs []domain.BusinessRule) string {
	var sb strings.Builder

	sb.WriteString("# Business Rules\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("---\n\n")

	for _, br := range brs {
		sb.WriteString(fmt.Sprintf("## %s – %s\n\n", br.ID, br.Title))
		sb.WriteString(fmt.Sprintf("**Rule:** %s\n\n", br.Rule))
		sb.WriteString(fmt.Sprintf("**Invariant:** %s\n\n", br.Invariant))
		sb.WriteString(fmt.Sprintf("**Enforcement:** %s\n\n", br.Enforcement))

		if br.ErrorCode != "" {
			sb.WriteString(fmt.Sprintf("**Error Code:** `%s`\n\n", br.ErrorCode))
		}

		sb.WriteString("**Traceability:**\n")
		for _, ref := range br.SourceRefs {
			sb.WriteString(fmt.Sprintf("- Source: %s\n", ref))
		}
		for _, ref := range br.DecisionRefs {
			sb.WriteString(fmt.Sprintf("- Decision: %s\n", ref))
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatDecisions(decisions []domain.Decision) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n## Decisions from %s\n\n", time.Now().Format("2006-01-02")))

	for _, d := range decisions {
		sb.WriteString(fmt.Sprintf("- **%s: %s**\n", d.ID, d.Subject))
		sb.WriteString(fmt.Sprintf("  - Q: %s\n", d.Question))
		sb.WriteString(fmt.Sprintf("  - A: %s\n", d.Answer))
		sb.WriteString(fmt.Sprintf("  - Decided: %s by %s\n\n", d.DecidedAt.Format("2006-01-02 15:04"), d.Source))
	}

	return sb.String()
}

func printDeriveSummary(cfg *config.Config, dm *domain.Domain, decisions []domain.Decision, result *domain.DerivationResult) {
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "           DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "========================================")

	fmt.Fprintln(os.Stderr, "\nDomain:")
	fmt.Fprintf(os.Stderr, "  Entities:      %d\n", len(dm.Entities))
	fmt.Fprintf(os.Stderr, "  Operations:    %d\n", len(dm.Operations))
	fmt.Fprintf(os.Stderr, "  Relationships: %d\n", len(dm.Relationships))

	fmt.Fprintln(os.Stderr, "\nDecisions used:")
	fmt.Fprintf(os.Stderr, "  Total: %d\n", len(decisions))

	fmt.Fprintln(os.Stderr, "\nGenerated:")
	fmt.Fprintf(os.Stderr, "  Acceptance Criteria: %d\n", len(result.AcceptanceCriteria))
	fmt.Fprintf(os.Stderr, "  Business Rules:      %d\n", len(result.BusinessRules))

	fmt.Fprintln(os.Stderr, "\nOutput:")
	fmt.Fprintf(os.Stderr, "  %s/acceptance-criteria.md\n", cfg.OutputDir)
	fmt.Fprintf(os.Stderr, "  %s/business-rules.md\n", cfg.OutputDir)
}
