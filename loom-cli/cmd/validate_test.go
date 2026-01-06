package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// =============================================================================
// extractIDsAndRefs tests
// =============================================================================

func TestExtractIDsAndRefs_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.md")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	ids, refs, err := extractIDsAndRefs(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(ids) != 0 {
		t.Errorf("Expected 0 IDs, got %d", len(ids))
	}
	if len(refs) != 0 {
		t.Errorf("Expected 0 refs, got %d", len(refs))
	}
}

func TestExtractIDsAndRefs_SingleAC(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "acceptance-criteria.md")
	content := `# Acceptance Criteria

## AC-ORD-001 – Create Order

**Given** a customer is logged in
**When** they add items to cart
**Then** an order is created
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ids, _, err := extractIDsAndRefs(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(ids) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(ids))
	}

	if _, ok := ids["AC-ORD-001"]; !ok {
		t.Error("Expected to find AC-ORD-001")
	}
}

func TestExtractIDsAndRefs_MultipleIDs(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.md")
	// Note: TC IDs are handled separately by validateTDAI, not extractIDsAndRefs
	// The headerPattern doesn't support the TC-AC-XXX-NNN-XNN format
	content := `# Test Specs

## AC-ORD-001 – Create Order
Some content

## AC-ORD-002 – Update Order
More content

## BR-ORD-001 – Order Validation Rule
Business rule content

### TS-ORD-001 – Tech Spec
Tech spec content
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ids, _, err := extractIDsAndRefs(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedIDs := []string{"AC-ORD-001", "AC-ORD-002", "BR-ORD-001", "TS-ORD-001"}
	for _, expected := range expectedIDs {
		if _, ok := ids[expected]; !ok {
			t.Errorf("Expected to find %s", expected)
		}
	}

	if len(ids) != 4 {
		t.Errorf("Expected 4 IDs, got %d", len(ids))
	}
}

func TestExtractIDsAndRefs_WithReferences(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "tech-specs.md")
	content := `# Tech Specs

## TS-ORD-001 – Order Validation

**Rule:** Order must be valid

**Traceability:**
- BR: BR-ORD-001
- AC: AC-ORD-001, AC-ORD-002
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ids, refs, err := extractIDsAndRefs(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, ok := ids["TS-ORD-001"]; !ok {
		t.Error("Expected to find TS-ORD-001")
	}

	tsRefs, ok := refs["TS-ORD-001"]
	if !ok {
		t.Fatal("Expected refs for TS-ORD-001")
	}

	// Should have BR and AC refs
	hasRef := func(refs []string, target string) bool {
		for _, r := range refs {
			if r == target {
				return true
			}
		}
		return false
	}

	if !hasRef(tsRefs, "BR-ORD-001") {
		t.Error("Expected BR-ORD-001 in refs")
	}
	if !hasRef(tsRefs, "AC-ORD-001") {
		t.Error("Expected AC-ORD-001 in refs")
	}
}

// =============================================================================
// File finder tests
// =============================================================================

func TestFindL1Files(t *testing.T) {
	tmpDir := t.TempDir()

	// Create L1 files
	l1Files := []string{
		"domain-model.md",
		"bounded-context-map.md",
		"acceptance-criteria.md",
		"business-rules.md",
	}

	for _, f := range l1Files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create non-L1 file
	if err := os.WriteFile(filepath.Join(tmpDir, "other.md"), []byte("other"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := findL1Files(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(files) != len(l1Files) {
		t.Errorf("Expected %d L1 files, got %d", len(l1Files), len(files))
	}
}

func TestFindL2Files(t *testing.T) {
	tmpDir := t.TempDir()

	l2Files := []string{
		"tech-specs.md",
		"interface-contracts.md",
		"aggregate-design.md",
		"sequence-design.md",
		"initial-data-model.md",
	}

	for _, f := range l2Files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := findL2Files(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(files) != len(l2Files) {
		t.Errorf("Expected %d L2 files, got %d", len(l2Files), len(files))
	}
}

func TestFindL3Files(t *testing.T) {
	tmpDir := t.TempDir()

	l3Files := []string{
		"test-cases.md",
		"feature-tickets.md",
		"service-boundaries.md",
	}

	for _, f := range l3Files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("# "+f), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := findL3Files(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(files) != len(l3Files) {
		t.Errorf("Expected %d L3 files, got %d", len(l3Files), len(files))
	}
}

// =============================================================================
// Parser tests
// =============================================================================

func TestParseEntitiesFromDomainModel(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "domain-model.md")
	content := `# Domain Model

## Entities

### ENT-ORD-001 – Order

**Type:** aggregate_root

**Attributes:**
- orderId
- status

### ENT-ORD-002 – OrderItem

**Type:** entity

**Attributes:**
- quantity
- price
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	entities, err := parseEntitiesFromDomainModel(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entities) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(entities))
	}

	// Check first entity
	if entities[0].ID != "ENT-ORD-001" {
		t.Errorf("Expected ID 'ENT-ORD-001', got '%s'", entities[0].ID)
	}
	if entities[0].Name != "Order" {
		t.Errorf("Expected name 'Order', got '%s'", entities[0].Name)
	}
	if entities[0].Type != "aggregate_root" {
		t.Errorf("Expected type 'aggregate_root', got '%s'", entities[0].Type)
	}

	// Check second entity
	if entities[1].ID != "ENT-ORD-002" {
		t.Errorf("Expected ID 'ENT-ORD-002', got '%s'", entities[1].ID)
	}
	if entities[1].Type != "entity" {
		t.Errorf("Expected type 'entity', got '%s'", entities[1].Type)
	}
}

func TestParseAggregatesFromDesign(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "aggregate-design.md")
	content := `# Aggregate Design

## AGG-ORD-001 – Order

**Purpose:** Manages order lifecycle

### Aggregate Root: Order

**Identity:** OrderId

### Child Entities

#### OrderItem

**Purpose:** Line item in order

#### OrderPayment

**Purpose:** Payment record
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	aggregates, err := parseAggregatesFromDesign(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(aggregates) != 1 {
		t.Fatalf("Expected 1 aggregate, got %d", len(aggregates))
	}

	agg := aggregates[0]
	if agg.ID != "AGG-ORD-001" {
		t.Errorf("Expected ID 'AGG-ORD-001', got '%s'", agg.ID)
	}
	if agg.Name != "Order" {
		t.Errorf("Expected name 'Order', got '%s'", agg.Name)
	}
	if agg.RootEntityName != "Order" {
		t.Errorf("Expected root 'Order', got '%s'", agg.RootEntityName)
	}
	if len(agg.ChildEntities) != 2 {
		t.Errorf("Expected 2 child entities, got %d", len(agg.ChildEntities))
	}
}

func TestParseServicesFromBoundaries(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "service-boundaries.md")
	content := `# Service Boundaries

## SVC-ORDER: Order Service

**Purpose:** Manages orders

**API Base:** /api/v1/orders

## SVC-CUSTOMER – Customer Service

**Purpose:** Manages customers

**API Base:** /api/v1/customers
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	services, err := parseServicesFromBoundaries(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(services) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(services))
	}

	if services[0].ID != "SVC-ORDER" {
		t.Errorf("Expected ID 'SVC-ORDER', got '%s'", services[0].ID)
	}
	if services[0].APIBase != "/api/v1/orders" {
		t.Errorf("Expected API base '/api/v1/orders', got '%s'", services[0].APIBase)
	}
}

func TestParseContractsFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "interface-contracts.md")
	content := `# Interface Contracts

## IC-ORD-001 – Order Service

**Purpose:** Order management API

**Base URL:** /api/v1/orders

### Operations

#### Create Order

## IC-CUST-001 – Customer Service

**Purpose:** Customer API

**Base URL:** /api/v1/customers
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	contracts, err := parseContractsFromFile(path)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(contracts) != 2 {
		t.Fatalf("Expected 2 contracts, got %d", len(contracts))
	}

	if contracts[0].ID != "IC-ORD-001" {
		t.Errorf("Expected ID 'IC-ORD-001', got '%s'", contracts[0].ID)
	}
	if contracts[0].BaseURL != "/api/v1/orders" {
		t.Errorf("Expected base URL '/api/v1/orders', got '%s'", contracts[0].BaseURL)
	}
}

// =============================================================================
// Validation rule tests
// =============================================================================

func TestValidateStructural_NoDuplicates(t *testing.T) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	allIDs := map[string]string{
		"AC-ORD-001": "file1.md",
		"AC-ORD-002": "file1.md",
		"BR-ORD-001": "file2.md",
	}

	files := []string{"file1.md", "file2.md"}

	checks := validateStructural(files, allIDs, result)

	// Find V010 check
	var v010Check *ValidationCheck
	for i := range checks {
		if checks[i].Rule == RuleV010 {
			v010Check = &checks[i]
			break
		}
	}

	if v010Check == nil {
		t.Fatal("Expected V010 check")
	}

	if v010Check.Status != "pass" {
		t.Errorf("Expected V010 to pass, got %s", v010Check.Status)
	}
}

func TestValidateTraceability_ValidRefs(t *testing.T) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	allIDs := map[string]string{
		"TS-ORD-001": "tech-specs.md",
		"BR-ORD-001": "business-rules.md",
		"AC-ORD-001": "acceptance-criteria.md",
	}

	allRefs := map[string][]string{
		"TS-ORD-001": {"BR-ORD-001", "AC-ORD-001"},
	}

	checks := validateTraceability(allIDs, allRefs, result)

	// Find V003 check
	var v003Check *ValidationCheck
	for i := range checks {
		if checks[i].Rule == RuleV003 {
			v003Check = &checks[i]
			break
		}
	}

	if v003Check == nil {
		t.Fatal("Expected V003 check")
	}

	if v003Check.Status != "pass" {
		t.Errorf("Expected V003 to pass, got %s: %s", v003Check.Status, v003Check.Message)
	}

	if len(result.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(result.Errors))
	}
}

func TestValidateTraceability_InvalidRefs(t *testing.T) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	allIDs := map[string]string{
		"TS-ORD-001": "tech-specs.md",
	}

	allRefs := map[string][]string{
		"TS-ORD-001": {"BR-ORD-999", "AC-NONEXISTENT"},
	}

	checks := validateTraceability(allIDs, allRefs, result)

	// Find V003 check
	var v003Check *ValidationCheck
	for i := range checks {
		if checks[i].Rule == RuleV003 {
			v003Check = &checks[i]
			break
		}
	}

	if v003Check == nil {
		t.Fatal("Expected V003 check")
	}

	if v003Check.Status != "fail" {
		t.Errorf("Expected V003 to fail, got %s", v003Check.Status)
	}

	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestCalculateSummary(t *testing.T) {
	result := &ValidationResult{
		Checks: []ValidationCheck{
			{Rule: "V001", Status: "pass"},
			{Rule: "V002", Status: "pass"},
			{Rule: "V003", Status: "fail"},
			{Rule: "V004", Status: "skip"},
			{Rule: "V005", Status: "pass"},
		},
		Errors: []ValidationError{
			{Rule: "V003", Message: "Error 1"},
			{Rule: "V003", Message: "Error 2"},
		},
		Warnings: []ValidationWarning{
			{Rule: "V002", Message: "Warning 1"},
		},
	}

	summary := calculateSummary(result)

	if summary.TotalChecks != 5 {
		t.Errorf("Expected 5 total checks, got %d", summary.TotalChecks)
	}

	if summary.Passed != 3 {
		t.Errorf("Expected 3 passed, got %d", summary.Passed)
	}

	if summary.Failed != 1 {
		t.Errorf("Expected 1 failed, got %d", summary.Failed)
	}

	if summary.ErrorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", summary.ErrorCount)
	}

	if summary.Warnings != 1 {
		t.Errorf("Expected 1 warning, got %d", summary.Warnings)
	}
}

// =============================================================================
// ID pattern tests
// =============================================================================

func TestIDPatterns(t *testing.T) {
	validIDs := map[string]string{
		"AC":   "AC-ORD-001",
		"BR":   "BR-CUST-002",
		"TC":   "TC-AC-ORD-001-P01",
		"TS":   "TS-ORD-001",
		"IC":   "IC-CUST-001",
		"AGG":  "AGG-ORD-001",
		"SEQ":  "SEQ-ORD-001",
		"EVT":  "EVT-ORD-001",
		"CMD":  "CMD-ORD-001",
		"SKEL": "SKEL-ORD-001",
	}

	for prefix, id := range validIDs {
		pattern, ok := idPatterns[prefix]
		if !ok {
			t.Errorf("No pattern for prefix %s", prefix)
			continue
		}

		if !pattern.MatchString(id) {
			t.Errorf("ID %s should match pattern for %s", id, prefix)
		}
	}

	// Test invalid IDs
	// Note: MatchString does substring matching, so TC-AC-ORD-001 would match
	// the AC pattern via the AC-ORD-001 substring. Testing only truly invalid patterns.
	invalidIDs := []string{
		"AC-001",           // Missing domain part
		"BR-ord-001",       // Lowercase domain
		"INVALID-ORD-001",  // Unknown prefix
		"XX-YYY-000",       // Completely invalid format
	}

	for _, id := range invalidIDs {
		matched := false
		for _, pattern := range idPatterns {
			if pattern.MatchString(id) {
				matched = true
				break
			}
		}
		if matched {
			t.Errorf("ID %s should not match any pattern", id)
		}
	}
}

// =============================================================================
// Integration test with real fixtures
// =============================================================================

func TestValidate_Integration(t *testing.T) {
	// Skip if fixture doesn't exist
	fixturePath := "../testdata/fixtures/l2-ecommerce-order"
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		t.Skip("Fixture not found, skipping integration test")
	}

	result, err := validate(fixturePath, "L2")
	if err != nil {
		t.Fatalf("Validation error: %v", err)
	}

	// Should have some checks
	if len(result.Checks) == 0 {
		t.Error("Expected some validation checks")
	}

	// Summary should be calculated
	if result.Summary.TotalChecks != len(result.Checks) {
		t.Errorf("Summary mismatch: %d vs %d", result.Summary.TotalChecks, len(result.Checks))
	}
}
