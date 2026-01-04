// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatAggregateDesign formats aggregate design as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: aggregate-design.md
func FormatAggregateDesign(aggregates []AggregateDesign) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Aggregate Design", "L2"))
	sb.WriteString("# Aggregate Design\n\n")

	for _, agg := range aggregates {
		sb.WriteString(formatAggregate(agg))
	}

	return sb.String()
}

func formatAggregate(agg AggregateDesign) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("## %s – %s\n\n", agg.ID, agg.Name))
	sb.WriteString(FormatAnchor(agg.ID))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", agg.Purpose))

	// Invariants
	if len(agg.Invariants) > 0 {
		sb.WriteString("### Invariants\n\n")
		for _, inv := range agg.Invariants {
			sb.WriteString(fmt.Sprintf("- **%s:** %s\n", inv.ID, inv.Description))
			if inv.Enforcement != "" {
				sb.WriteString(fmt.Sprintf("  - *Enforcement:* %s\n", inv.Enforcement))
			}
		}
		sb.WriteString("\n")
	}

	// Aggregate Root
	sb.WriteString("### Aggregate Root\n\n")
	sb.WriteString(fmt.Sprintf("**Entity:** `%s`\n\n", agg.Root.Name))

	if len(agg.Root.Attributes) > 0 {
		sb.WriteString("**Attributes:**\n\n")
		sb.WriteString("| Name | Type | Required |\n")
		sb.WriteString("|------|------|----------|\n")
		for _, attr := range agg.Root.Attributes {
			required := "No"
			if attr.Required {
				required = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				attr.Name, attr.Type, required))
		}
		sb.WriteString("\n")
	}

	if len(agg.Root.Methods) > 0 {
		sb.WriteString("**Methods:**\n")
		for _, method := range agg.Root.Methods {
			sb.WriteString(fmt.Sprintf("- `%s`\n", method))
		}
		sb.WriteString("\n")
	}

	// Child Entities
	if len(agg.Entities) > 0 {
		sb.WriteString("### Child Entities\n\n")
		for _, entity := range agg.Entities {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", entity.Name))
			if len(entity.Attributes) > 0 {
				sb.WriteString("| Name | Type | Required |\n")
				sb.WriteString("|------|------|----------|\n")
				for _, attr := range entity.Attributes {
					required := "No"
					if attr.Required {
						required = "Yes"
					}
					sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
						attr.Name, attr.Type, required))
				}
				sb.WriteString("\n")
			}
		}
	}

	// Value Objects
	if len(agg.ValueObjects) > 0 {
		sb.WriteString("### Value Objects\n\n")
		for _, vo := range agg.ValueObjects {
			sb.WriteString(fmt.Sprintf("- `%s`\n", vo))
		}
		sb.WriteString("\n")
	}

	// Behaviors
	if len(agg.Behaviors) > 0 {
		sb.WriteString("### Behaviors\n\n")
		for _, beh := range agg.Behaviors {
			sb.WriteString(fmt.Sprintf("#### `%s`\n\n", beh.Name))
			sb.WriteString(fmt.Sprintf("%s\n\n", beh.Description))
			if len(beh.Parameters) > 0 {
				sb.WriteString(fmt.Sprintf("**Parameters:** %s\n\n", strings.Join(beh.Parameters, ", ")))
			}
			if beh.Returns != "" {
				sb.WriteString(fmt.Sprintf("**Returns:** %s\n\n", beh.Returns))
			}
			if len(beh.Raises) > 0 {
				sb.WriteString(fmt.Sprintf("**Raises:** %s\n\n", strings.Join(beh.Raises, ", ")))
			}
		}
	}

	// Events
	if len(agg.Events) > 0 {
		sb.WriteString("### Domain Events\n\n")
		for _, event := range agg.Events {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", event.Name))
			sb.WriteString(fmt.Sprintf("**Trigger:** %s\n\n", event.Trigger))
			if len(event.Payload) > 0 {
				sb.WriteString(fmt.Sprintf("**Payload:** %s\n\n", strings.Join(event.Payload, ", ")))
			}
		}
	}

	// Repository
	sb.WriteString("### Repository\n\n")
	sb.WriteString(fmt.Sprintf("**Interface:** `%s`\n\n", agg.Repository.Name))
	if len(agg.Repository.Methods) > 0 {
		sb.WriteString("**Methods:**\n")
		for _, method := range agg.Repository.Methods {
			sb.WriteString(fmt.Sprintf("- `%s`\n", method))
		}
		sb.WriteString("\n")
	}

	// External References
	if len(agg.ExternalReferences) > 0 {
		sb.WriteString("### External References\n\n")
		sb.WriteString("| Aggregate | Type | Via |\n")
		sb.WriteString("|-----------|------|-----|\n")
		for _, ref := range agg.ExternalReferences {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				ref.Aggregate, ref.Type, ref.Via))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	return sb.String()
}
