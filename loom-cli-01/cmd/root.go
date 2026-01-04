// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md (CLI interface)
// Implements: l2/tech-specs.md TS-ARCH-001 (Command router architecture)
// See: l0/decisions.md DEC-L1-001 (stdlib only, no cobra)
package cmd

import (
	"flag"
	"fmt"
	"os"

	"loom-cli/internal/domain"
)

// Version is set at build time.
var Version = "dev"

// Global flags
var (
	verbose bool
)

// Execute runs the CLI and returns the exit code.
// Implements: l2/tech-specs.md TS-ARCH-001
func Execute() int {
	// Parse global flags first
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")

	// Help and version flags
	showHelp := flag.Bool("help", false, "Show help")
	showHelpShort := flag.Bool("h", false, "Show help (shorthand)")
	showVersion := flag.Bool("version", false, "Show version")

	// Don't parse flags yet - we need to find the subcommand first
	if len(os.Args) < 2 {
		printUsage()
		return domain.ExitCodeError
	}

	// Check for version/help before subcommand
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-version" {
			fmt.Printf("loom-cli version %s\n", Version)
			return domain.ExitCodeSuccess
		}
		if arg == "--help" || arg == "-help" || arg == "-h" {
			printUsage()
			return domain.ExitCodeSuccess
		}
	}

	// Get subcommand
	subcommand := os.Args[1]

	// Check if first arg is a flag (not a subcommand)
	if subcommand[0] == '-' {
		// Parse flags for root command
		flag.Parse()
		if *showHelp || *showHelpShort {
			printUsage()
			return domain.ExitCodeSuccess
		}
		if *showVersion {
			fmt.Printf("loom-cli version %s\n", Version)
			return domain.ExitCodeSuccess
		}
		fmt.Fprintf(os.Stderr, "Error: unknown flag: %s\n", subcommand)
		printUsage()
		return domain.ExitCodeError
	}

	// Route to subcommand
	// Implements: l2/interface-contracts.md (command routing)
	switch subcommand {
	case "analyze":
		return runAnalyze(os.Args[2:])
	case "interview":
		return runInterview(os.Args[2:])
	case "derive":
		return runDerive(os.Args[2:])
	case "derive-l2":
		return runDeriveL2(os.Args[2:])
	case "derive-l3":
		return runDeriveL3(os.Args[2:])
	case "validate":
		return runValidate(os.Args[2:])
	case "sync-links":
		return runSyncLinks(os.Args[2:])
	case "cascade":
		return runCascade(os.Args[2:])
	case "help":
		printUsage()
		return domain.ExitCodeSuccess
	case "version":
		fmt.Printf("loom-cli version %s\n", Version)
		return domain.ExitCodeSuccess
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n", subcommand)
		printUsage()
		return domain.ExitCodeError
	}
}

// printUsage prints CLI usage information.
// Implements: l2/interface-contracts.md IC-HLP-001
func printUsage() {
	fmt.Println(`loom-cli - AI-assisted specification derivation tool

Usage:
  loom-cli <command> [flags]

Commands:
  analyze      Extract domain model from L0 specification
  interview    Conduct structured interview to resolve ambiguities
  derive       Generate L1 documents (AC, BR, domain model)
  derive-l2    Generate L2 documents (tech specs, contracts)
  derive-l3    Generate L3 documents (test cases, API spec)
  validate     Validate document consistency and traceability
  sync-links   Synchronize cross-document references
  cascade      Run full derivation pipeline (L0→L1→L2→L3)
  help         Show this help message
  version      Show version information

Global Flags:
  -v, --verbose    Enable verbose output
  -h, --help       Show help
  --version        Show version

Exit Codes:
  0    Success
  1    Error
  100  Interview has more questions (interview command only)

Examples:
  # Analyze L0 specification
  loom-cli analyze --input-file story.md

  # Run interview to resolve ambiguities
  loom-cli interview --init analysis.json
  loom-cli interview --answer '{"question_id":"Q1","answer":"a"}'

  # Full cascade derivation
  loom-cli cascade --input-file story.md --output-dir ./output --skip-interview

  # Validate generated documents
  loom-cli validate --input-dir ./output --level ALL

For more information on each command, run:
  loom-cli <command> --help`)
}

// Verbose returns whether verbose mode is enabled.
func Verbose() bool {
	return verbose
}
