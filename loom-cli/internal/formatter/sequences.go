package formatter

import (
	"fmt"
	"strings"
)

// FormatSequenceDesign formats sequence design as markdown
func FormatSequenceDesign(sequences []SequenceDesign, timestamp string) string {
	var sb strings.Builder

	fm := DefaultFrontmatter("Sequence Design", timestamp, "L2")
	sb.WriteString(FormatHeaderWithFrontmatter(fm))
	sb.WriteString("---\n\n")

	for _, seq := range sequences {
		sb.WriteString(formatSequence(seq))
	}

	return sb.String()
}

// formatSequence formats a single sequence
func formatSequence(seq SequenceDesign) string {
	var sb strings.Builder

	sb.WriteString(FormatSectionHeader(2, seq.ID, seq.Name))
	sb.WriteString(fmt.Sprintf("%s\n\n", seq.Description))

	// Trigger
	sb.WriteString("### Trigger\n\n")
	sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", seq.Trigger.Type))
	sb.WriteString(fmt.Sprintf("%s\n\n", seq.Trigger.Description))

	// Participants
	if len(seq.Participants) > 0 {
		sb.WriteString("### Participants\n\n")
		for _, p := range seq.Participants {
			sb.WriteString(fmt.Sprintf("- **%s** (%s)\n", p.Name, p.Type))
		}
		sb.WriteString("\n")
	}

	// Steps
	sb.WriteString("### Sequence\n\n")
	for _, step := range seq.Steps {
		sb.WriteString(fmt.Sprintf("%d. **%s** → %s: %s\n", step.Step, step.Actor, step.Target, step.Action))
		if step.Event != "" {
			sb.WriteString(fmt.Sprintf("   - Emits: `%s`\n", step.Event))
		}
		if step.Returns != "" {
			sb.WriteString(fmt.Sprintf("   - Returns: %s\n", step.Returns))
		}
	}
	sb.WriteString("\n")

	// Mermaid Diagram
	sb.WriteString("### Sequence Diagram\n\n")
	sb.WriteString("```mermaid\nsequenceDiagram\n")
	for _, p := range seq.Participants {
		sb.WriteString(fmt.Sprintf("    participant %s\n", p.Name))
	}
	for _, step := range seq.Steps {
		sb.WriteString(fmt.Sprintf("    %s->>%s: %s\n", step.Actor, step.Target, step.Action))
		if step.Returns != "" {
			sb.WriteString(fmt.Sprintf("    %s-->>%s: %s\n", step.Target, step.Actor, step.Returns))
		}
	}
	sb.WriteString("```\n\n")

	// Outcome
	sb.WriteString("### Outcome\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", seq.Outcome.Success))
	if len(seq.Outcome.StateChanges) > 0 {
		sb.WriteString("**State Changes:**\n")
		for _, sc := range seq.Outcome.StateChanges {
			sb.WriteString(fmt.Sprintf("- %s\n", sc))
		}
		sb.WriteString("\n")
	}

	// Exceptions
	if len(seq.Exceptions) > 0 {
		sb.WriteString("### Exceptions\n\n")
		for _, ex := range seq.Exceptions {
			sb.WriteString(fmt.Sprintf("- **%s** (step %d): %s\n", ex.Condition, ex.Step, ex.Handling))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	return sb.String()
}
