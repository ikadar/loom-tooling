// Package config provides configuration management.
//
// Implements: l2/package-structure.md PKG-004
// See: l2/internal-api.md
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds CLI configuration.
//
// Implements: l2/internal-api.md
type Config struct {
	InputFile      string
	InputDir       string
	OutputDir      string
	DecisionsFile  string
	AnalysisFile   string
	VocabularyFile string
	NFRFile        string
	Format         string // "text" or "json"
	BatchMode      bool
	Verbose        bool
}

// ParseArgsForAnalyze parses analyze command arguments.
//
// Requires: --input-file OR --input-dir
// Optional: --decisions, --verbose
//
// Implements: l2/internal-api.md
func ParseArgsForAnalyze(args []string) (*Config, error) {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	cfg := &Config{}

	fs.StringVar(&cfg.InputFile, "input-file", "", "Input file path")
	fs.StringVar(&cfg.InputDir, "input-dir", "", "Input directory path")
	fs.StringVar(&cfg.DecisionsFile, "decisions", "", "Existing decisions file")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.InputFile == "" && cfg.InputDir == "" {
		return nil, fmt.Errorf("either --input-file or --input-dir is required")
	}

	return cfg, nil
}

// ParseArgsForDerive parses derive command arguments.
//
// Requires: --output-dir
// Optional: --decisions, --analysis-file, --vocabulary, --nfr, --verbose
//
// Implements: l2/internal-api.md
func ParseArgsForDerive(args []string) (*Config, error) {
	fs := flag.NewFlagSet("derive", flag.ContinueOnError)
	cfg := &Config{}

	fs.StringVar(&cfg.OutputDir, "output-dir", "", "Output directory path (required)")
	fs.StringVar(&cfg.DecisionsFile, "decisions", "", "Decisions file")
	fs.StringVar(&cfg.AnalysisFile, "analysis-file", "", "Analysis JSON file")
	fs.StringVar(&cfg.VocabularyFile, "vocabulary", "", "Domain vocabulary file")
	fs.StringVar(&cfg.NFRFile, "nfr", "", "Non-functional requirements file")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("--output-dir is required")
	}

	return cfg, nil
}

// ParseArgs is legacy alias for ParseArgsForDerive.
//
// Implements: l2/internal-api.md
func ParseArgs(args []string) (*Config, error) {
	return ParseArgsForDerive(args)
}

// ReadInputFiles reads all L0 input markdown files.
// Returns: combined content, list of file paths, error
//
// Implements: l2/internal-api.md
func (cfg *Config) ReadInputFiles() (string, []string, error) {
	var files []string
	var contents []string

	if cfg.InputFile != "" {
		content, err := os.ReadFile(cfg.InputFile)
		if err != nil {
			return "", nil, fmt.Errorf("failed to read input file: %w", err)
		}
		files = append(files, cfg.InputFile)
		contents = append(contents, string(content))
	}

	if cfg.InputDir != "" {
		err := filepath.Walk(cfg.InputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".md") {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				files = append(files, path)
				contents = append(contents, string(content))
			}
			return nil
		})
		if err != nil {
			return "", nil, fmt.Errorf("failed to read input directory: %w", err)
		}
	}

	if len(files) == 0 {
		return "", nil, fmt.Errorf("no input files found")
	}

	return strings.Join(contents, "\n\n---\n\n"), files, nil
}

// ReadVocabulary reads optional domain vocabulary file.
//
// Implements: l2/internal-api.md
func (cfg *Config) ReadVocabulary() (string, error) {
	if cfg.VocabularyFile == "" {
		return "", nil
	}

	content, err := os.ReadFile(cfg.VocabularyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read vocabulary file: %w", err)
	}

	return string(content), nil
}

// ReadNFR reads optional non-functional requirements file.
//
// Implements: l2/internal-api.md
func (cfg *Config) ReadNFR() (string, error) {
	if cfg.NFRFile == "" {
		return "", nil
	}

	content, err := os.ReadFile(cfg.NFRFile)
	if err != nil {
		return "", fmt.Errorf("failed to read NFR file: %w", err)
	}

	return string(content), nil
}
