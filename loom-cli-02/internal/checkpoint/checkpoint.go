// Package checkpoint provides resume functionality for long operations.
//
// Implements: l2/package-structure.md PKG-003
// See: l2/internal-api.md
package checkpoint

import (
	"encoding/json"
	"os"
	"time"
)

// Checkpoint holds state for resumable operations.
//
// Implements: l2/internal-api.md
type Checkpoint struct {
	Phase     string    `json:"phase"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

// Save persists checkpoint to file.
//
// Implements: l2/internal-api.md
func Save(path string, cp *Checkpoint) error {
	cp.Timestamp = time.Now()

	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Load reads checkpoint from file.
//
// Implements: l2/internal-api.md
func Load(path string) (*Checkpoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, err
	}

	return &cp, nil
}

// Exists checks if checkpoint file exists.
//
// Implements: l2/internal-api.md
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Clear removes checkpoint file.
//
// Implements: l2/internal-api.md
func Clear(path string) error {
	if !Exists(path) {
		return nil
	}
	return os.Remove(path)
}
