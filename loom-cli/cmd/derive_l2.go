package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ikadar/loom-cli/internal/claude"
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

	// Create Claude client
	client := claude.NewClient()

	// Phase 1: Generate Test Cases from ACs (TDAI methodology)
	fmt.Fprintln(os.Stderr, "\nPhase L2-1: Generating TDAI Test Cases from Acceptance Criteria...")

	tcPrompt := prompts.DeriveTestCases + "\n" + string(acContent)

	var tcResult struct {
		TestSuites []TestSuite `json:"test_suites"`
		Summary    TDAISummary `json:"summary"`
	}
	if err := client.CallJSON(tcPrompt, &tcResult); err != nil {
		return fmt.Errorf("failed to generate test cases: %w", err)
	}

	// Flatten test suites into test cases for compatibility
	var allTestCases []TestCase
	for _, suite := range tcResult.TestSuites {
		allTestCases = append(allTestCases, suite.Tests...)
	}

	fmt.Fprintf(os.Stderr, "  Generated: %d Test Cases (P:%d N:%d B:%d H:%d)\n",
		tcResult.Summary.Total,
		tcResult.Summary.ByCategory.Positive,
		tcResult.Summary.ByCategory.Negative,
		tcResult.Summary.ByCategory.Boundary,
		tcResult.Summary.ByCategory.Hallucination)

	// Phase 2: Generate Tech Specs from BRs
	fmt.Fprintln(os.Stderr, "\nPhase L2-2: Generating Tech Specs from Business Rules...")

	tsPrompt := prompts.DeriveTechSpecs + "\n" + string(brContent)

	var tsResult struct {
		TechSpecs []TechSpec `json:"tech_specs"`
		Summary   struct {
			Total int `json:"total"`
		} `json:"summary"`
	}
	if err := client.CallJSON(tsPrompt, &tsResult); err != nil {
		return fmt.Errorf("failed to generate tech specs: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Tech Specs\n", len(tsResult.TechSpecs))

	// Phase 3: Generate Interface Contracts
	fmt.Fprintln(os.Stderr, "\nPhase L2-3: Generating Interface Contracts...")

	l1Input := string(dmContent) + "\n\n---\n\n" + string(brContent) + "\n\n---\n\n" + string(acContent)
	icPrompt := prompts.DeriveInterfaceContracts + "\n" + l1Input

	var icResult struct {
		InterfaceContracts []InterfaceContract `json:"interface_contracts"`
		SharedTypes        []SharedType        `json:"shared_types"`
		Summary            struct {
			TotalContracts   int `json:"total_contracts"`
			TotalOperations  int `json:"total_operations"`
			TotalEvents      int `json:"total_events"`
			TotalSharedTypes int `json:"total_shared_types"`
		} `json:"summary"`
	}
	if err := client.CallJSON(icPrompt, &icResult); err != nil {
		return fmt.Errorf("failed to generate interface contracts: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Interface Contracts, %d Operations\n",
		len(icResult.InterfaceContracts), icResult.Summary.TotalOperations)

	// Phase 4: Generate Aggregate Design
	fmt.Fprintln(os.Stderr, "\nPhase L2-4: Generating Aggregate Design...")

	aggInput := string(dmContent) + "\n\n---\n\n" + string(brContent)
	aggPrompt := prompts.DeriveAggregateDesign + "\n" + aggInput

	var aggResult struct {
		Aggregates []AggregateDesign `json:"aggregates"`
		Summary    struct {
			TotalAggregates  int `json:"total_aggregates"`
			TotalInvariants  int `json:"total_invariants"`
			TotalBehaviors   int `json:"total_behaviors"`
			TotalEvents      int `json:"total_events"`
		} `json:"summary"`
	}
	if err := client.CallJSON(aggPrompt, &aggResult); err != nil {
		return fmt.Errorf("failed to generate aggregate design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Aggregates, %d Behaviors\n",
		len(aggResult.Aggregates), aggResult.Summary.TotalBehaviors)

	// Phase 5: Generate Sequence Design
	fmt.Fprintln(os.Stderr, "\nPhase L2-5: Generating Sequence Design...")

	seqInput := string(dmContent) + "\n\n---\n\n" + string(brContent)
	seqPrompt := prompts.DeriveSequenceDesign + "\n" + seqInput

	var seqResult struct {
		Sequences []SequenceDesign `json:"sequences"`
		Summary   struct {
			TotalSequences int `json:"total_sequences"`
		} `json:"summary"`
	}
	if err := client.CallJSON(seqPrompt, &seqResult); err != nil {
		return fmt.Errorf("failed to generate sequence design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Sequences\n", len(seqResult.Sequences))

	// Phase 6: Generate Data Model
	fmt.Fprintln(os.Stderr, "\nPhase L2-6: Generating Initial Data Model...")

	dataInput := string(dmContent)
	dataPrompt := prompts.DeriveDataModel + "\n" + dataInput

	var dataResult struct {
		Tables  []DataTable `json:"tables"`
		Enums   []DataEnum  `json:"enums"`
		Summary struct {
			TotalTables      int `json:"total_tables"`
			TotalIndexes     int `json:"total_indexes"`
			TotalForeignKeys int `json:"total_foreign_keys"`
			TotalEnums       int `json:"total_enums"`
		} `json:"summary"`
	}
	if err := client.CallJSON(dataPrompt, &dataResult); err != nil {
		return fmt.Errorf("failed to generate data model: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Tables, %d Indexes\n",
		len(dataResult.Tables), dataResult.Summary.TotalIndexes)

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

	// Store TDAI summary for output
	tdaiSummary := tcResult.Summary

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write Test Cases (TDAI format)
	fmt.Fprintln(os.Stderr, "\nPhase L2-3: Writing output...")

	tcPath := filepath.Join(outputDir, "test-cases.md")
	if err := writeTestCases(tcPath, result.TestCases, tdaiSummary); err != nil {
		return fmt.Errorf("failed to write test cases: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", tcPath)

	// Write Tech Specs
	tsPath := filepath.Join(outputDir, "tech-specs.md")
	if err := writeTechSpecs(tsPath, result.TechSpecs); err != nil {
		return fmt.Errorf("failed to write tech specs: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", tsPath)

	// Write Interface Contracts
	icPath := filepath.Join(outputDir, "interface-contracts.md")
	if err := writeInterfaceContracts(icPath, icResult.InterfaceContracts, icResult.SharedTypes); err != nil {
		return fmt.Errorf("failed to write interface contracts: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", icPath)

	// Write Aggregate Design
	aggPath := filepath.Join(outputDir, "aggregate-design.md")
	if err := writeAggregateDesign(aggPath, aggResult.Aggregates); err != nil {
		return fmt.Errorf("failed to write aggregate design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", aggPath)

	// Write Sequence Design
	seqPath := filepath.Join(outputDir, "sequence-design.md")
	if err := writeSequenceDesign(seqPath, seqResult.Sequences); err != nil {
		return fmt.Errorf("failed to write sequence design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", seqPath)

	// Write Data Model
	dataPath := filepath.Join(outputDir, "initial-data-model.md")
	if err := writeDataModel(dataPath, dataResult.Tables, dataResult.Enums); err != nil {
		return fmt.Errorf("failed to write data model: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", dataPath)

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
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# TDAI Test Cases\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "**Methodology:** Test-Driven AI Development (TDAI)\n\n")

	// TDAI Summary
	fmt.Fprintf(f, "## Summary\n\n")
	fmt.Fprintf(f, "| Category | Count | Ratio |\n")
	fmt.Fprintf(f, "|----------|-------|-------|\n")
	fmt.Fprintf(f, "| Positive | %d | %.1f%% |\n", summary.ByCategory.Positive, summary.Coverage.PositiveRatio*100)
	fmt.Fprintf(f, "| Negative | %d | %.1f%% |\n", summary.ByCategory.Negative, summary.Coverage.NegativeRatio*100)
	fmt.Fprintf(f, "| Boundary | %d | - |\n", summary.ByCategory.Boundary)
	fmt.Fprintf(f, "| Hallucination Prevention | %d | - |\n", summary.ByCategory.Hallucination)
	fmt.Fprintf(f, "| **Total** | **%d** | - |\n\n", summary.Total)

	fmt.Fprintf(f, "**Coverage:** %d ACs covered\n", summary.Coverage.ACsCovered)
	if summary.Coverage.HasHallucinationTests {
		fmt.Fprintf(f, "**Hallucination Prevention:** ✓ Enabled\n")
	}
	fmt.Fprintf(f, "\n---\n\n")

	// Group tests by category
	categories := []string{"positive", "negative", "boundary", "hallucination"}
	categoryNames := map[string]string{
		"positive":      "Positive Tests (Happy Path)",
		"negative":      "Negative Tests (Error Cases)",
		"boundary":      "Boundary Tests",
		"hallucination": "Hallucination Prevention Tests",
	}

	for _, cat := range categories {
		var catTests []TestCase
		for _, tc := range testCases {
			if tc.Category == cat {
				catTests = append(catTests, tc)
			}
		}

		if len(catTests) == 0 {
			continue
		}

		fmt.Fprintf(f, "## %s\n\n", categoryNames[cat])

		for _, tc := range catTests {
			// Add anchor: ### TC-AC-CUST-001-P01 – Title {#tc-ac-cust-001-p01}
			fmt.Fprintf(f, "### %s – %s {#%s}\n\n", tc.ID, tc.Name, toAnchorL2(tc.ID))

			// Special handling for hallucination tests
			if tc.Category == "hallucination" && tc.ShouldNot != "" {
				fmt.Fprintf(f, "**⚠️ Should NOT:** %s\n\n", tc.ShouldNot)
			}

			fmt.Fprintf(f, "**Preconditions:**\n")
			for _, p := range tc.Preconditions {
				fmt.Fprintf(f, "- %s\n", p)
			}
			fmt.Fprintf(f, "\n")

			if len(tc.TestData) > 0 {
				fmt.Fprintf(f, "**Test Data:**\n")
				fmt.Fprintf(f, "| Field | Value | Notes |\n")
				fmt.Fprintf(f, "|-------|-------|-------|\n")
				for _, td := range tc.TestData {
					fmt.Fprintf(f, "| %s | %v | %s |\n", td.Field, td.Value, td.Notes)
				}
				fmt.Fprintf(f, "\n")
			}

			fmt.Fprintf(f, "**Steps:**\n")
			for i, s := range tc.Steps {
				fmt.Fprintf(f, "%d. %s\n", i+1, s)
			}
			fmt.Fprintf(f, "\n")

			fmt.Fprintf(f, "**Expected Result:**\n")
			for _, r := range tc.ExpectedResults {
				fmt.Fprintf(f, "- %s\n", r)
			}
			fmt.Fprintf(f, "\n")

			fmt.Fprintf(f, "**Traceability:**\n")
			// Link to L1 acceptance-criteria.md
			fmt.Fprintf(f, "- AC: %s\n", toLinkL2(tc.ACRef, "../l1/acceptance-criteria.md"))
			if len(tc.BRRefs) > 0 {
				fmt.Fprintf(f, "- BR: ")
				for i, br := range tc.BRRefs {
					if i > 0 {
						fmt.Fprintf(f, ", ")
					}
					fmt.Fprintf(f, "%s", toLinkL2(br, "../l1/business-rules.md"))
				}
				fmt.Fprintf(f, "\n")
			}
			fmt.Fprintf(f, "\n---\n\n")
		}
	}

	return nil
}

func writeTechSpecs(path string, techSpecs []TechSpec) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Technical Specifications\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, ts := range techSpecs {
		// Add anchor: ## TS-BR-CUST-001 – Title {#ts-br-cust-001}
		fmt.Fprintf(f, "## %s – %s {#%s}\n\n", ts.ID, ts.Name, toAnchorL2(ts.ID))
		fmt.Fprintf(f, "**Rule:** %s\n\n", ts.Rule)
		fmt.Fprintf(f, "**Implementation Approach:**\n%s\n\n", ts.Implementation)

		if len(ts.ValidationPoints) > 0 {
			fmt.Fprintf(f, "**Validation Points:**\n")
			for _, vp := range ts.ValidationPoints {
				fmt.Fprintf(f, "- %s\n", vp)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(ts.DataRequirements) > 0 {
			fmt.Fprintf(f, "**Data Requirements:**\n")
			fmt.Fprintf(f, "| Field | Type | Constraints | Source |\n")
			fmt.Fprintf(f, "|-------|------|-------------|--------|\n")
			for _, dr := range ts.DataRequirements {
				fmt.Fprintf(f, "| %s | %s | %s | %s |\n", dr.Field, dr.Type, dr.Constraints, dr.Source)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(ts.ErrorHandling) > 0 {
			fmt.Fprintf(f, "**Error Handling:**\n")
			fmt.Fprintf(f, "| Condition | Error Code | Message | HTTP Status |\n")
			fmt.Fprintf(f, "|-----------|------------|---------|-------------|\n")
			for _, eh := range ts.ErrorHandling {
				fmt.Fprintf(f, "| %s | %s | %s | %d |\n", eh.Condition, eh.ErrorCode, eh.Message, eh.HTTPStatus)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "**Traceability:**\n")
		// Link to L1 business-rules.md
		fmt.Fprintf(f, "- BR: %s\n", toLinkL2(ts.BRRef, "../l1/business-rules.md"))
		if len(ts.RelatedACs) > 0 {
			fmt.Fprintf(f, "- Related ACs: ")
			for i, ac := range ts.RelatedACs {
				if i > 0 {
					fmt.Fprintf(f, ", ")
				}
				fmt.Fprintf(f, "%s", toLinkL2(ac, "../l1/acceptance-criteria.md"))
			}
			fmt.Fprintf(f, "\n")
		}
		fmt.Fprintf(f, "\n---\n\n")
	}

	return nil
}

func writeInterfaceContracts(path string, contracts []InterfaceContract, sharedTypes []SharedType) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Interface Contracts\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	// Shared Types section
	if len(sharedTypes) > 0 {
		fmt.Fprintf(f, "## Shared Types\n\n")
		for _, st := range sharedTypes {
			fmt.Fprintf(f, "### %s\n\n", st.Name)
			fmt.Fprintf(f, "| Field | Type | Constraints |\n")
			fmt.Fprintf(f, "|-------|------|-------------|\n")
			for _, field := range st.Fields {
				fmt.Fprintf(f, "| %s | %s | %s |\n", field.Name, field.Type, field.Constraints)
			}
			fmt.Fprintf(f, "\n")
		}
		fmt.Fprintf(f, "---\n\n")
	}

	// Contracts
	for _, ic := range contracts {
		fmt.Fprintf(f, "## %s – %s {#%s}\n\n", ic.ID, ic.ServiceName, toAnchorL2(ic.ID))
		fmt.Fprintf(f, "**Purpose:** %s\n\n", ic.Purpose)
		fmt.Fprintf(f, "**Base URL:** `%s`\n\n", ic.BaseURL)

		if ic.SecurityRequirements.Authentication != "" {
			fmt.Fprintf(f, "**Security:**\n")
			fmt.Fprintf(f, "- Authentication: %s\n", ic.SecurityRequirements.Authentication)
			fmt.Fprintf(f, "- Authorization: %s\n\n", ic.SecurityRequirements.Authorization)
		}

		// Operations
		fmt.Fprintf(f, "### Operations\n\n")
		for _, op := range ic.Operations {
			fmt.Fprintf(f, "#### %s `%s %s`\n\n", op.Name, op.Method, op.Path)
			fmt.Fprintf(f, "%s\n\n", op.Description)

			if len(op.InputSchema) > 0 {
				fmt.Fprintf(f, "**Input:**\n")
				for name, field := range op.InputSchema {
					req := ""
					if field.Required {
						req = " (required)"
					}
					fmt.Fprintf(f, "- `%s`: %s%s\n", name, field.Type, req)
				}
				fmt.Fprintf(f, "\n")
			}

			if len(op.OutputSchema) > 0 {
				fmt.Fprintf(f, "**Output:**\n")
				for name, field := range op.OutputSchema {
					fmt.Fprintf(f, "- `%s`: %s\n", name, field.Type)
				}
				fmt.Fprintf(f, "\n")
			}

			if len(op.Errors) > 0 {
				fmt.Fprintf(f, "**Errors:**\n")
				fmt.Fprintf(f, "| Code | HTTP | Message |\n")
				fmt.Fprintf(f, "|------|------|----------|\n")
				for _, e := range op.Errors {
					fmt.Fprintf(f, "| %s | %d | %s |\n", e.Code, e.HTTPStatus, e.Message)
				}
				fmt.Fprintf(f, "\n")
			}

			if len(op.Preconditions) > 0 {
				fmt.Fprintf(f, "**Preconditions:** %v\n\n", op.Preconditions)
			}
			if len(op.Postconditions) > 0 {
				fmt.Fprintf(f, "**Postconditions:** %v\n\n", op.Postconditions)
			}
		}

		// Events
		if len(ic.Events) > 0 {
			fmt.Fprintf(f, "### Events\n\n")
			for _, ev := range ic.Events {
				fmt.Fprintf(f, "- **%s**: %s (payload: %v)\n", ev.Name, ev.Description, ev.Payload)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeAggregateDesign(path string, aggregates []AggregateDesign) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Aggregate Design\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, agg := range aggregates {
		fmt.Fprintf(f, "## %s – %s {#%s}\n\n", agg.ID, agg.Name, toAnchorL2(agg.ID))
		fmt.Fprintf(f, "**Purpose:** %s\n\n", agg.Purpose)

		// Root
		fmt.Fprintf(f, "### Aggregate Root: %s\n\n", agg.Root.Entity)
		fmt.Fprintf(f, "**Identity:** %s\n\n", agg.Root.Identity)
		if len(agg.Root.Attributes) > 0 {
			fmt.Fprintf(f, "**Attributes:**\n")
			fmt.Fprintf(f, "| Name | Type | Mutable |\n")
			fmt.Fprintf(f, "|------|------|----------|\n")
			for _, attr := range agg.Root.Attributes {
				fmt.Fprintf(f, "| %s | %s | %v |\n", attr.Name, attr.Type, attr.Mutable)
			}
			fmt.Fprintf(f, "\n")
		}

		// Invariants
		if len(agg.Invariants) > 0 {
			fmt.Fprintf(f, "### Invariants\n\n")
			for _, inv := range agg.Invariants {
				fmt.Fprintf(f, "- **%s**: %s\n  - Enforcement: %s\n", inv.ID, inv.Rule, inv.Enforcement)
			}
			fmt.Fprintf(f, "\n")
		}

		// Child Entities
		if len(agg.Entities) > 0 {
			fmt.Fprintf(f, "### Child Entities\n\n")
			for _, ent := range agg.Entities {
				fmt.Fprintf(f, "#### %s\n\n", ent.Name)
				fmt.Fprintf(f, "**Identity:** %s\n\n", ent.Identity)
				fmt.Fprintf(f, "**Purpose:** %s\n\n", ent.Purpose)
			}
		}

		// Value Objects
		if len(agg.ValueObjects) > 0 {
			fmt.Fprintf(f, "### Value Objects\n\n")
			fmt.Fprintf(f, "%v\n\n", agg.ValueObjects)
		}

		// Behaviors
		if len(agg.Behaviors) > 0 {
			fmt.Fprintf(f, "### Behaviors\n\n")
			fmt.Fprintf(f, "| Command | Pre | Post | Emits |\n")
			fmt.Fprintf(f, "|---------|-----|------|-------|\n")
			for _, b := range agg.Behaviors {
				fmt.Fprintf(f, "| %s | %v | %v | %s |\n", b.Command, b.Preconditions, b.Postconditions, b.Emits)
			}
			fmt.Fprintf(f, "\n")
		}

		// Events
		if len(agg.Events) > 0 {
			fmt.Fprintf(f, "### Events\n\n")
			for _, ev := range agg.Events {
				fmt.Fprintf(f, "- **%s**: %v\n", ev.Name, ev.Payload)
			}
			fmt.Fprintf(f, "\n")
		}

		// Repository
		fmt.Fprintf(f, "### Repository: %s\n\n", agg.Repository.Name)
		fmt.Fprintf(f, "- Load Strategy: %s\n", agg.Repository.LoadStrategy)
		fmt.Fprintf(f, "- Concurrency: %s\n\n", agg.Repository.Concurrency)

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeSequenceDesign(path string, sequences []SequenceDesign) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Sequence Design\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	for _, seq := range sequences {
		fmt.Fprintf(f, "## %s – %s {#%s}\n\n", seq.ID, seq.Name, toAnchorL2(seq.ID))
		fmt.Fprintf(f, "%s\n\n", seq.Description)

		// Trigger
		fmt.Fprintf(f, "### Trigger\n\n")
		fmt.Fprintf(f, "**Type:** %s\n\n", seq.Trigger.Type)
		fmt.Fprintf(f, "%s\n\n", seq.Trigger.Description)

		// Participants
		if len(seq.Participants) > 0 {
			fmt.Fprintf(f, "### Participants\n\n")
			for _, p := range seq.Participants {
				fmt.Fprintf(f, "- **%s** (%s)\n", p.Name, p.Type)
			}
			fmt.Fprintf(f, "\n")
		}

		// Steps
		fmt.Fprintf(f, "### Sequence\n\n")
		for _, step := range seq.Steps {
			fmt.Fprintf(f, "%d. **%s** → %s: %s\n", step.Step, step.Actor, step.Target, step.Action)
			if step.Event != "" {
				fmt.Fprintf(f, "   - Emits: `%s`\n", step.Event)
			}
			if step.Returns != "" {
				fmt.Fprintf(f, "   - Returns: %s\n", step.Returns)
			}
		}
		fmt.Fprintf(f, "\n")

		// Mermaid Diagram
		fmt.Fprintf(f, "### Sequence Diagram\n\n")
		fmt.Fprintf(f, "```mermaid\nsequenceDiagram\n")
		for _, p := range seq.Participants {
			fmt.Fprintf(f, "    participant %s\n", p.Name)
		}
		for _, step := range seq.Steps {
			fmt.Fprintf(f, "    %s->>%s: %s\n", step.Actor, step.Target, step.Action)
			if step.Returns != "" {
				fmt.Fprintf(f, "    %s-->>%s: %s\n", step.Target, step.Actor, step.Returns)
			}
		}
		fmt.Fprintf(f, "```\n\n")

		// Outcome
		fmt.Fprintf(f, "### Outcome\n\n")
		fmt.Fprintf(f, "%s\n\n", seq.Outcome.Success)
		if len(seq.Outcome.StateChanges) > 0 {
			fmt.Fprintf(f, "**State Changes:**\n")
			for _, sc := range seq.Outcome.StateChanges {
				fmt.Fprintf(f, "- %s\n", sc)
			}
			fmt.Fprintf(f, "\n")
		}

		// Exceptions
		if len(seq.Exceptions) > 0 {
			fmt.Fprintf(f, "### Exceptions\n\n")
			for _, ex := range seq.Exceptions {
				fmt.Fprintf(f, "- **%s** (step %d): %s\n", ex.Condition, ex.Step, ex.Handling)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeDataModel(path string, tables []DataTable, enums []DataEnum) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Initial Data Model\n\n")
	fmt.Fprintf(f, "Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "---\n\n")

	// Enums
	if len(enums) > 0 {
		fmt.Fprintf(f, "## Enumerations\n\n")
		for _, enum := range enums {
			fmt.Fprintf(f, "### %s\n\n", enum.Name)
			fmt.Fprintf(f, "Values: `%v`\n\n", enum.Values)
		}
		fmt.Fprintf(f, "---\n\n")
	}

	// Tables
	fmt.Fprintf(f, "## Tables\n\n")
	for _, tbl := range tables {
		fmt.Fprintf(f, "### %s – %s {#%s}\n\n", tbl.ID, tbl.Name, toAnchorL2(tbl.ID))
		fmt.Fprintf(f, "**Aggregate:** %s\n\n", tbl.Aggregate)
		fmt.Fprintf(f, "**Purpose:** %s\n\n", tbl.Purpose)

		// Fields
		fmt.Fprintf(f, "**Columns:**\n\n")
		fmt.Fprintf(f, "| Name | Type | Constraints | Default |\n")
		fmt.Fprintf(f, "|------|------|-------------|----------|\n")
		for _, field := range tbl.Fields {
			fmt.Fprintf(f, "| %s | %s | %s | %s |\n", field.Name, field.Type, field.Constraints, field.Default)
		}
		fmt.Fprintf(f, "\n")

		// Primary Key
		fmt.Fprintf(f, "**Primary Key:** %v\n\n", tbl.PrimaryKey.Columns)

		// Indexes
		if len(tbl.Indexes) > 0 {
			fmt.Fprintf(f, "**Indexes:**\n")
			for _, idx := range tbl.Indexes {
				fmt.Fprintf(f, "- `%s`: %v\n", idx.Name, idx.Columns)
			}
			fmt.Fprintf(f, "\n")
		}

		// Foreign Keys
		if len(tbl.ForeignKeys) > 0 {
			fmt.Fprintf(f, "**Foreign Keys:**\n")
			for _, fk := range tbl.ForeignKeys {
				fmt.Fprintf(f, "- %v → %s (ON DELETE %s)\n", fk.Columns, fk.References, fk.OnDelete)
			}
			fmt.Fprintf(f, "\n")
		}

		// Check Constraints
		if len(tbl.CheckConstraints) > 0 {
			fmt.Fprintf(f, "**Constraints:**\n")
			for _, c := range tbl.CheckConstraints {
				fmt.Fprintf(f, "- `%s`: %s\n", c.Name, c.Expression)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "---\n\n")
	}

	// ER Diagram
	fmt.Fprintf(f, "## Entity-Relationship Diagram\n\n")
	fmt.Fprintf(f, "```mermaid\nerDiagram\n")
	for _, tbl := range tables {
		fmt.Fprintf(f, "    %s {\n", tbl.Name)
		for _, field := range tbl.Fields {
			fmt.Fprintf(f, "        %s %s\n", field.Type, field.Name)
		}
		fmt.Fprintf(f, "    }\n")
	}
	// Add relationships
	for _, tbl := range tables {
		for _, fk := range tbl.ForeignKeys {
			// Extract target table from references like "customers(id)"
			target := fk.References
			if idx := strings.Index(target, "("); idx > 0 {
				target = target[:idx]
			}
			fmt.Fprintf(f, "    %s ||--o{ %s : \"FK\"\n", target, tbl.Name)
		}
	}
	fmt.Fprintf(f, "```\n")

	return nil
}
