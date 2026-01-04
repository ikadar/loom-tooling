// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

import (
	"fmt"
	"time"
)

// FormatFrontmatter generates YAML frontmatter for documents.
//
// Implements: l2/package-structure.md PKG-007
func FormatFrontmatter(title string, level string) string {
	return fmt.Sprintf(`---
title: "%s"
generated: %s
status: draft
level: %s
---

`, title, time.Now().Format(time.RFC3339), level)
}
