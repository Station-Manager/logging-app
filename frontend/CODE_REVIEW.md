# Frontend Code Review - Logging App (Svelte 5 / TypeScript)

**Review Date:** January 8, 2026

## Executive Summary

This is a well-structured Svelte 5 application with a Wails backend for amateur radio QSO logging. The codebase demonstrates good practices in many areas, with some opportunities for improvement. Overall code quality is **good to very good**.

---

## 1. Architecture & Project Structure

### ✅ Strengths

- **Clean separation of concerns**: The `$lib` structure is logical with `states/`, `stores/`, `utils/`, `ui/`, and `constants/` directories
- **Svelte 5 runes adoption**: Good use of `$state`, `$derived`, and `$props` throughout
- **Static adapter**: Appropriate for Wails desktop application (SSR disabled, prerender enabled)
- **TypeScript strict mode**: Enabled with proper configuration
- **Tailwind CSS v4**: Modern CSS approach with utility classes
- **Intentional state management pattern**: Svelte 5 runes (`$state`) are used for frequently-changing reactive data, while Svelte stores (`writable`/`derived`) are used for infrequently-changing or static configuration data. This hybrid approach provides positive performance benefits by avoiding unnecessary reactive overhead for stable values.
    - `logging-mode-store.ts` - mode changes rarely during a session
    - `cat-state-store.ts` - CAT state values are loaded once and cached

---

## 2. State Management

### ✅ Strengths

- **Well-defined interfaces**: `QsoState`, `CatState`, `SessionState`, etc. are properly typed
- **Methods on state objects**: Good pattern of co-locating mutations with state (e.g., `qsoState.reset()`)
- **Cleanup patterns**: `cleanupSessionState()` properly clears intervals to prevent memory leaks
- **Timer management**: Module-scoped interval IDs prevent duplicate timers
- **Performance-conscious design**: Using stores for infrequently-changing data (logging mode, CAT state values) and runes for highly-reactive data (QSO fields, timers) is an appropriate optimization

### ⚠️ Issues Identified

1. ~~**ClipboardState class vs object pattern inconsistency**~~ ✅ FIXED
   Converted `ClipboardState` from class pattern to object pattern with `$state` interface to match other state files:

    ```typescript
    export interface ClipboardState {
        list: string[];
        maxLength: number;
        add(this: ClipboardState, item: string): void;
        setMaxLength(this: ClipboardState, len: number): void;
    }

    export const clipboardState: ClipboardState = $state({ ... });
    ```

2. ~~**Potential memory leak in ContestTimersClass**~~ ✅ FIXED
   Changed timer ID types from `number` to `number | null` for cleaner semantics:

    ```typescript
    sinceStartTimerID: number | null = null;
    sinceLastQsoTimerID: number | null = null;
    ```

    Now properly checks for `null` before clearing intervals and sets to `null` after clearing.

3. **QsoState.stopTimer() clears time fields**:
    ```typescript
    stopTimer(this: QsoState): void {
        // ...
        qsoState.time_on = getTimeUTC();  // This resets time_on!
        qsoState.time_off = qsoState.time_on;
    }
    ```
    This may be intentional but could be surprising if the user just wants to pause timing without resetting.

---

## 3. Component Design

### ✅ Strengths

- **Reusable input components**: `TextInput`, `DateInput`, `TimeInput`, `Rst`, `Mode` are well-abstracted
- **Customizable via props**: Components accept CSS customization props (`labelCss`, `divCss`, `inputCss`)
- **Snippets for conditional rendering**: Good use of Svelte 5 snippets in `QsoPanel.svelte` for `normalLogging`/`contextLogging`

### ⚠️ Issues Identified

1. **Direct DOM manipulation in several places**:

    ```svelte
    // LoggingCardHeader.svelte
    const operatorField = (): void => {
        appState.activePanel = STATION_PANEL
        const opField = document.getElementById("operator_callsign") as HTMLInputElement;
        if (opField) {
            opField.focus();
            opField.select();
        }
    }
    ```

    Consider using `bind:this` and passing refs through context or props instead of `getElementById`.

2. **Mixed shortcut handling approaches**:
    - Some components use `@svelte-put/shortcut` action on `<svelte:window>`
    - `TimerControls.svelte` defines `toggleTimer` for F5 but buttons have separate click handlers
    - Consider centralizing keyboard shortcuts in a single location

3. **Hardcoded dimensions throughout**:

    ```svelte
    <div class="w-[1024px] h-[651px]">
    <div class="flex flex-col h-[370px]">
    ```

    These fixed dimensions may cause issues if the app window is resized. Consider using more flexible layouts or CSS variables.

4. **Duplicated CSS classes**: Many components have nearly identical CSS class strings. Consider extracting common patterns:
    ```typescript
    // Example: Create a ui/styles.ts
    export const inputBase =
        'block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600';
    ```

---

## 4. Type Safety

### ✅ Strengths

- **Strict TypeScript**: `strict: true` in tsconfig
- **Well-typed Wails bindings**: Generated `models.ts` provides type safety for Go structs
- **Explicit interface definitions**: All state objects have corresponding interfaces

### ⚠️ Issues Identified

1. **`any` usage in tests**:

    ```typescript
    sessionState.list = [{} as never, {} as never];
    ```

    Using `never` cast is acceptable for tests, but could use factory functions for cleaner test data.

2. **Type assertion overuse**:

    ```typescript
    const target = event.currentTarget as HTMLInputElement;
    ```

    This is often necessary for event handlers but could benefit from type guards in some cases.

3. **Unsafe tag mapping in cat-state.svelte.ts**:
    ```typescript
    const tag = key as unknown as tags.CatStateTag;
    ```
    This double cast (`as unknown as`) bypasses type checking. Consider validating the tag first.

---

## 5. Error Handling

### ✅ Strengths

- **Centralized error handler**: `handleAsyncError()` provides consistent error logging and toast notifications
- **Wails LogError integration**: Errors are logged to the backend
- **Toast notifications**: User-friendly error display with appropriate duration and styling

### ⚠️ Issues Identified

1. **Swallowed errors in callsignLookup**:

    ```typescript
    const callsignLookup = async (): Promise<void> => {
        // ...
        try {
            // await OpenInBrowser($configState.qrz_view_url + qsoState.call.toUpperCase());
        } catch (e: unknown) {
            handleAsyncError(e, 'Failed to open callsign lookup');
        }
    };
    ```

    The actual call is commented out but the catch block exists - dead code.

2. **Missing error handling in some async operations**:
    - `+layout.svelte onMount()` has a try-catch but only wraps `Ready()`, not `sessionState.start()`

3. **Silent failures**: `qsoState.fromQso()` and `qsoState.updateFromCAT()` return early on null/undefined without logging.

---

## 6. Performance Considerations

### ✅ Strengths

- **Derived values**: Good use of `$derived` and `$derived.by()` for computed properties
- **Efficient list filtering**: `ClipboardState.add()` creates new arrays for reactivity instead of mutating
- **Hybrid state management**: Using Svelte stores for infrequently-changing data and runes for reactive data reduces unnecessary reactive overhead

### ⚠️ Potential Issues

1. **Frequent timer callbacks**:
    - `QSO_TIMER_INTERVAL_MS = 60_000` (1 minute) - reasonable
    - `CONTEST_TIMER_INTERVAL_MS = 1_000` (1 second) - could cause frequent reactivity updates

    Consider batching updates or using `requestAnimationFrame` for display-only updates.

2. ~~**Re-computation in templates**~~ ✅ FIXED

    ```svelte
    <div>{formatTimeSecondsToHHColonMMColonSS(contestTimers.elapsedSinceLastQso)}</div>
    ```

    Added `$derived` properties to `ContestTimersClass` for formatted values:

    ```typescript
    formattedSinceStart: string = $derived(
        formatTimeSecondsToHHColonMMColonSS(this.elapsedSinceStart)
    );
    formattedSinceLastQso: string = $derived(
        formatTimeSecondsToHHColonMMColonSS(this.elapsedSinceLastQso)
    );
    ```

    Components now use `{contestTimers.formattedSinceLastQso}` instead of calling the function directly.

3. **Large state objects**: `QsoState` has ~50+ properties. Consider grouping related fields into sub-objects.

---

## 7. Testing

### ✅ Strengths

- **Good test coverage for state logic**: `session-state.svelte.test.ts`, `contest-state.svelte.test.ts`, `contest-timers.svelte.test.ts`, `clipboard-state.svelte.test.ts`, `focus-context.svelte.test.ts`
- **Vitest with fake timers**: Proper testing of interval-based logic
- **Isolated tests**: Each test properly resets state in `beforeEach`
- **Mocking for context-based modules**: Focus context tests mock Svelte's context API
- **Co-located tests**: Tests are placed next to the files they test (e.g., `frequency.ts` → `frequency.test.ts`)
- **Component testing setup**: `@testing-library/svelte` configured with proper Svelte 5 support

### ⚠️ Areas for Improvement

1. **Limited component tests**: Only `TextInput.test.ts` exists. Consider adding tests for:
    - Form validation behavior in `Callsign.svelte`
    - User interactions in `FormControls.svelte`

2. **Missing utility tests**: `time-date.ts` has no tests (only `frequency.ts` is tested)

~~3. **Test file location inconsistency**~~ ✅ RESOLVED
Removed duplicate `tests/frequency.test.ts` and empty `tests/` directory. All tests now follow the co-located pattern:

```
src/lib/utils/frequency.ts      → src/lib/utils/frequency.test.ts
src/lib/states/session-state.svelte.ts → src/lib/states/session-state.svelte.test.ts
```

---

## 8. Accessibility

### ✅ Strengths

- **Labels on form inputs**: Most inputs have associated `<label>` elements with `for` attributes
- **aria-label on icon buttons**: Buttons with only icons have `aria-label`
- **disabled states**: Form elements properly use `disabled` attribute

### ⚠️ Issues Identified

1. **Missing aria-invalid propagation**:

    ```svelte
    aria-invalid={invalid}
    ```

    Only `Callsign.svelte` has this; other inputs with validation don't.

2. ~~**No focus management**~~ ✅ FIXED
   Implemented cross-component focus management via Svelte Context:
    - Created `focus-context.svelte.ts` with `FocusContext` interface and `FocusRefsClass`
    - Components now use `bind:this={focusContext.refs.elementName}` to register focusable elements
    - Focus is triggered via `focusContext.focus('elementName', select?)` which uses `tick()` for proper timing
    - Replaced all `document.getElementById()` calls in:
        - `LoggingCardHeader.svelte` (operator callsign)
        - `FormControls.svelte` (callsign input)
        - `SessionPanel.svelte` (edit callsign)
        - `QsoPanel.svelte` (srx_rcvd)
        - `InfoPanelHeader.svelte` (email input)

3. **Color contrast concerns**: Red validation states use `outline-red-600` which should be verified for WCAG compliance.

4. **Missing skip links and landmarks**: No `<main>` role or skip navigation for keyboard users.

---

## 9. Security

### ✅ No Critical Issues

- **No XSS vulnerabilities**: Svelte's default escaping handles user input
- **Backend validation**: Form data is validated before sending to Go backend
- **Email validation**: Uses regex pattern for email validation in `InfoPanelHeader.svelte`

### ⚠️ Minor Concerns

1. **Client-side only validation**: Callsign and RST validation is client-side; ensure backend also validates
2. **Email regex complexity**: The email regex is complex; consider using a library or simpler pattern

---

## 10. Code Quality & Maintainability

### ✅ Strengths

- **Consistent formatting**: Prettier configured and enforced
- **ESLint configuration**: Proper TypeScript and Svelte linting
- **Good documentation**: JSDoc comments on important functions (e.g., `FormControls.svelte`)
- **Named constants**: Magic numbers extracted to `constants/timers.ts`

### ⚠️ Issues Identified

1. **Commented-out code**:

    ```typescript
    // QsoPanel.svelte
    // if (target.value.length < 1) {
    //     target.classList.add('outline-red-500', 'outline-2');
    // }
    ```

    Remove or implement properly.

2. **Inconsistent import ordering**: Some files have Wails imports first, others have Svelte imports first.

3. **Large files**: `new-qso-state.svelte.ts` is 365 lines; consider splitting helper functions.

4. **Unused import potential**: `isValid` in `StationPanel.svelte` is always `true` - appears unused.

---

## 11. Bug Reports

### 🐛 Confirmed Bugs (FIXED)

1. ~~**SessionPanel shows wrong RST value**~~ ✅ FIXED

    ```svelte
    <div class={rstCss}>{entry.rst_sent}</div>
    <div class={rstCss}>{entry.rst_rcvd}</div>  // Fixed: was rst_sent
    ```

    Line ~103 in SessionPanel.svelte now correctly displays `rst_rcvd`.

2. ~~**QsoEditState missing tx_pwr assignment**~~ ✅ FIXED

    ```typescript
    // Added tx_pwr to both fromQso() and toQso() methods
    this.tx_pwr = qso.tx_pwr;  // in fromQso
    tx_pwr: this.tx_pwr,       // in toQso
    ```

3. ~~**VFO frequency inputs not updating qsoEditState**~~ ✅ FIXED
   In `SessionPanel.svelte`'s `vfos` snippet, changed to two-way binding:
    ```svelte
    bind:value={qsoEditState.freq_rx}
    bind:value={qsoEditState.freq}
    ```

---

## 12. Recommendations Summary

### High Priority

~~1. Fix the `rst_rcvd` bug in SessionPanel.svelte~~ ✅ FIXED
~~2. Fix the missing `tx_pwr` in qsoEditState.toQso()~~ ✅ FIXED
~~3. Fix VFO frequency inputs in edit modal to use two-way binding~~ ✅ FIXED

### Medium Priority

4. Add component tests with `@testing-library/svelte`
5. Extract common CSS patterns to reduce duplication
   ~~6. Replace `document.getElementById()` with `bind:this` refs~~ ✅ FIXED - Implemented focus context
6. Add tests for `time-date.ts` utilities
7. Centralize keyboard shortcut definitions

### Low Priority

9. Consider responsive/flexible layouts instead of fixed dimensions
10. Add aria-invalid to all validating inputs
11. Remove commented-out code blocks
12. Split large state files into smaller modules

---

## Conclusion

The codebase is well-organized and demonstrates good understanding of Svelte 5 patterns. The intentional use of Svelte stores for infrequently-changing data alongside runes for reactive data is a smart performance optimization. The main concerns are:

- A few actual bugs that need fixing
- Limited test coverage for components
- Some accessibility gaps

Overall, this is a solid foundation that would benefit from the recommended improvements.
