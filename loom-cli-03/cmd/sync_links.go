// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md IC-SYN-001
// See: l2/sequence-design.md SEQ-SYN-001
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"loom-cli/internal/domain"
	"loom-cli/prompts"
)

// Ensure prompts package is imported
var _ = prompts.Derivation

// runSyncLinks handles the sync-links command.
//
// Implements: IC-SYN-001
// Automatically fixes missing bidirectional references.
func runSyncLinks(args []string) int {
	fs := flag.NewFlagSet("sync-links", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "Directory containing documents (required)")
	dryRun := fs.Bool("dry-run", false, "Preview changes without modifying files")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	if *inputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --input-dir is required")
		return domain.ExitCodeError
	}

	// Collect all IDs and their locations
	idLocations := make(map[string][]string) // ID -> list of files where it appears
	fileContents := make(map[string]string)  // file -> content

	// Scan all markdown files
	err := filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(*inputDir, path)
		fileContents[relPath] = string(content)

		// Find all IDs in this file
		for _, pattern := range idPatterns {
			matches := pattern.FindAllString(string(content), -1)
			for _, id := range matches {
				idLocations[id] = append(idLocations[id], relPath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to scan directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Find missing links
	fixes := findMissingLinks(idLocations, fileContents)

	if len(fixes) == 0 {
		fmt.Println("No missing links found.")
		return domain.ExitCodeSuccess
	}

	// Report or apply fixes
	fmt.Printf("Found %d missing links:\n", len(fixes))
	for _, fix := range fixes {
		fmt.Printf("  %s: Add reference to %s\n", fix.File, fix.MissingID)
	}

	if *dryRun {
		fmt.Println("\nDry run - no changes made.")
		return domain.ExitCodeSuccess
	}

	// Apply fixes
	fixCount := 0
	for _, fix := range fixes {
		content := fileContents[fix.File]
		updated := addReference(content, fix.MissingID, fix.Section)

		if updated != content {
			fullPath := filepath.Join(*inputDir, fix.File)
			if err := os.WriteFile(fullPath, []byte(updated), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to update %s: %v\n", fix.File, err)
				continue
			}
			fixCount++
			if *verbose {
				fmt.Printf("  Updated: %s\n", fix.File)
			}
		}
	}

	fmt.Printf("\nFixed %d missing links.\n", fixCount)
	return domain.ExitCodeSuccess
}

// linkFix represents a missing link that needs to be added.
type linkFix struct {
	File      string
	MissingID string
	Section   string
}

// findMissingLinks identifies where bidirectional links are missing.
func findMissingLinks(idLocations map[string][]string, fileContents map[string]string) []linkFix {
	var fixes []linkFix

	// For each ID, check if all files that reference it also have a back-reference
	for id, files := range idLocations {
		if len(files) < 2 {
			continue // Only one file has this ID
		}

		// Find the "owner" file (usually the one that defines the ID)
		ownerFile := findOwnerFile(id, files, fileContents)

		// Check other files for back-references
		for _, file := range files {
			if file == ownerFile {
				continue
			}

			// Check if owner file references this file's IDs
			if !hasBackReference(ownerFile, file, fileContents) {
				// Find an ID from the other file to reference
				if refID := findIDFromFile(file, fileContents); refID != "" {
					fixes = append(fixes, linkFix{
						File:      ownerFile,
						MissingID: refID,
						Section:   "Related Documents",
					})
				}
			}
		}
	}

	return dedupeFixes(fixes)
}

// findOwnerFile determines which file "owns" an ID (where it's defined).
func findOwnerFile(id string, files []string, fileContents map[string]string) string {
	// Heuristic: The owner is the file where the ID appears in a heading
	for _, file := range files {
		content := fileContents[file]
		// Check if ID appears in a heading (### ID: ...)
		if strings.Contains(content, fmt.Sprintf("### %s", id)) ||
			strings.Contains(content, fmt.Sprintf("## %s", id)) {
			return file
		}
	}
	// Default to first file
	return files[0]
}

// hasBackReference checks if file1 references any ID from file2.
func hasBackReference(file1, file2 string, fileContents map[string]string) bool {
	content1 := fileContents[file1]
	content2 := fileContents[file2]

	// Find IDs in file2
	for _, pattern := range idPatterns {
		matches := pattern.FindAllString(content2, -1)
		for _, id := range matches {
			if strings.Contains(content1, id) {
				return true
			}
		}
	}
	return false
}

// findIDFromFile returns the first significant ID from a file.
func findIDFromFile(file string, fileContents map[string]string) string {
	content := fileContents[file]

	// Priority order: AC, BR, TS, IC, AGG, SEQ, TC
	priorityOrder := []string{"AC", "BR", "TS", "IC", "AGG", "SEQ", "TC"}

	for _, prefix := range priorityOrder {
		pattern := idPatterns[prefix]
		if pattern == nil {
			continue
		}
		matches := pattern.FindAllString(content, -1)
		if len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// dedupeFixes removes duplicate fixes.
func dedupeFixes(fixes []linkFix) []linkFix {
	seen := make(map[string]bool)
	var result []linkFix

	for _, fix := range fixes {
		key := fix.File + ":" + fix.MissingID
		if !seen[key] {
			seen[key] = true
			result = append(result, fix)
		}
	}

	return result
}

// addReference adds a reference to an ID in a document.
func addReference(content, id, section string) string {
	// Look for "Related Documents" or similar section
	sectionPatterns := []string{
		"## Related Documents",
		"## Related",
		"## References",
		"## See Also",
	}

	for _, pattern := range sectionPatterns {
		if idx := strings.Index(content, pattern); idx != -1 {
			// Find the end of the section (next ## or end of file)
			sectionEnd := len(content)
			nextSection := regexp.MustCompile(`\n## `)
			if loc := nextSection.FindStringIndex(content[idx+len(pattern):]); loc != nil {
				sectionEnd = idx + len(pattern) + loc[0]
			}

			// Check if ID is already referenced
			if strings.Contains(content[idx:sectionEnd], id) {
				return content
			}

			// Add reference at end of section
			insertion := fmt.Sprintf("- See: %s\n", id)
			return content[:sectionEnd] + insertion + content[sectionEnd:]
		}
	}

	// No related section found, add one at the end
	return content + fmt.Sprintf("\n\n## Related Documents\n\n- See: %s\n", id)
}
