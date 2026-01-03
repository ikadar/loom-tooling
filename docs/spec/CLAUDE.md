# Loom CLI Specification

Ez a mappa a loom-cli teljes specifikációját tartalmazza L0→L1→L2→L3 hierarchiában.

## Szereped

Te egy Go fejlesztő vagy, aki a specifikáció alapján implementálja a loom-cli-t.

**FONTOS:**
- Mielőtt bármit implementálsz, olvasd el az ÖSSZES spec dokumentumot!
- Készíts tervet és kérd a user jóváhagyását!
- Minden implementált kód TRACE-ELJEN a specifikációhoz!

---

## Implementációs Roadmap

### Fázis 1: Spec Megértése (NE UGORJ ÁT!)

**Olvasd el az ÖSSZES dokumentumot:**

```
L0 (Requirements):
□ l0/loom-cli.md              # User stories - MIT csinál a CLI
□ l0/nfr.md                   # Non-functional requirements
□ l0/domain-vocabulary.md     # Domain fogalmak
□ l0/decisions.md             # MIÉRT döntések (DEC-001 - DEC-031)

L1 (Strategic Design):
□ l1/domain-model.md          # Entitások, kapcsolatok
□ l1/acceptance-criteria.md   # Elfogadási kritériumok
□ l1/business-rules.md        # Üzleti szabályok
□ l1/bounded-context-map.md   # Kontextus határok
□ l1/decisions.md             # L1→L2 döntések (DEC-L1-001 - DEC-L1-017)

L2 (Tactical Design):
□ l2/interface-contracts.md   # CLI parancsok, opciók
□ l2/package-structure.md     # Go package struktúra
□ l2/internal-api.md          # Internal Go APIs
□ l2/aggregate-design.md      # Core Go structs
□ l2/tech-specs.md            # Algoritmusok, implementáció
□ l2/sequence-design.md       # Végrehajtási flow-k
□ l2/initial-data-model.md    # JSON sémák
□ l2/decisions.md             # L2→L3 döntések
□ l2/prompts/                 # 21 prompt file

L3 (Operational Design):
□ l3/test-cases.md            # Teszt esetek referencia
□ l3/implementation-skeletons.md  # Kód vázlatok
```

### Fázis 1b: Terv Készítése

1. Készíts implementációs tervet a spec alapján
2. Minden komponenshez írd meg melyik spec dokumentumból származik
3. **Kérd a user jóváhagyását a tervhez!**

### Fázis 2: Core Infrastructure

| Fájl | Spec Forrás |
|------|-------------|
| `main.go` + `go.mod` | l2/package-structure.md PKG-001 |
| `cmd/root.go` | l2/interface-contracts.md, l2/tech-specs.md TS-ARCH-001 |
| `internal/domain/types.go` | l2/internal-api.md, l2/aggregate-design.md |
| `internal/claude/client.go` | l2/internal-api.md, l2/tech-specs.md TS-ARCH-002 |
| `internal/claude/retry.go` | l2/tech-specs.md TS-RETRY-001/002, DEC-L1-015 |
| `prompts/prompts.go` | l2/internal-api.md PKG-011 |
| `prompts/*.md` | l2/prompts/ (mind a 21 fájl) |

### Fázis 3: Parancsok

| Fájl | Spec Forrás |
|------|-------------|
| `cmd/analyze.go` | IC-ANL-001, TS-ARCH-001a, SEQ-CAS-001 |
| `cmd/interview.go` | IC-INT-001, TS-ARCH-001c, SEQ-INT-001/002 |
| `cmd/derive_new.go` | IC-DRV-001, TS-ARCH-001b |
| `cmd/derive_l2.go` | IC-DRV-002, SEQ-DRV-001 |
| `cmd/derive_l3.go` | IC-DRV-003 |
| `cmd/validate.go` | IC-VAL-001, SEQ-VAL-001, BR-VAL-001/002/003 |
| `cmd/sync_links.go` | IC-SYN-001, SEQ-SYN-001 |
| `cmd/cascade.go` | IC-CAS-001, SEQ-CAS-001/002, AGG-CAS-001 |
| `cmd/interactive.go` | TS-ARCH-004, DEC-L1-016/017 |

### Fázis 4: Tesztelés

```
□ go build -o loom-cli .
□ Teszt: cascade --skip-interview (benchmark data)
□ Teszt: validate --level ALL
□ Teszt: interactive mode (-i flag)
□ Összehasonlítás: referencia output vs generált output
```

---

## Traceability Szabályok

**KÖTELEZŐ:** Minden implementált komponensnek KELL trace-elnie a spec-hez!

### Kód kommentek

```go
// Implements: IC-ANL-001, TS-ARCH-001a
// See: l2/interface-contracts.md, l2/tech-specs.md
func runAnalyze() error {
```

### Ha valami NINCS a spec-ben

1. **NE implementáld** önkényesen
2. Kérdezd meg a usert
3. Ha kell, először ADD HOZZÁ a spec-hez (decisions.md)
4. Csak utána implementáld

### Ha ellentmondást találsz

1. ÁLLJ MEG
2. Jelezd a usernek melyik dokumentumok mondanak ellent
3. Várd meg a tisztázást

---

## Spec Olvasási Sorrend

```
1. l2/interface-contracts.md     # CLI parancsok, opciók, exit codes
2. l2/package-structure.md       # Go package struktúra
3. l2/internal-api.md            # Internal Go APIs
4. l2/aggregate-design.md        # Core Go structs
5. l2/tech-specs.md              # Implementation algorithms
6. l2/sequence-design.md         # Execution flows
7. l2/prompts/                   # 21 prompt file (go:embed)
```

---

## Kritikus Döntések

| Döntés | Érték | Forrás |
|--------|-------|--------|
| Go verzió | 1.21+ | DEC-L1-014 |
| Dependencies | stdlib only | DEC-L1-001 |
| Claude integráció | `claude -p` CLI | DEC-L1-001 |
| Retry config | 3 attempt, 2s base, 30s max | DEC-L1-015 |
| Preview limits | 20 lines, 80 chars | DEC-L1-016 |
| Max output tokens | 100,000 | DEC-L1-002 |
| Max group size | 5 questions | l2/tech-specs.md |
| Chunk size | 5 ACs per batch | l2/internal-api.md |

---

## Teszt Adat

**Input (L0):**
```
test/benchmark/01-ecommerce-order/input-l0.md
```

**Referencia output:**
```
test/benchmark/01-ecommerce-order/
├── l1-anchor-test/          # L1 reference output
├── l2-interactive-test/     # L2 reference output
└── l3-test-20251229/        # L3 reference output
```

### Teszt parancsok

```bash
# Build
cd loom-cli && go build -o loom-cli .

# Full cascade test
./loom-cli cascade \
  --input-file ../test/benchmark/01-ecommerce-order/input-l0.md \
  --output-dir /tmp/test-output \
  --skip-interview

# Validation
./loom-cli validate --input-dir /tmp/test-output --level ALL

# Expected: Exit 0, no validation errors
```

---

## Spec Struktúra

```
l0/                          # User Stories, Requirements
├── loom-cli.md              # User stories (US-001 - US-009)
├── nfr.md                   # Non-functional requirements
├── domain-vocabulary.md     # Domain terms
└── decisions.md             # L0→L1 decisions (DEC-001 - DEC-031)

l1/                          # Strategic Design
├── domain-model.md          # Entities, value objects
├── bounded-context-map.md   # Context boundaries
├── acceptance-criteria.md   # AC-XXX-NNN
├── business-rules.md        # BR-XXX-NNN
└── decisions.md             # L1→L2 decisions (DEC-L1-001 - DEC-L1-017)

l2/                          # Tactical Design
├── tech-specs.md            # TS-XXX-NNN - Algorithms
├── interface-contracts.md   # IC-XXX-NNN - CLI interface
├── aggregate-design.md      # AGG-XXX-NNN - Go structs
├── sequence-design.md       # SEQ-XXX-NNN - Flows
├── initial-data-model.md    # JSON schemas
├── internal-api.md          # Go package APIs
├── package-structure.md     # Directory layout
├── build-test-guide.md      # Build/test commands
├── prompt-catalog.md        # Prompt overview
├── decisions.md             # L2→L3 decisions
└── prompts/                 # 21 prompt files

l3/                          # Operational Design
├── test-cases.md            # TC-AC-XXX-NNN-{T}NN
├── implementation-skeletons.md
├── feature-tickets.md
├── service-boundaries.md
├── event-message-design.md
├── dependency-graph.md
└── openapi-spec.md
```

---

## Traceability

Minden dokumentum trace-el a forrásához:
- L1 → L0 (user stories, decisions)
- L2 → L1 (acceptance criteria, business rules)
- L3 → L2 (tech specs, aggregates)

ID pattern-ek: `l2/tech-specs.md` TS-FMT-002 szekció.
