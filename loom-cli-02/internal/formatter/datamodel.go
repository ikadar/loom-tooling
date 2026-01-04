// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatDataModel formats data model as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: initial-data-model.md
func FormatDataModel(tables []DataTable, enums []DataEnum) string {
	var sb strings.Builder

	sb.WriteString(FormatFrontmatter("Initial Data Model", "L2"))
	sb.WriteString("# Initial Data Model\n\n")

	// Enums section
	if len(enums) > 0 {
		sb.WriteString("## Enum Types\n\n")
		for _, enum := range enums {
			sb.WriteString(formatEnum(enum))
		}
		sb.WriteString("---\n\n")
	}

	// Tables section
	sb.WriteString("## Tables\n\n")
	for _, table := range tables {
		sb.WriteString(formatTable(table))
	}

	return sb.String()
}

func formatEnum(enum DataEnum) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### %s\n\n", enum.Name))
	sb.WriteString("```sql\n")
	sb.WriteString(fmt.Sprintf("CREATE TYPE %s AS ENUM (\n", enum.Name))
	for i, val := range enum.Values {
		comma := ","
		if i == len(enum.Values)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    '%s'%s\n", val, comma))
	}
	sb.WriteString(");\n")
	sb.WriteString("```\n\n")

	return sb.String()
}

func formatTable(table DataTable) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("## %s – %s\n\n", table.ID, table.Name))
	sb.WriteString(FormatAnchor(table.ID))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("**Aggregate:** %s\n\n", table.Aggregate))
	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", table.Purpose))

	// Fields table
	sb.WriteString("### Fields\n\n")
	sb.WriteString("| Name | Type | Nullable | Default | Constraints |\n")
	sb.WriteString("|------|------|----------|---------|-------------|\n")
	for _, field := range table.Fields {
		nullable := "NO"
		if field.Nullable {
			nullable = "YES"
		}
		defaultVal := field.Default
		if defaultVal == "" {
			defaultVal = "-"
		}
		constraints := field.Constraints
		if constraints == "" {
			constraints = "-"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			field.Name, field.Type, nullable, defaultVal, constraints))
	}
	sb.WriteString("\n")

	// Primary Key
	sb.WriteString("### Primary Key\n\n")
	sb.WriteString(fmt.Sprintf("- **Type:** %s\n", table.PrimaryKey.Type))
	sb.WriteString(fmt.Sprintf("- **Columns:** %s\n\n", strings.Join(table.PrimaryKey.Columns, ", ")))

	// Indexes
	if len(table.Indexes) > 0 {
		sb.WriteString("### Indexes\n\n")
		sb.WriteString("| Name | Columns | Unique |\n")
		sb.WriteString("|------|---------|--------|\n")
		for _, idx := range table.Indexes {
			unique := "No"
			if idx.Unique {
				unique = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				idx.Name, strings.Join(idx.Columns, ", "), unique))
		}
		sb.WriteString("\n")
	}

	// Foreign Keys
	if len(table.ForeignKeys) > 0 {
		sb.WriteString("### Foreign Keys\n\n")
		sb.WriteString("| Name | Columns | References | On Delete | On Update |\n")
		sb.WriteString("|------|---------|------------|-----------|----------|\n")
		for _, fk := range table.ForeignKeys {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				fk.Name, strings.Join(fk.Columns, ", "), fk.References,
				fk.OnDelete, fk.OnUpdate))
		}
		sb.WriteString("\n")
	}

	// Check Constraints
	if len(table.CheckConstraints) > 0 {
		sb.WriteString("### Check Constraints\n\n")
		sb.WriteString("| Name | Expression |\n")
		sb.WriteString("|------|------------|\n")
		for _, cc := range table.CheckConstraints {
			sb.WriteString(fmt.Sprintf("| %s | `%s` |\n", cc.Name, cc.Expression))
		}
		sb.WriteString("\n")
	}

	// SQL DDL
	sb.WriteString("### DDL\n\n")
	sb.WriteString("```sql\n")
	sb.WriteString(generateTableDDL(table))
	sb.WriteString("```\n\n")

	sb.WriteString("---\n\n")

	return sb.String()
}

// generateTableDDL generates SQL CREATE TABLE statement.
func generateTableDDL(table DataTable) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", table.Name))

	// Fields
	for i, field := range table.Fields {
		nullStr := "NOT NULL"
		if field.Nullable {
			nullStr = "NULL"
		}
		defaultStr := ""
		if field.Default != "" {
			defaultStr = " DEFAULT " + field.Default
		}
		constraintStr := ""
		if field.Constraints != "" {
			constraintStr = " " + field.Constraints
		}

		comma := ","
		if i == len(table.Fields)-1 && len(table.PrimaryKey.Columns) == 0 &&
			len(table.ForeignKeys) == 0 && len(table.CheckConstraints) == 0 {
			comma = ""
		}

		sb.WriteString(fmt.Sprintf("    %s %s %s%s%s%s\n",
			field.Name, field.Type, nullStr, defaultStr, constraintStr, comma))
	}

	// Primary Key
	if len(table.PrimaryKey.Columns) > 0 {
		comma := ","
		if len(table.ForeignKeys) == 0 && len(table.CheckConstraints) == 0 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    PRIMARY KEY (%s)%s\n",
			strings.Join(table.PrimaryKey.Columns, ", "), comma))
	}

	// Foreign Keys
	for i, fk := range table.ForeignKeys {
		comma := ","
		if i == len(table.ForeignKeys)-1 && len(table.CheckConstraints) == 0 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s ON DELETE %s ON UPDATE %s%s\n",
			fk.Name, strings.Join(fk.Columns, ", "), fk.References,
			fk.OnDelete, fk.OnUpdate, comma))
	}

	// Check Constraints
	for i, cc := range table.CheckConstraints {
		comma := ","
		if i == len(table.CheckConstraints)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    CONSTRAINT %s CHECK (%s)%s\n",
			cc.Name, cc.Expression, comma))
	}

	sb.WriteString(");\n")

	// Indexes (separate statements)
	for _, idx := range table.Indexes {
		uniqueStr := ""
		if idx.Unique {
			uniqueStr = "UNIQUE "
		}
		sb.WriteString(fmt.Sprintf("\nCREATE %sINDEX %s ON %s (%s);\n",
			uniqueStr, idx.Name, table.Name, strings.Join(idx.Columns, ", ")))
	}

	return sb.String()
}
