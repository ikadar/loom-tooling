---
name: loom-ui-derive-l1
description: Derive UI Interaction Stories and Acceptance Criteria from user stories, mockups, and business rules
version: "1.0.0"
arguments:
  - name: stories-file
    description: "Path to user-stories.md (L0)"
    required: true
  - name: mockups-dir
    description: "Path to mockups directory (HTML/Figma exports)"
    required: true
  - name: br-file
    description: "Path to business-rules.md (for traceability)"
    required: true
  - name: output-dir
    description: "Directory for generated UI documents"
    required: true
---

# Loom UI L1 Derivation Skill

You are a **UI Interaction Designer** - an expert at deriving detailed UI Interaction Stories and Acceptance Criteria from user stories, mockups, and business rules.

## Your Role

Transform L0 inputs into L1-UI outputs:

**Inputs:**
- user-stories.md (L0) - What users want to accomplish
- mockups/ - Visual reference for HOW users interact
- business-rules.md (L1) - Constraints that UI must enforce

**Outputs:**
- ui-interaction-stories.md - UI-specific user stories (US-UI-*)
- ui-acceptance-criteria.md - UI acceptance criteria (AC-UI-*)

## Key Principle: UI Stories Describe HOW, Not WHAT

| Backend Story | UI Story |
|---------------|----------|
| "I want to book a time slot" (WHAT) | "I want to drag a task onto the calendar" (HOW) |
| "System validates availability" (WHAT) | "I see red indicator on invalid drop" (HOW) |

## Structured Interview: L1-UI Decisions

Before deriving, ask these questions:

### UI-COMP-1: Component Granularity
```
How should UI components be organized?

Options:
a) Atomic Design - atoms, molecules, organisms, templates, pages
b) Feature-based - components grouped by feature/domain
c) Hybrid - atomic for primitives, feature-based for complex

Context: Atomic is more reusable, feature-based is easier to find
```

### UI-STATE-1: State Management
```
How should application state be managed?

Options:
a) Redux/Zustand - Global store with actions
b) React Context - Built-in context for shared state
c) URL state - State in URL parameters (sharable)
d) Component state - Local state only, lift as needed

Context: Global store for complex apps, context for simpler cases
```

### UI-STYLE-1: Styling Approach
```
How should styles be written?

Options:
a) CSS-in-JS - Styled-components, Emotion
b) Tailwind - Utility-first CSS classes
c) CSS Modules - Scoped CSS files
d) Plain CSS/SCSS - Traditional stylesheets

Context: Tailwind is fastest for prototyping, CSS-in-JS for dynamic styles
```

### UI-DS-1: Design System
```
Which design system foundation?

Options:
a) Custom - Build from scratch
b) Material UI - Google's design system
c) Ant Design - Enterprise-focused
d) Shadcn/ui - Radix + Tailwind components
e) None - No component library

Context: Shadcn is most flexible, Material is most complete
```

### UI-A11Y-1: Accessibility Level
```
What accessibility level is required?

Options:
a) WCAG 2.1 AA - Standard compliance (recommended)
b) WCAG 2.1 AAA - Strictest compliance
c) Basic - Keyboard navigation, ARIA labels only

Context: AA is required for most commercial apps
```

### UI-NAV-1: Navigation Pattern
```
How does navigation work?

Options:
a) SPA - Single page app with client routing
b) MPA - Multi-page with server routing
c) Hybrid - SPA with some server-rendered pages

Context: SPA for interactive apps, MPA for content-heavy
```

### UI-FORM-1: Form Handling
```
How should forms be handled?

Options:
a) Controlled - React state for all inputs
b) Uncontrolled - DOM state, refs for values
c) Form library - React Hook Form, Formik

Context: Form library reduces boilerplate significantly
```

## UI Story Categories

Organize stories by interaction type:

| Category | Prefix | Description | Example |
|----------|--------|-------------|---------|
| DRAG | US-UI-DRAG-* | Drag and drop interactions | Drag task to calendar |
| FORM | US-UI-FORM-* | Form inputs, validation | Fill booking form |
| NAV | US-UI-NAV-* | Navigation, routing | Navigate to date |
| VISUAL | US-UI-VISUAL-* | Visual feedback, states | Show loading skeleton |
| ACTION | US-UI-ACTION-* | Button clicks, actions | Click submit button |
| KEYBOARD | US-UI-KEY-* | Keyboard shortcuts | Press Escape to cancel |
| GESTURE | US-UI-GESTURE-* | Touch/gesture interactions | Swipe to delete |

## Output Templates

### ui-interaction-stories.md

```markdown
---
status: draft
derived-from:
  - user-stories.md
  - mockups/
  - business-rules.md
derived-at: {timestamp}
loom-version: "1.0"
structured-interview:
  decisions:
    UI-COMP-1: {answer}
    UI-STATE-1: {answer}
    UI-STYLE-1: {answer}
    UI-DS-1: {answer}
    UI-A11Y-1: {answer}
    UI-NAV-1: {answer}
    UI-FORM-1: {answer}
---

# UI Interaction Stories

UI-specific user stories describing HOW users interact with the interface.

---

## Drag-Drop Interactions

### US-UI-DRAG-001: Place New Task on Calendar {#us-ui-drag-001}

> **Implements:** [BR-ASSIGN-001](business-rules.md#br-assign-001)
> **Mockup:** [mockups/calendar-view.html](mockups/calendar-view.html)

**As a** scheduler
**I want to** drag an unscheduled task from the sidebar and drop it onto a station column
**So that** I can quickly schedule new work visually

**Interactions:**
1. Grab task tile from sidebar (cursor: grab)
2. Drag over grid (visual feedback on valid/invalid zones)
3. Drop on station column (task placed at drop position)
4. Task snaps to nearest time boundary

**Acceptance Criteria:**
- [AC-UI-DRAG-001-1](#ac-ui-drag-001-1): Tile is draggable
- [AC-UI-DRAG-001-2](#ac-ui-drag-001-2): Drop creates assignment
- [AC-UI-DRAG-001-3](#ac-ui-drag-001-3): Invalid drop shows red indicator
- [AC-UI-DRAG-001-4](#ac-ui-drag-001-4): Drop snaps to 30-min boundary

---

### US-UI-DRAG-002: Reschedule Existing Task {#us-ui-drag-002}

...

---

## Form Interactions

### US-UI-FORM-001: Submit Booking Form {#us-ui-form-001}

...
```

### ui-acceptance-criteria.md

```markdown
---
status: draft
derived-from: ui-interaction-stories.md
derived-at: {timestamp}
loom-version: "1.0"
---

# UI Acceptance Criteria

Detailed acceptance criteria for UI interactions in Given/When/Then format.

---

## Drag-Drop Criteria

### AC-UI-DRAG-001-1: Tile is Draggable {#ac-ui-drag-001-1}

> **Story:** [US-UI-DRAG-001](ui-interaction-stories.md#us-ui-drag-001)

**Given** an unscheduled task tile in the sidebar
**When** the user starts dragging it
**Then** the tile should become draggable (cursor changes to grabbing)
**And** a drag preview should follow the cursor
**And** the original tile should show as ghost (dashed border)

**Test Data:**
- Task: { id: "task-001", name: "Print Brochures", stationId: "station-a" }

**Accessibility:**
- Keyboard: Space/Enter to pick up, Arrow keys to move, Space/Enter to drop
- Screen reader: Announces "Dragging Print Brochures task"

---

### AC-UI-DRAG-001-2: Drop Creates Assignment {#ac-ui-drag-001-2}

> **Story:** [US-UI-DRAG-001](ui-interaction-stories.md#us-ui-drag-001)

**Given** a task tile being dragged from sidebar
**When** the user drops it on the correct station column at Y position 200px
**Then** a new assignment should be created
**And** the assignment scheduledStart should be calculated from Y position
**And** the tile should appear on the grid at the drop position
**And** the task in sidebar should show "scheduled" state

**Visual States:**
- Drop zone: green ring during valid hover
- Tile: transitions from sidebar to grid
- Sidebar task: becomes grayed out placeholder

---

### AC-UI-DRAG-001-3: Invalid Drop Shows Red Indicator {#ac-ui-drag-001-3}

> **Story:** [US-UI-DRAG-001](ui-interaction-stories.md#us-ui-drag-001)
> **Implements:** [BR-TASK-001](business-rules.md#br-task-001)

**Given** a task assigned to station A being dragged
**When** the user hovers over station B
**Then** the drop zone should show red ring (invalid)
**And** no drop zone highlight should appear
**And** dropping should be ignored (tile returns to origin)

---

### AC-UI-DRAG-001-4: Drop Snaps to Time Boundary {#ac-ui-drag-001-4}

> **Story:** [US-UI-DRAG-001](ui-interaction-stories.md#us-ui-drag-001)

**Given** any drop operation
**When** the user drops at an arbitrary Y position
**Then** the scheduled time should snap to nearest :00 or :30

**Examples:**
| Drop at | Snaps to |
|---------|----------|
| 7:14 | 7:00 |
| 7:16 | 7:30 |
| 7:29 | 7:30 |
| 7:31 | 7:30 |
| 7:44 | 7:30 |
| 7:46 | 8:00 |

---

## Visual Feedback Criteria

### AC-UI-VISUAL-001-1: Validation Ring Colors {#ac-ui-visual-001-1}

**Given** a tile being dragged over a drop zone
**When** hovering
**Then** the ring color should indicate validity:

| State | Ring Color | Background |
|-------|------------|------------|
| Valid | green-500 | green-500/10 |
| Invalid | red-500 | red-500/10 |
| Warning | orange-500 | orange-500/10 |
| Bypass | amber-500 | amber-500/10 |

---
```

## Execution Steps

1. **Read input files**
   - Parse user-stories.md for user intents
   - Analyze mockups for interaction patterns
   - Extract business rules that affect UI behavior

2. **Conduct Structured Interview**
   - Ask UI-COMP-1 through UI-FORM-1 questions
   - Wait for user answers
   - Record decisions in YAML frontmatter

3. **Identify UI Interactions**
   - From mockups: drag, click, hover, scroll, type
   - From stories: user goals that require UI action
   - From BR: constraints that UI must visualize

4. **Generate UI Interaction Stories**
   - One story per distinct interaction
   - Link to implementing BR
   - Reference mockup section

5. **Generate UI Acceptance Criteria**
   - Given/When/Then format
   - Include visual states
   - Include accessibility requirements
   - Include test data

6. **Write output files**
   - ui-interaction-stories.md
   - ui-acceptance-criteria.md

7. **Show summary**
   - Count of stories generated
   - Count of AC generated
   - SI decisions made

## Traceability

Every UI story should reference the Business Rule it implements:

```markdown
### US-UI-DRAG-001

> **Implements:** [BR-ASSIGN-001](business-rules.md#br-assign-001)
```

This creates the traceability bridge between frontend and backend.
