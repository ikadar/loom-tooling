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
	"github.com/ikadar/loom-cli/prompts"
)

// L3Result is the output of the derive-l3 command
type L3Result struct {
	APISpec                 APISpec                  `json:"api_spec"`
	ImplementationSkeletons []ImplementationSkeleton `json:"implementation_skeletons"`
	Summary                 L3Summary                `json:"summary"`
}

type L3Summary struct {
	EndpointsCount int `json:"endpoints_count"`
	ServicesCount  int `json:"services_count"`
	FunctionsCount int `json:"functions_count"`
}

type APISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    map[string]interface{} `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
}

type ImplementationSkeleton struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Functions    []FunctionSpec `json:"functions"`
	Dependencies []string       `json:"dependencies"`
	RelatedSpecs []string       `json:"related_specs"`
}

type FunctionSpec struct {
	Name        string   `json:"name"`
	Signature   string   `json:"signature"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	ErrorCases  []string `json:"error_cases"`
}

// Feature Ticket types
type FeatureTicket struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	Status                string   `json:"status"`
	BusinessGoal          string   `json:"business_goal"`
	UserStory             string   `json:"user_story"`
	AcceptanceCriteriaRefs []string `json:"acceptance_criteria_refs"`
	NFR                   []string `json:"nfr"`
	Dependencies          []string `json:"dependencies"`
	ImpactAreas           []string `json:"impact_areas"`
	OutOfScope            []string `json:"out_of_scope"`
	Priority              string   `json:"priority"`
	EstimatedComplexity   string   `json:"estimated_complexity"`
}

// Service Boundary types
type ServiceBoundary struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Purpose          string              `json:"purpose"`
	Capabilities     []string            `json:"capabilities"`
	Inputs           []ServiceIO         `json:"inputs"`
	Outputs          []ServiceIO         `json:"outputs"`
	OwnedAggregates  []string            `json:"owned_aggregates"`
	Dependencies     []ServiceDependency `json:"dependencies"`
	APIBase          string              `json:"api_base"`
	SeparationReason string              `json:"separation_reason"`
}

type ServiceIO struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ServiceDependency struct {
	Service string `json:"service"`
	Type    string `json:"type"`
	Reason  string `json:"reason"`
}

// Event Design types
type DomainEvent struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Purpose             string         `json:"purpose"`
	Trigger             string         `json:"trigger"`
	Aggregate           string         `json:"aggregate"`
	Payload             []EventField   `json:"payload"`
	InvariantsReflected []string       `json:"invariants_reflected"`
	Consumers           []string       `json:"consumers"`
	Version             string         `json:"version"`
}

type EventField struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type Command struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	Intent            string       `json:"intent"`
	RequiredData      []EventField `json:"required_data"`
	ExpectedOutcome   string       `json:"expected_outcome"`
	FailureConditions []string     `json:"failure_conditions"`
}

type IntegrationEvent struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Purpose   string   `json:"purpose"`
	Source    string   `json:"source"`
	Consumers []string `json:"consumers"`
	Payload   []string `json:"payload"`
}

// Dependency Graph types
type GraphComponent struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type GraphDependency struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func runDeriveL3() error {
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
		return fmt.Errorf("--input-dir is required (directory containing L2 documents)")
	}
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required")
	}

	// Read L1 and L2 documents
	fmt.Fprintln(os.Stderr, "Phase L3-0: Reading L1/L2 documents...")

	// L2 docs
	tsContent, err := os.ReadFile(filepath.Join(inputDir, "tech-specs.md"))
	if err != nil {
		return fmt.Errorf("failed to read tech-specs.md: %w", err)
	}

	// L1 docs (required for test case generation)
	acContent, err := os.ReadFile(filepath.Join(inputDir, "acceptance-criteria.md"))
	if err != nil {
		// Try parent directory (L1)
		acContent, err = os.ReadFile(filepath.Join(inputDir, "../acceptance-criteria.md"))
		if err != nil {
			return fmt.Errorf("failed to read acceptance-criteria.md: %w", err)
		}
	}

	aggContent, err := os.ReadFile(filepath.Join(inputDir, "aggregate-design.md"))
	if err != nil {
		aggContent = []byte{}
	}

	seqContent, err := os.ReadFile(filepath.Join(inputDir, "sequence-design.md"))
	if err != nil {
		seqContent = []byte{}
	}

	bcContent, err := os.ReadFile(filepath.Join(inputDir, "bounded-context-map.md"))
	if err != nil {
		// Try parent (L1 dir)
		bcContent, _ = os.ReadFile(filepath.Join(inputDir, "../bounded-context-map.md"))
	}

	dmContent, err := os.ReadFile(filepath.Join(inputDir, "domain-model.md"))
	if err != nil {
		// Try parent (L1 dir)
		dmContent, _ = os.ReadFile(filepath.Join(inputDir, "../domain-model.md"))
	}

	fmt.Fprintf(os.Stderr, "  Read: tech-specs.md (%d bytes)\n", len(tsContent))
	fmt.Fprintf(os.Stderr, "  Read: acceptance-criteria.md (%d bytes)\n", len(acContent))
	fmt.Fprintf(os.Stderr, "  Read: aggregate-design.md (%d bytes)\n", len(aggContent))
	fmt.Fprintf(os.Stderr, "  Read: sequence-design.md (%d bytes)\n", len(seqContent))
	fmt.Fprintf(os.Stderr, "  Read: bounded-context-map.md (%d bytes)\n", len(bcContent))
	fmt.Fprintf(os.Stderr, "  Read: domain-model.md (%d bytes)\n", len(dmContent))

	// Create Claude client
	client := claude.NewClient()

	// Phase L3-1: Generate Test Cases from Acceptance Criteria (TDAI)
	fmt.Fprintln(os.Stderr, "\nPhase L3-1: Generating TDAI Test Cases from Acceptance Criteria...")

	tcGenerator := generator.NewChunkedTestCaseGenerator(client)
	tcResult, err := tcGenerator.Generate(string(acContent))
	if err != nil {
		return fmt.Errorf("failed to generate test cases: %w", err)
	}

	// Flatten test suites for output
	allTestCases := generator.FlattenTestCases(tcResult.TestSuites)
	fmt.Fprintf(os.Stderr, "  Generated: %d Test Cases (P:%d N:%d B:%d H:%d)\n",
		tcResult.Summary.Total,
		tcResult.Summary.ByCategory.Positive,
		tcResult.Summary.ByCategory.Negative,
		tcResult.Summary.ByCategory.Boundary,
		tcResult.Summary.ByCategory.Hallucination)

	// Phase L3-2a: Generate API Spec
	fmt.Fprintln(os.Stderr, "\nPhase L3-2a: Generating API Specification...")

	apiPrompt := prompts.DeriveL3API + "\n\n" + string(tsContent)

	var apiResult APISpec
	if err := client.CallJSON(apiPrompt, &apiResult); err != nil {
		return fmt.Errorf("failed to generate API spec: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Generated: %d endpoints\n", len(apiResult.Paths))

	// Phase L3-2b: Generate Implementation Skeletons
	fmt.Fprintln(os.Stderr, "\nPhase L3-2b: Generating Implementation Skeletons...")

	skelPrompt := prompts.DeriveL3Skeletons + "\n\n" + string(tsContent)

	var skelResult struct {
		ImplementationSkeletons []ImplementationSkeleton `json:"implementation_skeletons"`
		Summary                 struct {
			ServicesCount  int `json:"services_count"`
			FunctionsCount int `json:"functions_count"`
		} `json:"summary"`
	}
	if err := client.CallJSON(skelPrompt, &skelResult); err != nil {
		return fmt.Errorf("failed to generate implementation skeletons: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  Generated: %d services\n", len(skelResult.ImplementationSkeletons))

	// Combine into L3Result
	var result L3Result
	result.APISpec = apiResult
	result.ImplementationSkeletons = skelResult.ImplementationSkeletons
	result.Summary.EndpointsCount = len(apiResult.Paths)
	result.Summary.ServicesCount = len(skelResult.ImplementationSkeletons)

	// Phase 3: Generate Feature Tickets
	fmt.Fprintln(os.Stderr, "\nPhase L3-3: Generating Feature Definition Tickets...")

	// Use AC + tech specs for feature tickets (test cases are now generated in this phase)
	ftInput := string(acContent) + "\n\n---\n\n" + string(tsContent)
	ftPrompt := prompts.DeriveFeatureTickets + ftInput

	var ftResult struct {
		FeatureTickets []FeatureTicket `json:"feature_tickets"`
		Summary        struct {
			TotalTickets int            `json:"total_tickets"`
			ByPriority   map[string]int `json:"by_priority"`
		} `json:"summary"`
	}
	if err := client.CallJSON(ftPrompt, &ftResult); err != nil {
		return fmt.Errorf("failed to generate feature tickets: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Feature Tickets\n", len(ftResult.FeatureTickets))

	// Phase 4: Generate Service Boundaries
	fmt.Fprintln(os.Stderr, "\nPhase L3-4: Generating Service Boundaries...")

	sbInput := string(bcContent) + "\n\n---\n\n" + string(aggContent)
	sbPrompt := prompts.DeriveServiceBoundaries + sbInput

	var sbResult struct {
		Services []ServiceBoundary `json:"services"`
		Summary  struct {
			TotalServices     int `json:"total_services"`
			TotalDependencies int `json:"total_dependencies"`
		} `json:"summary"`
	}
	if err := client.CallJSON(sbPrompt, &sbResult); err != nil {
		return fmt.Errorf("failed to generate service boundaries: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Service Boundaries\n", len(sbResult.Services))

	// Phase 5: Generate Event Design
	fmt.Fprintln(os.Stderr, "\nPhase L3-5: Generating Event & Message Design...")

	evInput := string(dmContent) + "\n\n---\n\n" + string(seqContent)
	evPrompt := prompts.DeriveEventDesign + evInput

	var evResult struct {
		DomainEvents      []DomainEvent      `json:"domain_events"`
		Commands          []Command          `json:"commands"`
		IntegrationEvents []IntegrationEvent `json:"integration_events"`
		Summary           struct {
			DomainEvents      int `json:"domain_events"`
			Commands          int `json:"commands"`
			IntegrationEvents int `json:"integration_events"`
		} `json:"summary"`
	}
	if err := client.CallJSON(evPrompt, &evResult); err != nil {
		return fmt.Errorf("failed to generate event design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Events, %d Commands\n",
		len(evResult.DomainEvents), len(evResult.Commands))

	// Phase 6: Generate Dependency Graph
	fmt.Fprintln(os.Stderr, "\nPhase L3-6: Generating Dependency Graph...")

	dgInput := ""
	for _, svc := range sbResult.Services {
		dgInput += fmt.Sprintf("Service: %s\nPurpose: %s\nDependencies: ", svc.Name, svc.Purpose)
		for _, dep := range svc.Dependencies {
			dgInput += fmt.Sprintf("%s (%s), ", dep.Service, dep.Type)
		}
		dgInput += "\n\n"
	}
	dgPrompt := prompts.DeriveDependencyGraph + dgInput

	var dgResult struct {
		Components   []GraphComponent   `json:"components"`
		Dependencies []GraphDependency  `json:"dependencies"`
		Summary      struct {
			TotalComponents   int            `json:"total_components"`
			TotalDependencies int            `json:"total_dependencies"`
			ByType            map[string]int `json:"by_type"`
		} `json:"summary"`
	}
	if err := client.CallJSON(dgPrompt, &dgResult); err != nil {
		return fmt.Errorf("failed to generate dependency graph: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Generated: %d Components, %d Dependencies\n",
		len(dgResult.Components), len(dgResult.Dependencies))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write output files
	fmt.Fprintln(os.Stderr, "\nPhase L3-W: Writing output...")

	// Write Test Cases (TDAI format)
	tcPath := filepath.Join(outputDir, "test-cases.md")
	if err := writeL3TestCases(tcPath, allTestCases, tcResult.Summary); err != nil {
		return fmt.Errorf("failed to write test cases: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", tcPath)

	// Write API Spec (OpenAPI JSON)

	apiPath := filepath.Join(outputDir, "openapi.json")
	apiContent, _ := json.MarshalIndent(result.APISpec, "", "  ")
	if err := os.WriteFile(apiPath, apiContent, 0644); err != nil {
		return fmt.Errorf("failed to write API spec: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", apiPath)

	// Write Implementation Skeletons
	implPath := filepath.Join(outputDir, "implementation-skeletons.md")
	if err := writeSkeletons(implPath, result.ImplementationSkeletons); err != nil {
		return fmt.Errorf("failed to write skeletons: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", implPath)

	// Write full JSON for further processing
	jsonPath := filepath.Join(outputDir, "l3-output.json")
	jsonContent, _ := json.MarshalIndent(result, "", "  ")
	if err := os.WriteFile(jsonPath, jsonContent, 0644); err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", jsonPath)

	// Write Feature Tickets
	ftPath := filepath.Join(outputDir, "feature-tickets.md")
	if err := writeFeatureTickets(ftPath, ftResult.FeatureTickets); err != nil {
		return fmt.Errorf("failed to write feature tickets: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", ftPath)

	// Write Service Boundaries
	sbPath := filepath.Join(outputDir, "service-boundaries.md")
	if err := writeServiceBoundaries(sbPath, sbResult.Services); err != nil {
		return fmt.Errorf("failed to write service boundaries: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", sbPath)

	// Write Event & Message Design
	evPath := filepath.Join(outputDir, "event-message-design.md")
	if err := writeEventDesign(evPath, evResult.DomainEvents, evResult.Commands, evResult.IntegrationEvents); err != nil {
		return fmt.Errorf("failed to write event design: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", evPath)

	// Write Dependency Graph
	dgPath := filepath.Join(outputDir, "dependency-graph.md")
	if err := writeDependencyGraph(dgPath, dgResult.Components, dgResult.Dependencies); err != nil {
		return fmt.Errorf("failed to write dependency graph: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Written: %s\n", dgPath)

	// Print summary
	fmt.Fprintln(os.Stderr, "\n========================================")
	fmt.Fprintln(os.Stderr, "   L3 DERIVATION COMPLETE")
	fmt.Fprintln(os.Stderr, "   (Operational Design Layer)")
	fmt.Fprintln(os.Stderr, "========================================")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Generated:")
	fmt.Fprintf(os.Stderr, "  Test Cases:          %d (P:%d N:%d B:%d H:%d)\n",
		tcResult.Summary.Total,
		tcResult.Summary.ByCategory.Positive,
		tcResult.Summary.ByCategory.Negative,
		tcResult.Summary.ByCategory.Boundary,
		tcResult.Summary.ByCategory.Hallucination)
	fmt.Fprintf(os.Stderr, "  API Endpoints:       %d\n", len(result.APISpec.Paths))
	fmt.Fprintf(os.Stderr, "  Impl Skeletons:      %d\n", len(result.ImplementationSkeletons))

	funcCount := 0
	for _, skel := range result.ImplementationSkeletons {
		funcCount += len(skel.Functions)
	}
	fmt.Fprintf(os.Stderr, "  Functions:           %d\n", funcCount)
	fmt.Fprintf(os.Stderr, "  Feature Tickets:     %d\n", len(ftResult.FeatureTickets))
	fmt.Fprintf(os.Stderr, "  Service Boundaries:  %d\n", len(sbResult.Services))
	fmt.Fprintf(os.Stderr, "  Domain Events:       %d\n", len(evResult.DomainEvents))
	fmt.Fprintf(os.Stderr, "  Commands:            %d\n", len(evResult.Commands))
	fmt.Fprintf(os.Stderr, "  Integration Events:  %d\n", len(evResult.IntegrationEvents))
	fmt.Fprintf(os.Stderr, "  Graph Components:    %d\n", len(dgResult.Components))
	fmt.Fprintf(os.Stderr, "  Graph Dependencies:  %d\n", len(dgResult.Dependencies))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Output:")
	fmt.Fprintf(os.Stderr, "  %s\n", tcPath)
	fmt.Fprintf(os.Stderr, "  %s\n", apiPath)
	fmt.Fprintf(os.Stderr, "  %s\n", implPath)
	fmt.Fprintf(os.Stderr, "  %s\n", ftPath)
	fmt.Fprintf(os.Stderr, "  %s\n", sbPath)
	fmt.Fprintf(os.Stderr, "  %s\n", evPath)
	fmt.Fprintf(os.Stderr, "  %s\n", dgPath)

	return nil
}

func writeSkeletons(path string, skeletons []ImplementationSkeleton) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	fm := formatter.DefaultFrontmatter("Implementation Skeletons", timestamp, "L3")
	fmt.Fprintf(f, "%s", formatter.FormatHeaderWithFrontmatter(fm))
	fmt.Fprintf(f, "---\n\n")

	for i, skel := range skeletons {
		// Generate SKEL-XXX-NNN ID from service name
		skelID := fmt.Sprintf("SKEL-%s-%03d", extractDomainCode(skel.Name), i+1)
		fmt.Fprintf(f, "## %s: %s (%s)\n\n", skelID, skel.Name, skel.Type)

		if len(skel.Dependencies) > 0 {
			fmt.Fprintf(f, "**Dependencies:** %v\n\n", skel.Dependencies)
		}

		fmt.Fprintf(f, "**Related Specs:** %v\n\n", skel.RelatedSpecs)

		fmt.Fprintf(f, "### Functions\n\n")
		for _, fn := range skel.Functions {
			fmt.Fprintf(f, "#### `%s`\n\n", fn.Name)
			fmt.Fprintf(f, "```typescript\n%s\n```\n\n", fn.Signature)
			fmt.Fprintf(f, "%s\n\n", fn.Description)

			if len(fn.Steps) > 0 {
				fmt.Fprintf(f, "**Implementation Steps:**\n")
				for i, step := range fn.Steps {
					fmt.Fprintf(f, "%d. %s\n", i+1, step)
				}
				fmt.Fprintf(f, "\n")
			}

			if len(fn.ErrorCases) > 0 {
				fmt.Fprintf(f, "**Error Cases:** %v\n\n", fn.ErrorCases)
			}
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeFeatureTickets(path string, tickets []FeatureTicket) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	fm := formatter.DefaultFrontmatter("Feature Definition Tickets", timestamp, "L3")
	fmt.Fprintf(f, "%s", formatter.FormatHeaderWithFrontmatter(fm))
	fmt.Fprintf(f, "---\n\n")

	for _, t := range tickets {
		fmt.Fprintf(f, "## %s: %s\n\n", t.ID, t.Title)
		fmt.Fprintf(f, "**Status:** %s | **Priority:** %s | **Complexity:** %s\n\n", t.Status, t.Priority, t.EstimatedComplexity)

		fmt.Fprintf(f, "### Business Goal\n%s\n\n", t.BusinessGoal)
		fmt.Fprintf(f, "### User Story\n%s\n\n", t.UserStory)

		if len(t.AcceptanceCriteriaRefs) > 0 {
			fmt.Fprintf(f, "### Acceptance Criteria References\n")
			for _, ac := range t.AcceptanceCriteriaRefs {
				fmt.Fprintf(f, "- %s\n", ac)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(t.NFR) > 0 {
			fmt.Fprintf(f, "### Non-Functional Requirements\n")
			for _, nfr := range t.NFR {
				fmt.Fprintf(f, "- %s\n", nfr)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(t.Dependencies) > 0 {
			fmt.Fprintf(f, "### Dependencies\n")
			for _, dep := range t.Dependencies {
				fmt.Fprintf(f, "- %s\n", dep)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(t.ImpactAreas) > 0 {
			fmt.Fprintf(f, "### Impact Areas\n")
			for _, area := range t.ImpactAreas {
				fmt.Fprintf(f, "- %s\n", area)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(t.OutOfScope) > 0 {
			fmt.Fprintf(f, "### Out of Scope\n")
			for _, oos := range t.OutOfScope {
				fmt.Fprintf(f, "- %s\n", oos)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeServiceBoundaries(path string, services []ServiceBoundary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	fm := formatter.DefaultFrontmatter("Service Boundaries", timestamp, "L3")
	fmt.Fprintf(f, "%s", formatter.FormatHeaderWithFrontmatter(fm))
	fmt.Fprintf(f, "---\n\n")

	for _, s := range services {
		fmt.Fprintf(f, "## %s: %s\n\n", s.ID, s.Name)
		fmt.Fprintf(f, "**Purpose:** %s\n\n", s.Purpose)
		fmt.Fprintf(f, "**API Base:** `%s`\n\n", s.APIBase)
		fmt.Fprintf(f, "**Separation Reason:** %s\n\n", s.SeparationReason)

		if len(s.Capabilities) > 0 {
			fmt.Fprintf(f, "### Capabilities\n")
			for _, cap := range s.Capabilities {
				fmt.Fprintf(f, "- %s\n", cap)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(s.Inputs) > 0 {
			fmt.Fprintf(f, "### Inputs\n")
			for _, in := range s.Inputs {
				fmt.Fprintf(f, "- **%s:** %s\n", in.Type, in.Name)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(s.Outputs) > 0 {
			fmt.Fprintf(f, "### Outputs\n")
			for _, out := range s.Outputs {
				fmt.Fprintf(f, "- **%s:** %s\n", out.Type, out.Name)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(s.OwnedAggregates) > 0 {
			fmt.Fprintf(f, "### Owned Aggregates\n")
			for _, agg := range s.OwnedAggregates {
				fmt.Fprintf(f, "- %s\n", agg)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(s.Dependencies) > 0 {
			fmt.Fprintf(f, "### Dependencies\n")
			for _, dep := range s.Dependencies {
				fmt.Fprintf(f, "- **%s** (%s): %s\n", dep.Service, dep.Type, dep.Reason)
			}
			fmt.Fprintf(f, "\n")
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeEventDesign(path string, events []DomainEvent, commands []Command, integrationEvents []IntegrationEvent) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	fm := formatter.DefaultFrontmatter("Event & Message Design", timestamp, "L3")
	fmt.Fprintf(f, "%s", formatter.FormatHeaderWithFrontmatter(fm))
	fmt.Fprintf(f, "---\n\n")

	// Domain Events
	fmt.Fprintf(f, "## Domain Events\n\n")
	for _, e := range events {
		fmt.Fprintf(f, "### %s: %s\n\n", e.ID, e.Name)
		fmt.Fprintf(f, "**Purpose:** %s\n\n", e.Purpose)
		fmt.Fprintf(f, "**Trigger:** %s\n\n", e.Trigger)
		fmt.Fprintf(f, "**Aggregate:** %s | **Version:** %s\n\n", e.Aggregate, e.Version)

		if len(e.Payload) > 0 {
			fmt.Fprintf(f, "**Payload:**\n")
			for _, p := range e.Payload {
				fmt.Fprintf(f, "- `%s`: %s\n", p.Field, p.Type)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(e.InvariantsReflected) > 0 {
			fmt.Fprintf(f, "**Invariants Reflected:** %v\n\n", e.InvariantsReflected)
		}

		if len(e.Consumers) > 0 {
			fmt.Fprintf(f, "**Consumers:** %v\n\n", e.Consumers)
		}

		fmt.Fprintf(f, "---\n\n")
	}

	// Commands
	fmt.Fprintf(f, "## Commands\n\n")
	for _, c := range commands {
		fmt.Fprintf(f, "### %s: %s\n\n", c.ID, c.Name)
		fmt.Fprintf(f, "**Intent:** %s\n\n", c.Intent)
		fmt.Fprintf(f, "**Expected Outcome:** %s\n\n", c.ExpectedOutcome)

		if len(c.RequiredData) > 0 {
			fmt.Fprintf(f, "**Required Data:**\n")
			for _, d := range c.RequiredData {
				fmt.Fprintf(f, "- `%s`: %s\n", d.Field, d.Type)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(c.FailureConditions) > 0 {
			fmt.Fprintf(f, "**Failure Conditions:** %v\n\n", c.FailureConditions)
		}

		fmt.Fprintf(f, "---\n\n")
	}

	// Integration Events
	fmt.Fprintf(f, "## Integration Events\n\n")
	for _, ie := range integrationEvents {
		fmt.Fprintf(f, "### %s: %s\n\n", ie.ID, ie.Name)
		fmt.Fprintf(f, "**Purpose:** %s\n\n", ie.Purpose)
		fmt.Fprintf(f, "**Source:** %s\n\n", ie.Source)

		if len(ie.Consumers) > 0 {
			fmt.Fprintf(f, "**Consumers:** %v\n\n", ie.Consumers)
		}

		if len(ie.Payload) > 0 {
			fmt.Fprintf(f, "**Payload:** %v\n\n", ie.Payload)
		}

		fmt.Fprintf(f, "---\n\n")
	}

	return nil
}

func writeDependencyGraph(path string, components []GraphComponent, dependencies []GraphDependency) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	fm := formatter.DefaultFrontmatter("Dependency Graph", timestamp, "L3")
	fmt.Fprintf(f, "%s", formatter.FormatHeaderWithFrontmatter(fm))
	fmt.Fprintf(f, "---\n\n")

	// Build component ID map: original name -> DEP-XXX-NNN
	compIDMap := make(map[string]string)
	for i, c := range components {
		depID := fmt.Sprintf("DEP-%s-%03d", extractDomainCode(c.ID), i+1)
		compIDMap[c.ID] = depID
	}

	// Mermaid diagram
	fmt.Fprintf(f, "## Visual Overview\n\n")
	fmt.Fprintf(f, "```mermaid\ngraph TD\n")
	for _, c := range components {
		shape := "([%s])"
		if c.Type == "external" {
			shape = "[[%s]]"
		} else if c.Type == "domain_service" {
			shape = "[%s]"
		}
		id := sanitizeID(compIDMap[c.ID])
		fmt.Fprintf(f, "    %s"+shape+"\n", id, c.ID)
	}
	fmt.Fprintf(f, "\n")
	for _, d := range dependencies {
		fromID := sanitizeID(compIDMap[d.From])
		toID := sanitizeID(compIDMap[d.To])
		if fromID == "" {
			fromID = sanitizeID(d.From)
		}
		if toID == "" {
			toID = sanitizeID(d.To)
		}
		arrow := "-->"
		if d.Type == "async" {
			arrow = "-.->|async|"
		} else if d.Type == "external" {
			arrow = "==>|ext|"
		}
		fmt.Fprintf(f, "    %s %s %s\n", fromID, arrow, toID)
	}
	fmt.Fprintf(f, "```\n\n")

	// Components section with headers for ID detection
	fmt.Fprintf(f, "## Components\n\n")
	for i, c := range components {
		depID := fmt.Sprintf("DEP-%s-%03d", extractDomainCode(c.ID), i+1)
		fmt.Fprintf(f, "### %s: %s\n\n", depID, c.ID)
		fmt.Fprintf(f, "- **Type:** %s\n", c.Type)
		fmt.Fprintf(f, "- **Description:** %s\n\n", c.Description)
	}

	// Dependencies section
	fmt.Fprintf(f, "## Dependencies\n\n")
	fmt.Fprintf(f, "| From | To | Type | Description |\n")
	fmt.Fprintf(f, "|---|---|---|---|\n")
	for _, d := range dependencies {
		fromID := compIDMap[d.From]
		toID := compIDMap[d.To]
		if fromID == "" {
			fromID = d.From
		}
		if toID == "" {
			toID = d.To
		}
		fmt.Fprintf(f, "| %s | %s | %s | %s |\n", fromID, toID, d.Type, d.Description)
	}
	fmt.Fprintf(f, "\n")

	return nil
}

func sanitizeID(id string) string {
	result := ""
	for _, c := range id {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			result += string(c)
		}
	}
	return result
}

// extractDomainCode extracts a short code from a service name
// e.g., "ProductService" -> "PROD", "CustomerService" -> "CUST"
func extractDomainCode(name string) string {
	// Remove common suffixes
	name = strings.TrimSuffix(name, "Service")
	name = strings.TrimSuffix(name, "Controller")
	name = strings.TrimSuffix(name, "Repository")

	// Take first 4 chars uppercase
	if len(name) > 4 {
		name = name[:4]
	}
	return strings.ToUpper(name)
}

// writeL3TestCases writes test cases to markdown file using the formatter
func writeL3TestCases(path string, testCases []generator.TestCase, summary generator.TDAISummary) error {
	timestamp := time.Now().Format(time.RFC3339)

	// Convert generator types to formatter types
	fmtCases := make([]formatter.TestCase, len(testCases))
	for i, tc := range testCases {
		testData := make([]formatter.TestData, len(tc.TestData))
		for j, td := range tc.TestData {
			testData[j] = formatter.TestData{
				Field: td.Field,
				Value: td.Value,
				Notes: td.Notes,
			}
		}
		fmtCases[i] = formatter.TestCase{
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

	// Convert summary
	var fmtSummary formatter.TDAISummary
	fmtSummary.Total = summary.Total
	fmtSummary.ByCategory.Positive = summary.ByCategory.Positive
	fmtSummary.ByCategory.Negative = summary.ByCategory.Negative
	fmtSummary.ByCategory.Boundary = summary.ByCategory.Boundary
	fmtSummary.ByCategory.Hallucination = summary.ByCategory.Hallucination
	fmtSummary.Coverage.ACsCovered = summary.Coverage.ACsCovered
	fmtSummary.Coverage.PositiveRatio = summary.Coverage.PositiveRatio
	fmtSummary.Coverage.NegativeRatio = summary.Coverage.NegativeRatio
	fmtSummary.Coverage.HasHallucinationTests = summary.Coverage.HasHallucinationTests

	content := formatter.FormatTestCases(fmtCases, fmtSummary, timestamp)
	return os.WriteFile(path, []byte(content), 0644)
}
