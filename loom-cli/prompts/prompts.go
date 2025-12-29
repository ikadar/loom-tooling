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
