---
title: "Loom CLI Data Model"
generated: 2025-01-03T14:30:00Z
status: draft
level: L2
---

# Loom CLI Data Model

## Overview

This document defines the data model for loom-cli, specifying the JSON schemas for all persisted state files and output formats.

**Traceability:** Derived from [domain-model.md](../l1/domain-model.md) and [aggregate-design.md](aggregate-design.md).

---

## State Files

### TBL-STA-001: Cascade State

**File:** `.cascade-state.json`

**Location:** `{output-dir}/.cascade-state.json`

**Purpose:** Tracks cascade derivation progress for resume capability.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "input_hash", "phases", "config", "timestamps"],
  "properties": {
    "version": {
      "type": "string",
      "const": "1.0"
    },
    "input_hash": {
      "type": "string",
      "description": "SHA256 hash of input file (first 16 chars)"
    },
    "phases": {
      "type": "object",
      "required": ["analyze", "interview", "derive-l1", "derive-l2", "derive-l3"],
      "additionalProperties": {
        "$ref": "#/definitions/PhaseState"
      }
    },
    "config": {
      "$ref": "#/definitions/CascadeConfig"
    },
    "timestamps": {
      "type": "object",
      "required": ["started"],
      "properties": {
        "started": { "type": "string", "format": "date-time" },
        "completed": { "type": "string", "format": "date-time" }
      }
    }
  },
  "definitions": {
    "PhaseState": {
      "type": "object",
      "required": ["status"],
      "properties": {
        "status": {
          "type": "string",
          "enum": ["pending", "running", "completed", "failed"]
        },
        "timestamp": { "type": "string", "format": "date-time" },
        "error": { "type": "string" }
      }
    },
    "CascadeConfig": {
      "type": "object",
      "properties": {
        "skip_interview": { "type": "boolean", "default": false },
        "interactive": { "type": "boolean", "default": false }
      }
    }
  }
}
```

**Example:**
```json
{
  "version": "1.0",
  "input_hash": "a1b2c3d4e5f6g7h8",
  "phases": {
    "analyze": { "status": "completed", "timestamp": "2024-01-15T10:00:00Z" },
    "interview": { "status": "completed", "timestamp": "2024-01-15T10:01:00Z" },
    "derive-l1": { "status": "completed", "timestamp": "2024-01-15T10:03:00Z" },
    "derive-l2": { "status": "running", "timestamp": "2024-01-15T10:05:00Z" },
    "derive-l3": { "status": "pending" }
  },
  "config": {
    "skip_interview": false,
    "interactive": false
  },
  "timestamps": {
    "started": "2024-01-15T10:00:00Z"
  }
}
```

**Related:** AGG-CAS-001, BR-DRV-003

---

### TBL-STA-002: Interview State

**File:** `.interview-state.json`

**Location:** `{output-dir}/.interview-state.json` or specified by `--state`

**Purpose:** Tracks interview session with questions and decisions.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["session_id", "questions", "current_index", "decisions", "complete"],
  "properties": {
    "session_id": { "type": "string" },
    "domain_model": { "$ref": "#/definitions/Domain" },
    "questions": {
      "type": "array",
      "items": { "$ref": "#/definitions/Ambiguity" }
    },
    "current_index": {
      "type": "integer",
      "minimum": 0
    },
    "decisions": {
      "type": "array",
      "items": { "$ref": "#/definitions/Decision" }
    },
    "skipped": {
      "type": "array",
      "items": { "type": "string" }
    },
    "input_content": { "type": "string" },
    "complete": { "type": "boolean" }
  },
  "definitions": {
    "Ambiguity": {
      "type": "object",
      "required": ["id", "category", "subject", "question", "severity"],
      "properties": {
        "id": { "type": "string" },
        "category": { "type": "string", "enum": ["entity", "operation", "ui"] },
        "subject": { "type": "string" },
        "question": { "type": "string" },
        "severity": { "type": "string", "enum": ["critical", "important", "minor"] },
        "suggested_answer": { "type": "string" },
        "options": { "type": "array", "items": { "type": "string" } },
        "checklist_item": { "type": "string" },
        "depends_on": {
          "type": "array",
          "items": { "$ref": "#/definitions/SkipCondition" }
        }
      }
    },
    "SkipCondition": {
      "type": "object",
      "required": ["question_id", "skip_if_answer"],
      "properties": {
        "question_id": { "type": "string" },
        "skip_if_answer": { "type": "array", "items": { "type": "string" } }
      }
    },
    "Decision": {
      "type": "object",
      "required": ["id", "question", "answer", "source"],
      "properties": {
        "id": { "type": "string" },
        "question": { "type": "string" },
        "answer": { "type": "string" },
        "decided_at": { "type": "string", "format": "date-time" },
        "source": {
          "type": "string",
          "enum": ["user", "default", "existing", "user_accepted_suggested"]
        },
        "category": { "type": "string" },
        "subject": { "type": "string" }
      }
    },
    "Domain": {
      "type": "object",
      "properties": {
        "entities": { "type": "array" },
        "operations": { "type": "array" },
        "relationships": { "type": "array" },
        "business_rules": { "type": "array", "items": { "type": "string" } },
        "ui_mentions": { "type": "array", "items": { "type": "string" } }
      }
    }
  }
}
```

**Example:**
```json
{
  "session_id": "interview-1705312800",
  "domain_model": { "entities": [], "operations": [] },
  "questions": [
    {
      "id": "AMB-ENT-ORDER-001",
      "category": "entity",
      "subject": "Order",
      "question": "Can orders be cancelled after shipping?",
      "severity": "critical",
      "suggested_answer": "Only within 24 hours",
      "options": ["Yes, always", "No, never", "Only within 24 hours"],
      "checklist_item": "Order cancellation policy defined"
    }
  ],
  "current_index": 1,
  "decisions": [
    {
      "id": "AMB-ENT-ORDER-001",
      "question": "Can orders be cancelled after shipping?",
      "answer": "Only within 24 hours",
      "decided_at": "2024-01-15T10:05:00Z",
      "source": "user",
      "category": "entity",
      "subject": "Order"
    }
  ],
  "skipped": [],
  "input_content": "...",
  "complete": false
}
```

**Related:** AGG-INT-001, BR-INT-002

**Related DEC:** DEC-L1-003, DEC-L1-004, DEC-L1-007, DEC-L1-008, DEC-L1-012

---

### TBL-STA-003: Analysis Output

**File:** `.analysis.json`

**Location:** `{output-dir}/.analysis.json` or stdout

**Purpose:** Contains discovered domain model and identified ambiguities.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["domain_model", "ambiguities"],
  "properties": {
    "domain_model": { "$ref": "#/definitions/Domain" },
    "ambiguities": {
      "type": "array",
      "items": { "$ref": "#/definitions/Ambiguity" }
    },
    "existing_decisions": {
      "type": "array",
      "items": { "$ref": "#/definitions/Decision" }
    },
    "input_files": {
      "type": "array",
      "items": { "type": "string" }
    },
    "input_content": { "type": "string" }
  },
  "definitions": {
    "Domain": {
      "type": "object",
      "properties": {
        "entities": {
          "type": "array",
          "items": { "$ref": "#/definitions/Entity" }
        },
        "operations": {
          "type": "array",
          "items": { "$ref": "#/definitions/Operation" }
        },
        "relationships": {
          "type": "array",
          "items": { "$ref": "#/definitions/Relationship" }
        },
        "business_rules": {
          "type": "array",
          "items": { "type": "string" }
        },
        "ui_mentions": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    },
    "Entity": {
      "type": "object",
      "required": ["name"],
      "properties": {
        "name": { "type": "string" },
        "mentioned_attributes": {
          "type": "array",
          "items": { "type": "string" }
        },
        "mentioned_operations": {
          "type": "array",
          "items": { "type": "string" }
        },
        "mentioned_states": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    },
    "Operation": {
      "type": "object",
      "required": ["name"],
      "properties": {
        "name": { "type": "string" },
        "actor": { "type": "string" },
        "trigger": { "type": "string" },
        "target": { "type": "string" },
        "mentioned_inputs": {
          "type": "array",
          "items": { "type": "string" }
        },
        "mentioned_rules": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    },
    "Relationship": {
      "type": "object",
      "required": ["from", "to", "type"],
      "properties": {
        "from": { "type": "string" },
        "to": { "type": "string" },
        "type": {
          "type": "string",
          "enum": ["has_many", "belongs_to", "has_one"]
        },
        "cardinality": {
          "type": "string",
          "enum": ["1:1", "1:N", "N:M"]
        }
      }
    },
    "Ambiguity": {
      "type": "object",
      "required": ["id", "category", "subject", "question", "severity"],
      "properties": {
        "id": { "type": "string" },
        "category": { "type": "string", "enum": ["entity", "operation", "ui"] },
        "subject": { "type": "string" },
        "question": { "type": "string" },
        "severity": { "type": "string", "enum": ["critical", "important", "minor"] },
        "suggested_answer": { "type": "string" },
        "options": { "type": "array", "items": { "type": "string" } },
        "checklist_item": { "type": "string" },
        "depends_on": {
          "type": "array",
          "items": { "$ref": "#/definitions/SkipCondition" }
        }
      }
    },
    "SkipCondition": {
      "type": "object",
      "properties": {
        "question_id": { "type": "string" },
        "skip_if_answer": { "type": "array", "items": { "type": "string" } }
      }
    },
    "Decision": {
      "type": "object",
      "properties": {
        "id": { "type": "string" },
        "question": { "type": "string" },
        "answer": { "type": "string" },
        "source": { "type": "string" }
      }
    }
  }
}
```

**Related:** AGG-ANL-001

**Related DEC:** DEC-L1-003, DEC-L1-005, DEC-L1-006, DEC-L1-007

---

## Output Formats

### TBL-OUT-001: Validation Result (JSON)

**Output:** stdout with `--json` flag

**Purpose:** Machine-readable validation results for CI/CD.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["level", "summary"],
  "properties": {
    "level": {
      "type": "string",
      "enum": ["L1", "L2", "L3", "ALL"]
    },
    "errors": {
      "type": "array",
      "items": { "$ref": "#/definitions/ValidationError" }
    },
    "warnings": {
      "type": "array",
      "items": { "$ref": "#/definitions/ValidationWarning" }
    },
    "checks": {
      "type": "array",
      "items": { "$ref": "#/definitions/ValidationCheck" }
    },
    "summary": { "$ref": "#/definitions/ValidationSummary" }
  },
  "definitions": {
    "ValidationError": {
      "type": "object",
      "required": ["rule", "message"],
      "properties": {
        "file": { "type": "string" },
        "line": { "type": "integer" },
        "rule": { "type": "string", "pattern": "^V0[01]\\d$" },
        "message": { "type": "string" },
        "ref_id": { "type": "string" }
      }
    },
    "ValidationWarning": {
      "type": "object",
      "required": ["rule", "message"],
      "properties": {
        "file": { "type": "string" },
        "line": { "type": "integer" },
        "rule": { "type": "string" },
        "message": { "type": "string" }
      }
    },
    "ValidationCheck": {
      "type": "object",
      "required": ["rule", "status", "message"],
      "properties": {
        "rule": { "type": "string" },
        "status": {
          "type": "string",
          "enum": ["pass", "fail", "skip"]
        },
        "message": { "type": "string" },
        "count": { "type": "integer" }
      }
    },
    "ValidationSummary": {
      "type": "object",
      "required": ["total_checks", "passed", "failed"],
      "properties": {
        "total_checks": { "type": "integer" },
        "passed": { "type": "integer" },
        "failed": { "type": "integer" },
        "warnings": { "type": "integer" },
        "error_count": { "type": "integer" }
      }
    }
  }
}
```

**Example:**
```json
{
  "level": "ALL",
  "errors": [
    {
      "file": "acceptance-criteria.md",
      "line": 45,
      "rule": "V003",
      "message": "Reference 'BR-ORD-999' not found",
      "ref_id": "BR-ORD-999"
    }
  ],
  "warnings": [],
  "checks": [
    { "rule": "V001", "status": "pass", "message": "All 5 documents have IDs", "count": 5 },
    { "rule": "V002", "status": "pass", "message": "All 42 IDs follow patterns", "count": 42 },
    { "rule": "V003", "status": "fail", "message": "1 invalid reference found", "count": 1 }
  ],
  "summary": {
    "total_checks": 10,
    "passed": 9,
    "failed": 1,
    "warnings": 0,
    "error_count": 1
  }
}
```

**Related:** AGG-VAL-001, BR-VAL-003

---

### TBL-OUT-002: Interview Question (JSON)

**Output:** stdout during interview

**Purpose:** Machine-readable question for scripted interviews.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["question_id", "text"],
  "properties": {
    "question_id": { "type": "string" },
    "text": { "type": "string" },
    "context": { "type": "string" },
    "options": {
      "type": "array",
      "items": { "type": "string" }
    },
    "remaining": { "type": "integer" },
    "total": { "type": "integer" }
  }
}
```

**Example:**
```json
{
  "question_id": "Q1",
  "text": "Should order cancellation be allowed after shipping?",
  "context": "The user story mentions order cancellation but doesn't specify timing constraints.",
  "options": ["Yes, always", "No, never", "Only within 24 hours"],
  "remaining": 4,
  "total": 5
}
```

**Related:** BR-INT-001

---

### TBL-OUT-003: Answer Input (JSON)

**Input:** `--answer` flag value

**Purpose:** Provide answer to interview question.

**Schema:**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["question_id", "answer", "source"],
  "properties": {
    "question_id": { "type": "string" },
    "answer": { "type": "string" },
    "source": {
      "type": "string",
      "enum": ["user", "default", "existing", "user_accepted_suggested"]
    }
  }
}
```

**Example:**
```json
{
  "question_id": "AMB-ENT-ORDER-001",
  "answer": "Only within 24 hours",
  "source": "user"
}
```

**Related:** BR-INT-002, DEC-011, DEC-L1-004

---

## File Locations

| File | Location | Created By |
|------|----------|------------|
| `.cascade-state.json` | `{output-dir}/` | cascade |
| `.interview-state.json` | `{output-dir}/` or `--state` | interview |
| `.analysis.json` | `{output-dir}/` | cascade (analyze phase) |
| `l1/*.md` | `{output-dir}/l1/` | derive, cascade |
| `l2/*.md` | `{output-dir}/l2/` | derive-l2, cascade |
| `l3/*.md` | `{output-dir}/l3/` | derive-l3, cascade |

---

## Related Documents

| Level | Document | Description |
|-------|----------|-------------|
| L0 | [domain-vocabulary.md](../l0/domain-vocabulary.md) | Domain Vocabulary |
| L0 | [loom-cli.md](../l0/loom-cli.md) | User Stories |
| L0 | [nfr.md](../l0/nfr.md) | Non-Functional Requirements |
| L0 | [decisions.md](../l0/decisions.md) | Design Decisions (L0→L1) |
| L1 | [domain-model.md](../l1/domain-model.md) | Domain Model (source) |
| L1 | [decisions.md](../l1/decisions.md) | Design Decisions (L1→L2) |
| L2 | [decisions.md](decisions.md) | Design Decisions (L2→L3) |
| L2 | [prompt-catalog.md](prompt-catalog.md) | Prompt Catalog |
| L1 | [bounded-context-map.md](../l1/bounded-context-map.md) | Bounded Context Map |
| L1 | [business-rules.md](../l1/business-rules.md) | Business Rules |
| L1 | [acceptance-criteria.md](../l1/acceptance-criteria.md) | Acceptance Criteria |
| L2 | [tech-specs.md](tech-specs.md) | Technical Specifications |
| L2 | [interface-contracts.md](interface-contracts.md) | CLI Interface Contract |
| L2 | [aggregate-design.md](aggregate-design.md) | Aggregate Design |
| L2 | [sequence-design.md](sequence-design.md) | Sequence Design |
| L2 | This document | Data Model |
