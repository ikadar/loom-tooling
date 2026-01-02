package formatter

import (
	"fmt"
	"strings"
)

// FormatInterfaceContracts formats interface contracts as markdown
func FormatInterfaceContracts(contracts []InterfaceContract, sharedTypes []SharedType, timestamp string) string {
	var sb strings.Builder

	sb.WriteString("# Interface Contracts\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", timestamp))
	sb.WriteString("---\n\n")

	// Shared Types section
	if len(sharedTypes) > 0 {
		sb.WriteString("## Shared Types\n\n")
		for _, st := range sharedTypes {
			sb.WriteString(fmt.Sprintf("### %s\n\n", st.Name))
			sb.WriteString("| Field | Type | Constraints |\n")
			sb.WriteString("|-------|------|-------------|\n")
			for _, field := range st.Fields {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", field.Name, field.Type, field.Constraints))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("---\n\n")
	}

	// Contracts
	for _, ic := range contracts {
		sb.WriteString(formatContract(ic))
	}

	return sb.String()
}

// formatContract formats a single interface contract
func formatContract(ic InterfaceContract) string {
	var sb strings.Builder

	sb.WriteString(FormatSectionHeader(2, ic.ID, ic.ServiceName))
	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", ic.Purpose))
	sb.WriteString(fmt.Sprintf("**Base URL:** `%s`\n\n", ic.BaseURL))

	if ic.SecurityRequirements.Authentication != "" {
		sb.WriteString("**Security:**\n")
		sb.WriteString(fmt.Sprintf("- Authentication: %s\n", ic.SecurityRequirements.Authentication))
		sb.WriteString(fmt.Sprintf("- Authorization: %s\n\n", ic.SecurityRequirements.Authorization))
	}

	// Operations
	sb.WriteString("### Operations\n\n")
	for _, op := range ic.Operations {
		sb.WriteString(fmt.Sprintf("#### %s `%s %s`\n\n", op.Name, op.Method, op.Path))
		sb.WriteString(fmt.Sprintf("%s\n\n", op.Description))

		if len(op.InputSchema) > 0 {
			sb.WriteString("**Input:**\n")
			for name, field := range op.InputSchema {
				req := ""
				if field.Required {
					req = " (required)"
				}
				sb.WriteString(fmt.Sprintf("- `%s`: %s%s\n", name, field.Type, req))
			}
			sb.WriteString("\n")
		}

		if len(op.OutputSchema) > 0 {
			sb.WriteString("**Output:**\n")
			for name, field := range op.OutputSchema {
				sb.WriteString(fmt.Sprintf("- `%s`: %s\n", name, field.Type))
			}
			sb.WriteString("\n")
		}

		if len(op.Errors) > 0 {
			sb.WriteString("**Errors:**\n")
			sb.WriteString("| Code | HTTP | Message |\n")
			sb.WriteString("|------|------|----------|\n")
			for _, e := range op.Errors {
				sb.WriteString(fmt.Sprintf("| %s | %d | %s |\n", e.Code, e.HTTPStatus, e.Message))
			}
			sb.WriteString("\n")
		}

		if len(op.Preconditions) > 0 {
			sb.WriteString(fmt.Sprintf("**Preconditions:** %v\n\n", op.Preconditions))
		}
		if len(op.Postconditions) > 0 {
			sb.WriteString(fmt.Sprintf("**Postconditions:** %v\n\n", op.Postconditions))
		}
	}

	// Events
	if len(ic.Events) > 0 {
		sb.WriteString("### Events\n\n")
		for _, ev := range ic.Events {
			sb.WriteString(fmt.Sprintf("- **%s**: %s (payload: %v)\n", ev.Name, ev.Description, ev.Payload))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	return sb.String()
}
