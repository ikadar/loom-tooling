package formatter

import (
	"fmt"
	"strings"
)

// Frontmatter represents YAML frontmatter metadata for derived documents
type Frontmatter struct {
	Title      string   // Document title
	Generated  string   // Timestamp in RFC3339 format
	Status     string   // Document status (default: "draft")
	Level      string   // Derivation level (L1, L2, L3)
	SourceDocs []string // Source document paths
	CLIVersion string   // loom-cli version
}

// DefaultFrontmatter creates a Frontmatter with default values
func DefaultFrontmatter(title, timestamp, level string) Frontmatter {
	return Frontmatter{
		Title:      title,
		Generated:  timestamp,
		Status:     "draft",
		Level:      level,
		CLIVersion: "0.3.0",
	}
}

// FormatFrontmatter formats the frontmatter as YAML
func FormatFrontmatter(fm Frontmatter) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: \"%s\"\n", fm.Title))
	sb.WriteString(fmt.Sprintf("generated: %s\n", fm.Generated))
	sb.WriteString(fmt.Sprintf("status: %s\n", fm.Status))
	sb.WriteString(fmt.Sprintf("level: %s\n", fm.Level))

	if len(fm.SourceDocs) > 0 {
		sb.WriteString("source:\n")
		for _, src := range fm.SourceDocs {
			sb.WriteString(fmt.Sprintf("  - %s\n", src))
		}
	}

	sb.WriteString(fmt.Sprintf("loom-cli-version: %s\n", fm.CLIVersion))
	sb.WriteString("---\n\n")

	return sb.String()
}

// FormatHeaderWithFrontmatter creates a document header with YAML frontmatter
func FormatHeaderWithFrontmatter(fm Frontmatter) string {
	var sb strings.Builder
	sb.WriteString(FormatFrontmatter(fm))
	sb.WriteString(fmt.Sprintf("# %s\n\n", fm.Title))
	return sb.String()
}
