# Phase 2: Claude Client (ACL)

## Cél

Implementáld az `internal/claude/` package-et: Anti-Corruption Layer a Claude Code CLI-hez.

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/package-structure.md | PKG-005 | Package structure |
| l2/internal-api.md | internal/claude | Client, Response, RetryConfig |
| l2/tech-specs.md | TS-ARCH-002 | CLI invocation |
| l2/tech-specs.md | TS-RETRY-001, TS-RETRY-002 | Retry logic |
| l2/tech-specs.md | SEQ-JSON-001 | JSON extraction |
| l1/decisions.md | DEC-L1-001 | Claude CLI integration |
| l1/decisions.md | DEC-L1-002 | Max output tokens 100,000 |
| l1/decisions.md | DEC-L1-015 | Retry config (3/2s/30s) |

---

## Implementálandó Fájlok

### 1. internal/claude/client.go

```
☐ Fájl: internal/claude/client.go
☐ Spec: l2/package-structure.md PKG-005, l2/tech-specs.md TS-ARCH-002
```

**Traceability komment:**
```go
// Package claude provides the Claude CLI integration layer.
//
// Implements: l2/tech-specs.md TS-ARCH-002 (Anti-Corruption Layer via Claude Code CLI)
// See: l1/bounded-context-map.md (ACL pattern)
package claude
```

**Implementálandó típusok és függvények:**

| Elem | Spec | Leírás |
|------|------|--------|
| `Client` struct | l2/internal-api.md | SessionID, Verbose fields |
| `NewClient()` | l2/internal-api.md | Constructor |
| `Call(prompt string)` | TS-ARCH-002 | Basic CLI call |
| `CallWithSystemPrompt()` | l2/internal-api.md | System + user prompt |
| `CallJSON()` | SEQ-JSON-001 | Parse JSON response |
| `CallJSONWithRetry()` | TS-RETRY-001 | With retry logic |
| `extractJSON()` | SEQ-JSON-001 | Extract JSON from response |
| `sanitizeJSON()` | SEQ-JSON-001 Step 3 | Fix LLM JSON issues |
| `BuildPrompt()` | l2/internal-api.md | Context injection |

**Call() implementáció (TS-ARCH-002):**
```go
func (c *Client) Call(prompt string) (string, error) {
    args := []string{"-p", prompt}
    if c.SessionID != "" {
        args = append(args, "--resume", c.SessionID)
    }

    cmd := exec.Command("claude", args...)

    // Implements: DEC-L1-002 (100,000 max output tokens)
    cmd.Env = append(os.Environ(), "CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000")

    // ... execution
}
```

**extractJSON() Strategy (SEQ-JSON-001):**
1. Markdown Code Block Extraction: ` ```json ... ``` `
2. Raw JSON Extraction: Find first `{` or `[`
3. JSON Sanitization: Fix unescaped newlines
4. Conversion Fallback: Ask Claude to fix if still invalid

**BuildPrompt() Pattern:**
```go
// All prompts use XML-style context markers:
// <context>
// </context>
func BuildPrompt(template string, documents ...string) string {
    context := strings.Join(documents, "\n\n---\n\n")
    return strings.Replace(template, "</context>", context+"\n</context>", 1)
}
```

---

### 2. internal/claude/retry.go

```
☐ Fájl: internal/claude/retry.go
☐ Spec: l2/tech-specs.md TS-RETRY-001, TS-RETRY-002, DEC-L1-015
```

**Traceability komment:**
```go
// Package claude provides the Claude CLI integration layer.
//
// This file implements retry logic with exponential backoff.
// Implements: l2/tech-specs.md TS-RETRY-001, TS-RETRY-002
// Implements: DEC-L1-015 (retry configuration)
package claude
```

**Implementálandó:**

| Elem | Spec | Leírás |
|------|------|--------|
| `RetryConfig` struct | TS-RETRY-001 | MaxAttempts, BaseDelay, MaxDelay |
| `DefaultRetryConfig()` | DEC-L1-015 | 3 attempts, 2s base, 30s max |
| `WithRetry(cfg, fn)` | TS-RETRY-001 | Generic retry wrapper |
| `calculateDelay()` | TS-RETRY-001 | Exponential backoff |
| `isRetryableError()` | TS-RETRY-002 | Error classification |

**DefaultRetryConfig (DEC-L1-015):**
```go
// Implements: DEC-L1-015 (3 attempts, 2s base, 30s max)
func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   2 * time.Second,
        MaxDelay:    30 * time.Second,
    }
}
```

**Backoff Formula (TS-RETRY-001):**
```
delay = min(BaseDelay * 2^(attempt-1), MaxDelay)

Attempt 1: 2s
Attempt 2: 4s
Attempt 3: 8s (capped at 30s)
```

**Error Classification (TS-RETRY-002):**

| Category | Patterns | Action |
|----------|----------|--------|
| Retryable | rate limit, timeout, 503, 502, 500, overloaded, connection refused, connection reset, service unavailable, gateway timeout | Retry |
| Non-retryable | invalid api key, unauthorized, bad request, 404, not found, forbidden, 401, 403 | Fail immediately |
| Unknown | Anything else | Default: retry |

---

## Definition of Done

```
☐ internal/claude/client.go létezik
☐ Client struct: SessionID, Verbose fields
☐ NewClient() implemented
☐ Call() uses exec.Command("claude", ...)
☐ Call() sets CLAUDE_CODE_MAX_OUTPUT_TOKENS=100000
☐ CallWithSystemPrompt() implemented
☐ CallJSON() with extractJSON() fallback
☐ extractJSON() handles markdown blocks AND raw JSON
☐ sanitizeJSON() fixes newlines in strings
☐ BuildPrompt() injects into </context>
☐ internal/claude/retry.go létezik
☐ RetryConfig with correct fields
☐ DefaultRetryConfig() returns 3/2s/30s
☐ WithRetry() implements exponential backoff
☐ isRetryableError() classifies errors correctly
☐ Minden fájl tartalmaz traceability kommentet
☐ `go build` HIBA NÉLKÜL fut
```

---

## Test

```go
// Manual test (requires claude CLI installed)
client := claude.NewClient()
response, err := client.Call("Say 'Hello' and nothing else")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response) // Should print: Hello
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
