package interview

import (
	"fmt"

	"github.com/ikadar/loom-cli/internal/domain"
)

// MaxGroupSize is the maximum number of questions in a group
const MaxGroupSize = 5

// GroupQuestions organizes questions into groups by subject
// Questions with the same subject are grouped together (up to MaxGroupSize)
func GroupQuestions(questions []domain.Ambiguity) []domain.QuestionGroup {
	// First pass: collect questions by subject
	bySubject := make(map[string][]domain.Ambiguity)
	subjectOrder := make([]string, 0)

	for _, q := range questions {
		key := q.Subject
		if key == "" {
			key = q.Category // Fallback to category if no subject
		}
		if _, exists := bySubject[key]; !exists {
			subjectOrder = append(subjectOrder, key)
		}
		bySubject[key] = append(bySubject[key], q)
	}

	// Second pass: create groups (respecting MaxGroupSize)
	var groups []domain.QuestionGroup
	groupCounter := 0

	for _, subject := range subjectOrder {
		qs := bySubject[subject]

		// Split into chunks if too large
		for i := 0; i < len(qs); i += MaxGroupSize {
			end := i + MaxGroupSize
			if end > len(qs) {
				end = len(qs)
			}

			chunk := qs[i:end]
			groupCounter++

			// Determine common category
			category := chunk[0].Category
			for _, q := range chunk[1:] {
				if q.Category != category {
					category = "mixed"
					break
				}
			}

			groups = append(groups, domain.QuestionGroup{
				ID:        fmt.Sprintf("GRP-%03d", groupCounter),
				Subject:   subject,
				Category:  category,
				Questions: chunk,
			})
		}
	}

	return groups
}

// FilterAnsweredQuestions removes questions that have already been answered
func FilterAnsweredQuestions(questions []domain.Ambiguity, decisions []domain.Decision) []domain.Ambiguity {
	answered := make(map[string]bool)
	for _, d := range decisions {
		answered[d.ID] = true
	}

	var remaining []domain.Ambiguity
	for _, q := range questions {
		if !answered[q.ID] {
			remaining = append(remaining, q)
		}
	}
	return remaining
}

// FilterSkippedQuestions removes questions that should be skipped
func FilterSkippedQuestions(questions []domain.Ambiguity, decisions []domain.Decision, skipped []string) []domain.Ambiguity {
	skippedMap := make(map[string]bool)
	for _, id := range skipped {
		skippedMap[id] = true
	}

	var remaining []domain.Ambiguity
	for _, q := range questions {
		if !skippedMap[q.ID] && !shouldSkipQuestion(&q, decisions) {
			remaining = append(remaining, q)
		}
	}
	return remaining
}

// shouldSkipQuestion checks if a question should be skipped based on dependencies
func shouldSkipQuestion(q *domain.Ambiguity, decisions []domain.Decision) bool {
	if len(q.DependsOn) == 0 {
		return false
	}

	for _, dep := range q.DependsOn {
		for _, d := range decisions {
			if d.ID == dep.QuestionID {
				// Check if answer matches any skip condition
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

// containsIgnoreCase checks if s contains substr (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	return contains(sLower, substrLower)
}

// Simple implementations to avoid importing strings package
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
