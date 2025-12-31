package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidationResult holds the complete validation output
type ValidationResult struct {
	Level    string             `json:"level"`
	Errors   []ValidationError  `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	Checks   []ValidationCheck  `json:"checks"`
	Summary  ValidationSummary  `json:"summary"`
}

type ValidationError struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
	RefID   string `json:"ref_id,omitempty"`
}

type ValidationWarning struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type ValidationCheck struct {
	Rule    string `json:"rule"`
	Status  string `json:"status"` // "pass", "fail", "skip"
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

type ValidationSummary struct {
	TotalChecks   int `json:"total_checks"`
	Passed        int `json:"passed"`
	Failed        int `json:"failed"`
	Warnings      int `json:"warnings"`
	ErrorCount    int `json:"error_count"`
}

// Validation rule IDs
const (
	RuleV001 = "V001" // Every doc has IDs
	RuleV002 = "V002" // IDs follow pattern
	RuleV003 = "V003" // References point to existing IDs
	RuleV004 = "V004" // Bidirectional links consistent
	RuleV005 = "V005" // Every AC has at least 1 test case
	RuleV006 = "V006" // Every Entity has aggregate
	RuleV007 = "V007" // Every Service has interface contract
	RuleV008 = "V008" // Negative test ratio >= 20%
	RuleV009 = "V009" // Every AC has hallucination prevention test
	RuleV010 = "V010" // No duplicate IDs
)

// ID patterns for validation
var idPatterns = map[string]*regexp.Regexp{
	"AC":  regexp.MustCompile(`AC-[A-Z]+-\d{3}`),
	"BR":  regexp.MustCompile(`BR-[A-Z]+-\d{3}`),
	"TC":  regexp.MustCompile(`TC-AC-[A-Z]+-\d{3}-[PNBH]\d{2}`),
	"TS":  regexp.MustCompile(`TS-[A-Z]+-\d{3}`),
	"IC":  regexp.MustCompile(`IC-[A-Z]+-\d{3}`),
	"AGG": regexp.MustCompile(`AGG-[A-Z]+-\d{3}`),
	"SEQ": regexp.MustCompile(`SEQ-[A-Z]+-\d{3}`),
	"ENT": regexp.MustCompile(`ENT-[A-Z]+`),
	"BC":  regexp.MustCompile(`BC-[A-Z]+`),
}

// Generic ID pattern to find any ID-like string
var genericIDPattern = regexp.MustCompile(`\b(AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC)-[A-Z]+-?\d*`)

func runValidate() error {
	args := os.Args[2:]

	var inputDir string
	var level string
	var jsonOutput bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input-dir":
			if i+1 < len(args) {
				i++
				inputDir = args[i]
			}
		case "--level":
			if i+1 < len(args) {
				i++
				level = strings.ToUpper(args[i])
			}
		case "--json":
			jsonOutput = true
		}
	}

	if inputDir == "" {
		return fmt.Errorf("--input-dir is required")
	}
	if level == "" {
		level = "ALL"
	}

	// Run validation
	result, err := validate(inputDir, level)
	if err != nil {
		return err
	}

	// Output results
	if jsonOutput {
		outputValidationJSON(result)
		if result.Summary.Failed > 0 || result.Summary.ErrorCount > 0 {
			return fmt.Errorf("validation failed with %d errors", result.Summary.ErrorCount)
		}
		return nil
	}
	return outputText(result)
}

func validate(inputDir string, level string) (*ValidationResult, error) {
	result := &ValidationResult{
		Level:    level,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
		Checks:   []ValidationCheck{},
	}

	// Collect all IDs from documents
	allIDs := make(map[string]string)      // ID -> file
	allRefs := make(map[string][]string)   // ID -> referenced IDs
	acIDs := []string{}
	tcByAC := make(map[string][]string)    // AC ID -> TC IDs
	tcCategories := make(map[string]string) // TC ID -> category

	// Determine which files to validate based on level
	var files []string
	var err error

	switch level {
	case "L1":
		files, err = findL1Files(inputDir)
	case "L2":
		files, err = findL2Files(inputDir)
	case "L3":
		files, err = findL3Files(inputDir)
	default: // ALL
		files, err = findAllFiles(inputDir)
	}

	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Validating %s documents in %s...\n\n", level, inputDir)

	// Phase 1: Collect all IDs
	fmt.Fprintln(os.Stderr, "Phase 1: Collecting IDs...")
	for _, file := range files {
		ids, refs, err := extractIDsAndRefs(file)
		if err != nil {
			result.Warnings = append(result.Warnings, ValidationWarning{
				File:    file,
				Rule:    "PARSE",
				Message: fmt.Sprintf("Could not parse file: %v", err),
			})
			continue
		}

		for id, line := range ids {
			if existing, ok := allIDs[id]; ok {
				// V010: Duplicate ID
				result.Errors = append(result.Errors, ValidationError{
					File:    file,
					Line:    line,
					Rule:    RuleV010,
					Message: fmt.Sprintf("Duplicate ID '%s' (also in %s)", id, existing),
					RefID:   id,
				})
			} else {
				allIDs[id] = file
			}

			// Collect AC IDs
			if strings.HasPrefix(id, "AC-") {
				acIDs = append(acIDs, id)
			}

			// Collect TC info
			if strings.HasPrefix(id, "TC-") {
				// Extract AC ref from TC ID (TC-AC-CUST-001-P01 -> AC-CUST-001)
				parts := strings.Split(id, "-")
				if len(parts) >= 5 {
					acRef := fmt.Sprintf("%s-%s-%s", parts[1], parts[2], parts[3])
					tcByAC[acRef] = append(tcByAC[acRef], id)
				}
			}
		}

		for id, refList := range refs {
			allRefs[id] = append(allRefs[id], refList...)
		}

		fmt.Fprintf(os.Stderr, "  %s: %d IDs found\n", filepath.Base(file), len(ids))
	}

	// Phase 2: Structural Validation
	fmt.Fprintln(os.Stderr, "\nPhase 2: Structural Validation...")
	structuralCheck := validateStructural(files, allIDs, result)
	result.Checks = append(result.Checks, structuralCheck...)

	// Phase 3: Traceability Validation
	fmt.Fprintln(os.Stderr, "\nPhase 3: Traceability Validation...")
	traceCheck := validateTraceability(allIDs, allRefs, result)
	result.Checks = append(result.Checks, traceCheck...)

	// Phase 4: Completeness Validation
	fmt.Fprintln(os.Stderr, "\nPhase 4: Completeness Validation...")
	completeCheck := validateCompleteness(inputDir, acIDs, tcByAC, allIDs, result)
	result.Checks = append(result.Checks, completeCheck...)

	// Phase 5: TDAI Validation
	fmt.Fprintln(os.Stderr, "\nPhase 5: TDAI Validation...")
	tdaiCheck := validateTDAI(inputDir, acIDs, tcByAC, tcCategories, result)
	result.Checks = append(result.Checks, tdaiCheck...)

	// Calculate summary
	result.Summary = calculateSummary(result)

	return result, nil
}

func findL1Files(dir string) ([]string, error) {
	patterns := []string{
		"domain-model.md",
		"bounded-context-map.md",
		"acceptance-criteria.md",
		"business-rules.md",
		"decisions.md",
	}
	return findFilesByPatterns(dir, patterns)
}

func findL2Files(dir string) ([]string, error) {
	patterns := []string{
		"test-cases.md",
		"tech-specs.md",
		"interface-contracts.md",
		"aggregate-design.md",
		"sequence-design.md",
		"initial-data-model.md",
	}
	return findFilesByPatterns(dir, patterns)
}

func findL3Files(dir string) ([]string, error) {
	patterns := []string{
		"openapi.json",
		"implementation-skeletons.md",
		"feature-tickets.md",
		"service-boundaries.md",
		"event-message-design.md",
		"dependency-graph.md",
	}
	return findFilesByPatterns(dir, patterns)
}

func findAllFiles(dir string) ([]string, error) {
	var files []string

	l1, _ := findL1Files(dir)
	l2, _ := findL2Files(dir)
	l3, _ := findL3Files(dir)

	files = append(files, l1...)
	files = append(files, l2...)
	files = append(files, l3...)

	return files, nil
}

func findFilesByPatterns(dir string, patterns []string) ([]string, error) {
	var files []string
	for _, pattern := range patterns {
		path := filepath.Join(dir, pattern)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}
	return files, nil
}

func extractIDsAndRefs(file string) (map[string]int, map[string][]string, error) {
	ids := make(map[string]int)
	refs := make(map[string][]string)

	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	currentID := ""

	// Pattern to find IDs in headers: ## AC-CUST-001 – Title
	headerPattern := regexp.MustCompile(`^#{1,4}\s+((?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC)-[A-Z]+-?\d*(?:-[A-Z]\d{2})?)`)
	// Pattern to find refs in traceability sections
	refPattern := regexp.MustCompile(`(?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC)-[A-Z]+-?\d*(?:-[A-Z]\d{2})?`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for ID in header
		if matches := headerPattern.FindStringSubmatch(line); len(matches) > 1 {
			id := matches[1]
			ids[id] = lineNum
			currentID = id
		}

		// Check for references in Traceability sections
		if strings.Contains(line, "Traceability") || strings.Contains(line, "AC:") ||
		   strings.Contains(line, "BR:") || strings.Contains(line, "Source:") {
			foundRefs := refPattern.FindAllString(line, -1)
			for _, ref := range foundRefs {
				if ref != currentID && currentID != "" {
					refs[currentID] = append(refs[currentID], ref)
				}
			}
		}
	}

	return ids, refs, scanner.Err()
}

func validateStructural(files []string, allIDs map[string]string, result *ValidationResult) []ValidationCheck {
	var checks []ValidationCheck

	// V001: Check that documents have IDs
	docWithIDs := 0
	for _, file := range files {
		hasID := false
		for _, f := range allIDs {
			if f == file {
				hasID = true
				break
			}
		}
		if hasID {
			docWithIDs++
		}
	}

	if docWithIDs == len(files) {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV001,
			Status:  "pass",
			Message: fmt.Sprintf("All %d documents have IDs", len(files)),
			Count:   len(files),
		})
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV001,
			Status:  "fail",
			Message: fmt.Sprintf("Only %d of %d documents have IDs", docWithIDs, len(files)),
			Count:   docWithIDs,
		})
	}

	// V002: Check ID patterns
	validPatterns := 0
	invalidPatterns := 0
	for id := range allIDs {
		valid := false
		for _, pattern := range idPatterns {
			if pattern.MatchString(id) {
				valid = true
				break
			}
		}
		if valid {
			validPatterns++
		} else {
			invalidPatterns++
			result.Warnings = append(result.Warnings, ValidationWarning{
				Rule:    RuleV002,
				Message: fmt.Sprintf("ID '%s' does not match expected pattern", id),
			})
		}
	}

	if invalidPatterns == 0 {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV002,
			Status:  "pass",
			Message: fmt.Sprintf("All %d IDs follow expected patterns", validPatterns),
			Count:   validPatterns,
		})
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV002,
			Status:  "fail",
			Message: fmt.Sprintf("%d IDs have invalid patterns", invalidPatterns),
			Count:   invalidPatterns,
		})
	}

	// V010: Duplicate check already done during collection
	duplicates := 0
	for _, e := range result.Errors {
		if e.Rule == RuleV010 {
			duplicates++
		}
	}
	if duplicates == 0 {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV010,
			Status:  "pass",
			Message: "No duplicate IDs found",
		})
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV010,
			Status:  "fail",
			Message: fmt.Sprintf("%d duplicate IDs found", duplicates),
			Count:   duplicates,
		})
	}

	return checks
}

func validateTraceability(allIDs map[string]string, allRefs map[string][]string, result *ValidationResult) []ValidationCheck {
	var checks []ValidationCheck

	// V003: Check that all references point to existing IDs
	validRefs := 0
	invalidRefs := 0
	for fromID, refs := range allRefs {
		for _, ref := range refs {
			if _, exists := allIDs[ref]; exists {
				validRefs++
			} else {
				invalidRefs++
				result.Errors = append(result.Errors, ValidationError{
					Rule:    RuleV003,
					Message: fmt.Sprintf("Reference '%s' from '%s' not found", ref, fromID),
					RefID:   ref,
				})
			}
		}
	}

	if invalidRefs == 0 {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV003,
			Status:  "pass",
			Message: fmt.Sprintf("All %d references are valid", validRefs),
			Count:   validRefs,
		})
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV003,
			Status:  "fail",
			Message: fmt.Sprintf("%d invalid references found", invalidRefs),
			Count:   invalidRefs,
		})
	}

	// V004: Bidirectional links (simplified - just check if refs exist both ways)
	// This is a simplified check - full bidirectional would require more complex logic
	checks = append(checks, ValidationCheck{
		Rule:    RuleV004,
		Status:  "skip",
		Message: "Bidirectional link check not yet implemented",
	})

	return checks
}

func validateCompleteness(inputDir string, acIDs []string, tcByAC map[string][]string, allIDs map[string]string, result *ValidationResult) []ValidationCheck {
	var checks []ValidationCheck

	// V005: Every AC has at least 1 test case
	acsWithTests := 0
	acsWithoutTests := 0
	for _, acID := range acIDs {
		if tcs, ok := tcByAC[acID]; ok && len(tcs) > 0 {
			acsWithTests++
		} else {
			acsWithoutTests++
			result.Errors = append(result.Errors, ValidationError{
				Rule:    RuleV005,
				Message: fmt.Sprintf("AC '%s' has no test cases", acID),
				RefID:   acID,
			})
		}
	}

	if acsWithoutTests == 0 && len(acIDs) > 0 {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV005,
			Status:  "pass",
			Message: fmt.Sprintf("All %d ACs have test cases", acsWithTests),
			Count:   acsWithTests,
		})
	} else if len(acIDs) == 0 {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV005,
			Status:  "skip",
			Message: "No ACs found to validate",
		})
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV005,
			Status:  "fail",
			Message: fmt.Sprintf("%d ACs have no test cases", acsWithoutTests),
			Count:   acsWithoutTests,
		})
	}

	// V006: Every Entity has aggregate (simplified)
	checks = append(checks, ValidationCheck{
		Rule:    RuleV006,
		Status:  "skip",
		Message: "Entity-aggregate validation not yet implemented",
	})

	// V007: Every Service has interface contract (simplified)
	checks = append(checks, ValidationCheck{
		Rule:    RuleV007,
		Status:  "skip",
		Message: "Service-interface contract validation not yet implemented",
	})

	return checks
}

func validateTDAI(inputDir string, acIDs []string, tcByAC map[string][]string, tcCategories map[string]string, result *ValidationResult) []ValidationCheck {
	var checks []ValidationCheck

	// Read test-cases.md to get category information
	testCasesPath := filepath.Join(inputDir, "test-cases.md")
	categories := make(map[string]int)
	hallucinationByAC := make(map[string]bool)

	if content, err := os.ReadFile(testCasesPath); err == nil {
		lines := strings.Split(string(content), "\n")
		currentCategory := ""

		for _, line := range lines {
			// Detect category sections
			if strings.Contains(line, "## Positive Tests") {
				currentCategory = "positive"
			} else if strings.Contains(line, "## Negative Tests") {
				currentCategory = "negative"
			} else if strings.Contains(line, "## Boundary Tests") {
				currentCategory = "boundary"
			} else if strings.Contains(line, "## Hallucination Prevention") {
				currentCategory = "hallucination"
			}

			// Count tests by category
			if strings.HasPrefix(line, "### TC-") {
				categories[currentCategory]++

				// Extract AC ref from TC ID
				parts := strings.Split(strings.TrimPrefix(line, "### "), " ")
				if len(parts) > 0 {
					tcID := parts[0]
					idParts := strings.Split(tcID, "-")
					if len(idParts) >= 5 {
						acRef := fmt.Sprintf("%s-%s-%s", idParts[1], idParts[2], idParts[3])
						if currentCategory == "hallucination" {
							hallucinationByAC[acRef] = true
						}
					}
				}
			}

			// Also check AC: line for reference
			if strings.Contains(line, "- AC:") {
				acMatch := regexp.MustCompile(`AC-[A-Z]+-\d{3}`).FindString(line)
				if acMatch != "" && currentCategory == "hallucination" {
					hallucinationByAC[acMatch] = true
				}
			}
		}
	}

	total := categories["positive"] + categories["negative"] + categories["boundary"] + categories["hallucination"]

	// V008: Negative test ratio >= 20%
	if total > 0 {
		negativeRatio := float64(categories["negative"]) / float64(total)
		if negativeRatio >= 0.20 {
			checks = append(checks, ValidationCheck{
				Rule:    RuleV008,
				Status:  "pass",
				Message: fmt.Sprintf("Negative test ratio: %.1f%% (>= 20%%)", negativeRatio*100),
				Count:   categories["negative"],
			})
		} else {
			checks = append(checks, ValidationCheck{
				Rule:    RuleV008,
				Status:  "fail",
				Message: fmt.Sprintf("Negative test ratio: %.1f%% (< 20%%)", negativeRatio*100),
				Count:   categories["negative"],
			})
			result.Errors = append(result.Errors, ValidationError{
				Rule:    RuleV008,
				Message: fmt.Sprintf("Negative test ratio %.1f%% is below required 20%%", negativeRatio*100),
			})
		}
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV008,
			Status:  "skip",
			Message: "No test cases found for TDAI validation",
		})
	}

	// V009: Every AC has hallucination prevention test
	acsWithHallucination := 0
	acsWithoutHallucination := 0

	for _, acID := range acIDs {
		if hallucinationByAC[acID] {
			acsWithHallucination++
		} else {
			acsWithoutHallucination++
			result.Warnings = append(result.Warnings, ValidationWarning{
				Rule:    RuleV009,
				Message: fmt.Sprintf("AC '%s' has no hallucination prevention test", acID),
			})
		}
	}

	if len(acIDs) > 0 {
		if acsWithoutHallucination == 0 {
			checks = append(checks, ValidationCheck{
				Rule:    RuleV009,
				Status:  "pass",
				Message: fmt.Sprintf("All %d ACs have hallucination prevention tests", acsWithHallucination),
				Count:   acsWithHallucination,
			})
		} else {
			checks = append(checks, ValidationCheck{
				Rule:    RuleV009,
				Status:  "fail",
				Message: fmt.Sprintf("%d ACs missing hallucination prevention tests", acsWithoutHallucination),
				Count:   acsWithoutHallucination,
			})
		}
	} else {
		checks = append(checks, ValidationCheck{
			Rule:    RuleV009,
			Status:  "skip",
			Message: "No ACs found for hallucination test validation",
		})
	}

	return checks
}

func calculateSummary(result *ValidationResult) ValidationSummary {
	summary := ValidationSummary{
		TotalChecks: len(result.Checks),
		ErrorCount:  len(result.Errors),
		Warnings:    len(result.Warnings),
	}

	for _, check := range result.Checks {
		switch check.Status {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
		}
	}

	return summary
}

func outputText(result *ValidationResult) error {
	// Print checks by category
	fmt.Println("\n========================================")
	fmt.Println("   VALIDATION RESULTS")
	fmt.Printf("   Level: %s\n", result.Level)
	fmt.Println("========================================\n")

	// Structural
	fmt.Println("Structural Validation:")
	for _, check := range result.Checks {
		if check.Rule == RuleV001 || check.Rule == RuleV002 || check.Rule == RuleV010 {
			printCheck(check)
		}
	}

	// Traceability
	fmt.Println("\nTraceability Validation:")
	for _, check := range result.Checks {
		if check.Rule == RuleV003 || check.Rule == RuleV004 {
			printCheck(check)
		}
	}

	// Completeness
	fmt.Println("\nCompleteness Validation:")
	for _, check := range result.Checks {
		if check.Rule == RuleV005 || check.Rule == RuleV006 || check.Rule == RuleV007 {
			printCheck(check)
		}
	}

	// TDAI
	fmt.Println("\nTDAI Validation:")
	for _, check := range result.Checks {
		if check.Rule == RuleV008 || check.Rule == RuleV009 {
			printCheck(check)
		}
	}

	// Errors
	if len(result.Errors) > 0 {
		fmt.Println("\n----------------------------------------")
		fmt.Printf("ERRORS (%d):\n", len(result.Errors))
		for _, err := range result.Errors {
			if err.Line > 0 {
				fmt.Printf("  ✗ [%s] %s:%d - %s\n", err.Rule, filepath.Base(err.File), err.Line, err.Message)
			} else {
				fmt.Printf("  ✗ [%s] %s\n", err.Rule, err.Message)
			}
		}
	}

	// Warnings
	if len(result.Warnings) > 0 {
		fmt.Println("\n----------------------------------------")
		fmt.Printf("WARNINGS (%d):\n", len(result.Warnings))
		for _, warn := range result.Warnings {
			fmt.Printf("  ⚠ [%s] %s\n", warn.Rule, warn.Message)
		}
	}

	// Summary
	fmt.Println("\n========================================")
	fmt.Printf("Summary: %d passed, %d failed, %d warnings\n",
		result.Summary.Passed, result.Summary.Failed, result.Summary.Warnings)
	fmt.Println("========================================")

	// Exit with error if failures
	if result.Summary.Failed > 0 || result.Summary.ErrorCount > 0 {
		return fmt.Errorf("validation failed with %d errors", result.Summary.ErrorCount)
	}

	return nil
}

func printCheck(check ValidationCheck) {
	icon := "✓"
	if check.Status == "fail" {
		icon = "✗"
	} else if check.Status == "skip" {
		icon = "○"
	}
	fmt.Printf("  %s [%s] %s\n", icon, check.Rule, check.Message)
}

func outputValidationJSON(result *ValidationResult) {
	// Would use encoding/json here but keeping it simple
	fmt.Println("{")
	fmt.Printf("  \"level\": \"%s\",\n", result.Level)
	fmt.Printf("  \"summary\": {\n")
	fmt.Printf("    \"total_checks\": %d,\n", result.Summary.TotalChecks)
	fmt.Printf("    \"passed\": %d,\n", result.Summary.Passed)
	fmt.Printf("    \"failed\": %d,\n", result.Summary.Failed)
	fmt.Printf("    \"warnings\": %d,\n", result.Summary.Warnings)
	fmt.Printf("    \"error_count\": %d\n", result.Summary.ErrorCount)
	fmt.Printf("  }\n")
	fmt.Println("}")
}
