package formatter

import (
	"strings"
	"testing"
)

func TestFormatDataModel_Empty(t *testing.T) {
	result := FormatDataModel([]DataTable{}, []DataEnum{}, "2024-01-15T10:00:00Z")

	if !strings.Contains(result, "# Initial Data Model") {
		t.Error("Expected document title")
	}

	// Should have ER diagram section even if empty
	if !strings.Contains(result, "## Entity-Relationship Diagram") {
		t.Error("Expected ER diagram section")
	}

	// Should have mermaid code block
	if !strings.Contains(result, "```mermaid") {
		t.Error("Expected mermaid code block")
	}
}

func TestFormatDataModel_WithEnums(t *testing.T) {
	enums := []DataEnum{
		{
			Name:   "OrderStatus",
			Values: []string{"DRAFT", "PENDING", "CONFIRMED", "SHIPPED", "DELIVERED", "CANCELLED"},
		},
		{
			Name:   "PaymentMethod",
			Values: []string{"CREDIT_CARD", "DEBIT_CARD", "BANK_TRANSFER", "CASH"},
		},
	}

	result := FormatDataModel([]DataTable{}, enums, "2024-01-15T10:00:00Z")

	// Check enums section
	if !strings.Contains(result, "## Enumerations") {
		t.Error("Expected Enumerations section")
	}

	if !strings.Contains(result, "### OrderStatus") {
		t.Error("Expected OrderStatus enum")
	}

	if !strings.Contains(result, "### PaymentMethod") {
		t.Error("Expected PaymentMethod enum")
	}

	if !strings.Contains(result, "DRAFT") {
		t.Error("Expected DRAFT value")
	}
}

func TestFormatDataModel_SingleTable(t *testing.T) {
	tables := []DataTable{
		{
			ID:        "TBL-ORDERS",
			Name:      "orders",
			Aggregate: "Order",
			Purpose:   "Stores order header information",
			Fields: []DataField{
				{Name: "id", Type: "UUID", Constraints: "NOT NULL", Default: "gen_random_uuid()"},
				{Name: "customer_id", Type: "UUID", Constraints: "NOT NULL", Default: ""},
				{Name: "status", Type: "order_status", Constraints: "NOT NULL", Default: "'DRAFT'"},
				{Name: "total_cents", Type: "INTEGER", Constraints: "NOT NULL, CHECK > 0", Default: "0"},
				{Name: "created_at", Type: "TIMESTAMPTZ", Constraints: "NOT NULL", Default: "NOW()"},
			},
			PrimaryKey: DataPrimaryKey{
				Columns: []string{"id"},
			},
			Indexes: []DataIndex{
				{Name: "idx_orders_customer", Columns: []string{"customer_id"}},
				{Name: "idx_orders_status", Columns: []string{"status", "created_at"}},
			},
			ForeignKeys: []DataForeignKey{
				{Columns: []string{"customer_id"}, References: "customers(id)", OnDelete: "RESTRICT"},
			},
			CheckConstraints: []DataConstraint{
				{Name: "chk_total_positive", Expression: "total_cents >= 0"},
			},
		},
	}

	result := FormatDataModel(tables, nil, "2024-01-15T10:00:00Z")

	// Check table header
	if !strings.Contains(result, "### TBL-ORDERS – orders") {
		t.Error("Expected table header")
	}

	// Check anchor
	if !strings.Contains(result, "{#tbl-orders}") {
		t.Error("Expected anchor")
	}

	// Check aggregate
	if !strings.Contains(result, "**Aggregate:** Order") {
		t.Error("Expected aggregate reference")
	}

	// Check purpose
	if !strings.Contains(result, "**Purpose:** Stores order header") {
		t.Error("Expected purpose")
	}

	// Check columns table
	if !strings.Contains(result, "**Columns:**") {
		t.Error("Expected Columns section")
	}
	if !strings.Contains(result, "| Name | Type | Constraints | Default |") {
		t.Error("Expected columns table header")
	}
	if !strings.Contains(result, "| id | UUID | NOT NULL | gen_random_uuid() |") {
		t.Error("Expected id column row")
	}

	// Check primary key
	if !strings.Contains(result, "**Primary Key:** [id]") {
		t.Error("Expected primary key")
	}

	// Check indexes
	if !strings.Contains(result, "**Indexes:**") {
		t.Error("Expected Indexes section")
	}
	if !strings.Contains(result, "`idx_orders_customer`") {
		t.Error("Expected customer index")
	}

	// Check foreign keys
	if !strings.Contains(result, "**Foreign Keys:**") {
		t.Error("Expected Foreign Keys section")
	}
	if !strings.Contains(result, "customers(id)") {
		t.Error("Expected FK reference")
	}
	if !strings.Contains(result, "ON DELETE RESTRICT") {
		t.Error("Expected ON DELETE clause")
	}

	// Check constraints
	if !strings.Contains(result, "**Constraints:**") {
		t.Error("Expected Constraints section")
	}
	if !strings.Contains(result, "`chk_total_positive`") {
		t.Error("Expected constraint name")
	}

	// Check ER diagram
	if !strings.Contains(result, "```mermaid") {
		t.Error("Expected mermaid diagram")
	}
	if !strings.Contains(result, "erDiagram") {
		t.Error("Expected erDiagram keyword")
	}
	if !strings.Contains(result, "orders {") {
		t.Error("Expected orders table in diagram")
	}
}

func TestFormatDataModel_ERDiagramRelationships(t *testing.T) {
	tables := []DataTable{
		{
			ID:   "TBL-ORDERS",
			Name: "orders",
			Fields: []DataField{
				{Name: "id", Type: "UUID"},
				{Name: "customer_id", Type: "UUID"},
			},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
			ForeignKeys: []DataForeignKey{
				{Columns: []string{"customer_id"}, References: "customers(id)", OnDelete: "RESTRICT"},
			},
		},
		{
			ID:   "TBL-ORDER-ITEMS",
			Name: "order_items",
			Fields: []DataField{
				{Name: "id", Type: "UUID"},
				{Name: "order_id", Type: "UUID"},
			},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
			ForeignKeys: []DataForeignKey{
				{Columns: []string{"order_id"}, References: "orders(id)", OnDelete: "CASCADE"},
			},
		},
	}

	result := FormatDataModel(tables, nil, "2024-01-15T10:00:00Z")

	// Check relationships in ER diagram
	if !strings.Contains(result, "customers ||--o{ orders") {
		t.Error("Expected customers -> orders relationship")
	}
	if !strings.Contains(result, "orders ||--o{ order_items") {
		t.Error("Expected orders -> order_items relationship")
	}
}

func TestFormatDataModel_MinimalTable(t *testing.T) {
	tables := []DataTable{
		{
			ID:         "TBL-SIMPLE",
			Name:       "simple_table",
			Aggregate:  "Simple",
			Purpose:    "A simple table",
			Fields:     []DataField{{Name: "id", Type: "UUID"}},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
			// No indexes, foreign keys, or check constraints
		},
	}

	result := FormatDataModel(tables, nil, "2024-01-15T10:00:00Z")

	// Should have basic structure
	if !strings.Contains(result, "TBL-SIMPLE") {
		t.Error("Expected table ID")
	}

	// Should not have optional sections
	if strings.Contains(result, "**Indexes:**") {
		t.Error("Should not have Indexes section when empty")
	}
	if strings.Contains(result, "**Foreign Keys:**") {
		t.Error("Should not have Foreign Keys section when empty")
	}
	if strings.Contains(result, "**Constraints:**") {
		t.Error("Should not have Constraints section when empty")
	}
}

func TestFormatDataModel_MultipleTables(t *testing.T) {
	tables := []DataTable{
		{
			ID:         "TBL-ORDERS",
			Name:       "orders",
			Fields:     []DataField{{Name: "id", Type: "UUID"}},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
		},
		{
			ID:         "TBL-CUSTOMERS",
			Name:       "customers",
			Fields:     []DataField{{Name: "id", Type: "UUID"}},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
		},
	}

	result := FormatDataModel(tables, nil, "2024-01-15T10:00:00Z")

	// Check both tables are present
	if !strings.Contains(result, "TBL-ORDERS") {
		t.Error("Expected orders table")
	}
	if !strings.Contains(result, "TBL-CUSTOMERS") {
		t.Error("Expected customers table")
	}

	// Check both are in ER diagram
	if !strings.Contains(result, "orders {") {
		t.Error("Expected orders in ER diagram")
	}
	if !strings.Contains(result, "customers {") {
		t.Error("Expected customers in ER diagram")
	}
}

func TestFormatDataModel_NoEnums(t *testing.T) {
	tables := []DataTable{
		{
			ID:         "TBL-TEST",
			Name:       "test",
			Fields:     []DataField{{Name: "id", Type: "UUID"}},
			PrimaryKey: DataPrimaryKey{Columns: []string{"id"}},
		},
	}

	result := FormatDataModel(tables, nil, "2024-01-15T10:00:00Z")

	// Should not have Enumerations section
	if strings.Contains(result, "## Enumerations") {
		t.Error("Should not have Enumerations section when no enums")
	}
}
