// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatInterfaceContracts formats contracts as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: interface-contracts.md
func FormatInterfaceContracts(contracts []InterfaceContract, sharedTypes []SharedType) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Interface Contracts", "L2"))
	sb.WriteString("# Interface Contracts\n\n")

	// Shared Types section
	if len(sharedTypes) > 0 {
		sb.WriteString("## Shared Types\n\n")
		for _, st := range sharedTypes {
			sb.WriteString(formatSharedType(st))
		}
		sb.WriteString("---\n\n")
	}

	// Contracts
	for _, contract := range contracts {
		sb.WriteString(formatContract(contract))
	}

	return sb.String()
}

func formatSharedType(st SharedType) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### %s\n\n", st.Name))

	if st.Schema != nil {
		schemaJSON, _ := json.MarshalIndent(st.Schema, "", "  ")
		sb.WriteString("```json\n")
		sb.WriteString(string(schemaJSON))
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}

func formatContract(contract InterfaceContract) string {
	var sb strings.Builder

	// Contract header
	sb.WriteString(fmt.Sprintf("## %s – %s\n\n", contract.ID, contract.ServiceName))
	sb.WriteString(FormatAnchor(contract.ID))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", contract.Purpose))
	if contract.BaseURL != "" {
		sb.WriteString(fmt.Sprintf("**Base URL:** `%s`\n\n", contract.BaseURL))
	}

	// Security Requirements
	sb.WriteString(formatSecurityRequirements(contract.SecurityRequirements))

	// Operations
	if len(contract.Operations) > 0 {
		sb.WriteString("### Operations\n\n")
		for _, op := range contract.Operations {
			sb.WriteString(formatOperation(op))
		}
	}

	// Events
	if len(contract.Events) > 0 {
		sb.WriteString("### Events\n\n")
		for _, event := range contract.Events {
			sb.WriteString(formatContractEvent(event))
		}
	}

	sb.WriteString("---\n\n")

	return sb.String()
}

func formatSecurityRequirements(sec SecurityRequirements) string {
	var sb strings.Builder

	sb.WriteString("### Security Requirements\n\n")
	if sec.Authentication != "" {
		sb.WriteString(fmt.Sprintf("- **Authentication:** %s\n", sec.Authentication))
	}
	if sec.Authorization != "" {
		sb.WriteString(fmt.Sprintf("- **Authorization:** %s\n", sec.Authorization))
	}
	if len(sec.Scopes) > 0 {
		sb.WriteString(fmt.Sprintf("- **Scopes:** %s\n", strings.Join(sec.Scopes, ", ")))
	}
	sb.WriteString("\n")

	return sb.String()
}

func formatOperation(op ContractOperation) string {
	var sb strings.Builder

	// Operation header
	sb.WriteString(fmt.Sprintf("#### %s `%s %s`\n\n", op.ID, op.Method, op.Path))
	sb.WriteString(FormatAnchor(op.ID))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("**Summary:** %s\n\n", op.Summary))
	if op.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", op.Description))
	}

	// Request
	sb.WriteString("**Request:**\n")
	if op.Request.ContentType != "" {
		sb.WriteString(fmt.Sprintf("- Content-Type: `%s`\n", op.Request.ContentType))
	}
	if op.Request.Schema != nil {
		schemaJSON, _ := json.MarshalIndent(op.Request.Schema, "", "  ")
		sb.WriteString("\n```json\n")
		sb.WriteString(string(schemaJSON))
		sb.WriteString("\n```\n")
	}
	sb.WriteString("\n")

	// Response
	sb.WriteString("**Response:**\n")
	sb.WriteString(fmt.Sprintf("- Status: `%d`\n", op.Response.StatusCode))
	if op.Response.ContentType != "" {
		sb.WriteString(fmt.Sprintf("- Content-Type: `%s`\n", op.Response.ContentType))
	}
	if op.Response.Schema != nil {
		schemaJSON, _ := json.MarshalIndent(op.Response.Schema, "", "  ")
		sb.WriteString("\n```json\n")
		sb.WriteString(string(schemaJSON))
		sb.WriteString("\n```\n")
	}
	sb.WriteString("\n")

	// Errors
	if len(op.Errors) > 0 {
		sb.WriteString("**Errors:**\n\n")
		sb.WriteString("| Status | Code | Description |\n")
		sb.WriteString("|--------|------|-------------|\n")
		for _, err := range op.Errors {
			sb.WriteString(fmt.Sprintf("| %d | %s | %s |\n",
				err.StatusCode, err.Code, err.Description))
		}
		sb.WriteString("\n")
	}

	// Related ACs
	if len(op.RelatedACs) > 0 {
		sb.WriteString(fmt.Sprintf("**Related ACs:** %s\n\n", strings.Join(op.RelatedACs, ", ")))
	}

	return sb.String()
}

func formatContractEvent(event ContractEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("#### Event: %s\n\n", event.Name))
	sb.WriteString(fmt.Sprintf("%s\n\n", event.Description))

	if event.Payload != nil {
		payloadJSON, _ := json.MarshalIndent(event.Payload, "", "  ")
		sb.WriteString("**Payload:**\n```json\n")
		sb.WriteString(string(payloadJSON))
		sb.WriteString("\n```\n\n")
	}

	return sb.String()
}
