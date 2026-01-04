// Package cmd provides CLI commands for loom-cli.
//
// This file implements the cascade command.
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
// Implements: IC-CAS-001
//
// Pipeline:
//  1. Analyze
//  2. Interview (if not --skip-interview)
//  3. Derive L1
//  4. Derive L2
//  5. Derive L3
//
// State: .cascade-state.json in output directory
func runCascade(args []string) int {
	fs := flag.NewFlagSet("cascade", flag.ExitOnError)

	inputFile := fs.String("input-file", "", "L0 specification file")
	inputDir := fs.String("input-dir", "", "L0 specification directory")
	outputDir := fs.String("output-dir", "./output", "Output directory")
	skipInterview := fs.Bool("skip-interview", false, "Use default answers for all questions")
	interactive := fs.Bool("interactive", false, "Enable interactive approval mode")
	resume := fs.Bool("resume", false, "Resume from last saved state")
	fromLevel := fs.String("from", "", "Resume from specific level: analyze, interview, l1, l2, l3")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	if *inputFile == "" && *inputDir == "" && !*resume {
		fmt.Fprintf(os.Stderr, "Error: either --input-file, --input-dir, or --resume is required\n")
		return domain.ExitCodeError
	}

	// Create output directory structure
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		return domain.ExitCodeError
	}

	l1Dir := filepath.Join(*outputDir, "l1")
	l2Dir := filepath.Join(*outputDir, "l2")
	l3Dir := filepath.Join(*outputDir, "l3")

	os.MkdirAll(l1Dir, 0755)
	os.MkdirAll(l2Dir, 0755)
	os.MkdirAll(l3Dir, 0755)

	statePath := filepath.Join(*outputDir, ".cascade-state.json")

	// Load or create state
	var state *domain.CascadeState
	if *resume {
		var err error
		state, err = loadCascadeState(statePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
			return domain.ExitCodeError
		}
	} else {
		// Calculate input hash
		inputContent, _, err := readInputContent(*inputFile, *inputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return domain.ExitCodeError
		}

		hash := sha256.Sum256([]byte(inputContent))
		inputHash := hex.EncodeToString(hash[:8])

		state = &domain.CascadeState{
			Version:   "1.0",
			InputHash: inputHash,
			Phases: map[string]*domain.PhaseState{
				"analyze":   {Status: "pending"},
				"interview": {Status: "pending"},
				"l1":        {Status: "pending"},
				"l2":        {Status: "pending"},
				"l3":        {Status: "pending"},
			},
			Config: domain.CascadeStateConfig{
				SkipInterview: *skipInterview,
				Interactive:   *interactive,
			},
		}
		state.Timestamps.Started = time.Now()
	}

	// Determine starting point
	startPhase := "analyze"
	if *fromLevel != "" {
		startPhase = *fromLevel
	} else if *resume {
		// Find first non-completed phase
		phases := []string{"analyze", "interview", "l1", "l2", "l3"}
		for _, p := range phases {
			if state.Phases[p].Status != "completed" {
				startPhase = p
				break
			}
		}
	}

	analysisFile := filepath.Join(*outputDir, "analysis.json")
	interviewFile := filepath.Join(*outputDir, "interview-state.json")

	// Phase 1: Analyze
	if shouldRunPhase(startPhase, "analyze") {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 1: Analyze...")
		state.Phases["analyze"].Status = "running"
		state.Phases["analyze"].Timestamp = time.Now()
		saveCascadeState(state, statePath)

		analyzeArgs := []string{}
		if *inputFile != "" {
			analyzeArgs = append(analyzeArgs, "--input-file", *inputFile)
		}
		if *inputDir != "" {
			analyzeArgs = append(analyzeArgs, "--input-dir", *inputDir)
		}
		analyzeArgs = append(analyzeArgs, "--output", analysisFile)

		if code := runAnalyze(analyzeArgs); code != domain.ExitCodeSuccess {
			state.Phases["analyze"].Status = "failed"
			state.Phases["analyze"].Error = "analyze command failed"
			saveCascadeState(state, statePath)
			return code
		}

		state.Phases["analyze"].Status = "completed"
		saveCascadeState(state, statePath)
	}

	// Phase 2: Interview
	if shouldRunPhase(startPhase, "interview") && !*skipInterview {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 2: Interview...")
		state.Phases["interview"].Status = "running"
		state.Phases["interview"].Timestamp = time.Now()
		saveCascadeState(state, statePath)

		// Initialize interview
		initArgs := []string{"--init", analysisFile, "--state", interviewFile}
		if code := runInterview(initArgs); code == domain.ExitCodeError {
			state.Phases["interview"].Status = "failed"
			saveCascadeState(state, statePath)
			return code
		}

		// Auto-answer with defaults if skip-interview
		// (This block won't run since we check !*skipInterview above)
		state.Phases["interview"].Status = "completed"
		saveCascadeState(state, statePath)
	} else if *skipInterview {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 2: Interview (skipped - using defaults)...")
		state.Phases["interview"].Status = "completed"
		saveCascadeState(state, statePath)
	}

	// Determine interview output file
	deriveInput := analysisFile
	if !*skipInterview {
		deriveInput = interviewFile
	}

	// Phase 3: Derive L1
	if shouldRunPhase(startPhase, "l1") {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 3: Derive L1...")
		state.Phases["l1"].Status = "running"
		state.Phases["l1"].Timestamp = time.Now()
		saveCascadeState(state, statePath)

		deriveArgs := []string{"--input-file", deriveInput, "--output-dir", l1Dir}
		if *interactive {
			deriveArgs = append(deriveArgs, "--interactive")
		}

		if code := runDerive(deriveArgs); code != domain.ExitCodeSuccess {
			state.Phases["l1"].Status = "failed"
			saveCascadeState(state, statePath)
			return code
		}

		state.Phases["l1"].Status = "completed"
		saveCascadeState(state, statePath)
	}

	// Phase 4: Derive L2
	if shouldRunPhase(startPhase, "l2") {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 4: Derive L2...")
		state.Phases["l2"].Status = "running"
		state.Phases["l2"].Timestamp = time.Now()
		saveCascadeState(state, statePath)

		deriveL2Args := []string{"--input-dir", l1Dir, "--output-dir", l2Dir}
		if *interactive {
			deriveL2Args = append(deriveL2Args, "--interactive")
		}

		if code := runDeriveL2(deriveL2Args); code != domain.ExitCodeSuccess {
			state.Phases["l2"].Status = "failed"
			saveCascadeState(state, statePath)
			return code
		}

		state.Phases["l2"].Status = "completed"
		saveCascadeState(state, statePath)
	}

	// Phase 5: Derive L3
	if shouldRunPhase(startPhase, "l3") {
		fmt.Fprintln(os.Stderr, "[cascade] Phase 5: Derive L3...")
		state.Phases["l3"].Status = "running"
		state.Phases["l3"].Timestamp = time.Now()
		saveCascadeState(state, statePath)

		deriveL3Args := []string{"--input-dir", l2Dir, "--l1-dir", l1Dir, "--output-dir", l3Dir}
		if *interactive {
			deriveL3Args = append(deriveL3Args, "--interactive")
		}

		if code := runDeriveL3(deriveL3Args); code != domain.ExitCodeSuccess {
			state.Phases["l3"].Status = "failed"
			saveCascadeState(state, statePath)
			return code
		}

		state.Phases["l3"].Status = "completed"
		saveCascadeState(state, statePath)
	}

	// Mark complete
	state.Timestamps.Completed = time.Now()
	saveCascadeState(state, statePath)

	fmt.Fprintln(os.Stderr, "[cascade] Cascade complete!")
	fmt.Fprintf(os.Stderr, "Output directory: %s\n", *outputDir)
	fmt.Fprintf(os.Stderr, "  L1: %s\n", l1Dir)
	fmt.Fprintf(os.Stderr, "  L2: %s\n", l2Dir)
	fmt.Fprintf(os.Stderr, "  L3: %s\n", l3Dir)

	return domain.ExitCodeSuccess
}

// shouldRunPhase determines if a phase should run based on starting point.
func shouldRunPhase(startPhase, currentPhase string) bool {
	phases := []string{"analyze", "interview", "l1", "l2", "l3"}

	startIdx := -1
	currentIdx := -1
	for i, p := range phases {
		if p == startPhase {
			startIdx = i
		}
		if p == currentPhase {
			currentIdx = i
		}
	}

	return currentIdx >= startIdx
}

// loadCascadeState loads state from file.
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

// saveCascadeState saves state to file.
func saveCascadeState(state *domain.CascadeState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
