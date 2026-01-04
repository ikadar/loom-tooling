// Package interview provides interview question grouping and processing.
//
// Implements: l2/package-structure.md PKG-008
// Implements: DEC-L1-011 (question grouping)
// See: l2/internal-api.md
package interview

import (
	"fmt"
	"time"

	"loom-cli/internal/domain"
)

// GroupQuestions groups ambiguities by subject/category.
//
// Grouping Strategy:
// 1. Group by Subject (entity/operation name)
// 2. Within subject, group by Category
// 3. Max 5 questions per group (MaxGroupSize)
// 4. If more, split into multiple groups
//
// Implements: DEC-L1-011
func GroupQuestions(ambiguities []domain.Ambiguity) []domain.QuestionGroup {
	// First, group by subject+category
	groupMap := make(map[string][]domain.Ambiguity)

	for _, amb := range ambiguities {
		key := fmt.Sprintf("%s|%s", amb.Subject, amb.Category)
		groupMap[key] = append(groupMap[key], amb)
	}

	var groups []domain.QuestionGroup
	groupNum := 1

	for key, qs := range groupMap {
		// Split subject and category from key
		subject := ""
		category := ""
		for i, c := range key {
			if c == '|' {
				subject = key[:i]
				category = key[i+1:]
				break
			}
		}

		// Split into chunks of MaxGroupSize
		for i := 0; i < len(qs); i += domain.MaxGroupSize {
			end := i + domain.MaxGroupSize
			if end > len(qs) {
				end = len(qs)
			}

			chunk := qs[i:end]
			group := domain.QuestionGroup{
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

// ShouldSkip checks if a question should be skipped based on previous answers.
//
// Implements: DEC-L1-008 (skip conditions)
func ShouldSkip(question domain.Ambiguity, decisions []domain.Decision) bool {
	if len(question.DependsOn) == 0 {
		return false
	}

	for _, cond := range question.DependsOn {
		// Find the decision for the dependency question
		for _, dec := range decisions {
			if dec.ID == cond.QuestionID {
				// Check if the answer matches any skip condition
				for _, skipAnswer := range cond.SkipIfAnswer {
					if dec.Answer == skipAnswer {
						return true
					}
				}
			}
		}
	}

	return false
}

// ProcessAnswer records an answer and returns updated state.
//
// Implements: l2/internal-api.md
func ProcessAnswer(state *domain.InterviewState, answer domain.AnswerInput) error {
	// Find the question
	var question *domain.Ambiguity
	for i := range state.Questions {
		if state.Questions[i].ID == answer.QuestionID {
			question = &state.Questions[i]
			break
		}
	}

	if question == nil {
		return fmt.Errorf("question not found: %s", answer.QuestionID)
	}

	// Create decision
	decision := domain.Decision{
		ID:        question.ID,
		Question:  question.Question,
		Answer:    answer.Answer,
		DecidedAt: time.Now(),
		Source:    answer.Source,
		Category:  question.Category,
		Subject:   question.Subject,
	}

	state.Decisions = append(state.Decisions, decision)
	state.CurrentIndex++

	// Check for skip conditions on remaining questions
	for i := state.CurrentIndex; i < len(state.Questions); i++ {
		q := state.Questions[i]
		if ShouldSkip(q, state.Decisions) {
			state.Skipped = append(state.Skipped, q.ID)
		}
	}

	// Check if complete
	if state.CurrentIndex >= len(state.Questions) {
		state.Complete = true
	}

	return nil
}

// GetNextQuestion returns the next unanswered, non-skipped question.
//
// Implements: l2/internal-api.md
func GetNextQuestion(state *domain.InterviewState) *domain.Ambiguity {
	for i := state.CurrentIndex; i < len(state.Questions); i++ {
		q := state.Questions[i]

		// Check if skipped
		skipped := false
		for _, id := range state.Skipped {
			if id == q.ID {
				skipped = true
				break
			}
		}

		if !skipped && !ShouldSkip(q, state.Decisions) {
			return &q
		}

		// Skip this one and move to next
		state.CurrentIndex = i + 1
	}

	return nil
}

// CountRemaining returns the number of remaining questions (excluding skipped).
func CountRemaining(state *domain.InterviewState) int {
	count := 0
	for i := state.CurrentIndex; i < len(state.Questions); i++ {
		q := state.Questions[i]

		skipped := false
		for _, id := range state.Skipped {
			if id == q.ID {
				skipped = true
				break
			}
		}

		if !skipped && !ShouldSkip(q, state.Decisions) {
			count++
		}
	}
	return count
}
