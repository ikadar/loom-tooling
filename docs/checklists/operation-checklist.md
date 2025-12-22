# Operation Completeness Checklist

Reference document for operation analysis. Used by `/loom-analyze-operations` command.

---

## A. Basic Definition

| Aspect | Question |
|--------|----------|
| Name | What is the operation called? |
| Purpose | What does this operation accomplish? |
| Actor | Who can perform this? (role, permission, owner, system) |
| Trigger | What initiates this? (button click, API call, scheduled, event) |
| Target | What entity/entities are affected? |
| Preconditions | What must be true before operation can start? |
| Postconditions | What is guaranteed to be true after success? |
| Atomicity | Is this all-or-nothing? Can it partially succeed? |
| Idempotency | What if called twice with same input? Same result? |

---

## B. Input Definition

For EACH input parameter:

| Aspect | Question |
|--------|----------|
| Name | Parameter name |
| Type | Data type (string, int, array, object, file) |
| Required | Is it required? |
| Default | Default value if not provided? |
| Min | Minimum value/length? |
| Max | Maximum value/length? |
| Format | Expected format/pattern? |
| Validation | Other validation rules? |
| Error | Error message if invalid? |

**Output format:**

```markdown
| Input | Type | Required | Default | Min | Max | Format | Validation |
|-------|------|----------|---------|-----|-----|--------|------------|
| task_id | UUID | yes | - | - | - | UUID v4 | must exist |
| start_time | datetime | yes | - | now | +1 year | ISO 8601 | must be in future |
```

---

## C. Output Definition

| Aspect | Question |
|--------|----------|
| Success Response | What is returned on success? |
| Response Type | Data type of response |
| Always Present | Which fields are always present? |
| Nullable Fields | Which fields can be null? |
| Partial Success | Is partial success possible? What's returned? |

---

## D. Error Handling

For EACH possible error:

| Error Type | Questions |
|------------|-----------|
| Input Validation | What if input is invalid? |
| Precondition Failed | What if precondition not met? |
| Not Found | What if target entity doesn't exist? |
| Conflict | What if operation conflicts with current state? |
| Permission Denied | What if actor lacks permission? |
| Concurrent Modification | What if entity changed during operation? |
| System Error | What if unexpected system error? |
| Timeout | What if operation times out? |
| Rate Limited | What if rate limit exceeded? |

For each error, specify:

| Aspect | Question |
|--------|----------|
| HTTP Status | Status code (400, 404, 409, 500, etc.) |
| Error Code | Machine-readable code (VALIDATION_ERROR, NOT_FOUND) |
| User Message | Human-readable message for UI |
| Technical Message | Detailed message for logs |
| Recovery Action | What can user/system do to recover? |
| Retry | Is retry appropriate? |

---

## E. Side Effects

| Side Effect | Questions |
|-------------|-----------|
| Audit Log | Is change logged? What details? |
| Notification | Who is notified? What channel? (email, push, in-app) |
| Entity Created | Is a related entity created? Which? |
| Entity Updated | Is a related entity updated? Which fields? |
| Entity Deleted | Is a related entity deleted? |
| Event Emitted | Is a domain event emitted? Event name and payload? |
| Cache Invalidation | Which caches need invalidation? |
| External System | Is an external system called? |

For each side effect:
- Is it conditional? On what?
- Is it synchronous or async?
- What if side effect fails? (rollback, ignore, retry)

---

## F. Performance & Limits

| Aspect | Question |
|--------|----------|
| Expected Duration | How long should this take? (p50, p99) |
| Timeout | What's the timeout? |
| Rate Limit | Requests per minute/hour? |
| Batch Size | Max items per request? |
| Payload Size | Max request size? |
| Retry Policy | How many retries? Backoff strategy? |
| Throttling | What if system under load? |

---

## G. Edge Cases (Auto-Generate)

For EACH operation, automatically ask about:

| Edge Case | Question |
|-----------|----------|
| All Optional Omitted | What if all optional inputs omitted? |
| Min Boundary Input | What if input at minimum value? |
| Max Boundary Input | What if input at maximum value? |
| Just Below Min | What if input just below minimum? |
| Just Above Max | What if input just above maximum? |
| Empty String | What if empty string where string expected? |
| Null Value | What if null where value expected? |
| Empty Array | What if empty array where array expected? |
| Target Not Found | What if target entity doesn't exist? |
| Target Just Deleted | What if target deleted between check and operation? |
| Rapid Repeat | What if operation called twice rapidly? |
| During Maintenance | What if called during system maintenance? |
| Network Failure | What if network fails mid-operation? |
| Partial Data | What if some but not all required data available? |
| Stale Data | What if operating on stale data? |
| Large Batch | What if batch at max size? |
| Concurrent Same Target | What if two users operate on same target simultaneously? |

---

## H. Transaction Boundaries

| Aspect | Question |
|--------|----------|
| Transaction Scope | What's included in the transaction? |
| Rollback Trigger | What causes rollback? |
| Compensation | If can't rollback, how to compensate? |
| Saga Pattern | Is this part of a larger saga? |
