package formatter

import (
	"strings"
	"testing"
)

func TestFormatTechSpecs_Empty(t *testing.T) {
	result := FormatTechSpecs([]TechSpec{}, "2024-01-15T10:00:00Z")

	// Should still have header
	if !strings.Contains(result, "# Technical Specifications") {
		t.Error("Expected document title")
	}

	// Should have frontmatter
	if !strings.Contains(result, "level: L2") {
		t.Error("Expected L2 level in frontmatter")
	}
}

func TestFormatTechSpecs_SingleSpec(t *testing.T) {
	specs := []TechSpec{
		{
			ID:             "TS-ORD-001",
			Name:           "Order Total Validation",
			BRRef:          "BR-ORD-001",
			Rule:           "Order total must be between $1 and $10,000",
			Implementation: "Validate in OrderService.validateOrder()",
			ValidationPoints: []string{
				"POST /orders endpoint",
				"OrderService.addItem()",
			},
			DataRequirements: []DataReq{
				{
					Field:       "total_cents",
					Type:        "integer",
					Constraints: "> 0, < 1000000",
					Source:      "calculated",
				},
			},
			ErrorHandling: []ErrorHandling{
				{
					Condition:  "total < 100 cents",
					ErrorCode:  "ORDER_BELOW_MINIMUM",
					Message:    "Order must be at least $1",
					HTTPStatus: 400,
				},
			},
			RelatedACs: []string{"AC-ORD-001", "AC-ORD-002"},
		},
	}

	result := FormatTechSpecs(specs, "2024-01-15T10:00:00Z")

	// Check ID and name in header
	if !strings.Contains(result, "## TS-ORD-001 – Order Total Validation") {
		t.Error("Expected section header with ID and name")
	}

	// Check anchor
	if !strings.Contains(result, "{#ts-ord-001}") {
		t.Error("Expected lowercase anchor")
	}

	// Check rule
	if !strings.Contains(result, "**Rule:** Order total must be between $1 and $10,000") {
		t.Error("Expected rule content")
	}

	// Check implementation
	if !strings.Contains(result, "**Implementation Approach:**") {
		t.Error("Expected implementation section")
	}

	// Check validation points
	if !strings.Contains(result, "**Validation Points:**") {
		t.Error("Expected validation points section")
	}
	if !strings.Contains(result, "- POST /orders endpoint") {
		t.Error("Expected validation point item")
	}

	// Check data requirements table
	if !strings.Contains(result, "**Data Requirements:**") {
		t.Error("Expected data requirements section")
	}
	if !strings.Contains(result, "| Field | Type | Constraints | Source |") {
		t.Error("Expected data requirements table header")
	}
	if !strings.Contains(result, "| total_cents | integer | > 0, < 1000000 | calculated |") {
		t.Error("Expected data requirements row")
	}

	// Check error handling table
	if !strings.Contains(result, "**Error Handling:**") {
		t.Error("Expected error handling section")
	}
	if !strings.Contains(result, "| Condition | Error Code | Message | HTTP Status |") {
		t.Error("Expected error handling table header")
	}
	if !strings.Contains(result, "ORDER_BELOW_MINIMUM") {
		t.Error("Expected error code in table")
	}
	if !strings.Contains(result, "400") {
		t.Error("Expected HTTP status in table")
	}

	// Check traceability
	if !strings.Contains(result, "**Traceability:**") {
		t.Error("Expected traceability section")
	}
	if !strings.Contains(result, "BR-ORD-001") {
		t.Error("Expected BR reference")
	}
	if !strings.Contains(result, "AC-ORD-001") {
		t.Error("Expected AC reference")
	}
}

func TestFormatTechSpecs_MultipleSpecs(t *testing.T) {
	specs := []TechSpec{
		{
			ID:   "TS-ORD-001",
			Name: "First Spec",
			Rule: "First rule",
		},
		{
			ID:   "TS-ORD-002",
			Name: "Second Spec",
			Rule: "Second rule",
		},
	}

	result := FormatTechSpecs(specs, "2024-01-15T10:00:00Z")

	// Check both specs are present
	if !strings.Contains(result, "TS-ORD-001") {
		t.Error("Expected first spec ID")
	}
	if !strings.Contains(result, "TS-ORD-002") {
		t.Error("Expected second spec ID")
	}

	// Check separators
	separatorCount := strings.Count(result, "---")
	if separatorCount < 3 { // At least: frontmatter end + spec separators
		t.Errorf("Expected at least 3 separators, got %d", separatorCount)
	}
}

func TestFormatTechSpecs_NoOptionalSections(t *testing.T) {
	specs := []TechSpec{
		{
			ID:   "TS-ORD-001",
			Name: "Minimal Spec",
			Rule: "Simple rule",
			// No validation points, data requirements, error handling, or related ACs
		},
	}

	result := FormatTechSpecs(specs, "2024-01-15T10:00:00Z")

	// Should have basic structure
	if !strings.Contains(result, "TS-ORD-001") {
		t.Error("Expected spec ID")
	}
	if !strings.Contains(result, "Simple rule") {
		t.Error("Expected rule")
	}

	// Should not have empty sections
	// (The formatter should handle empty slices gracefully)
}

func TestFormatTechSpecs_LinkFormat(t *testing.T) {
	specs := []TechSpec{
		{
			ID:         "TS-ORD-001",
			Name:       "Test Spec",
			BRRef:      "BR-ORD-001",
			Rule:       "Test rule",
			RelatedACs: []string{"AC-ORD-001"},
		},
	}

	result := FormatTechSpecs(specs, "2024-01-15T10:00:00Z")

	// Check BR link format (should use L1BasePath)
	if !strings.Contains(result, "[BR-ORD-001]") {
		t.Error("Expected BR link text")
	}

	// Check AC link format
	if !strings.Contains(result, "[AC-ORD-001]") {
		t.Error("Expected AC link text")
	}
}
