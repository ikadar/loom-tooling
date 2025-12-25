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
  loom-cli analyze [options]     # Phases 0-3: analyze and find ambiguities
  loom-cli interview [options]   # Phase 4: iterative structured interview
  loom-cli derive [options]      # Phases 5-6: derive AC/BR from decisions
  loom-cli version
  loom-cli help

Commands:
  analyze    Analyze L0 input, discover domain model, find ambiguities
  interview  Conduct iterative structured interview (one question at a time)
  derive     Derive L1 documents from domain model + decisions
  version    Show version information
  help       Show this help message

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

Derive Options:
  --output-dir <path>     Directory for generated L1 documents (required)
  --decisions <path>      Path to decisions.md (to append new decisions)
  --analysis-file <path>  Path to analysis JSON or interview state

Flow:
  1. loom-cli analyze --input-file story.md > analysis.json
  2. loom-cli interview --init analysis.json --state state.json  # exits 100
  3. [Claude presents question, user answers]
  4. loom-cli interview --state state.json --answer '{"question_id":"...","answer":"..."}'
  5. Repeat 3-4 until exit code 0
  6. loom-cli derive --output-dir ./out --analysis-file state.json

Examples:
  # Initialize interview
  loom-cli interview --init /tmp/analysis.json --state /tmp/interview.json

  # Answer a question
  loom-cli interview --state /tmp/interview.json \
    --answer '{"question_id":"AMB-ENT-001","answer":"Yes, add Draft status","source":"user"}'`)
}
