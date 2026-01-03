# Interview Prompt

You are conducting a Structured Interview to resolve domain model and specification ambiguities.

## Your Task

Present questions clearly, grouped by concept. Include decision point IDs for traceability.
Process answers to determine classifications.

## Question Categories

### 1. Entity/Value Object Classification (EVO)
For each concept needing EVO classification:
- Present relevant EVO decision point questions
- Collect answers
- Apply decision logic

### 2. Aggregate Boundary Decisions (AGG)
For concepts needing aggregate decisions:
- Present AGG decision point questions
- Determine same vs separate aggregate

### 3. Reference Type Decisions (REF)
For relationships needing clarification:
- Present REF decision point questions
- Determine embed vs reference

### 4. General Ambiguities (AMB)
For specification ambiguities:
- Present question with options
- Document answer for derivation

## Question Presentation Format

```markdown
---
**[Category: Subject]** (Decision Point: {ID})

Q: {Clear question about the decision}

Options:
a) {option1} → implies {classification/outcome}
b) {option2} → implies {classification/outcome}
c) {option3}
d) Other (please specify)

Suggested default: {suggestion with rationale}

---
```

## Answer Processing

After receiving answers, output:
```json
{
  "decisions": [
    {
      "concept": "ConceptName",
      "decision_point": "EVO-1",
      "question_asked": "Does X need independent identity?",
      "answer_selected": "b",
      "answer_text": "Yes, for shipping tracking",
      "classification_result": "entity",
      "rationale": "User confirmed shipping tracks individual items"
    }
  ],
  "classification_summary": [
    {
      "concept": "ConceptName",
      "final_classification": "entity|value_object",
      "decision_points_resolved": ["EVO-1", "EVO-4"],
      "confidence": "high",
      "rationale": "Based on user answers: ..."
    }
  ]
}
```

## Decision Logic

### EVO Classification:
```
IF (EVO-1 OR EVO-2 OR EVO-3 OR EVO-4) = Yes AND EVO-5 = No:
  → ENTITY
ELSE IF EVO-5 = Yes AND (EVO-1 AND EVO-2 AND EVO-3 AND EVO-4) = No:
  → VALUE OBJECT
ELSE:
  → ASK FOLLOW-UP QUESTIONS
```

### AGG Classification:
```
IF (AGG-1 OR AGG-2) = Yes:
  → SAME AGGREGATE
IF (AGG-3 OR AGG-4) = Yes:
  → SEPARATE AGGREGATE
```

## Quality Checklist
- [ ] Every concept with low/medium confidence has interview questions
- [ ] Decision point IDs included in every question
- [ ] Options include clear implications
- [ ] Answers mapped to classification results
- [ ] Rationale documented for each decision

---

AMBIGUITIES TO ASK:
