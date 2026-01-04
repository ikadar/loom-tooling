// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatSequenceDesign formats sequences as markdown with Mermaid diagrams.
//
// Implements: l2/package-structure.md PKG-007
// Output: sequence-design.md
func FormatSequenceDesign(sequences []SequenceDesign) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Sequence Design", "L2"))
	sb.WriteString("# Sequence Design\n\n")

	for _, seq := range sequences {
		sb.WriteString(formatSequence(seq))
	}

	return sb.String()
}

func formatSequence(seq SequenceDesign) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("## %s – %s\n\n", seq.ID, seq.Name))
	sb.WriteString(FormatAnchor(seq.ID))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("%s\n\n", seq.Description))

	// Trigger
	sb.WriteString("### Trigger\n\n")
	sb.WriteString(fmt.Sprintf("**Actor:** %s\n", seq.Trigger.Actor))
	sb.WriteString(fmt.Sprintf("**Action:** %s\n\n", seq.Trigger.Action))

	// Mermaid sequence diagram
	sb.WriteString("### Sequence Diagram\n\n")
	sb.WriteString(generateMermaidSequence(seq))

	// Participants
	sb.WriteString("### Participants\n\n")
	sb.WriteString("| Name | Type |\n")
	sb.WriteString("|------|------|\n")
	for _, p := range seq.Participants {
		sb.WriteString(fmt.Sprintf("| %s | %s |\n", p.Name, p.Type))
	}
	sb.WriteString("\n")

	// Steps (textual description)
	if len(seq.Steps) > 0 {
		sb.WriteString("### Steps\n\n")
		for i, step := range seq.Steps {
			asyncLabel := ""
			if step.Async {
				asyncLabel = " *(async)*"
			}
			if step.Returns != "" {
				sb.WriteString(fmt.Sprintf("%d. **%s** → **%s**: %s → %s%s\n",
					i+1, step.From, step.To, step.Action, step.Returns, asyncLabel))
			} else {
				sb.WriteString(fmt.Sprintf("%d. **%s** → **%s**: %s%s\n",
					i+1, step.From, step.To, step.Action, asyncLabel))
			}
		}
		sb.WriteString("\n")
	}

	// Outcome
	sb.WriteString("### Outcome\n\n")
	sb.WriteString(fmt.Sprintf("**Success:** %s\n", seq.Outcome.Success))
	sb.WriteString(fmt.Sprintf("**Result:** %s\n\n", seq.Outcome.Result))

	// Exceptions
	if len(seq.Exceptions) > 0 {
		sb.WriteString("### Exception Flows\n\n")
		sb.WriteString("| Condition | Handler | Result |\n")
		sb.WriteString("|-----------|---------|--------|\n")
		for _, exc := range seq.Exceptions {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				exc.Condition, exc.Handler, exc.Result))
		}
		sb.WriteString("\n")
	}

	// Traceability
	var sources []string
	for _, ac := range seq.RelatedACs {
		sources = append(sources, "AC: "+ac)
	}
	for _, br := range seq.RelatedBRs {
		sources = append(sources, "BR: "+br)
	}
	if len(sources) > 0 {
		sb.WriteString("### Traceability\n\n")
		for _, src := range sources {
			sb.WriteString(fmt.Sprintf("- %s\n", src))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	return sb.String()
}

// generateMermaidSequence creates a Mermaid sequence diagram.
func generateMermaidSequence(seq SequenceDesign) string {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("sequenceDiagram\n")

	// Participants
	for _, p := range seq.Participants {
		participantType := getParticipantKeyword(p.Type)
		sb.WriteString(fmt.Sprintf("    %s %s\n", participantType, p.Name))
	}
	sb.WriteString("\n")

	// Steps
	for _, step := range seq.Steps {
		arrow := "->>"
		if step.Async {
			arrow = "--))"
		}

		if step.Returns != "" {
			// Request
			sb.WriteString(fmt.Sprintf("    %s%s%s: %s\n",
				step.From, arrow, step.To, step.Action))
			// Response (dashed line)
			sb.WriteString(fmt.Sprintf("    %s-->>%s: %s\n",
				step.To, step.From, step.Returns))
		} else {
			sb.WriteString(fmt.Sprintf("    %s%s%s: %s\n",
				step.From, arrow, step.To, step.Action))
		}
	}

	sb.WriteString("```\n\n")

	return sb.String()
}

func getParticipantKeyword(pType string) string {
	switch pType {
	case "actor":
		return "actor"
	case "database":
		return "participant"
	case "external":
		return "participant"
	default:
		return "participant"
	}
}
