// Package formatter provides JSON to Markdown formatting utilities.
//
// This file implements YAML frontmatter generation.
//
// Implements: l2/tech-specs.md TS-FMT-001
package formatter

import (
	"fmt"
	"time"

	"loom-cli/internal/domain"
)

// FormatFrontmatter generates YAML frontmatter for a document.
//
// Implements: TS-FMT-001
func FormatFrontmatter(title, level string) string {
	return fmt.Sprintf(`---
title: "%s"
generated: %s
status: draft
level: %s
---

`, title, time.Now().Format(time.RFC3339), level)
}

// FormatAcceptanceCriteria formats acceptance criteria as markdown.
func FormatAcceptanceCriteria(acs []domain.AcceptanceCriteria) string {
	var result string
	result += FormatFrontmatter("Acceptance Criteria", "L1")
	result += "# Acceptance Criteria\n\n"

	for _, ac := range acs {
		result += fmt.Sprintf("## %s: %s\n\n", ac.ID, ac.Title)
		result += fmt.Sprintf("**Given:** %s\n\n", ac.Given)
		result += fmt.Sprintf("**When:** %s\n\n", ac.When)
		result += fmt.Sprintf("**Then:** %s\n\n", ac.Then)

		if len(ac.ErrorCases) > 0 {
			result += "**Error Cases:**\n"
			for _, e := range ac.ErrorCases {
				result += fmt.Sprintf("- %s\n", e)
			}
			result += "\n"
		}

		if len(ac.SourceRefs) > 0 {
			result += fmt.Sprintf("**Source Refs:** %v\n\n", ac.SourceRefs)
		}

		result += "---\n\n"
	}

	return result
}

// FormatBusinessRules formats business rules as markdown.
func FormatBusinessRules(brs []domain.BusinessRule) string {
	var result string
	result += FormatFrontmatter("Business Rules", "L1")
	result += "# Business Rules\n\n"

	for _, br := range brs {
		result += fmt.Sprintf("## %s: %s\n\n", br.ID, br.Title)
		result += fmt.Sprintf("**Rule:** %s\n\n", br.Rule)
		result += fmt.Sprintf("**Invariant:** %s\n\n", br.Invariant)
		result += fmt.Sprintf("**Enforcement:** %s\n\n", br.Enforcement)

		if br.ErrorCode != "" {
			result += fmt.Sprintf("**Error Code:** %s\n\n", br.ErrorCode)
		}

		if len(br.SourceRefs) > 0 {
			result += fmt.Sprintf("**Source Refs:** %v\n\n", br.SourceRefs)
		}

		result += "---\n\n"
	}

	return result
}

// FormatDecisions formats decisions as markdown.
func FormatDecisions(decisions []domain.Decision) string {
	var result string
	result += FormatFrontmatter("Design Decisions", "L1")
	result += "# Design Decisions\n\n"

	for _, d := range decisions {
		result += fmt.Sprintf("## %s\n\n", d.ID)
		result += fmt.Sprintf("**Question:** %s\n\n", d.Question)
		result += fmt.Sprintf("**Answer:** %s\n\n", d.Answer)
		result += fmt.Sprintf("**Source:** %s\n\n", d.Source)
		result += fmt.Sprintf("**Decided:** %s\n\n", d.DecidedAt.Format(time.RFC3339))
		result += "---\n\n"
	}

	return result
}
