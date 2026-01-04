// Package cmd provides CLI commands for loom-cli.
//
// This file implements the interview command.
// Implements: l2/interface-contracts.md IC-INT-001
// Implements: l2/tech-specs.md TS-ARCH-001c
// See: l2/sequence-design.md SEQ-INT-001, SEQ-INT-002
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"loom-cli/internal/domain"
	"loom-cli/internal/interview"
)

// runInterview implements the interview command.
// Implements: IC-INT-001
//
// Modes:
//
//	--init         Initialize new interview from analysis.json
//	(no flag)      Show next question
//	--answer       Record single answer
//	--answers      Record multiple answers (batch mode)
//	--grouped      Enable grouped question mode
//
// Exit Codes:
//
//	0   - Interview complete (no more questions)
//	1   - Error
//	100 - More questions available
func runInterview(args []string) int {
	fs := flag.NewFlagSet("interview", flag.ExitOnError)

	initFile := fs.String("init", "", "Initialize interview from analysis.json")
	stateFile := fs.String("state", "interview-state.json", "Interview state file")
	answerJSON := fs.String("answer", "", "Answer JSON: {\"question_id\":\"...\",\"answer\":\"...\",\"source\":\"user\"}")
	answersJSON := fs.String("answers", "", "Multiple answers JSON array")
	grouped := fs.Bool("grouped", false, "Show questions grouped by subject")
	outputJSON := fs.Bool("json", true, "Output as JSON")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// Mode: Initialize
	if *initFile != "" {
		return initInterview(*initFile, *stateFile)
	}

	// Load existing state
	state, err := loadInterviewState(*stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'loom-cli interview --init <analysis.json>' first\n")
		return domain.ExitCodeError
	}

	// Mode: Record answer
	if *answerJSON != "" {
		return recordAnswer(state, *answerJSON, *stateFile)
	}

	// Mode: Record multiple answers
	if *answersJSON != "" {
		return recordAnswers(state, *answersJSON, *stateFile)
	}

	// Mode: Show next question
	return showNextQuestion(state, *grouped, *outputJSON)
}

// initInterview creates a new interview state from analysis result.
func initInterview(analysisFile, stateFile string) int {
	// Load analysis result
	data, err := os.ReadFile(analysisFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading analysis file: %v\n", err)
		return domain.ExitCodeError
	}

	var analysis domain.AnalyzeResult
	if err := json.Unmarshal(data, &analysis); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing analysis file: %v\n", err)
		return domain.ExitCodeError
	}

	// Add dependencies to questions
	// Implements: DEC-L1-008 (automatic dependency inference)
	questions := interview.AddDependencies(analysis.Ambiguities)

	// Create interview state
	state := &domain.InterviewState{
		SessionID:    fmt.Sprintf("session-%d", time.Now().Unix()),
		DomainModel:  analysis.DomainModel,
		Questions:    questions,
		Decisions:    analysis.Decisions,
		CurrentIndex: 0,
		InputContent: analysis.InputContent,
		Complete:     len(questions) == 0,
	}

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		return domain.ExitCodeError
	}

	fmt.Fprintf(os.Stderr, "Interview initialized with %d questions\n", len(questions))

	if state.Complete {
		return domain.ExitCodeSuccess
	}
	return domain.ExitCodeQuestion
}

// recordAnswer records a single answer.
func recordAnswer(state *domain.InterviewState, answerJSON, stateFile string) int {
	var input domain.AnswerInput
	if err := json.Unmarshal([]byte(answerJSON), &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing answer: %v\n", err)
		return domain.ExitCodeError
	}

	// Find the question
	var question *domain.Ambiguity
	for _, q := range state.Questions {
		if q.ID == input.QuestionID {
			question = &q
			break
		}
	}

	if question == nil {
		fmt.Fprintf(os.Stderr, "Error: question not found: %s\n", input.QuestionID)
		return domain.ExitCodeError
	}

	// Create decision
	decision := domain.Decision{
		ID:        input.QuestionID,
		Question:  question.Question,
		Answer:    input.Answer,
		DecidedAt: time.Now(),
		Source:    input.Source,
		Category:  question.Category,
		Subject:   question.Subject,
	}

	state.Decisions = append(state.Decisions, decision)

	// Check if any questions should now be skipped
	remaining := interview.FilterSkippedQuestions(state.Questions, state.Decisions)
	answeredCount := len(state.Decisions)

	// Mark skipped questions
	for _, q := range state.Questions {
		if interview.ShouldSkipQuestion(q, state.Decisions) {
			alreadySkipped := false
			for _, s := range state.Skipped {
				if s == q.ID {
					alreadySkipped = true
					break
				}
			}
			if !alreadySkipped {
				state.Skipped = append(state.Skipped, q.ID)
			}
		}
	}

	// Check if complete
	state.Complete = len(remaining) <= answeredCount

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		return domain.ExitCodeError
	}

	if state.Complete {
		output := domain.InterviewOutput{
			Status:         "complete",
			RemainingCount: 0,
			SkippedCount:   len(state.Skipped),
			Message:        "Interview complete",
		}
		outputJSON, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJSON))
		return domain.ExitCodeSuccess
	}

	return domain.ExitCodeQuestion
}

// recordAnswers records multiple answers.
func recordAnswers(state *domain.InterviewState, answersJSON, stateFile string) int {
	var inputs []domain.AnswerInput
	if err := json.Unmarshal([]byte(answersJSON), &inputs); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing answers: %v\n", err)
		return domain.ExitCodeError
	}

	for _, input := range inputs {
		// Find the question
		var question *domain.Ambiguity
		for _, q := range state.Questions {
			if q.ID == input.QuestionID {
				question = &q
				break
			}
		}

		if question == nil {
			continue
		}

		decision := domain.Decision{
			ID:        input.QuestionID,
			Question:  question.Question,
			Answer:    input.Answer,
			DecidedAt: time.Now(),
			Source:    input.Source,
			Category:  question.Category,
			Subject:   question.Subject,
		}

		state.Decisions = append(state.Decisions, decision)
	}

	// Update skipped and completion status
	for _, q := range state.Questions {
		if interview.ShouldSkipQuestion(q, state.Decisions) {
			alreadySkipped := false
			for _, s := range state.Skipped {
				if s == q.ID {
					alreadySkipped = true
					break
				}
			}
			if !alreadySkipped {
				state.Skipped = append(state.Skipped, q.ID)
			}
		}
	}

	remaining := interview.FilterSkippedQuestions(state.Questions, state.Decisions)
	answeredCount := len(state.Decisions)
	state.Complete = len(remaining) <= answeredCount

	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving state: %v\n", err)
		return domain.ExitCodeError
	}

	if state.Complete {
		return domain.ExitCodeSuccess
	}
	return domain.ExitCodeQuestion
}

// showNextQuestion shows the next unanswered question.
func showNextQuestion(state *domain.InterviewState, grouped, outputJSON bool) int {
	if state.Complete {
		output := domain.InterviewOutput{
			Status:         "complete",
			RemainingCount: 0,
			SkippedCount:   len(state.Skipped),
			Message:        "Interview complete",
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return domain.ExitCodeSuccess
	}

	// Get remaining questions
	remaining := interview.FilterSkippedQuestions(state.Questions, state.Decisions)

	// Build answered set
	answered := make(map[string]bool)
	for _, d := range state.Decisions {
		answered[d.ID] = true
	}

	// Filter to unanswered
	var unanswered []domain.Ambiguity
	for _, q := range remaining {
		if !answered[q.ID] {
			unanswered = append(unanswered, q)
		}
	}

	if len(unanswered) == 0 {
		state.Complete = true
		output := domain.InterviewOutput{
			Status:         "complete",
			RemainingCount: 0,
			SkippedCount:   len(state.Skipped),
			Message:        "Interview complete",
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return domain.ExitCodeSuccess
	}

	var output domain.InterviewOutput

	if grouped {
		// Group mode
		groups := interview.GroupQuestions(unanswered)
		if len(groups) > 0 {
			output = domain.InterviewOutput{
				Status:         "group",
				Group:          &groups[0],
				Progress:       interview.GetProgress(state),
				RemainingCount: len(unanswered),
				SkippedCount:   len(state.Skipped),
			}
		}
	} else {
		// Single question mode
		output = domain.InterviewOutput{
			Status:         "question",
			Question:       &unanswered[0],
			Progress:       interview.GetProgress(state),
			RemainingCount: len(unanswered),
			SkippedCount:   len(state.Skipped),
		}
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
	return domain.ExitCodeQuestion
}

// loadInterviewState loads state from file.
func loadInterviewState(path string) (*domain.InterviewState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state domain.InterviewState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// saveInterviewState saves state to file.
func saveInterviewState(state *domain.InterviewState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
