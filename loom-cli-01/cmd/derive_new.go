// Package cmd provides CLI commands for loom-cli.
//
// This file implements the derive command (L1 derivation).
// Implements: l2/interface-contracts.md IC-DRV-001
// Implements: l2/tech-specs.md TS-ARCH-001b
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
	"loom-cli/prompts"
)

// runDerive implements the derive command (L1 generation).
// Implements: IC-DRV-001
//
// Flags:
//
//	--input-file      Interview state or analysis file
//	--output-dir      Output directory for L1 documents
//	--interactive     Enable interactive approval mode
//
// Output files:
//   - domain-model.md
//   - bounded-context-map.md
//   - acceptance-criteria.md
//   - business-rules.md
//   - decisions.md
func runDerive(args []string) int {
	fs := flag.NewFlagSet("derive", flag.ExitOnError)

	inputFile := fs.String("input-file", "", "Interview state or analysis file")
	outputDir := fs.String("output-dir", "./l1", "Output directory")
	interactive := fs.Bool("interactive", false, "Enable interactive approval")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --input-file is required\n")
		return domain.ExitCodeError
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Load input
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return domain.ExitCodeError
	}

	// Try to parse as interview state first, then as analysis result
	var domainModel *domain.Domain
	var decisions []domain.Decision
	var inputContent string

	var state domain.InterviewState
	if err := json.Unmarshal(data, &state); err == nil && state.DomainModel != nil {
		domainModel = state.DomainModel
		decisions = state.Decisions
		inputContent = state.InputContent
	} else {
		var analysis domain.AnalyzeResult
		if err := json.Unmarshal(data, &analysis); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing input file: %v\n", err)
			return domain.ExitCodeError
		}
		domainModel = analysis.DomainModel
		decisions = analysis.Decisions
		inputContent = analysis.InputContent
	}

	client := claude.NewClient()
	client.Verbose = Verbose()

	// Build context from domain model and decisions
	modelJSON, _ := json.MarshalIndent(domainModel, "", "  ")
	decisionsJSON, _ := json.MarshalIndent(decisions, "", "  ")
	context := fmt.Sprintf("Domain Model:\n%s\n\nDecisions:\n%s\n\nOriginal Input:\n%s",
		string(modelJSON), string(decisionsJSON), inputContent)

	// Phase 1: Generate Domain Model markdown
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive] Phase 1: Generating domain-model.md...")
	}

	domainModelMD, err := generateL1Document(client, prompts.DeriveDomainModel, context, "Domain Model", *interactive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating domain model: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeL1File(*outputDir, "domain-model.md", domainModelMD); err != nil {
		return domain.ExitCodeError
	}

	// Phase 2: Generate Bounded Context Map
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive] Phase 2: Generating bounded-context-map.md...")
	}

	boundedContextMD, err := generateL1Document(client, prompts.DeriveBoundedContext, context, "Bounded Context Map", *interactive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating bounded context map: %v\n", err)
		return domain.ExitCodeError
	}
	if err := writeL1File(*outputDir, "bounded-context-map.md", boundedContextMD); err != nil {
		return domain.ExitCodeError
	}

	// Phase 3: Generate Acceptance Criteria and Business Rules
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive] Phase 3: Generating acceptance-criteria.md and business-rules.md...")
	}

	derivationPrompt := claude.BuildPrompt(prompts.Derivation, context)

	var derivation struct {
		AcceptanceCriteria []struct {
			ID         string   `json:"id"`
			Title      string   `json:"title"`
			Given      string   `json:"given"`
			When       string   `json:"when"`
			Then       string   `json:"then"`
			ErrorCases []string `json:"error_cases"`
		} `json:"acceptance_criteria"`
		BusinessRules []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Rule        string `json:"rule"`
			Invariant   string `json:"invariant"`
			Enforcement string `json:"enforcement"`
			ErrorCode   string `json:"error_code"`
		} `json:"business_rules"`
	}

	if err := client.CallJSONWithRetry(derivationPrompt, &derivation, claude.DefaultRetryConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating AC/BR: %v\n", err)
		return domain.ExitCodeError
	}

	// Generate acceptance-criteria.md
	acMD := generateACMarkdown(derivation.AcceptanceCriteria)
	if err := writeL1File(*outputDir, "acceptance-criteria.md", acMD); err != nil {
		return domain.ExitCodeError
	}

	// Generate business-rules.md
	brMD := generateBRMarkdown(derivation.BusinessRules)
	if err := writeL1File(*outputDir, "business-rules.md", brMD); err != nil {
		return domain.ExitCodeError
	}

	// Phase 4: Generate decisions.md from interview decisions
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive] Phase 4: Generating decisions.md...")
	}

	decisionsMD := generateDecisionsMarkdown(decisions)
	if err := writeL1File(*outputDir, "decisions.md", decisionsMD); err != nil {
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive] L1 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// generateL1Document generates a markdown document from a prompt.
func generateL1Document(client *claude.Client, prompt, context, name string, interactive bool) (string, error) {
	fullPrompt := claude.BuildPrompt(prompt, context)

	var result struct {
		Content string `json:"content,omitempty"`
	}

	// First try JSON response
	response, err := client.Call(fullPrompt)
	if err != nil {
		return "", err
	}

	// Try to parse as JSON with content field
	if err := json.Unmarshal([]byte(response), &result); err == nil && result.Content != "" {
		return result.Content, nil
	}

	// Otherwise use raw response (it's likely already markdown)
	return response, nil
}

// generateACMarkdown generates acceptance-criteria.md content.
func generateACMarkdown(acs []struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Given      string   `json:"given"`
	When       string   `json:"when"`
	Then       string   `json:"then"`
	ErrorCases []string `json:"error_cases"`
}) string {
	md := fmt.Sprintf(`---
title: "Acceptance Criteria"
generated: %s
status: draft
level: L1
---

# Acceptance Criteria

`, time.Now().Format(time.RFC3339))

	for _, ac := range acs {
		md += fmt.Sprintf("## %s – %s\n\n", ac.ID, ac.Title)
		md += fmt.Sprintf("**Given** %s\n", ac.Given)
		md += fmt.Sprintf("**When** %s\n", ac.When)
		md += fmt.Sprintf("**Then** %s\n\n", ac.Then)

		if len(ac.ErrorCases) > 0 {
			md += "**Error Cases:**\n"
			for _, ec := range ac.ErrorCases {
				md += fmt.Sprintf("- %s\n", ec)
			}
			md += "\n"
		}
		md += "---\n\n"
	}

	return md
}

// generateBRMarkdown generates business-rules.md content.
func generateBRMarkdown(brs []struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rule        string `json:"rule"`
	Invariant   string `json:"invariant"`
	Enforcement string `json:"enforcement"`
	ErrorCode   string `json:"error_code"`
}) string {
	md := fmt.Sprintf(`---
title: "Business Rules"
generated: %s
status: draft
level: L1
---

# Business Rules

`, time.Now().Format(time.RFC3339))

	for _, br := range brs {
		md += fmt.Sprintf("## %s – %s\n\n", br.ID, br.Title)
		md += fmt.Sprintf("**Rule:** %s\n\n", br.Rule)
		md += fmt.Sprintf("**Invariant:** %s\n\n", br.Invariant)
		md += fmt.Sprintf("**Enforcement:** %s\n\n", br.Enforcement)
		if br.ErrorCode != "" {
			md += fmt.Sprintf("**Error Code:** `%s`\n\n", br.ErrorCode)
		}
		md += "---\n\n"
	}

	return md
}

// generateDecisionsMarkdown generates decisions.md from interview decisions.
func generateDecisionsMarkdown(decisions []domain.Decision) string {
	md := fmt.Sprintf(`---
title: "Design Decisions"
generated: %s
status: draft
level: L1
---

# Design Decisions

This document records decisions made during the structured interview phase.

`, time.Now().Format(time.RFC3339))

	// Group by category
	byCategory := make(map[string][]domain.Decision)
	for _, d := range decisions {
		cat := d.Category
		if cat == "" {
			cat = "general"
		}
		byCategory[cat] = append(byCategory[cat], d)
	}

	for cat, decs := range byCategory {
		md += fmt.Sprintf("## %s Decisions\n\n", cat)
		for _, d := range decs {
			md += fmt.Sprintf("### %s\n\n", d.ID)
			md += fmt.Sprintf("**Question:** %s\n\n", d.Question)
			md += fmt.Sprintf("**Answer:** %s\n\n", d.Answer)
			md += fmt.Sprintf("**Source:** %s\n\n", d.Source)
			if d.Subject != "" {
				md += fmt.Sprintf("**Subject:** %s\n\n", d.Subject)
			}
			md += "---\n\n"
		}
	}

	return md
}

// writeL1File writes a file to the output directory.
func writeL1File(outputDir, filename, content string) error {
	path := filepath.Join(outputDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filename, err)
		return err
	}
	if Verbose() {
		fmt.Fprintf(os.Stderr, "[derive] Wrote %s\n", path)
	}
	return nil
}
