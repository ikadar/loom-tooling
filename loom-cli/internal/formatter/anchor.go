package formatter

import (
	"fmt"
	"strings"
)

// ToAnchor converts an ID to a lowercase anchor (e.g., "TC-AC-CUST-001-P01" -> "tc-ac-cust-001-p01")
func ToAnchor(id string) string {
	return strings.ToLower(id)
}

// ToLink creates a markdown link with anchor
func ToLink(id, file string) string {
	return fmt.Sprintf("[%s](%s#%s)", id, file, ToAnchor(id))
}

// FormatHeader creates a standard document header with title and timestamp
func FormatHeader(title, timestamp string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", timestamp))
	sb.WriteString("---\n\n")
	return sb.String()
}

// FormatSectionHeader creates a section header with anchor
func FormatSectionHeader(level int, id, name string) string {
	prefix := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s – %s {#%s}\n\n", prefix, id, name, ToAnchor(id))
}
