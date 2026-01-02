package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/internal/workflow"
	"github.com/ikadar/loom-cli/prompts"
)

// TestCase represents a single test case
type TestCase struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Category        string     `json:"category"` // positive, negative, boundary, hallucination
	ACRef           string     `json:"ac_ref"`
	BRRefs          []string   `json:"br_refs"`
	Preconditions   []string   `json:"preconditions"`
	TestData        []TestData `json:"test_data"`
	Steps           []string   `json:"steps"`
	ExpectedResults []string   `json:"expected_results"`
	ShouldNot       string     `json:"should_not,omitempty"` // For hallucination prevention tests
}

// TestData represents test data for a test case
type TestData struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Notes string      `json:"notes"`
}

// TestSuite represents a suite of tests for one AC
type TestSuite struct {
	ACRef   string     `json:"ac_ref"`
	ACTitle string     `json:"ac_title"`
	Tests   []TestCase `json:"tests"`
}

// TDAISummary holds statistics about generated tests
type TDAISummary struct {
	Total      int `json:"total"`
	ByCategory struct {
		Positive      int `json:"positive"`
		Negative      int `json:"negative"`
		Boundary      int `json:"boundary"`
		Hallucination int `json:"hallucination"`
	} `json:"by_category"`
	Coverage struct {
		ACsCovered            int     `json:"acs_covered"`
		PositiveRatio         float64 `json:"positive_ratio"`
		NegativeRatio         float64 `json:"negative_ratio"`
		HasHallucinationTests bool    `json:"has_hallucination_tests"`
	} `json:"coverage"`
}

// TestCaseResult holds the result of test case generation
type TestCaseResult struct {
	TestSuites []TestSuite `json:"test_suites"`
	Summary    TDAISummary `json:"summary"`
}

// ChunkedTestCaseGenerator generates test cases in batches
type ChunkedTestCaseGenerator struct {
	Client    *claude.Client
	ChunkSize int // Number of ACs per chunk
}

// NewChunkedTestCaseGenerator creates a new generator with default settings
func NewChunkedTestCaseGenerator(client *claude.Client) *ChunkedTestCaseGenerator {
	return &ChunkedTestCaseGenerator{
		Client:    client,
		ChunkSize: 5, // 5 ACs per batch
	}
}

// Generate generates test cases from AC content in chunks
func (g *ChunkedTestCaseGenerator) Generate(acContent string) (*TestCaseResult, error) {
	// Parse ACs from content
	acSections := parseACsections(acContent)

	if len(acSections) == 0 {
		return nil, fmt.Errorf("no ACs found in content")
	}

	fmt.Fprintf(os.Stderr, "  Found %d ACs, processing in chunks of %d\n", len(acSections), g.ChunkSize)

	// Create progress tracker
	numChunks := (len(acSections) + g.ChunkSize - 1) / g.ChunkSize
	progress := workflow.NewProgress("Generating", numChunks)

	// Process in chunks
	var allSuites []TestSuite
	var errors []string
	retryCfg := claude.DefaultRetryConfig()

	for i := 0; i < len(acSections); i += g.ChunkSize {
		end := i + g.ChunkSize
		if end > len(acSections) {
			end = len(acSections)
		}
		chunk := acSections[i:end]
		chunkNum := i/g.ChunkSize + 1

		// Build chunk content
		chunkContent := buildChunkContent(chunk)

		// Generate for this chunk (inject document into <context> tag)
		prompt := buildPrompt(prompts.DeriveTestCases, chunkContent)

		var result struct {
			TestSuites []TestSuite `json:"test_suites"`
			Summary    TDAISummary `json:"summary"`
		}

		err := g.Client.CallJSONWithRetry(prompt, &result, retryCfg)
		if err != nil {
			errors = append(errors, fmt.Sprintf("chunk %d: %v", chunkNum, err))
			progress.Increment()
			continue
		}

		allSuites = append(allSuites, result.TestSuites...)
		progress.Increment()
	}

	progress.Done()

	// Report errors if any
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "  Warnings: %d chunk(s) failed:\n", len(errors))
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "    - %s\n", e)
		}
	}

	if len(allSuites) == 0 {
		return nil, fmt.Errorf("all chunks failed: %v", errors)
	}

	// Calculate combined summary
	summary := calculateSummary(allSuites)

	return &TestCaseResult{
		TestSuites: allSuites,
		Summary:    summary,
	}, nil
}

// parseACsections splits AC content into individual AC sections
func parseACsections(content string) []string {
	// Split by AC header pattern: ## AC-XXX-NNN
	pattern := regexp.MustCompile(`(?m)^## AC-[A-Z]+-\d+`)
	indices := pattern.FindAllStringIndex(content, -1)

	if len(indices) == 0 {
		// Fallback: try splitting by --- separator
		parts := strings.Split(content, "\n---\n")
		var filtered []string
		for _, p := range parts {
			if strings.Contains(p, "AC-") {
				filtered = append(filtered, strings.TrimSpace(p))
			}
		}
		return filtered
	}

	var sections []string
	for i, idx := range indices {
		start := idx[0]
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(content)
		}
		section := strings.TrimSpace(content[start:end])
		if section != "" {
			sections = append(sections, section)
		}
	}

	return sections
}

// buildChunkContent creates prompt content from AC sections
func buildChunkContent(sections []string) string {
	var sb strings.Builder
	sb.WriteString("# Acceptance Criteria\n\n")

	for i, section := range sections {
		if i > 0 {
			sb.WriteString("\n---\n\n")
		}
		sb.WriteString(section)
		sb.WriteString("\n")
	}

	return sb.String()
}

// calculateSummary aggregates statistics from all test suites
func calculateSummary(suites []TestSuite) TDAISummary {
	var summary TDAISummary

	for _, suite := range suites {
		for _, tc := range suite.Tests {
			summary.Total++
			switch tc.Category {
			case "positive":
				summary.ByCategory.Positive++
			case "negative":
				summary.ByCategory.Negative++
			case "boundary":
				summary.ByCategory.Boundary++
			case "hallucination":
				summary.ByCategory.Hallucination++
			}
		}
	}

	summary.Coverage.ACsCovered = len(suites)
	if summary.Total > 0 {
		summary.Coverage.PositiveRatio = float64(summary.ByCategory.Positive) / float64(summary.Total)
		summary.Coverage.NegativeRatio = float64(summary.ByCategory.Negative) / float64(summary.Total)
	}
	summary.Coverage.HasHallucinationTests = summary.ByCategory.Hallucination > 0

	return summary
}

// FlattenTestCases extracts all test cases from suites
func FlattenTestCases(suites []TestSuite) []TestCase {
	var all []TestCase
	for _, suite := range suites {
		all = append(all, suite.Tests...)
	}
	return all
}

// buildPrompt injects document content into the <context> tag of an XML-structured prompt.
func buildPrompt(promptTemplate, document string) string {
	return strings.Replace(promptTemplate, "</context>", document+"\n</context>", 1)
}
