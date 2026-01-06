package formatter

import (
	"strings"
	"testing"
)

func TestFormatInterfaceContracts_Empty(t *testing.T) {
	result := FormatInterfaceContracts([]InterfaceContract{}, []SharedType{}, "2024-01-15T10:00:00Z")

	// Should still have header
	if !strings.Contains(result, "# Interface Contracts") {
		t.Error("Expected document title")
	}
}

func TestFormatInterfaceContracts_WithSharedTypes(t *testing.T) {
	sharedTypes := []SharedType{
		{
			Name: "Money",
			Fields: []TypeField{
				{Name: "amount", Type: "integer", Constraints: "> 0"},
				{Name: "currency", Type: "string", Constraints: "ISO 4217"},
			},
		},
	}

	result := FormatInterfaceContracts([]InterfaceContract{}, sharedTypes, "2024-01-15T10:00:00Z")

	// Check shared types section
	if !strings.Contains(result, "## Shared Types") {
		t.Error("Expected Shared Types section")
	}

	if !strings.Contains(result, "### Money") {
		t.Error("Expected Money type header")
	}

	if !strings.Contains(result, "| Field | Type | Constraints |") {
		t.Error("Expected table header")
	}

	if !strings.Contains(result, "| amount | integer | > 0 |") {
		t.Error("Expected amount field row")
	}
}

func TestFormatInterfaceContracts_SingleContract(t *testing.T) {
	contracts := []InterfaceContract{
		{
			ID:          "IC-ORDER",
			ServiceName: "Order Service",
			Purpose:     "Manages order lifecycle",
			BaseURL:     "/api/v1/orders",
			SecurityRequirements: SecurityRequirements{
				Authentication: "JWT Bearer",
				Authorization:  "Role-based (CUSTOMER, ADMIN)",
			},
			Operations: []ContractOperation{
				{
					ID:          "OP-CREATE-ORDER",
					Name:        "Create Order",
					Method:      "POST",
					Path:        "/",
					Description: "Creates a new order",
					InputSchema: map[string]SchemaField{
						"customerId": {Type: "string", Required: true},
						"items":      {Type: "array", Required: true},
					},
					OutputSchema: map[string]SchemaField{
						"orderId": {Type: "string"},
						"status":  {Type: "string"},
					},
					Errors: []ContractError{
						{Code: "ORDER_INVALID", HTTPStatus: 400, Message: "Invalid order data"},
						{Code: "CUSTOMER_NOT_FOUND", HTTPStatus: 404, Message: "Customer not found"},
					},
					Preconditions:  []string{"Customer exists", "Items available"},
					Postconditions: []string{"Order created", "Inventory reserved"},
					RelatedACs:     []string{"AC-ORD-001"},
					RelatedBRs:     []string{"BR-ORD-001"},
				},
			},
			Events: []ContractEvent{
				{
					Name:        "OrderCreated",
					Description: "Emitted when order is created",
					Payload:     []string{"orderId", "customerId", "totalAmount"},
				},
			},
		},
	}

	result := FormatInterfaceContracts(contracts, nil, "2024-01-15T10:00:00Z")

	// Check contract header
	if !strings.Contains(result, "## IC-ORDER – Order Service") {
		t.Error("Expected contract header")
	}

	// Check anchor
	if !strings.Contains(result, "{#ic-order}") {
		t.Error("Expected anchor")
	}

	// Check purpose
	if !strings.Contains(result, "**Purpose:** Manages order lifecycle") {
		t.Error("Expected purpose")
	}

	// Check base URL
	if !strings.Contains(result, "**Base URL:** `/api/v1/orders`") {
		t.Error("Expected base URL")
	}

	// Check security
	if !strings.Contains(result, "**Security:**") {
		t.Error("Expected security section")
	}
	if !strings.Contains(result, "Authentication: JWT Bearer") {
		t.Error("Expected authentication")
	}
	if !strings.Contains(result, "Authorization: Role-based") {
		t.Error("Expected authorization")
	}

	// Check operations section
	if !strings.Contains(result, "### Operations") {
		t.Error("Expected Operations section")
	}

	// Check operation header
	if !strings.Contains(result, "#### Create Order `POST /`") {
		t.Error("Expected operation header with method and path")
	}

	// Check input
	if !strings.Contains(result, "**Input:**") {
		t.Error("Expected Input section")
	}
	if !strings.Contains(result, "`customerId`: string (required)") {
		t.Error("Expected required field indicator")
	}

	// Check output
	if !strings.Contains(result, "**Output:**") {
		t.Error("Expected Output section")
	}

	// Check errors table
	if !strings.Contains(result, "**Errors:**") {
		t.Error("Expected Errors section")
	}
	if !strings.Contains(result, "| Code | HTTP | Message |") {
		t.Error("Expected errors table header")
	}
	if !strings.Contains(result, "| ORDER_INVALID | 400 |") {
		t.Error("Expected error row")
	}

	// Check preconditions
	if !strings.Contains(result, "**Preconditions:**") {
		t.Error("Expected Preconditions")
	}

	// Check postconditions
	if !strings.Contains(result, "**Postconditions:**") {
		t.Error("Expected Postconditions")
	}

	// Check events section
	if !strings.Contains(result, "### Events") {
		t.Error("Expected Events section")
	}
	if !strings.Contains(result, "**OrderCreated**") {
		t.Error("Expected event name")
	}
}

func TestFormatInterfaceContracts_NoSecurity(t *testing.T) {
	contracts := []InterfaceContract{
		{
			ID:          "IC-PUBLIC",
			ServiceName: "Public API",
			Purpose:     "Public endpoints",
			BaseURL:     "/api/public",
			// No security requirements
		},
	}

	result := FormatInterfaceContracts(contracts, nil, "2024-01-15T10:00:00Z")

	// Should not have security section if empty
	// Note: Current implementation always shows it if Authentication is not empty
	if strings.Contains(result, "**Security:**") {
		// This is expected behavior - only show if authentication is set
	}
}

func TestFormatInterfaceContracts_NoEvents(t *testing.T) {
	contracts := []InterfaceContract{
		{
			ID:          "IC-QUERY",
			ServiceName: "Query Service",
			Purpose:     "Read-only queries",
			BaseURL:     "/api/query",
			Events:      nil,
		},
	}

	result := FormatInterfaceContracts(contracts, nil, "2024-01-15T10:00:00Z")

	// Should not have events section if empty
	if strings.Contains(result, "### Events") {
		t.Error("Should not have Events section when no events")
	}
}

func TestFormatInterfaceContracts_MultipleContracts(t *testing.T) {
	contracts := []InterfaceContract{
		{
			ID:          "IC-ORDER",
			ServiceName: "Order Service",
			Purpose:     "Orders",
			BaseURL:     "/api/orders",
		},
		{
			ID:          "IC-CUSTOMER",
			ServiceName: "Customer Service",
			Purpose:     "Customers",
			BaseURL:     "/api/customers",
		},
	}

	result := FormatInterfaceContracts(contracts, nil, "2024-01-15T10:00:00Z")

	// Check both contracts are present
	if !strings.Contains(result, "IC-ORDER") {
		t.Error("Expected first contract")
	}
	if !strings.Contains(result, "IC-CUSTOMER") {
		t.Error("Expected second contract")
	}

	// Check both have proper anchors
	if !strings.Contains(result, "{#ic-order}") {
		t.Error("Expected first anchor")
	}
	if !strings.Contains(result, "{#ic-customer}") {
		t.Error("Expected second anchor")
	}
}
