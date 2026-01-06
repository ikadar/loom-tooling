package cmd

import (
	"fmt"
	"os"
)

// Version information
const Version = "0.3.0"

// Execute runs the CLI
func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]

	switch command {
	case "init":
		return runInit()
	case "analyze":
		return runAnalyze()
	case "analyze-v2":
		return runAnalyzeV2()
	case "interview":
		return runInterview()
	case "derive":
		return runDeriveNew()
	case "derive-l2":
		return runDeriveL2()
	case "derive-l3":
		return runDeriveL3()
	case "validate":
		return runValidate()
	case "sync-links":
		return runSyncLinks()
	case "cascade":
		return runCascade()
	case "status":
		return runStatus()
	case "rederive":
		return runRederive()
	case "migrate":
		return runMigrate()
	case "version":
		fmt.Printf("loom-cli v%s\n", Version)
		return nil
	case "help", "--help", "-h":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func printUsage() {
	fmt.Println(`loom-cli - AI-DOP Documentation Derivation CLI

Usage:
  loom-cli init [options]        # Initialize project state tracking
  loom-cli cascade [options]     # Full derivation: L0 → L1 → L2 → L3 (one command!)
  loom-cli analyze [options]     # L0 → ambiguities
  loom-cli interview [options]   # ambiguities → decisions
  loom-cli derive [options]      # L0+decisions → L1 (Strategic Design)
  loom-cli derive-l2 [options]   # L1 → L2 (Tactical Design)
  loom-cli derive-l3 [options]   # L2 → L3 (Operational Design)
  loom-cli status [options]      # Show derivation status (stale artifacts)
  loom-cli rederive [options]    # Re-derive stale artifacts
  loom-cli migrate [options]     # Migrate existing project to LOOM format
  loom-cli validate [options]    # Validate generated documents
  loom-cli sync-links [options]  # Fix missing bidirectional links
  loom-cli version
  loom-cli help

Commands:
  init       Initialize derivation state tracking (.loom directory)
  cascade    Full cascade derivation: L0 → L1 → L2 → L3 in one command
  analyze    Analyze L0 input, discover domain model, find ambiguities
  interview  Conduct iterative structured interview (one question at a time)
  derive     Derive L1 Strategic Design (Domain Model, Bounded Contexts, AC, BR)
  derive-l2  Derive L2 Tactical Design (Tech Specs, Contracts, Aggregates, Sequences)
  derive-l3  Derive L3 Operational Design (Test Cases, API Spec, Skeletons, Events)
  status     Show derivation status and stale artifacts
  rederive   Re-derive stale artifacts (update from upstream changes)
  migrate    Migrate existing project to LOOM-marked format
  validate   Validate documents (structure, traceability, completeness, TDAI)
  sync-links Add missing bidirectional references between documents
  version    Show version information
  help       Show this help message

Derivation Flow:
  L0 (User Stories) → analyze → interview → derive → L1 (Strategic Design)
  L1 (Strategic) → derive-l2 → L2 (Tactical Design)
  L2 (Tactical) → derive-l3 → L3 (Operational Design)

Init Options:
  --project-dir <path>  Project root directory (default: current directory)
  --force               Overwrite existing state file
  --scan                Scan existing documents and build initial state

Cascade Options (Full Derivation):
  --input-file <path>     L0 input file (user story)
  --input-dir <path>      L0 input directory
  --output-dir <path>     Base output directory (creates l1/, l2/, l3/)
  --skip-interview        Skip interview, use AI defaults
  --decisions <path>      Use existing decisions.md
  --interactive, -i       Interactive approval mode
  --resume                Resume from previous state
  --from <level>          Re-derive from level (l1, l2, l3)

Analyze Options:
  --input-file <path>     Path to single L0 input file
  --input-dir <path>      Path to directory with L0 files
  --decisions <path>      Path to existing decisions.md

Interview Options:
  --init <path>           Initialize interview from analysis JSON
  --state <path>          Path to interview state file
  --answer <json>         JSON with answer: {"question_id":"...", "answer":"...", "source":"user"}

  Exit codes:
    0   = Interview complete, no more questions
    1   = Error
    100 = Question available (output contains question JSON)

Derive Options (L0 → L1):
  --output-dir <path>     Directory for generated L1 documents (required)
  --decisions <path>      Path to decisions.md (to append new decisions)
  --analysis-file <path>  Path to analysis JSON or interview state
  --vocabulary <path>     Optional domain vocabulary file (enhances domain model)
  --nfr <path>            Optional non-functional requirements file (adds to BRs/ACs)

Derive-L2 Options (L1 → L2):
  --input-dir <path>      Directory containing L1 docs (acceptance-criteria.md, business-rules.md)
  --output-dir <path>     Directory for generated L2 documents (required)
  --interactive, -i       Interactive approval mode (preview each file before writing)

Derive-L3 Options (L2 → L3):
  --input-dir <path>      Directory containing L2 docs (test-cases.md, tech-specs.md)
  --output-dir <path>     Directory for generated L3 documents (required)

Status Options:
  --project-dir <path>    Project root directory (default: current directory)
  --verbose               Show detailed artifact information
  --layer <l0|l1|l2|l3>   Filter by layer
  --format <text|json>    Output format (default: text)
  --plan                  Show derivation plan for stale artifacts

Rederive Options:
  --project-dir <path>    Project root directory (default: current directory)
  --layer <l1|l2|l3>      Only derive artifacts in this layer
  --all                   Derive all stale artifacts
  --dry-run               Preview without making changes
  --verbose               Show detailed output
  --preserve-manual       Keep manual sections during re-derivation (default: true)
  --interactive           Confirm each derivation
  <artifact-ids>          Specific artifact IDs to derive (positional args)

Migrate Options:
  --project-dir <path>    Project root directory (default: current directory)
  --dry-run               Preview without making changes
  --backup-dir <path>     Directory for backups (default: .loom/backups)
  --verbose               Show detailed output
  --force                 Force migration even if markers exist

Validate Options:
  --input-dir <path>      Directory containing documents to validate (required)
  --level <L1|L2|L3|ALL>  Validation level (default: ALL)
  --json                  Output results as JSON

Sync-Links Options:
  --input-dir <path>      Directory containing documents to sync (required)
  --dry-run               Show what would be changed without modifying files

Validation Rules:
  V001  Every document has IDs
  V002  IDs follow expected patterns (AC-XXX-NNN, BR-XXX-NNN, etc.)
  V003  All references point to existing IDs
  V004  Bidirectional links are consistent
  V005  Every AC has at least 1 test case
  V006  Every Entity has an aggregate
  V007  Every Service has an interface contract
  V008  Negative test ratio >= 20% (TDAI)
  V009  Every AC has hallucination prevention test (TDAI)
  V010  No duplicate IDs

Cascade Example (Recommended):
  loom-cli cascade --input-file story.md --output-dir ./specs --skip-interview

Full Flow Example (Manual):
  1. loom-cli analyze --input-file story.md > analysis.json
  2. loom-cli interview --init analysis.json --state state.json
  3. [Answer questions until exit code 0]
  4. loom-cli derive --output-dir ./l1 --analysis-file state.json
  5. loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2
  6. loom-cli derive-l3 --input-dir ./l2 --output-dir ./l3`)
}
