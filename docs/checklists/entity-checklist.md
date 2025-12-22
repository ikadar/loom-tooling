# Entity Completeness Checklist

Reference document for entity analysis. Used by `/loom-analyze-entities` command.

---

## A. Attribute Definition

For EACH attribute, verify:

| Aspect | Question |
|--------|----------|
| Type | What is the exact data type? (string, int, float, date, datetime, enum, ref, array) |
| Required on Create | Is this field required when creating? |
| Required on Update | Is this field required when updating? |
| Default | What is the default value if not provided? |
| Unique | Must it be unique? Within what scope? |
| Min Value/Length | What is the minimum? |
| Max Value/Length | What is the maximum? |
| Format/Pattern | What format must it match? (regex, date format) |
| Validation Rules | What other validation applies? (cannot be in past, must be positive) |

**Output format:**

```markdown
| Attribute | Type | Required | Default | Unique | Min | Max | Format | Validation |
|-----------|------|----------|---------|--------|-----|-----|--------|------------|
| name | string | yes | - | yes (global) | 1 | 100 | - | no whitespace only |
```

---

## B. Enum Values

For EACH enum/status attribute:

| Aspect | Question |
|--------|----------|
| Values | What are ALL valid values? |
| Initial | What is the initial value on creation? |
| Transitions | What transitions are allowed? (from → to) |
| Transition Trigger | What triggers each transition? (user action, system event, time) |
| Transition Actor | Who/what can trigger each transition? |

**Output format:**

```markdown
### {AttributeName} Enum

Values: `value1`, `value2`, `value3`
Initial: `value1`

| From | To | Trigger | Actor |
|------|----|---------|-------|
| value1 | value2 | User clicks "Approve" | Admin only |
| value2 | value3 | System detects completion | System |
```

---

## C. Lifecycle

| Phase | Questions |
|-------|-----------|
| **Creation** | |
| - Trigger | What triggers creation? (user action, API, import, cascade) |
| - Actor | Who can create? (any user, specific role, system only) |
| - Validation | What must be true to create successfully? |
| - Side Effects | What else happens? (audit log, notification, related entities created) |
| **Update** | |
| - Trigger | What triggers update? |
| - Actor | Who can update? (creator, owner, admin, anyone) |
| - Validation | What must be true to update? |
| - Partial Update | Can individual fields be updated or full replacement only? |
| - Side Effects | What else happens on update? |
| **Deletion** | |
| - Allowed | Can this entity be deleted at all? |
| - Trigger | What triggers deletion? |
| - Actor | Who can delete? |
| - Preconditions | What must be true? (no references, specific state) |
| - Behavior | Hard delete or soft delete? |
| - Cascade | What happens to related entities? |
| - Side Effects | What else happens? |

---

## D. Relationships

For EACH related entity:

| Aspect | Question |
|--------|----------|
| Type | contains, references, belongs_to, many_to_many? |
| Cardinality | 1:1, 1:N, N:1, N:M? |
| Required | Is the relationship required or optional? |
| Mutable | Can relationship be changed after creation? |
| On Parent Delete | What happens to this entity when parent is deleted? |
| On Parent Update | What happens to this entity when parent is updated? |
| On Child Delete | What happens when child is deleted? |
| On Child Update | What happens when child is updated? |

**Output format:**

```markdown
| Related Entity | Type | Cardinality | Required | Mutable | On Parent Delete | On Child Delete |
|----------------|------|-------------|----------|---------|------------------|-----------------|
| Job | belongs_to | N:1 | yes | no | cascade delete | block |
```

---

## E. Constraints & Business Rules

| Aspect | Question |
|--------|----------|
| Scope | Single entity, relationship, or cross-entity? |
| Condition | What condition triggers the rule? |
| Enforcement | Where enforced? (UI only, API, database, multiple) |
| Violation | What happens on violation? (block, warn, auto-fix) |
| Error Message | What does user see? |
| Error Code | Machine-readable error code? |

---

## F. Concurrent Access

| Scenario | Questions |
|----------|-----------|
| Simultaneous Edit | What if two users edit same entity? (first wins, last wins, merge, conflict UI) |
| Edit While Referenced | What if entity modified while being used elsewhere? |
| Delete While Viewed | What if entity deleted while another user viewing? |
| Optimistic Locking | Is version/etag used? |

---

## G. History & Audit

| Aspect | Question |
|--------|----------|
| Logged | Are changes logged? |
| Fields Tracked | All fields or specific ones? |
| View Access | Who can view history? |
| Retention | How long is history kept? |
| Revert | Can changes be reverted? |

---

## H. Edge Cases (Auto-Generate)

For EACH entity, automatically ask about:

| Edge Case | Question |
|-----------|----------|
| All Optional Empty | What if all optional fields are empty/null? |
| Whitespace Only | What if string field contains only whitespace? |
| Min Boundary | What happens at minimum value/length? |
| Max Boundary | What happens at maximum value/length? |
| Just Below Min | What happens just below minimum? |
| Just Above Max | What happens just above maximum? |
| Distant Past Date | What if date is in distant past? |
| Distant Future Date | What if date is in distant future? |
| Invalid Reference | What if reference points to non-existent entity? |
| Circular Reference | What if circular reference is created? (if applicable) |
| Unicode/Special Chars | Are unicode and special characters allowed in strings? |
| Very Long Valid Input | What if input is at max length? Performance? |
| Duplicate | What if exact duplicate is attempted? |
