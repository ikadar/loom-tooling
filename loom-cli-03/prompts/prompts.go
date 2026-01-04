// Package prompts provides embedded prompt templates for Claude API calls.
//
// Implements: l2/package-structure.md PKG-011, l2/prompt-catalog.md
// See: l2/tech-specs.md TS-ARCH-002 (context injection pattern)
//
// All prompts use XML-style context markers:
//
//	<context>
//	</context>
//
// Context is injected using claude.BuildPrompt(template, documents...).
package prompts

import (
	"embed"
)

//go:embed *.md
var promptFS embed.FS

// Prompt template variables loaded from embedded markdown files.
var (
	// ==========================================================================
	// Analyze prompts (3)
	// ==========================================================================

	// DomainDiscovery discovers entities, operations, relationships from L0 input.
	// Implements: PRM-ANL-001
	DomainDiscovery string

	// EntityAnalysis analyzes discovered entities for completeness.
	// Implements: PRM-ANL-002
	EntityAnalysis string

	// OperationAnalysis analyzes discovered operations for completeness.
	// Implements: PRM-ANL-003
	OperationAnalysis string

	// ==========================================================================
	// Interview prompts (1)
	// ==========================================================================

	// InterviewPrompt generates follow-up questions for ambiguities.
	// Implements: PRM-INT-001
	InterviewPrompt string

	// ==========================================================================
	// L1 Derivation prompts (3)
	// ==========================================================================

	// DeriveDomainModel generates formal L1 domain model.
	// Implements: PRM-DRV-002
	DeriveDomainModel string

	// DeriveBoundedContext generates bounded context map.
	// Implements: PRM-DRV-003
	DeriveBoundedContext string

	// Derivation generates acceptance criteria and business rules.
	// Implements: PRM-DRV-001
	Derivation string

	// ==========================================================================
	// L2 Derivation prompts (6)
	// ==========================================================================

	// DeriveL2 is the master L2 derivation prompt.
	// Implements: PRM-L2-001
	DeriveL2 string

	// DeriveTechSpecs generates technical specifications.
	// Implements: PRM-L2-002
	DeriveTechSpecs string

	// DeriveInterfaceContracts generates interface contracts.
	// Implements: PRM-L2-003
	DeriveInterfaceContracts string

	// DeriveAggregateDesign generates aggregate design.
	// Implements: PRM-L2-004
	DeriveAggregateDesign string

	// DeriveSequenceDesign generates sequence diagrams.
	// Implements: PRM-L2-005
	DeriveSequenceDesign string

	// DeriveDataModel generates initial data model.
	// Implements: PRM-L2-006
	DeriveDataModel string

	// ==========================================================================
	// L3 Derivation prompts (8)
	// ==========================================================================

	// DeriveTestCases generates test cases from acceptance criteria.
	// Implements: PRM-L3-001
	DeriveTestCases string

	// DeriveL3 is the master L3 derivation prompt.
	// Implements: PRM-L3-COMBINED
	DeriveL3 string

	// DeriveL3API generates OpenAPI specification.
	// Implements: PRM-L3-002
	DeriveL3API string

	// DeriveL3Skeletons generates implementation skeletons.
	// Implements: PRM-L3-003
	DeriveL3Skeletons string

	// DeriveFeatureTickets generates feature tickets.
	// Implements: PRM-L3-004
	DeriveFeatureTickets string

	// DeriveServiceBoundaries generates service boundaries.
	// Implements: PRM-L3-005
	DeriveServiceBoundaries string

	// DeriveEventDesign generates event/message design.
	// Implements: PRM-L3-006
	DeriveEventDesign string

	// DeriveDependencyGraph generates dependency graph.
	// Implements: PRM-L3-007
	DeriveDependencyGraph string
)

func init() {
	// Load all prompts from embedded files
	DomainDiscovery = mustLoad("domain-discovery.md")
	EntityAnalysis = mustLoad("entity-analysis.md")
	OperationAnalysis = mustLoad("operation-analysis.md")

	InterviewPrompt = mustLoad("interview.md")

	DeriveDomainModel = mustLoad("derive-domain-model.md")
	DeriveBoundedContext = mustLoad("derive-bounded-context.md")
	Derivation = mustLoad("derivation.md")

	DeriveL2 = mustLoad("derive-l2.md")
	DeriveTechSpecs = mustLoad("derive-tech-specs.md")
	DeriveInterfaceContracts = mustLoad("derive-interface-contracts.md")
	DeriveAggregateDesign = mustLoad("derive-aggregate-design.md")
	DeriveSequenceDesign = mustLoad("derive-sequence-design.md")
	DeriveDataModel = mustLoad("derive-data-model.md")

	DeriveTestCases = mustLoad("derive-test-cases.md")
	DeriveL3 = mustLoad("derive-l3.md")
	DeriveL3API = mustLoad("derive-l3-api.md")
	DeriveL3Skeletons = mustLoad("derive-l3-skeletons.md")
	DeriveFeatureTickets = mustLoad("derive-feature-tickets.md")
	DeriveServiceBoundaries = mustLoad("derive-service-boundaries.md")
	DeriveEventDesign = mustLoad("derive-event-design.md")
	DeriveDependencyGraph = mustLoad("derive-dependency-graph.md")
}

// mustLoad reads a file from the embedded filesystem or panics.
func mustLoad(filename string) string {
	data, err := promptFS.ReadFile(filename)
	if err != nil {
		panic("failed to load prompt " + filename + ": " + err.Error())
	}
	return string(data)
}
