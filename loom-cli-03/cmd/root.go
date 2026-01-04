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

// Version is the CLI version
const Version = "0.3.0"

// Execute is the main entry point for the CLI.
// Returns exit code: 0=success, 1=error, 100=interview question available.
//
// Implements: l2/interface-contracts.md (Exit codes)
func Execute() int {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return domain.ExitCodeSuccess
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "analyze":
		return runAnalyze(commandArgs)
	case "interview":
		return runInterview(commandArgs)
	case "derive":
		return runDerive(commandArgs)
	case "derive-l2":
		return runDeriveL2(commandArgs)
	case "derive-l3":
		return runDeriveL3(commandArgs)
	case "validate":
		return runValidate(commandArgs)
	case "sync-links":
		return runSyncLinks(commandArgs)
	case "cascade":
		return runCascade(commandArgs)
	case "help", "--help", "-h":
		printUsage()
		return domain.ExitCodeSuccess
	case "version":
		fmt.Printf("loom-cli v%s\n", Version)
		return domain.ExitCodeSuccess
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		return domain.ExitCodeError
	}
}

// printUsage prints the help message.
//
// Implements: IC-HLP-001
func printUsage() {
	fmt.Println(`loom-cli - AI-assisted specification derivation tool

Usage:
  loom-cli <command> [options]

Commands:
  analyze       Analyze user stories to discover domain model (IC-ANL-001)
  interview     Conduct structured interview to resolve ambiguities (IC-INT-001)
  derive        Derive L1 (Strategic Design) documents (IC-DRV-001)
  derive-l2     Derive L2 (Tactical Design) documents (IC-DRV-002)
  derive-l3     Derive L3 (Operational Design) documents (IC-DRV-003)
  validate      Validate derived documents (IC-VAL-001)
  sync-links    Fix missing bidirectional references (IC-SYN-001)
  cascade       Run entire derivation pipeline (IC-CAS-001)
  version       Show version information (IC-VER-001)
  help          Show this help message (IC-HLP-001)

Exit Codes:
  0    Success
  1    Error
  100  Interview: question available

For command-specific help:
  loom-cli <command> --help`)
}

// Command implementations are in separate files:
// - analyze.go: runAnalyze()
// - interview.go: runInterview()
// - derive_new.go: runDerive()
// - derive_l2.go: runDeriveL2()
// - derive_l3.go: runDeriveL3()
// - validate.go: runValidate()
// - sync_links.go: runSyncLinks()
// - cascade.go: runCascade()
