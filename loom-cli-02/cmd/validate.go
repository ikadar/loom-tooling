// Implements: l2/interface-contracts.md IC-VAL-001
// See: l2/sequence-design.md SEQ-VAL-001
// See: l1/business-rules.md BR-VAL-001/002/003
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"loom-cli/internal/domain"
)

// Validation rule patterns
//
// Implements: l2/tech-specs.md TS-FMT-002
var idPatterns = map[string]*regexp.Regexp{
	// L1 patterns
	"AC":  regexp.MustCompile(`AC-[A-Z]+-\d{3}`),
	"BR":  regexp.MustCompile(`BR-[A-Z]+-\d{3}`),
	"ENT": regexp.MustCompile(`ENT-[A-Z]+`),
	"BC":  regexp.MustCompile(`BC-[A-Z]+`),
	// L2 patterns
	"TC":  regexp.MustCompile(`TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2}`),
	"TS":  regexp.MustCompile(`TS-[A-Z]+-\d{3}`),
	"IC":  regexp.MustCompile(`IC-[A-Z]+-\d{3}`),
	"AGG": regexp.MustCompile(`AGG-[A-Z]+-\d{3}`),
	"SEQ": regexp.MustCompile(`SEQ-[A-Z]+-\d{3}`),
	// L3 patterns
	"EVT":  regexp.MustCompile(`EVT-[A-Z]+-\d{3}`),
	"CMD":  regexp.MustCompile(`CMD-[A-Z]+-\d{3}`),
	"INT":  regexp.MustCompile(`INT-[A-Z]+-\d{3}`),
	"SVC":  regexp.MustCompile(`SVC-[A-Z]+`),
	"FDT":  regexp.MustCompile(`FDT-\d{3}`),
	"SKEL": regexp.MustCompile(`SKEL-[A-Z]+-\d{3}`),
	"DEP":  regexp.MustCompile(`DEP-[A-Z]+-\d{3}`),
}

// runValidate implements the validate command.
//
// Implements: IC-VAL-001
// Validation Rules:
//   - V001: Documents have IDs
//   - V002: IDs follow patterns
//   - V003: References exist
//   - V004: Bidirectional links (deferred - DEC-L1-013)
//   - V005: AC has test cases
//   - V006: Entity has aggregate
//   - V007: Service has contract
//   - V008: Negative test ratio >= 20%
//   - V009: Hallucination tests exist
//   - V010: No duplicate IDs
func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	inputDir := fs.String("input-dir", "", "Directory to validate (required)")
	level := fs.String("level", "ALL", "Validation level: L1, L2, L3, or ALL")
	jsonOutput := fs.Bool("json", false, "Output as JSON")

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

	// Build ID registry
	idRegistry := make(map[string]string)  // ID -> file
	allIDs := make(map[string][]string)    // ID -> list of files (for duplicate detection)
	references := make(map[string][]string) // file -> list of referenced IDs
	fileContents := make(map[string]string) // file -> content

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		contentStr := string(content)
		fileContents[file] = contentStr

		// Extract IDs defined in this file
		ids := extractIDs(contentStr)
		for _, id := range ids {
			idRegistry[id] = file
			allIDs[id] = append(allIDs[id], file)
		}

		// Extract references to other IDs
		refs := extractReferences(contentStr)
		references[file] = refs
	}

	// Run validation checks
	result := &domain.ValidationResult{
		Level:    *level,
		Errors:   []domain.ValidationError{},
		Warnings: []domain.ValidationWarning{},
		Checks:   []domain.ValidationCheck{},
	}

	// V001: Documents have IDs
	v001 := runV001(files, fileContents)
	result.Checks = append(result.Checks, v001)
	if v001.Status == "fail" {
		result.Errors = append(result.Errors, domain.ValidationError{
			Rule:    "V001",
			Message: "Some documents lack proper IDs",
		})
	}

	// V002: IDs follow patterns
	v002 := runV002(idRegistry)
	result.Checks = append(result.Checks, v002)

	// V003: References exist
	v003 := runV003(references, idRegistry)
	result.Checks = append(result.Checks, v003)
	for _, ref := range v003.Message {
		if strings.Contains(string(ref), "not found") {
			result.Errors = append(result.Errors, domain.ValidationError{
				Rule:    "V003",
				Message: string(ref),
			})
		}
	}

	// V004: Bidirectional links (deferred - DEC-L1-013)
	v004 := domain.ValidationCheck{
		Rule:    "V004",
		Status:  "skip",
		Message: "Deferred per DEC-L1-013",
	}
	result.Checks = append(result.Checks, v004)

	// V005-V009: Level-specific checks
	if *level == "L2" || *level == "L3" || *level == "ALL" {
		v005 := runV005(idRegistry, references)
		result.Checks = append(result.Checks, v005)

		v006 := runV006(idRegistry)
		result.Checks = append(result.Checks, v006)

		v007 := runV007(idRegistry)
		result.Checks = append(result.Checks, v007)
	}

	if *level == "L3" || *level == "ALL" {
		v008 := runV008(idRegistry)
		result.Checks = append(result.Checks, v008)

		v009 := runV009(idRegistry)
		result.Checks = append(result.Checks, v009)
	}

	// V010: No duplicate IDs
	v010 := runV010(allIDs)
	result.Checks = append(result.Checks, v010)

	// Calculate summary
	result.Summary = calculateSummary(result)

	// Output
	if *jsonOutput {
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
	} else {
		printValidationResult(result)
	}

	if result.Summary.ErrorCount > 0 {
		return domain.ExitCodeError
	}
	return domain.ExitCodeSuccess
}

// collectMarkdownFiles collects all .md files from directory and subdirectories.
func collectMarkdownFiles(dir string) ([]string, error) {
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

// extractIDs extracts all defined IDs from content.
func extractIDs(content string) []string {
	var ids []string
	for _, pattern := range idPatterns {
		matches := pattern.FindAllString(content, -1)
		ids = append(ids, matches...)
	}
	return uniqueStrings(ids)
}

// extractReferences extracts all referenced IDs from content.
func extractReferences(content string) []string {
	// Same as extractIDs for now - could be refined to distinguish definitions vs references
	return extractIDs(content)
}

// uniqueStrings returns unique strings from a slice.
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// V001: Documents have IDs
func runV001(files []string, contents map[string]string) domain.ValidationCheck {
	missingCount := 0
	for _, file := range files {
		content := contents[file]
		ids := extractIDs(content)
		if len(ids) == 0 {
			missingCount++
		}
	}

	status := "pass"
	if missingCount > 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V001",
		Status:  status,
		Message: fmt.Sprintf("%d/%d documents have IDs", len(files)-missingCount, len(files)),
		Count:   len(files) - missingCount,
	}
}

// V002: IDs follow patterns
func runV002(idRegistry map[string]string) domain.ValidationCheck {
	invalidCount := 0
	for id := range idRegistry {
		valid := false
		for _, pattern := range idPatterns {
			if pattern.MatchString(id) {
				valid = true
				break
			}
		}
		if !valid {
			invalidCount++
		}
	}

	status := "pass"
	if invalidCount > 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V002",
		Status:  status,
		Message: fmt.Sprintf("%d/%d IDs follow patterns", len(idRegistry)-invalidCount, len(idRegistry)),
		Count:   len(idRegistry) - invalidCount,
	}
}

// V003: References exist
func runV003(references map[string][]string, idRegistry map[string]string) domain.ValidationCheck {
	missingCount := 0
	for _, refs := range references {
		for _, ref := range refs {
			if _, exists := idRegistry[ref]; !exists {
				missingCount++
			}
		}
	}

	status := "pass"
	if missingCount > 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V003",
		Status:  status,
		Message: fmt.Sprintf("%d unresolved references", missingCount),
		Count:   missingCount,
	}
}

// V005: AC has test cases
func runV005(idRegistry map[string]string, references map[string][]string) domain.ValidationCheck {
	acPattern := idPatterns["AC"]
	tcPattern := idPatterns["TC"]

	acCount := 0
	coveredCount := 0

	for id := range idRegistry {
		if acPattern.MatchString(id) {
			acCount++
			// Check if there's a TC referencing this AC
			for tcID := range idRegistry {
				if tcPattern.MatchString(tcID) && strings.Contains(tcID, id[3:]) {
					coveredCount++
					break
				}
			}
		}
	}

	status := "pass"
	if acCount > 0 && coveredCount < acCount {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V005",
		Status:  status,
		Message: fmt.Sprintf("%d/%d ACs have test cases", coveredCount, acCount),
		Count:   coveredCount,
	}
}

// V006: Entity has aggregate
func runV006(idRegistry map[string]string) domain.ValidationCheck {
	entPattern := idPatterns["ENT"]
	aggPattern := idPatterns["AGG"]

	entCount := 0
	for id := range idRegistry {
		if entPattern.MatchString(id) {
			entCount++
		}
	}

	aggCount := 0
	for id := range idRegistry {
		if aggPattern.MatchString(id) {
			aggCount++
		}
	}

	status := "pass"
	if entCount > 0 && aggCount == 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V006",
		Status:  status,
		Message: fmt.Sprintf("%d entities, %d aggregates", entCount, aggCount),
		Count:   aggCount,
	}
}

// V007: Service has contract
func runV007(idRegistry map[string]string) domain.ValidationCheck {
	svcPattern := idPatterns["SVC"]
	icPattern := idPatterns["IC"]

	svcCount := 0
	for id := range idRegistry {
		if svcPattern.MatchString(id) {
			svcCount++
		}
	}

	icCount := 0
	for id := range idRegistry {
		if icPattern.MatchString(id) {
			icCount++
		}
	}

	status := "pass"
	if svcCount > 0 && icCount == 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V007",
		Status:  status,
		Message: fmt.Sprintf("%d services, %d contracts", svcCount, icCount),
		Count:   icCount,
	}
}

// V008: Negative test ratio >= 20%
func runV008(idRegistry map[string]string) domain.ValidationCheck {
	tcPattern := idPatterns["TC"]

	totalTests := 0
	negativeTests := 0

	for id := range idRegistry {
		if tcPattern.MatchString(id) {
			totalTests++
			if strings.Contains(id, "-N") {
				negativeTests++
			}
		}
	}

	ratio := 0.0
	if totalTests > 0 {
		ratio = float64(negativeTests) / float64(totalTests) * 100
	}

	status := "pass"
	if totalTests > 0 && ratio < 20 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V008",
		Status:  status,
		Message: fmt.Sprintf("%.1f%% negative tests (%d/%d)", ratio, negativeTests, totalTests),
		Count:   negativeTests,
	}
}

// V009: Hallucination tests exist
func runV009(idRegistry map[string]string) domain.ValidationCheck {
	tcPattern := idPatterns["TC"]

	hallucinationTests := 0

	for id := range idRegistry {
		if tcPattern.MatchString(id) && strings.Contains(id, "-H") {
			hallucinationTests++
		}
	}

	status := "pass"
	if hallucinationTests == 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V009",
		Status:  status,
		Message: fmt.Sprintf("%d hallucination prevention tests", hallucinationTests),
		Count:   hallucinationTests,
	}
}

// V010: No duplicate IDs
func runV010(allIDs map[string][]string) domain.ValidationCheck {
	duplicateCount := 0
	for _, files := range allIDs {
		if len(files) > 1 {
			duplicateCount++
		}
	}

	status := "pass"
	if duplicateCount > 0 {
		status = "fail"
	}

	return domain.ValidationCheck{
		Rule:    "V010",
		Status:  status,
		Message: fmt.Sprintf("%d duplicate IDs", duplicateCount),
		Count:   duplicateCount,
	}
}

func calculateSummary(result *domain.ValidationResult) domain.ValidationSummary {
	passed := 0
	failed := 0
	for _, check := range result.Checks {
		switch check.Status {
		case "pass":
			passed++
		case "fail":
			failed++
		}
	}

	return domain.ValidationSummary{
		TotalChecks: len(result.Checks),
		Passed:      passed,
		Failed:      failed,
		Warnings:    len(result.Warnings),
		ErrorCount:  len(result.Errors),
	}
}

func printValidationResult(result *domain.ValidationResult) {
	fmt.Printf("Validation Results (Level: %s)\n", result.Level)
	fmt.Println(strings.Repeat("=", 50))

	for _, check := range result.Checks {
		status := "✓"
		if check.Status == "fail" {
			status = "✗"
		} else if check.Status == "skip" {
			status = "-"
		}
		fmt.Printf("[%s] %s: %s\n", status, check.Rule, check.Message)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Summary: %d passed, %d failed, %d warnings\n",
		result.Summary.Passed, result.Summary.Failed, result.Summary.Warnings)

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range result.Errors {
			fmt.Printf("  [%s] %s\n", err.Rule, err.Message)
		}
	}
}
