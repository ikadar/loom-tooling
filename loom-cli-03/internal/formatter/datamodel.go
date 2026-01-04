// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements data model formatting.
//
// Implements: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatDataModel formats data model as markdown.
func FormatDataModel(tables []DataTable, enums []DataEnum) string {
	var result strings.Builder

	result.WriteString(FormatFrontmatter("Initial Data Model", "L2"))
	result.WriteString("# Initial Data Model\n\n")

	// Enums
	if len(enums) > 0 {
		result.WriteString("## Enumerations\n\n")
		for _, e := range enums {
			result.WriteString(fmt.Sprintf("### %s\n\n", e.Name))
			result.WriteString("| Value |\n")
			result.WriteString("|-------|\n")
			for _, v := range e.Values {
				result.WriteString(fmt.Sprintf("| %s |\n", v))
			}
			result.WriteString("\n")
		}
	}

	// Tables
	result.WriteString("## Tables\n\n")
	for _, table := range tables {
		result.WriteString(fmt.Sprintf("### %s: %s\n\n", table.ID, table.Name))
		result.WriteString(fmt.Sprintf("**Aggregate:** %s\n\n", table.Aggregate))
		result.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", table.Purpose))

		// Fields
		result.WriteString("#### Fields\n\n")
		result.WriteString("| Field | Type | Nullable | Default | Description |\n")
		result.WriteString("|-------|------|----------|---------|-------------|\n")
		for _, f := range table.Fields {
			nullable := "NO"
			if f.Nullable {
				nullable = "YES"
			}
			defaultVal := f.Default
			if defaultVal == "" {
				defaultVal = "-"
			}
			result.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				f.Name, f.Type, nullable, defaultVal, f.Description))
		}
		result.WriteString("\n")

		// Primary Key
		result.WriteString("#### Primary Key\n\n")
		result.WriteString(fmt.Sprintf("- **Fields:** %s\n", strings.Join(table.PrimaryKey.Fields, ", ")))
		result.WriteString(fmt.Sprintf("- **Type:** %s\n\n", table.PrimaryKey.Type))

		// Indexes
		if len(table.Indexes) > 0 {
			result.WriteString("#### Indexes\n\n")
			result.WriteString("| Name | Fields | Unique |\n")
			result.WriteString("|------|--------|--------|\n")
			for _, idx := range table.Indexes {
				unique := "NO"
				if idx.Unique {
					unique = "YES"
				}
				result.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
					idx.Name, strings.Join(idx.Fields, ", "), unique))
			}
			result.WriteString("\n")
		}

		// Foreign Keys
		if len(table.ForeignKeys) > 0 {
			result.WriteString("#### Foreign Keys\n\n")
			result.WriteString("| Name | Fields | References | On Delete |\n")
			result.WriteString("|------|--------|------------|------------|\n")
			for _, fk := range table.ForeignKeys {
				result.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					fk.Name, strings.Join(fk.Fields, ", "), fk.References, fk.OnDelete))
			}
			result.WriteString("\n")
		}

		// Check Constraints
		if len(table.CheckConstraints) > 0 {
			result.WriteString("#### Check Constraints\n\n")
			for _, cc := range table.CheckConstraints {
				result.WriteString(fmt.Sprintf("- **%s:** `%s`\n", cc.Name, cc.Expression))
			}
			result.WriteString("\n")
		}

		result.WriteString("---\n\n")
	}

	return result.String()
}
