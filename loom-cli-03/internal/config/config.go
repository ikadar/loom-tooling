// Package config provides configuration management.
//
// Implements: l2/package-structure.md PKG-004
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds CLI configuration.
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
func ParseArgsForAnalyze(args []string) (*Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)

	fs.StringVar(&cfg.InputFile, "input-file", "", "Path to input markdown file")
	fs.StringVar(&cfg.InputDir, "input-dir", "", "Path to directory containing markdown files")
	fs.StringVar(&cfg.VocabularyFile, "vocabulary", "", "Path to domain vocabulary file")
	fs.StringVar(&cfg.NFRFile, "nfr", "", "Path to non-functional requirements file")
	fs.StringVar(&cfg.Format, "format", "json", "Output format: json or text")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.InputFile == "" && cfg.InputDir == "" {
		return nil, fmt.Errorf("either --input-file or --input-dir is required")
	}

	return cfg, nil
}

// ParseArgsForDerive parses derive command arguments.
func ParseArgsForDerive(args []string) (*Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("derive", flag.ContinueOnError)

	fs.StringVar(&cfg.InputFile, "input-file", "", "Path to input markdown file")
	fs.StringVar(&cfg.InputDir, "input-dir", "", "Path to directory containing markdown files")
	fs.StringVar(&cfg.OutputDir, "output-dir", "", "Output directory for generated files")
	fs.StringVar(&cfg.DecisionsFile, "decisions", "", "Path to decisions file")
	fs.StringVar(&cfg.AnalysisFile, "analysis", "", "Path to analysis JSON file")
	fs.StringVar(&cfg.VocabularyFile, "vocabulary", "", "Path to domain vocabulary file")
	fs.StringVar(&cfg.NFRFile, "nfr", "", "Path to non-functional requirements file")
	fs.BoolVar(&cfg.BatchMode, "batch", false, "Batch mode (no interactive prompts)")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.OutputDir == "" {
		return nil, fmt.Errorf("--output-dir is required")
	}

	return cfg, nil
}

// ReadInputFiles reads all L0 input markdown files.
// Returns: combined content, list of file paths, error
func (cfg *Config) ReadInputFiles() (string, []string, error) {
	var files []string
	var contents []string

	if cfg.InputFile != "" {
		data, err := os.ReadFile(cfg.InputFile)
		if err != nil {
			return "", nil, err
		}
		files = append(files, cfg.InputFile)
		contents = append(contents, string(data))
	}

	if cfg.InputDir != "" {
		entries, err := os.ReadDir(cfg.InputDir)
		if err != nil {
			return "", nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}

			path := filepath.Join(cfg.InputDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				return "", nil, err
			}
			files = append(files, path)
			contents = append(contents, string(data))
		}
	}

	if len(contents) == 0 {
		return "", nil, fmt.Errorf("no input files found")
	}

	return strings.Join(contents, "\n\n---\n\n"), files, nil
}

// ReadVocabulary reads optional domain vocabulary file.
func (cfg *Config) ReadVocabulary() (string, error) {
	if cfg.VocabularyFile == "" {
		return "", nil
	}

	data, err := os.ReadFile(cfg.VocabularyFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadNFR reads optional non-functional requirements file.
func (cfg *Config) ReadNFR() (string, error) {
	if cfg.NFRFile == "" {
		return "", nil
	}

	data, err := os.ReadFile(cfg.NFRFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
