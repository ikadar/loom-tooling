package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds all CLI configuration
type Config struct {
	InputFile     string
	InputDir      string
	OutputDir     string
	DecisionsFile string
	AnalysisFile  string // For derive command
	Format        string // "text" or "json"
	BatchMode     bool   // Non-interactive mode
	Verbose       bool
}

// ParseArgsForAnalyze parses arguments for the analyze command
func ParseArgsForAnalyze(args []string) (*Config, error) {
	cfg := &Config{
		Format: "json",
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--input-file":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--input-file requires a value")
			}
			i++
			cfg.InputFile = args[i]

		case "--input-dir":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--input-dir requires a value")
			}
			i++
			cfg.InputDir = args[i]

		case "--decisions":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--decisions requires a value")
			}
			i++
			cfg.DecisionsFile = args[i]

		case "--verbose", "-v":
			cfg.Verbose = true
		}
	}

	// Validate
	if cfg.InputFile == "" && cfg.InputDir == "" {
		return nil, fmt.Errorf("either --input-file or --input-dir is required")
	}

	if cfg.InputFile != "" && cfg.InputDir != "" {
		return nil, fmt.Errorf("cannot specify both --input-file and --input-dir")
	}

	// Set default decisions file
	if cfg.DecisionsFile == "" {
		if cfg.InputDir != "" {
			cfg.DecisionsFile = filepath.Join(cfg.InputDir, "decisions.md")
		} else {
			cfg.DecisionsFile = filepath.Join(filepath.Dir(cfg.InputFile), "decisions.md")
		}
	}

	return cfg, nil
}

// ParseArgsForDerive parses arguments for the derive command
func ParseArgsForDerive(args []string) (*Config, error) {
	cfg := &Config{
		Format: "text",
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--output-dir":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--output-dir requires a value")
			}
			i++
			cfg.OutputDir = args[i]

		case "--decisions":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--decisions requires a value")
			}
			i++
			cfg.DecisionsFile = args[i]

		case "--analysis-file":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--analysis-file requires a value")
			}
			i++
			cfg.AnalysisFile = args[i]

		case "--verbose", "-v":
			cfg.Verbose = true
		}
	}

	// Validate
	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("--output-dir is required")
	}

	// Set default decisions file if not specified
	if cfg.DecisionsFile == "" {
		cfg.DecisionsFile = filepath.Join(cfg.OutputDir, "decisions.md")
	}

	return cfg, nil
}

// ParseArgs is the legacy function - now redirects to ParseArgsForDerive
func ParseArgs(args []string) (*Config, error) {
	return ParseArgsForDerive(args)
}

// ReadInputFiles reads all L0 input files
func (cfg *Config) ReadInputFiles() (string, []string, error) {
	var files []string
	var contents []string

	if cfg.InputFile != "" {
		// Single file mode
		content, err := os.ReadFile(cfg.InputFile)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read input file: %w", err)
		}
		files = append(files, cfg.InputFile)
		contents = append(contents, string(content))
	} else {
		// Directory mode
		entries, err := os.ReadDir(cfg.InputDir)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read input directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			// Skip decisions.md
			if entry.Name() == "decisions.md" {
				continue
			}

			path := filepath.Join(cfg.InputDir, entry.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				return "", nil, fmt.Errorf("failed to read %s: %w", path, err)
			}

			files = append(files, path)
			contents = append(contents, fmt.Sprintf("<!-- SOURCE: %s -->\n%s", entry.Name(), string(content)))
		}
	}

	if len(files) == 0 {
		return "", nil, fmt.Errorf("no markdown files found in input")
	}

	// Combine all contents
	combined := strings.Join(contents, "\n\n---\n\n")
	return combined, files, nil
}
