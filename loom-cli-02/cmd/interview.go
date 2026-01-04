// Implements: l2/interface-contracts.md IC-INT-001
// See: l2/sequence-design.md SEQ-INT-001, SEQ-INT-002
// See: l2/aggregate-design.md AGG-INT-001
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"loom-cli/internal/domain"
)

// runInterview implements the interview command.
//
// Implements: IC-INT-001
// Exit codes:
//   - 0: Interview complete
//   - 1: Error
//   - 100: More questions available
func runInterview(args []string) int {
	fs := flag.NewFlagSet("interview", flag.ContinueOnError)
	initFile := fs.String("init", "", "Initialize from analysis JSON file")
	stateFile := fs.String("state", "", "Interview state file path (required)")
	answerJSON := fs.String("answer", "", "Answer JSON: {\"question_id\":\"...\",\"answer\":\"...\",\"source\":\"user\"}")
	answersJSON := fs.String("answers", "", "Batch answers as JSON array")
	grouped := fs.Bool("grouped", false, "Show all questions at once (grouped mode)")
	fs.Bool("g", false, "Alias for --grouped")
	skip := fs.Bool("skip", false, "Skip current question with AI default")

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
		return initInterview(*initFile, *stateFile, *grouped)
	}

	// Load existing state
	state, err := loadInterviewState(*stateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read state file: %v\n", err)
		return domain.ExitCodeError
	}

	// Batch answer mode
	if *answersJSON != "" {
		return processBatchAnswers(state, *answersJSON, *stateFile)
	}

	// Single answer mode
	if *answerJSON != "" {
		return processAnswer(state, *answerJSON, *stateFile)
	}

	// Skip mode
	if *skip {
		return skipQuestion(state, *stateFile)
	}

	// Default: show current question
	return showCurrentQuestion(state)
}

// initInterview initializes a new interview from analysis JSON.
func initInterview(analysisFile, stateFile string, grouped bool) int {
	// Read analysis file
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

	// Add automatic dependencies to questions
	questions := addDependencies(analysis.Ambiguities)

	// Create initial state
	state := &domain.InterviewState{
		SessionID:    generateSessionID(),
		DomainModel:  analysis.DomainModel,
		Questions:    questions,
		Decisions:    analysis.Decisions,
		CurrentIndex: 0,
		Skipped:      []string{},
		InputContent: analysis.InputContent,
		Complete:     false,
	}

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	// Check if there are any questions
	if len(state.Questions) == 0 {
		state.Complete = true
		saveInterviewState(state, stateFile)
		output := domain.InterviewOutput{
			Status:  "complete",
			Message: "No questions to answer",
		}
		outputJSON(output)
		return domain.ExitCodeSuccess
	}

	// Output first question or all questions (grouped mode)
	if grouped {
		return showGroupedQuestions(state)
	}
	return showCurrentQuestion(state)
}

// processAnswer processes a single answer.
func processAnswer(state *domain.InterviewState, answerJSON, stateFile string) int {
	var answer domain.AnswerInput
	if err := json.Unmarshal([]byte(answerJSON), &answer); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse answer: %v\n", err)
		return domain.ExitCodeError
	}

	// Find the question
	var question *domain.Ambiguity
	for i := range state.Questions {
		if state.Questions[i].ID == answer.QuestionID {
			question = &state.Questions[i]
			break
		}
	}

	if question == nil {
		fmt.Fprintf(os.Stderr, "Error: question not found: %s\n", answer.QuestionID)
		return domain.ExitCodeError
	}

	// Record decision
	decision := domain.Decision{
		ID:        answer.QuestionID,
		Question:  question.Question,
		Answer:    answer.Answer,
		DecidedAt: time.Now(),
		Source:    answer.Source,
		Category:  question.Category,
		Subject:   question.Subject,
	}
	state.Decisions = append(state.Decisions, decision)

	// Advance to next unskipped question
	advanceToNextQuestion(state)

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	// Check if complete
	if state.Complete {
		output := domain.InterviewOutput{
			Status:  "complete",
			Message: "Interview complete",
		}
		outputJSON(output)
		return domain.ExitCodeSuccess
	}

	return showCurrentQuestion(state)
}

// processBatchAnswers processes multiple answers at once.
func processBatchAnswers(state *domain.InterviewState, answersJSON, stateFile string) int {
	var answers []domain.AnswerInput
	if err := json.Unmarshal([]byte(answersJSON), &answers); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse answers: %v\n", err)
		return domain.ExitCodeError
	}

	for _, answer := range answers {
		// Find the question
		var question *domain.Ambiguity
		for i := range state.Questions {
			if state.Questions[i].ID == answer.QuestionID {
				question = &state.Questions[i]
				break
			}
		}

		if question == nil {
			continue // Skip unknown questions
		}

		// Record decision
		decision := domain.Decision{
			ID:        answer.QuestionID,
			Question:  question.Question,
			Answer:    answer.Answer,
			DecidedAt: time.Now(),
			Source:    answer.Source,
			Category:  question.Category,
			Subject:   question.Subject,
		}
		state.Decisions = append(state.Decisions, decision)
	}

	// Mark as complete
	state.Complete = true
	state.CurrentIndex = len(state.Questions)

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	output := domain.InterviewOutput{
		Status:  "complete",
		Message: "Interview complete",
	}
	outputJSON(output)
	return domain.ExitCodeSuccess
}

// skipQuestion skips the current question with AI default.
func skipQuestion(state *domain.InterviewState, stateFile string) int {
	if state.CurrentIndex >= len(state.Questions) {
		output := domain.InterviewOutput{
			Status:  "complete",
			Message: "Interview complete",
		}
		outputJSON(output)
		return domain.ExitCodeSuccess
	}

	question := &state.Questions[state.CurrentIndex]

	// Record decision with default source
	decision := domain.Decision{
		ID:        question.ID,
		Question:  question.Question,
		Answer:    question.SuggestedAnswer,
		DecidedAt: time.Now(),
		Source:    "default",
		Category:  question.Category,
		Subject:   question.Subject,
	}
	state.Decisions = append(state.Decisions, decision)
	state.Skipped = append(state.Skipped, question.ID)

	// Advance to next question
	advanceToNextQuestion(state)

	// Save state
	if err := saveInterviewState(state, stateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save state: %v\n", err)
		return domain.ExitCodeError
	}

	if state.Complete {
		output := domain.InterviewOutput{
			Status:  "complete",
			Message: "Interview complete",
		}
		outputJSON(output)
		return domain.ExitCodeSuccess
	}

	return showCurrentQuestion(state)
}

// showCurrentQuestion outputs the current question.
func showCurrentQuestion(state *domain.InterviewState) int {
	// Skip questions that should be auto-skipped
	for state.CurrentIndex < len(state.Questions) {
		if !shouldSkipQuestion(&state.Questions[state.CurrentIndex], state.Decisions) {
			break
		}
		state.Skipped = append(state.Skipped, state.Questions[state.CurrentIndex].ID)
		state.CurrentIndex++
	}

	if state.CurrentIndex >= len(state.Questions) {
		state.Complete = true
		output := domain.InterviewOutput{
			Status:  "complete",
			Message: "Interview complete",
		}
		outputJSON(output)
		return domain.ExitCodeSuccess
	}

	question := &state.Questions[state.CurrentIndex]
	output := domain.InterviewOutput{
		Status:         "question",
		Question:       question,
		Progress:       fmt.Sprintf("%d/%d", state.CurrentIndex+1, len(state.Questions)),
		RemainingCount: len(state.Questions) - state.CurrentIndex - 1,
		SkippedCount:   len(state.Skipped),
	}
	outputJSON(output)
	return domain.ExitCodeQuestion
}

// showGroupedQuestions outputs all questions grouped.
func showGroupedQuestions(state *domain.InterviewState) int {
	groups := groupQuestions(state.Questions)

	output := domain.InterviewOutput{
		Status:         "group",
		Group:          &groups[0], // Simplified: show first group
		Progress:       fmt.Sprintf("0/%d", len(state.Questions)),
		RemainingCount: len(state.Questions),
		SkippedCount:   0,
	}
	outputJSON(output)
	return domain.ExitCodeQuestion
}

// advanceToNextQuestion moves to the next unskipped question.
func advanceToNextQuestion(state *domain.InterviewState) {
	state.CurrentIndex++

	for state.CurrentIndex < len(state.Questions) {
		if !shouldSkipQuestion(&state.Questions[state.CurrentIndex], state.Decisions) {
			return
		}
		state.Skipped = append(state.Skipped, state.Questions[state.CurrentIndex].ID)
		state.CurrentIndex++
	}

	state.Complete = true
}

// shouldSkipQuestion checks if a question should be auto-skipped based on previous answers.
//
// Implements: DEC-L1-008 (Question dependency skip logic)
func shouldSkipQuestion(q *domain.Ambiguity, decisions []domain.Decision) bool {
	if len(q.DependsOn) == 0 {
		return false
	}

	for _, dep := range q.DependsOn {
		for _, d := range decisions {
			if d.ID == dep.QuestionID {
				// Check if answer matches any skip phrase (case-insensitive)
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

// addDependencies adds automatic dependency inference to questions.
//
// Implements: l2/aggregate-design.md AGG-INT-001 (Automatic Dependency Inference)
func addDependencies(questions []domain.Ambiguity) []domain.Ambiguity {
	result := make([]domain.Ambiguity, len(questions))
	copy(result, questions)

	// Phase 1: Build capability question map
	questionMap := make(map[string]string) // key -> question ID
	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		if strings.Contains(qLower, "can") && strings.Contains(qLower, "deleted") {
			questionMap[q.Subject+"_delete"] = q.ID
		}
		if strings.Contains(qLower, "can") && strings.Contains(qLower, "modified") {
			questionMap[q.Subject+"_modify"] = q.ID
		}
		if (strings.Contains(qLower, "have") || strings.Contains(qLower, "support")) &&
			strings.Contains(qLower, "expir") {
			questionMap[q.Subject+"_expire"] = q.ID
		}
	}

	// Phase 2: Add dependencies to follow-up questions
	for i := range result {
		q := &result[i]
		qLower := strings.ToLower(q.Question)

		// Deletion follow-up patterns
		if strings.Contains(qLower, "after delet") ||
			strings.Contains(qLower, "when delet") ||
			strings.Contains(qLower, "deletion cascade") ||
			strings.Contains(qLower, "upon deletion") {
			if depID, ok := questionMap[q.Subject+"_delete"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"cannot be deleted",
						"no deletion",
						"not deletable",
						"cannot delete",
						"no, ",
						"soft delete only",
						"immutable",
					},
				})
			}
		}

		// Modification follow-up patterns
		if strings.Contains(qLower, "after modif") ||
			strings.Contains(qLower, "when modif") ||
			strings.Contains(qLower, "modification trigger") {
			if depID, ok := questionMap[q.Subject+"_modify"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"cannot be modified",
						"immutable",
						"no modification",
						"cannot modify",
					},
				})
			}
		}

		// Expiration follow-up patterns
		if strings.Contains(qLower, "when expir") ||
			strings.Contains(qLower, "after expir") ||
			strings.Contains(qLower, "expiration notification") {
			if depID, ok := questionMap[q.Subject+"_expire"]; ok {
				q.DependsOn = append(q.DependsOn, domain.SkipCondition{
					QuestionID: depID,
					SkipIfAnswer: []string{
						"no expiration",
						"does not expire",
						"never expires",
					},
				})
			}
		}
	}

	return result
}

// groupQuestions groups questions by subject.
//
// Implements: DEC-L1-011 (max 5 questions per group)
func groupQuestions(questions []domain.Ambiguity) []domain.QuestionGroup {
	// Group by subject
	subjectMap := make(map[string][]domain.Ambiguity)
	var subjectOrder []string

	for _, q := range questions {
		if _, exists := subjectMap[q.Subject]; !exists {
			subjectOrder = append(subjectOrder, q.Subject)
		}
		subjectMap[q.Subject] = append(subjectMap[q.Subject], q)
	}

	// Create groups (max 5 per group)
	var groups []domain.QuestionGroup
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
			for _, q := range chunk[1:] {
				if q.Category != category {
					category = "mixed"
					break
				}
			}

			groups = append(groups, domain.QuestionGroup{
				ID:        fmt.Sprintf("GRP-%03d", groupNum),
				Subject:   subject,
				Category:  category,
				Questions: chunk,
			})
			groupNum++
		}
	}

	return groups
}

// loadInterviewState loads interview state from file.
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

// saveInterviewState saves interview state to file.
func saveInterviewState(state *domain.InterviewState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// generateSessionID generates a unique session ID.
func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}

// outputJSON outputs a value as JSON to stdout.
func outputJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
