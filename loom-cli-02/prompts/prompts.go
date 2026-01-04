// Package prompts provides embedded prompt templates for Claude.
//
// Implements: l2/package-structure.md PKG-005
// See: l2/prompt-catalog.md
package prompts

import (
	_ "embed"
)

// =============================================================================
// Analyze Phase Prompts (3)
// =============================================================================

//go:embed domain-discovery.md
var DomainDiscovery string

//go:embed entity-analysis.md
var EntityAnalysis string

//go:embed operation-analysis.md
var OperationAnalysis string

// =============================================================================
// Interview Phase Prompt (1)
// =============================================================================

//go:embed interview.md
var Interview string

// =============================================================================
// L1 Derivation Prompts (3)
// =============================================================================

//go:embed derivation.md
var Derivation string

//go:embed derive-domain-model.md
var DeriveDomainModel string

//go:embed derive-bounded-context.md
var DeriveBoundedContext string

// =============================================================================
// L2 Derivation Prompts (6)
// =============================================================================

//go:embed derive-l2.md
var DeriveL2 string

//go:embed derive-tech-specs.md
var DeriveTechSpecs string

//go:embed derive-interface-contracts.md
var DeriveInterfaceContracts string

//go:embed derive-aggregate-design.md
var DeriveAggregateDesign string

//go:embed derive-sequence-design.md
var DeriveSequenceDesign string

//go:embed derive-data-model.md
var DeriveDataModel string

// =============================================================================
// L3 Derivation Prompts (8)
// =============================================================================

//go:embed derive-test-cases.md
var DeriveTestCases string

//go:embed derive-l3.md
var DeriveL3 string

//go:embed derive-l3-api.md
var DeriveL3API string

//go:embed derive-l3-skeletons.md
var DeriveL3Skeletons string

//go:embed derive-feature-tickets.md
var DeriveFeatureTickets string

//go:embed derive-service-boundaries.md
var DeriveServiceBoundaries string

//go:embed derive-event-design.md
var DeriveEventDesign string

//go:embed derive-dependency-graph.md
var DeriveDependencyGraph string

// =============================================================================
// Prompt Helpers
// =============================================================================

// All returns all prompt names for validation.
func All() []string {
	return []string{
		// Analyze (3)
		"domain-discovery",
		"entity-analysis",
		"operation-analysis",
		// Interview (1)
		"interview",
		// L1 Derivation (3)
		"derivation",
		"derive-domain-model",
		"derive-bounded-context",
		// L2 Derivation (6)
		"derive-l2",
		"derive-tech-specs",
		"derive-interface-contracts",
		"derive-aggregate-design",
		"derive-sequence-design",
		"derive-data-model",
		// L3 Derivation (8)
		"derive-test-cases",
		"derive-l3",
		"derive-l3-api",
		"derive-l3-skeletons",
		"derive-feature-tickets",
		"derive-service-boundaries",
		"derive-event-design",
		"derive-dependency-graph",
	}
}

// Count returns the total number of prompts.
func Count() int {
	return len(All())
}

// Get returns a prompt by name.
func Get(name string) string {
	switch name {
	case "domain-discovery":
		return DomainDiscovery
	case "entity-analysis":
		return EntityAnalysis
	case "operation-analysis":
		return OperationAnalysis
	case "interview":
		return Interview
	case "derivation":
		return Derivation
	case "derive-domain-model":
		return DeriveDomainModel
	case "derive-bounded-context":
		return DeriveBoundedContext
	case "derive-l2":
		return DeriveL2
	case "derive-tech-specs":
		return DeriveTechSpecs
	case "derive-interface-contracts":
		return DeriveInterfaceContracts
	case "derive-aggregate-design":
		return DeriveAggregateDesign
	case "derive-sequence-design":
		return DeriveSequenceDesign
	case "derive-data-model":
		return DeriveDataModel
	case "derive-test-cases":
		return DeriveTestCases
	case "derive-l3":
		return DeriveL3
	case "derive-l3-api":
		return DeriveL3API
	case "derive-l3-skeletons":
		return DeriveL3Skeletons
	case "derive-feature-tickets":
		return DeriveFeatureTickets
	case "derive-service-boundaries":
		return DeriveServiceBoundaries
	case "derive-event-design":
		return DeriveEventDesign
	case "derive-dependency-graph":
		return DeriveDependencyGraph
	default:
		return ""
	}
}
