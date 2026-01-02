package formatter

import (
	"fmt"
	"strings"
)

// L1BasePath is the relative path to L1 documents from L2
const L1BasePath = "../l1"

// FormatTestCases formats test cases as markdown
func FormatTestCases(testCases []TestCase, summary TDAISummary, timestamp string) string {
	var sb strings.Builder

	// Header with frontmatter (test-cases is L3 per documentation-derivation-strategy.md)
	fm := DefaultFrontmatter("TDAI Test Cases", timestamp, "L3")
	sb.WriteString(FormatHeaderWithFrontmatter(fm))
	sb.WriteString("**Methodology:** Test-Driven AI Development (TDAI)\n\n")

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Category | Count | Ratio |\n")
	sb.WriteString("|----------|-------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Positive | %d | %.1f%% |\n", summary.ByCategory.Positive, summary.Coverage.PositiveRatio*100))
	sb.WriteString(fmt.Sprintf("| Negative | %d | %.1f%% |\n", summary.ByCategory.Negative, summary.Coverage.NegativeRatio*100))
	sb.WriteString(fmt.Sprintf("| Boundary | %d | - |\n", summary.ByCategory.Boundary))
	sb.WriteString(fmt.Sprintf("| Hallucination Prevention | %d | - |\n", summary.ByCategory.Hallucination))
	sb.WriteString(fmt.Sprintf("| **Total** | **%d** | - |\n\n", summary.Total))

	sb.WriteString(fmt.Sprintf("**Coverage:** %d ACs covered\n", summary.Coverage.ACsCovered))
	if summary.Coverage.HasHallucinationTests {
		sb.WriteString("**Hallucination Prevention:** ✓ Enabled\n")
	}
	sb.WriteString("\n---\n\n")

	// Group tests by category
	categories := []string{"positive", "negative", "boundary", "hallucination"}
	categoryNames := map[string]string{
		"positive":      "Positive Tests (Happy Path)",
		"negative":      "Negative Tests (Error Cases)",
		"boundary":      "Boundary Tests",
		"hallucination": "Hallucination Prevention Tests",
	}

	for _, cat := range categories {
		var catTests []TestCase
		for _, tc := range testCases {
			if tc.Category == cat {
				catTests = append(catTests, tc)
			}
		}

		if len(catTests) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", categoryNames[cat]))

		for _, tc := range catTests {
			sb.WriteString(formatTestCase(tc))
		}
	}

	return sb.String()
}

// formatTestCase formats a single test case
func formatTestCase(tc TestCase) string {
	var sb strings.Builder

	// Header with anchor
	sb.WriteString(FormatSectionHeader(3, tc.ID, tc.Name))

	// Special handling for hallucination tests
	if tc.Category == "hallucination" && tc.ShouldNot != "" {
		sb.WriteString(fmt.Sprintf("**⚠️ Should NOT:** %s\n\n", tc.ShouldNot))
	}

	// Preconditions
	sb.WriteString("**Preconditions:**\n")
	for _, p := range tc.Preconditions {
		sb.WriteString(fmt.Sprintf("- %s\n", p))
	}
	sb.WriteString("\n")

	// Test Data
	if len(tc.TestData) > 0 {
		sb.WriteString("**Test Data:**\n")
		sb.WriteString("| Field | Value | Notes |\n")
		sb.WriteString("|-------|-------|-------|\n")
		for _, td := range tc.TestData {
			sb.WriteString(fmt.Sprintf("| %s | %v | %s |\n", td.Field, td.Value, td.Notes))
		}
		sb.WriteString("\n")
	}

	// Steps
	sb.WriteString("**Steps:**\n")
	for i, s := range tc.Steps {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, s))
	}
	sb.WriteString("\n")

	// Expected Results
	sb.WriteString("**Expected Result:**\n")
	for _, r := range tc.ExpectedResults {
		sb.WriteString(fmt.Sprintf("- %s\n", r))
	}
	sb.WriteString("\n")

	// Traceability
	sb.WriteString("**Traceability:**\n")
	sb.WriteString(fmt.Sprintf("- AC: %s\n", ToLink(tc.ACRef, L1BasePath+"/acceptance-criteria.md")))
	if len(tc.BRRefs) > 0 {
		sb.WriteString("- BR: ")
		for i, br := range tc.BRRefs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(ToLink(br, L1BasePath+"/business-rules.md"))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}
