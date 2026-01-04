// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements aggregate design formatting.
//
// Implements: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatAggregateDesign formats aggregate designs as markdown.
func FormatAggregateDesign(aggregates []AggregateDesign) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Aggregate Design", "L2"))
	result.WriteString("# Aggregate Design\n\n")

	for _, agg := range aggregates {
		result.WriteString(fmt.Sprintf("## %s: %s\n\n", agg.ID, agg.Name))
		result.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", agg.Purpose))

		// Invariants
		if len(agg.Invariants) > 0 {
			result.WriteString("### Invariants\n\n")
			for _, inv := range agg.Invariants {
				result.WriteString(fmt.Sprintf("- **%s:** %s\n", inv.ID, inv.Description))
				if inv.Enforcement != "" {
					result.WriteString(fmt.Sprintf("  - Enforcement: %s\n", inv.Enforcement))
				}
			}
			result.WriteString("\n")
		}

		// Root
		result.WriteString("### Aggregate Root\n\n")
		result.WriteString(fmt.Sprintf("**%s** (ID: `%s`)\n\n", agg.Root.Name, agg.Root.Identifier))
		if len(agg.Root.Fields) > 0 {
			result.WriteString("Fields:\n")
			for _, f := range agg.Root.Fields {
				result.WriteString(fmt.Sprintf("- %s\n", f))
			}
			result.WriteString("\n")
		}

		// Entities
		if len(agg.Entities) > 0 {
			result.WriteString("### Entities\n\n")
			for _, ent := range agg.Entities {
				result.WriteString(fmt.Sprintf("#### %s\n\n", ent.Name))
				if len(ent.Fields) > 0 {
					result.WriteString("Fields:\n")
					for _, f := range ent.Fields {
						result.WriteString(fmt.Sprintf("- %s\n", f))
					}
					result.WriteString("\n")
				}
			}
		}

		// Value Objects
		if len(agg.ValueObjects) > 0 {
			result.WriteString("### Value Objects\n\n")
			for _, vo := range agg.ValueObjects {
				result.WriteString(fmt.Sprintf("- %s\n", vo))
			}
			result.WriteString("\n")
		}

		// Behaviors
		if len(agg.Behaviors) > 0 {
			result.WriteString("### Behaviors\n\n")
			for _, b := range agg.Behaviors {
				result.WriteString(fmt.Sprintf("#### %s\n\n", b.Name))
				result.WriteString(fmt.Sprintf("%s\n\n", b.Description))
				if len(b.ACRefs) > 0 {
					result.WriteString(fmt.Sprintf("**ACs:** %s\n\n", strings.Join(b.ACRefs, ", ")))
				}
			}
		}

		// Events
		if len(agg.Events) > 0 {
			result.WriteString("### Domain Events\n\n")
			for _, e := range agg.Events {
				result.WriteString(fmt.Sprintf("- **%s**\n", e.Name))
				if len(e.Payload) > 0 {
					result.WriteString(fmt.Sprintf("  - Payload: %s\n", strings.Join(e.Payload, ", ")))
				}
			}
			result.WriteString("\n")
		}

		// Repository
		if len(agg.Repository.Operations) > 0 {
			result.WriteString("### Repository\n\n")
			for _, op := range agg.Repository.Operations {
				result.WriteString(fmt.Sprintf("- %s\n", op))
			}
			result.WriteString("\n")
		}

		// External References
		if len(agg.ExternalReferences) > 0 {
			result.WriteString("### External References\n\n")
			for _, ref := range agg.ExternalReferences {
				result.WriteString(fmt.Sprintf("- %s (%s)\n", ref.Aggregate, ref.Type))
			}
			result.WriteString("\n")
		}

		result.WriteString("---\n\n")
	}

	return result.String()
}
