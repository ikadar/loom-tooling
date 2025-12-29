# Benchmark Teszt Készlet

> **Cél:** CLI promptok minőségének validálása különböző domaineken

## Struktúra

```
benchmark/
├── README.md
├── 01-ecommerce-order/      # Közepes komplexitás
│   ├── input-l0.md          # L0 input dokumentum
│   ├── expected-entities.json
│   ├── expected-ambiguities.json
│   └── expected-severity.json
├── 02-scheduling-calendar/   # Magas - temporal
├── 03-fintech-payment/       # Magas - state machine
├── 04-simple-crud-todo/      # Alacsony - baseline
└── 05-multitenant-saas/      # Magas - authorization
```

## Benchmark Dokumentumok

| # | Domain | Komplexitás | Fókusz |
|---|--------|-------------|--------|
| 1 | E-commerce (Order) | Közepes | Entity lifecycle, relationships |
| 2 | Scheduling (Calendar) | Magas | Temporal constraints, conflicts |
| 3 | Fintech (Payment) | Magas | State machine, error handling |
| 4 | Simple CRUD (Todo) | Alacsony | Basic coverage, baseline |
| 5 | Multi-tenant SaaS | Magas | Authorization, isolation |

## Használat

### Manuális Teszt

```bash
# 1. Futtatsd a CLI-t az input-l0.md-vel
loom analyze-entities benchmark/01-ecommerce-order/input-l0.md

# 2. Hasonlítsd össze az expected-*.json fájlokkal
```

### Értékelési Kritériumok

| Metrika | Pass Criteria |
|---------|---------------|
| Entity detection | >= 90% of expected entities found |
| Ambiguity count | >= minimum per severity |
| Critical coverage | 100% of expected critical ambiguities |
| False positive rate | < 20% irrelevant ambiguities |
| Edge case generation | >= 5 auto-generated per entity |

## Expected Output Formátumok

### expected-entities.json

```json
{
  "entities": [
    {
      "name": "Order",
      "type": "entity",
      "confidence": "high",
      "attributes_expected": ["id", "status", "total"],
      "relationships_expected": ["Customer", "OrderItem"]
    }
  ]
}
```

### expected-ambiguities.json

```json
{
  "minimum_ambiguities": [
    {
      "id_pattern": "AMB-ENT-*",
      "subject": "Order",
      "aspect_contains": "deletion",
      "severity": "critical"
    }
  ],
  "minimum_count": {
    "critical": 5,
    "important": 10,
    "minor": 5
  }
}
```

### expected-severity.json

```json
{
  "severity_distribution": {
    "critical_min": 5,
    "critical_max": 15,
    "important_min": 10,
    "important_max": 30,
    "minor_min": 5,
    "minor_max": 20
  }
}
```

## Milestone Kapcsolat

- **M4 (Platform Implementation):** Benchmark 1-5 létrehozása és manuális validálás
- **M6 (Beta & Iteration):** Automatizált teszt runner, A/B testing valós userekkel

---

*Létrehozva: 2025-12-29*
