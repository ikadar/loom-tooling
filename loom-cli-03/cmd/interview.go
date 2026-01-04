// Package cmd provides CLI commands for loom-cli.
//
// Implements: l2/interface-contracts.md IC-INT-001
// See: l2/sequence-design.md SEQ-INT-001, SEQ-INT-002
// See: l2/tech-specs.md TS-ARCH-001c
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"loom-cli/internal/domain"
	"loom-cli/prompts"
)

// Ensure prompts package is imported (used in interview generation)
var _ = prompts.InterviewPrompt

// runInterview handles the interview command.
//
// Implements: IC-INT-001
// Exit Codes:
//   - 0: Interview complete, no more questions
//   - 1: Error
//   - 100: Question available (output contains question JSON)
func runInterview(args []string) int {
	fs := flag.NewFlagSet("interview", flag.ContinueOnError)
	initFile := fs.String("init", "", "Initialize interview from analysis JSON file")
	stateFile := fs.String("state", "", "Path to interview state file (required)")
	answerJSON := fs.String("answer", "", "Answer as JSON: {\"question_id\":\"...\",\"answer\":\"...\",\"source\":\"user\"}")
	answersJSON := fs.String("answers", "", "Batch answers as JSON array")
	grouped := fs.Bool("grouped", false, "Show all questions grouped by subject")
	skip := fs.Bool("skip", false, "Skip current question with default answer")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return domain.ExitCodeError
	}

	// State file is always required
	if *stateFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --state is required")
		return domain.ExitCodeError
	}

	// Initialize mode
	if *initFile != "" {
		return initializeInterview(*initFile, *stateFile)
	}

	// Load existing state
	state, err := loadInterviewState(*stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read state file: %v\n", err)
		return domain.ExitCodeError
	}

	// Handle batch answers
	if *answersJSON != "" {
		return processBatchAnswers(state, *stateFile, *answersJSON)
	}

	// Handle single answer
	if *answerJSON != "" {
		return processSingleAnswer(state, *stateFile, *answerJSON)
	}

	// Handle skip
	if *skip {
		return processSkip(state, *stateFile)
	}

	// No answer provided - output current question or group
	if *grouped {
		return outputGroupedQuestions(state)
	}
	return outputNextQuestion(state)
}

// initializeInterview creates a new interview state from analysis JSON.
func initializeInterview(analysisFile, stateFile string) int {
	// Read analysis JSON
	data, err := os.ReadFile(analysisFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read analysis file: %v\n", err)
		return domain.ExitCodeError
	}

	var analysis domain.AnalyzeResult
	if err := json.Unmarshal(data, &analysis); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse analysis file: %v\n", err)
		return domain.ExitCodeError
	}

	// Create interview state
	state := domain.InterviewState{
		SessionID:    fmt.Sprintf("interview-%d", time.Now().Unix()),
		DomainModel:  analysis.DomainModel,
		Questions:    analysis.Ambiguities,
		Decisions:    analysis.Decisions,
		CurrentIndex: 0,
		Skipped:      []string{},
		InputContent: analysis.InputContent,
		Complete:     false,
	}

	// Save state
	if err := saveInterviewState(stateFile, &state); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	// Output first question
	return outputNextQuestion(&state)
}

// loadInterviewState loads state from JSON file.
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

// saveInterviewState saves state to JSON file.
func saveInterviewState(path string, state *domain.InterviewState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// processSingleAnswer records an answer and advances to next question.
func processSingleAnswer(state *domain.InterviewState, stateFile, answerJSON string) int {
	var input domain.AnswerInput
	if err := json.Unmarshal([]byte(answerJSON), &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse answer: %v\n", err)
		return domain.ExitCodeError
	}

	// Find the question
	var question *domain.Ambiguity
	for i := range state.Questions {
		if state.Questions[i].ID == input.QuestionID {
			question = &state.Questions[i]
			break
		}
	}

	if question == nil {
		fmt.Fprintf(os.Stderr, "Error: question %s not found\n", input.QuestionID)
		return domain.ExitCodeError
	}

	// Record decision
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

	// Advance index
	state.CurrentIndex++

	// Skip dependent questions based on this answer
	skipDependentQuestions(state)

	// Save state
	if err := saveInterviewState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	return outputNextQuestion(state)
}

// processBatchAnswers records multiple answers at once.
func processBatchAnswers(state *domain.InterviewState, stateFile, answersJSON string) int {
	var inputs []domain.AnswerInput
	if err := json.Unmarshal([]byte(answersJSON), &inputs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse answers: %v\n", err)
		return domain.ExitCodeError
	}

	// Process each answer
	for _, input := range inputs {
		var question *domain.Ambiguity
		for i := range state.Questions {
			if state.Questions[i].ID == input.QuestionID {
				question = &state.Questions[i]
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

	state.CurrentIndex = len(state.Questions)
	skipDependentQuestions(state)

	// Save state
	if err := saveInterviewState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	return outputNextQuestion(state)
}

// processSkip skips the current question with default answer.
func processSkip(state *domain.InterviewState, stateFile string) int {
	if state.CurrentIndex >= len(state.Questions) {
		return outputNextQuestion(state)
	}

	question := state.Questions[state.CurrentIndex]

	// Use suggested answer or mark as skipped
	answer := question.SuggestedAnswer
	if answer == "" {
		answer = "[SKIPPED]"
	}

	decision := domain.Decision{
		ID:        question.ID,
		Question:  question.Question,
		Answer:    answer,
		DecidedAt: time.Now(),
		Source:    "default",
		Category:  question.Category,
		Subject:   question.Subject,
	}
	state.Decisions = append(state.Decisions, decision)
	state.Skipped = append(state.Skipped, question.ID)
	state.CurrentIndex++

	// Save state
	if err := saveInterviewState(stateFile, state); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	return outputNextQuestion(state)
}

// outputNextQuestion outputs the next unanswered question.
func outputNextQuestion(state *domain.InterviewState) int {
	// Find next unanswered, non-skipped question
	for state.CurrentIndex < len(state.Questions) {
		q := state.Questions[state.CurrentIndex]

		// Check if already answered
		answered := false
		for _, d := range state.Decisions {
			if d.ID == q.ID {
				answered = true
				break
			}
		}
		if answered {
			state.CurrentIndex++
			continue
		}

		// Check if should be skipped based on dependencies
		if shouldSkipQuestion(&q, state.Decisions) {
			state.Skipped = append(state.Skipped, q.ID)
			state.CurrentIndex++
			continue
		}

		break
	}

	// Check if complete
	if state.CurrentIndex >= len(state.Questions) {
		state.Complete = true
		output := domain.InterviewOutput{
			Status:         "complete",
			Progress:       fmt.Sprintf("%d/%d", len(state.Decisions), len(state.Questions)),
			RemainingCount: 0,
			SkippedCount:   len(state.Skipped),
			Message:        "Interview complete. All questions answered.",
		}

		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return domain.ExitCodeSuccess
	}

	// Output current question
	question := state.Questions[state.CurrentIndex]
	output := domain.InterviewOutput{
		Status:         "question",
		Question:       &question,
		Progress:       fmt.Sprintf("%d/%d", state.CurrentIndex+1, len(state.Questions)),
		RemainingCount: len(state.Questions) - state.CurrentIndex - 1,
		SkippedCount:   len(state.Skipped),
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
	return domain.ExitCodeQuestion
}

// outputGroupedQuestions outputs all questions grouped by subject.
func outputGroupedQuestions(state *domain.InterviewState) int {
	groups := groupQuestions(state.Questions)

	output := domain.InterviewOutput{
		Status:         "group",
		Group:          groups[0], // First group
		Progress:       fmt.Sprintf("1/%d groups", len(groups)),
		RemainingCount: len(state.Questions) - len(state.Decisions),
		SkippedCount:   len(state.Skipped),
	}

	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
	return domain.ExitCodeQuestion
}

// groupQuestions groups questions by subject.
//
// Implements: TS-ARCH-001c (Question Grouping)
func groupQuestions(questions []domain.Ambiguity) []*domain.QuestionGroup {
	// Group by subject
	subjectOrder := []string{}
	subjectMap := make(map[string][]domain.Ambiguity)

	for _, q := range questions {
		if _, exists := subjectMap[q.Subject]; !exists {
			subjectOrder = append(subjectOrder, q.Subject)
		}
		subjectMap[q.Subject] = append(subjectMap[q.Subject], q)
	}

	// Create groups (max 5 per group)
	var groups []*domain.QuestionGroup
	groupNum := 1

	for _, subject := range subjectOrder {
		qs := subjectMap[subject]

		// Split into chunks of MaxGroupSize
		for i := 0; i < len(qs); i += domain.MaxGroupSize {
			end := i + domain.MaxGroupSize
			if end > len(qs) {
				end = len(qs)
			}

			chunk := qs[i:end]

			// Determine category
			category := chunk[0].Category
			for _, q := range chunk {
				if q.Category != category {
					category = "mixed"
					break
				}
			}

			group := &domain.QuestionGroup{
				ID:        fmt.Sprintf("GRP-%03d", groupNum),
				Subject:   subject,
				Category:  category,
				Questions: chunk,
			}
			groups = append(groups, group)
			groupNum++
		}
	}

	return groups
}

// shouldSkipQuestion checks if a question should be skipped based on dependencies.
//
// Implements: TS-ARCH-001c (Skip Condition Evaluation)
// See: AGG-INT-001
func shouldSkipQuestion(q *domain.Ambiguity, decisions []domain.Decision) bool {
	if len(q.DependsOn) == 0 {
		return false
	}

	for _, dep := range q.DependsOn {
		for _, d := range decisions {
			if d.ID == dep.QuestionID {
				// Check if answer matches any skip condition (case-insensitive)
				for _, skipPhrase := range dep.SkipIfAnswer {
					if containsIgnoreCase(d.Answer, skipPhrase) {
						return true
					}
				}
			}
		}
	}

	return false
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// skipDependentQuestions updates state by marking questions that should be skipped.
func skipDependentQuestions(state *domain.InterviewState) {
	for i := state.CurrentIndex; i < len(state.Questions); i++ {
		q := &state.Questions[i]

		// Check if already in decisions
		answered := false
		for _, d := range state.Decisions {
			if d.ID == q.ID {
				answered = true
				break
			}
		}
		if answered {
			continue
		}

		// Check if should be skipped
		if shouldSkipQuestion(q, state.Decisions) {
			alreadySkipped := false
			for _, id := range state.Skipped {
				if id == q.ID {
					alreadySkipped = true
					break
				}
			}
			if !alreadySkipped {
				state.Skipped = append(state.Skipped, q.ID)
			}
		}
	}
}
