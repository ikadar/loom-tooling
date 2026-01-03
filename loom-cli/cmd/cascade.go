package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CascadeConfig holds configuration for the cascade command
type CascadeConfig struct {
	InputFile     string
	InputDir      string
	OutputDir     string
	SkipInterview bool
	DecisionsFile string
	Interactive   bool
	Resume        bool
	FromLevel     string
}

// CascadeState tracks the progress of cascade derivation
type CascadeState struct {
	Version   string                  `json:"version"`
	InputHash string                  `json:"input_hash"`
	Phases    map[string]*PhaseState  `json:"phases"`
	Config    CascadeStateConfig      `json:"config"`
	Timestamps struct {
		Started   time.Time `json:"started"`
		Completed time.Time `json:"completed,omitempty"`
	} `json:"timestamps"`
}

type PhaseState struct {
	Status    string    `json:"status"` // pending, running, completed, failed
	Timestamp time.Time `json:"timestamp,omitempty"`
	Error     string    `json:"error,omitempty"`
}

type CascadeStateConfig struct {
	SkipInterview bool `json:"skip_interview"`
	Interactive   bool `json:"interactive"`
}

const cascadeStateFile = ".cascade-state.json"

func runCascade() error {
	cfg, err := parseCascadeArgs()
	if err != nil {
		return err
	}

	// Validate config
	if cfg.InputFile == "" && cfg.InputDir == "" {
		return fmt.Errorf("either --input-file or --input-dir is required")
	}
	if cfg.OutputDir == "" {
		return fmt.Errorf("--output-dir is required")
	}

	// Create output directories
	l1Dir := filepath.Join(cfg.OutputDir, "l1")
	l2Dir := filepath.Join(cfg.OutputDir, "l2")
	l3Dir := filepath.Join(cfg.OutputDir, "l3")

	for _, dir := range []string{cfg.OutputDir, l1Dir, l2Dir, l3Dir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Load or create state
	state, err := loadCascadeState(cfg)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load state: %w", err)
	}
	if state == nil {
		state = newCascadeState(cfg)
	}

	// Check if we should skip to a specific level
	if cfg.FromLevel != "" {
		resetFromLevel(state, cfg.FromLevel)
	}

	fmt.Fprintf(os.Stderr, "╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(os.Stderr, "║               LOOM CASCADE DERIVATION                        ║\n")
	fmt.Fprintf(os.Stderr, "╚══════════════════════════════════════════════════════════════╝\n\n")

	// Phase 1: Analyze
	if shouldRunPhase(state, "analyze", cfg) {
		fmt.Fprintf(os.Stderr, "━━━ Phase 1/5: Analyze ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		if err := runCascadeAnalyze(cfg, state); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "━━━ Phase 1/5: Analyze [SKIPPED - already completed] ━━━━━━━━━\n")
	}

	// Phase 2: Interview (optional)
	if !cfg.SkipInterview && shouldRunPhase(state, "interview", cfg) {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 2/5: Interview ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		if err := runCascadeInterview(cfg, state); err != nil {
			return err
		}
	} else if cfg.SkipInterview {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 2/5: Interview [SKIPPED - using AI defaults] ━━━━━━━\n")
		state.Phases["interview"].Status = "completed"
		state.Phases["interview"].Timestamp = time.Now()
	} else {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 2/5: Interview [SKIPPED - already completed] ━━━━━━\n")
	}

	// Phase 3: Derive L1
	if shouldRunPhase(state, "derive-l1", cfg) {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 3/5: Derive L1 (Strategic Design) ━━━━━━━━━━━━━━━━━━\n")
		if err := runCascadeDeriveL1(cfg, state); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 3/5: Derive L1 [SKIPPED - already completed] ━━━━━━━\n")
	}

	// Phase 4: Derive L2
	if shouldRunPhase(state, "derive-l2", cfg) {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 4/5: Derive L2 (Tactical Design) ━━━━━━━━━━━━━━━━━━━\n")
		if err := runCascadeDeriveL2(cfg, state); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 4/5: Derive L2 [SKIPPED - already completed] ━━━━━━━\n")
	}

	// Phase 5: Derive L3
	if shouldRunPhase(state, "derive-l3", cfg) {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 5/5: Derive L3 (Operational Design) ━━━━━━━━━━━━━━━━\n")
		if err := runCascadeDeriveL3(cfg, state); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(os.Stderr, "\n━━━ Phase 5/5: Derive L3 [SKIPPED - already completed] ━━━━━━━\n")
	}

	// Mark complete
	state.Timestamps.Completed = time.Now()
	if err := saveCascadeState(cfg.OutputDir, state); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save final state: %v\n", err)
	}

	// Print summary
	fmt.Fprintf(os.Stderr, "\n╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(os.Stderr, "║               CASCADE COMPLETE                               ║\n")
	fmt.Fprintf(os.Stderr, "╚══════════════════════════════════════════════════════════════╝\n")
	fmt.Fprintf(os.Stderr, "\nOutput directories:\n")
	fmt.Fprintf(os.Stderr, "  L1 (Strategic): %s\n", l1Dir)
	fmt.Fprintf(os.Stderr, "  L2 (Tactical):  %s\n", l2Dir)
	fmt.Fprintf(os.Stderr, "  L3 (Operational): %s\n", l3Dir)

	return nil
}

func parseCascadeArgs() (*CascadeConfig, error) {
	cfg := &CascadeConfig{}
	args := os.Args[2:] // Skip "loom-cli" and "cascade"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input-file":
			if i+1 < len(args) {
				cfg.InputFile = args[i+1]
				i++
			}
		case "--input-dir":
			if i+1 < len(args) {
				cfg.InputDir = args[i+1]
				i++
			}
		case "--output-dir":
			if i+1 < len(args) {
				cfg.OutputDir = args[i+1]
				i++
			}
		case "--skip-interview":
			cfg.SkipInterview = true
		case "--decisions":
			if i+1 < len(args) {
				cfg.DecisionsFile = args[i+1]
				i++
			}
		case "--interactive", "-i":
			cfg.Interactive = true
		case "--resume", "-r":
			cfg.Resume = true
		case "--from":
			if i+1 < len(args) {
				cfg.FromLevel = args[i+1]
				i++
			}
		}
	}

	return cfg, nil
}

func newCascadeState(cfg *CascadeConfig) *CascadeState {
	return &CascadeState{
		Version:   "1.0",
		InputHash: computeInputHash(cfg),
		Phases: map[string]*PhaseState{
			"analyze":    {Status: "pending"},
			"interview":  {Status: "pending"},
			"derive-l1":  {Status: "pending"},
			"derive-l2":  {Status: "pending"},
			"derive-l3":  {Status: "pending"},
		},
		Config: CascadeStateConfig{
			SkipInterview: cfg.SkipInterview,
			Interactive:   cfg.Interactive,
		},
		Timestamps: struct {
			Started   time.Time `json:"started"`
			Completed time.Time `json:"completed,omitempty"`
		}{
			Started: time.Now(),
		},
	}
}

func computeInputHash(cfg *CascadeConfig) string {
	h := sha256.New()
	if cfg.InputFile != "" {
		if content, err := os.ReadFile(cfg.InputFile); err == nil {
			h.Write(content)
		}
	}
	if cfg.InputDir != "" {
		h.Write([]byte(cfg.InputDir))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func loadCascadeState(cfg *CascadeConfig) (*CascadeState, error) {
	if !cfg.Resume {
		return nil, os.ErrNotExist
	}

	path := filepath.Join(cfg.OutputDir, cascadeStateFile)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state CascadeState
	if err := json.Unmarshal(content, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func saveCascadeState(outputDir string, state *CascadeState) error {
	path := filepath.Join(outputDir, cascadeStateFile)
	content, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func shouldRunPhase(state *CascadeState, phase string, cfg *CascadeConfig) bool {
	if !cfg.Resume {
		return true
	}
	ps := state.Phases[phase]
	return ps == nil || ps.Status != "completed"
}

func resetFromLevel(state *CascadeState, level string) {
	levels := []string{"l1", "l2", "l3"}
	phases := []string{"derive-l1", "derive-l2", "derive-l3"}

	startIdx := -1
	for i, l := range levels {
		if l == level {
			startIdx = i
			break
		}
	}

	if startIdx >= 0 {
		for i := startIdx; i < len(phases); i++ {
			state.Phases[phases[i]].Status = "pending"
		}
	}
}

// Phase runners

func runCascadeAnalyze(cfg *CascadeConfig, state *CascadeState) error {
	state.Phases["analyze"].Status = "running"
	saveCascadeState(cfg.OutputDir, state)

	// Build args for analyze command
	analyzeArgs := []string{"loom-cli", "analyze"}
	if cfg.InputFile != "" {
		analyzeArgs = append(analyzeArgs, "--input-file", cfg.InputFile)
	}
	if cfg.InputDir != "" {
		analyzeArgs = append(analyzeArgs, "--input-dir", cfg.InputDir)
	}
	if cfg.DecisionsFile != "" {
		analyzeArgs = append(analyzeArgs, "--decisions", cfg.DecisionsFile)
	}

	// Save original args and restore after
	origArgs := os.Args
	os.Args = analyzeArgs

	// Redirect stdout to capture analysis output
	analysisFile := filepath.Join(cfg.OutputDir, ".analysis.json")
	origStdout := os.Stdout
	f, err := os.Create(analysisFile)
	if err != nil {
		os.Args = origArgs
		state.Phases["analyze"].Status = "failed"
		state.Phases["analyze"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("failed to create analysis file: %w", err)
	}
	os.Stdout = f

	err = runAnalyze()

	// Restore stdout and args
	os.Stdout = origStdout
	f.Close()
	os.Args = origArgs

	if err != nil {
		state.Phases["analyze"].Status = "failed"
		state.Phases["analyze"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("analyze failed: %w", err)
	}

	state.Phases["analyze"].Status = "completed"
	state.Phases["analyze"].Timestamp = time.Now()
	saveCascadeState(cfg.OutputDir, state)
	return nil
}

func runCascadeInterview(cfg *CascadeConfig, state *CascadeState) error {
	state.Phases["interview"].Status = "running"
	saveCascadeState(cfg.OutputDir, state)

	analysisFile := filepath.Join(cfg.OutputDir, ".analysis.json")
	stateFile := filepath.Join(cfg.OutputDir, ".interview-state.json")

	// Initialize interview
	origArgs := os.Args
	os.Args = []string{"loom-cli", "interview", "--init", analysisFile, "--state", stateFile}

	err := runInterview()
	os.Args = origArgs

	if err != nil {
		state.Phases["interview"].Status = "failed"
		state.Phases["interview"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("interview init failed: %w", err)
	}

	// For cascade with skip-interview, we just use AI defaults
	// The interview state file now contains the decisions
	// TODO: For interactive mode, loop through questions

	state.Phases["interview"].Status = "completed"
	state.Phases["interview"].Timestamp = time.Now()
	saveCascadeState(cfg.OutputDir, state)
	return nil
}

func runCascadeDeriveL1(cfg *CascadeConfig, state *CascadeState) error {
	state.Phases["derive-l1"].Status = "running"
	saveCascadeState(cfg.OutputDir, state)

	l1Dir := filepath.Join(cfg.OutputDir, "l1")

	// Determine which state file to use
	stateFile := filepath.Join(cfg.OutputDir, ".interview-state.json")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		// Fall back to analysis file if no interview
		stateFile = filepath.Join(cfg.OutputDir, ".analysis.json")
	}

	origArgs := os.Args
	os.Args = []string{"loom-cli", "derive", "--output-dir", l1Dir, "--analysis-file", stateFile}

	if cfg.DecisionsFile != "" {
		os.Args = append(os.Args, "--decisions", cfg.DecisionsFile)
	}

	err := runDeriveNew()
	os.Args = origArgs

	if err != nil {
		state.Phases["derive-l1"].Status = "failed"
		state.Phases["derive-l1"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("derive L1 failed: %w", err)
	}

	state.Phases["derive-l1"].Status = "completed"
	state.Phases["derive-l1"].Timestamp = time.Now()
	saveCascadeState(cfg.OutputDir, state)
	return nil
}

func runCascadeDeriveL2(cfg *CascadeConfig, state *CascadeState) error {
	state.Phases["derive-l2"].Status = "running"
	saveCascadeState(cfg.OutputDir, state)

	l1Dir := filepath.Join(cfg.OutputDir, "l1")
	l2Dir := filepath.Join(cfg.OutputDir, "l2")

	origArgs := os.Args
	args := []string{"loom-cli", "derive-l2", "--input-dir", l1Dir, "--output-dir", l2Dir}
	if cfg.Interactive {
		args = append(args, "--interactive")
	}
	os.Args = args

	err := runDeriveL2()
	os.Args = origArgs

	if err != nil {
		state.Phases["derive-l2"].Status = "failed"
		state.Phases["derive-l2"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("derive L2 failed: %w", err)
	}

	state.Phases["derive-l2"].Status = "completed"
	state.Phases["derive-l2"].Timestamp = time.Now()
	saveCascadeState(cfg.OutputDir, state)
	return nil
}

func runCascadeDeriveL3(cfg *CascadeConfig, state *CascadeState) error {
	state.Phases["derive-l3"].Status = "running"
	saveCascadeState(cfg.OutputDir, state)

	l2Dir := filepath.Join(cfg.OutputDir, "l2")
	l3Dir := filepath.Join(cfg.OutputDir, "l3")

	origArgs := os.Args
	os.Args = []string{"loom-cli", "derive-l3", "--input-dir", l2Dir, "--output-dir", l3Dir}

	err := runDeriveL3()
	os.Args = origArgs

	if err != nil {
		state.Phases["derive-l3"].Status = "failed"
		state.Phases["derive-l3"].Error = err.Error()
		saveCascadeState(cfg.OutputDir, state)
		return fmt.Errorf("derive L3 failed: %w", err)
	}

	state.Phases["derive-l3"].Status = "completed"
	state.Phases["derive-l3"].Timestamp = time.Now()
	saveCascadeState(cfg.OutputDir, state)
	return nil
}
