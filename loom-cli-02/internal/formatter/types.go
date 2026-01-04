// Package formatter provides JSON to Markdown formatting.
//
// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter

// =============================================================================
// Test Case Types
// =============================================================================

// TestCase represents a single test case.
type TestCase struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Category        string     `json:"category"` // positive, negative, boundary, hallucination
	ACRef           string     `json:"ac_ref"`
	BRRefs          []string   `json:"br_refs"`
	Preconditions   []string   `json:"preconditions"`
	TestData        []TestData `json:"test_data"`
	Steps           []string   `json:"steps"`
	ExpectedResults []string   `json:"expected_results"`
	ShouldNot       string     `json:"should_not,omitempty"`
}

// TestData represents test data for a test case.
type TestData struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Notes string      `json:"notes"`
}

// TestSuite groups test cases by AC.
type TestSuite struct {
	ACRef   string     `json:"ac_ref"`
	ACTitle string     `json:"ac_title"`
	Tests   []TestCase `json:"tests"`
}

// TDAISummary contains test case statistics.
type TDAISummary struct {
	Total      int `json:"total"`
	ByCategory struct {
		Positive      int `json:"positive"`
		Negative      int `json:"negative"`
		Boundary      int `json:"boundary"`
		Hallucination int `json:"hallucination"`
	} `json:"by_category"`
	Coverage struct {
		ACsCovered            int     `json:"acs_covered"`
		PositiveRatio         float64 `json:"positive_ratio"`
		NegativeRatio         float64 `json:"negative_ratio"`
		HasHallucinationTests bool    `json:"has_hallucination_tests"`
	} `json:"coverage"`
}

// =============================================================================
// Tech Spec Types
// =============================================================================

// TechSpec represents a technical specification.
type TechSpec struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	BRRef            string          `json:"br_ref"`
	Rule             string          `json:"rule"`
	Implementation   string          `json:"implementation"`
	ValidationPoints []string        `json:"validation_points"`
	DataRequirements []DataReq       `json:"data_requirements"`
	ErrorHandling    []ErrorHandling `json:"error_handling"`
	RelatedACs       []string        `json:"related_acs"`
}

// DataReq represents a data requirement.
type DataReq struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	Source      string `json:"source"`
}

// ErrorHandling represents error handling specification.
type ErrorHandling struct {
	Condition  string `json:"condition"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status"`
}

// =============================================================================
// Interface Contract Types
// =============================================================================

// InterfaceContract represents an API contract.
type InterfaceContract struct {
	ID                   string               `json:"id"`
	ServiceName          string               `json:"serviceName"`
	Purpose              string               `json:"purpose"`
	BaseURL              string               `json:"baseUrl"`
	Operations           []ContractOperation  `json:"operations"`
	Events               []ContractEvent      `json:"events"`
	SecurityRequirements SecurityRequirements `json:"securityRequirements"`
}

// ContractOperation represents an API operation.
type ContractOperation struct {
	ID          string              `json:"id"`
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Request     ContractRequest     `json:"request"`
	Response    ContractResponse    `json:"response"`
	Errors      []ContractError     `json:"errors"`
	RelatedACs  []string            `json:"relatedACs"`
}

// ContractRequest represents a request specification.
type ContractRequest struct {
	ContentType string                 `json:"contentType"`
	Schema      map[string]interface{} `json:"schema"`
	Example     interface{}            `json:"example"`
}

// ContractResponse represents a response specification.
type ContractResponse struct {
	StatusCode  int                    `json:"statusCode"`
	ContentType string                 `json:"contentType"`
	Schema      map[string]interface{} `json:"schema"`
	Example     interface{}            `json:"example"`
}

// ContractError represents an error response.
type ContractError struct {
	StatusCode  int    `json:"statusCode"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ContractEvent represents an event specification.
type ContractEvent struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Payload     map[string]interface{} `json:"payload"`
}

// SecurityRequirements represents security requirements.
type SecurityRequirements struct {
	Authentication string   `json:"authentication"`
	Authorization  string   `json:"authorization"`
	Scopes         []string `json:"scopes"`
}

// SharedType represents a shared type definition.
type SharedType struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
}

// =============================================================================
// Aggregate Design Types
// =============================================================================

// AggregateDesign represents an aggregate design.
type AggregateDesign struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Purpose            string           `json:"purpose"`
	Invariants         []AggInvariant   `json:"invariants"`
	Root               AggRoot          `json:"root"`
	Entities           []AggEntity      `json:"entities"`
	ValueObjects       []string         `json:"valueObjects"`
	Behaviors          []AggBehavior    `json:"behaviors"`
	Events             []AggEvent       `json:"events"`
	Repository         AggRepository    `json:"repository"`
	ExternalReferences []AggExternalRef `json:"externalReferences"`
}

// AggInvariant represents an aggregate invariant.
type AggInvariant struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Enforcement string `json:"enforcement"`
}

// AggRoot represents the aggregate root entity.
type AggRoot struct {
	Name       string            `json:"name"`
	Attributes []AggAttribute    `json:"attributes"`
	Methods    []string          `json:"methods"`
}

// AggAttribute represents an entity attribute.
type AggAttribute struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// AggEntity represents an entity within the aggregate.
type AggEntity struct {
	Name       string         `json:"name"`
	Attributes []AggAttribute `json:"attributes"`
}

// AggBehavior represents a behavior/method.
type AggBehavior struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parameters  []string `json:"parameters"`
	Returns     string   `json:"returns"`
	Raises      []string `json:"raises"`
}

// AggEvent represents a domain event.
type AggEvent struct {
	Name    string   `json:"name"`
	Trigger string   `json:"trigger"`
	Payload []string `json:"payload"`
}

// AggRepository represents repository operations.
type AggRepository struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

// AggExternalRef represents an external reference.
type AggExternalRef struct {
	Aggregate string `json:"aggregate"`
	Type      string `json:"type"` // reference, lookup
	Via       string `json:"via"`
}

// =============================================================================
// Sequence Design Types
// =============================================================================

// SequenceDesign represents a sequence diagram design.
type SequenceDesign struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Trigger      SequenceTrigger     `json:"trigger"`
	Participants []SeqParticipant    `json:"participants"`
	Steps        []SequenceStep      `json:"steps"`
	Outcome      SequenceOutcome     `json:"outcome"`
	Exceptions   []SequenceException `json:"exceptions"`
	RelatedACs   []string            `json:"relatedACs"`
	RelatedBRs   []string            `json:"relatedBRs"`
}

// SequenceTrigger represents what triggers the sequence.
type SequenceTrigger struct {
	Actor  string `json:"actor"`
	Action string `json:"action"`
}

// SeqParticipant represents a sequence participant.
type SeqParticipant struct {
	Name string `json:"name"`
	Type string `json:"type"` // actor, service, database, external
}

// SequenceStep represents a step in the sequence.
type SequenceStep struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Action  string `json:"action"`
	Returns string `json:"returns,omitempty"`
	Async   bool   `json:"async,omitempty"`
}

// SequenceOutcome represents the sequence outcome.
type SequenceOutcome struct {
	Success string `json:"success"`
	Result  string `json:"result"`
}

// SequenceException represents an exception flow.
type SequenceException struct {
	Condition string `json:"condition"`
	Handler   string `json:"handler"`
	Result    string `json:"result"`
}

// =============================================================================
// Data Model Types
// =============================================================================

// DataTable represents a database table.
type DataTable struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Aggregate        string           `json:"aggregate"`
	Purpose          string           `json:"purpose"`
	Fields           []DataField      `json:"fields"`
	PrimaryKey       DataPrimaryKey   `json:"primaryKey"`
	Indexes          []DataIndex      `json:"indexes"`
	ForeignKeys      []DataForeignKey `json:"foreignKeys"`
	CheckConstraints []DataConstraint `json:"checkConstraints"`
}

// DataField represents a table field.
type DataField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Nullable    bool   `json:"nullable"`
	Default     string `json:"default"`
	Constraints string `json:"constraints"`
}

// DataPrimaryKey represents a primary key.
type DataPrimaryKey struct {
	Columns []string `json:"columns"`
	Type    string   `json:"type"` // single, composite
}

// DataIndex represents a database index.
type DataIndex struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// DataForeignKey represents a foreign key.
type DataForeignKey struct {
	Name       string   `json:"name"`
	Columns    []string `json:"columns"`
	References string   `json:"references"`
	OnDelete   string   `json:"onDelete"`
	OnUpdate   string   `json:"onUpdate"`
}

// DataConstraint represents a check constraint.
type DataConstraint struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

// DataEnum represents an enum type.
type DataEnum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
