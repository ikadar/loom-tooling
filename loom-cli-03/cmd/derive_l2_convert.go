// Package cmd provides CLI commands for loom-cli.
//
// This file contains JSON to Markdown conversion utilities for L2 documents.
//
// Implements: l2/package-structure.md
// See: l2/internal-api.md (formatter package)
package cmd

// Note: JSON to Markdown conversion is now handled by the internal/formatter package.
// This file exists for backwards compatibility and contains helper utilities.

// convertTechSpecsToMarkdown is a passthrough to formatter.FormatTechSpecs.
// Deprecated: Use formatter.FormatTechSpecs directly.
func convertTechSpecsToMarkdown() {
	// Implemented in internal/formatter/techspecs.go
}

// convertContractsToMarkdown is a passthrough to formatter.FormatInterfaceContracts.
// Deprecated: Use formatter.FormatInterfaceContracts directly.
func convertContractsToMarkdown() {
	// Implemented in internal/formatter/contracts.go
}

// convertAggregateToMarkdown is a passthrough to formatter.FormatAggregateDesign.
// Deprecated: Use formatter.FormatAggregateDesign directly.
func convertAggregateToMarkdown() {
	// Implemented in internal/formatter/aggregates.go
}

// convertSequenceToMarkdown is a passthrough to formatter.FormatSequenceDesign.
// Deprecated: Use formatter.FormatSequenceDesign directly.
func convertSequenceToMarkdown() {
	// Implemented in internal/formatter/sequences.go
}

// convertDataModelToMarkdown is a passthrough to formatter.FormatDataModel.
// Deprecated: Use formatter.FormatDataModel directly.
func convertDataModelToMarkdown() {
	// Implemented in internal/formatter/datamodel.go
}
