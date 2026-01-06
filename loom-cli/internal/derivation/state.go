package derivation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// =============================================================================
// Constants
// =============================================================================

const (
	// StateFileName is the name of the derivation state file
	StateFileName = "derivation-state.json"

	// LockFileName is the name of the lock file
	LockFileName = "derivation-state.lock"

	// LoomDirName is the name of the .loom directory
	LoomDirName = ".loom"

	// StateVersion is the current state file format version
	StateVersion = "1.0"

	// LockTimeout is how long to wait for a lock
	LockTimeout = 30 * time.Second

	// LockStaleTimeout is how long before a lock is considered stale
	LockStaleTimeout = 5 * time.Minute
)

// =============================================================================
// DerivationState
// =============================================================================

// DerivationState holds the complete state of a derivation project
type DerivationState struct {
	// Version is the state file format version
	Version string `json:"version"`

	// Project is the project name/identifier
	Project string `json:"project"`

	// LastFullDerive is when a full derivation was last run
	LastFullDerive time.Time `json:"last_full_derive,omitempty"`

	// LoomVersion is the loom-cli version that last modified this state
	LoomVersion string `json:"loom_version"`

	// Artifacts maps artifact IDs to their state
	Artifacts map[string]*Artifact `json:"artifacts"`

	// Decisions maps decision IDs to their state
	Decisions map[string]*Decision `json:"decisions"`

	// DependencyGraph holds the dependency relationships
	DependencyGraph *DependencyGraph `json:"dependency_graph"`

	// FileHashes caches file content hashes for incremental updates
	FileHashes map[string]*FileHashInfo `json:"file_hashes,omitempty"`

	// mu protects concurrent access to the state
	mu sync.RWMutex `json:"-"`
}

// FileHashInfo caches hash information for a file
type FileHashInfo struct {
	Path    string    `json:"path"`
	Hash    string    `json:"hash"`
	ModTime time.Time `json:"mod_time"`
	Size    int64     `json:"size"`
}

// =============================================================================
// State Manager
// =============================================================================

// StateManager handles loading, saving, and locking derivation state
type StateManager struct {
	// ProjectDir is the root directory of the project
	ProjectDir string

	// LoomDir is the path to the .loom directory
	LoomDir string

	// StatePath is the path to the state file
	StatePath string

	// LockPath is the path to the lock file
	LockPath string

	// lockFile holds the lock file handle when locked
	lockFile *os.File

	// mu protects the lock state
	mu sync.Mutex
}

// NewStateManager creates a new state manager for a project directory
func NewStateManager(projectDir string) *StateManager {
	loomDir := filepath.Join(projectDir, LoomDirName)
	return &StateManager{
		ProjectDir: projectDir,
		LoomDir:    loomDir,
		StatePath:  filepath.Join(loomDir, StateFileName),
		LockPath:   filepath.Join(loomDir, LockFileName),
	}
}

// =============================================================================
// State Loading
// =============================================================================

// Load loads the derivation state from disk
// Returns a new empty state if the file doesn't exist
func (sm *StateManager) Load() (*DerivationState, error) {
	// Check if state file exists
	if _, err := os.Stat(sm.StatePath); os.IsNotExist(err) {
		return sm.NewState(), nil
	}

	// Read the file
	data, err := os.ReadFile(sm.StatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	// Parse JSON
	var state DerivationState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	// Migrate if needed
	if state.Version != StateVersion {
		if err := sm.migrateState(&state); err != nil {
			return nil, fmt.Errorf("failed to migrate state: %w", err)
		}
	}

	// Initialize maps if nil
	if state.Artifacts == nil {
		state.Artifacts = make(map[string]*Artifact)
	}
	if state.Decisions == nil {
		state.Decisions = make(map[string]*Decision)
	}
	if state.FileHashes == nil {
		state.FileHashes = make(map[string]*FileHashInfo)
	}
	if state.DependencyGraph == nil {
		state.DependencyGraph = NewDependencyGraph()
	}

	return &state, nil
}

// NewState creates a new empty derivation state
func (sm *StateManager) NewState() *DerivationState {
	return &DerivationState{
		Version:         StateVersion,
		Project:         filepath.Base(sm.ProjectDir),
		LoomVersion:     "0.3.0", // TODO: Get from version constant
		Artifacts:       make(map[string]*Artifact),
		Decisions:       make(map[string]*Decision),
		DependencyGraph: NewDependencyGraph(),
		FileHashes:      make(map[string]*FileHashInfo),
	}
}

// migrateState handles migration from older state file versions
func (sm *StateManager) migrateState(state *DerivationState) error {
	// Currently only version 1.0, no migration needed
	// Future versions will add migration logic here
	state.Version = StateVersion
	return nil
}

// =============================================================================
// State Saving
// =============================================================================

// Save saves the derivation state to disk
func (sm *StateManager) Save(state *DerivationState) error {
	// Ensure .loom directory exists
	if err := os.MkdirAll(sm.LoomDir, 0755); err != nil {
		return fmt.Errorf("failed to create .loom directory: %w", err)
	}

	// Update version
	state.Version = StateVersion

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write atomically using temp file + rename
	tmpPath := sm.StatePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp state file: %w", err)
	}

	if err := os.Rename(tmpPath, sm.StatePath); err != nil {
		os.Remove(tmpPath) // Clean up temp file
		return fmt.Errorf("failed to rename state file: %w", err)
	}

	return nil
}

// =============================================================================
// File Locking
// =============================================================================

// Lock acquires an exclusive lock on the state file
// This prevents concurrent modifications from multiple processes
func (sm *StateManager) Lock() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.lockFile != nil {
		return nil // Already locked by this process
	}

	// Ensure .loom directory exists
	if err := os.MkdirAll(sm.LoomDir, 0755); err != nil {
		return fmt.Errorf("failed to create .loom directory: %w", err)
	}

	// Check for stale lock
	if info, err := os.Stat(sm.LockPath); err == nil {
		if time.Since(info.ModTime()) > LockStaleTimeout {
			// Lock is stale, remove it
			os.Remove(sm.LockPath)
		}
	}

	// Try to create lock file exclusively
	deadline := time.Now().Add(LockTimeout)
	for {
		f, err := os.OpenFile(sm.LockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			// Got the lock
			sm.lockFile = f

			// Write lock info
			lockInfo := struct {
				PID       int       `json:"pid"`
				Hostname  string    `json:"hostname"`
				LockedAt  time.Time `json:"locked_at"`
			}{
				PID:      os.Getpid(),
				Hostname: hostname(),
				LockedAt: time.Now(),
			}
			json.NewEncoder(f).Encode(lockInfo)

			return nil
		}

		if !os.IsExist(err) {
			return fmt.Errorf("failed to create lock file: %w", err)
		}

		// Lock exists, check if we've timed out
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for lock (another process may be running)")
		}

		// Wait and retry
		time.Sleep(100 * time.Millisecond)
	}
}

// Unlock releases the lock on the state file
func (sm *StateManager) Unlock() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.lockFile == nil {
		return nil // Not locked
	}

	// Close and remove lock file
	sm.lockFile.Close()
	sm.lockFile = nil

	if err := os.Remove(sm.LockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}

	return nil
}

// IsLocked checks if the state is currently locked
func (sm *StateManager) IsLocked() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.lockFile != nil
}

// hostname returns the current hostname or "unknown"
func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

// =============================================================================
// State Operations
// =============================================================================

// GetArtifact returns an artifact by ID (thread-safe)
func (s *DerivationState) GetArtifact(id string) *Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Artifacts[id]
}

// SetArtifact sets an artifact (thread-safe)
func (s *DerivationState) SetArtifact(artifact *Artifact) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Artifacts[artifact.ID] = artifact
}

// RemoveArtifact removes an artifact by ID (thread-safe)
func (s *DerivationState) RemoveArtifact(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Artifacts, id)
}

// GetDecision returns a decision by ID (thread-safe)
func (s *DerivationState) GetDecision(id string) *Decision {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Decisions[id]
}

// SetDecision sets a decision (thread-safe)
func (s *DerivationState) SetDecision(decision *Decision) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Decisions[decision.ID] = decision
}

// GetArtifactsByLayer returns all artifacts in a layer
func (s *DerivationState) GetArtifactsByLayer(layer string) []*Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Artifact
	for _, a := range s.Artifacts {
		if a.Layer == layer {
			result = append(result, a)
		}
	}
	return result
}

// GetArtifactsByStatus returns all artifacts with a given status
func (s *DerivationState) GetArtifactsByStatus(status ArtifactStatus) []*Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Artifact
	for _, a := range s.Artifacts {
		if a.Status == status {
			result = append(result, a)
		}
	}
	return result
}

// GetStaleArtifacts returns all stale and affected artifacts
func (s *DerivationState) GetStaleArtifacts() []*Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Artifact
	for _, a := range s.Artifacts {
		if a.Status == StatusStale || a.Status == StatusAffected {
			result = append(result, a)
		}
	}
	return result
}

// GetFileHash returns cached hash info for a file
func (s *DerivationState) GetFileHash(path string) *FileHashInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.FileHashes[path]
}

// SetFileHash sets cached hash info for a file
func (s *DerivationState) SetFileHash(info *FileHashInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FileHashes[info.Path] = info
}

// =============================================================================
// State Statistics
// =============================================================================

// GetStatusReport generates a status report for the current state
func (s *DerivationState) GetStatusReport() *StatusReport {
	s.mu.RLock()
	defer s.mu.RUnlock()

	report := &StatusReport{
		Project:        s.Project,
		LastFullDerive: s.LastFullDerive,
		LayerSummary:   make(map[string]*LayerStatus),
		TotalArtifacts: len(s.Artifacts),
	}

	// Initialize layer summaries
	for _, layer := range []string{"l0", "l1", "l2", "l3"} {
		report.LayerSummary[layer] = &LayerStatus{}
	}

	// Count artifacts by layer and status
	for _, a := range s.Artifacts {
		ls := report.LayerSummary[a.Layer]
		if ls == nil {
			ls = &LayerStatus{}
			report.LayerSummary[a.Layer] = ls
		}

		switch a.Status {
		case StatusCurrent:
			ls.Current++
		case StatusStale:
			ls.Stale++
			report.StaleArtifacts = append(report.StaleArtifacts, ArtifactSummary{
				ID:       a.ID,
				Type:     a.Type,
				Location: a.Location.File,
				Status:   a.Status,
			})
		case StatusAffected:
			ls.Affected++
			report.AffectedArtifacts = append(report.AffectedArtifacts, ArtifactSummary{
				ID:       a.ID,
				Type:     a.Type,
				Location: a.Location.File,
				Status:   a.Status,
			})
		case StatusModified:
			ls.Modified++
			report.ModifiedArtifacts = append(report.ModifiedArtifacts, ArtifactSummary{
				ID:       a.ID,
				Type:     a.Type,
				Location: a.Location.File,
				Status:   a.Status,
				Message:  fmt.Sprintf("Manual sections: %v", a.ManualSections),
			})
		case StatusNew:
			ls.New++
		case StatusOrphaned:
			ls.Orphaned++
		}
	}

	return report
}
