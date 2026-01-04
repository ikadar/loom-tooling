// Package cmd provides CLI commands for loom-cli.
//
// This file implements the derive-l3 command.
// Implements: l2/interface-contracts.md IC-DRV-003
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

// runDeriveL3 implements the derive-l3 command.
// Implements: IC-DRV-003
//
// Output files:
//   - test-cases.md
//   - openapi-spec.md
//   - implementation-skeletons.md
//   - feature-tickets.md
//   - service-boundaries.md
//   - event-message-design.md
//   - dependency-graph.md
func runDeriveL3(args []string) int {
	fs := flag.NewFlagSet("derive-l3", flag.ExitOnError)

	inputDir := fs.String("input-dir", "./l2", "L2 documents directory")
	l1Dir := fs.String("l1-dir", "./l1", "L1 documents directory (for AC reference)")
	outputDir := fs.String("output-dir", "./l3", "Output directory")
	interactive := fs.Bool("interactive", false, "Enable interactive approval")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Read L1 and L2 documents
	l1Content, _ := readL1Documents(*l1Dir)
	l2Content, err := readL1Documents(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading L2 documents: %v\n", err)
		return domain.ExitCodeError
	}

	context := l1Content + "\n\n---\n\n" + l2Content

	client := claude.NewClient()
	client.Verbose = Verbose()
	_ = interactive // TODO: implement interactive mode

	// Phase L3-1: Test Cases (TDAI)
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-1: Generating test-cases.md...")
	}
	if err := generateL3Document(client, prompts.DeriveTestCases, l1Content, *outputDir, "test-cases.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-2a: OpenAPI Spec
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-2a: Generating openapi-spec.md...")
	}
	if err := generateL3Document(client, prompts.DeriveL3API, context, *outputDir, "openapi-spec.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-2b: Implementation Skeletons
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-2b: Generating implementation-skeletons.md...")
	}
	if err := generateL3Document(client, prompts.DeriveL3Skeletons, l2Content, *outputDir, "implementation-skeletons.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-3: Feature Tickets
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-3: Generating feature-tickets.md...")
	}
	if err := generateL3Document(client, prompts.DeriveFeatureTickets, context, *outputDir, "feature-tickets.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-4: Service Boundaries
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-4: Generating service-boundaries.md...")
	}
	if err := generateL3Document(client, prompts.DeriveServiceBoundaries, context, *outputDir, "service-boundaries.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-5: Event/Message Design
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-5: Generating event-message-design.md...")
	}
	if err := generateL3Document(client, prompts.DeriveEventDesign, context, *outputDir, "event-message-design.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L3-6: Dependency Graph
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l3] Phase L3-6: Generating dependency-graph.md...")
	}
	if err := generateL3Document(client, prompts.DeriveDependencyGraph, context, *outputDir, "dependency-graph.md"); err != nil {
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive-l3] L3 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// generateL3Document generates an L3 document and writes it to file.
func generateL3Document(client *claude.Client, prompt, context, outputDir, filename string) error {
	fullPrompt := claude.BuildPrompt(prompt, context)

	response, err := client.Call(fullPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", filename, err)
		return err
	}

	// Try to extract markdown content from JSON response
	content := extractL3Content(response, filename)

	// Add frontmatter if not present
	if len(content) < 3 || content[:3] != "---" {
		frontmatter := fmt.Sprintf(`---
title: "%s"
generated: %s
status: draft
level: L3
---

`, filenameToTitle(filename), time.Now().Format(time.RFC3339))
		content = frontmatter + content
	}

	path := filepath.Join(outputDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filename, err)
		return err
	}

	if Verbose() {
		fmt.Fprintf(os.Stderr, "[derive-l3] Wrote %s\n", path)
	}

	return nil
}

// extractL3Content extracts content from response.
func extractL3Content(response, filename string) string {
	// Try to parse as JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Not JSON, return as-is
		return response
	}

	// Special handling for test-cases
	if filename == "test-cases.md" {
		return generateTestCasesMarkdown(result)
	}

	// Default: convert JSON to markdown
	return jsonToMarkdown(result, filename)
}

// generateTestCasesMarkdown generates test-cases.md from JSON.
func generateTestCasesMarkdown(data map[string]interface{}) string {
	md := `# TDAI Test Cases

**Methodology:** Test-Driven AI Development (TDAI)

`

	// Summary section
	if summary, ok := data["summary"].(map[string]interface{}); ok {
		md += "## Summary\n\n"
		md += "| Category | Count |\n"
		md += "|----------|-------|\n"
		if byCat, ok := summary["by_category"].(map[string]interface{}); ok {
			for cat, count := range byCat {
				md += fmt.Sprintf("| %s | %v |\n", cat, count)
			}
		}
		md += "\n"
	}

	// Test suites
	if suites, ok := data["test_suites"].([]interface{}); ok {
		for _, suite := range suites {
			s, ok := suite.(map[string]interface{})
			if !ok {
				continue
			}
			md += fmt.Sprintf("## %s: %s\n\n", s["ac_ref"], s["ac_title"])

			if tests, ok := s["tests"].([]interface{}); ok {
				for _, test := range tests {
					t, ok := test.(map[string]interface{})
					if !ok {
						continue
					}
					md += fmt.Sprintf("### %s: %s\n\n", t["id"], t["name"])
					md += fmt.Sprintf("**Category:** %s\n\n", t["category"])
					if sq, ok := t["source_quote"].(string); ok && sq != "" {
						md += fmt.Sprintf("**Source Quote:** %s\n\n", sq)
					}
					if shouldNot, ok := t["should_not"].(string); ok && shouldNot != "" {
						md += fmt.Sprintf("**Should NOT:** %s\n\n", shouldNot)
					}
					md += "---\n\n"
				}
			}
		}
	}

	return md
}
