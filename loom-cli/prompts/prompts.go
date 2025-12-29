package prompts

import (
	_ "embed"
)

// Prompts are loaded from external .md files using go:embed
// This allows easy editing without recompiling

//go:embed domain-discovery.md
var DomainDiscovery string

//go:embed entity-analysis.md
var EntityAnalysis string

//go:embed operation-analysis.md
var OperationAnalysis string

//go:embed derivation.md
var DerivationPrompt string

//go:embed interview.md
var InterviewPrompt string

//go:embed derive-l2.md
var DeriveL2 string

//go:embed derive-test-cases.md
var DeriveTestCases string

//go:embed derive-tech-specs.md
var DeriveTechSpecs string

//go:embed derive-l3.md
var DeriveL3 string

//go:embed derive-domain-model.md
var DeriveDomainModel string

//go:embed derive-bounded-context.md
var DeriveBoundedContext string

//go:embed derive-interface-contracts.md
var DeriveInterfaceContracts string

//go:embed derive-sequence-design.md
var DeriveSequenceDesign string

//go:embed derive-aggregate-design.md
var DeriveAggregateDesign string

//go:embed derive-data-model.md
var DeriveDataModel string

//go:embed derive-feature-tickets.md
var DeriveFeatureTickets string

//go:embed derive-service-boundaries.md
var DeriveServiceBoundaries string

//go:embed derive-event-design.md
var DeriveEventDesign string

//go:embed derive-dependency-graph.md
var DeriveDependencyGraph string
