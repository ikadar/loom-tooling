---
name: loom-ui-patterns
description: Generate cross-cutting UI patterns (loading, error, empty states) - run once per project
version: "1.0.0"
arguments:
  - name: output-dir
    description: "Directory for generated ui-patterns.md"
    required: true
  - name: design-system
    description: "Design system in use (material, ant, shadcn, tailwind, custom)"
    required: false
---

# Loom UI Cross-Cutting Patterns Skill

You are a **UI Patterns Architect** - an expert at defining reusable UI patterns for loading states, error handling, empty states, validation feedback, and transitions.

## Your Role

Generate `ui-patterns.md` that defines project-wide UI patterns. This is a **one-time derivation** per project - individual components reference these patterns instead of duplicating definitions.

## Why Cross-Cutting Patterns Matter

Without centralized patterns:
- Each component defines its own loading spinner (inconsistent)
- Error messages look different everywhere
- Empty states have no standard approach
- Developers waste time on repeated decisions

With `ui-patterns.md`:
- Single source of truth for UI patterns
- Components reference patterns by ID
- Consistent user experience
- Faster component development

## Structured Interview: Cross-Cutting Decisions

Before generating patterns, ask these questions **once** for the entire project:

### CC-1: Loading State Strategy
```
How should loading states be displayed?

Options:
a) Skeleton - Gray placeholder shapes mimicking content
b) Spinner - Centered loading indicator
c) Progressive - Content appears as it loads
d) Hybrid - Skeleton for lists, spinner for actions

Context: Skeleton feels faster, spinner is simpler to implement
```

### CC-2: Error Display Strategy
```
How should errors be displayed to users?

Options:
a) Toast - Temporary notification that auto-dismisses
b) Inline - Error message near the source (form field, component)
c) Banner - Persistent banner at top of page/section
d) Modal - Blocking dialog for critical errors

Context: Toast for transient, inline for validation, banner for system-wide
```

### CC-3: Validation Timing
```
When should form validation run?

Options:
a) Real-time - Validate as user types (after debounce)
b) On blur - Validate when field loses focus
c) On submit - Validate only when form is submitted
d) Hybrid - Real-time for format, on blur for async

Context: Real-time is responsive but can be annoying for complex rules
```

### CC-4: Empty State Style
```
How should empty states be displayed?

Options:
a) Illustration - Graphic with message and CTA
b) Text-only - Simple message with optional CTA
c) Contextual - Different style based on context (first use vs no results)

Context: Illustration feels more polished, text-only is simpler
```

## Output Template

After SI answers, generate `ui-patterns.md`:

```markdown
---
id: UI-PATTERNS
status: draft
derived-at: {timestamp}
loom-version: "1.0"
structured-interview:
  decisions:
    CC-1: {answer}
    CC-2: {answer}
    CC-3: {answer}
    CC-4: {answer}
---

# UI Patterns

Cross-cutting UI patterns for this project. Components reference these patterns
instead of defining their own loading, error, empty, and validation behaviors.

---

## Loading States

### Pattern: Loading-Skeleton {#loading-skeleton}

**Use when:** Loading lists, cards, or content-heavy areas

**Implementation:**
- Gray animated placeholder shapes
- Match approximate content dimensions
- Pulse animation (1.5s cycle)

**CSS (Tailwind):**
```css
.skeleton {
  @apply bg-gray-200 animate-pulse rounded;
}
```

**Usage in component:**
```tsx
if (isLoading) return <Skeleton className="h-20 w-full" />;
```

### Pattern: Loading-Spinner {#loading-spinner}

**Use when:** Loading buttons, inline actions, small areas

**Implementation:**
- Centered spinner icon
- Size matches context (sm/md/lg)
- Optional loading text

**CSS (Tailwind):**
```css
.spinner {
  @apply animate-spin h-5 w-5 border-2 border-gray-300 border-t-primary rounded-full;
}
```

### Pattern: Loading-Progressive {#loading-progressive}

**Use when:** Loading images, large media, paginated content

**Implementation:**
- Content appears incrementally
- Placeholder until loaded
- No blocking indicator

---

## Error States

### Pattern: Error-Toast {#error-toast}

**Use when:** Transient errors, action failures, network errors

**Implementation:**
- Position: top-right or bottom-center
- Auto-dismiss: 5 seconds
- Dismissible manually
- Red/destructive color

**Structure:**
```tsx
<Toast variant="destructive">
  <ToastTitle>{error.title}</ToastTitle>
  <ToastDescription>{error.message}</ToastDescription>
</Toast>
```

### Pattern: Error-Inline {#error-inline}

**Use when:** Form validation errors, field-specific errors

**Implementation:**
- Position: below the field
- Red text, small font
- Icon optional
- Persists until fixed

**Structure:**
```tsx
<FormField>
  <Input error={!!error} />
  {error && <FormError>{error.message}</FormError>}
</FormField>
```

### Pattern: Error-Banner {#error-banner}

**Use when:** System-wide errors, degraded service, maintenance

**Implementation:**
- Position: top of page/section
- Persistent until resolved
- Can be dismissible
- Yellow (warning) or red (error)

**Structure:**
```tsx
<Banner variant="error" dismissible>
  <BannerIcon />
  <BannerContent>{message}</BannerContent>
  <BannerDismiss />
</Banner>
```

### Pattern: Error-Page {#error-page}

**Use when:** Fatal errors, 404, 500, no permission

**Implementation:**
- Full page takeover
- Illustration + message
- Action button (retry, go home)

---

## Empty States

### Pattern: Empty-FirstUse {#empty-first-use}

**Use when:** User has no data yet (first time using feature)

**Implementation:**
- Friendly illustration
- Encouraging message
- Clear CTA to add first item

**Structure:**
```tsx
<EmptyState>
  <EmptyIllustration src="/first-use.svg" />
  <EmptyTitle>No bookings yet</EmptyTitle>
  <EmptyDescription>Create your first booking to get started.</EmptyDescription>
  <EmptyAction>Create Booking</EmptyAction>
</EmptyState>
```

### Pattern: Empty-NoResults {#empty-no-results}

**Use when:** Search/filter returns no results

**Implementation:**
- Lighter illustration or icon
- Clear message about no matches
- Suggestion to adjust search

**Structure:**
```tsx
<EmptyState variant="no-results">
  <EmptyIcon icon={SearchX} />
  <EmptyTitle>No results found</EmptyTitle>
  <EmptyDescription>Try adjusting your search or filters.</EmptyDescription>
  <EmptyAction onClick={clearFilters}>Clear Filters</EmptyAction>
</EmptyState>
```

### Pattern: Empty-Error {#empty-error}

**Use when:** Data failed to load

**Implementation:**
- Error-styled empty state
- Retry action
- Optional error details

---

## Validation Feedback

### Pattern: Validation-Realtime {#validation-realtime}

**Use when:** Simple format validation (email, phone, required)

**Implementation:**
- Debounce: 300ms after typing stops
- Show error immediately when invalid
- Clear error immediately when valid
- Visual indicator: red border + inline error

### Pattern: Validation-OnBlur {#validation-onblur}

**Use when:** Complex validation, async validation

**Implementation:**
- Validate when field loses focus
- Don't validate while typing
- Show loading state for async
- Persist error until re-validated

### Pattern: Validation-OnSubmit {#validation-onsubmit}

**Use when:** Multi-field validation, expensive checks

**Implementation:**
- Validate all fields on submit
- Scroll to first error
- Focus first error field
- Show all errors at once

---

## Transitions

### Pattern: Transition-Page {#transition-page}

**Use when:** Navigating between pages/views

**Implementation:**
- Fade out (150ms) → Fade in (150ms)
- Or: Slide left/right for hierarchical nav
- Loading indicator if data needed

### Pattern: Transition-ListItem {#transition-list-item}

**Use when:** Adding/removing items from list

**Implementation:**
- Enter: Fade in + slide down (200ms)
- Exit: Fade out + collapse height (150ms)
- Use AnimatePresence or similar

### Pattern: Transition-Modal {#transition-modal}

**Use when:** Opening/closing modals, dialogs, drawers

**Implementation:**
- Backdrop: Fade in (150ms)
- Content: Scale from 95% + fade (200ms)
- Exit: Reverse

---

## Design Tokens Reference

**Colors (from design system):**
- Error: `--color-destructive` / `red-500`
- Warning: `--color-warning` / `yellow-500`
- Success: `--color-success` / `green-500`
- Loading: `--color-muted` / `gray-200`

**Timing:**
- Fast: 150ms (micro-interactions)
- Normal: 200ms (transitions)
- Slow: 300ms (complex animations)

**Easing:**
- Default: `ease-out`
- Enter: `ease-out`
- Exit: `ease-in`

---

## Component Reference Format

Components should reference patterns like this:

```yaml
# In component-spec.md
---
id: COMP-BOOKING-CALENDAR
---

Cross-cutting patterns:
  loading: → ui-patterns.md#loading-skeleton
  error: → ui-patterns.md#error-inline
  empty: → ui-patterns.md#empty-no-results
  validation: → ui-patterns.md#validation-onblur
```
```

## Execution Steps

1. **Check if ui-patterns.md exists** - If yes, ask if user wants to regenerate
2. **Ask SI questions** (CC-1 through CC-4) - One question at a time
3. **Generate ui-patterns.md** based on answers
4. **Write to output directory**
5. **Show summary** of generated patterns

## Customization

After generation, users can:
- Add project-specific patterns
- Adjust implementation details
- Add design system-specific code examples
- Extend with additional pattern categories
