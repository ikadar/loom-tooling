package formatter

import (
	"fmt"
	"strings"
)

// FormatTechSpecs formats technical specifications as markdown
func FormatTechSpecs(techSpecs []TechSpec, timestamp string) string {
	var sb strings.Builder

	fm := DefaultFrontmatter("Technical Specifications", timestamp, "L2")
	sb.WriteString(FormatHeaderWithFrontmatter(fm))
	sb.WriteString("---\n\n")

	for _, ts := range techSpecs {
		sb.WriteString(formatTechSpec(ts))
	}

	return sb.String()
}

// formatTechSpec formats a single tech spec
func formatTechSpec(ts TechSpec) string {
	var sb strings.Builder

	// Header with anchor
	sb.WriteString(FormatSectionHeader(2, ts.ID, ts.Name))
	sb.WriteString(fmt.Sprintf("**Rule:** %s\n\n", ts.Rule))
	sb.WriteString(fmt.Sprintf("**Implementation Approach:**\n%s\n\n", ts.Implementation))

	// Validation Points
	if len(ts.ValidationPoints) > 0 {
		sb.WriteString("**Validation Points:**\n")
		for _, vp := range ts.ValidationPoints {
			sb.WriteString(fmt.Sprintf("- %s\n", vp))
		}
		sb.WriteString("\n")
	}

	// Data Requirements
	if len(ts.DataRequirements) > 0 {
		sb.WriteString("**Data Requirements:**\n")
		sb.WriteString("| Field | Type | Constraints | Source |\n")
		sb.WriteString("|-------|------|-------------|--------|\n")
		for _, dr := range ts.DataRequirements {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", dr.Field, dr.Type, dr.Constraints, dr.Source))
		}
		sb.WriteString("\n")
	}

	// Error Handling
	if len(ts.ErrorHandling) > 0 {
		sb.WriteString("**Error Handling:**\n")
		sb.WriteString("| Condition | Error Code | Message | HTTP Status |\n")
		sb.WriteString("|-----------|------------|---------|-------------|\n")
		for _, eh := range ts.ErrorHandling {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n", eh.Condition, eh.ErrorCode, eh.Message, eh.HTTPStatus))
		}
		sb.WriteString("\n")
	}

	// Traceability
	sb.WriteString("**Traceability:**\n")
	sb.WriteString(fmt.Sprintf("- BR: %s\n", ToLink(ts.BRRef, L1BasePath+"/business-rules.md")))
	if len(ts.RelatedACs) > 0 {
		sb.WriteString("- Related ACs: ")
		for i, ac := range ts.RelatedACs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(ToLink(ac, L1BasePath+"/acceptance-criteria.md"))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}
