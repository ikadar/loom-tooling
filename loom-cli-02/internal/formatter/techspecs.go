// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatTechSpecs formats tech specs as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: tech-specs.md
func FormatTechSpecs(specs []TechSpec) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Technical Specifications", "L2"))
	sb.WriteString("# Technical Specifications\n\n")

	for _, spec := range specs {
		sb.WriteString(formatTechSpec(spec))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatTechSpec(spec TechSpec) string {
	var sb strings.Builder

	// Header with anchor
	sb.WriteString(fmt.Sprintf("## %s – %s\n\n", spec.ID, spec.Name))
	sb.WriteString(FormatAnchor(spec.ID))
	sb.WriteString("\n\n")

	// Business Rule reference
	if spec.BRRef != "" {
		sb.WriteString(fmt.Sprintf("**Business Rule:** %s\n\n", spec.BRRef))
	}

	// Rule
	if spec.Rule != "" {
		sb.WriteString(fmt.Sprintf("**Rule:** %s\n\n", spec.Rule))
	}

	// Implementation
	if spec.Implementation != "" {
		sb.WriteString(fmt.Sprintf("**Implementation:** %s\n\n", spec.Implementation))
	}

	// Validation Points
	if len(spec.ValidationPoints) > 0 {
		sb.WriteString("**Validation Points:**\n")
		for _, point := range spec.ValidationPoints {
			sb.WriteString(fmt.Sprintf("- %s\n", point))
		}
		sb.WriteString("\n")
	}

	// Data Requirements
	if len(spec.DataRequirements) > 0 {
		sb.WriteString("**Data Requirements:**\n\n")
		sb.WriteString("| Field | Type | Constraints | Source |\n")
		sb.WriteString("|-------|------|-------------|--------|\n")
		for _, dr := range spec.DataRequirements {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dr.Field, dr.Type, dr.Constraints, dr.Source))
		}
		sb.WriteString("\n")
	}

	// Error Handling
	if len(spec.ErrorHandling) > 0 {
		sb.WriteString("**Error Handling:**\n\n")
		sb.WriteString("| Condition | Code | Message | HTTP |\n")
		sb.WriteString("|-----------|------|---------|------|\n")
		for _, eh := range spec.ErrorHandling {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
				eh.Condition, eh.ErrorCode, eh.Message, eh.HTTPStatus))
		}
		sb.WriteString("\n")
	}

	// Traceability (Related ACs)
	if len(spec.RelatedACs) > 0 {
		sb.WriteString("**Traceability:**\n")
		for _, ac := range spec.RelatedACs {
			sb.WriteString(fmt.Sprintf("- AC: %s\n", ac))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
