package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ikadar/loom-cli/internal/domain"
	"github.com/ikadar/loom-cli/internal/interview"
)

const (
	ExitCodeComplete  = 0   // Interview complete, no more questions
	ExitCodeError     = 1   // Error occurred
	ExitCodeQuestion  = 100 // There's a question to answer
)

func runInterview() error {
	// Parse arguments
	args := os.Args[2:]

	var stateFile string
	var answerJSON string
	var answersJSON string // For batch answers (grouped mode)
	var initFile string
	var grouped bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--state":
			if i+1 < len(args) {
				i++
				stateFile = args[i]
			}
		case "--answer":
			if i+1 < len(args) {
				i++
				answerJSON = args[i]
			}
		case "--answers":
			if i+1 < len(args) {
				i++
				answersJSON = args[i]
			}
		case "--init":
			if i+1 < len(args) {
				i++
				initFile = args[i]
			}
		case "--grouped", "-g":
			grouped = true
		}
	}

	// Mode 1: Initialize from analysis file
	if initFile != "" {
		return initInterview(initFile, stateFile, grouped)
	}

	// Mode 2: Continue interview with answer
	if stateFile == "" {
		return fmt.Errorf("--state is required")
	}

	return continueInterview(stateFile, answerJSON, answersJSON)
}

// initInterview creates a new interview state from analysis output
func initInterview(analysisFile, stateFile string, grouped bool) error {
	// Read analysis file
	content, err := os.ReadFile(analysisFile)
	if err != nil {
		return fmt.Errorf("failed to read analysis file: %w", err)
	}

	var analysis struct {
		DomainModel  *domain.Domain     `json:"domain_model"`
		Ambiguities  []domain.Ambiguity `json:"ambiguities"`
		Decisions    []domain.Decision  `json:"existing_decisions"`
		InputContent string             `json:"input_content"`
	}

	if err := json.Unmarshal(content, &analysis); err != nil {
		return fmt.Errorf("failed to parse analysis file: %w", err)
	}

	// Add dependency information to questions
	questions := addDependencies(analysis.Ambiguities)

	// Create initial state
	state := domain.InterviewState{
		SessionID:    fmt.Sprintf("interview-%d", time.Now().Unix()),
		DomainModel:  analysis.DomainModel,
		Questions:    questions,
		Decisions:    analysis.Decisions,
		CurrentIndex: 0,
		Skipped:      []string{},
		InputContent: analysis.InputContent,
		Complete:     false,
	}

	// Save state
	if stateFile == "" {
		stateFile = fmt.Sprintf("/tmp/%s.json", state.SessionID)
	}

	if err := saveState(&state, stateFile); err != nil {
		return err
	}

	// Output first question (or group)
	if grouped {
		return outputNextGroup(&state, stateFile)
	}
	return outputNextQuestion(&state, stateFile)
}

// continueInterview processes an answer and returns the next question
func continueInterview(stateFile, answerJSON, answersJSON string) error {
	// Load state
	state, err := loadState(stateFile)
	if err != nil {
		return err
	}

	// Determine if we're in grouped mode based on answersJSON
	grouped := answersJSON != ""

	// Process batch answers if provided (grouped mode)
	if answersJSON != "" {
		var answers []struct {
			QuestionID string `json:"question_id"`
			Answer     string `json:"answer"`
			Source     string `json:"source"`
		}

		if err := json.Unmarshal([]byte(answersJSON), &answers); err != nil {
			return fmt.Errorf("failed to parse answers: %w", err)
		}

		for _, answer := range answers {
			for _, q := range state.Questions {
				if q.ID == answer.QuestionID {
					decision := domain.Decision{
						ID:        q.ID,
						Question:  q.Question,
						Answer:    answer.Answer,
						DecidedAt: time.Now(),
						Source:    answer.Source,
						Category:  q.Category,
						Subject:   q.Subject,
					}
					state.Decisions = append(state.Decisions, decision)
					state.CurrentIndex++ // Advance for each answer
					break
				}
			}
		}
	} else if answerJSON != "" {
		// Process single answer (legacy mode)
		var answer struct {
			QuestionID string `json:"question_id"`
			Answer     string `json:"answer"`
			Source     string `json:"source"` // "user" or "user_accepted_suggested"
		}

		if err := json.Unmarshal([]byte(answerJSON), &answer); err != nil {
			return fmt.Errorf("failed to parse answer: %w", err)
		}

		// Find the question and record decision
		for _, q := range state.Questions {
			if q.ID == answer.QuestionID {
				decision := domain.Decision{
					ID:        q.ID,
					Question:  q.Question,
					Answer:    answer.Answer,
					DecidedAt: time.Now(),
					Source:    answer.Source,
					Category:  q.Category,
					Subject:   q.Subject,
				}
				state.Decisions = append(state.Decisions, decision)
				break
			}
		}

		// Move to next question
		state.CurrentIndex++
	}

	// Save updated state
	if err := saveState(state, stateFile); err != nil {
		return err
	}

	// Output next question/group (or complete)
	if grouped {
		return outputNextGroup(state, stateFile)
	}
	return outputNextQuestion(state, stateFile)
}

// outputNextQuestion finds and outputs the next unanswered, non-skipped question
func outputNextQuestion(state *domain.InterviewState, stateFile string) error {
	totalQuestions := len(state.Questions)

	for state.CurrentIndex < totalQuestions {
		q := state.Questions[state.CurrentIndex]

		// Check if this question should be skipped
		if shouldSkip(&q, state.Decisions) {
			state.Skipped = append(state.Skipped, q.ID)
			state.CurrentIndex++

			// Save state after skip
			saveState(state, stateFile)
			continue
		}

		// Found a question to ask
		answeredCount := len(state.Decisions) - countExisting(state.Decisions)
		remaining := totalQuestions - state.CurrentIndex - len(state.Skipped)

		output := domain.InterviewOutput{
			Status:         "question",
			Question:       &q,
			Progress:       fmt.Sprintf("%d/%d", answeredCount+1, totalQuestions-len(state.Skipped)),
			RemainingCount: remaining,
			SkippedCount:   len(state.Skipped),
		}

		outputJSON(output)
		os.Exit(ExitCodeQuestion)
		return nil
	}

	// No more questions - interview complete
	state.Complete = true
	saveState(state, stateFile)

	output := domain.InterviewOutput{
		Status:         "complete",
		RemainingCount: 0,
		SkippedCount:   len(state.Skipped),
		Message:        fmt.Sprintf("Interview complete. %d decisions recorded, %d questions skipped.",
			len(state.Decisions)-countExisting(state.Decisions), len(state.Skipped)),
	}

	outputJSON(output)
	os.Exit(ExitCodeComplete)
	return nil
}

// outputNextGroup finds and outputs the next group of related questions
func outputNextGroup(state *domain.InterviewState, stateFile string) error {
	totalQuestions := len(state.Questions)

	// Filter out answered and skipped questions
	remaining := interview.FilterAnsweredQuestions(state.Questions, state.Decisions)
	remaining = interview.FilterSkippedQuestions(remaining, state.Decisions, state.Skipped)

	if len(remaining) == 0 {
		// No more questions - interview complete
		state.Complete = true
		saveState(state, stateFile)

		output := domain.InterviewOutput{
			Status:         "complete",
			RemainingCount: 0,
			SkippedCount:   len(state.Skipped),
			Message: fmt.Sprintf("Interview complete. %d decisions recorded, %d questions skipped.",
				len(state.Decisions)-countExisting(state.Decisions), len(state.Skipped)),
		}

		outputJSON(output)
		os.Exit(ExitCodeComplete)
		return nil
	}

	// Group the remaining questions
	groups := interview.GroupQuestions(remaining)

	if len(groups) == 0 {
		// Fallback: shouldn't happen, but handle gracefully
		return outputNextQuestion(state, stateFile)
	}

	// Output the first group
	group := groups[0]
	answeredCount := len(state.Decisions) - countExisting(state.Decisions)

	output := domain.InterviewOutput{
		Status:         "group",
		Group:          &group,
		Progress:       fmt.Sprintf("%d/%d", answeredCount+1, totalQuestions-len(state.Skipped)),
		RemainingCount: len(remaining),
		SkippedCount:   len(state.Skipped),
		Message:        fmt.Sprintf("%d questions about %s", len(group.Questions), group.Subject),
	}

	outputJSON(output)
	os.Exit(ExitCodeQuestion)
	return nil
}

// shouldSkip checks if a question should be skipped based on previous answers
func shouldSkip(q *domain.Ambiguity, decisions []domain.Decision) bool {
	if len(q.DependsOn) == 0 {
		return false
	}

	for _, dep := range q.DependsOn {
		for _, d := range decisions {
			if d.ID == dep.QuestionID {
				// Check if answer matches any skip condition
				answerLower := strings.ToLower(d.Answer)
				for _, skipPhrase := range dep.SkipIfAnswer {
					if strings.Contains(answerLower, strings.ToLower(skipPhrase)) {
						return true
					}
				}
			}
		}
	}

	return false
}

// addDependencies adds dependency information to questions
// This uses heuristics to detect common patterns
func addDependencies(questions []domain.Ambiguity) []domain.Ambiguity {
	result := make([]domain.Ambiguity, len(questions))
	copy(result, questions)

	// Build a map of questions by subject and keywords
	questionMap := make(map[string]string) // keyword -> question ID

	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		// Detect "can X be deleted" type questions
		if strings.Contains(qLower, "can") && strings.Contains(qLower, "deleted") {
			questionMap[q.Subject+"_delete"] = q.ID
		}
		if strings.Contains(qLower, "can") && strings.Contains(qLower, "modified") {
			questionMap[q.Subject+"_modify"] = q.ID
		}
	}

	// Add dependencies for follow-up questions
	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		// If question is about "after deletion" or "when deleted", depend on delete question
		if strings.Contains(qLower, "after delet") ||
		   strings.Contains(qLower, "when delet") ||
		   strings.Contains(qLower, "deletion cascade") ||
		   strings.Contains(qLower, "upon deletion") {
			if depID, ok := questionMap[q.Subject+"_delete"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID:   depID,
					SkipIfAnswer: []string{"cannot be deleted", "no deletion", "not deletable", "cannot delete", "no, ", "soft delete only"},
				})
			}
		}

		// If question is about "after modification", depend on modify question
		if strings.Contains(qLower, "after modif") ||
		   strings.Contains(qLower, "when modif") ||
		   strings.Contains(qLower, "modification trigger") {
			if depID, ok := questionMap[q.Subject+"_modify"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID:   depID,
					SkipIfAnswer: []string{"cannot be modified", "immutable", "no modification", "cannot modify"},
				})
			}
		}

		// Questions about expiration depend on whether expiration exists
		if strings.Contains(qLower, "when expir") ||
		   strings.Contains(qLower, "after expir") ||
		   strings.Contains(qLower, "expiration notification") {
			// Look for "does X have expiration" question
			for j := range result {
				if j >= i {
					break
				}
				other := &result[j]
				if other.Subject == q.Subject &&
				   strings.Contains(strings.ToLower(other.Question), "expir") &&
				   (strings.Contains(strings.ToLower(other.Question), "have") ||
				    strings.Contains(strings.ToLower(other.Question), "support")) {
					q.DependsOn = append(q.DependsOn, domain.SkipCondition{
						QuestionID:   other.ID,
						SkipIfAnswer: []string{"no expiration", "does not expire", "never expires"},
					})
				}
			}
		}
	}

	return result
}

func countExisting(decisions []domain.Decision) int {
	count := 0
	for _, d := range decisions {
		if d.Source == "existing" {
			count++
		}
	}
	return count
}

func loadState(path string) (*domain.InterviewState, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state domain.InterviewState
	if err := json.Unmarshal(content, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

func saveState(state *domain.InterviewState, path string) error {
	content, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func outputJSON(v interface{}) {
	content, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(content))
}
