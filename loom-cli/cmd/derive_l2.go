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

// L2Result is the output of the derive-l2 command
type L2Result struct {
	Summary    L2Summary    `json:"summary"`
	TestCases  []TestCase   `json:"test_cases"`
	TechSpecs  []TechSpec   `json:"tech_specs"`
}

type L2Summary struct {
	TestCasesGenerated int         `json:"test_cases_generated"`
	TechSpecsGenerated int         `json:"tech_specs_generated"`
	Coverage           L2Coverage  `json:"coverage"`
}

type L2Coverage struct {
	ACsCovered      int `json:"acs_covered"`
	BRsCovered      int `json:"brs_covered"`
	HappyPathTests  int `json:"happy_path_tests"`
	ErrorTests      int `json:"error_tests"`
	EdgeCaseTests   int `json:"edge_case_tests"`
}

type TestCase struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	ACRef           string     `json:"ac_ref"`
	BRRefs          []string   `json:"br_refs"`
	Preconditions   []string   `json:"preconditions"`
	TestData        []TestData `json:"test_data"`
	Steps           []string   `json:"steps"`
	ExpectedResults []string   `json:"expected_results"`
}

type TestData struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Notes string      `json:"notes"`
}

type TechSpec struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	BRRef            string          `json:"br_ref"`
	Rule             string          `json:"rule"`
	Implementation   string          `json:"implementation"`
	ValidationPoints []string        `json:"validation_points"`
	DataRequirements []DataReq       `json:"data_requirements"`
	ErrorHandling    []ErrorHandling `json:"error_handling"`
	RelatedACs       []string        `json:"related_acs"`
}

type DataReq struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	Source      string `json:"source"`
}

type ErrorHandling struct {
	Condition  string `json:"condition"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status"`
}

func runDeriveL2() error {
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
		return fmt.Errorf("--input-dir is required (directory containing L1 documents)")
	}
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required")
	}

	// Read L1 documents
	fmt.Fprintln(os.Stderr, "Phase L2-0: Reading L1 documents...")

	acContent, err := os.ReadFile(filepath.Join(inputDir, "acceptance-criteria.md"))
	if err != nil {
		return fmt.Errorf("failed to read acceptance-criteria.md: %w", err)
	}

	brContent, err := os.ReadFile(filepath.Join(inputDir, "business-rules.md"))
	if err != nil {
		return fmt.Errorf("failed to read business-rules.md: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Read: acceptance-criteria.md (%d bytes)\n", len(acContent))
	fmt.Fprintf(os.Stderr, "  Read: business-rules.md (%d bytes)\n", len(brContent))

	// Create Claude client
	client := claude.NewClient()

	// Phase 1: Generate Test Cases from ACs
	fmt.Fprintln(os.Stderr, "\nPhase L2-1: Generating Test Cases from Acceptance Criteria...")

	tcPrompt := prompts.DeriveTestCases + "\n" + string(acContent)

	var tcResult struct {
		TestCases []TestCase `json:"test_cases"`
		Summary   struct {
			Total     int `json:"total"`
			HappyPath int `json:"happy_path"`
			ErrorCase int `json:"error_case"`
			EdgeCase  int `json:"edge_case"`
		} `json:"summary"`
	}
	if err := client.CallJSON(tcPrompt, &tcResult); err != nil {
		return fmt.Errorf("failed to generate test cases: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Test Cases\n", len(tcResult.TestCases))

	// Phase 2: Generate Tech Specs from BRs
	fmt.Fprintln(os.Stderr, "\nPhase L2-2: Generating Tech Specs from Business Rules...")

	tsPrompt := prompts.DeriveTechSpecs + "\n" + string(brContent)

	var tsResult struct {
		TechSpecs []TechSpec `json:"tech_specs"`
		Summary   struct {
			Total int `json:"total"`
		} `json:"summary"`
	}
	if err := client.CallJSON(tsPrompt, &tsResult); err != nil {
		return fmt.Errorf("failed to generate tech specs: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Tech Specs\n", len(tsResult.TechSpecs))

	// Combine results
	var result L2Result
	result.TestCases = tcResult.TestCases
	result.TechSpecs = tsResult.TechSpecs
	result.Summary = L2Summary{
		TestCasesGenerated: len(tcResult.TestCases),
		TechSpecsGenerated: len(tsResult.TechSpecs),
		Coverage: L2Coverage{
			ACsCovered:     tcResult.Summary.Total,
			BRsCovered:     tsResult.Summary.Total,
			HappyPathTests: tcResult.Summary.HappyPath,
			ErrorTests:     tcResult.Summary.ErrorCase,
			EdgeCaseTests:  tcResult.Summary.EdgeCase,
		},
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write Test Cases
	fmt.Fprintln(os.Stderr, "\nPhase L2-3: Writing output...")

	tcPath := filepath.Join(outputDir, "test-cases.md")
	if err := writeTestCases(tcPath, result.TestCases); err != nil {
		return fmt.Errorf("failed to write test cases: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", tcPath)

	// Write Tech Specs
	tsPath := filepath.Join(outputDir, "tech-specs.md")
	if err := writeTechSpecs(tsPath, result.TechSpecs); err != nil {
		return fmt.Errorf("failed to write tech specs: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", tsPath)

	// Write JSON for further processing
	jsonPath := filepath.Join(outputDir, "l2-output.json")
	jsonContent, _ := json.MarshalIndent(result, "", "  ")
	if err := os.WriteFile(jsonPath, jsonContent, 0644); err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", jsonPath)

	// Print summary
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "        L2 DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Generated:")
	fmt.Fprintf(os.Stderr, "  Test Cases:  %d\n", len(result.TestCases))
	fmt.Fprintf(os.Stderr, "  Tech Specs:  %d\n", len(result.TechSpecs))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Coverage:")
	fmt.Fprintf(os.Stderr, "  ACs covered:      %d\n", result.Summary.Coverage.ACsCovered)
	fmt.Fprintf(os.Stderr, "  BRs covered:      %d\n", result.Summary.Coverage.BRsCovered)
	fmt.Fprintf(os.Stderr, "  Happy path tests: %d\n", result.Summary.Coverage.HappyPathTests)
	fmt.Fprintf(os.Stderr, "  Error tests:      %d\n", result.Summary.Coverage.ErrorTests)
	fmt.Fprintf(os.Stderr, "  Edge case tests:  %d\n", result.Summary.Coverage.EdgeCaseTests)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Output:")
	fmt.Fprintf(os.Stderr, "  %s\n", tcPath)
	fmt.Fprintf(os.Stderr, "  %s\n", tsPath)

	return nil
}

func writeTestCases(path string, testCases []TestCase) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Test Cases\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, tc := range testCases {
		fmt.Fprintf(f, "## %s – %s\n\n", tc.ID, tc.Name)
		fmt.Fprintf(f, "**Type:** %s\n\n", tc.Type)

		fmt.Fprintf(f, "**Preconditions:**\n")
		for _, p := range tc.Preconditions {
			fmt.Fprintf(f, "- %s\n", p)
		}
		fmt.Fprintf(f, "\n")

		if len(tc.TestData) > 0 {
			fmt.Fprintf(f, "**Test Data:**\n")
			fmt.Fprintf(f, "| Field | Value | Notes |\n")
			fmt.Fprintf(f, "|-------|-------|-------|\n")
			for _, td := range tc.TestData {
				fmt.Fprintf(f, "| %s | %v | %s |\n", td.Field, td.Value, td.Notes)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "**Steps:**\n")
		for i, s := range tc.Steps {
			fmt.Fprintf(f, "%d. %s\n", i+1, s)
		}
		fmt.Fprintf(f, "\n")

		fmt.Fprintf(f, "**Expected Result:**\n")
		for _, r := range tc.ExpectedResults {
			fmt.Fprintf(f, "- %s\n", r)
		}
		fmt.Fprintf(f, "\n")

		fmt.Fprintf(f, "**Traceability:**\n")
		fmt.Fprintf(f, "- AC: %s\n", tc.ACRef)
		if len(tc.BRRefs) > 0 {
			fmt.Fprintf(f, "- BR: %v\n", tc.BRRefs)
		}
		fmt.Fprintf(f, "\n---\n\n")
	}

	return nil
}

func writeTechSpecs(path string, techSpecs []TechSpec) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Technical Specifications\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, ts := range techSpecs {
		fmt.Fprintf(f, "## %s – %s\n\n", ts.ID, ts.Name)
		fmt.Fprintf(f, "**Rule:** %s\n\n", ts.Rule)
		fmt.Fprintf(f, "**Implementation Approach:**\n%s\n\n", ts.Implementation)

		if len(ts.ValidationPoints) > 0 {
			fmt.Fprintf(f, "**Validation Points:**\n")
			for _, vp := range ts.ValidationPoints {
				fmt.Fprintf(f, "- %s\n", vp)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(ts.DataRequirements) > 0 {
			fmt.Fprintf(f, "**Data Requirements:**\n")
			fmt.Fprintf(f, "| Field | Type | Constraints | Source |\n")
			fmt.Fprintf(f, "|-------|------|-------------|--------|\n")
			for _, dr := range ts.DataRequirements {
				fmt.Fprintf(f, "| %s | %s | %s | %s |\n", dr.Field, dr.Type, dr.Constraints, dr.Source)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(ts.ErrorHandling) > 0 {
			fmt.Fprintf(f, "**Error Handling:**\n")
			fmt.Fprintf(f, "| Condition | Error Code | Message | HTTP Status |\n")
			fmt.Fprintf(f, "|-----------|------------|---------|-------------|\n")
			for _, eh := range ts.ErrorHandling {
				fmt.Fprintf(f, "| %s | %s | %s | %d |\n", eh.Condition, eh.ErrorCode, eh.Message, eh.HTTPStatus)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "**Traceability:**\n")
		fmt.Fprintf(f, "- BR: %s\n", ts.BRRef)
		if len(ts.RelatedACs) > 0 {
			fmt.Fprintf(f, "- Related ACs: %v\n", ts.RelatedACs)
		}
		fmt.Fprintf(f, "\n---\n\n")
	}

	return nil
}
