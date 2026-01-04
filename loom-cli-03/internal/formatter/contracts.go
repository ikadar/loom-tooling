// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements interface contracts formatting.
//
// Implements: l2/internal-api.md
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatInterfaceContracts formats interface contracts as markdown.
func FormatInterfaceContracts(contracts []InterfaceContract, sharedTypes []SharedType) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Interface Contracts", "L2"))
	result.WriteString("# Interface Contracts\n\n")

	for _, contract := range contracts {
		result.WriteString(fmt.Sprintf("## %s: %s\n\n", contract.ID, contract.ServiceName))
		result.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", contract.Purpose))

		if contract.BaseURL != "" {
			result.WriteString(fmt.Sprintf("**Base URL:** `%s`\n\n", contract.BaseURL))
		}

		// Operations
		if len(contract.Operations) > 0 {
			result.WriteString("### Operations\n\n")
			for _, op := range contract.Operations {
				result.WriteString(fmt.Sprintf("#### %s\n\n", op.Name))
				result.WriteString(fmt.Sprintf("**%s** `%s`\n\n", op.Method, op.Path))
				result.WriteString(fmt.Sprintf("%s\n\n", op.Description))

				// Request
				if op.Request.Schema != nil {
					result.WriteString("**Request:**\n\n")
					result.WriteString(fmt.Sprintf("Content-Type: `%s`\n\n", op.Request.ContentType))
					schemaJSON, _ := json.MarshalIndent(op.Request.Schema, "", "  ")
					result.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(schemaJSON)))
				}

				// Response
				if op.Response.Schema != nil {
					result.WriteString("**Response:**\n\n")
					result.WriteString(fmt.Sprintf("Status: `%d`\n\n", op.Response.StatusCode))
					schemaJSON, _ := json.MarshalIndent(op.Response.Schema, "", "  ")
					result.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(schemaJSON)))
				}

				// Errors
				if len(op.Errors) > 0 {
					result.WriteString("**Errors:**\n\n")
					result.WriteString("| Status | Code | Message |\n")
					result.WriteString("|--------|------|----------|\n")
					for _, e := range op.Errors {
						result.WriteString(fmt.Sprintf("| %d | %s | %s |\n", e.StatusCode, e.Code, e.Message))
					}
					result.WriteString("\n")
				}
			}
		}

		// Events
		if len(contract.Events) > 0 {
			result.WriteString("### Events\n\n")
			for _, evt := range contract.Events {
				result.WriteString(fmt.Sprintf("#### %s\n\n", evt.Name))
				result.WriteString(fmt.Sprintf("**Type:** %s\n\n", evt.Type))
				if evt.Payload != nil {
					payloadJSON, _ := json.MarshalIndent(evt.Payload, "", "  ")
					result.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(payloadJSON)))
				}
			}
		}

		// Security
		result.WriteString("### Security Requirements\n\n")
		result.WriteString(fmt.Sprintf("- **Authentication:** %s\n", contract.SecurityRequirements.Authentication))
		result.WriteString(fmt.Sprintf("- **Authorization:** %s\n", contract.SecurityRequirements.Authorization))
		if len(contract.SecurityRequirements.Scopes) > 0 {
			result.WriteString(fmt.Sprintf("- **Scopes:** %s\n", strings.Join(contract.SecurityRequirements.Scopes, ", ")))
		}
		result.WriteString("\n---\n\n")
	}

	// Shared Types
	if len(sharedTypes) > 0 {
		result.WriteString("## Shared Types\n\n")
		for _, st := range sharedTypes {
			result.WriteString(fmt.Sprintf("### %s\n\n", st.Name))
			result.WriteString(fmt.Sprintf("**Type:** %s\n\n", st.Type))
			if st.Properties != nil {
				propsJSON, _ := json.MarshalIndent(st.Properties, "", "  ")
				result.WriteString(fmt.Sprintf("```json\n%s\n```\n\n", string(propsJSON)))
			}
		}
	}

	return result.String()
}
