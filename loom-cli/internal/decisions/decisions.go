package decisions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ikadar/loom-cli/internal/domain"
)

// AmbiguityDecision represents a resolved ambiguity
type AmbiguityDecision struct {
	AmbiguityID string    `json:"ambiguity_id"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	Source      string    `json:"source"` // "user", "default", "existing"
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	DecidedAt   time.Time `json:"decided_at"`
}

// DecisionSet holds all decisions
type DecisionSet struct {
	Decisions []AmbiguityDecision `json:"decisions"`
}

// LoadFromFile loads existing decisions from a decisions.md file
func LoadFromFile(path string) (*DecisionSet, error) {
	ds := &DecisionSet{
		Decisions: []AmbiguityDecision{},
	}

	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// No existing file, return empty set
		return ds, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read decisions file: %w", err)
	}

	// Parse decisions from markdown
	// Looking for pattern: ### AMB-XXX-NNN: ...
	// **Question:** ...
	// **Decision:** ...
	// **Source:** ...

	lines := strings.Split(string(content), "\n")
	var currentID, currentQuestion, currentAnswer, currentSource string

	ambIDRegex := regexp.MustCompile(`^###\s+(AMB-[A-Z]+-\d+)`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check for ambiguity header
		if matches := ambIDRegex.FindStringSubmatch(line); matches != nil {
			// Save previous if exists
			if currentID != "" && currentAnswer != "" {
				ds.Decisions = append(ds.Decisions, AmbiguityDecision{
					AmbiguityID: currentID,
					Question:    currentQuestion,
					Answer:      currentAnswer,
					Source:      currentSource,
				})
			}
			currentID = matches[1]
			currentQuestion = ""
			currentAnswer = ""
			currentSource = "existing"
			continue
		}

		// Check for Question line
		if strings.HasPrefix(line, "**Question:**") {
			currentQuestion = strings.TrimSpace(strings.TrimPrefix(line, "**Question:**"))
			continue
		}

		// Check for Decision line
		if strings.HasPrefix(line, "**Decision:**") {
			currentAnswer = strings.TrimSpace(strings.TrimPrefix(line, "**Decision:**"))
			continue
		}

		// Check for Source line
		if strings.HasPrefix(line, "**Source:**") {
			currentSource = strings.TrimSpace(strings.TrimPrefix(line, "**Source:**"))
			continue
		}
	}

	// Save last one
	if currentID != "" && currentAnswer != "" {
		ds.Decisions = append(ds.Decisions, AmbiguityDecision{
			AmbiguityID: currentID,
			Question:    currentQuestion,
			Answer:      currentAnswer,
			Source:      currentSource,
		})
	}

	return ds, nil
}

// HasDecision checks if a decision exists for the given ambiguity ID
func (ds *DecisionSet) HasDecision(ambiguityID string) bool {
	for _, d := range ds.Decisions {
		if d.AmbiguityID == ambiguityID {
			return true
		}
	}
	return false
}

// GetDecision returns the decision for the given ambiguity ID
func (ds *DecisionSet) GetDecision(ambiguityID string) *AmbiguityDecision {
	for _, d := range ds.Decisions {
		if d.AmbiguityID == ambiguityID {
			return &d
		}
	}
	return nil
}

// AddDecision adds a new decision
func (ds *DecisionSet) AddDecision(d AmbiguityDecision) {
	// Update if exists
	for i, existing := range ds.Decisions {
		if existing.AmbiguityID == d.AmbiguityID {
			ds.Decisions[i] = d
			return
		}
	}
	ds.Decisions = append(ds.Decisions, d)
}

// ResolveAmbiguities resolves ambiguities either interactively or with defaults
func ResolveAmbiguities(
	ambiguities []domain.Step1_5Ambiguity,
	existingDecisions *DecisionSet,
	interactive bool,
	criticalOnly bool,
) (*DecisionSet, error) {
	result := &DecisionSet{
		Decisions: make([]AmbiguityDecision, 0, len(ambiguities)),
	}

	// Copy existing decisions
	for _, d := range existingDecisions.Decisions {
		result.Decisions = append(result.Decisions, d)
	}

	for _, amb := range ambiguities {
		// Check if already resolved
		if existingDecisions.HasDecision(amb.ID) {
			existing := existingDecisions.GetDecision(amb.ID)
			fmt.Fprintf(os.Stderr, "  [%s] %s (existing: %s)\n", amb.ID, truncate(amb.Question, 50), truncate(existing.Answer, 30))
			continue
		}

		// Determine if we should ask interactively
		shouldAsk := interactive
		if criticalOnly && string(amb.Severity) != "critical" {
			shouldAsk = false
		}

		var answer string
		var source string

		if shouldAsk {
			// Interactive mode
			fmt.Fprintf(os.Stderr, "\n%s [%s] (%s)\n", amb.ID, amb.Category, amb.Severity)
			fmt.Fprintf(os.Stderr, "Question: %s\n", amb.Question)
			fmt.Fprintf(os.Stderr, "Why important: %s\n", amb.WhyImportant)

			if len(amb.SuggestedOptions) > 0 {
				fmt.Fprintf(os.Stderr, "\nOptions:\n")
				for i, opt := range amb.SuggestedOptions {
					fmt.Fprintf(os.Stderr, "  [%d] %s\n", i+1, opt)
				}
				if amb.DefaultSuggestion != "" {
					fmt.Fprintf(os.Stderr, "  [Enter] Use default: %s\n", amb.DefaultSuggestion)
				}
			}

			fmt.Fprintf(os.Stderr, "\nYour answer (or number, or Enter for default): ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("failed to read input: %w", err)
			}
			input = strings.TrimSpace(input)

			if input == "" {
				// Use default
				answer = amb.DefaultSuggestion
				source = "default"
			} else if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
				// Number selection
				idx := int(input[0] - '1')
				if idx < len(amb.SuggestedOptions) {
					answer = amb.SuggestedOptions[idx]
					source = "user"
				} else {
					answer = input
					source = "user"
				}
			} else {
				answer = input
				source = "user"
			}
		} else {
			// Use default
			answer = amb.DefaultSuggestion
			source = "default"
			fmt.Fprintf(os.Stderr, "  [%s] %s -> %s (default)\n", amb.ID, truncate(amb.Question, 40), truncate(answer, 30))
		}

		result.AddDecision(AmbiguityDecision{
			AmbiguityID: amb.ID,
			Question:    amb.Question,
			Answer:      answer,
			Source:      source,
			Category:    string(amb.Category),
			Severity:    string(amb.Severity),
			DecidedAt:   time.Now(),
		})
	}

	return result, nil
}

// WriteToFile writes decisions to a markdown file
func (ds *DecisionSet) WriteToFile(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	var sb strings.Builder

	// Header
	sb.WriteString("---\n")
	sb.WriteString("title: \"Ambiguity Decisions\"\n")
	sb.WriteString(fmt.Sprintf("generated: %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("status: draft\n")
	sb.WriteString("level: L0\n")
	sb.WriteString("---\n\n")

	sb.WriteString("# Ambiguity Decisions\n\n")
	sb.WriteString("This document records decisions that resolve ambiguities identified during L0 analysis.\n")
	sb.WriteString("These decisions inform the domain modeling process.\n\n")
	sb.WriteString("---\n\n")

	// Group by category
	categories := make(map[string][]AmbiguityDecision)
	categoryOrder := []string{}

	for _, d := range ds.Decisions {
		if _, exists := categories[d.Category]; !exists {
			categoryOrder = append(categoryOrder, d.Category)
		}
		categories[d.Category] = append(categories[d.Category], d)
	}

	for _, cat := range categoryOrder {
		decisions := categories[cat]
		sb.WriteString(fmt.Sprintf("## %s\n\n", formatCategory(cat)))

		for _, d := range decisions {
			sb.WriteString(fmt.Sprintf("### %s\n\n", d.AmbiguityID))
			sb.WriteString(fmt.Sprintf("**Question:** %s\n\n", d.Question))
			sb.WriteString(fmt.Sprintf("**Decision:** %s\n\n", d.Answer))
			sb.WriteString(fmt.Sprintf("**Source:** %s\n\n", d.Source))
			sb.WriteString("---\n\n")
		}
	}

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| ID | Category | Severity | Source |\n")
	sb.WriteString("|----|----------|----------|--------|\n")

	userCount, defaultCount, existingCount := 0, 0, 0
	for _, d := range ds.Decisions {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", d.AmbiguityID, d.Category, d.Severity, d.Source))
		switch d.Source {
		case "user":
			userCount++
		case "default":
			defaultCount++
		case "existing":
			existingCount++
		}
	}

	sb.WriteString(fmt.Sprintf("\n**Statistics:**\n"))
	sb.WriteString(fmt.Sprintf("- Total decisions: %d\n", len(ds.Decisions)))
	sb.WriteString(fmt.Sprintf("- User: %d\n", userCount))
	sb.WriteString(fmt.Sprintf("- Default: %d\n", defaultCount))
	sb.WriteString(fmt.Sprintf("- Existing: %d\n", existingCount))

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatCategory(cat string) string {
	switch cat {
	case "missing_definition":
		return "Missing Definitions"
	case "unclear_relationship":
		return "Unclear Relationships"
	case "synonym_resolution":
		return "Synonym Resolution"
	case "boundary_ambiguity":
		return "Boundary Ambiguities"
	case "business_rule_gap":
		return "Business Rule Gaps"
	case "state_lifecycle":
		return "State & Lifecycle"
	default:
		return strings.Title(strings.ReplaceAll(cat, "_", " "))
	}
}
