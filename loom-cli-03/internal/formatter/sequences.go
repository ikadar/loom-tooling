// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements sequence design formatting with Mermaid diagrams.
//
// Implements: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatSequenceDesign formats sequence designs as markdown with Mermaid diagrams.
func FormatSequenceDesign(sequences []SequenceDesign) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Sequence Design", "L2"))
	result.WriteString("# Sequence Design\n\n")

	for _, seq := range sequences {
		result.WriteString(fmt.Sprintf("## %s: %s\n\n", seq.ID, seq.Name))
		result.WriteString(fmt.Sprintf("%s\n\n", seq.Description))

		// Trigger
		result.WriteString(fmt.Sprintf("**Trigger:** %s %s\n\n", seq.Trigger.Actor, seq.Trigger.Action))

		// Related refs
		if len(seq.RelatedACs) > 0 {
			result.WriteString(fmt.Sprintf("**Related ACs:** %s\n\n", strings.Join(seq.RelatedACs, ", ")))
		}
		if len(seq.RelatedBRs) > 0 {
			result.WriteString(fmt.Sprintf("**Related BRs:** %s\n\n", strings.Join(seq.RelatedBRs, ", ")))
		}

		// Mermaid sequence diagram
		result.WriteString("### Sequence Diagram\n\n")
		result.WriteString("```mermaid\nsequenceDiagram\n")

		// Participants
		for _, p := range seq.Participants {
			result.WriteString(fmt.Sprintf("    participant %s as %s\n", sanitizeMermaidID(p.Name), p.Name))
		}
		result.WriteString("\n")

		// Steps
		for _, step := range seq.Steps {
			from := sanitizeMermaidID(step.From)
			to := sanitizeMermaidID(step.To)
			result.WriteString(fmt.Sprintf("    %s->>%s: %s\n", from, to, step.Action))
			if step.Returns != "" {
				result.WriteString(fmt.Sprintf("    %s-->>%s: %s\n", to, from, step.Returns))
			}
		}

		result.WriteString("```\n\n")

		// Outcome
		result.WriteString("### Outcomes\n\n")
		result.WriteString(fmt.Sprintf("**Success:** %s\n\n", seq.Outcome.Success))
		result.WriteString(fmt.Sprintf("**Failure:** %s\n\n", seq.Outcome.Failure))

		// Exceptions
		if len(seq.Exceptions) > 0 {
			result.WriteString("### Exceptions\n\n")
			result.WriteString("| Condition | Handling |\n")
			result.WriteString("|-----------|----------|\n")
			for _, e := range seq.Exceptions {
				result.WriteString(fmt.Sprintf("| %s | %s |\n", e.Condition, e.Handling))
			}
			result.WriteString("\n")
		}

		result.WriteString("---\n\n")
	}

	return result.String()
}

// sanitizeMermaidID converts a name to a valid Mermaid participant ID.
func sanitizeMermaidID(name string) string {
	// Remove spaces and special characters
	result := strings.ReplaceAll(name, " ", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, ".", "_")
	return result
}
