// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements test case formatting.
//
// Implements: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatTestCases formats test cases as markdown.
func FormatTestCases(suites []TestSuite, summary TDAISummary) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Test Cases", "L3"))
	result.WriteString("# Test Cases\n\n")

	// Summary section
	result.WriteString("## Summary\n\n")
	result.WriteString(fmt.Sprintf("- **Total Test Cases:** %d\n", summary.Total))
	result.WriteString(fmt.Sprintf("- **Positive:** %d\n", summary.ByCategory.Positive))
	result.WriteString(fmt.Sprintf("- **Negative:** %d\n", summary.ByCategory.Negative))
	result.WriteString(fmt.Sprintf("- **Boundary:** %d\n", summary.ByCategory.Boundary))
	result.WriteString(fmt.Sprintf("- **Hallucination Prevention:** %d\n\n", summary.ByCategory.Hallucination))

	// Test suites
	for _, suite := range suites {
		result.WriteString(fmt.Sprintf("## %s: %s\n\n", suite.ACRef, suite.ACTitle))

		for _, tc := range suite.Tests {
			result.WriteString(fmt.Sprintf("### %s: %s\n\n", tc.ID, tc.Name))
			result.WriteString(fmt.Sprintf("**Category:** %s\n\n", tc.Category))

			if len(tc.BRRefs) > 0 {
				result.WriteString(fmt.Sprintf("**Business Rules:** %s\n\n", strings.Join(tc.BRRefs, ", ")))
			}

			if len(tc.Preconditions) > 0 {
				result.WriteString("**Preconditions:**\n")
				for _, p := range tc.Preconditions {
					result.WriteString(fmt.Sprintf("- %s\n", p))
				}
				result.WriteString("\n")
			}

			if len(tc.TestData) > 0 {
				result.WriteString("**Test Data:**\n\n")
				result.WriteString("| Field | Value | Notes |\n")
				result.WriteString("|-------|-------|-------|\n")
				for _, td := range tc.TestData {
					result.WriteString(fmt.Sprintf("| %s | %v | %s |\n", td.Field, td.Value, td.Notes))
				}
				result.WriteString("\n")
			}

			if len(tc.Steps) > 0 {
				result.WriteString("**Steps:**\n")
				for i, s := range tc.Steps {
					result.WriteString(fmt.Sprintf("%d. %s\n", i+1, s))
				}
				result.WriteString("\n")
			}

			if len(tc.ExpectedResults) > 0 {
				result.WriteString("**Expected Results:**\n")
				for _, e := range tc.ExpectedResults {
					result.WriteString(fmt.Sprintf("- %s\n", e))
				}
				result.WriteString("\n")
			}

			if tc.ShouldNot != "" {
				result.WriteString(fmt.Sprintf("**Should NOT:** %s\n\n", tc.ShouldNot))
			}

			result.WriteString("---\n\n")
		}
	}

	return result.String()
}
