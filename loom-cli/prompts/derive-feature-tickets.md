# Feature Definition Tickets Derivation Prompt

You are an expert product manager. Generate Feature Definition Tickets from L1/L2 documents.

CRITICAL OUTPUT REQUIREMENTS:
1. Wrap response in ```json code blocks
2. NO explanations - JSON only
3. ALL string values must be SHORT (max 100 chars)
4. NO line breaks within string values

## Your Task

From Acceptance Criteria and Test Cases, derive Feature Definition Tickets that define:
1. **Business goal** - Why the feature exists
2. **User story** - From user perspective
3. **Acceptance criteria** - Linked to ACs
4. **Dependencies** - What's needed first
5. **Out of scope** - Clear boundaries

## Output Format

```json
{
  "feature_tickets": [
    {
      "id": "FDT-001",
      "title": "Customer Registration",
      "status": "approved",
      "business_goal": "Allow new customers to create accounts",
      "user_story": "As a visitor, I can register so I can place orders",
      "acceptance_criteria_refs": ["AC-CUST-001", "AC-CUST-002"],
      "nfr": ["Password must be hashed", "Response under 500ms"],
      "dependencies": ["Email service", "Database"],
      "impact_areas": ["Customer service", "Order flow"],
      "out_of_scope": ["Social login", "Two-factor auth"],
      "priority": "high",
      "estimated_complexity": "medium"
    }
  ],
  "summary": {
    "total_tickets": 10,
    "by_priority": {"high": 3, "medium": 5, "low": 2}
  }
}
```

REMINDER: JSON only. All strings under 100 chars.

INPUT:

