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
	"time"

	"loom-cli/internal/claude"
	"loom-cli/internal/domain"
)

// runDerive implements the derive command for L1 generation.
//
// Implements: IC-DRV-001
// Output files:
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
	vocabularyFile := fs.String("vocabulary", "", "Domain vocabulary file")
	nfrFile := fs.String("nfr", "", "Non-functional requirements file")
	verbose := fs.Bool("verbose", false, "Verbose output")
	interactive := fs.Bool("interactive", false, "Interactive approval mode")
	fs.Bool("i", false, "Alias for --interactive")

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
		fmt.Fprintf(os.Stderr, "Error: failed to create output directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Load analysis/interview state
	var domainModel *domain.Domain
	var decisions []domain.Decision
	var inputContent string

	if *analysisFile != "" {
		data, err := os.ReadFile(*analysisFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read analysis file: %v\n", err)
			return domain.ExitCodeError
		}

		// Try interview state first
		var interviewState domain.InterviewState
		if err := json.Unmarshal(data, &interviewState); err == nil && interviewState.DomainModel != nil {
			domainModel = interviewState.DomainModel
			decisions = interviewState.Decisions
			inputContent = interviewState.InputContent
		} else {
			// Try analysis result
			var analysis domain.AnalyzeResult
			if err := json.Unmarshal(data, &analysis); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to parse analysis file: %v\n", err)
				return domain.ExitCodeError
			}
			domainModel = analysis.DomainModel
			decisions = analysis.Decisions
			inputContent = analysis.InputContent
		}
	}

	if domainModel == nil {
		fmt.Fprintln(os.Stderr, "Error: domain_model is required in input")
		return domain.ExitCodeError
	}

	// Load optional files
	var vocabulary, nfr string
	if *vocabularyFile != "" {
		data, _ := os.ReadFile(*vocabularyFile)
		vocabulary = string(data)
	}
	if *nfrFile != "" {
		data, _ := os.ReadFile(*nfrFile)
		nfr = string(data)
	}

	// Load existing decisions
	if *decisionsFile != "" {
		// TODO: Parse existing decisions and merge
	}

	// Create Claude client
	client := claude.NewClient()
	client.Verbose = *verbose

	// Phase 1: Generate Domain Model
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Generating domain-model.md...")
	}
	domainModelMD, err := deriveDomainModel(client, domainModel, vocabulary, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: domain model derivation failed: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		domainModelMD, _ = interactiveApproval("domain-model.md", domainModelMD, client, func() (string, error) {
			return deriveDomainModel(client, domainModel, vocabulary, inputContent)
		})
	}

	// Phase 2: Generate Bounded Context Map
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Generating bounded-context-map.md...")
	}
	boundedContextMD, err := deriveBoundedContext(client, domainModel, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: bounded context map derivation failed: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		boundedContextMD, _ = interactiveApproval("bounded-context-map.md", boundedContextMD, client, func() (string, error) {
			return deriveBoundedContext(client, domainModel, inputContent)
		})
	}

	// Phase 3: Generate Acceptance Criteria & Business Rules
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Generating acceptance-criteria.md and business-rules.md...")
	}
	acMD, brMD, err := deriveACAndBR(client, domainModel, decisions, nfr, inputContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: derivation failed: %v\n", err)
		return domain.ExitCodeError
	}
	if *interactive {
		acMD, _ = interactiveApproval("acceptance-criteria.md", acMD, client, func() (string, error) {
			ac, _, e := deriveACAndBR(client, domainModel, decisions, nfr, inputContent)
			return ac, e
		})
		brMD, _ = interactiveApproval("business-rules.md", brMD, client, func() (string, error) {
			_, br, e := deriveACAndBR(client, domainModel, decisions, nfr, inputContent)
			return br, e
		})
	}

	// Phase 4: Generate Decisions
	if *verbose {
		fmt.Fprintln(os.Stderr, "[derive] Generating decisions.md...")
	}
	decisionsMD := formatDecisions(decisions)

	// Write all files
	files := map[string]string{
		"domain-model.md":        domainModelMD,
		"bounded-context-map.md": boundedContextMD,
		"acceptance-criteria.md": acMD,
		"business-rules.md":      brMD,
		"decisions.md":           decisionsMD,
	}

	for name, content := range files {
		path := filepath.Join(*outputDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to write %s: %v\n", name, err)
			return domain.ExitCodeError
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "[derive] Written: %s\n", path)
		}
	}

	fmt.Fprintf(os.Stderr, "L1 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// deriveDomainModel generates domain-model.md content.
func deriveDomainModel(client *claude.Client, dm *domain.Domain, vocabulary, inputContent string) (string, error) {
	dmJSON, _ := json.MarshalIndent(dm, "", "  ")

	prompt := fmt.Sprintf(`You are a DDD expert. Generate a domain-model.md document from the following analysis.

<context>
Domain Model Analysis:
%s

Original Input:
%s

Vocabulary:
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L1)
2. Overview section
3. Entities section with ID pattern: ENT-{NAME}
4. Value Objects section
5. Domain Events section
6. Traceability section

Use proper markdown formatting. Include traceability IDs.`, string(dmJSON), inputContent, vocabulary)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}

// deriveBoundedContext generates bounded-context-map.md content.
func deriveBoundedContext(client *claude.Client, dm *domain.Domain, inputContent string) (string, error) {
	dmJSON, _ := json.MarshalIndent(dm, "", "  ")

	prompt := fmt.Sprintf(`You are a DDD expert. Generate a bounded-context-map.md document.

<context>
Domain Model:
%s

Original Input:
%s
</context>

Generate a markdown document with:
1. YAML frontmatter (title, generated date, status: draft, level: L1)
2. Overview section
3. Bounded Contexts section with ID pattern: BC-{NAME}
4. Context Relationships (upstream/downstream, ACL, etc.)
5. Mermaid diagram showing context map
6. Traceability section

Use proper markdown formatting.`, string(dmJSON), inputContent)

	response, err := client.Call(prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}

// deriveACAndBR generates acceptance-criteria.md and business-rules.md content.
func deriveACAndBR(client *claude.Client, dm *domain.Domain, decisions []domain.Decision, nfr, inputContent string) (string, string, error) {
	dmJSON, _ := json.MarshalIndent(dm, "", "  ")
	decJSON, _ := json.MarshalIndent(decisions, "", "  ")

	prompt := fmt.Sprintf(`You are a requirements analyst. Generate acceptance-criteria.md and business-rules.md documents.

<context>
Domain Model:
%s

Decisions:
%s

NFR:
%s

Original Input:
%s
</context>

Generate TWO documents separated by "---DOCUMENT_SEPARATOR---":

DOCUMENT 1: acceptance-criteria.md
1. YAML frontmatter (title, generated date, status: draft, level: L1)
2. Acceptance Criteria with ID pattern: AC-{CTX}-{NNN}
3. Each AC should have: ID, Title, Given/When/Then format
4. Error cases for each AC
5. Traceability to source requirements

---DOCUMENT_SEPARATOR---

DOCUMENT 2: business-rules.md
1. YAML frontmatter (title, generated date, status: draft, level: L1)
2. Business Rules with ID pattern: BR-{CTX}-{NNN}
3. Each BR should have: ID, Title, Rule statement, Invariant, Enforcement
4. Traceability to source requirements`, string(dmJSON), string(decJSON), nfr, inputContent)

	response, err := client.Call(prompt)
	if err != nil {
		return "", "", err
	}

	// Split response
	parts := splitDocuments(response, "---DOCUMENT_SEPARATOR---")
	if len(parts) < 2 {
		return response, "", nil
	}

	return parts[0], parts[1], nil
}

// formatDecisions formats decisions as markdown.
func formatDecisions(decisions []domain.Decision) string {
	frontmatter := fmt.Sprintf(`---
title: "Design Decisions"
generated: %s
status: draft
level: L1
---

# Design Decisions

## Overview

This document records design decisions made during L0→L1 derivation.

## Decisions

`, time.Now().Format(time.RFC3339))

	var content string
	for i, d := range decisions {
		content += fmt.Sprintf(`### DEC-%03d: %s

**Question:** %s

**Answer:** %s

**Source:** %s

**Decided:** %s

---

`, i+1, d.Subject, d.Question, d.Answer, d.Source, d.DecidedAt.Format(time.RFC3339))
	}

	return frontmatter + content
}

// splitDocuments splits a response into multiple documents.
func splitDocuments(response, separator string) []string {
	parts := make([]string, 0)
	for _, part := range splitString(response, separator) {
		trimmed := trimString(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// splitString splits a string by separator.
func splitString(s, sep string) []string {
	var parts []string
	for {
		idx := indexOf(s, sep)
		if idx == -1 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	return parts
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// trimString trims whitespace from a string.
func trimString(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// interactiveApproval handles interactive approval for a file.
// Returns the approved content and any error.
func interactiveApproval(filename, content string, client *claude.Client, regenerate func() (string, error)) (string, error) {
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
			newContent, err := regenerate()
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
