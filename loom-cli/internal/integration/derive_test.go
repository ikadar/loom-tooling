package integration

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/internal/domain"
)

// =============================================================================
// Integration tests for derive workflow using mock LLM responses
// =============================================================================

// TestDeriveWorkflow_DomainModelFormatting tests the complete flow from
// LLM response to formatted markdown output
func TestDeriveWorkflow_DomainModelFormatting(t *testing.T) {
	// Load synthetic domain model response
	responsePath := "../../testdata/responses/domain-model-response.json"
	if _, err := os.Stat(responsePath); os.IsNotExist(err) {
		t.Skip("Synthetic response not found, skipping integration test")
	}

	responseData, err := os.ReadFile(responsePath)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Parse the response as DomainModelDoc would be parsed
	var domainModel struct {
		DomainModel struct {
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			BoundedContexts []string `json:"bounded_contexts"`
		} `json:"domain_model"`
		Entities []struct {
			ID            string `json:"id"`
			Name          string `json:"name"`
			Type          string `json:"type"`
			Purpose       string `json:"purpose"`
			Attributes    []struct {
				Name        string `json:"name"`
				Type        string `json:"type"`
				Constraints string `json:"constraints"`
			} `json:"attributes"`
			Invariants    []string `json:"invariants"`
			Operations    []struct {
				Name           string   `json:"name"`
				Signature      string   `json:"signature"`
				Preconditions  []string `json:"preconditions"`
				Postconditions []string `json:"postconditions"`
			} `json:"operations"`
			Events []struct {
				Name    string   `json:"name"`
				Trigger string   `json:"trigger"`
				Payload []string `json:"payload"`
			} `json:"events"`
			Relationships []struct {
				Target      string `json:"target"`
				Type        string `json:"type"`
				Cardinality string `json:"cardinality"`
			} `json:"relationships"`
		} `json:"entities"`
		ValueObjects []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Purpose    string `json:"purpose"`
			Attributes []struct {
				Name        string `json:"name"`
				Type        string `json:"type"`
				Constraints string `json:"constraints"`
			} `json:"attributes"`
			Operations []string `json:"operations"`
		} `json:"value_objects"`
		Summary struct {
			AggregateRoots  int `json:"aggregate_roots"`
			Entities        int `json:"entities"`
			ValueObjects    int `json:"value_objects"`
			TotalOperations int `json:"total_operations"`
			TotalEvents     int `json:"total_events"`
		} `json:"summary"`
	}

	if err := json.Unmarshal(responseData, &domainModel); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify domain model structure
	if domainModel.DomainModel.Name != "E-Commerce Order Management" {
		t.Errorf("Expected domain name 'E-Commerce Order Management', got '%s'", domainModel.DomainModel.Name)
	}

	if len(domainModel.Entities) < 2 {
		t.Errorf("Expected at least 2 entities, got %d", len(domainModel.Entities))
	}

	// Check entity structure
	foundOrder := false
	for _, entity := range domainModel.Entities {
		if entity.Name == "Order" {
			foundOrder = true
			if entity.Type != "aggregate_root" {
				t.Errorf("Expected Order to be aggregate_root, got '%s'", entity.Type)
			}
			if len(entity.Attributes) < 4 {
				t.Errorf("Expected at least 4 attributes on Order, got %d", len(entity.Attributes))
			}
			if len(entity.Operations) < 3 {
				t.Errorf("Expected at least 3 operations on Order, got %d", len(entity.Operations))
			}
		}
	}
	if !foundOrder {
		t.Error("Expected to find Order entity")
	}

	// Verify value objects
	if len(domainModel.ValueObjects) < 2 {
		t.Errorf("Expected at least 2 value objects, got %d", len(domainModel.ValueObjects))
	}

	// Verify summary
	if domainModel.Summary.AggregateRoots < 1 {
		t.Error("Expected at least 1 aggregate root in summary")
	}
}

// TestDeriveWorkflow_DerivationResponse tests parsing of AC and BR derivation results
func TestDeriveWorkflow_DerivationResponse(t *testing.T) {
	responsePath := "../../testdata/responses/derivation-response.json"
	if _, err := os.Stat(responsePath); os.IsNotExist(err) {
		t.Skip("Synthetic response not found, skipping integration test")
	}

	responseData, err := os.ReadFile(responsePath)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var result struct {
		AcceptanceCriteria []domain.AcceptanceCriteria `json:"acceptance_criteria"`
		BusinessRules      []domain.BusinessRule       `json:"business_rules"`
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify acceptance criteria
	if len(result.AcceptanceCriteria) < 4 {
		t.Errorf("Expected at least 4 acceptance criteria, got %d", len(result.AcceptanceCriteria))
	}

	// Check AC structure
	for _, ac := range result.AcceptanceCriteria {
		if ac.ID == "" {
			t.Error("AC should have an ID")
		}
		if !strings.HasPrefix(ac.ID, "AC-") {
			t.Errorf("AC ID should start with 'AC-', got '%s'", ac.ID)
		}
		if ac.Given == "" || ac.When == "" || ac.Then == "" {
			t.Errorf("AC %s should have Given/When/Then", ac.ID)
		}
	}

	// Verify business rules
	if len(result.BusinessRules) < 4 {
		t.Errorf("Expected at least 4 business rules, got %d", len(result.BusinessRules))
	}

	// Check BR structure
	for _, br := range result.BusinessRules {
		if br.ID == "" {
			t.Error("BR should have an ID")
		}
		if !strings.HasPrefix(br.ID, "BR-") {
			t.Errorf("BR ID should start with 'BR-', got '%s'", br.ID)
		}
		if br.Rule == "" {
			t.Errorf("BR %s should have a rule", br.ID)
		}
	}
}

// TestMockClient_Integration tests that MockClient correctly simulates LLM responses
func TestMockClient_Integration(t *testing.T) {
	// Create mock client
	mock := claude.NewMockClient()

	// Load domain model response
	responsePath := "../../testdata/responses/domain-model-response.json"
	if _, err := os.Stat(responsePath); os.IsNotExist(err) {
		t.Skip("Synthetic response not found")
	}

	responseData, err := os.ReadFile(responsePath)
	if err != nil {
		t.Fatal(err)
	}

	// Add response for domain model derivation
	mock.AddContainsResponse("domain model", string(responseData))

	// Test Call
	result, err := mock.Call("Please derive the domain model from the input")
	if err != nil {
		t.Fatalf("MockClient.Call failed: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty response")
	}

	// Test CallJSON
	var parsed map[string]interface{}
	err = mock.CallJSON("derive domain model", &parsed)
	if err != nil {
		t.Fatalf("MockClient.CallJSON failed: %v", err)
	}

	// Verify parsed content
	if _, ok := parsed["entities"]; !ok {
		t.Error("Expected 'entities' in parsed response")
	}
	if _, ok := parsed["value_objects"]; !ok {
		t.Error("Expected 'value_objects' in parsed response")
	}

	// Check call count
	if mock.GetCallCount() != 2 {
		t.Errorf("Expected 2 calls, got %d", mock.GetCallCount())
	}
}

// TestDomainTypes_AcceptanceCriteria tests domain type construction and serialization
func TestDomainTypes_AcceptanceCriteria(t *testing.T) {
	// Create sample acceptance criteria using domain types
	acs := []domain.AcceptanceCriteria{
		{
			ID:         "AC-ORD-001",
			Title:      "Create Order",
			Given:      "A registered customer with valid payment",
			When:       "The customer creates a new order",
			Then:       "The order is created with PENDING status",
			ErrorCases: []string{"ERROR when customer not found"},
			SourceRefs: []string{"US-001"},
		},
		{
			ID:         "AC-ORD-002",
			Title:      "Submit Order",
			Given:      "An order with items",
			When:       "The customer submits the order",
			Then:       "The order status changes to SUBMITTED",
			ErrorCases: []string{"ERROR when no items"},
			SourceRefs: []string{"US-002"},
		},
	}

	// Verify structure
	if len(acs) != 2 {
		t.Errorf("Expected 2 ACs, got %d", len(acs))
	}

	// Verify ID format
	for _, ac := range acs {
		if !strings.HasPrefix(ac.ID, "AC-") {
			t.Errorf("AC ID should start with 'AC-', got '%s'", ac.ID)
		}
		if ac.Given == "" || ac.When == "" || ac.Then == "" {
			t.Errorf("AC %s should have Given/When/Then", ac.ID)
		}
	}

	// Test JSON serialization
	data, err := json.Marshal(acs)
	if err != nil {
		t.Fatalf("Failed to marshal ACs: %v", err)
	}

	var parsed []domain.AcceptanceCriteria
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal ACs: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected 2 parsed ACs, got %d", len(parsed))
	}
}

// TestDomainTypes_BusinessRules tests business rule domain types
func TestDomainTypes_BusinessRules(t *testing.T) {
	brs := []domain.BusinessRule{
		{
			ID:          "BR-ORD-001",
			Title:       "Order Total Calculation",
			Rule:        "Order total must equal sum of line items",
			Invariant:   "total = sum(items.subtotal)",
			Enforcement: "Calculated on item add/update",
			ErrorCode:   "ORDER_TOTAL_MISMATCH",
			SourceRefs:  []string{"US-001"},
		},
	}

	// Verify structure
	br := brs[0]
	if !strings.HasPrefix(br.ID, "BR-") {
		t.Errorf("BR ID should start with 'BR-', got '%s'", br.ID)
	}
	if br.Rule == "" {
		t.Error("BR should have a rule")
	}
	if br.ErrorCode == "" {
		t.Error("BR should have an error code")
	}

	// Test JSON serialization
	data, err := json.Marshal(brs)
	if err != nil {
		t.Fatalf("Failed to marshal BRs: %v", err)
	}

	var parsed []domain.BusinessRule
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal BRs: %v", err)
	}

	if parsed[0].ErrorCode != "ORDER_TOTAL_MISMATCH" {
		t.Errorf("Expected error code 'ORDER_TOTAL_MISMATCH', got '%s'", parsed[0].ErrorCode)
	}
}

// TestDomainTypes_DerivationResult tests full derivation result
func TestDomainTypes_DerivationResult(t *testing.T) {
	responsePath := "../../testdata/responses/derivation-response.json"
	if _, err := os.Stat(responsePath); os.IsNotExist(err) {
		t.Skip("Synthetic response not found")
	}

	responseData, err := os.ReadFile(responsePath)
	if err != nil {
		t.Fatal(err)
	}

	// Parse response as derivation result
	var result struct {
		AcceptanceCriteria []domain.AcceptanceCriteria `json:"acceptance_criteria"`
		BusinessRules      []domain.BusinessRule       `json:"business_rules"`
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		t.Fatal(err)
	}

	// Verify acceptance criteria
	if len(result.AcceptanceCriteria) < 4 {
		t.Errorf("Expected at least 4 ACs, got %d", len(result.AcceptanceCriteria))
	}

	// Verify business rules
	if len(result.BusinessRules) < 4 {
		t.Errorf("Expected at least 4 BRs, got %d", len(result.BusinessRules))
	}

	// Verify all ACs have required fields
	for _, ac := range result.AcceptanceCriteria {
		if ac.ID == "" || ac.Title == "" || ac.Given == "" || ac.When == "" || ac.Then == "" {
			t.Errorf("AC %s missing required fields", ac.ID)
		}
	}

	// Verify all BRs have required fields
	for _, br := range result.BusinessRules {
		if br.ID == "" || br.Title == "" || br.Rule == "" {
			t.Errorf("BR %s missing required fields", br.ID)
		}
	}
}
