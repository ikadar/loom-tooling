// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatTestCases formats test cases as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: test-cases.md
func FormatTestCases(suites []TestSuite, summary TDAISummary) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Test Cases", "L3"))
	sb.WriteString("# Test Cases\n\n")

	// Summary section
	sb.WriteString(formatTDAISummary(summary))
	sb.WriteString("\n---\n\n")

	// Test suites grouped by AC
	for _, suite := range suites {
		sb.WriteString(formatTestSuite(suite))
	}

	return sb.String()
}

func formatTDAISummary(summary TDAISummary) string {
	var sb strings.Builder

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("**Total Test Cases:** %d\n\n", summary.Total))

	sb.WriteString("**By Category:**\n")
	sb.WriteString(fmt.Sprintf("- Positive: %d\n", summary.ByCategory.Positive))
	sb.WriteString(fmt.Sprintf("- Negative: %d\n", summary.ByCategory.Negative))
	sb.WriteString(fmt.Sprintf("- Boundary: %d\n", summary.ByCategory.Boundary))
	sb.WriteString(fmt.Sprintf("- Hallucination: %d\n", summary.ByCategory.Hallucination))
	sb.WriteString("\n")

	sb.WriteString("**Coverage:**\n")
	sb.WriteString(fmt.Sprintf("- ACs Covered: %d\n", summary.Coverage.ACsCovered))
	sb.WriteString(fmt.Sprintf("- Positive Ratio: %.1f%%\n", summary.Coverage.PositiveRatio*100))
	sb.WriteString(fmt.Sprintf("- Negative Ratio: %.1f%%\n", summary.Coverage.NegativeRatio*100))
	sb.WriteString(fmt.Sprintf("- Has Hallucination Tests: %t\n", summary.Coverage.HasHallucinationTests))
	sb.WriteString("\n")

	return sb.String()
}

func formatTestSuite(suite TestSuite) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s: %s\n\n", suite.ACRef, suite.ACTitle))

	for _, tc := range suite.Tests {
		sb.WriteString(formatTestCase(tc))
	}

	return sb.String()
}

func formatTestCase(tc TestCase) string {
	var sb strings.Builder

	// Header with category badge
	categoryBadge := getCategoryBadge(tc.Category)
	sb.WriteString(fmt.Sprintf("### %s: %s %s\n\n", tc.ID, tc.Name, categoryBadge))
	sb.WriteString(FormatAnchor(tc.ID))
	sb.WriteString("\n\n")

	// AC Reference
	sb.WriteString(fmt.Sprintf("**AC:** %s\n\n", tc.ACRef))

	// BR References
	if len(tc.BRRefs) > 0 {
		sb.WriteString(fmt.Sprintf("**Business Rules:** %s\n\n", strings.Join(tc.BRRefs, ", ")))
	}

	// Preconditions
	if len(tc.Preconditions) > 0 {
		sb.WriteString("**Preconditions:**\n")
		for _, pre := range tc.Preconditions {
			sb.WriteString(fmt.Sprintf("- %s\n", pre))
		}
		sb.WriteString("\n")
	}

	// Test Data
	if len(tc.TestData) > 0 {
		sb.WriteString("**Test Data:**\n\n")
		sb.WriteString("| Field | Value | Notes |\n")
		sb.WriteString("|-------|-------|-------|\n")
		for _, td := range tc.TestData {
			sb.WriteString(fmt.Sprintf("| %s | %v | %s |\n", td.Field, td.Value, td.Notes))
		}
		sb.WriteString("\n")
	}

	// Steps
	if len(tc.Steps) > 0 {
		sb.WriteString("**Steps:**\n")
		for i, step := range tc.Steps {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
		}
		sb.WriteString("\n")
	}

	// Expected Results
	if len(tc.ExpectedResults) > 0 {
		sb.WriteString("**Expected Results:**\n")
		for _, result := range tc.ExpectedResults {
			sb.WriteString(fmt.Sprintf("- %s\n", result))
		}
		sb.WriteString("\n")
	}

	// Should Not (for hallucination tests)
	if tc.ShouldNot != "" {
		sb.WriteString(fmt.Sprintf("**Should NOT:** %s\n\n", tc.ShouldNot))
	}

	sb.WriteString("---\n\n")

	return sb.String()
}

func getCategoryBadge(category string) string {
	switch category {
	case "positive":
		return "[✓ Positive]"
	case "negative":
		return "[✗ Negative]"
	case "boundary":
		return "[⊡ Boundary]"
	case "hallucination":
		return "[⚠ Hallucination]"
	default:
		return ""
	}
}
