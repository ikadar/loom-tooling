// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"strings"
)

// FormatAnchor creates a markdown anchor link.
//
// Implements: l2/package-structure.md PKG-007
func FormatAnchor(id string) string {
	return fmt.Sprintf("<a id=\"%s\"></a>", id)
}

// FormatReference creates a reference to another document.
//
// Implements: l2/package-structure.md PKG-007
func FormatReference(id string, docPath string) string {
	// Generate anchor from ID (lowercase, hyphens)
	anchor := strings.ToLower(id)
	return fmt.Sprintf("[%s](%s#%s)", id, docPath, anchor)
}

// FormatTraceability creates a traceability section.
//
// Implements: l2/package-structure.md PKG-007
func FormatTraceability(sources []string, decisions []string) string {
	var sb strings.Builder

	sb.WriteString("## Traceability\n\n")

	if len(sources) > 0 {
		sb.WriteString("**Sources:**\n")
		for _, src := range sources {
			sb.WriteString(fmt.Sprintf("- %s\n", src))
		}
		sb.WriteString("\n")
	}

	if len(decisions) > 0 {
		sb.WriteString("**Decisions:**\n")
		for _, dec := range decisions {
			sb.WriteString(fmt.Sprintf("- %s\n", dec))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
