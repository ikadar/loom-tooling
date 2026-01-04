// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements anchor/link formatting utilities.
//
// Implements: l2/package-structure.md PKG-007
package formatter

import (
	"fmt"
	"strings"
)

// FormatAnchor creates a markdown anchor link.
func FormatAnchor(id string) string {
	return fmt.Sprintf("[%s](#%s)", id, strings.ToLower(strings.ReplaceAll(id, "_", "-")))
}

// FormatCrossReference creates a cross-reference link to another document.
func FormatCrossReference(id, document string) string {
	anchor := strings.ToLower(strings.ReplaceAll(id, "_", "-"))
	return fmt.Sprintf("[%s](%s#%s)", id, document, anchor)
}

// FormatSourceRefs formats source references as a comma-separated list.
func FormatSourceRefs(refs []string) string {
	if len(refs) == 0 {
		return ""
	}
	return strings.Join(refs, ", ")
}

// FormatIDList formats a list of IDs as markdown links.
func FormatIDList(ids []string) string {
	if len(ids) == 0 {
		return "-"
	}

	links := make([]string, len(ids))
	for i, id := range ids {
		links[i] = FormatAnchor(id)
	}
	return strings.Join(links, ", ")
}
