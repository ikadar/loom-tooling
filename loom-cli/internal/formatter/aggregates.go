package formatter

import (
	"fmt"
	"strings"
)

// FormatAggregateDesign formats aggregate design as markdown
func FormatAggregateDesign(aggregates []AggregateDesign, timestamp string) string {
	var sb strings.Builder

	fm := DefaultFrontmatter("Aggregate Design", timestamp, "L2")
	sb.WriteString(FormatHeaderWithFrontmatter(fm))
	sb.WriteString("---\n\n")

	for _, agg := range aggregates {
		sb.WriteString(formatAggregate(agg))
	}

	return sb.String()
}

// formatAggregate formats a single aggregate
func formatAggregate(agg AggregateDesign) string {
	var sb strings.Builder

	sb.WriteString(FormatSectionHeader(2, agg.ID, agg.Name))
	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", agg.Purpose))

	// Root
	sb.WriteString(fmt.Sprintf("### Aggregate Root: %s\n\n", agg.Root.Entity))
	sb.WriteString(fmt.Sprintf("**Identity:** %s\n\n", agg.Root.Identity))
	if len(agg.Root.Attributes) > 0 {
		sb.WriteString("**Attributes:**\n")
		sb.WriteString("| Name | Type | Mutable |\n")
		sb.WriteString("|------|------|----------|\n")
		for _, attr := range agg.Root.Attributes {
			sb.WriteString(fmt.Sprintf("| %s | %s | %v |\n", attr.Name, attr.Type, attr.Mutable))
		}
		sb.WriteString("\n")
	}

	// Invariants
	if len(agg.Invariants) > 0 {
		sb.WriteString("### Invariants\n\n")
		for _, inv := range agg.Invariants {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n  - Enforcement: %s\n", inv.ID, inv.Rule, inv.Enforcement))
		}
		sb.WriteString("\n")
	}

	// Child Entities
	if len(agg.Entities) > 0 {
		sb.WriteString("### Child Entities\n\n")
		for _, ent := range agg.Entities {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", ent.Name))
			sb.WriteString(fmt.Sprintf("**Identity:** %s\n\n", ent.Identity))
			sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", ent.Purpose))
		}
	}

	// Value Objects
	if len(agg.ValueObjects) > 0 {
		sb.WriteString("### Value Objects\n\n")
		sb.WriteString(fmt.Sprintf("%v\n\n", agg.ValueObjects))
	}

	// Behaviors
	if len(agg.Behaviors) > 0 {
		sb.WriteString("### Behaviors\n\n")
		sb.WriteString("| Command | Pre | Post | Emits |\n")
		sb.WriteString("|---------|-----|------|-------|\n")
		for _, b := range agg.Behaviors {
			sb.WriteString(fmt.Sprintf("| %s | %v | %v | %s |\n", b.Command, b.Preconditions, b.Postconditions, b.Emits))
		}
		sb.WriteString("\n")
	}

	// Events
	if len(agg.Events) > 0 {
		sb.WriteString("### Events\n\n")
		for _, ev := range agg.Events {
			sb.WriteString(fmt.Sprintf("- **%s**: %v\n", ev.Name, ev.Payload))
		}
		sb.WriteString("\n")
	}

	// Repository
	sb.WriteString(fmt.Sprintf("### Repository: %s\n\n", agg.Repository.Name))
	sb.WriteString(fmt.Sprintf("- Load Strategy: %s\n", agg.Repository.LoadStrategy))
	sb.WriteString(fmt.Sprintf("- Concurrency: %s\n\n", agg.Repository.Concurrency))

	sb.WriteString("---\n\n")

	return sb.String()
}
