// Package generator provides document generation utilities.
//
// Implements: l2/package-structure.md PKG-006
// See: l2/internal-api.md (generator section)
package generator

import (
	"loom-cli/internal/claude"
	"loom-cli/internal/formatter"
	"loom-cli/prompts"
)

// ChunkedTestCaseGenerator generates test cases in batches.
//
// Implements: l2/internal-api.md
type ChunkedTestCaseGenerator struct {
	Client    *claude.Client
	ChunkSize int // Default: 5 ACs per chunk
}

// TestCaseResult holds the complete test case generation result.
type TestCaseResult struct {
	TestSuites []formatter.TestSuite `json:"test_suites"`
	Summary    formatter.TDAISummary `json:"summary"`
}

// NewChunkedTestCaseGenerator creates a generator with default settings.
func NewChunkedTestCaseGenerator(client *claude.Client) *ChunkedTestCaseGenerator {
	return &ChunkedTestCaseGenerator{
		Client:    client,
		ChunkSize: 5,
	}
}

// Generate generates test cases from AC content.
func (g *ChunkedTestCaseGenerator) Generate(acContent string) (*TestCaseResult, error) {
	prompt := claude.BuildPrompt(prompts.DeriveTestCases, acContent)

	var result TestCaseResult
	if err := g.Client.CallJSONWithRetry(prompt, &result, claude.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	// Calculate summary if not provided
	if result.Summary.Total == 0 {
		result.Summary = calculateSummary(result.TestSuites)
	}

	return &result, nil
}

// calculateSummary calculates test case statistics.
func calculateSummary(suites []formatter.TestSuite) formatter.TDAISummary {
	var summary formatter.TDAISummary

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
		summary.Coverage.ACsCovered++
	}

	if summary.Total > 0 {
		summary.Coverage.PositiveRatio = float64(summary.ByCategory.Positive) / float64(summary.Total)
		summary.Coverage.NegativeRatio = float64(summary.ByCategory.Negative) / float64(summary.Total)
		summary.Coverage.HasHallucinationTests = summary.ByCategory.Hallucination > 0
	}

	return summary
}

// FlattenTestCases extracts all test cases from suites.
func FlattenTestCases(suites []formatter.TestSuite) []formatter.TestCase {
	var result []formatter.TestCase
	for _, suite := range suites {
		result = append(result, suite.Tests...)
	}
	return result
}
