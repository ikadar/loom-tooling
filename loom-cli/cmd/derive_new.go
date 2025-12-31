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
	DomainModel  *domain.Domain    `json:"domain_model"`
	Decisions    []domain.Decision `json:"decisions"`
	InputContent string            `json:"input_content"`
}

// DomainModelDoc represents the domain-model.md document structure
type DomainModelDoc struct {
	DomainModel struct {
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		BoundedContexts []string `json:"bounded_contexts"`
	} `json:"domain_model"`
	Entities     []DomainEntity     `json:"entities"`
	ValueObjects []DomainValueObject `json:"value_objects"`
	Summary      struct {
		AggregateRoots  int `json:"aggregate_roots"`
		Entities        int `json:"entities"`
		ValueObjects    int `json:"value_objects"`
		TotalOperations int `json:"total_operations"`
		TotalEvents     int `json:"total_events"`
	} `json:"summary"`
}

type DomainEntity struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Type          string              `json:"type"`
	Purpose       string              `json:"purpose"`
	Attributes    []EntityAttribute   `json:"attributes"`
	Invariants    []string            `json:"invariants"`
	Operations    []EntityOperation   `json:"operations"`
	Events        []EntityEvent       `json:"events"`
	Relationships []EntityRelationship `json:"relationships"`
}

type EntityAttribute struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
}

type EntityOperation struct {
	Name           string   `json:"name"`
	Signature      string   `json:"signature"`
	Preconditions  []string `json:"preconditions"`
	Postconditions []string `json:"postconditions"`
}

type EntityEvent struct {
	Name    string   `json:"name"`
	Trigger string   `json:"trigger"`
	Payload []string `json:"payload"`
}

type EntityRelationship struct {
	Target      string `json:"target"`
	Type        string `json:"type"`
	Cardinality string `json:"cardinality"`
}

type DomainValueObject struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Purpose    string            `json:"purpose"`
	Attributes []EntityAttribute `json:"attributes"`
	Operations []string          `json:"operations"`
}

// BoundedContextMap represents the bounded-context-map.md document structure
type BoundedContextMap struct {
	BoundedContexts      []BoundedContext      `json:"bounded_contexts"`
	ContextRelationships []ContextRelationship `json:"context_relationships"`
	ContextMapDiagram    string                `json:"context_map_diagram"`
	Summary              struct {
		TotalContexts            int      `json:"total_contexts"`
		TotalRelationships       int      `json:"total_relationships"`
		IntegrationPatternsUsed []string `json:"integration_patterns_used"`
	} `json:"summary"`
}

type BoundedContext struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Purpose            string            `json:"purpose"`
	CoreEntities       []string          `json:"core_entities"`
	Aggregates         []string          `json:"aggregates"`
	Capabilities       []string          `json:"capabilities"`
	UbiquitousLanguage map[string]string `json:"ubiquitous_language"`
}

type ContextRelationship struct {
	Upstream           string   `json:"upstream"`
	Downstream         string   `json:"downstream"`
	RelationshipType   string   `json:"relationship_type"`
	Description        string   `json:"description"`
	IntegrationPattern string   `json:"integration_pattern"`
	SharedData         []string `json:"shared_data"`
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

	// === PHASE 5a: Derive Domain Model Document ===
	fmt.Fprintln(os.Stderr, "Phase 5a: Deriving domain-model.md...")

	domainModelDoc, err := deriveDomainModelDoc(client, input.DomainModel, input.InputContent)
	if err != nil {
		return fmt.Errorf("domain model derivation failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d entities, %d value objects\n",
		len(domainModelDoc.Entities), len(domainModelDoc.ValueObjects))

	// === PHASE 5b: Derive Bounded Context Map ===
	fmt.Fprintln(os.Stderr, "\nPhase 5b: Deriving bounded-context-map.md...")

	boundedContextMap, err := deriveBoundedContextMap(client, domainModelDoc)
	if err != nil {
		return fmt.Errorf("bounded context map derivation failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d contexts, %d relationships\n",
		len(boundedContextMap.BoundedContexts), len(boundedContextMap.ContextRelationships))

	// === PHASE 5c: Derive AC and BR ===
	fmt.Fprintln(os.Stderr, "\nPhase 5c: Deriving acceptance-criteria.md and business-rules.md...")

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

	if err := writeOutputFiles(cfg, result, domainModelDoc, boundedContextMap, newDecisions); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Print summary
	printDeriveSummary(cfg, input.DomainModel, input.Decisions, result, domainModelDoc, boundedContextMap)

	return nil
}

// Phase 5a: Derive Domain Model document
func deriveDomainModelDoc(client *claude.Client, dm *domain.Domain, input string) (*DomainModelDoc, error) {
	domainJSON, _ := json.MarshalIndent(dm, "", "  ")

	prompt := prompts.DeriveDomainModel + "\n" + input + "\n\nEXISTING DOMAIN ANALYSIS:\n" + string(domainJSON)

	var result DomainModelDoc
	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Phase 5b: Derive Bounded Context Map
func deriveBoundedContextMap(client *claude.Client, domainModelDoc *DomainModelDoc) (*BoundedContextMap, error) {
	domainJSON, _ := json.MarshalIndent(domainModelDoc, "", "  ")

	prompt := prompts.DeriveBoundedContext + "\n" + string(domainJSON)

	var result BoundedContextMap
	if err := client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Phase 5c: Derive AC and BR
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
func writeOutputFiles(cfg *config.Config, result *domain.DerivationResult, domainModelDoc *DomainModelDoc, boundedContextMap *BoundedContextMap, newDecisions []domain.Decision) error {
	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return err
	}

	// Write domain-model.md
	dmPath := cfg.OutputDir + "/domain-model.md"
	dmContent := formatDomainModel(domainModelDoc)
	if err := os.WriteFile(dmPath, []byte(dmContent), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", dmPath)

	// Write bounded-context-map.md
	bcPath := cfg.OutputDir + "/bounded-context-map.md"
	bcContent := formatBoundedContextMap(boundedContextMap)
	if err := os.WriteFile(bcPath, []byte(bcContent), 0644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", bcPath)

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

// toAnchor converts an ID to a lowercase anchor (e.g., "AC-CUST-001" -> "ac-cust-001")
func toAnchor(id string) string {
	return strings.ToLower(id)
}

// toLink creates a markdown link with anchor (e.g., [AC-CUST-001](acceptance-criteria.md#ac-cust-001))
func toLink(id, file string) string {
	return fmt.Sprintf("[%s](%s#%s)", id, file, toAnchor(id))
}

func formatAC(acs []domain.AcceptanceCriteria) string {
	var sb strings.Builder

	sb.WriteString("# Acceptance Criteria\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("---\n\n")

	for _, ac := range acs {
		// Add anchor to heading: ## AC-CUST-001 – Title {#ac-cust-001}
		sb.WriteString(fmt.Sprintf("## %s – %s {#%s}\n\n", ac.ID, ac.Title, toAnchor(ac.ID)))
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
			// Link to decisions.md
			sb.WriteString(fmt.Sprintf("- Decision: %s\n", toLink(ref, "decisions.md")))
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
		// Add anchor to heading: ## BR-CUST-001 – Title {#br-cust-001}
		sb.WriteString(fmt.Sprintf("## %s – %s {#%s}\n\n", br.ID, br.Title, toAnchor(br.ID)))
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
			// Link to decisions.md
			sb.WriteString(fmt.Sprintf("- Decision: %s\n", toLink(ref, "decisions.md")))
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

func formatDomainModel(doc *DomainModelDoc) string {
	var sb strings.Builder

	sb.WriteString("# Domain Model\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Domain:** %s\n\n", doc.DomainModel.Name))
	sb.WriteString(fmt.Sprintf("%s\n\n", doc.DomainModel.Description))
	sb.WriteString("---\n\n")

	// Entities
	sb.WriteString("## Entities\n\n")
	for _, e := range doc.Entities {
		// Add anchor: ### ENT-Customer – Customer {#ent-customer}
		sb.WriteString(fmt.Sprintf("### %s – %s {#%s}\n\n", e.ID, e.Name, toAnchor(e.ID)))
		sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", e.Type))
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", e.Purpose))

		if len(e.Attributes) > 0 {
			sb.WriteString("**Attributes:**\n")
			sb.WriteString("| Name | Type | Constraints |\n")
			sb.WriteString("|------|------|-------------|\n")
			for _, attr := range e.Attributes {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", attr.Name, attr.Type, attr.Constraints))
			}
			sb.WriteString("\n")
		}

		if len(e.Invariants) > 0 {
			sb.WriteString("**Invariants:**\n")
			for _, inv := range e.Invariants {
				sb.WriteString(fmt.Sprintf("- %s\n", inv))
			}
			sb.WriteString("\n")
		}

		if len(e.Operations) > 0 {
			sb.WriteString("**Operations:**\n")
			for _, op := range e.Operations {
				sb.WriteString(fmt.Sprintf("- `%s`\n", op.Signature))
				if len(op.Preconditions) > 0 {
					sb.WriteString(fmt.Sprintf("  - Pre: %v\n", op.Preconditions))
				}
				if len(op.Postconditions) > 0 {
					sb.WriteString(fmt.Sprintf("  - Post: %v\n", op.Postconditions))
				}
			}
			sb.WriteString("\n")
		}

		if len(e.Events) > 0 {
			sb.WriteString("**Events:**\n")
			for _, evt := range e.Events {
				sb.WriteString(fmt.Sprintf("- `%s` (trigger: %s)\n", evt.Name, evt.Trigger))
			}
			sb.WriteString("\n")
		}

		if len(e.Relationships) > 0 {
			sb.WriteString("**Relationships:**\n")
			for _, rel := range e.Relationships {
				sb.WriteString(fmt.Sprintf("- %s %s (%s)\n", rel.Type, rel.Target, rel.Cardinality))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	// Value Objects
	if len(doc.ValueObjects) > 0 {
		sb.WriteString("## Value Objects\n\n")
		for _, vo := range doc.ValueObjects {
			// Add anchor: ### VO-Money – Money {#vo-money}
			sb.WriteString(fmt.Sprintf("### %s – %s {#%s}\n\n", vo.ID, vo.Name, toAnchor(vo.ID)))
			sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", vo.Purpose))

			if len(vo.Attributes) > 0 {
				sb.WriteString("**Attributes:**\n")
				for _, attr := range vo.Attributes {
					sb.WriteString(fmt.Sprintf("- `%s` (%s): %s\n", attr.Name, attr.Type, attr.Constraints))
				}
				sb.WriteString("\n")
			}

			if len(vo.Operations) > 0 {
				sb.WriteString(fmt.Sprintf("**Operations:** %v\n\n", vo.Operations))
			}

			sb.WriteString("---\n\n")
		}
	}

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Aggregate Roots: %d\n", doc.Summary.AggregateRoots))
	sb.WriteString(fmt.Sprintf("- Entities: %d\n", doc.Summary.Entities))
	sb.WriteString(fmt.Sprintf("- Value Objects: %d\n", doc.Summary.ValueObjects))
	sb.WriteString(fmt.Sprintf("- Total Operations: %d\n", doc.Summary.TotalOperations))
	sb.WriteString(fmt.Sprintf("- Total Events: %d\n", doc.Summary.TotalEvents))

	return sb.String()
}

func formatBoundedContextMap(bcm *BoundedContextMap) string {
	var sb strings.Builder

	sb.WriteString("# Bounded Context Map\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("---\n\n")

	// Bounded Contexts
	sb.WriteString("## Bounded Contexts\n\n")
	for _, bc := range bcm.BoundedContexts {
		// Add anchor: ### BC-Customer – Customer Context {#bc-customer}
		sb.WriteString(fmt.Sprintf("### %s – %s {#%s}\n\n", bc.ID, bc.Name, toAnchor(bc.ID)))
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", bc.Purpose))

		if len(bc.CoreEntities) > 0 {
			sb.WriteString(fmt.Sprintf("**Core Entities:** %v\n\n", bc.CoreEntities))
		}

		if len(bc.Aggregates) > 0 {
			sb.WriteString(fmt.Sprintf("**Aggregates:** %v\n\n", bc.Aggregates))
		}

		if len(bc.Capabilities) > 0 {
			sb.WriteString("**Capabilities:**\n")
			for _, cap := range bc.Capabilities {
				sb.WriteString(fmt.Sprintf("- %s\n", cap))
			}
			sb.WriteString("\n")
		}

		if len(bc.UbiquitousLanguage) > 0 {
			sb.WriteString("**Ubiquitous Language:**\n")
			for term, def := range bc.UbiquitousLanguage {
				sb.WriteString(fmt.Sprintf("- **%s**: %s\n", term, def))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	// Context Relationships
	if len(bcm.ContextRelationships) > 0 {
		sb.WriteString("## Context Relationships\n\n")
		sb.WriteString("| Upstream | Downstream | Type | Pattern |\n")
		sb.WriteString("|----------|------------|------|----------|\n")
		for _, rel := range bcm.ContextRelationships {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				rel.Upstream, rel.Downstream, rel.RelationshipType, rel.IntegrationPattern))
		}
		sb.WriteString("\n")
	}

	// Generate Context Map Diagram programmatically
	sb.WriteString("## Context Map Diagram\n\n")
	sb.WriteString("```mermaid\nflowchart TB\n")
	for _, bc := range bcm.BoundedContexts {
		// Create subgraph for each context
		sb.WriteString(fmt.Sprintf("    subgraph %s[\"%s\"]\n", bc.ID, bc.Name))
		for _, entity := range bc.CoreEntities {
			sb.WriteString(fmt.Sprintf("        %s_%s[%s]\n", bc.ID, entity, entity))
		}
		sb.WriteString("    end\n")
	}
	sb.WriteString("\n")
	// Add relationships
	for _, rel := range bcm.ContextRelationships {
		label := rel.IntegrationPattern
		if label == "" {
			label = rel.RelationshipType
		}
		sb.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", rel.Upstream, label, rel.Downstream))
	}
	sb.WriteString("```\n\n")

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Total Contexts: %d\n", bcm.Summary.TotalContexts))
	sb.WriteString(fmt.Sprintf("- Total Relationships: %d\n", bcm.Summary.TotalRelationships))
	if len(bcm.Summary.IntegrationPatternsUsed) > 0 {
		sb.WriteString(fmt.Sprintf("- Integration Patterns: %v\n", bcm.Summary.IntegrationPatternsUsed))
	}

	return sb.String()
}

func printDeriveSummary(cfg *config.Config, dm *domain.Domain, decisions []domain.Decision, result *domain.DerivationResult, domainModelDoc *DomainModelDoc, boundedContextMap *BoundedContextMap) {
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "   L1 DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "   (Strategic Design Layer)")
	fmt.Fprintln(os.Stderr, "========================================")

	fmt.Fprintln(os.Stderr, "\nInput:")
	fmt.Fprintf(os.Stderr, "  Entities (analyzed):    %d\n", len(dm.Entities))
	fmt.Fprintf(os.Stderr, "  Operations (analyzed):  %d\n", len(dm.Operations))
	fmt.Fprintf(os.Stderr, "  Decisions used:         %d\n", len(decisions))

	fmt.Fprintln(os.Stderr, "\nGenerated L1 Documents:")
	fmt.Fprintf(os.Stderr, "  Domain Model:           %d entities, %d value objects\n",
		len(domainModelDoc.Entities), len(domainModelDoc.ValueObjects))
	fmt.Fprintf(os.Stderr, "  Bounded Contexts:       %d contexts, %d relationships\n",
		len(boundedContextMap.BoundedContexts), len(boundedContextMap.ContextRelationships))
	fmt.Fprintf(os.Stderr, "  Acceptance Criteria:    %d\n", len(result.AcceptanceCriteria))
	fmt.Fprintf(os.Stderr, "  Business Rules:         %d\n", len(result.BusinessRules))

	fmt.Fprintln(os.Stderr, "\nOutput Files:")
	fmt.Fprintf(os.Stderr, "  %s/domain-model.md\n", cfg.OutputDir)
	fmt.Fprintf(os.Stderr, "  %s/bounded-context-map.md\n", cfg.OutputDir)
	fmt.Fprintf(os.Stderr, "  %s/acceptance-criteria.md\n", cfg.OutputDir)
	fmt.Fprintf(os.Stderr, "  %s/business-rules.md\n", cfg.OutputDir)
}
