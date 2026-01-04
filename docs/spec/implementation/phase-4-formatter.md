# Phase 4: Formatter Package

## Cél

Implementáld az `internal/formatter/` package-et: JSON→Markdown konverziók.

**FONTOS:** Ez 9 fájl! Ne hagyj ki egyet sem!

---

## Spec Források (OLVASD EL ELŐBB!)

| Dokumentum | Szekció | Mit definiál |
|------------|---------|--------------|
| l2/package-structure.md | PKG-007 | Package structure, 9 files |
| l2/internal-api.md | internal/formatter | Type definitions, Format functions |

---

## Implementálandó Fájlok (9 db!)

### 1. internal/formatter/types.go

```
☐ Fájl: internal/formatter/types.go
☐ Spec: l2/package-structure.md PKG-007, l2/internal-api.md
```

**Traceability:**
```go
// Package formatter provides JSON to Markdown formatting.
//
// Implements: l2/package-structure.md PKG-007
// See: l2/internal-api.md
package formatter
```

**Típusok (l2/internal-api.md):**

```go
// Test Case types
type TestCase struct {
    ID              string     `json:"id"`
    Name            string     `json:"name"`
    Category        string     `json:"category"` // positive, negative, boundary, hallucination
    ACRef           string     `json:"ac_ref"`
    BRRefs          []string   `json:"br_refs"`
    Preconditions   []string   `json:"preconditions"`
    TestData        []TestData `json:"test_data"`
    Steps           []string   `json:"steps"`
    ExpectedResults []string   `json:"expected_results"`
    ShouldNot       string     `json:"should_not,omitempty"`
}

type TestData struct {
    Field string      `json:"field"`
    Value interface{} `json:"value"`
    Notes string      `json:"notes"`
}

type TDAISummary struct {
    Total      int
    ByCategory struct {
        Positive, Negative, Boundary, Hallucination int
    }
    Coverage struct {
        ACsCovered            int
        PositiveRatio         float64
        NegativeRatio         float64
        HasHallucinationTests bool
    }
}

// Tech Spec types
type TechSpec struct {
    ID               string          `json:"id"`
    Name             string          `json:"name"`
    BRRef            string          `json:"br_ref"`
    Rule             string          `json:"rule"`
    Implementation   string          `json:"implementation"`
    ValidationPoints []string        `json:"validation_points"`
    DataRequirements []DataReq       `json:"data_requirements"`
    ErrorHandling    []ErrorHandling `json:"error_handling"`
    RelatedACs       []string        `json:"related_acs"`
}

type DataReq struct {
    Field       string `json:"field"`
    Type        string `json:"type"`
    Constraints string `json:"constraints"`
    Source      string `json:"source"`
}

type ErrorHandling struct {
    Condition  string `json:"condition"`
    ErrorCode  string `json:"error_code"`
    Message    string `json:"message"`
    HTTPStatus int    `json:"http_status"`
}

// Interface Contract types
type InterfaceContract struct {
    ID                   string               `json:"id"`
    ServiceName          string               `json:"serviceName"`
    Purpose              string               `json:"purpose"`
    BaseURL              string               `json:"baseUrl"`
    Operations           []ContractOperation  `json:"operations"`
    Events               []ContractEvent      `json:"events"`
    SecurityRequirements SecurityRequirements `json:"securityRequirements"`
}

// Aggregate Design types
type AggregateDesign struct {
    ID                 string           `json:"id"`
    Name               string           `json:"name"`
    Purpose            string           `json:"purpose"`
    Invariants         []AggInvariant   `json:"invariants"`
    Root               AggRoot          `json:"root"`
    Entities           []AggEntity      `json:"entities"`
    ValueObjects       []string         `json:"valueObjects"`
    Behaviors          []AggBehavior    `json:"behaviors"`
    Events             []AggEvent       `json:"events"`
    Repository         AggRepository    `json:"repository"`
    ExternalReferences []AggExternalRef `json:"externalReferences"`
}

// Sequence Design types
type SequenceDesign struct {
    ID           string              `json:"id"`
    Name         string              `json:"name"`
    Description  string              `json:"description"`
    Trigger      SequenceTrigger     `json:"trigger"`
    Participants []SeqParticipant    `json:"participants"`
    Steps        []SequenceStep      `json:"steps"`
    Outcome      SequenceOutcome     `json:"outcome"`
    Exceptions   []SequenceException `json:"exceptions"`
    RelatedACs   []string            `json:"relatedACs"`
    RelatedBRs   []string            `json:"relatedBRs"`
}

// Data Model types
type DataTable struct {
    ID               string           `json:"id"`
    Name             string           `json:"name"`
    Aggregate        string           `json:"aggregate"`
    Purpose          string           `json:"purpose"`
    Fields           []DataField      `json:"fields"`
    PrimaryKey       DataPrimaryKey   `json:"primaryKey"`
    Indexes          []DataIndex      `json:"indexes"`
    ForeignKeys      []DataForeignKey `json:"foreignKeys"`
    CheckConstraints []DataConstraint `json:"checkConstraints"`
}

type DataEnum struct {
    Name   string   `json:"name"`
    Values []string `json:"values"`
}
```

---

### 2. internal/formatter/frontmatter.go

```
☐ Fájl: internal/formatter/frontmatter.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatFrontmatter generates YAML frontmatter for documents.
//
// Implements: l2/package-structure.md PKG-007
func FormatFrontmatter(title string, level string) string
```

**Output:**
```yaml
---
title: "Document Title"
generated: 2024-01-15T10:30:00Z
status: draft
level: L2
---
```

---

### 3. internal/formatter/anchor.go

```
☐ Fájl: internal/formatter/anchor.go
☐ Spec: l2/package-structure.md PKG-007
```

**Functions:**
```go
// FormatAnchor creates a markdown anchor link.
func FormatAnchor(id string) string

// FormatReference creates a reference to another document.
func FormatReference(id string, docPath string) string

// FormatTraceability creates a traceability section.
func FormatTraceability(sources []string, decisions []string) string
```

---

### 4. internal/formatter/techspecs.go

```
☐ Fájl: internal/formatter/techspecs.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatTechSpecs formats tech specs as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: tech-specs.md
func FormatTechSpecs(specs []TechSpec) string
```

**Output Format:**
```markdown
## TS-BR-XXX-001 – Spec Name

**Business Rule:** BR-XXX-001

**Rule:** The rule statement

**Implementation:** How to implement

**Validation Points:**
- Point 1
- Point 2

**Data Requirements:**
| Field | Type | Constraints | Source |
|-------|------|-------------|--------|

**Error Handling:**
| Condition | Code | Message | HTTP |
|-----------|------|---------|------|

**Traceability:**
- AC: AC-XXX-001
```

---

### 5. internal/formatter/testcases.go

```
☐ Fájl: internal/formatter/testcases.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatTestCases formats test cases as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: test-cases.md
func FormatTestCases(suites []TestSuite, summary TDAISummary) string
```

---

### 6. internal/formatter/contracts.go

```
☐ Fájl: internal/formatter/contracts.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatInterfaceContracts formats contracts as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: interface-contracts.md
func FormatInterfaceContracts(contracts []InterfaceContract, sharedTypes []SharedType) string
```

---

### 7. internal/formatter/aggregates.go

```
☐ Fájl: internal/formatter/aggregates.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatAggregateDesign formats aggregate design as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: aggregate-design.md
func FormatAggregateDesign(aggregates []AggregateDesign) string
```

---

### 8. internal/formatter/sequences.go

```
☐ Fájl: internal/formatter/sequences.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatSequenceDesign formats sequences as markdown with Mermaid diagrams.
//
// Implements: l2/package-structure.md PKG-007
// Output: sequence-design.md
func FormatSequenceDesign(sequences []SequenceDesign) string
```

**Mermaid Output:**
```markdown
```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Service

    Client->>API: POST /resource
    API->>Service: validate()
    Service-->>API: OK
    API-->>Client: 201 Created
```​
```

---

### 9. internal/formatter/datamodel.go

```
☐ Fájl: internal/formatter/datamodel.go
☐ Spec: l2/package-structure.md PKG-007
```

**Function:**
```go
// FormatDataModel formats data model as markdown.
//
// Implements: l2/package-structure.md PKG-007
// Output: initial-data-model.md
func FormatDataModel(tables []DataTable, enums []DataEnum) string
```

---

## Definition of Done

```
☐ internal/formatter/types.go - ALL types defined
☐ internal/formatter/frontmatter.go - FormatFrontmatter()
☐ internal/formatter/anchor.go - FormatAnchor(), FormatReference(), FormatTraceability()
☐ internal/formatter/techspecs.go - FormatTechSpecs()
☐ internal/formatter/testcases.go - FormatTestCases()
☐ internal/formatter/contracts.go - FormatInterfaceContracts()
☐ internal/formatter/aggregates.go - FormatAggregateDesign()
☐ internal/formatter/sequences.go - FormatSequenceDesign() with Mermaid
☐ internal/formatter/datamodel.go - FormatDataModel()
☐ 9 fájl összesen (nem 8, nem 7, KILENC!)
☐ Minden fájl tartalmaz traceability kommentet
☐ `go build` HIBA NÉLKÜL fut
```

---

## NE LÉPJ TOVÁBB amíg MIND A 9 FÁJL KÉSZ!
