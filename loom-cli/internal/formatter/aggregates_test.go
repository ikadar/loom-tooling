package formatter

import (
	"strings"
	"testing"
)

func TestFormatAggregateDesign_Empty(t *testing.T) {
	result := FormatAggregateDesign([]AggregateDesign{}, "2024-01-15T10:00:00Z")

	if !strings.Contains(result, "# Aggregate Design") {
		t.Error("Expected document title")
	}
}

func TestFormatAggregateDesign_SingleAggregate(t *testing.T) {
	aggregates := []AggregateDesign{
		{
			ID:      "AGG-ORDER",
			Name:    "Order Aggregate",
			Purpose: "Manages order lifecycle and consistency",
			Root: AggRoot{
				Entity:   "Order",
				Identity: "OrderId (UUID)",
				Attributes: []AggAttribute{
					{Name: "orderId", Type: "UUID", Mutable: false},
					{Name: "status", Type: "OrderStatus", Mutable: true},
					{Name: "totalAmount", Type: "Money", Mutable: true},
				},
			},
			Invariants: []AggInvariant{
				{
					ID:          "INV-ORD-001",
					Rule:        "Order total must be positive",
					Enforcement: "Immediate on add/update item",
				},
			},
			Entities: []AggEntity{
				{
					Name:     "OrderItem",
					Identity: "LineNumber",
					Purpose:  "Represents a single item in the order",
				},
			},
			ValueObjects: []string{"Money", "Address", "OrderStatus"},
			Behaviors: []AggBehavior{
				{
					Name:           "AddItem",
					Command:        "AddItemCommand",
					Preconditions:  []string{"Order is draft"},
					Postconditions: []string{"Item added", "Total recalculated"},
					Emits:          "ItemAddedEvent",
				},
			},
			Events: []AggEvent{
				{
					Name:    "OrderCreated",
					Payload: []string{"orderId", "customerId", "createdAt"},
				},
			},
			Repository: AggRepository{
				Name:         "OrderRepository",
				LoadStrategy: "Eager (include items)",
				Concurrency:  "Optimistic locking via version",
			},
		},
	}

	result := FormatAggregateDesign(aggregates, "2024-01-15T10:00:00Z")

	// Check header
	if !strings.Contains(result, "## AGG-ORDER – Order Aggregate") {
		t.Error("Expected aggregate header")
	}

	// Check anchor
	if !strings.Contains(result, "{#agg-order}") {
		t.Error("Expected anchor")
	}

	// Check purpose
	if !strings.Contains(result, "**Purpose:** Manages order lifecycle") {
		t.Error("Expected purpose")
	}

	// Check root
	if !strings.Contains(result, "### Aggregate Root: Order") {
		t.Error("Expected aggregate root section")
	}
	if !strings.Contains(result, "**Identity:** OrderId (UUID)") {
		t.Error("Expected identity")
	}

	// Check attributes table
	if !strings.Contains(result, "| Name | Type | Mutable |") {
		t.Error("Expected attributes table header")
	}
	if !strings.Contains(result, "| orderId | UUID | false |") {
		t.Error("Expected orderId attribute row")
	}

	// Check invariants
	if !strings.Contains(result, "### Invariants") {
		t.Error("Expected invariants section")
	}
	if !strings.Contains(result, "**INV-ORD-001**") {
		t.Error("Expected invariant ID")
	}
	if !strings.Contains(result, "Enforcement: Immediate") {
		t.Error("Expected enforcement")
	}

	// Check child entities
	if !strings.Contains(result, "### Child Entities") {
		t.Error("Expected child entities section")
	}
	if !strings.Contains(result, "#### OrderItem") {
		t.Error("Expected OrderItem entity")
	}

	// Check value objects
	if !strings.Contains(result, "### Value Objects") {
		t.Error("Expected value objects section")
	}
	if !strings.Contains(result, "Money") {
		t.Error("Expected Money value object")
	}

	// Check behaviors table
	if !strings.Contains(result, "### Behaviors") {
		t.Error("Expected behaviors section")
	}
	if !strings.Contains(result, "| Command | Pre | Post | Emits |") {
		t.Error("Expected behaviors table header")
	}
	if !strings.Contains(result, "AddItemCommand") {
		t.Error("Expected AddItemCommand in behaviors")
	}

	// Check events
	if !strings.Contains(result, "### Events") {
		t.Error("Expected events section")
	}
	if !strings.Contains(result, "**OrderCreated**") {
		t.Error("Expected OrderCreated event")
	}

	// Check repository
	if !strings.Contains(result, "### Repository: OrderRepository") {
		t.Error("Expected repository section")
	}
	if !strings.Contains(result, "Load Strategy: Eager") {
		t.Error("Expected load strategy")
	}
	if !strings.Contains(result, "Concurrency: Optimistic") {
		t.Error("Expected concurrency strategy")
	}
}

func TestFormatAggregateDesign_MinimalAggregate(t *testing.T) {
	aggregates := []AggregateDesign{
		{
			ID:      "AGG-SIMPLE",
			Name:    "Simple Aggregate",
			Purpose: "Minimal aggregate",
			Root: AggRoot{
				Entity:   "SimpleEntity",
				Identity: "Id",
			},
			Repository: AggRepository{
				Name:         "SimpleRepository",
				LoadStrategy: "Lazy",
				Concurrency:  "None",
			},
			// No invariants, entities, value objects, behaviors, or events
		},
	}

	result := FormatAggregateDesign(aggregates, "2024-01-15T10:00:00Z")

	// Should have basic structure
	if !strings.Contains(result, "AGG-SIMPLE") {
		t.Error("Expected aggregate ID")
	}

	// Should not have optional sections
	if strings.Contains(result, "### Invariants") {
		t.Error("Should not have Invariants section when empty")
	}
	if strings.Contains(result, "### Child Entities") {
		t.Error("Should not have Child Entities section when empty")
	}
	if strings.Contains(result, "### Value Objects") {
		t.Error("Should not have Value Objects section when empty")
	}
	if strings.Contains(result, "### Behaviors") {
		t.Error("Should not have Behaviors section when empty")
	}
	if strings.Contains(result, "### Events") {
		t.Error("Should not have Events section when empty")
	}
}

func TestFormatAggregateDesign_MultipleAggregates(t *testing.T) {
	aggregates := []AggregateDesign{
		{
			ID:   "AGG-ORDER",
			Name: "Order",
			Root: AggRoot{Entity: "Order", Identity: "OrderId"},
			Repository: AggRepository{
				Name:         "OrderRepository",
				LoadStrategy: "Eager",
				Concurrency:  "Optimistic",
			},
		},
		{
			ID:   "AGG-CUSTOMER",
			Name: "Customer",
			Root: AggRoot{Entity: "Customer", Identity: "CustomerId"},
			Repository: AggRepository{
				Name:         "CustomerRepository",
				LoadStrategy: "Lazy",
				Concurrency:  "Pessimistic",
			},
		},
	}

	result := FormatAggregateDesign(aggregates, "2024-01-15T10:00:00Z")

	// Check both aggregates are present
	if !strings.Contains(result, "AGG-ORDER") {
		t.Error("Expected first aggregate")
	}
	if !strings.Contains(result, "AGG-CUSTOMER") {
		t.Error("Expected second aggregate")
	}

	// Check both have anchors
	if !strings.Contains(result, "{#agg-order}") {
		t.Error("Expected first anchor")
	}
	if !strings.Contains(result, "{#agg-customer}") {
		t.Error("Expected second anchor")
	}
}
