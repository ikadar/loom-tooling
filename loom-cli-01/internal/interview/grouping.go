// Package interview provides question grouping and skip logic for interviews.
//
// Implements: l2/aggregate-design.md AGG-INT-001
// Implements: l2/tech-specs.md (Question grouping, skip logic)
// Implements: DEC-L1-008 (automatic dependency inference)
// Implements: DEC-L1-011 (question grouping)
package interview

import (
	"fmt"

	"loom-cli/internal/domain"
)

// GroupQuestions organizes ambiguities into groups by subject.
// Implements: DEC-L1-011 (question grouping)
// Implements: l2/tech-specs.md TS-ARCH-001c
//
// Groups are formed by:
// 1. Same subject (entity/operation name)
// 2. Same category (entity, operation, ui)
// 3. Max group size of 5 questions
func GroupQuestions(ambiguities []domain.Ambiguity) []domain.QuestionGroup {
	// Group by subject and category
	groupMap := make(map[string]*domain.QuestionGroup)

	for _, amb := range ambiguities {
		key := amb.Category + ":" + amb.Subject
		group, exists := groupMap[key]
		if !exists {
			group = &domain.QuestionGroup{
				ID:        "GRP-" + amb.Subject,
				Subject:   amb.Subject,
				Category:  amb.Category,
				Questions: []domain.Ambiguity{},
			}
			groupMap[key] = group
		}

		// Enforce max group size
		if len(group.Questions) < domain.MaxGroupSize {
			group.Questions = append(group.Questions, amb)
		}
	}

	// Convert to slice
	var groups []domain.QuestionGroup
	for _, group := range groupMap {
		if len(group.Questions) > 0 {
			groups = append(groups, *group)
		}
	}

	return groups
}

// AddDependencies infers skip conditions between related questions.
// Implements: DEC-L1-008 (automatic dependency inference)
//
// Skip rules:
// - EVO-1 (identity) = yes → skip EVO-5 (value equality)
// - EVO-5 (value equality) = yes → skip EVO-1,2,3,4
// - AGG-1 (transactional) = yes → skip AGG-3,4 (independent)
// - AGG-3 (independent) = yes → skip AGG-1,2 (same aggregate)
func AddDependencies(ambiguities []domain.Ambiguity) []domain.Ambiguity {
	// Build index by ID
	byID := make(map[string]int)
	for i, amb := range ambiguities {
		byID[amb.ID] = i
	}

	// Define skip rules based on decision point patterns
	skipRules := map[string]map[string][]string{
		// If EVO-1 (identity) answered yes → concept is entity → skip value object questions
		"EVO-1": {"yes": {"EVO-5"}},
		// If EVO-5 (value equality) answered yes → concept is value object → skip entity questions
		"EVO-5": {"yes": {"EVO-1", "EVO-2", "EVO-3", "EVO-4"}},
		// If AGG-1 (transactional boundary) answered yes → same aggregate → skip separate aggregate questions
		"AGG-1": {"yes": {"AGG-3", "AGG-4"}},
		"AGG-2": {"yes": {"AGG-3", "AGG-4"}},
		// If AGG-3 (independent lifecycle) answered yes → separate aggregate → skip same aggregate questions
		"AGG-3": {"yes": {"AGG-1", "AGG-2"}},
		"AGG-4": {"yes": {"AGG-1", "AGG-2"}},
	}

	// Apply rules
	for i := range ambiguities {
		amb := &ambiguities[i]

		// Check if this question's checklist item has skip rules
		for sourceItem, rules := range skipRules {
			if amb.ChecklistItem == sourceItem {
				continue
			}

			// If a question for sourceItem exists, add dependency
			for _, otherAmb := range ambiguities {
				if otherAmb.ChecklistItem == sourceItem && otherAmb.Subject == amb.Subject {
					for answer, skipItems := range rules {
						for _, skipItem := range skipItems {
							if amb.ChecklistItem == skipItem {
								amb.DependsOn = append(amb.DependsOn, domain.SkipCondition{
									QuestionID:   otherAmb.ID,
									SkipIfAnswer: []string{answer},
								})
							}
						}
					}
				}
			}
		}
	}

	return ambiguities
}

// ShouldSkipQuestion determines if a question should be skipped based on existing decisions.
// Implements: DEC-L1-008 (automatic dependency inference)
func ShouldSkipQuestion(question domain.Ambiguity, decisions []domain.Decision) bool {
	if len(question.DependsOn) == 0 {
		return false
	}

	// Build decision lookup
	decisionAnswers := make(map[string]string)
	for _, d := range decisions {
		decisionAnswers[d.ID] = d.Answer
	}

	// Check all skip conditions
	for _, dep := range question.DependsOn {
		answer, exists := decisionAnswers[dep.QuestionID]
		if !exists {
			continue
		}

		for _, skipAnswer := range dep.SkipIfAnswer {
			if answer == skipAnswer {
				return true
			}
		}
	}

	return false
}

// FilterSkippedQuestions removes questions that should be skipped.
func FilterSkippedQuestions(questions []domain.Ambiguity, decisions []domain.Decision) []domain.Ambiguity {
	var result []domain.Ambiguity
	for _, q := range questions {
		if !ShouldSkipQuestion(q, decisions) {
			result = append(result, q)
		}
	}
	return result
}

// GetNextQuestion returns the next unanswered question.
// Returns nil if all questions are answered.
func GetNextQuestion(state *domain.InterviewState) *domain.Ambiguity {
	// Filter out skipped questions
	remaining := FilterSkippedQuestions(state.Questions, state.Decisions)

	// Build answered set
	answered := make(map[string]bool)
	for _, d := range state.Decisions {
		answered[d.ID] = true
	}

	// Find first unanswered
	for _, q := range remaining {
		if !answered[q.ID] {
			return &q
		}
	}

	return nil
}

// GetProgress returns a progress string like "3/10".
func GetProgress(state *domain.InterviewState) string {
	total := len(state.Questions)
	answered := len(state.Decisions)
	skipped := len(state.Skipped)

	return fmt.Sprintf("%d/%d (skipped: %d)", answered, total, skipped)
}
