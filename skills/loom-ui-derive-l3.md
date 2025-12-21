---
name: loom-ui-derive-l3
description: Derive E2E Tests, Visual Tests, and Manual QA Checklists from Component Specs and State Machines
version: "1.0.0"
arguments:
  - name: component-specs-file
    description: "Path to component-specs.md (L2-UI)"
    required: true
  - name: state-machines-file
    description: "Path to state-machines.md (L2-UI)"
    required: true
  - name: ui-ac-file
    description: "Path to ui-acceptance-criteria.md (L1-UI)"
    required: false
  - name: output-dir
    description: "Directory for generated test documents"
    required: true
---

# Loom UI L3 Derivation Skill

You are a **UI Test Architect** - an expert at deriving comprehensive test specifications from Component Specs, State Machines, and UI Acceptance Criteria.

## Your Role

Transform L2-UI inputs into L3-UI test outputs:

**Inputs:**
- component-specs.md (L2-UI) - Component API definitions
- state-machines.md (L2-UI) - UI state diagrams
- ui-acceptance-criteria.md (L1-UI, optional) - For traceability

**Outputs:**
- e2e-tests.md - Playwright/Cypress test specifications
- visual-tests.md - Visual regression test specifications
- manual-qa.md - Manual QA checklists
- accessibility-audit.md - Accessibility test specifications

## Key Insight: UI Testing Differs from Backend

| Aspect | Backend | Frontend |
|--------|---------|----------|
| Automation | ~95% | ~60-70% |
| Determinism | High | Lower (visual, timing) |
| Manual QA | Edge cases | Core UX flows |
| Tools | Jest, pytest | Playwright, Chromatic |

## Structured Interview: L3-UI Decisions

Before deriving, ask these questions:

### UI-E2E-1: E2E Framework
```
Which E2E testing framework?

Options:
a) Playwright - Microsoft, multi-browser, fast
b) Cypress - Developer-friendly, single browser focus
c) None - No E2E tests (not recommended)

Context: Playwright is more powerful, Cypress has better DX
```

### UI-VIS-1: Visual Regression
```
Use visual regression testing?

Options:
a) Chromatic - Storybook integration, hosted
b) Percy - Standalone, BrowserStack owned
c) None - No visual regression tests

Context: Chromatic if using Storybook, Percy otherwise
```

### UI-STORY-1: Storybook Integration
```
Use Storybook for component documentation?

Options:
a) Yes (CSF3) - Component Story Format 3
b) Yes (MDX) - Markdown with JSX
c) No - No Storybook

Context: Storybook enables isolated component testing
```

### UI-QA-1: Manual QA Depth
```
How much manual QA coverage?

Options:
a) Full regression - Test all features every release
b) Critical paths only - Test happy paths + critical flows
c) None - Rely on automated tests only

Context: Critical paths is practical for most teams
```

### UI-QA-2: Device Coverage
```
Which devices to test?

Options:
a) All major - Desktop (Chrome, Firefox, Safari) + Mobile (iOS, Android)
b) Mobile-first - Prioritize mobile browsers
c) Desktop only - Desktop browsers only

Context: All major is ideal, mobile-first for mobile apps
```

## Output Templates

### e2e-tests.md

```markdown
---
status: draft
derived-from:
  - component-specs.md
  - state-machines.md
  - ui-acceptance-criteria.md
derived-at: {timestamp}
loom-version: "1.0"
structured-interview:
  decisions:
    UI-E2E-1: {answer}
    UI-VIS-1: {answer}
    UI-STORY-1: {answer}
    UI-QA-1: {answer}
    UI-QA-2: {answer}
---

# E2E Test Specifications

Playwright/Cypress test specifications derived from UI acceptance criteria.

---

## Test File: drag-drop.spec.ts

> **Implements:** [SM-DRAG](state-machines.md#sm-drag)

### E2E-UI-DRAG-001: New Task Placement {#e2e-ui-drag-001}

> **AC:** [AC-UI-DRAG-001-1](ui-acceptance-criteria.md#ac-ui-drag-001-1)

```typescript
test('should create assignment when task dropped on station', async ({ page }) => {
  // Arrange
  await page.goto('/scheduler');
  const sidebarTile = page.getByTestId('task-tile-task-001');
  const stationColumn = page.getByTestId('station-column-station-a');

  // Act
  await sidebarTile.dragTo(stationColumn, {
    targetPosition: { x: 50, y: 200 } // 7:00 AM position
  });

  // Assert
  await expect(page.getByTestId('tile-assignment-new')).toBeVisible();
  await expect(page.getByTestId('tile-assignment-new'))
    .toHaveAttribute('data-scheduled-start', /07:00/);
});
```

**Test Data:**
- Task: `{ id: 'task-001', name: 'Print Brochures', stationId: 'station-a' }`
- Expected time: 07:00 (based on Y=200px)

---

### E2E-UI-DRAG-002: Invalid Station Rejection {#e2e-ui-drag-002}

> **AC:** [AC-UI-DRAG-001-3](ui-acceptance-criteria.md#ac-ui-drag-001-3)

```typescript
test('should reject drop on wrong station', async ({ page }) => {
  // Arrange
  await page.goto('/scheduler');
  const sidebarTile = page.getByTestId('task-tile-task-001'); // station-a task
  const wrongStation = page.getByTestId('station-column-station-b');

  // Act
  await sidebarTile.dragTo(wrongStation);

  // Assert - tile should return to sidebar
  await expect(sidebarTile).toBeVisible();
  await expect(page.getByTestId('tile-assignment-new')).not.toBeVisible();
});
```

---

### E2E-UI-DRAG-003: Time Snapping {#e2e-ui-drag-003}

> **AC:** [AC-UI-DRAG-001-4](ui-acceptance-criteria.md#ac-ui-drag-001-4)

```typescript
test.describe('time snapping', () => {
  const snapCases = [
    { dropY: 170, expectedTime: '07:00' }, // 7:14 snaps to 7:00
    { dropY: 180, expectedTime: '07:30' }, // 7:16 snaps to 7:30
    { dropY: 220, expectedTime: '07:30' }, // 7:29 snaps to 7:30
    { dropY: 240, expectedTime: '08:00' }, // 7:46 snaps to 8:00
  ];

  for (const { dropY, expectedTime } of snapCases) {
    test(`should snap Y=${dropY} to ${expectedTime}`, async ({ page }) => {
      await page.goto('/scheduler');
      const tile = page.getByTestId('task-tile-task-001');
      const station = page.getByTestId('station-column-station-a');

      await tile.dragTo(station, { targetPosition: { x: 50, y: dropY } });

      await expect(page.getByTestId('tile-assignment-new'))
        .toHaveAttribute('data-scheduled-start', new RegExp(expectedTime));
    });
  }
});
```

---

## Test File: tile-states.spec.ts

> **Implements:** [SM-TILE](state-machines.md#sm-tile)

### E2E-UI-TILE-001: Tile Selection {#e2e-ui-tile-001}

```typescript
test('should select job when tile clicked', async ({ page }) => {
  await page.goto('/scheduler');
  const tile = page.getByTestId('tile-assignment-001');

  await tile.click();

  await expect(tile).toHaveClass(/selected/);
  await expect(page.getByTestId('job-details-panel'))
    .toContainText('Job: 12345');
});
```

### E2E-UI-TILE-002: Tile Recall {#e2e-ui-tile-002}

```typescript
test('should recall tile on double-click', async ({ page }) => {
  await page.goto('/scheduler');
  const tile = page.getByTestId('tile-assignment-001');

  await tile.dblclick();

  await expect(tile).not.toBeVisible();
  await expect(page.getByTestId('task-tile-task-001'))
    .not.toHaveClass(/scheduled/);
});
```

---

## Test Coverage Matrix

| AC ID | Test ID | Status |
|-------|---------|--------|
| AC-UI-DRAG-001-1 | E2E-UI-DRAG-001 | ✅ Covered |
| AC-UI-DRAG-001-2 | E2E-UI-DRAG-001 | ✅ Covered |
| AC-UI-DRAG-001-3 | E2E-UI-DRAG-002 | ✅ Covered |
| AC-UI-DRAG-001-4 | E2E-UI-DRAG-003 | ✅ Covered |
| AC-UI-TILE-001-1 | E2E-UI-TILE-001 | ✅ Covered |
| AC-UI-TILE-002-1 | E2E-UI-TILE-002 | ✅ Covered |
```

### visual-tests.md

```markdown
---
status: draft
derived-from: component-specs.md
derived-at: {timestamp}
loom-version: "1.0"
---

# Visual Regression Test Specifications

Chromatic/Percy visual test specifications for component states.

---

## Tile Component Visual Tests

> **Component:** [COMP-TILE](component-specs.md#comp-tile)

### VIS-UI-TILE-001: Tile States {#vis-ui-tile-001}

**Capture all tile visual states:**

| State | Description | Storybook Story |
|-------|-------------|-----------------|
| Idle | Default appearance | `Tile/Idle` |
| Hovered | Swap buttons visible | `Tile/Hovered` |
| Selected | Glow effect | `Tile/Selected` |
| Dragging | Ghost placeholder | `Tile/Dragging` |
| Muted | Desaturated | `Tile/Muted` |

**Storybook Stories:**

```typescript
// Tile.stories.tsx
export const Idle: Story = {
  args: {
    assignment: mockAssignment,
    task: mockTask,
    job: mockJob,
  },
};

export const Selected: Story = {
  args: {
    ...Idle.args,
    isSelected: true,
  },
};

export const Muted: Story = {
  args: {
    ...Idle.args,
    activeJobId: 'other-job',
  },
};
```

---

## Drop Zone Visual Tests

> **Component:** [COMP-STATION-COLUMN](component-specs.md#comp-station-column)

### VIS-UI-DROP-001: Drop Zone States {#vis-ui-drop-001}

| State | Ring Color | Background |
|-------|------------|------------|
| Normal | none | slate-900 |
| Valid | green-500 | green-500/10 |
| Invalid | red-500 | red-500/10 |
| Warning | orange-500 | orange-500/10 |

---

## Responsive Visual Tests

### VIS-UI-RESP-001: Breakpoint Snapshots {#vis-ui-resp-001}

Capture at standard breakpoints:

| Breakpoint | Width | Story Suffix |
|------------|-------|--------------|
| Mobile | 375px | `-mobile` |
| Tablet | 768px | `-tablet` |
| Desktop | 1280px | `-desktop` |
| Wide | 1920px | `-wide` |
```

### manual-qa.md

```markdown
---
status: draft
derived-from:
  - component-specs.md
  - state-machines.md
derived-at: {timestamp}
loom-version: "1.0"
---

# Manual QA Checklist

Manual testing checklist for aspects that cannot be automated.

---

## QA Focus Areas

Manual QA focuses on what automation cannot verify:

1. **UX Feel** - Does interaction feel smooth and natural?
2. **Real Devices** - Behavior on actual phones/tablets
3. **Accessibility** - Screen reader experience
4. **Visual Polish** - Subtle alignment, spacing issues
5. **Error Recovery** - Can user recover from errors?

---

## Drag-Drop Manual Tests

### QA-UI-DRAG-001: Drag Feel {#qa-ui-drag-001}

**Objective:** Verify drag interaction feels smooth

**Steps:**
1. [ ] Drag a tile from sidebar to grid
2. [ ] Observe cursor change (grab → grabbing)
3. [ ] Verify drag preview follows cursor smoothly
4. [ ] Check for any jank or stuttering

**Expected:**
- Immediate response to drag start (<100ms)
- Smooth 60fps cursor tracking
- No visual glitches during drag

**Devices:**
- [ ] Desktop Chrome
- [ ] Desktop Safari
- [ ] Desktop Firefox
- [ ] MacBook trackpad
- [ ] External mouse

---

### QA-UI-DRAG-002: Touch Drag {#qa-ui-drag-002}

**Objective:** Verify drag works on touch devices

**Steps:**
1. [ ] Long-press tile to initiate drag
2. [ ] Move finger to target location
3. [ ] Release to drop

**Expected:**
- Long-press (500ms) initiates drag
- Drag preview appears under finger
- Drop on release

**Devices:**
- [ ] iPhone Safari
- [ ] iPad Safari
- [ ] Android Chrome
- [ ] Android tablet

---

## Accessibility Manual Tests

### QA-UI-A11Y-001: Keyboard Navigation {#qa-ui-a11y-001}

**Objective:** Verify full keyboard operability

**Steps:**
1. [ ] Tab to first tile
2. [ ] Press Space/Enter to pick up
3. [ ] Use Arrow keys to move
4. [ ] Press Space/Enter to drop
5. [ ] Press Escape to cancel

**Expected:**
- Clear focus indicators
- Logical tab order
- All actions possible without mouse

---

### QA-UI-A11Y-002: Screen Reader {#qa-ui-a11y-002}

**Objective:** Verify screen reader experience

**Steps:**
1. [ ] Navigate to scheduler with VoiceOver/NVDA
2. [ ] Hear tile information read aloud
3. [ ] Perform drag operation
4. [ ] Hear feedback on drop

**Expected Announcements:**
- "Task: Print Brochures, Station A, Unscheduled"
- "Dragging Print Brochures"
- "Over Station A, 7:00 AM, Valid drop zone"
- "Dropped at 7:00 AM"

---

## Device Matrix

| Device | Browser | Priority | Tester |
|--------|---------|----------|--------|
| MacBook Pro | Chrome | High | - |
| MacBook Pro | Safari | High | - |
| MacBook Pro | Firefox | Medium | - |
| iPhone 14 | Safari | High | - |
| iPhone SE | Safari | Medium | - |
| iPad Pro | Safari | Medium | - |
| Samsung S23 | Chrome | High | - |
| Pixel 7 | Chrome | Medium | - |

---

## Regression Checklist

Run before each release:

### Critical Path

- [ ] Create new booking (drag from sidebar)
- [ ] Reschedule booking (drag within grid)
- [ ] Cancel booking (double-click)
- [ ] Job selection (click tile)
- [ ] Navigation (date picker, scroll)

### Edge Cases

- [ ] Drag to wrong station (should reject)
- [ ] Drag to past time (should reject)
- [ ] Rapid clicks (no double-processing)
- [ ] Network error during save
- [ ] Session timeout handling
```

### accessibility-audit.md

```markdown
---
status: draft
derived-from: component-specs.md
derived-at: {timestamp}
loom-version: "1.0"
---

# Accessibility Audit Specifications

Automated and manual accessibility testing specifications.

---

## Automated Checks (axe-core)

### A11Y-AUTO-001: Component Audit {#a11y-auto-001}

Run axe-core on all component stories:

```typescript
// In Storybook test-runner
import { checkA11y } from '@storybook/addon-a11y';

test('Tile should have no violations', async () => {
  await checkA11y(page);
});
```

### A11Y-AUTO-002: Page Audit {#a11y-auto-002}

Run Lighthouse accessibility audit:

```bash
lighthouse http://localhost:3000/scheduler \
  --only-categories=accessibility \
  --output=json
```

**Target Score:** ≥90

---

## WCAG 2.1 AA Checklist

### Perceivable

- [ ] **1.1.1** Non-text content has text alternatives
- [ ] **1.3.1** Info and relationships programmatically determinable
- [ ] **1.4.1** Color not sole means of conveying information
- [ ] **1.4.3** Contrast ratio ≥4.5:1 for text
- [ ] **1.4.11** Contrast ratio ≥3:1 for UI components

### Operable

- [ ] **2.1.1** All functionality keyboard accessible
- [ ] **2.1.2** No keyboard traps
- [ ] **2.4.3** Logical focus order
- [ ] **2.4.7** Visible focus indicator

### Understandable

- [ ] **3.2.1** No unexpected context changes on focus
- [ ] **3.3.1** Error identification
- [ ] **3.3.2** Labels for inputs

### Robust

- [ ] **4.1.2** Name, role, value for UI components
```

## Execution Steps

1. **Read input files**
   - Parse component-specs.md for components
   - Parse state-machines.md for state transitions
   - Parse ui-acceptance-criteria.md for test scenarios

2. **Conduct Structured Interview**
   - Ask UI-E2E-1 through UI-QA-2 questions
   - Wait for user answers
   - Record decisions in YAML frontmatter

3. **Generate E2E Tests**
   - One test per AC
   - Follow Playwright/Cypress best practices
   - Include test data

4. **Generate Visual Tests**
   - One visual test per component state
   - Include responsive breakpoints
   - Storybook story references

5. **Generate Manual QA**
   - Focus on what automation cannot verify
   - Include device matrix
   - Include regression checklist

6. **Generate Accessibility Audit**
   - Automated checks (axe-core)
   - WCAG 2.1 AA checklist

7. **Write output files**
   - e2e-tests.md
   - visual-tests.md
   - manual-qa.md
   - accessibility-audit.md

8. **Show summary**
   - Count of E2E tests
   - Count of visual tests
   - Coverage percentage
