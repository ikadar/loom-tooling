package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Reference represents a reference from one ID to another
type Reference struct {
	FromID   string
	ToID     string
	FromFile string
}

// DocumentInfo holds parsed information about a document
type DocumentInfo struct {
	FilePath     string
	IDs          map[string]int      // ID -> line number
	References   map[string][]string // ID -> referenced IDs
	BackRefs     map[string][]string // ID -> IDs that reference this ID
	TraceSection map[string]int      // ID -> line number of Traceability section
}

func runSyncLinks() error {
	// Parse arguments
	inputDir := ""
	dryRun := false

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--input-dir":
			if i+1 < len(os.Args) {
				inputDir = os.Args[i+1]
				i++
			}
		case "--dry-run":
			dryRun = true
		}
	}

	if inputDir == "" {
		return fmt.Errorf("--input-dir is required")
	}

	fmt.Fprintf(os.Stderr, "Syncing bidirectional links in %s...\n\n", inputDir)

	// Find all markdown files
	files, err := findAllMarkdownFiles(inputDir)
	if err != nil {
		return err
	}

	// Phase 1: Parse all documents
	fmt.Fprintln(os.Stderr, "Phase 1: Parsing documents...")
	docs := make(map[string]*DocumentInfo)
	allIDs := make(map[string]string) // ID -> file

	for _, file := range files {
		doc, err := parseDocumentForSync(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Warning: Could not parse %s: %v\n", filepath.Base(file), err)
			continue
		}
		docs[file] = doc
		for id := range doc.IDs {
			allIDs[id] = file
		}
		fmt.Fprintf(os.Stderr, "  %s: %d IDs, %d references\n", filepath.Base(file), len(doc.IDs), countRefs(doc.References))
	}

	// Phase 2: Build reference graph
	fmt.Fprintln(os.Stderr, "\nPhase 2: Building reference graph...")
	allRefs := []Reference{}

	for file, doc := range docs {
		for fromID, toIDs := range doc.References {
			for _, toID := range toIDs {
				// Only add if target ID exists
				if _, exists := allIDs[toID]; exists {
					allRefs = append(allRefs, Reference{
						FromID:   fromID,
						ToID:     toID,
						FromFile: file,
					})
				}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "  Found %d valid references\n", len(allRefs))

	// Phase 3: Find missing back-references
	fmt.Fprintln(os.Stderr, "\nPhase 3: Finding missing back-references...")
	missingBackRefs := findMissingBackRefs(allRefs, docs, allIDs)
	fmt.Fprintf(os.Stderr, "  Found %d missing back-references\n", len(missingBackRefs))

	if len(missingBackRefs) == 0 {
		fmt.Fprintln(os.Stderr, "\n✓ All links are bidirectional!")
		return nil
	}

	// Group by target file for reporting
	byFile := make(map[string][]Reference)
	for _, ref := range missingBackRefs {
		targetFile := allIDs[ref.ToID]
		byFile[targetFile] = append(byFile[targetFile], ref)
	}

	fmt.Fprintln(os.Stderr, "\nMissing back-references:")
	for file, refs := range byFile {
		fmt.Fprintf(os.Stderr, "  %s:\n", filepath.Base(file))
		for _, ref := range refs {
			fmt.Fprintf(os.Stderr, "    %s should reference back to %s\n", ref.ToID, ref.FromID)
		}
	}

	if dryRun {
		fmt.Fprintln(os.Stderr, "\n--dry-run: No files modified")
		return nil
	}

	// Phase 4: Add missing back-references
	fmt.Fprintln(os.Stderr, "\nPhase 4: Adding missing back-references...")
	filesModified := 0
	linksAdded := 0

	for file, refs := range byFile {
		added, err := addBackReferences(file, refs, docs[file])
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error updating %s: %v\n", filepath.Base(file), err)
			continue
		}
		if added > 0 {
			filesModified++
			linksAdded += added
			fmt.Fprintf(os.Stderr, "  Updated: %s (+%d back-references)\n", filepath.Base(file), added)
		}
	}

	fmt.Fprintf(os.Stderr, "\n========================================\n")
	fmt.Fprintf(os.Stderr, "   SYNC COMPLETE\n")
	fmt.Fprintf(os.Stderr, "========================================\n")
	fmt.Fprintf(os.Stderr, "Files modified:      %d\n", filesModified)
	fmt.Fprintf(os.Stderr, "Back-references added: %d\n", linksAdded)

	return nil
}

func findAllMarkdownFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseDocumentForSync(file string) (*DocumentInfo, error) {
	doc := &DocumentInfo{
		FilePath:     file,
		IDs:          make(map[string]int),
		References:   make(map[string][]string),
		BackRefs:     make(map[string][]string),
		TraceSection: make(map[string]int),
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	currentID := ""
	inTraceSection := false

	// Pattern to find IDs in headers: ## ID – Title {#anchor} or ### ID: Title
	// Captures everything from prefix to the separator (–, :, or {#)
	headerPattern := regexp.MustCompile(`^#{1,4}\s+((?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC|EVT|CMD|INT|SVC|FDT|SKEL|DEP|VO)-[A-Z0-9-]+(?:-[A-Z]\d{2})?)(?:\s+[–:—]|\s+\{#)`)
	// Pattern to find refs - matches ID patterns in text
	refPattern := regexp.MustCompile(`\b((?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC|EVT|CMD|INT|SVC|FDT|SKEL|DEP|VO)-[A-Z0-9]+-[A-Z0-9-]+(?:-[A-Z]\d{2})?)\b`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for ID in header
		if matches := headerPattern.FindStringSubmatch(line); len(matches) > 1 {
			id := matches[1]
			doc.IDs[id] = lineNum
			currentID = id
			inTraceSection = false
		}

		// Check for Traceability section
		if strings.Contains(line, "**Traceability:**") || strings.Contains(line, "## Traceability") {
			if currentID != "" {
				doc.TraceSection[currentID] = lineNum
			}
			inTraceSection = true
		}

		// Check for section end (next header or ---)
		if inTraceSection && (strings.HasPrefix(line, "#") || strings.HasPrefix(line, "---")) {
			inTraceSection = false
		}

		// Extract references in traceability sections
		if inTraceSection || strings.Contains(line, "- Source:") || strings.Contains(line, "- AC:") ||
			strings.Contains(line, "- BR:") || strings.Contains(line, "- Related") ||
			strings.Contains(line, "- Referenced by:") {
			matches := refPattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 1 {
					ref := match[1]
					if ref != currentID && currentID != "" {
						// Avoid duplicates
						if !contains(doc.References[currentID], ref) {
							doc.References[currentID] = append(doc.References[currentID], ref)
						}
					}
				}
			}
		}
	}

	return doc, scanner.Err()
}

func countRefs(refs map[string][]string) int {
	count := 0
	for _, r := range refs {
		count += len(r)
	}
	return count
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func findMissingBackRefs(allRefs []Reference, docs map[string]*DocumentInfo, allIDs map[string]string) []Reference {
	var missing []Reference

	for _, ref := range allRefs {
		// Check if ToID has a back-reference to FromID
		targetFile := allIDs[ref.ToID]
		if targetFile == "" {
			continue
		}

		targetDoc := docs[targetFile]
		if targetDoc == nil {
			continue
		}

		// Check if target ID has a reference back to source ID
		hasBackRef := false
		if backRefs, ok := targetDoc.References[ref.ToID]; ok {
			for _, backRef := range backRefs {
				if backRef == ref.FromID {
					hasBackRef = true
					break
				}
			}
		}

		if !hasBackRef {
			missing = append(missing, ref)
		}
	}

	return missing
}

func addBackReferences(file string, refs []Reference, doc *DocumentInfo) (int, error) {
	// Read the file
	content, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(content), "\n")

	// Group refs by target ID
	refsByTarget := make(map[string][]string)
	for _, ref := range refs {
		refsByTarget[ref.ToID] = append(refsByTarget[ref.ToID], ref.FromID)
	}

	// Sort target IDs for consistent output
	var targetIDs []string
	for id := range refsByTarget {
		targetIDs = append(targetIDs, id)
	}
	sort.Strings(targetIDs)

	added := 0

	// For each target ID, find its Traceability section and add back-refs
	for _, targetID := range targetIDs {
		fromIDs := refsByTarget[targetID]
		sort.Strings(fromIDs)

		// Find the section for this ID
		traceLine, hasTrace := doc.TraceSection[targetID]
		if !hasTrace {
			// Need to find where to add Traceability section
			// Find the end of this ID's section (next --- or next header)
			idLine := doc.IDs[targetID]
			insertLine := findSectionEnd(lines, idLine)

			// Insert Traceability section
			newLines := []string{
				"",
				"**Traceability:**",
			}
			for _, fromID := range fromIDs {
				newLines = append(newLines, fmt.Sprintf("- Referenced by: %s", fromID))
				added++
			}

			lines = insertLines(lines, insertLine, newLines)
		} else {
			// Find where to insert in existing Traceability section
			// Look for the end of the traceability section
			insertLine := traceLine + 1
			for insertLine < len(lines) {
				line := lines[insertLine]
				if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "---") || line == "" {
					break
				}
				insertLine++
			}

			// Add back-references
			var newLines []string
			for _, fromID := range fromIDs {
				backRefLine := fmt.Sprintf("- Referenced by: %s", fromID)
				// Check if this back-ref already exists
				alreadyExists := false
				for i := traceLine; i < insertLine && i < len(lines); i++ {
					if strings.Contains(lines[i], fromID) && strings.Contains(lines[i], "Referenced by") {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					newLines = append(newLines, backRefLine)
					added++
				}
			}

			if len(newLines) > 0 {
				lines = insertLines(lines, insertLine, newLines)
			}
		}
	}

	if added > 0 {
		// Write the file
		err = os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
		if err != nil {
			return 0, err
		}
	}

	return added, nil
}

func findSectionEnd(lines []string, startLine int) int {
	for i := startLine; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "---") {
			return i
		}
		if i > startLine && strings.HasPrefix(lines[i], "#") {
			return i
		}
	}
	return len(lines)
}

func insertLines(lines []string, at int, newLines []string) []string {
	if at >= len(lines) {
		return append(lines, newLines...)
	}
	result := make([]string, 0, len(lines)+len(newLines))
	result = append(result, lines[:at]...)
	result = append(result, newLines...)
	result = append(result, lines[at:]...)
	return result
}
