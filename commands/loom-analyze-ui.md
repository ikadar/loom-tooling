---
name: loom-analyze-ui
description: Analyze UI/UX for completeness and generate ambiguity list
version: "1.0.0"
arguments:
  - name: input
    description: "Path to UI/UX specification file OR indication that UI is mentioned in domain input"
    required: true
  - name: domain-input
    description: "Path to domain input file (to check if UI is mentioned but not specified)"
    required: false
  - name: output-format
    description: "Output format: 'ambiguities' (default) or 'full-report'"
    required: false
---

# UI/UX Completeness Analysis

Analyze UI/UX specifications for completeness using the UI checklist.

## Reference

Use checklist from: `.claude/docs/checklists/ui-checklist.md`

## Process

### Step 0: UI Input Detection

**CRITICAL:** First check if UI/UX specification exists.

If `domain-input` is provided:
1. Scan for UI-related mentions (screens, components, interactions, layout)
2. Check if detailed UI spec exists

```markdown
## UI Input Assessment

| Aspect | Mentioned in Domain? | Detailed Spec Exists? |
|--------|---------------------|----------------------|
| Screens/pages | ✓ (3 mentioned) | ✗ |
| Components | ✓ (5 mentioned) | ✗ |
| Interactions | ✓ (drag-drop, filter) | ✗ |
| States | ✗ | ✗ |
| Error handling | ✗ | ✗ |
| Loading states | ✗ | ✗ |
| Responsive | ✗ | ✗ |
| Accessibility | ✗ | ✗ |
```

**If UI is mentioned but no detailed spec:**

```markdown
## ⚠️ UI/UX Specification Missing

UI components are mentioned in the domain input but detailed specification is missing.

**Mentioned UI elements:**
- Scheduling Grid (drag-drop interaction)
- Left Panel (job list, search)
- Right Panel (late jobs)
- Similarity Indicators (●/○)

**To proceed, please provide:**
1. Screen/page list with purpose and layout
2. Component breakdown with states
3. Interaction specifications (triggers, feedback, outcomes)
4. Error, loading, and empty state definitions
5. Responsive behavior requirements
6. Accessibility requirements

**Options:**
a) Provide UI/UX specification file
b) Answer UI questions interactively
c) Mark UI as out of scope for this derivation
```

### Step 1: Extract UI Elements

From UI specification, extract:

```yaml
screens:
  - name: "Scheduling Page"
    purpose: "Main scheduling interface"
    components: [header, left_panel, grid, right_panel]

components:
  - name: "Scheduling Grid"
    location: ["Scheduling Page"]
    interactions: [drag_drop, click, scroll]
    data_source: "tasks, stations, assignments"

interactions:
  - name: "Drag Task to Grid"
    trigger: "mouse drag"
    component: "Task Tile"
    target: "Station Column"
```

### Step 2: Apply Checklist

For EACH element:

#### Screens (Section 1)
- Basic definition (purpose, entry/exit, URL, permissions)
- Layout (type, breakpoints, components, scroll, sticky)
- Initial state

#### Components (Section 2)
- Basic definition
- Visual states (12 states: default, hover, focus, active, disabled, loading, error, success, empty, selected, dragging, drop target)
- Content states (no data, loading, partial, error, stale)
- Sizing (width, height, overflow)

#### Interactions (Section 3)
- Triggers (mouse, keyboard, touch, programmatic)
- Feedback during (start, progress, valid, invalid)
- Outcomes (success, partial, failure, cancelled)
- Modifiers (Shift, Ctrl, Alt)
- Constraints (valid targets, timing)
- Edge cases

#### Navigation (Section 4)
- Global nav
- In-page nav
- Transitions
- Deep linking

#### Cross-cutting (Section 5)
- Loading patterns
- Error patterns
- Empty state patterns
- Notification patterns

#### Accessibility (Section 6)
- Keyboard navigation
- Screen reader
- Color contrast
- Touch targets

#### Responsive (Section 7)
- Each breakpoint behavior

### Step 3: Generate Ambiguity List

```yaml
ambiguities:
  - id: "AMB-UI-001"
    element: "Scheduling Grid"
    category: "interaction"
    aspect: "drag.feedback"
    question: "What visual feedback during task drag?"
    severity: "important"

  - id: "AMB-UI-002"
    element: "Task Tile"
    category: "state"
    aspect: "loading"
    question: "What does a task tile look like while loading?"
    severity: "important"

  - id: "AMB-UI-003"
    element: "Left Panel"
    category: "responsive"
    aspect: "mobile"
    question: "How does left panel behave on mobile?"
    severity: "important"
```

### Step 4: Severity Classification

| Severity | Criteria |
|----------|----------|
| **critical** | Core interaction undefined, blocks implementation |
| **important** | State or feedback undefined, affects UX |
| **minor** | Polish details, has sensible default |

## Output Format

### Ambiguities Only (default)

```markdown
## UI/UX Analysis: Ambiguities Found

**Screens Analyzed:** 1
**Components Analyzed:** 8
**Interactions Analyzed:** 5
**Total Ambiguities:** 52

### Critical (8)

| ID | Element | Aspect | Question |
|----|---------|--------|----------|
| AMB-UI-001 | Grid | drop_behavior | What happens on drop? |
| ... |

### Important (31)

| ID | Element | Aspect | Question |
|----|---------|--------|----------|
| ... |

### Minor (13)

| ID | Element | Aspect | Question | Suggested Default |
|----|---------|--------|----------|-------------------|
| ... |
```

### UI Specification Missing

```markdown
## ⚠️ UI/UX Analysis Cannot Proceed

**Status:** UI elements mentioned but specification missing

**Action Required:** Provide UI/UX specification or mark as out of scope

**Mentioned Elements:**
- Scheduling Grid
- Left Panel
- Right Panel
- Task Tiles
- Similarity Indicators

**Missing Specifications:**
- Component state definitions (12 states each)
- Interaction specifications
- Error/loading/empty states
- Responsive behavior
- Accessibility requirements
```

## Quality Criteria

Before outputting, verify:
- [ ] UI input existence checked
- [ ] Every screen analyzed
- [ ] Every component has all 12 visual states checked
- [ ] Every interaction has triggers, feedback, outcomes checked
- [ ] Navigation fully specified
- [ ] Cross-cutting patterns defined
- [ ] Accessibility checked
- [ ] Responsive behavior per breakpoint checked
- [ ] Severities correctly assigned
