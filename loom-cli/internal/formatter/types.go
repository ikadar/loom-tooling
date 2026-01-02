package formatter

// TestCase represents a single test case for formatting
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
	ShouldNot       string     `json:"should_not,omitempty"` // For hallucination prevention tests
}

// TestData represents test data for a test case
type TestData struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Notes string      `json:"notes"`
}

// TDAISummary holds statistics about generated tests
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

// TechSpec represents a technical specification for formatting
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

// DataReq represents data requirements for a tech spec
type DataReq struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	Source      string `json:"source"`
}

// ErrorHandling represents error handling for a tech spec
type ErrorHandling struct {
	Condition  string `json:"condition"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status"`
}

// InterfaceContract types
type InterfaceContract struct {
	ID                   string               `json:"id"`
	ServiceName          string               `json:"serviceName"`
	Purpose              string               `json:"purpose"`
	BaseURL              string               `json:"baseUrl"`
	Operations           []ContractOperation  `json:"operations"`
	Events               []ContractEvent      `json:"events"`
	SecurityRequirements SecurityRequirements `json:"securityRequirements"`
}

type ContractOperation struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Method         string                 `json:"method"`
	Path           string                 `json:"path"`
	Description    string                 `json:"description"`
	InputSchema    map[string]SchemaField `json:"inputSchema"`
	OutputSchema   map[string]SchemaField `json:"outputSchema"`
	Errors         []ContractError        `json:"errors"`
	Preconditions  []string               `json:"preconditions"`
	Postconditions []string               `json:"postconditions"`
	RelatedACs     []string               `json:"relatedACs"`
	RelatedBRs     []string               `json:"relatedBRs"`
}

type SchemaField struct {
	Type     string `json:"type"`
	Required bool   `json:"required,omitempty"`
}

type ContractError struct {
	Code       string `json:"code"`
	HTTPStatus int    `json:"httpStatus"`
	Message    string `json:"message"`
}

type ContractEvent struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Payload     []string `json:"payload"`
}

type SecurityRequirements struct {
	Authentication string `json:"authentication"`
	Authorization  string `json:"authorization"`
}

type SharedType struct {
	Name   string      `json:"name"`
	Fields []TypeField `json:"fields"`
}

type TypeField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
}

// Aggregate Design types
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

type AggInvariant struct {
	ID          string `json:"id"`
	Rule        string `json:"rule"`
	Enforcement string `json:"enforcement"`
}

type AggRoot struct {
	Entity     string         `json:"entity"`
	Identity   string         `json:"identity"`
	Attributes []AggAttribute `json:"attributes"`
}

type AggAttribute struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Mutable bool   `json:"mutable"`
}

type AggEntity struct {
	Name       string         `json:"name"`
	Identity   string         `json:"identity"`
	Purpose    string         `json:"purpose"`
	Attributes []AggAttribute `json:"attributes"`
}

type AggBehavior struct {
	Name           string   `json:"name"`
	Command        string   `json:"command"`
	Preconditions  []string `json:"preconditions"`
	Postconditions []string `json:"postconditions"`
	Emits          string   `json:"emits"`
}

type AggEvent struct {
	Name    string   `json:"name"`
	Payload []string `json:"payload"`
}

type AggRepository struct {
	Name         string       `json:"name"`
	Methods      []RepoMethod `json:"methods"`
	LoadStrategy string       `json:"loadStrategy"`
	Concurrency  string       `json:"concurrency"`
}

type RepoMethod struct {
	Name    string `json:"name"`
	Params  string `json:"params"`
	Returns string `json:"returns"`
}

type AggExternalRef struct {
	Aggregate string `json:"aggregate"`
	Via       string `json:"via"`
	Type      string `json:"type"`
}

// Sequence Design types
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

type SequenceTrigger struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type SeqParticipant struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SequenceStep struct {
	Step    int      `json:"step"`
	Actor   string   `json:"actor"`
	Action  string   `json:"action"`
	Target  string   `json:"target"`
	Data    []string `json:"data,omitempty"`
	Returns string   `json:"returns,omitempty"`
	Event   string   `json:"event,omitempty"`
}

type SequenceOutcome struct {
	Success      string   `json:"success"`
	StateChanges []string `json:"state_changes"`
}

type SequenceException struct {
	Condition string `json:"condition"`
	Step      int    `json:"step"`
	Handling  string `json:"handling"`
}

// Data Model types
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

type DataField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	Default     string `json:"default,omitempty"`
}

type DataPrimaryKey struct {
	Columns []string `json:"columns"`
}

type DataIndex struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type DataForeignKey struct {
	Columns    []string `json:"columns"`
	References string   `json:"references"`
	OnDelete   string   `json:"onDelete"`
}

type DataConstraint struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

type DataEnum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
