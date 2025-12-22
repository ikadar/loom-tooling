# UI/UX Completeness Checklist

Reference document for UI/UX analysis. Used by `/loom-analyze-ui` command.

---

## 0. UI Input Detection

Before analyzing, check if UI/UX input exists:

| Aspect | Present? |
|--------|----------|
| Screen/page list | ? |
| Layout diagrams/wireframes | ? |
| Component specifications | ? |
| Interaction patterns | ? |
| State definitions | ? |
| Error states | ? |
| Loading states | ? |
| Empty states | ? |
| Responsive behavior | ? |
| Accessibility requirements | ? |

**If more than 3 aspects missing:** STOP and request UX-UI input.

---

## 1. Screen/Page Checklist

For EACH screen/page:

### A. Basic Definition

| Aspect | Question |
|--------|----------|
| Name | Screen identifier |
| Purpose | Primary goal of this screen |
| Entry Points | How does user get here? (nav, link, redirect, deep link) |
| Exit Points | Where can user go from here? |
| URL/Route | What URL pattern? Parameters? |
| Title | Page title for browser tab/SEO |
| Permissions | Who can access? What if unauthorized? |

### B. Layout

| Aspect | Question |
|--------|----------|
| Layout Type | Fixed width, fluid, or responsive? |
| Breakpoints | Mobile, tablet, desktop sizes? |
| Grid System | Column count, gutter size? |
| Component Placement | Header, sidebar, main, footer positions? |
| Scroll Behavior | Standard scroll, infinite scroll, pagination? |
| Sticky Elements | What stays visible on scroll? |
| Z-index Layers | Layering order for overlays, modals? |

### C. Initial State

| Aspect | Question |
|--------|----------|
| Default Values | What are default filter/sort/view settings? |
| Data Loading | What loads immediately vs lazy? |
| Focus | Where is initial focus? |
| Scroll Position | Top, remembered, or specific element? |

---

## 2. Component Checklist

For EACH component:

### A. Basic Definition

| Aspect | Question |
|--------|----------|
| Name | Component identifier |
| Purpose | What does it do? |
| Location | Which screens? |
| Data Source | What data does it display/edit? |
| Reusability | Single use or reusable across screens? |

### B. Visual States

| State | Questions |
|-------|-----------|
| Default/Idle | Base appearance |
| Hover | Mouse over appearance (desktop) |
| Focus | Keyboard focus appearance |
| Active/Pressed | Click/tap appearance |
| Disabled | Inactive appearance, cursor |
| Loading | Spinner, skeleton, shimmer? |
| Error | Error indicator, color, icon |
| Success | Success indicator, duration |
| Empty | No data appearance |
| Selected | Selection indicator |
| Dragging | Appearance while being dragged |
| Drop Target | Valid/invalid drop zone appearance |
| Partially Selected | For multi-select scenarios |
| Read-only | Editable vs view-only distinction |

### C. Content States

| State | Questions |
|-------|-----------|
| No Data Yet | First-time empty state (illustration, CTA) |
| No Data (searched) | No results message, suggestions |
| No Data (filtered) | Clear filters CTA |
| Loading Initial | Full component loading |
| Loading More | Incremental loading (pagination) |
| Loading Refresh | Background refresh indicator |
| Partial Data | Show available, indicate missing |
| Error Loading | Error message, retry button |
| Stale Data | Stale indicator, refresh prompt |
| Offline | Cached data indicator |

### D. Sizing

| Aspect | Question |
|--------|----------|
| Width | Fixed, fluid, or content-based? |
| Height | Fixed, fluid, or content-based? |
| Min/Max | Minimum and maximum dimensions? |
| Overflow | Scroll, truncate, or wrap? |

---

## 3. Interaction Checklist

For EACH interaction:

### A. Triggers

| Trigger Type | Questions |
|--------------|-----------|
| Mouse | Click, double-click, right-click, drag, hover? |
| Keyboard | Key or shortcut? Tab order? |
| Touch | Tap, long press, swipe, pinch, multi-touch? |
| Voice | Voice command? (if applicable) |
| Gesture | Custom gestures? |
| Programmatic | API or event trigger? |

### B. Feedback During

| Phase | Questions |
|-------|-----------|
| Start | Visual indication interaction started |
| In Progress | Progress indicator, animation |
| Valid Target | Appearance of valid drop/target |
| Invalid Target | Appearance of invalid drop/target |
| Near Boundary | Edge/boundary indication |

### C. Outcomes

| Outcome | Questions |
|---------|-----------|
| Success | Visual feedback, duration, next state |
| Partial Success | What's shown, what's missing |
| Failure | Error display, recovery options |
| Cancelled | Return to previous state, cleanup |

### D. Modifiers

| Modifier | Questions |
|----------|-----------|
| Shift+Action | Modified behavior? |
| Ctrl/Cmd+Action | Modified behavior? |
| Alt+Action | Modified behavior? |

### E. Constraints

| Constraint | Questions |
|------------|-----------|
| Valid Targets | Where can interaction complete? |
| Invalid Targets | What happens if invalid target? |
| Boundaries | Screen edges, container edges? |
| Timing | Debounce, throttle, timeout? |
| Conflicts | What if conflicting interaction? |

---

## 4. Navigation Checklist

### A. Global Navigation

| Element | Questions |
|---------|-----------|
| Primary Nav | Location, items, active state |
| Secondary Nav | Location, relationship to primary |
| Breadcrumbs | When shown, format, clickable? |
| Back Button | Browser vs in-app, behavior |
| Home Link | Logo? Text? Behavior |
| User Menu | Items, location, avatar |
| Search | Global or contextual, location |
| Notifications | Badge, dropdown, behavior |

### B. In-Page Navigation

| Element | Questions |
|---------|-----------|
| Tabs | Style, switching behavior, lazy load? |
| Accordion | Single or multi-open, animation |
| Anchor Links | Scroll behavior, offset for header |
| Pagination | Style, items per page, jump-to |
| Infinite Scroll | Load trigger, loading indicator |
| Table of Contents | Sticky, highlight current |

### C. Transitions

| Transition | Questions |
|------------|-----------|
| Page to Page | Animation type, duration, easing |
| Panel Open/Close | Animation, direction |
| Modal Open/Close | Backdrop, animation, focus trap |
| Drawer Open/Close | Direction, overlay |
| List Add/Remove | Animation for items |
| Collapse/Expand | Height animation |

### D. Deep Linking & History

| Aspect | Questions |
|--------|-----------|
| URL State | What app state is in URL? |
| Shareable URLs | Can current view be shared? |
| Bookmarks | Bookmark support |
| History | Browser back/forward behavior |
| Restore | State restoration on return |

---

## 5. Cross-Cutting Patterns

### A. Loading Patterns

| Scenario | Questions |
|----------|-----------|
| Initial Page Load | Skeleton, spinner, or blank? |
| Data Refresh | In-place, overlay, or toast? |
| Action In Progress | Button loading, full overlay? |
| Background Sync | Silent or indicator? |
| Lazy Load | Trigger point, placeholder |

### B. Error Patterns

| Scenario | Questions |
|----------|-----------|
| Validation Error | Inline, toast, banner? Timing? |
| API Error | Toast, modal, inline? Retry? |
| Network Error | Banner, offline mode, retry? |
| Permission Error | Redirect, message, login prompt? |
| Not Found | 404 page design |
| Server Error | 500 page design, support contact |

### C. Empty State Patterns

| Scenario | Questions |
|----------|-----------|
| First Use | Illustration, onboarding CTA |
| No Results | Message, suggestions, clear filters |
| Deleted | Confirmation, undo option |

### D. Notification Patterns

| Type | Questions |
|------|-----------|
| Success | Toast, inline, duration, dismissible? |
| Warning | Style, persistence, action? |
| Error | Style, persistence, retry action? |
| Info | Style, duration, dismissible? |
| Confirmation | Modal, inline, required action? |

---

## 6. Accessibility Checklist

| Aspect | Questions |
|--------|-----------|
| Keyboard Navigation | Full keyboard access? Tab order? |
| Focus Indicators | Visible focus state? |
| Screen Reader | ARIA labels, roles, live regions? |
| Color Contrast | WCAG AA (4.5:1) or AAA (7:1)? |
| Color Independence | Info not conveyed by color alone? |
| Touch Targets | Minimum 44x44px? |
| Text Sizing | Scalable? Minimum size? |
| Reduced Motion | Respects prefers-reduced-motion? |
| High Contrast | Supports high contrast mode? |
| Error Identification | Errors identified without color alone? |

---

## 7. Responsive Checklist

| Breakpoint | Questions |
|------------|-----------|
| Desktop (>1200px) | Full layout |
| Laptop (1024-1200px) | Layout adjustments |
| Tablet (768-1024px) | Navigation changes, touch optimization |
| Mobile (480-768px) | Single column, hamburger menu |
| Small Mobile (<480px) | Minimum viable layout |

For each breakpoint:
- Layout changes
- Hidden elements
- Changed interactions
- Touch vs mouse
- Typography scaling
