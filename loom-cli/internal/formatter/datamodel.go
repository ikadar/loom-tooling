package formatter

import (
	"fmt"
	"strings"
)

// FormatDataModel formats data model as markdown
func FormatDataModel(tables []DataTable, enums []DataEnum, timestamp string) string {
	var sb strings.Builder

	sb.WriteString("# Initial Data Model\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", timestamp))
	sb.WriteString("---\n\n")

	// Enums
	if len(enums) > 0 {
		sb.WriteString("## Enumerations\n\n")
		for _, enum := range enums {
			sb.WriteString(fmt.Sprintf("### %s\n\n", enum.Name))
			sb.WriteString(fmt.Sprintf("Values: `%v`\n\n", enum.Values))
		}
		sb.WriteString("---\n\n")
	}

	// Tables
	sb.WriteString("## Tables\n\n")
	for _, tbl := range tables {
		sb.WriteString(formatTable(tbl))
	}

	// ER Diagram
	sb.WriteString("## Entity-Relationship Diagram\n\n")
	sb.WriteString("```mermaid\nerDiagram\n")
	for _, tbl := range tables {
		sb.WriteString(fmt.Sprintf("    %s {\n", tbl.Name))
		for _, field := range tbl.Fields {
			sb.WriteString(fmt.Sprintf("        %s %s\n", field.Type, field.Name))
		}
		sb.WriteString("    }\n")
	}
	// Add relationships
	for _, tbl := range tables {
		for _, fk := range tbl.ForeignKeys {
			// Extract target table from references like "customers(id)"
			target := fk.References
			if idx := strings.Index(target, "("); idx > 0 {
				target = target[:idx]
			}
			sb.WriteString(fmt.Sprintf("    %s ||--o{ %s : \"FK\"\n", target, tbl.Name))
		}
	}
	sb.WriteString("```\n")

	return sb.String()
}

// formatTable formats a single table
func formatTable(tbl DataTable) string {
	var sb strings.Builder

	sb.WriteString(FormatSectionHeader(3, tbl.ID, tbl.Name))
	sb.WriteString(fmt.Sprintf("**Aggregate:** %s\n\n", tbl.Aggregate))
	sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", tbl.Purpose))

	// Fields
	sb.WriteString("**Columns:**\n\n")
	sb.WriteString("| Name | Type | Constraints | Default |\n")
	sb.WriteString("|------|------|-------------|----------|\n")
	for _, field := range tbl.Fields {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", field.Name, field.Type, field.Constraints, field.Default))
	}
	sb.WriteString("\n")

	// Primary Key
	sb.WriteString(fmt.Sprintf("**Primary Key:** %v\n\n", tbl.PrimaryKey.Columns))

	// Indexes
	if len(tbl.Indexes) > 0 {
		sb.WriteString("**Indexes:**\n")
		for _, idx := range tbl.Indexes {
			sb.WriteString(fmt.Sprintf("- `%s`: %v\n", idx.Name, idx.Columns))
		}
		sb.WriteString("\n")
	}

	// Foreign Keys
	if len(tbl.ForeignKeys) > 0 {
		sb.WriteString("**Foreign Keys:**\n")
		for _, fk := range tbl.ForeignKeys {
			sb.WriteString(fmt.Sprintf("- %v → %s (ON DELETE %s)\n", fk.Columns, fk.References, fk.OnDelete))
		}
		sb.WriteString("\n")
	}

	// Check Constraints
	if len(tbl.CheckConstraints) > 0 {
		sb.WriteString("**Constraints:**\n")
		for _, c := range tbl.CheckConstraints {
			sb.WriteString(fmt.Sprintf("- `%s`: %s\n", c.Name, c.Expression))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	return sb.String()
}
