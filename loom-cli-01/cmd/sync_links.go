// Package cmd provides CLI commands for loom-cli.
//
// This file implements the sync-links command.
// Implements: l2/interface-contracts.md IC-SYN-001
// See: l2/sequence-design.md SEQ-SYN-001
package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"loom-cli/internal/domain"
)

// runSyncLinks implements the sync-links command.
// Implements: IC-SYN-001
//
// Actions:
//   - Parse document IDs
//   - Find missing back-references
//   - Add to Traceability sections
func runSyncLinks(args []string) int {
	fs := flag.NewFlagSet("sync-links", flag.ExitOnError)

	inputDir := fs.String("input-dir", ".", "Directory containing documents")
	dryRun := fs.Bool("dry-run", false, "Show changes without applying")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Read all documents
	docs, err := readAllDocuments(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading documents: %v\n", err)
		return domain.ExitCodeError
	}

	// Build reference graph
	// id -> [files that define it]
	definitions := make(map[string][]string)
	// id -> [files that reference it]
	references := make(map[string][]string)

	for _, doc := range docs {
		for _, id := range doc.IDs {
			definitions[id] = append(definitions[id], doc.Path)
		}
		for _, ref := range doc.Refs {
			references[ref] = append(references[ref], doc.Path)
		}
	}

	// Find missing back-references
	type MissingLink struct {
		File   string
		ID     string
		Source string
	}
	var missing []MissingLink

	for id, defFiles := range definitions {
		refFiles := references[id]
		for _, refFile := range refFiles {
			// Check if the defining file references back
			for _, defFile := range defFiles {
				if refFile == defFile {
					continue
				}
				// Check if defFile references refFile's IDs
				found := false
				for _, doc := range docs {
					if doc.Path == defFile {
						for _, ref := range doc.Refs {
							for _, doc2 := range docs {
								if doc2.Path == refFile {
									for _, id2 := range doc2.IDs {
										if ref == id2 {
											found = true
											break
										}
									}
								}
							}
						}
					}
				}
				if !found {
					missing = append(missing, MissingLink{
						File:   defFile,
						ID:     id,
						Source: refFile,
					})
				}
			}
		}
	}

	if len(missing) == 0 {
		fmt.Println("All links are synchronized.")
		return domain.ExitCodeSuccess
	}

	fmt.Printf("Found %d missing back-references:\n", len(missing))
	for _, m := range missing {
		fmt.Printf("  - %s should reference %s (referenced by %s)\n", m.File, m.ID, m.Source)
	}

	if *dryRun {
		fmt.Println("\n(dry-run mode - no changes made)")
		return domain.ExitCodeSuccess
	}

	// Apply fixes
	fixed := 0
	for _, m := range missing {
		if addBackReference(m.File, m.ID, m.Source) {
			fixed++
		}
	}

	fmt.Printf("\nFixed %d/%d links.\n", fixed, len(missing))
	return domain.ExitCodeSuccess
}

// addBackReference adds a back-reference to a file's traceability section.
func addBackReference(file, id, source string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		return false
	}

	contentStr := string(content)

	// Find traceability section
	traceRe := regexp.MustCompile(`(?i)(## Related Documents|## Traceability|## References)\n`)
	loc := traceRe.FindStringIndex(contentStr)

	if loc == nil {
		// No traceability section, add one at the end
		contentStr += fmt.Sprintf("\n\n## Related Documents\n\n- Referenced by: [%s](%s)\n", id, source)
	} else {
		// Insert after section header
		insertPos := loc[1]
		// Find end of section (next ## or end of file)
		nextSection := strings.Index(contentStr[insertPos:], "\n## ")
		if nextSection == -1 {
			nextSection = len(contentStr) - insertPos
		}

		// Check if already referenced
		sectionEnd := insertPos + nextSection
		if strings.Contains(contentStr[insertPos:sectionEnd], id) {
			return false // Already referenced
		}

		// Insert reference
		ref := fmt.Sprintf("- Referenced by: %s\n", id)
		contentStr = contentStr[:insertPos] + ref + contentStr[insertPos:]
	}

	if err := os.WriteFile(file, []byte(contentStr), 0644); err != nil {
		return false
	}

	return true
}
