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

// Generic ID pattern to find any ID-like string
var genericIDPattern = regexp.MustCompile(`\b(AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z]*-?\d*`)

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

	// Pattern to find IDs in headers: ## AC-CUST-001 – Title or ### EVT-CUST-001: EventName
	headerPattern := regexp.MustCompile(`^#{1,4}\s+((?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z]*-?\d*(?:-[A-Z]\d{2})?)`)
	// Pattern to find refs in traceability sections
	refPattern := regexp.MustCompile(`(?:AC|BR|TC|TS|IC|AGG|SEQ|ENT|BC|EVT|CMD|INT|SVC|FDT|SKEL|DEP)-[A-Z]*-?\d*(?:-[A-Z]\d{2})?`)

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

	// V006: Every Entity has aggregate
	v006Check := validateEntityAggregates(inputDir, result)
	checks = append(checks, v006Check)

	// V007: Every Service has interface contract
	v007Check := validateServiceContracts(inputDir, result)
	checks = append(checks, v007Check)

	return checks
}

// EntityInfo holds parsed entity information from domain-model.md
type EntityInfo struct {
	ID       string
	Name     string
	Type     string // "aggregate_root" or "entity"
	LineNum  int
}

// AggregateInfo holds parsed aggregate information from aggregate-design.md
type AggregateInfo struct {
	ID              string
	Name            string
	RootEntityName  string
	ChildEntities   []string
	LineNum         int
}

func validateEntityAggregates(inputDir string, result *ValidationResult) ValidationCheck {
	// Try to find domain-model.md and aggregate-design.md
	domainModelPath := filepath.Join(inputDir, "domain-model.md")
	aggregateDesignPath := filepath.Join(inputDir, "aggregate-design.md")

	// Check if files exist
	if _, err := os.Stat(domainModelPath); os.IsNotExist(err) {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "skip",
			Message: "domain-model.md not found",
		}
	}
	if _, err := os.Stat(aggregateDesignPath); os.IsNotExist(err) {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "skip",
			Message: "aggregate-design.md not found",
		}
	}

	// Parse entities from domain-model.md
	entities, err := parseEntitiesFromDomainModel(domainModelPath)
	if err != nil {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "skip",
			Message: fmt.Sprintf("Could not parse domain-model.md: %v", err),
		}
	}

	if len(entities) == 0 {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "skip",
			Message: "No entities found in domain-model.md",
		}
	}

	// Parse aggregates from aggregate-design.md
	aggregates, err := parseAggregatesFromDesign(aggregateDesignPath)
	if err != nil {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "skip",
			Message: fmt.Sprintf("Could not parse aggregate-design.md: %v", err),
		}
	}

	// Build lookup maps
	aggByName := make(map[string]*AggregateInfo)
	childEntityNames := make(map[string]bool)
	for i := range aggregates {
		agg := &aggregates[i]
		aggByName[strings.ToLower(agg.Name)] = agg
		aggByName[strings.ToLower(agg.RootEntityName)] = agg
		for _, child := range agg.ChildEntities {
			childEntityNames[strings.ToLower(child)] = true
		}
	}

	// Validate: every aggregate_root entity has a corresponding aggregate
	// Validate: every entity is referenced as child in some aggregate
	entitiesWithAggregate := 0
	entitiesWithoutAggregate := 0

	for _, ent := range entities {
		hasAggregate := false
		entityNameLower := strings.ToLower(ent.Name)

		if ent.Type == "aggregate_root" {
			// Check if there's an aggregate for this root
			if _, ok := aggByName[entityNameLower]; ok {
				hasAggregate = true
			}
		} else {
			// Check if this entity is a child in some aggregate
			if childEntityNames[entityNameLower] {
				hasAggregate = true
			}
		}

		if hasAggregate {
			entitiesWithAggregate++
		} else {
			entitiesWithoutAggregate++
			result.Errors = append(result.Errors, ValidationError{
				Rule:    RuleV006,
				Message: fmt.Sprintf("Entity '%s' (%s) has no aggregate", ent.ID, ent.Type),
				RefID:   ent.ID,
			})
		}
	}

	if entitiesWithoutAggregate == 0 {
		return ValidationCheck{
			Rule:    RuleV006,
			Status:  "pass",
			Message: fmt.Sprintf("All %d entities have aggregates", entitiesWithAggregate),
			Count:   entitiesWithAggregate,
		}
	}

	return ValidationCheck{
		Rule:    RuleV006,
		Status:  "fail",
		Message: fmt.Sprintf("%d entities have no aggregates", entitiesWithoutAggregate),
		Count:   entitiesWithoutAggregate,
	}
}

func parseEntitiesFromDomainModel(filePath string) ([]EntityInfo, error) {
	var entities []EntityInfo

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Pattern: ### ENT-XXX – Name or ### ENT-XXX: Name
	entityPattern := regexp.MustCompile(`^###\s+(ENT-[A-Z]+-\d{3})\s*[–:—]\s*(\w+)`)
	typePattern := regexp.MustCompile(`^\*\*Type:\*\*\s*(\w+)`)

	var currentEntity *EntityInfo

	for i, line := range lines {
		// Check for entity header
		if matches := entityPattern.FindStringSubmatch(line); len(matches) > 2 {
			if currentEntity != nil {
				entities = append(entities, *currentEntity)
			}
			currentEntity = &EntityInfo{
				ID:      matches[1],
				Name:    matches[2],
				LineNum: i + 1,
				Type:    "entity", // Default
			}
		}

		// Check for Type: line
		if currentEntity != nil {
			if matches := typePattern.FindStringSubmatch(line); len(matches) > 1 {
				currentEntity.Type = strings.ToLower(matches[1])
			}
		}
	}

	// Don't forget the last entity
	if currentEntity != nil {
		entities = append(entities, *currentEntity)
	}

	return entities, nil
}

func parseAggregatesFromDesign(filePath string) ([]AggregateInfo, error) {
	var aggregates []AggregateInfo

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Pattern: ## AGG-XXX – Name
	aggPattern := regexp.MustCompile(`^##\s+(AGG-[A-Z]+-\d{3})\s*[–:—]\s*(\w+)`)
	// Pattern: ### Aggregate Root: Name
	rootPattern := regexp.MustCompile(`^###\s+Aggregate Root:\s*(\w+)`)
	// Pattern: #### ChildEntityName
	childPattern := regexp.MustCompile(`^####\s+(\w+)`)

	var currentAgg *AggregateInfo
	inChildSection := false

	for i, line := range lines {
		// Check for aggregate header
		if matches := aggPattern.FindStringSubmatch(line); len(matches) > 2 {
			if currentAgg != nil {
				aggregates = append(aggregates, *currentAgg)
			}
			currentAgg = &AggregateInfo{
				ID:            matches[1],
				Name:          matches[2],
				LineNum:       i + 1,
				ChildEntities: []string{},
			}
			inChildSection = false
		}

		if currentAgg != nil {
			// Check for aggregate root
			if matches := rootPattern.FindStringSubmatch(line); len(matches) > 1 {
				currentAgg.RootEntityName = matches[1]
			}

			// Check for Child Entities section
			if strings.Contains(line, "### Child Entities") {
				inChildSection = true
			}

			// Check for child entity (#### EntityName)
			if inChildSection {
				if matches := childPattern.FindStringSubmatch(line); len(matches) > 1 {
					currentAgg.ChildEntities = append(currentAgg.ChildEntities, matches[1])
				}
			}

			// End of child section on next ## or ###
			if inChildSection && (strings.HasPrefix(line, "### ") && !strings.Contains(line, "Child Entities")) {
				inChildSection = false
			}
		}
	}

	// Don't forget the last aggregate
	if currentAgg != nil {
		aggregates = append(aggregates, *currentAgg)
	}

	return aggregates, nil
}

// ServiceInfo holds parsed service information from service-boundaries.md
type ServiceInfo struct {
	ID      string
	Name    string
	APIBase string
	LineNum int
}

// ContractInfo holds parsed interface contract information
type ContractInfo struct {
	ID      string
	Name    string
	BaseURL string
	LineNum int
}

// serviceContractMapping maps service domain keywords to IC domain prefixes
var serviceContractMapping = map[string][]string{
	"CUSTOMER":  {"CUST"},
	"CATALOG":   {"PROD", "CAT"},
	"SHOPPING":  {"CART"},
	"CART":      {"CART"},
	"ORDER":     {"ORDER"},
	"INVENTORY": {"INV"},
	"PAYMENT":   {"PAY"},
	"SHIPPING":  {"SHIP"},
}

func validateServiceContracts(inputDir string, result *ValidationResult) ValidationCheck {
	// Try to find service-boundaries.md and interface-contracts.md
	serviceBoundariesPath := filepath.Join(inputDir, "service-boundaries.md")
	interfaceContractsPath := filepath.Join(inputDir, "interface-contracts.md")

	// Check if files exist
	if _, err := os.Stat(serviceBoundariesPath); os.IsNotExist(err) {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "skip",
			Message: "service-boundaries.md not found",
		}
	}
	if _, err := os.Stat(interfaceContractsPath); os.IsNotExist(err) {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "skip",
			Message: "interface-contracts.md not found",
		}
	}

	// Parse services from service-boundaries.md
	services, err := parseServicesFromBoundaries(serviceBoundariesPath)
	if err != nil {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "skip",
			Message: fmt.Sprintf("Could not parse service-boundaries.md: %v", err),
		}
	}

	if len(services) == 0 {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "skip",
			Message: "No services found in service-boundaries.md",
		}
	}

	// Parse contracts from interface-contracts.md
	contracts, err := parseContractsFromFile(interfaceContractsPath)
	if err != nil {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "skip",
			Message: fmt.Sprintf("Could not parse interface-contracts.md: %v", err),
		}
	}

	// Build lookup set for contract domain prefixes
	contractDomains := make(map[string]bool)
	for _, contract := range contracts {
		// Extract domain from IC-XXX-001 format (e.g., IC-CUST-001 -> CUST)
		parts := strings.Split(contract.ID, "-")
		if len(parts) >= 2 {
			contractDomains[parts[1]] = true
		}
	}

	// Validate: every service has at least one matching contract
	servicesWithContract := 0
	servicesWithoutContract := 0

	for _, svc := range services {
		hasContract := false

		// Extract service domain (SVC-CUSTOMER -> CUSTOMER)
		svcDomain := strings.TrimPrefix(svc.ID, "SVC-")

		// Check if any mapped IC domain exists
		if mappedDomains, ok := serviceContractMapping[svcDomain]; ok {
			for _, icDomain := range mappedDomains {
				if contractDomains[icDomain] {
					hasContract = true
					break
				}
			}
		}

		// Fallback: direct match (SVC-ORDER -> IC-ORDER)
		if !hasContract && contractDomains[svcDomain] {
			hasContract = true
		}

		if hasContract {
			servicesWithContract++
		} else {
			servicesWithoutContract++
			result.Errors = append(result.Errors, ValidationError{
				Rule:    RuleV007,
				Message: fmt.Sprintf("Service '%s' has no interface contract", svc.ID),
				RefID:   svc.ID,
			})
		}
	}

	if servicesWithoutContract == 0 {
		return ValidationCheck{
			Rule:    RuleV007,
			Status:  "pass",
			Message: fmt.Sprintf("All %d services have interface contracts", servicesWithContract),
			Count:   servicesWithContract,
		}
	}

	return ValidationCheck{
		Rule:    RuleV007,
		Status:  "fail",
		Message: fmt.Sprintf("%d services have no interface contracts", servicesWithoutContract),
		Count:   servicesWithoutContract,
	}
}

func parseServicesFromBoundaries(filePath string) ([]ServiceInfo, error) {
	var services []ServiceInfo

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Pattern: ## SVC-XXX: Name or ## SVC-XXX – Name
	svcPattern := regexp.MustCompile(`^##\s+(SVC-[A-Z]+)[\s:–—]+(.+)`)
	apiPattern := regexp.MustCompile(`^\*\*API Base:\*\*\s*(.+)`)

	var currentSvc *ServiceInfo

	for i, line := range lines {
		// Check for service header
		if matches := svcPattern.FindStringSubmatch(line); len(matches) > 2 {
			if currentSvc != nil {
				services = append(services, *currentSvc)
			}
			currentSvc = &ServiceInfo{
				ID:      matches[1],
				Name:    strings.TrimSpace(matches[2]),
				LineNum: i + 1,
			}
		}

		// Check for API Base line
		if currentSvc != nil {
			if matches := apiPattern.FindStringSubmatch(line); len(matches) > 1 {
				currentSvc.APIBase = strings.TrimSpace(matches[1])
			}
		}
	}

	// Don't forget the last service
	if currentSvc != nil {
		services = append(services, *currentSvc)
	}

	return services, nil
}

func parseContractsFromFile(filePath string) ([]ContractInfo, error) {
	var contracts []ContractInfo

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Pattern: ## IC-XXX-001 – Name or ## IC-XXX-001: Name
	icPattern := regexp.MustCompile(`^##\s+(IC-[A-Z]+-\d{3})\s*[–:—]+\s*(.+)`)
	urlPattern := regexp.MustCompile(`^\*\*Base URL:\*\*\s*(.+)`)

	var currentIC *ContractInfo

	for i, line := range lines {
		// Check for contract header
		if matches := icPattern.FindStringSubmatch(line); len(matches) > 2 {
			if currentIC != nil {
				contracts = append(contracts, *currentIC)
			}
			currentIC = &ContractInfo{
				ID:      matches[1],
				Name:    strings.TrimSpace(matches[2]),
				LineNum: i + 1,
			}
		}

		// Check for Base URL line
		if currentIC != nil {
			if matches := urlPattern.FindStringSubmatch(line); len(matches) > 1 {
				currentIC.BaseURL = strings.TrimSpace(matches[1])
			}
		}
	}

	// Don't forget the last contract
	if currentIC != nil {
		contracts = append(contracts, *currentIC)
	}

	return contracts, nil
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
	fmt.Println("========================================")

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
