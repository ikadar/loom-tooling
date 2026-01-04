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
//
// Implements: IC-SYN-001
// Finds and fixes missing bidirectional references.
func runSyncLinks(args []string) int {
	fs := flag.NewFlagSet("sync-links", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "Directory to sync (required)")
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

	// Collect all markdown files
	files, err := collectMarkdownFiles(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read directory: %v\n", err)
		return domain.ExitCodeError
	}

	// Build reference map: sourceID -> []targetID
	references := make(map[string][]string)
	idToFile := make(map[string]string)
	fileContents := make(map[string]string)

	// ID pattern to match all IDs
	allIDPattern := regexp.MustCompile(`(AC|BR|ENT|BC|TC|TS|IC|AGG|SEQ|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z0-9-]+`)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		contentStr := string(content)
		fileContents[file] = contentStr

		// Extract IDs and their file locations
		ids := allIDPattern.FindAllString(contentStr, -1)
		for _, id := range ids {
			if _, exists := idToFile[id]; !exists {
				idToFile[id] = file
			}
		}
	}

	// Find references between IDs
	for file, content := range fileContents {
		// Find IDs defined in this file (in headers or as definitions)
		definedIDs := findDefinedIDs(content)

		// Find IDs referenced in this file
		referencedIDs := allIDPattern.FindAllString(content, -1)

		for _, defID := range definedIDs {
			for _, refID := range referencedIDs {
				if defID != refID && !contains(references[defID], refID) {
					references[defID] = append(references[defID], refID)
				}
			}
		}
		_ = file // file is used implicitly through content
	}

	// Find missing back-references
	type missingRef struct {
		sourceID   string
		targetID   string
		sourceFile string
		targetFile string
	}
	var missing []missingRef

	for sourceID, targets := range references {
		for _, targetID := range targets {
			// Check if target references back to source
			if targetRefs, ok := references[targetID]; ok {
				if !contains(targetRefs, sourceID) {
					missing = append(missing, missingRef{
						sourceID:   sourceID,
						targetID:   targetID,
						sourceFile: idToFile[sourceID],
						targetFile: idToFile[targetID],
					})
				}
			}
		}
	}

	if len(missing) == 0 {
		fmt.Println("All bidirectional links are in sync.")
		return domain.ExitCodeSuccess
	}

	if *dryRun {
		fmt.Println("Missing back-references (dry run):")
		for _, m := range missing {
			fmt.Printf("  %s -> %s (in %s)\n", m.targetID, m.sourceID, m.targetFile)
		}
		fmt.Printf("\n%d missing back-references found.\n", len(missing))
		return domain.ExitCodeSuccess
	}

	// Apply fixes
	fileFixes := make(map[string][]string) // file -> list of IDs to add

	for _, m := range missing {
		if m.targetFile != "" {
			fileFixes[m.targetFile] = append(fileFixes[m.targetFile], m.sourceID)
		}
	}

	fixCount := 0
	for file, idsToAdd := range fileFixes {
		content := fileContents[file]

		// Add references to traceability section or end of file
		additions := "\n\n## Added References\n\n"
		for _, id := range uniqueStrings(idsToAdd) {
			additions += fmt.Sprintf("- %s\n", id)
		}

		// Check if file has a Traceability section
		if idx := strings.Index(content, "## Traceability"); idx != -1 {
			// Insert before next section or at end
			nextSection := strings.Index(content[idx+1:], "\n## ")
			if nextSection != -1 {
				insertPoint := idx + 1 + nextSection
				content = content[:insertPoint] + "\n**Added by sync-links:**\n" +
					strings.Join(uniqueStrings(idsToAdd), ", ") + "\n" + content[insertPoint:]
			} else {
				content += "\n**Added by sync-links:** " + strings.Join(uniqueStrings(idsToAdd), ", ") + "\n"
			}
		} else {
			content += additions
		}

		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to update %s: %v\n", file, err)
			continue
		}
		fixCount += len(idsToAdd)

		if *verbose {
			fmt.Printf("Updated: %s (added %d references)\n", file, len(idsToAdd))
		}
	}

	fmt.Printf("Sync complete. Added %d back-references to %d files.\n", fixCount, len(fileFixes))
	return domain.ExitCodeSuccess
}

// findDefinedIDs finds IDs that are defined (not just referenced) in content.
func findDefinedIDs(content string) []string {
	var ids []string

	// Look for IDs in markdown headers: ### ID: Title or ### ID
	headerPattern := regexp.MustCompile(`###\s+((?:AC|BR|ENT|BC|TC|TS|IC|AGG|SEQ|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z0-9-]+)`)
	matches := headerPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			ids = append(ids, match[1])
		}
	}

	// Look for IDs at start of lines: AC-ORD-001: Title
	linePattern := regexp.MustCompile(`(?m)^((?:AC|BR|ENT|BC|TC|TS|IC|AGG|SEQ|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z0-9-]+):`)
	matches = linePattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			ids = append(ids, match[1])
		}
	}

	return uniqueStrings(ids)
}

// contains checks if a string slice contains a value.
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// collectMarkdownFiles is defined in validate.go
// uniqueStrings is defined in validate.go
