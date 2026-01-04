// Package cmd provides CLI commands for loom-cli.
//
// This file provides JSON to Markdown conversion utilities for L2 documents.
//
// Implements: l2/package-structure.md
package cmd

import (
	"fmt"
	"strings"
	"time"
)

// L2 Document Types for JSON→Markdown conversion

// TechSpec represents a technical specification.
type TechSpec struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Implementation   string   `json:"implementation"`
	ValidationPoints []string `json:"validation_points"`
	ErrorHandling    []string `json:"error_handling"`
	RelatedACs       []string `json:"related_acs"`
	RelatedBRs       []string `json:"related_brs"`
}

// InterfaceContract represents an API contract.
type InterfaceContract struct {
	ID          string              `json:"id"`
	ServiceName string              `json:"service_name"`
	Purpose     string              `json:"purpose"`
	BaseURL     string              `json:"base_url"`
	Operations  []ContractOperation `json:"operations"`
}

// ContractOperation represents an API operation.
type ContractOperation struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Request     string   `json:"request"`
	Response    string   `json:"response"`
	ErrorCodes  []string `json:"error_codes"`
}

// AggregateSpec represents an aggregate design.
type AggregateSpec struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Purpose    string   `json:"purpose"`
	Root       string   `json:"root"`
	Invariants []string `json:"invariants"`
	Entities   []string `json:"entities"`
	ValueObjs  []string `json:"value_objects"`
	Behaviors  []string `json:"behaviors"`
	Events     []string `json:"events"`
}

// SequenceSpec represents a sequence design.
type SequenceSpec struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Trigger      string   `json:"trigger"`
	Participants []string `json:"participants"`
	Steps        []string `json:"steps"`
	Outcome      string   `json:"outcome"`
	MermaidCode  string   `json:"mermaid_code"`
	RelatedACs   []string `json:"related_acs"`
}

// DataTable represents a data model table.
type DataTable struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Aggregate   string      `json:"aggregate"`
	Purpose     string      `json:"purpose"`
	Fields      []DataField `json:"fields"`
	PrimaryKey  string      `json:"primary_key"`
	Indexes     []string    `json:"indexes"`
	ForeignKeys []string    `json:"foreign_keys"`
}

// DataField represents a table field.
type DataField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Nullable    bool   `json:"nullable"`
	Default     string `json:"default"`
	Constraints string `json:"constraints"`
}

// convertTechSpecsToMarkdown converts tech specs to markdown.
func convertTechSpecsToMarkdown(specs []TechSpec) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(fmt.Sprintf(`---
title: "Technical Specifications"
generated: %s
status: draft
level: L2
---

# Technical Specifications

## Overview

This document contains the technical specifications for the system.

`, time.Now().Format(time.RFC3339)))

	// Specs
	sb.WriteString("## Specifications\n\n")
	for _, spec := range specs {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", spec.ID, spec.Name))
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", spec.Description))
		sb.WriteString(fmt.Sprintf("**Implementation:**\n%s\n\n", spec.Implementation))

		if len(spec.ValidationPoints) > 0 {
			sb.WriteString("**Validation Points:**\n")
			for _, vp := range spec.ValidationPoints {
				sb.WriteString(fmt.Sprintf("- %s\n", vp))
			}
			sb.WriteString("\n")
		}

		if len(spec.ErrorHandling) > 0 {
			sb.WriteString("**Error Handling:**\n")
			for _, eh := range spec.ErrorHandling {
				sb.WriteString(fmt.Sprintf("- %s\n", eh))
			}
			sb.WriteString("\n")
		}

		sb.WriteString(fmt.Sprintf("**Related ACs:** %s\n\n", strings.Join(spec.RelatedACs, ", ")))
		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// convertContractsToMarkdown converts interface contracts to markdown.
func convertContractsToMarkdown(contracts []InterfaceContract) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(fmt.Sprintf(`---
title: "Interface Contracts"
generated: %s
status: draft
level: L2
---

# Interface Contracts

## Overview

This document defines the API contracts for the system.

`, time.Now().Format(time.RFC3339)))

	// Contracts
	sb.WriteString("## Contracts\n\n")
	for _, contract := range contracts {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", contract.ID, contract.ServiceName))
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", contract.Purpose))
		sb.WriteString(fmt.Sprintf("**Base URL:** `%s`\n\n", contract.BaseURL))

		sb.WriteString("**Operations:**\n\n")
		for _, op := range contract.Operations {
			sb.WriteString(fmt.Sprintf("#### %s %s\n\n", op.Method, op.Path))
			sb.WriteString(fmt.Sprintf("%s\n\n", op.Description))
			sb.WriteString(fmt.Sprintf("**Request:** `%s`\n\n", op.Request))
			sb.WriteString(fmt.Sprintf("**Response:** `%s`\n\n", op.Response))

			if len(op.ErrorCodes) > 0 {
				sb.WriteString("**Error Codes:**\n")
				for _, ec := range op.ErrorCodes {
					sb.WriteString(fmt.Sprintf("- %s\n", ec))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// convertAggregateToMarkdown converts aggregate designs to markdown.
func convertAggregateToMarkdown(aggregates []AggregateSpec) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(fmt.Sprintf(`---
title: "Aggregate Design"
generated: %s
status: draft
level: L2
---

# Aggregate Design

## Overview

This document defines the aggregate designs following DDD principles.

`, time.Now().Format(time.RFC3339)))

	// Aggregates
	sb.WriteString("## Aggregates\n\n")
	for _, agg := range aggregates {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", agg.ID, agg.Name))
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", agg.Purpose))
		sb.WriteString(fmt.Sprintf("**Aggregate Root:** %s\n\n", agg.Root))

		if len(agg.Invariants) > 0 {
			sb.WriteString("**Invariants:**\n")
			for _, inv := range agg.Invariants {
				sb.WriteString(fmt.Sprintf("- %s\n", inv))
			}
			sb.WriteString("\n")
		}

		if len(agg.Entities) > 0 {
			sb.WriteString(fmt.Sprintf("**Entities:** %s\n\n", strings.Join(agg.Entities, ", ")))
		}

		if len(agg.ValueObjs) > 0 {
			sb.WriteString(fmt.Sprintf("**Value Objects:** %s\n\n", strings.Join(agg.ValueObjs, ", ")))
		}

		if len(agg.Behaviors) > 0 {
			sb.WriteString("**Behaviors:**\n")
			for _, b := range agg.Behaviors {
				sb.WriteString(fmt.Sprintf("- %s\n", b))
			}
			sb.WriteString("\n")
		}

		if len(agg.Events) > 0 {
			sb.WriteString(fmt.Sprintf("**Events:** %s\n\n", strings.Join(agg.Events, ", ")))
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// convertSequenceToMarkdown converts sequence designs to markdown.
func convertSequenceToMarkdown(sequences []SequenceSpec) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(fmt.Sprintf(`---
title: "Sequence Design"
generated: %s
status: draft
level: L2
---

# Sequence Design

## Overview

This document defines the sequence diagrams for key operations.

`, time.Now().Format(time.RFC3339)))

	// Sequences
	sb.WriteString("## Sequences\n\n")
	for _, seq := range sequences {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", seq.ID, seq.Name))
		sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", seq.Description))
		sb.WriteString(fmt.Sprintf("**Trigger:** %s\n\n", seq.Trigger))
		sb.WriteString(fmt.Sprintf("**Participants:** %s\n\n", strings.Join(seq.Participants, ", ")))

		if len(seq.Steps) > 0 {
			sb.WriteString("**Steps:**\n")
			for i, step := range seq.Steps {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
			}
			sb.WriteString("\n")
		}

		sb.WriteString(fmt.Sprintf("**Outcome:** %s\n\n", seq.Outcome))

		if seq.MermaidCode != "" {
			sb.WriteString("**Diagram:**\n\n")
			sb.WriteString("```mermaid\n")
			sb.WriteString(seq.MermaidCode)
			sb.WriteString("\n```\n\n")
		}

		sb.WriteString(fmt.Sprintf("**Related ACs:** %s\n\n", strings.Join(seq.RelatedACs, ", ")))
		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// convertDataModelToMarkdown converts data model to markdown.
func convertDataModelToMarkdown(tables []DataTable, enums map[string][]string) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString(fmt.Sprintf(`---
title: "Initial Data Model"
generated: %s
status: draft
level: L2
---

# Initial Data Model

## Overview

This document defines the initial data model for the system.

`, time.Now().Format(time.RFC3339)))

	// Tables
	sb.WriteString("## Tables\n\n")
	for _, table := range tables {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", table.ID, table.Name))
		sb.WriteString(fmt.Sprintf("**Aggregate:** %s\n\n", table.Aggregate))
		sb.WriteString(fmt.Sprintf("**Purpose:** %s\n\n", table.Purpose))

		sb.WriteString("**Fields:**\n\n")
		sb.WriteString("| Name | Type | Nullable | Default | Constraints |\n")
		sb.WriteString("|------|------|----------|---------|-------------|\n")
		for _, field := range table.Fields {
			nullable := "No"
			if field.Nullable {
				nullable = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				field.Name, field.Type, nullable, field.Default, field.Constraints))
		}
		sb.WriteString("\n")

		sb.WriteString(fmt.Sprintf("**Primary Key:** %s\n\n", table.PrimaryKey))

		if len(table.Indexes) > 0 {
			sb.WriteString("**Indexes:**\n")
			for _, idx := range table.Indexes {
				sb.WriteString(fmt.Sprintf("- %s\n", idx))
			}
			sb.WriteString("\n")
		}

		if len(table.ForeignKeys) > 0 {
			sb.WriteString("**Foreign Keys:**\n")
			for _, fk := range table.ForeignKeys {
				sb.WriteString(fmt.Sprintf("- %s\n", fk))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	// Enums
	if len(enums) > 0 {
		sb.WriteString("## Enums\n\n")
		for name, values := range enums {
			sb.WriteString(fmt.Sprintf("### %s\n\n", name))
			sb.WriteString("Values:\n")
			for _, v := range values {
				sb.WriteString(fmt.Sprintf("- `%s`\n", v))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
