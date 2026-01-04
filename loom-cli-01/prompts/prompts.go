// Package prompts provides embedded prompt templates for AI-driven derivation.
//
// Implements: l2/prompt-catalog.md
// See: l2/prompts/ for full prompt content
//
// All prompts use XML-style context injection:
//
//	<context>
//	</context>
//
// Use claude.BuildPrompt() to inject documents into the context section.
package prompts

import (
	"embed"
)

// Embed all prompt markdown files.
// Implements: l2/prompt-catalog.md (Embedding Strategy)
//
//go:embed *.md
var promptFS embed.FS

// Prompt variables loaded at init time.
var (
	// Analyze phase prompts (PRM-ANL-*)
	DomainDiscovery   string // domain-discovery.md - Extract domain model from L0
	EntityAnalysis    string // entity-analysis.md - Enhance entity details
	OperationAnalysis string // operation-analysis.md - Enhance operation details

	// Interview phase prompt (PRM-INT-001)
	Interview string // interview.md - Resolve ambiguities

	// L1 derivation prompts (PRM-DRV-*)
	Derivation           string // derivation.md - Generate AC and BR
	DeriveDomainModel    string // derive-domain-model.md - Generate L1 domain model
	DeriveBoundedContext string // derive-bounded-context.md - Generate L1 bounded context map

	// L2 derivation prompts (PRM-L2-*)
	DeriveL2                 string // derive-l2.md - Combined L2 derivation
	DeriveTechSpecs          string // derive-tech-specs.md - Generate tech specs
	DeriveInterfaceContracts string // derive-interface-contracts.md - Generate interface contracts
	DeriveAggregateDesign    string // derive-aggregate-design.md - Generate aggregate design
	DeriveSequenceDesign     string // derive-sequence-design.md - Generate sequence diagrams
	DeriveDataModel          string // derive-data-model.md - Generate data model

	// L3 derivation prompts (PRM-L3-*)
	DeriveTestCases         string // derive-test-cases.md - Generate TDAI test cases
	DeriveL3                string // derive-l3.md - Combined L3 derivation
	DeriveL3API             string // derive-l3-api.md - Generate OpenAPI spec
	DeriveL3Skeletons       string // derive-l3-skeletons.md - Generate implementation skeletons
	DeriveFeatureTickets    string // derive-feature-tickets.md - Generate feature tickets
	DeriveServiceBoundaries string // derive-service-boundaries.md - Generate service boundaries
	DeriveEventDesign       string // derive-event-design.md - Generate event/message design
	DeriveDependencyGraph   string // derive-dependency-graph.md - Generate dependency graph
)

func init() {
	// Load all prompts from embedded filesystem
	mustLoad := func(name string) string {
		content, err := promptFS.ReadFile(name)
		if err != nil {
			panic("failed to load prompt " + name + ": " + err.Error())
		}
		return string(content)
	}

	// Analyze phase
	DomainDiscovery = mustLoad("domain-discovery.md")
	EntityAnalysis = mustLoad("entity-analysis.md")
	OperationAnalysis = mustLoad("operation-analysis.md")

	// Interview phase
	Interview = mustLoad("interview.md")

	// L1 derivation
	Derivation = mustLoad("derivation.md")
	DeriveDomainModel = mustLoad("derive-domain-model.md")
	DeriveBoundedContext = mustLoad("derive-bounded-context.md")

	// L2 derivation
	DeriveL2 = mustLoad("derive-l2.md")
	DeriveTechSpecs = mustLoad("derive-tech-specs.md")
	DeriveInterfaceContracts = mustLoad("derive-interface-contracts.md")
	DeriveAggregateDesign = mustLoad("derive-aggregate-design.md")
	DeriveSequenceDesign = mustLoad("derive-sequence-design.md")
	DeriveDataModel = mustLoad("derive-data-model.md")

	// L3 derivation
	DeriveTestCases = mustLoad("derive-test-cases.md")
	DeriveL3 = mustLoad("derive-l3.md")
	DeriveL3API = mustLoad("derive-l3-api.md")
	DeriveL3Skeletons = mustLoad("derive-l3-skeletons.md")
	DeriveFeatureTickets = mustLoad("derive-feature-tickets.md")
	DeriveServiceBoundaries = mustLoad("derive-service-boundaries.md")
	DeriveEventDesign = mustLoad("derive-event-design.md")
	DeriveDependencyGraph = mustLoad("derive-dependency-graph.md")
}
