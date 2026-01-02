package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikadar/loom-cli/internal/claude"
	"github.com/ikadar/loom-cli/internal/formatter"
	"github.com/ikadar/loom-cli/internal/generator"
	"github.com/ikadar/loom-cli/internal/workflow"
	"github.com/ikadar/loom-cli/prompts"
)

// toAnchorL2 converts an ID to a lowercase anchor (e.g., "TC-AC-CUST-001-P01" -> "tc-ac-cust-001-p01")
func toAnchorL2(id string) string {
	return strings.ToLower(id)
}

// toLinkL2 creates a markdown link with anchor
func toLinkL2(id, file string) string {
	return fmt.Sprintf("[%s](%s#%s)", id, file, toAnchorL2(id))
}

// L2Result is the output of the derive-l2 command
type L2Result struct {
	Summary    L2Summary    `json:"summary"`
	TestCases  []TestCase   `json:"test_cases"`
	TechSpecs  []TechSpec   `json:"tech_specs"`
}

type L2Summary struct {
	TestCasesGenerated int         `json:"test_cases_generated"`
	TechSpecsGenerated int         `json:"tech_specs_generated"`
	Coverage           L2Coverage  `json:"coverage"`
}

type L2Coverage struct {
	ACsCovered      int `json:"acs_covered"`
	BRsCovered      int `json:"brs_covered"`
	HappyPathTests  int `json:"happy_path_tests"`
	ErrorTests      int `json:"error_tests"`
	EdgeCaseTests   int `json:"edge_case_tests"`
}

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

type TestSuite struct {
	ACRef   string     `json:"ac_ref"`
	ACTitle string     `json:"ac_title"`
	Tests   []TestCase `json:"tests"`
}

type TestData struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Notes string      `json:"notes"`
}

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

// flattenTestSuites converts generator.TestSuite to local TestCase slice
func flattenTestSuites(suites []generator.TestSuite) []TestCase {
	var result []TestCase
	for _, suite := range suites {
		for _, tc := range suite.Tests {
			result = append(result, TestCase{
				ID:              tc.ID,
				Name:            tc.Name,
				Category:        tc.Category,
				ACRef:           tc.ACRef,
				BRRefs:          tc.BRRefs,
				Preconditions:   tc.Preconditions,
				TestData:        convertTestData(tc.TestData),
				Steps:           tc.Steps,
				ExpectedResults: tc.ExpectedResults,
				ShouldNot:       tc.ShouldNot,
			})
		}
	}
	return result
}

// convertTestData converts generator.TestData to local TestData
func convertTestData(data []generator.TestData) []TestData {
	result := make([]TestData, len(data))
	for i, d := range data {
		result[i] = TestData{
			Field: d.Field,
			Value: d.Value,
			Notes: d.Notes,
		}
	}
	return result
}

// convertTDAISummary converts generator.TDAISummary to local TDAISummary
func convertTDAISummary(s generator.TDAISummary) TDAISummary {
	var result TDAISummary
	result.Total = s.Total
	result.ByCategory.Positive = s.ByCategory.Positive
	result.ByCategory.Negative = s.ByCategory.Negative
	result.ByCategory.Boundary = s.ByCategory.Boundary
	result.ByCategory.Hallucination = s.ByCategory.Hallucination
	result.Coverage.ACsCovered = s.Coverage.ACsCovered
	result.Coverage.PositiveRatio = s.Coverage.PositiveRatio
	result.Coverage.NegativeRatio = s.Coverage.NegativeRatio
	result.Coverage.HasHallucinationTests = s.Coverage.HasHallucinationTests
	return result
}

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

type DataReq struct {
	Field       string `json:"field"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
	Source      string `json:"source"`
}

type ErrorHandling struct {
	Condition  string `json:"condition"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status"`
}

// Interface Contract types
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
	ID             string                       `json:"id"`
	Name           string                       `json:"name"`
	Method         string                       `json:"method"`
	Path           string                       `json:"path"`
	Description    string                       `json:"description"`
	InputSchema    map[string]SchemaField       `json:"inputSchema"`
	OutputSchema   map[string]SchemaField       `json:"outputSchema"`
	Errors         []ContractError              `json:"errors"`
	Preconditions  []string                     `json:"preconditions"`
	Postconditions []string                     `json:"postconditions"`
	RelatedACs     []string                     `json:"relatedACs"`
	RelatedBRs     []string                     `json:"relatedBRs"`
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
	Name   string       `json:"name"`
	Fields []TypeField  `json:"fields"`
}

type TypeField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Constraints string `json:"constraints"`
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
	Step    int    `json:"step"`
	Actor   string `json:"actor"`
	Action  string `json:"action"`
	Target  string `json:"target"`
	Data    []string `json:"data,omitempty"`
	Returns string `json:"returns,omitempty"`
	Event   string `json:"event,omitempty"`
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

// Aggregate Design types
type AggregateDesign struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	Purpose            string               `json:"purpose"`
	Invariants         []AggInvariant       `json:"invariants"`
	Root               AggRoot              `json:"root"`
	Entities           []AggEntity          `json:"entities"`
	ValueObjects       []string             `json:"valueObjects"`
	Behaviors          []AggBehavior        `json:"behaviors"`
	Events             []AggEvent           `json:"events"`
	Repository         AggRepository        `json:"repository"`
	ExternalReferences []AggExternalRef     `json:"externalReferences"`
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
	Name         string         `json:"name"`
	Methods      []RepoMethod   `json:"methods"`
	LoadStrategy string         `json:"loadStrategy"`
	Concurrency  string         `json:"concurrency"`
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

// Data Model types
type DataTable struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Aggregate        string            `json:"aggregate"`
	Purpose          string            `json:"purpose"`
	Fields           []DataField       `json:"fields"`
	PrimaryKey       DataPrimaryKey    `json:"primaryKey"`
	Indexes          []DataIndex       `json:"indexes"`
	ForeignKeys      []DataForeignKey  `json:"foreignKeys"`
	CheckConstraints []DataConstraint  `json:"checkConstraints"`
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

func runDeriveL2() error {
	// Parse arguments
	args := os.Args[2:]

	var inputDir string
	var outputDir string
	var interactive bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input-dir":
			if i+1 < len(args) {
				i++
				inputDir = args[i]
			}
		case "--output-dir":
			if i+1 < len(args) {
				i++
				outputDir = args[i]
			}
		case "--interactive", "-i":
			interactive = true
		}
	}

	if inputDir == "" {
		return fmt.Errorf("--input-dir is required (directory containing L1 documents)")
	}
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required")
	}

	// Read L1 documents
	fmt.Fprintln(os.Stderr, "Phase L2-0: Reading L1 documents...")

	acContent, err := os.ReadFile(filepath.Join(inputDir, "acceptance-criteria.md"))
	if err != nil {
		return fmt.Errorf("failed to read acceptance-criteria.md: %w", err)
	}

	brContent, err := os.ReadFile(filepath.Join(inputDir, "business-rules.md"))
	if err != nil {
		return fmt.Errorf("failed to read business-rules.md: %w", err)
	}

	dmContent, err := os.ReadFile(filepath.Join(inputDir, "domain-model.md"))
	if err != nil {
		return fmt.Errorf("failed to read domain-model.md: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Read: acceptance-criteria.md (%d bytes)\n", len(acContent))
	fmt.Fprintf(os.Stderr, "  Read: business-rules.md (%d bytes)\n", len(brContent))
	fmt.Fprintf(os.Stderr, "  Read: domain-model.md (%d bytes)\n", len(dmContent))

	// Show interactive mode header
	if interactive {
		workflow.PrintInteractiveHeader()
	}

	// Create Claude client
	client := claude.NewClient()

	// Phase 1: Generate Test Cases from ACs (TDAI methodology) - CHUNKED
	fmt.Fprintln(os.Stderr, "\nPhase L2-1: Generating TDAI Test Cases from Acceptance Criteria...")

	tcGenerator := generator.NewChunkedTestCaseGenerator(client)
	tcResult, err := tcGenerator.Generate(string(acContent))
	if err != nil {
		return fmt.Errorf("failed to generate test cases: %w", err)
	}

	// Flatten test suites into test cases for compatibility
	allTestCases := flattenTestSuites(tcResult.TestSuites)

	fmt.Fprintf(os.Stderr, "  Generated: %d Test Cases (P:%d N:%d B:%d H:%d)\n",
		tcResult.Summary.Total,
		tcResult.Summary.ByCategory.Positive,
		tcResult.Summary.ByCategory.Negative,
		tcResult.Summary.ByCategory.Boundary,
		tcResult.Summary.ByCategory.Hallucination)

	// Phases 2-6: Run in parallel
	fmt.Fprintln(os.Stderr, "\nPhases L2-2 to L2-6: Generating remaining L2 artifacts in parallel...")

	// Prepare inputs
	l1Input := string(dmContent) + "\n\n---\n\n" + string(brContent) + "\n\n---\n\n" + string(acContent)
	aggSeqInput := string(dmContent) + "\n\n---\n\n" + string(brContent)

	// Define result types for each phase
	type TechSpecsResult struct {
		TechSpecs []TechSpec `json:"tech_specs"`
		Summary   struct {
			Total int `json:"total"`
		} `json:"summary"`
	}

	type InterfaceContractsResult struct {
		InterfaceContracts []InterfaceContract `json:"interface_contracts"`
		SharedTypes        []SharedType        `json:"shared_types"`
		Summary            struct {
			TotalContracts   int `json:"total_contracts"`
			TotalOperations  int `json:"total_operations"`
			TotalEvents      int `json:"total_events"`
			TotalSharedTypes int `json:"total_shared_types"`
		} `json:"summary"`
	}

	type AggregateResult struct {
		Aggregates []AggregateDesign `json:"aggregates"`
		Summary    struct {
			TotalAggregates int `json:"total_aggregates"`
			TotalInvariants int `json:"total_invariants"`
			TotalBehaviors  int `json:"total_behaviors"`
			TotalEvents     int `json:"total_events"`
		} `json:"summary"`
	}

	type SequenceResult struct {
		Sequences []SequenceDesign `json:"sequences"`
		Summary   struct {
			TotalSequences int `json:"total_sequences"`
		} `json:"summary"`
	}

	type DataModelResult struct {
		Tables  []DataTable `json:"tables"`
		Enums   []DataEnum  `json:"enums"`
		Summary struct {
			TotalTables      int `json:"total_tables"`
			TotalIndexes     int `json:"total_indexes"`
			TotalForeignKeys int `json:"total_foreign_keys"`
			TotalEnums       int `json:"total_enums"`
		} `json:"summary"`
	}

	// Create parallel phases
	phases := []generator.Phase{
		{
			Name: "Tech Specs",
			Execute: func() (interface{}, error) {
				var result TechSpecsResult
				prompt := prompts.DeriveTechSpecs + "\n" + string(brContent)
				if err := client.CallJSON(prompt, &result); err != nil {
					return nil, err
				}
				return result, nil
			},
		},
		{
			Name: "Interface Contracts",
			Execute: func() (interface{}, error) {
				var result InterfaceContractsResult
				prompt := prompts.DeriveInterfaceContracts + "\n" + l1Input
				if err := client.CallJSON(prompt, &result); err != nil {
					return nil, err
				}
				return result, nil
			},
		},
		{
			Name: "Aggregate Design",
			Execute: func() (interface{}, error) {
				var result AggregateResult
				prompt := prompts.DeriveAggregateDesign + "\n" + aggSeqInput
				if err := client.CallJSON(prompt, &result); err != nil {
					return nil, err
				}
				return result, nil
			},
		},
		{
			Name: "Sequence Design",
			Execute: func() (interface{}, error) {
				var result SequenceResult
				prompt := prompts.DeriveSequenceDesign + "\n" + aggSeqInput
				if err := client.CallJSON(prompt, &result); err != nil {
					return nil, err
				}
				return result, nil
			},
		},
		{
			Name: "Data Model",
			Execute: func() (interface{}, error) {
				var result DataModelResult
				prompt := prompts.DeriveDataModel + "\n" + string(dmContent)
				if err := client.CallJSON(prompt, &result); err != nil {
					return nil, err
				}
				return result, nil
			},
		},
	}

	// Execute in parallel (max 3 concurrent to respect rate limits)
	executor := generator.NewParallelExecutor(3)
	phaseResults := executor.Execute(phases)

	// Extract results (with error handling)
	var tsResult TechSpecsResult
	var icResult InterfaceContractsResult
	var aggResult AggregateResult
	var seqResult SequenceResult
	var dataResult DataModelResult

	for _, pr := range phaseResults {
		if pr.Error != nil {
			return fmt.Errorf("failed to generate %s: %w", pr.Name, pr.Error)
		}
		switch pr.Name {
		case "Tech Specs":
			tsResult = pr.Data.(TechSpecsResult)
		case "Interface Contracts":
			icResult = pr.Data.(InterfaceContractsResult)
		case "Aggregate Design":
			aggResult = pr.Data.(AggregateResult)
		case "Sequence Design":
			seqResult = pr.Data.(SequenceResult)
		case "Data Model":
			dataResult = pr.Data.(DataModelResult)
		}
	}

	// Print generation summary
	fmt.Fprintf(os.Stderr, "\nGeneration Summary:\n")
	fmt.Fprintf(os.Stderr, "  Tech Specs:           %d specs\n", len(tsResult.TechSpecs))
	fmt.Fprintf(os.Stderr, "  Interface Contracts:  %d contracts, %d operations\n", len(icResult.InterfaceContracts), icResult.Summary.TotalOperations)
	fmt.Fprintf(os.Stderr, "  Aggregates:           %d aggregates, %d behaviors\n", len(aggResult.Aggregates), aggResult.Summary.TotalBehaviors)
	fmt.Fprintf(os.Stderr, "  Sequences:            %d sequences\n", len(seqResult.Sequences))
	fmt.Fprintf(os.Stderr, "  Data Tables:          %d tables, %d indexes\n", len(dataResult.Tables), dataResult.Summary.TotalIndexes)

	// Combine results
	var result L2Result
	result.TestCases = allTestCases
	result.TechSpecs = tsResult.TechSpecs
	result.Summary = L2Summary{
		TestCasesGenerated: len(allTestCases),
		TechSpecsGenerated: len(tsResult.TechSpecs),
		Coverage: L2Coverage{
			ACsCovered:     tcResult.Summary.Coverage.ACsCovered,
			BRsCovered:     tsResult.Summary.Total,
			HappyPathTests: tcResult.Summary.ByCategory.Positive,
			ErrorTests:     tcResult.Summary.ByCategory.Negative,
			EdgeCaseTests:  tcResult.Summary.ByCategory.Boundary,
		},
	}

	// Store TDAI summary for output (convert from generator type)
	tdaiSummary := convertTDAISummary(tcResult.Summary)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Track which files were written (for interactive mode skips)
	writtenFiles := make(map[string]bool)

	// Write Test Cases (TDAI format)
	fmt.Fprintln(os.Stderr, "\nPhase L2-W1: Writing Test Cases...")

	tcPath := filepath.Join(outputDir, "test-cases.md")
	if err := writeTestCases(tcPath, result.TestCases, tdaiSummary); err != nil {
		return fmt.Errorf("failed to write test cases: %w", err)
	}

	tcSummaryStr := fmt.Sprintf("(P:%d N:%d B:%d H:%d)", tdaiSummary.ByCategory.Positive, tdaiSummary.ByCategory.Negative, tdaiSummary.ByCategory.Boundary, tdaiSummary.ByCategory.Hallucination)
	tcWriteResult, _, err := workflow.HandleFileApproval(tcPath, "Test Cases (TDAI)", len(result.TestCases), "test cases", tcSummaryStr, interactive)
	if err != nil {
		return err
	}
	if tcWriteResult.Written {
		writtenFiles[tcPath] = true
	}

	// Write Tech Specs
	fmt.Fprintln(os.Stderr, "\nPhase L2-W2: Writing Tech Specs...")

	tsPath := filepath.Join(outputDir, "tech-specs.md")
	if err := writeTechSpecs(tsPath, result.TechSpecs); err != nil {
		return fmt.Errorf("failed to write tech specs: %w", err)
	}

	tsWriteResult, _, err := workflow.HandleFileApproval(tsPath, "Tech Specs", len(result.TechSpecs), "tech specs", "", interactive)
	if err != nil {
		return err
	}
	if tsWriteResult.Written {
		writtenFiles[tsPath] = true
	}

	// Write Interface Contracts
	fmt.Fprintln(os.Stderr, "\nPhase L2-W3: Writing Interface Contracts...")

	icPath := filepath.Join(outputDir, "interface-contracts.md")
	if err := writeInterfaceContracts(icPath, icResult.InterfaceContracts, icResult.SharedTypes); err != nil {
		return fmt.Errorf("failed to write interface contracts: %w", err)
	}

	icSummary := fmt.Sprintf("(%d operations)", icResult.Summary.TotalOperations)
	icWriteResult, _, err := workflow.HandleFileApproval(icPath, "Interface Contracts", len(icResult.InterfaceContracts), "contracts", icSummary, interactive)
	if err != nil {
		return err
	}
	if icWriteResult.Written {
		writtenFiles[icPath] = true
	}

	// Write Aggregate Design
	fmt.Fprintln(os.Stderr, "\nPhase L2-W4: Writing Aggregate Design...")

	aggPath := filepath.Join(outputDir, "aggregate-design.md")
	if err := writeAggregateDesign(aggPath, aggResult.Aggregates); err != nil {
		return fmt.Errorf("failed to write aggregate design: %w", err)
	}

	aggSummary := fmt.Sprintf("(%d behaviors)", aggResult.Summary.TotalBehaviors)
	aggWriteResult, _, err := workflow.HandleFileApproval(aggPath, "Aggregate Design", len(aggResult.Aggregates), "aggregates", aggSummary, interactive)
	if err != nil {
		return err
	}
	if aggWriteResult.Written {
		writtenFiles[aggPath] = true
	}

	// Write Sequence Design
	fmt.Fprintln(os.Stderr, "\nPhase L2-W5: Writing Sequence Design...")

	seqPath := filepath.Join(outputDir, "sequence-design.md")
	if err := writeSequenceDesign(seqPath, seqResult.Sequences); err != nil {
		return fmt.Errorf("failed to write sequence design: %w", err)
	}

	seqWriteResult, _, err := workflow.HandleFileApproval(seqPath, "Sequence Design", len(seqResult.Sequences), "sequences", "", interactive)
	if err != nil {
		return err
	}
	if seqWriteResult.Written {
		writtenFiles[seqPath] = true
	}

	// Write Data Model
	fmt.Fprintln(os.Stderr, "\nPhase L2-W6: Writing Data Model...")

	dataPath := filepath.Join(outputDir, "initial-data-model.md")
	if err := writeDataModel(dataPath, dataResult.Tables, dataResult.Enums); err != nil {
		return fmt.Errorf("failed to write data model: %w", err)
	}

	dataSummary := fmt.Sprintf("(%d indexes)", dataResult.Summary.TotalIndexes)
	dataWriteResult, _, err := workflow.HandleFileApproval(dataPath, "Initial Data Model", len(dataResult.Tables), "tables", dataSummary, interactive)
	if err != nil {
		return err
	}
	if dataWriteResult.Written {
		writtenFiles[dataPath] = true
	}

	// Write JSON for further processing
	jsonPath := filepath.Join(outputDir, "l2-output.json")
	l2Output := map[string]interface{}{
		"test_cases":          result.TestCases,
		"tech_specs":          result.TechSpecs,
		"interface_contracts": icResult.InterfaceContracts,
		"shared_types":        icResult.SharedTypes,
		"aggregates":          aggResult.Aggregates,
		"sequences":           seqResult.Sequences,
		"tables":              dataResult.Tables,
		"enums":               dataResult.Enums,
	}
	jsonContent, _ := json.MarshalIndent(l2Output, "", "  ")
	if err := os.WriteFile(jsonPath, jsonContent, 0644); err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", jsonPath)

	// Print summary
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "   L2 DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "   (Tactical Design Layer)")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Generated:")
	fmt.Fprintf(os.Stderr, "  Test Cases:           %d\n", len(result.TestCases))
	fmt.Fprintf(os.Stderr, "  Tech Specs:           %d\n", len(result.TechSpecs))
	fmt.Fprintf(os.Stderr, "  Interface Contracts:  %d (%d operations)\n", len(icResult.InterfaceContracts), icResult.Summary.TotalOperations)
	fmt.Fprintf(os.Stderr, "  Aggregates:           %d (%d behaviors)\n", len(aggResult.Aggregates), aggResult.Summary.TotalBehaviors)
	fmt.Fprintf(os.Stderr, "  Sequences:            %d\n", len(seqResult.Sequences))
	fmt.Fprintf(os.Stderr, "  Data Tables:          %d\n", len(dataResult.Tables))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Output Files:")
	fmt.Fprintf(os.Stderr, "  %s\n", tcPath)
	fmt.Fprintf(os.Stderr, "  %s\n", tsPath)
	fmt.Fprintf(os.Stderr, "  %s\n", icPath)
	fmt.Fprintf(os.Stderr, "  %s\n", aggPath)
	fmt.Fprintf(os.Stderr, "  %s\n", seqPath)
	fmt.Fprintf(os.Stderr, "  %s\n", dataPath)

	return nil
}

func writeTestCases(path string, testCases []TestCase, summary TDAISummary) error {
	timestamp := time.Now().Format(time.RFC3339)

	// Convert to formatter types
	fmtCases := convertTestCasesToFormatter(testCases)
	fmtSummary := convertSummaryToFormatter(summary)

	content := formatter.FormatTestCases(fmtCases, fmtSummary, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}

// convertTestCasesToFormatter converts local TestCase slice to formatter types
func convertTestCasesToFormatter(testCases []TestCase) []formatter.TestCase {
	result := make([]formatter.TestCase, len(testCases))
	for i, tc := range testCases {
		testData := make([]formatter.TestData, len(tc.TestData))
		for j, td := range tc.TestData {
			testData[j] = formatter.TestData{
				Field: td.Field,
				Value: td.Value,
				Notes: td.Notes,
			}
		}
		result[i] = formatter.TestCase{
			ID:              tc.ID,
			Name:            tc.Name,
			Category:        tc.Category,
			ACRef:           tc.ACRef,
			BRRefs:          tc.BRRefs,
			Preconditions:   tc.Preconditions,
			TestData:        testData,
			Steps:           tc.Steps,
			ExpectedResults: tc.ExpectedResults,
			ShouldNot:       tc.ShouldNot,
		}
	}
	return result
}

// convertSummaryToFormatter converts local TDAISummary to formatter type
func convertSummaryToFormatter(s TDAISummary) formatter.TDAISummary {
	var result formatter.TDAISummary
	result.Total = s.Total
	result.ByCategory.Positive = s.ByCategory.Positive
	result.ByCategory.Negative = s.ByCategory.Negative
	result.ByCategory.Boundary = s.ByCategory.Boundary
	result.ByCategory.Hallucination = s.ByCategory.Hallucination
	result.Coverage.ACsCovered = s.Coverage.ACsCovered
	result.Coverage.PositiveRatio = s.Coverage.PositiveRatio
	result.Coverage.NegativeRatio = s.Coverage.NegativeRatio
	result.Coverage.HasHallucinationTests = s.Coverage.HasHallucinationTests
	return result
}

func writeTechSpecs(path string, techSpecs []TechSpec) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmtSpecs := convertTechSpecsToFormatter(techSpecs)
	content := formatter.FormatTechSpecs(fmtSpecs, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}

// convertTechSpecsToFormatter converts local TechSpec slice to formatter types
func convertTechSpecsToFormatter(techSpecs []TechSpec) []formatter.TechSpec {
	result := make([]formatter.TechSpec, len(techSpecs))
	for i, ts := range techSpecs {
		dataReqs := make([]formatter.DataReq, len(ts.DataRequirements))
		for j, dr := range ts.DataRequirements {
			dataReqs[j] = formatter.DataReq{
				Field:       dr.Field,
				Type:        dr.Type,
				Constraints: dr.Constraints,
				Source:      dr.Source,
			}
		}
		errHandling := make([]formatter.ErrorHandling, len(ts.ErrorHandling))
		for j, eh := range ts.ErrorHandling {
			errHandling[j] = formatter.ErrorHandling{
				Condition:  eh.Condition,
				ErrorCode:  eh.ErrorCode,
				Message:    eh.Message,
				HTTPStatus: eh.HTTPStatus,
			}
		}
		result[i] = formatter.TechSpec{
			ID:               ts.ID,
			Name:             ts.Name,
			BRRef:            ts.BRRef,
			Rule:             ts.Rule,
			Implementation:   ts.Implementation,
			ValidationPoints: ts.ValidationPoints,
			DataRequirements: dataReqs,
			ErrorHandling:    errHandling,
			RelatedACs:       ts.RelatedACs,
		}
	}
	return result
}

func writeInterfaceContracts(path string, contracts []InterfaceContract, sharedTypes []SharedType) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmtContracts := convertContractsToFormatter(contracts)
	fmtSharedTypes := convertSharedTypesToFormatter(sharedTypes)
	content := formatter.FormatInterfaceContracts(fmtContracts, fmtSharedTypes, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}

func writeAggregateDesign(path string, aggregates []AggregateDesign) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmtAggs := convertAggregatesToFormatter(aggregates)
	content := formatter.FormatAggregateDesign(fmtAggs, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}

func writeSequenceDesign(path string, sequences []SequenceDesign) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmtSeqs := convertSequencesToFormatter(sequences)
	content := formatter.FormatSequenceDesign(fmtSeqs, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}

func writeDataModel(path string, tables []DataTable, enums []DataEnum) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmtTables := convertTablesToFormatter(tables)
	fmtEnums := convertEnumsToFormatter(enums)
	content := formatter.FormatDataModel(fmtTables, fmtEnums, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}
