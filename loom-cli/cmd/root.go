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
	case "analyze":
		return runAnalyze()
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
  loom-cli analyze [options]     # L0 → ambiguities
  loom-cli interview [options]   # ambiguities → decisions
  loom-cli derive [options]      # L0+decisions → L1 (Strategic Design)
  loom-cli derive-l2 [options]   # L1 → L2 (Tactical Design)
  loom-cli derive-l3 [options]   # L2 → L3 (Operational Design)
  loom-cli validate [options]    # Validate generated documents
  loom-cli version
  loom-cli help

Commands:
  analyze    Analyze L0 input, discover domain model, find ambiguities
  interview  Conduct iterative structured interview (one question at a time)
  derive     Derive L1 Strategic Design (Domain Model, Bounded Contexts, AC, BR)
  derive-l2  Derive L2 Tactical Design (Test Cases, Tech Specs, Aggregates, Sequences)
  derive-l3  Derive L3 Operational Design (API Spec, Services, Events, Dependencies)
  validate   Validate documents (structure, traceability, completeness, TDAI)
  version    Show version information
  help       Show this help message

Derivation Flow:
  L0 (User Stories) → analyze → interview → derive → L1 (Strategic Design)
  L1 (Strategic) → derive-l2 → L2 (Tactical Design)
  L2 (Tactical) → derive-l3 → L3 (Operational Design)

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

Derive-L2 Options (L1 → L2):
  --input-dir <path>      Directory containing L1 docs (acceptance-criteria.md, business-rules.md)
  --output-dir <path>     Directory for generated L2 documents (required)
  --interactive, -i       Interactive approval mode (preview each file before writing)

Derive-L3 Options (L2 → L3):
  --input-dir <path>      Directory containing L2 docs (test-cases.md, tech-specs.md)
  --output-dir <path>     Directory for generated L3 documents (required)

Validate Options:
  --input-dir <path>      Directory containing documents to validate (required)
  --level <L1|L2|L3|ALL>  Validation level (default: ALL)
  --json                  Output results as JSON

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

Full Flow Example:
  1. loom-cli analyze --input-file story.md > analysis.json
  2. loom-cli interview --init analysis.json --state state.json
  3. [Answer questions until exit code 0]
  4. loom-cli derive --output-dir ./l1 --analysis-file state.json
  5. loom-cli derive-l2 --input-dir ./l1 --output-dir ./l2
  6. loom-cli derive-l3 --input-dir ./l2 --output-dir ./l3`)
}
