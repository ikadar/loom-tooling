// Package cmd provides CLI commands for loom-cli.
//
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
	"loom-cli/prompts"
)

// Ensure prompts package is imported
var _ = prompts.Derivation

// runCascade handles the cascade command.
//
// Implements: IC-CAS-001
// Phases:
//  1. Analyze → .analysis.json
//  2. Interview (if not --skip-interview) → .interview-state.json
//  3. Derive L1 → l1/
//  4. Derive L2 → l2/
//  5. Derive L3 → l3/
func runCascade(args []string) int {
	fs := flag.NewFlagSet("cascade", flag.ContinueOnError)
	inputFile := fs.String("input-file", "", "L0 input file")
	inputDir := fs.String("input-dir", "", "L0 input directory")
	outputDir := fs.String("output-dir", "", "Output directory (required)")
	skipInterview := fs.Bool("skip-interview", false, "Use AI defaults, skip interview")
	decisionsFile := fs.String("decisions", "", "Existing decisions file")
	interactive := fs.Bool("interactive", false, "Interactive approval mode")
	fs.BoolVar(interactive, "i", false, "Interactive approval mode (shorthand)")
	resume := fs.Bool("resume", false, "Resume from saved state")
	fs.BoolVar(resume, "r", false, "Resume from saved state (shorthand)")
	fromLevel := fs.String("from", "", "Resume from level (l1, l2, l3)")
	verbose := fs.Bool("verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Validate input
	if !*resume && *inputFile == "" && *inputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: either --input-file or --input-dir is required")
		return domain.ExitCodeError
	}
	if *outputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --output-dir is required")
		return domain.ExitCodeError
	}

	// Create output directory structure
	if err := createOutputDirs(*outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Load or create cascade state
	stateFile := filepath.Join(*outputDir, ".cascade-state.json")
	var state *domain.CascadeState
	var err error

	if *resume || *fromLevel != "" {
		state, err = loadCascadeState(stateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to load state: %v\n", err)
			return domain.ExitCodeError
		}
		if *verbose {
			fmt.Fprintln(os.Stderr, "[cascade] Resuming from saved state...")
		}
	} else {
		// Read input to compute hash
		inputContent, _, err := readInputFiles(*inputFile, *inputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return domain.ExitCodeError
		}

		state = newCascadeState(inputContent, *skipInterview, *interactive)
	}

	// Handle --from flag to restart from specific level
	if *fromLevel != "" {
		resetFromLevel(state, *fromLevel)
	}

	// Run phases
	ctx := &cascadeContext{
		inputFile:     *inputFile,
		inputDir:      *inputDir,
		outputDir:     *outputDir,
		decisionsFile: *decisionsFile,
		skipInterview: *skipInterview,
		interactive:   *interactive,
		verbose:       *verbose,
		state:         state,
		stateFile:     stateFile,
	}

	// Phase 1: Analyze
	if state.Phases["analyze"].Status != "completed" {
		if err := ctx.runAnalyzePhase(); err != nil {
			return domain.ExitCodeError
		}
	}

	// Phase 2: Interview
	if !*skipInterview && state.Phases["interview"].Status != "completed" {
		if err := ctx.runInterviewPhase(); err != nil {
			return domain.ExitCodeError
		}
	}

	// Phase 3: Derive L1
	if state.Phases["derive-l1"].Status != "completed" {
		if err := ctx.runDeriveL1Phase(); err != nil {
			return domain.ExitCodeError
		}
	}

	// Phase 4: Derive L2
	if state.Phases["derive-l2"].Status != "completed" {
		if err := ctx.runDeriveL2Phase(); err != nil {
			return domain.ExitCodeError
		}
	}

	// Phase 5: Derive L3
	if state.Phases["derive-l3"].Status != "completed" {
		if err := ctx.runDeriveL3Phase(); err != nil {
			return domain.ExitCodeError
		}
	}

	// Mark complete
	state.Timestamps.Completed = time.Now()
	if err := saveCascadeState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save final state: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "\n[cascade] Complete! Output written to %s\n", *outputDir)
	return domain.ExitCodeSuccess
}

// cascadeContext holds context for cascade execution.
type cascadeContext struct {
	inputFile     string
	inputDir      string
	outputDir     string
	decisionsFile string
	skipInterview bool
	interactive   bool
	verbose       bool
	state         *domain.CascadeState
	stateFile     string
}

// createOutputDirs creates the output directory structure.
func createOutputDirs(outputDir string) error {
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "l1"),
		filepath.Join(outputDir, "l2"),
		filepath.Join(outputDir, "l3"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// newCascadeState creates a new cascade state.
//
// Implements: AGG-CAS-001
func newCascadeState(inputContent string, skipInterview, interactive bool) *domain.CascadeState {
	hash := sha256.Sum256([]byte(inputContent))

	return &domain.CascadeState{
		Version:   "1.0",
		InputHash: hex.EncodeToString(hash[:8]), // First 16 hex chars = 8 bytes
		Phases: map[string]*domain.PhaseState{
			"analyze":    {Status: "pending"},
			"interview":  {Status: "pending"},
			"derive-l1":  {Status: "pending"},
			"derive-l2":  {Status: "pending"},
			"derive-l3":  {Status: "pending"},
		},
		Config: domain.CascadeStateConfig{
			SkipInterview: skipInterview,
			Interactive:   interactive,
		},
		Timestamps: struct {
			Started   time.Time `json:"started"`
			Completed time.Time `json:"completed,omitempty"`
		}{
			Started: time.Now(),
		},
	}
}

// loadCascadeState loads cascade state from JSON file.
func loadCascadeState(path string) (*domain.CascadeState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no checkpoint found in %s, cannot resume", filepath.Dir(path))
	}

	var state domain.CascadeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// saveCascadeState saves cascade state to JSON file.
func saveCascadeState(path string, state *domain.CascadeState) error {
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

// Phase execution methods

func (ctx *cascadeContext) runAnalyzePhase() error {
	if ctx.verbose {
		fmt.Fprintln(os.Stderr, "\n[cascade] Phase 1/5: Analyze")
	}

	ctx.state.Phases["analyze"].Status = "running"
	ctx.state.Phases["analyze"].Timestamp = time.Now()
	saveCascadeState(ctx.stateFile, ctx.state)

	// Build analyze args
	args := []string{}
	if ctx.inputFile != "" {
		args = append(args, "--input-file", ctx.inputFile)
	}
	if ctx.inputDir != "" {
		args = append(args, "--input-dir", ctx.inputDir)
	}
	args = append(args, "--output", filepath.Join(ctx.outputDir, ".analysis.json"))
	if ctx.decisionsFile != "" {
		args = append(args, "--decisions", ctx.decisionsFile)
	}
	if ctx.verbose {
		args = append(args, "--verbose")
	}

	// Run analyze
	if code := runAnalyze(args); code != 0 {
		ctx.state.Phases["analyze"].Status = "failed"
		ctx.state.Phases["analyze"].Error = "analyze failed"
		saveCascadeState(ctx.stateFile, ctx.state)
		fmt.Fprintln(os.Stderr, "Error: analyze failed")
		return fmt.Errorf("analyze failed")
	}

	ctx.state.Phases["analyze"].Status = "completed"
	saveCascadeState(ctx.stateFile, ctx.state)
	return nil
}

func (ctx *cascadeContext) runInterviewPhase() error {
	if ctx.verbose {
		fmt.Fprintln(os.Stderr, "\n[cascade] Phase 2/5: Interview")
	}

	ctx.state.Phases["interview"].Status = "running"
	ctx.state.Phases["interview"].Timestamp = time.Now()
	saveCascadeState(ctx.stateFile, ctx.state)

	analysisFile := filepath.Join(ctx.outputDir, ".analysis.json")
	stateFile := filepath.Join(ctx.outputDir, ".interview-state.json")

	// Initialize interview
	args := []string{
		"--init", analysisFile,
		"--state", stateFile,
	}

	code := runInterview(args)

	if code == domain.ExitCodeQuestion {
		// For non-interactive mode, auto-answer with defaults
		if !ctx.interactive {
			if err := autoAnswerInterview(stateFile, ctx.verbose); err != nil {
				ctx.state.Phases["interview"].Status = "failed"
				ctx.state.Phases["interview"].Error = err.Error()
				saveCascadeState(ctx.stateFile, ctx.state)
				return err
			}
		} else {
			// Interactive mode would require terminal interaction
			// For now, mark as needing manual completion
			fmt.Fprintln(os.Stderr, "Interview requires manual completion. Run:")
			fmt.Fprintf(os.Stderr, "  loom-cli interview --state %s\n", stateFile)
			return fmt.Errorf("interview incomplete")
		}
	} else if code != 0 {
		ctx.state.Phases["interview"].Status = "failed"
		saveCascadeState(ctx.stateFile, ctx.state)
		return fmt.Errorf("interview failed")
	}

	ctx.state.Phases["interview"].Status = "completed"
	saveCascadeState(ctx.stateFile, ctx.state)
	return nil
}

func (ctx *cascadeContext) runDeriveL1Phase() error {
	if ctx.verbose {
		fmt.Fprintln(os.Stderr, "\n[cascade] Phase 3/5: Derive L1")
	}

	ctx.state.Phases["derive-l1"].Status = "running"
	ctx.state.Phases["derive-l1"].Timestamp = time.Now()
	saveCascadeState(ctx.stateFile, ctx.state)

	// Use interview state if available, otherwise analysis
	analysisFile := filepath.Join(ctx.outputDir, ".interview-state.json")
	if _, err := os.Stat(analysisFile); os.IsNotExist(err) {
		analysisFile = filepath.Join(ctx.outputDir, ".analysis.json")
	}

	args := []string{
		"--output-dir", filepath.Join(ctx.outputDir, "l1"),
		"--analysis-file", analysisFile,
	}
	if ctx.decisionsFile != "" {
		args = append(args, "--decisions", ctx.decisionsFile)
	}
	if ctx.verbose {
		args = append(args, "--verbose")
	}

	if code := runDerive(args); code != 0 {
		ctx.state.Phases["derive-l1"].Status = "failed"
		saveCascadeState(ctx.stateFile, ctx.state)
		return fmt.Errorf("derive L1 failed")
	}

	ctx.state.Phases["derive-l1"].Status = "completed"
	saveCascadeState(ctx.stateFile, ctx.state)
	return nil
}

func (ctx *cascadeContext) runDeriveL2Phase() error {
	if ctx.verbose {
		fmt.Fprintln(os.Stderr, "\n[cascade] Phase 4/5: Derive L2")
	}

	ctx.state.Phases["derive-l2"].Status = "running"
	ctx.state.Phases["derive-l2"].Timestamp = time.Now()
	saveCascadeState(ctx.stateFile, ctx.state)

	args := []string{
		"--input-dir", filepath.Join(ctx.outputDir, "l1"),
		"--output-dir", filepath.Join(ctx.outputDir, "l2"),
	}
	if ctx.interactive {
		args = append(args, "--interactive")
	}
	if ctx.verbose {
		args = append(args, "--verbose")
	}

	if code := runDeriveL2(args); code != 0 {
		ctx.state.Phases["derive-l2"].Status = "failed"
		saveCascadeState(ctx.stateFile, ctx.state)
		return fmt.Errorf("derive L2 failed")
	}

	ctx.state.Phases["derive-l2"].Status = "completed"
	saveCascadeState(ctx.stateFile, ctx.state)
	return nil
}

func (ctx *cascadeContext) runDeriveL3Phase() error {
	if ctx.verbose {
		fmt.Fprintln(os.Stderr, "\n[cascade] Phase 5/5: Derive L3")
	}

	ctx.state.Phases["derive-l3"].Status = "running"
	ctx.state.Phases["derive-l3"].Timestamp = time.Now()
	saveCascadeState(ctx.stateFile, ctx.state)

	args := []string{
		"--input-dir", filepath.Join(ctx.outputDir, "l2"),
		"--l1-dir", filepath.Join(ctx.outputDir, "l1"),
		"--output-dir", filepath.Join(ctx.outputDir, "l3"),
	}
	if ctx.verbose {
		args = append(args, "--verbose")
	}

	if code := runDeriveL3(args); code != 0 {
		ctx.state.Phases["derive-l3"].Status = "failed"
		saveCascadeState(ctx.stateFile, ctx.state)
		return fmt.Errorf("derive L3 failed")
	}

	ctx.state.Phases["derive-l3"].Status = "completed"
	saveCascadeState(ctx.stateFile, ctx.state)
	return nil
}

// autoAnswerInterview automatically answers all interview questions with defaults.
func autoAnswerInterview(stateFile string, verbose bool) error {
	state, err := loadInterviewState(stateFile)
	if err != nil {
		return err
	}

	// Answer all remaining questions with suggested answers
	for i := state.CurrentIndex; i < len(state.Questions); i++ {
		q := state.Questions[i]

		// Check if should skip
		if shouldSkipQuestion(&q, state.Decisions) {
			state.Skipped = append(state.Skipped, q.ID)
			continue
		}

		answer := q.SuggestedAnswer
		if answer == "" {
			answer = "[AI-suggested default]"
		}

		decision := domain.Decision{
			ID:        q.ID,
			Question:  q.Question,
			Answer:    answer,
			DecidedAt: time.Now(),
			Source:    "default",
			Category:  q.Category,
			Subject:   q.Subject,
		}
		state.Decisions = append(state.Decisions, decision)

		if verbose {
			fmt.Fprintf(os.Stderr, "  Auto-answered: %s\n", q.ID)
		}
	}

	state.Complete = true
	return saveInterviewState(stateFile, state)
}
