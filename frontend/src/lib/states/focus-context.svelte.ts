/**
 * Focus Context for Cross-Component Focus Management
 *
 * This module provides a Svelte context-based approach for managing focus
 * across components, replacing direct DOM manipulation via document.getElementById().
 *
 * Usage:
 * 1. In root layout: call setFocusContext() to initialize
 * 2. In components with focusable elements: bind:this={focusRefs.elementName}
 * 3. In components needing to focus: call focusRefs.focusElement('elementName')
 */

import { getContext, setContext, tick } from 'svelte';

const FOCUS_CONTEXT_KEY = Symbol('focus-context');

/**
 * Interface defining all focusable element refs managed by the context.
 * Add new focusable elements here as needed.
 */
export interface FocusRefs {
    // QSO Panel inputs
    callsignInput: HTMLInputElement | null;
    srxRcvdInput: HTMLInputElement | null;

    // Station Panel inputs
    operatorCallsignInput: HTMLInputElement | null;

    // Info Panel inputs
    fwdSessionEmailInput: HTMLInputElement | null;

    // Edit modal inputs (SessionPanel)
    editCallsignInput: HTMLInputElement | null;
}

/**
 * Interface for the focus context including refs and helper methods
 */
export interface FocusContext {
    refs: FocusRefs;

    /**
     * Focus an element by its ref name, optionally selecting its content.
     * Uses tick() to ensure DOM is updated before focusing.
     *
     * @param refName - The name of the ref in FocusRefs
     * @param select - Whether to also select the input content (default: false)
     */
    focus(refName: keyof FocusRefs, select?: boolean): Promise<void>;
}

/**
 * Class-based focus refs to enable $state usage
 */
class FocusRefsClass implements FocusRefs {
    callsignInput: HTMLInputElement | null = $state(null);
    srxRcvdInput: HTMLInputElement | null = $state(null);
    operatorCallsignInput: HTMLInputElement | null = $state(null);
    fwdSessionEmailInput: HTMLInputElement | null = $state(null);
    editCallsignInput: HTMLInputElement | null = $state(null);
}

/**
 * Creates and sets the focus context. Call this in the root layout.
 */
export function setFocusContext(): FocusContext {
    const refs = new FocusRefsClass();

    const context: FocusContext = {
        refs,
        async focus(refName: keyof FocusRefs, select: boolean = false): Promise<void> {
            // Wait for any pending DOM updates
            await tick();

            const element = refs[refName];
            if (element) {
                element.focus();
                if (select && 'select' in element) {
                    element.select();
                }
            }
        },
    };

    setContext(FOCUS_CONTEXT_KEY, context);
    return context;
}

/**
 * Gets the focus context. Call this in any component that needs to focus elements
 * or register focusable elements.
 */
export function getFocusContext(): FocusContext {
    const context = getContext<FocusContext>(FOCUS_CONTEXT_KEY);
    if (!context) {
        throw new Error(
            'FocusContext not found. Make sure setFocusContext() is called in a parent component.'
        );
    }
    return context;
}
