# Loom CLI Specification

Ez a mappa a loom-cli teljes specifikációját tartalmazza L0→L1→L2→L3 hierarchiában.

## Szereped

Te egy Go fejlesztő vagy, aki a specifikáció alapján implementálja a loom-cli-t.

**KRITIKUS SZABÁLYOK:**
1. Fázisonként haladj - NE ugorj előre!
2. Minden fázishoz van részletes prompt az `implementation/` mappában
3. Minden fázis végén ellenőrizd a checklistet
4. Minden kód TRACE-ELJEN a specifikációhoz!

---

## Implementációs Fázisok

**Master Checklist:** `implementation/checklist.md`

| Fázis | Prompt | Mit implementál |
|-------|--------|-----------------|
| 1 | `implementation/phase-1-core.md` | main.go, cmd/root.go, domain/types.go |
| 2 | `implementation/phase-2-claude.md` | claude/client.go, claude/retry.go |
| 3 | `implementation/phase-3-commands.md` | 10 cmd/*.go fájl |
| 4 | `implementation/phase-4-formatter.md` | 9 formatter fájl |
| 5 | `implementation/phase-5-generator.md` | generator/testcases.go, parallel.go |
| 6 | `implementation/phase-6-workflow.md` | workflow, config, checkpoint, interview |
| 7 | `implementation/phase-7-prompts.md` | 21 prompt fájl |

### Hogyan használd

1. **Olvasd el** az adott fázis promptját
2. **Implementáld** az ott felsorolt fájlokat
3. **Ellenőrizd** a Definition of Done checklistet
4. **Futtasd** `go build` - HIBA NÉLKÜL kell fusson
5. **Jelöld meg** a checklist.md-ben
6. **Lépj tovább** a következő fázisra

### Spec Olvasási Sorrend (fázisonként!)

Minden fázis promptja megmondja mely spec dokumentumokat kell előtte elolvasni.
A legfontosabbak:

```
l2/interface-contracts.md     # CLI parancsok, opciók, exit codes
l2/package-structure.md       # Go package struktúra
l2/internal-api.md            # Internal Go APIs
l2/tech-specs.md              # Implementation algorithms
l2/prompt-catalog.md          # Prompt struktúra szabvány
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
