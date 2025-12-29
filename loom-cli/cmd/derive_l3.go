package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/prompts"
)

// L3Result is the output of the derive-l3 command
type L3Result struct {
	APISpec                 APISpec                  `json:"api_spec"`
	ImplementationSkeletons []ImplementationSkeleton `json:"implementation_skeletons"`
	Summary                 L3Summary                `json:"summary"`
}

type L3Summary struct {
	EndpointsCount int `json:"endpoints_count"`
	ServicesCount  int `json:"services_count"`
	FunctionsCount int `json:"functions_count"`
}

type APISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    map[string]interface{} `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
}

type ImplementationSkeleton struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Functions    []FunctionSpec `json:"functions"`
	Dependencies []string       `json:"dependencies"`
	RelatedSpecs []string       `json:"related_specs"`
}

type FunctionSpec struct {
	Name        string   `json:"name"`
	Signature   string   `json:"signature"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	ErrorCases  []string `json:"error_cases"`
}

func runDeriveL3() error {
	// Parse arguments
	args := os.Args[2:]

	var inputDir string
	var outputDir string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input-dir":
			if i+1 < len(args) {
				i++
				inputDir = args[i]
			}
		case "--output-dir":
			if i+1 < len(args) {
				i++
				outputDir = args[i]
			}
		}
	}

	if inputDir == "" {
		return fmt.Errorf("--input-dir is required (directory containing L2 documents)")
	}
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required")
	}

	// Read L2 documents
	fmt.Fprintln(os.Stderr, "Phase L3-0: Reading L2 documents...")

	tcContent, err := os.ReadFile(filepath.Join(inputDir, "test-cases.md"))
	if err != nil {
		return fmt.Errorf("failed to read test-cases.md: %w", err)
	}

	tsContent, err := os.ReadFile(filepath.Join(inputDir, "tech-specs.md"))
	if err != nil {
		return fmt.Errorf("failed to read tech-specs.md: %w", err)
	}

	l2Input := fmt.Sprintf("## Test Cases\n\n%s\n\n## Technical Specifications\n\n%s",
		string(tcContent), string(tsContent))

	fmt.Fprintf(os.Stderr, "  Read: test-cases.md (%d bytes)\n", len(tcContent))
	fmt.Fprintf(os.Stderr, "  Read: tech-specs.md (%d bytes)\n", len(tsContent))

	// Create Claude client
	client := claude.NewClient()

	// Generate L3 documents
	fmt.Fprintln(os.Stderr, "\nPhase L3-1: Generating API Spec and Implementation Skeletons...")

	prompt := prompts.DeriveL3 + l2Input

	var result L3Result
	if err := client.CallJSON(prompt, &result); err != nil {
		return fmt.Errorf("failed to generate L3 documents: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Generated: %d endpoints, %d services\n",
		len(result.APISpec.Paths), len(result.ImplementationSkeletons))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write API Spec (OpenAPI JSON)
	fmt.Fprintln(os.Stderr, "\nPhase L3-2: Writing output...")

	apiPath := filepath.Join(outputDir, "openapi.json")
	apiContent, _ := json.MarshalIndent(result.APISpec, "", "  ")
	if err := os.WriteFile(apiPath, apiContent, 0644); err != nil {
		return fmt.Errorf("failed to write API spec: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", apiPath)

	// Write Implementation Skeletons
	implPath := filepath.Join(outputDir, "implementation-skeletons.md")
	if err := writeSkeletons(implPath, result.ImplementationSkeletons); err != nil {
		return fmt.Errorf("failed to write skeletons: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", implPath)

	// Write full JSON for further processing
	jsonPath := filepath.Join(outputDir, "l3-output.json")
	jsonContent, _ := json.MarshalIndent(result, "", "  ")
	if err := os.WriteFile(jsonPath, jsonContent, 0644); err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", jsonPath)

	// Print summary
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "        L3 DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Generated:")
	fmt.Fprintf(os.Stderr, "  API Endpoints: %d\n", len(result.APISpec.Paths))
	fmt.Fprintf(os.Stderr, "  Services:      %d\n", len(result.ImplementationSkeletons))

	funcCount := 0
	for _, skel := range result.ImplementationSkeletons {
		funcCount += len(skel.Functions)
	}
	fmt.Fprintf(os.Stderr, "  Functions:     %d\n", funcCount)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Output:")
	fmt.Fprintf(os.Stderr, "  %s\n", apiPath)
	fmt.Fprintf(os.Stderr, "  %s\n", implPath)

	return nil
}

func writeSkeletons(path string, skeletons []ImplementationSkeleton) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Implementation Skeletons\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, skel := range skeletons {
		fmt.Fprintf(f, "## %s (%s)\n\n", skel.Name, skel.Type)

		if len(skel.Dependencies) > 0 {
			fmt.Fprintf(f, "**Dependencies:** %v\n\n", skel.Dependencies)
		}

		fmt.Fprintf(f, "**Related Specs:** %v\n\n", skel.RelatedSpecs)

		fmt.Fprintf(f, "### Functions\n\n")
		for _, fn := range skel.Functions {
			fmt.Fprintf(f, "#### `%s`\n\n", fn.Name)
			fmt.Fprintf(f, "```typescript\n%s\n```\n\n", fn.Signature)
			fmt.Fprintf(f, "%s\n\n", fn.Description)

			if len(fn.Steps) > 0 {
				fmt.Fprintf(f, "**Implementation Steps:**\n")
				for i, step := range fn.Steps {
					fmt.Fprintf(f, "%d. %s\n", i+1, step)
				}
				fmt.Fprintf(f, "\n")
			}

			if len(fn.ErrorCases) > 0 {
				fmt.Fprintf(f, "**Error Cases:** %v\n\n", fn.ErrorCases)
			}
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}
