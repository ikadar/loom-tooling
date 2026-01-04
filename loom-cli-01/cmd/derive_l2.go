// Package cmd provides CLI commands for loom-cli.
//
// This file implements the derive-l2 command.
// Implements: l2/interface-contracts.md IC-DRV-002
// See: l2/sequence-design.md SEQ-DRV-001
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

// runDeriveL2 implements the derive-l2 command.
// Implements: IC-DRV-002
//
// Output files:
//   - tech-specs.md
//   - interface-contracts.md
//   - aggregate-design.md
//   - sequence-design.md
//   - initial-data-model.md
func runDeriveL2(args []string) int {
	fs := flag.NewFlagSet("derive-l2", flag.ExitOnError)

	inputDir := fs.String("input-dir", "./l1", "L1 documents directory")
	outputDir := fs.String("output-dir", "./l2", "Output directory")
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

	// Read L1 documents
	l1Content, err := readL1Documents(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading L1 documents: %v\n", err)
		return domain.ExitCodeError
	}

	client := claude.NewClient()
	client.Verbose = Verbose()
	_ = interactive // TODO: implement interactive mode

	// Phase L2-1: Tech Specs
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase L2-1: Generating tech-specs.md...")
	}
	if err := generateL2Document(client, prompts.DeriveTechSpecs, l1Content, *outputDir, "tech-specs.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L2-2: Interface Contracts
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase L2-2: Generating interface-contracts.md...")
	}
	if err := generateL2Document(client, prompts.DeriveInterfaceContracts, l1Content, *outputDir, "interface-contracts.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L2-3: Aggregate Design
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase L2-3: Generating aggregate-design.md...")
	}
	if err := generateL2Document(client, prompts.DeriveAggregateDesign, l1Content, *outputDir, "aggregate-design.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L2-4: Sequence Design
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase L2-4: Generating sequence-design.md...")
	}
	if err := generateL2Document(client, prompts.DeriveSequenceDesign, l1Content, *outputDir, "sequence-design.md"); err != nil {
		return domain.ExitCodeError
	}

	// Phase L2-5: Data Model
	if Verbose() {
		fmt.Fprintln(os.Stderr, "[derive-l2] Phase L2-5: Generating initial-data-model.md...")
	}
	if err := generateL2Document(client, prompts.DeriveDataModel, l1Content, *outputDir, "initial-data-model.md"); err != nil {
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "[derive-l2] L2 derivation complete. Output: %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// readL1Documents reads all markdown files from L1 directory.
func readL1Documents(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var content string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return "", err
		}

		if content != "" {
			content += "\n\n---\n\n"
		}
		content += fmt.Sprintf("# File: %s\n\n%s", entry.Name(), string(data))
	}

	return content, nil
}

// generateL2Document generates an L2 document and writes it to file.
func generateL2Document(client *claude.Client, prompt, context, outputDir, filename string) error {
	fullPrompt := claude.BuildPrompt(prompt, context)

	response, err := client.Call(fullPrompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating %s: %v\n", filename, err)
		return err
	}

	// Try to extract markdown content from JSON response
	content := extractMarkdownFromResponse(response, filename)

	// Add frontmatter if not present
	if len(content) < 3 || content[:3] != "---" {
		frontmatter := fmt.Sprintf(`---
title: "%s"
generated: %s
status: draft
level: L2
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
		fmt.Fprintf(os.Stderr, "[derive-l2] Wrote %s\n", path)
	}

	return nil
}

// extractMarkdownFromResponse extracts markdown content from a JSON response.
func extractMarkdownFromResponse(response, filename string) string {
	// Try to parse as JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Not JSON, return as-is (likely already markdown)
		return response
	}

	// Convert JSON to markdown format
	return jsonToMarkdown(result, filename)
}

// jsonToMarkdown converts JSON response to markdown format.
func jsonToMarkdown(data map[string]interface{}, filename string) string {
	md := fmt.Sprintf("# %s\n\n", filenameToTitle(filename))

	// Generic JSON to markdown conversion
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	md += "```json\n" + string(prettyJSON) + "\n```\n"

	return md
}

// filenameToTitle converts filename to title.
func filenameToTitle(filename string) string {
	name := filepath.Base(filename)
	name = name[:len(name)-len(filepath.Ext(name))]

	// Convert kebab-case to Title Case
	result := ""
	capitalize := true
	for _, r := range name {
		if r == '-' {
			result += " "
			capitalize = true
		} else {
			if capitalize {
				result += string(r - 32) // uppercase
				capitalize = false
			} else {
				result += string(r)
			}
		}
	}

	return result
}
