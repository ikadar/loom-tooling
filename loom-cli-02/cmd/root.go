// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md (CLI interface)
// Implements: l2/tech-specs.md TS-ARCH-001 (Command router architecture)
// See: l0/decisions.md DEC-L1-001 (stdlib only, no cobra)
package cmd

import (
	"fmt"
	"os"

	"loom-cli/internal/domain"
)

// Version is set at build time
var Version = "v0.3.0"

// Execute runs the CLI and returns an exit code.
//
// Implements: l2/interface-contracts.md (Exit codes)
// Exit codes:
//   - 0: Success
//   - 1: Error
//   - 100: Interview has more questions
func Execute() int {
	if len(os.Args) < 2 {
		printUsage()
		return domain.ExitCodeError
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "analyze":
		return runAnalyze(args)
	case "interview":
		return runInterview(args)
	case "derive":
		return runDerive(args)
	case "derive-l2":
		return runDeriveL2(args)
	case "derive-l3":
		return runDeriveL3(args)
	case "validate":
		return runValidate(args)
	case "sync-links":
		return runSyncLinks(args)
	case "cascade":
		return runCascade(args)
	case "help", "--help", "-h":
		printUsage()
		return domain.ExitCodeSuccess
	case "version", "--version", "-v":
		fmt.Printf("loom-cli %s\n", Version)
		return domain.ExitCodeSuccess
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", command)
		printUsage()
		return domain.ExitCodeError
	}
}

func printUsage() {
	fmt.Println(`loom-cli - AI-assisted specification derivation tool

Usage:
  loom-cli <command> [options]

Commands:
  analyze      Analyze L0 documents to discover domain model and ambiguities
  interview    Conduct structured interview to resolve ambiguities
  derive       Derive L1 (Strategic Design) documents from analysis
  derive-l2    Derive L2 (Tactical Design) documents from L1
  derive-l3    Derive L3 (Operational Design) documents from L2
  validate     Validate derived documents for consistency
  sync-links   Fix missing bidirectional references
  cascade      Run entire derivation pipeline (L0→L1→L2→L3)
  help         Show this help message
  version      Show version information

Run 'loom-cli <command> --help' for more information on a command.`)
}
