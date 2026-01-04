# Phase 7: Prompt Content

## Cél

Implementáld az összes prompt markdown fájlt a `prompts/` könyvtárban.

**FONTOS:** Minden prompt KÖTELEZŐEN tartalmazza a spec szerinti szekciókat!

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/prompt-catalog.md | Prompt Structure Standard | Kötelező szekciók |
| l2/prompt-catalog.md | PRM-* | Minden prompt specifikációja |
| l2/prompts/*.md | Reference | Meglévő prompt fájlok (ha vannak) |

---

## Kötelező Prompt Struktúra

**l2/prompt-catalog.md szerint MINDEN prompt tartalmazza:**

```markdown
<role>
Expert persona with specific experience areas
Priority list (what to optimize for)
Approach description
</role>

<task>
Clear objective statement
What to produce
</task>

<thinking_process>
Step-by-step analysis instructions
What to consider before generating output
</thinking_process>

<instructions>
Detailed requirements
Format specifications
Coverage requirements
</instructions>

<output_format>
CRITICAL REQUIREMENTS for JSON output
JSON Schema with examples
Field descriptions
</output_format>

<examples>
Named examples with description
Input → Analysis → Output pattern
Multiple complexity levels (at least 2-3)
</examples>

<self_review>
Completeness checks
Consistency checks
Format checks
Fix instructions
</self_review>

<critical_output_format>
Final reminder: PURE JSON ONLY
Start/end requirements
</critical_output_format>

<context>
</context>
```

---

## Implementálandó Prompt Fájlok (21 db!)

### Analyze Phase (3 prompt)

```
☐ prompts/domain-discovery.md     PRM-ANL-001
☐ prompts/entity-analysis.md      PRM-ANL-002
☐ prompts/operation-analysis.md   PRM-ANL-003
```

### Interview Phase (1 prompt)

```
☐ prompts/interview.md            PRM-INT-001
```

### L1 Derivation (3 prompt)

```
☐ prompts/derivation.md           PRM-DRV-001
☐ prompts/derive-domain-model.md  PRM-DRV-002
☐ prompts/derive-bounded-context.md PRM-DRV-003
```

### L2 Derivation (6 prompt)

```
☐ prompts/derive-l2.md                  PRM-L2-001
☐ prompts/derive-tech-specs.md          PRM-L2-002
☐ prompts/derive-interface-contracts.md PRM-L2-003
☐ prompts/derive-aggregate-design.md    PRM-L2-004
☐ prompts/derive-sequence-design.md     PRM-L2-005
☐ prompts/derive-data-model.md          PRM-L2-006
```

### L3 Derivation (8 prompt)

```
☐ prompts/derive-test-cases.md          PRM-L3-001
☐ prompts/derive-l3.md                  PRM-L3-COMBINED
☐ prompts/derive-l3-api.md              PRM-L3-002
☐ prompts/derive-l3-skeletons.md        PRM-L3-003
☐ prompts/derive-feature-tickets.md     PRM-L3-004
☐ prompts/derive-service-boundaries.md  PRM-L3-005
☐ prompts/derive-event-design.md        PRM-L3-006
☐ prompts/derive-dependency-graph.md    PRM-L3-007
```

---

## Prompt Minőségi Követelmények

### 1. Examples Section (KÖTELEZŐ!)

Minden prompt tartalmazzon **legalább 2 példát**:

```markdown
<examples>
<example name="simple_case" description="Basic scenario">
Input: ...
Analysis:
- Point 1
- Point 2
Output:
{...}
</example>

<example name="complex_case" description="Edge case handling">
Input: ...
Analysis:
- ...
Output:
{...}
</example>
</examples>
```

### 2. Self-Review Section (KÖTELEZŐ!)

```markdown
<self_review>
After generating output, verify these criteria. If any fail, fix before outputting:

COMPLETENESS CHECK:
- [ ] Every input item has corresponding output
- [ ] No items were skipped

CONSISTENCY CHECK:
- [ ] All IDs follow pattern
- [ ] All references exist

FORMAT CHECK:
- [ ] JSON is valid
- [ ] Starts with { character
- [ ] No trailing commas

If issues found, fix before outputting.
</self_review>
```

### 3. String Length Limits (l2/prompt-catalog.md)

| Context | Max Length |
|---------|------------|
| Test case names | 60 chars |
| Service descriptions | 80 chars |
| Error messages | 100 chars |

Include in `<output_format>`:
```markdown
STRING LENGTH LIMITS:
- name: max 60 characters
- description: max 80 characters
- message: max 100 characters
```

---

## Prompt-Specifikus Követelmények

### derive-test-cases.md (PRM-L3-001)

**TDAI Coverage Requirements:**
- Positive tests: 2+ per AC
- Negative tests: 2+ per AC
- Boundary tests: 1+ per AC
- Hallucination tests: 1+ per AC
- Negative ratio: >= 30%

**Self-review checklist:**
```markdown
- [ ] Every AC has at least 6 tests (2P + 2N + 1B + 1H)
- [ ] All four categories represented per AC
- [ ] All IDs unique (TC-AC-XXX-NNN-{P|N|B|H}NN)
- [ ] Negative ratio >= 30%
- [ ] Every test has source_quote
```

### domain-discovery.md (PRM-ANL-001)

**Decision Points (l2/prompt-catalog.md):**
- EVO-1 through EVO-5: Entity vs Value Object
- AGG-1 through AGG-4: Aggregate boundaries
- REF-1 through REF-3: Reference type

**Output includes:**
```json
{
  "decision_points": {
    "EVO-1": {"answer": "yes|no|unclear", "evidence": "..."}
  },
  "needs_interview": true,
  "interview_questions": ["EVO-1", "AGG-2"]
}
```

---

## prompts/prompts.go Frissítés

Ellenőrizd, hogy `prompts.go` tartalmazza az összes promptot:

```go
var (
    // Analyze phase (3)
    DomainDiscovery   string
    EntityAnalysis    string
    OperationAnalysis string

    // Interview phase (1)
    Interview string

    // L1 derivation (3)
    Derivation           string
    DeriveDomainModel    string
    DeriveBoundedContext string

    // L2 derivation (6)
    DeriveL2                 string
    DeriveTechSpecs          string
    DeriveInterfaceContracts string
    DeriveAggregateDesign    string
    DeriveSequenceDesign     string
    DeriveDataModel          string

    // L3 derivation (8)
    DeriveTestCases         string
    DeriveL3                string
    DeriveL3API             string
    DeriveL3Skeletons       string
    DeriveFeatureTickets    string
    DeriveServiceBoundaries string
    DeriveEventDesign       string
    DeriveDependencyGraph   string
)
```

---

## Definition of Done

```
☐ 21 prompt fájl létezik a prompts/ könyvtárban
☐ MINDEN prompt tartalmaz <role> szekciót
☐ MINDEN prompt tartalmaz <task> szekciót
☐ MINDEN prompt tartalmaz <thinking_process> szekciót
☐ MINDEN prompt tartalmaz <instructions> szekciót
☐ MINDEN prompt tartalmaz <output_format> szekciót
☐ MINDEN prompt tartalmaz <examples> szekciót (min. 2 példa)
☐ MINDEN prompt tartalmaz <self_review> szekciót
☐ MINDEN prompt tartalmaz <critical_output_format> szekciót
☐ MINDEN prompt tartalmaz <context></context> szekciót
☐ derive-test-cases.md TDAI követelmények benne
☐ domain-discovery.md Decision Points benne
☐ prompts/prompts.go betölti mind a 21 promptot
☐ `go build` HIBA NÉLKÜL fut
```

---

## NE LÉPJ TOVÁBB amíg a checklist MINDEN eleme kész!
