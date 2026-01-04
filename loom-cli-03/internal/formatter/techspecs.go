// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements tech specs formatting.
//
// Implements: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatTechSpecs formats tech specs as markdown.
func FormatTechSpecs(specs []TechSpec) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Technical Specifications", "L2"))
	result.WriteString("# Technical Specifications\n\n")

	for _, spec := range specs {
		result.WriteString(fmt.Sprintf("## %s: %s\n\n", spec.ID, spec.Name))

		if spec.BRRef != "" {
			result.WriteString(fmt.Sprintf("**Business Rule:** %s\n\n", spec.BRRef))
		}

		result.WriteString(fmt.Sprintf("**Rule:** %s\n\n", spec.Rule))
		result.WriteString(fmt.Sprintf("**Implementation:** %s\n\n", spec.Implementation))

		if len(spec.ValidationPoints) > 0 {
			result.WriteString("**Validation Points:**\n")
			for _, vp := range spec.ValidationPoints {
				result.WriteString(fmt.Sprintf("- %s\n", vp))
			}
			result.WriteString("\n")
		}

		if len(spec.DataRequirements) > 0 {
			result.WriteString("**Data Requirements:**\n\n")
			result.WriteString("| Field | Type | Constraints |\n")
			result.WriteString("|-------|------|-------------|\n")
			for _, dr := range spec.DataRequirements {
				result.WriteString(fmt.Sprintf("| %s | %s | %s |\n", dr.Field, dr.Type, dr.Constraints))
			}
			result.WriteString("\n")
		}

		if len(spec.ErrorHandling) > 0 {
			result.WriteString("**Error Handling:**\n\n")
			result.WriteString("| Condition | Error Code | Message |\n")
			result.WriteString("|-----------|------------|---------|\n")
			for _, eh := range spec.ErrorHandling {
				result.WriteString(fmt.Sprintf("| %s | %s | %s |\n", eh.Condition, eh.ErrorCode, eh.Message))
			}
			result.WriteString("\n")
		}

		if len(spec.RelatedACs) > 0 {
			result.WriteString(fmt.Sprintf("**Related ACs:** %s\n\n", strings.Join(spec.RelatedACs, ", ")))
		}

		result.WriteString("---\n\n")
	}

	return result.String()
}
