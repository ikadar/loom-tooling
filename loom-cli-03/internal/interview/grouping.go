// Package interview provides interview question grouping and processing.
//
// Implements: l2/package-structure.md PKG-008
// Implements: DEC-L1-011 (question grouping)
package interview

import (
	"fmt"
	"strings"
	"time"

	"loom-cli/internal/domain"
)

// GroupQuestions groups ambiguities by subject/category.
//
// Implements: DEC-L1-011
// Max group size: 5 questions (domain.MaxGroupSize)
//
// Grouping Strategy:
// 1. Group by Subject (entity/operation name)
// 2. Within subject, group by Category
// 3. Max 5 questions per group (MaxGroupSize)
// 4. If more, split into multiple groups
func GroupQuestions(ambiguities []domain.Ambiguity) []domain.QuestionGroup {
	// Group by subject
	subjectOrder := []string{}
	subjectMap := make(map[string][]domain.Ambiguity)

	for _, q := range ambiguities {
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
			for _, q := range chunk {
				if q.Category != category {
					category = "mixed"
					break
				}
			}

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

	for _, dep := range question.DependsOn {
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

// ProcessAnswer records an answer and returns updated state.
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
		return fmt.Errorf("question %s not found", answer.QuestionID)
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

	// Update dependent question skip status
	updateSkippedQuestions(state)

	return nil
}

// updateSkippedQuestions marks questions that should be skipped.
func updateSkippedQuestions(state *domain.InterviewState) {
	for i := range state.Questions {
		q := &state.Questions[i]

		// Skip if already answered
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
		if ShouldSkip(*q, state.Decisions) {
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

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
