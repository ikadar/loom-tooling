package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint represents a saved state of L2 derivation
type Checkpoint struct {
	Version   string                 `json:"version"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	InputDir  string                 `json:"input_dir"`
	OutputDir string                 `json:"output_dir"`
	Phases    map[string]PhaseState  `json:"phases"`
	Data      map[string]interface{} `json:"data"`
}

// PhaseState tracks the completion state of a phase
type PhaseState struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"` // "pending", "running", "completed", "failed"
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"

	CheckpointVersion = "1.0"
)

// Manager handles checkpoint operations
type Manager struct {
	checkpoint *Checkpoint
	filePath   string
}

// NewManager creates a new checkpoint manager
func NewManager(outputDir string) *Manager {
	checkpointFile := filepath.Join(outputDir, ".loom-checkpoint.json")
	return &Manager{
		filePath: checkpointFile,
		checkpoint: &Checkpoint{
			Version:   CheckpointVersion,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			OutputDir: outputDir,
			Phases:    make(map[string]PhaseState),
			Data:      make(map[string]interface{}),
		},
	}
}

// SetInputDir sets the input directory in checkpoint
func (m *Manager) SetInputDir(inputDir string) {
	m.checkpoint.InputDir = inputDir
}

// Load attempts to load an existing checkpoint
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No checkpoint exists, start fresh
		}
		return fmt.Errorf("failed to read checkpoint: %w", err)
	}

	if err := json.Unmarshal(data, m.checkpoint); err != nil {
		return fmt.Errorf("failed to parse checkpoint: %w", err)
	}

	return nil
}

// Save persists the current checkpoint state
func (m *Manager) Save() error {
	m.checkpoint.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(m.checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	if err := os.WriteFile(m.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}

	return nil
}

// Delete removes the checkpoint file
func (m *Manager) Delete() error {
	if err := os.Remove(m.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete checkpoint: %w", err)
	}
	return nil
}

// StartPhase marks a phase as running
func (m *Manager) StartPhase(name string) {
	m.checkpoint.Phases[name] = PhaseState{
		Name:      name,
		Status:    StatusRunning,
		StartedAt: time.Now(),
	}
	m.Save() // Auto-save on state change
}

// CompletePhase marks a phase as completed and stores its data
func (m *Manager) CompletePhase(name string, data interface{}) {
	state := m.checkpoint.Phases[name]
	state.Status = StatusCompleted
	state.CompletedAt = time.Now()
	m.checkpoint.Phases[name] = state
	m.checkpoint.Data[name] = data
	m.Save() // Auto-save on state change
}

// FailPhase marks a phase as failed
func (m *Manager) FailPhase(name string, err error) {
	state := m.checkpoint.Phases[name]
	state.Status = StatusFailed
	state.CompletedAt = time.Now()
	state.Error = err.Error()
	m.checkpoint.Phases[name] = state
	m.Save() // Auto-save on state change
}

// IsPhaseCompleted checks if a phase was already completed
func (m *Manager) IsPhaseCompleted(name string) bool {
	if state, ok := m.checkpoint.Phases[name]; ok {
		return state.Status == StatusCompleted
	}
	return false
}

// GetPhaseData retrieves stored data for a completed phase
func (m *Manager) GetPhaseData(name string) (interface{}, bool) {
	data, ok := m.checkpoint.Data[name]
	return data, ok
}

// GetCompletedPhases returns list of completed phase names
func (m *Manager) GetCompletedPhases() []string {
	var completed []string
	for name, state := range m.checkpoint.Phases {
		if state.Status == StatusCompleted {
			completed = append(completed, name)
		}
	}
	return completed
}

// HasCheckpoint returns true if a checkpoint file exists
func (m *Manager) HasCheckpoint() bool {
	_, err := os.Stat(m.filePath)
	return err == nil
}

// GetFilePath returns the checkpoint file path
func (m *Manager) GetFilePath() string {
	return m.filePath
}

// PrintStatus prints the current checkpoint status
func (m *Manager) PrintStatus() {
	fmt.Fprintf(os.Stderr, "\nCheckpoint Status:\n")
	fmt.Fprintf(os.Stderr, "  File: %s\n", m.filePath)
	fmt.Fprintf(os.Stderr, "  Created: %s\n", m.checkpoint.CreatedAt.Format(time.RFC3339))

	if len(m.checkpoint.Phases) == 0 {
		fmt.Fprintf(os.Stderr, "  Phases: (none)\n")
		return
	}

	fmt.Fprintf(os.Stderr, "  Phases:\n")
	for name, state := range m.checkpoint.Phases {
		status := state.Status
		switch status {
		case StatusCompleted:
			status = "✓ " + status
		case StatusFailed:
			status = "✗ " + status
		case StatusRunning:
			status = "⟳ " + status
		default:
			status = "○ " + status
		}
		fmt.Fprintf(os.Stderr, "    - %s: %s\n", name, status)
	}
}
