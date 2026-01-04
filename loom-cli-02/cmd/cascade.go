// Implements: l2/interface-contracts.md IC-CAS-001
// See: l2/sequence-design.md SEQ-CAS-001, SEQ-CAS-002
// See: l2/aggregate-design.md AGG-CAS-001
package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"loom-cli/internal/domain"
)

// runCascade implements the cascade command.
//
// Implements: IC-CAS-001
// Runs the entire derivation pipeline: analyze → interview → derive L1 → derive L2 → derive L3
func runCascade(args []string) int {
	fs := flag.NewFlagSet("cascade", flag.ContinueOnError)
	inputFile := fs.String("input-file", "", "L0 input file")
	inputDir := fs.String("input-dir", "", "L0 input directory")
	outputDir := fs.String("output-dir", "", "Output directory (required)")
	skipInterview := fs.Bool("skip-interview", false, "Skip interview, use AI-suggested defaults")
	decisionsFile := fs.String("decisions", "", "Use existing decisions file")
	interactive := fs.Bool("interactive", false, "Interactive approval at each level")
	fs.Bool("i", false, "Alias for --interactive")
	resume := fs.Bool("resume", false, "Resume from interrupted state")
	fs.Bool("r", false, "Alias for --resume")
	fromLevel := fs.String("from", "", "Re-derive from specific level (l1, l2, l3)")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate required options
	if *inputFile == "" && *inputDir == "" && !*resume {
		fmt.Fprintln(os.Stderr, "Error: either --input-file or --input-dir is required")
		return domain.ExitCodeError
	}
	if *outputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --output-dir is required")
		return domain.ExitCodeError
	}

	// Create output directories
	l1Dir := filepath.Join(*outputDir, "l1")
	l2Dir := filepath.Join(*outputDir, "l2")
	l3Dir := filepath.Join(*outputDir, "l3")

	for _, dir := range []string{*outputDir, l1Dir, l2Dir, l3Dir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create directory %s: %v\n", dir, err)
			return domain.ExitCodeError
		}
	}

	stateFile := filepath.Join(*outputDir, ".cascade-state.json")

	// Load or create state
	var state *domain.CascadeState
	var err error

	if *resume || *fromLevel != "" {
		state, err = loadCascadeState(stateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to load state: %v\n", err)
			return domain.ExitCodeError
		}

		// Handle --from flag
		if *fromLevel != "" {
			resetFromLevel(state, *fromLevel)
		}
	} else {
		// Compute input hash
		var inputContent string
		if *inputFile != "" {
			data, err := os.ReadFile(*inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to read input file: %v\n", err)
				return domain.ExitCodeError
			}
			inputContent = string(data)
		} else if *inputDir != "" {
			inputContent, _, err = readInputDir(*inputDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to read input directory: %v\n", err)
				return domain.ExitCodeError
			}
		}

		state = newCascadeState(inputContent, *skipInterview, *interactive)
	}

	// Save state after each phase
	defer saveCascadeState(state, stateFile)

	// Phase 1: Analyze
	if state.Phases["analyze"].Status != "completed" {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 1: Analyze...")
		state.Phases["analyze"].Status = "running"
		state.Phases["analyze"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)

		analyzeArgs := []string{}
		if *inputFile != "" {
			analyzeArgs = append(analyzeArgs, "--input-file", *inputFile)
		} else if *inputDir != "" {
			analyzeArgs = append(analyzeArgs, "--input-dir", *inputDir)
		}
		if *decisionsFile != "" {
			analyzeArgs = append(analyzeArgs, "--decisions", *decisionsFile)
		}
		analyzeArgs = append(analyzeArgs, "--output", filepath.Join(*outputDir, ".analysis.json"))
		if *verbose {
			analyzeArgs = append(analyzeArgs, "--verbose")
		}

		if code := runAnalyze(analyzeArgs); code != domain.ExitCodeSuccess {
			state.Phases["analyze"].Status = "failed"
			state.Phases["analyze"].Error = "analyze failed"
			saveCascadeState(state, stateFile)
			fmt.Fprintln(os.Stderr, "Error: analyze failed")
			return domain.ExitCodeError
		}

		state.Phases["analyze"].Status = "completed"
		state.Phases["analyze"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)
	}

	// Phase 2: Interview (if not skipped)
	if !state.Config.SkipInterview && state.Phases["interview"].Status != "completed" {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 2: Interview...")
		state.Phases["interview"].Status = "running"
		state.Phases["interview"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)

		// Initialize interview
		interviewArgs := []string{
			"--init", filepath.Join(*outputDir, ".analysis.json"),
			"--state", filepath.Join(*outputDir, ".interview-state.json"),
		}

		if code := runInterview(interviewArgs); code == domain.ExitCodeError {
			state.Phases["interview"].Status = "failed"
			state.Phases["interview"].Error = "interview failed"
			saveCascadeState(state, stateFile)
			fmt.Fprintln(os.Stderr, "Error: interview failed")
			return domain.ExitCodeError
		}

		// In cascade with --skip-interview, auto-answer with defaults
		if state.Config.SkipInterview {
			// Skip all questions with defaults
			for {
				skipArgs := []string{
					"--state", filepath.Join(*outputDir, ".interview-state.json"),
					"--skip",
				}
				code := runInterview(skipArgs)
				if code == domain.ExitCodeSuccess {
					break
				}
				if code == domain.ExitCodeError {
					state.Phases["interview"].Status = "failed"
					saveCascadeState(state, stateFile)
					return domain.ExitCodeError
				}
			}
		}

		state.Phases["interview"].Status = "completed"
		state.Phases["interview"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)
	} else if state.Config.SkipInterview {
		state.Phases["interview"].Status = "completed"
		state.Phases["interview"].Timestamp = time.Now()
	}

	// Phase 3: Derive L1
	if state.Phases["derive-l1"].Status != "completed" {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 3: Derive L1...")
		state.Phases["derive-l1"].Status = "running"
		state.Phases["derive-l1"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)

		deriveArgs := []string{
			"--output-dir", l1Dir,
		}

		// Use interview state if available, otherwise analysis
		interviewState := filepath.Join(*outputDir, ".interview-state.json")
		analysisFile := filepath.Join(*outputDir, ".analysis.json")
		if _, err := os.Stat(interviewState); err == nil {
			deriveArgs = append(deriveArgs, "--analysis-file", interviewState)
		} else {
			deriveArgs = append(deriveArgs, "--analysis-file", analysisFile)
		}

		if *verbose {
			deriveArgs = append(deriveArgs, "--verbose")
		}
		if *interactive || state.Config.Interactive {
			deriveArgs = append(deriveArgs, "--interactive")
		}

		if code := runDerive(deriveArgs); code != domain.ExitCodeSuccess {
			state.Phases["derive-l1"].Status = "failed"
			state.Phases["derive-l1"].Error = "derive L1 failed"
			saveCascadeState(state, stateFile)
			fmt.Fprintln(os.Stderr, "Error: derive L1 failed")
			return domain.ExitCodeError
		}

		state.Phases["derive-l1"].Status = "completed"
		state.Phases["derive-l1"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)
	}

	// Phase 4: Derive L2
	if state.Phases["derive-l2"].Status != "completed" {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 4: Derive L2...")
		state.Phases["derive-l2"].Status = "running"
		state.Phases["derive-l2"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)

		deriveL2Args := []string{
			"--input-dir", l1Dir,
			"--output-dir", l2Dir,
		}
		if *verbose {
			deriveL2Args = append(deriveL2Args, "--verbose")
		}
		if *interactive || state.Config.Interactive {
			deriveL2Args = append(deriveL2Args, "--interactive")
		}

		if code := runDeriveL2(deriveL2Args); code != domain.ExitCodeSuccess {
			state.Phases["derive-l2"].Status = "failed"
			state.Phases["derive-l2"].Error = "derive L2 failed"
			saveCascadeState(state, stateFile)
			fmt.Fprintln(os.Stderr, "Error: derive L2 failed")
			return domain.ExitCodeError
		}

		state.Phases["derive-l2"].Status = "completed"
		state.Phases["derive-l2"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)
	}

	// Phase 5: Derive L3
	if state.Phases["derive-l3"].Status != "completed" {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 5: Derive L3...")
		state.Phases["derive-l3"].Status = "running"
		state.Phases["derive-l3"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)

		deriveL3Args := []string{
			"--input-dir", l2Dir,
			"--l1-dir", l1Dir,
			"--output-dir", l3Dir,
		}
		if *verbose {
			deriveL3Args = append(deriveL3Args, "--verbose")
		}

		if code := runDeriveL3(deriveL3Args); code != domain.ExitCodeSuccess {
			state.Phases["derive-l3"].Status = "failed"
			state.Phases["derive-l3"].Error = "derive L3 failed"
			saveCascadeState(state, stateFile)
			fmt.Fprintln(os.Stderr, "Error: derive L3 failed")
			return domain.ExitCodeError
		}

		state.Phases["derive-l3"].Status = "completed"
		state.Phases["derive-l3"].Timestamp = time.Now()
		saveCascadeState(state, stateFile)
	}

	// Mark complete
	state.Timestamps.Completed = time.Now()
	saveCascadeState(state, stateFile)

	fmt.Fprintf(os.Stderr, "\nCascade complete!\n")
	fmt.Fprintf(os.Stderr, "Output: %s\n", *outputDir)
	fmt.Fprintf(os.Stderr, "  - L1: %s\n", l1Dir)
	fmt.Fprintf(os.Stderr, "  - L2: %s\n", l2Dir)
	fmt.Fprintf(os.Stderr, "  - L3: %s\n", l3Dir)

	return domain.ExitCodeSuccess
}

// newCascadeState creates a new cascade state.
//
// Implements: AGG-CAS-001
func newCascadeState(inputContent string, skipInterview, interactive bool) *domain.CascadeState {
	hash := sha256.Sum256([]byte(inputContent))
	inputHash := hex.EncodeToString(hash[:])[:16] // First 16 chars per DEC-L1-010

	state := &domain.CascadeState{
		Version:   "1.0",
		InputHash: inputHash,
		Phases: map[string]*domain.PhaseState{
			"analyze":   {Status: "pending"},
			"interview": {Status: "pending"},
			"derive-l1": {Status: "pending"},
			"derive-l2": {Status: "pending"},
			"derive-l3": {Status: "pending"},
		},
		Config: domain.CascadeStateConfig{
			SkipInterview: skipInterview,
			Interactive:   interactive,
		},
	}
	state.Timestamps.Started = time.Now()

	return state
}

// loadCascadeState loads cascade state from file.
func loadCascadeState(path string) (*domain.CascadeState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state domain.CascadeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// saveCascadeState saves cascade state to file.
func saveCascadeState(state *domain.CascadeState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// resetFromLevel resets phases from a specific level.
func resetFromLevel(state *domain.CascadeState, level string) {
	switch level {
	case "l1":
		state.Phases["derive-l1"].Status = "pending"
		state.Phases["derive-l2"].Status = "pending"
		state.Phases["derive-l3"].Status = "pending"
	case "l2":
		state.Phases["derive-l2"].Status = "pending"
		state.Phases["derive-l3"].Status = "pending"
	case "l3":
		state.Phases["derive-l3"].Status = "pending"
	}
}
