// Package generator provides document generation orchestration.
//
// Implements: l2/package-structure.md PKG-006
// See: l2/internal-api.md
package generator

import (
	"encoding/json"
	"regexp"
	"strings"

	"loom-cli/internal/claude"
	"loom-cli/internal/formatter"
)

// Default chunk size for AC batching
const DefaultChunkSize = 5

// AC represents a parsed acceptance criterion for chunking.
type AC struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Given string `json:"given"`
	When  string `json:"when"`
	Then  string `json:"then"`
}

// TestSuite groups test cases by AC.
type TestSuite struct {
	ACRef   string               `json:"ac_ref"`
	ACTitle string               `json:"ac_title"`
	Tests   []formatter.TestCase `json:"tests"`
}

// TestCaseResult is the result of test case generation.
type TestCaseResult struct {
	TestSuites []TestSuite          `json:"test_suites"`
	Summary    formatter.TDAISummary `json:"summary"`
}

// ChunkedTestCaseGenerator generates test cases in batches.
//
// Implements: l2/internal-api.md
type ChunkedTestCaseGenerator struct {
	Client    *claude.Client
	ChunkSize int // Default: 5 ACs per chunk
}

// NewChunkedTestCaseGenerator creates generator with default settings.
//
// Implements: l2/internal-api.md
func NewChunkedTestCaseGenerator(client *claude.Client) *ChunkedTestCaseGenerator {
	return &ChunkedTestCaseGenerator{
		Client:    client,
		ChunkSize: DefaultChunkSize,
	}
}

// Generate generates test cases from AC markdown content.
//
// Strategy:
// 1. Parse ACs from markdown
// 2. Chunk into groups of ChunkSize
// 3. Generate test cases per chunk (parallel if possible)
// 4. Merge results
//
// Implements: l2/internal-api.md
func (g *ChunkedTestCaseGenerator) Generate(acContent string) (*TestCaseResult, error) {
	// Parse ACs from markdown
	acs := ParseACs(acContent)

	if len(acs) == 0 {
		return &TestCaseResult{
			TestSuites: []TestSuite{},
			Summary:    formatter.TDAISummary{},
		}, nil
	}

	// Chunk ACs
	chunks := ChunkACs(acs, g.ChunkSize)

	// Generate test cases for each chunk
	var allSuites []TestSuite
	for _, chunk := range chunks {
		suites, err := g.generateForChunk(chunk)
		if err != nil {
			return nil, err
		}
		allSuites = append(allSuites, suites...)
	}

	// Build summary
	summary := buildSummary(allSuites)

	return &TestCaseResult{
		TestSuites: allSuites,
		Summary:    summary,
	}, nil
}

// generateForChunk generates test cases for a chunk of ACs.
func (g *ChunkedTestCaseGenerator) generateForChunk(acs []AC) ([]TestSuite, error) {
	// Build prompt with ACs
	acJSON, _ := json.MarshalIndent(acs, "", "  ")
	prompt := claude.BuildPrompt("Generate test cases for these ACs", string(acJSON))

	var result struct {
		Suites []TestSuite `json:"test_suites"`
	}

	if err := g.Client.CallJSON(prompt, &result); err != nil {
		return nil, err
	}

	return result.Suites, nil
}

// ParseACs extracts acceptance criteria from markdown content.
func ParseACs(content string) []AC {
	var acs []AC

	// Pattern to match AC sections: ### AC-XXX-NNN: Title
	acPattern := regexp.MustCompile(`(?m)^###\s+(AC-[A-Z0-9-]+):\s*(.+)$`)
	matches := acPattern.FindAllStringSubmatchIndex(content, -1)

	for i, match := range matches {
		if len(match) < 6 {
			continue
		}

		id := content[match[2]:match[3]]
		title := content[match[4]:match[5]]

		// Get content between this match and next (or end)
		start := match[1]
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		section := content[start:end]

		// Parse Given/When/Then
		ac := AC{
			ID:    strings.TrimSpace(id),
			Title: strings.TrimSpace(title),
			Given: extractGWT(section, "Given"),
			When:  extractGWT(section, "When"),
			Then:  extractGWT(section, "Then"),
		}
		acs = append(acs, ac)
	}

	return acs
}

// extractGWT extracts Given/When/Then section from AC content.
func extractGWT(content, keyword string) string {
	pattern := regexp.MustCompile(`(?im)^\*\*` + keyword + `:\*\*\s*(.+)$`)
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// ChunkACs splits ACs into groups of specified size.
//
// Implements: l2/internal-api.md
func ChunkACs(acs []AC, chunkSize int) [][]AC {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	var chunks [][]AC
	for i := 0; i < len(acs); i += chunkSize {
		end := i + chunkSize
		if end > len(acs) {
			end = len(acs)
		}
		chunks = append(chunks, acs[i:end])
	}

	return chunks
}

// FlattenTestCases extracts all test cases from suites.
//
// Implements: l2/internal-api.md
func FlattenTestCases(suites []TestSuite) []formatter.TestCase {
	var all []formatter.TestCase
	for _, suite := range suites {
		all = append(all, suite.Tests...)
	}
	return all
}

// buildSummary creates a summary from test suites.
func buildSummary(suites []TestSuite) formatter.TDAISummary {
	var summary formatter.TDAISummary

	testCases := FlattenTestCases(suites)
	summary.Total = len(testCases)
	summary.Coverage.ACsCovered = len(suites)

	for _, tc := range testCases {
		switch tc.Category {
		case "positive":
			summary.ByCategory.Positive++
		case "negative":
			summary.ByCategory.Negative++
		case "boundary":
			summary.ByCategory.Boundary++
		case "hallucination":
			summary.ByCategory.Hallucination++
			summary.Coverage.HasHallucinationTests = true
		}
	}

	if summary.Total > 0 {
		summary.Coverage.PositiveRatio = float64(summary.ByCategory.Positive) / float64(summary.Total)
		summary.Coverage.NegativeRatio = float64(summary.ByCategory.Negative) / float64(summary.Total)
	}

	return summary
}
