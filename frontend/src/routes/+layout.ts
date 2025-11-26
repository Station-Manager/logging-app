import type { LayoutData } from './$types';
import { loadConfig } from '$lib/stores/config-store.js';
import { loadCatStateValues } from '$lib/stores/cat-state-store';
import { handleAsyncError } from '$lib/utils/error-handler';

export const prerender = true;
export const ssr = false;

/**
 * Asynchronously loads and initializes application layout data.
 *
 * This function attempts to load necessary configuration data for the
 * application. If an error occurs during the configuration loading process,
 * it logs the error message for debugging purposes. In all cases, it
 * resolves to an object, which represents the layout data.
 *
 * @returns {Promise<object>} A promise that resolves to an object containing
 * the application layout data.
 */
export const load: LayoutData = async (): Promise<object> => {
    try {
        await loadConfig();
        await loadCatStateValues();
    } catch (e: unknown) {
        handleAsyncError(e, '+layout.ts->load: LayoutData');
    }
    return {};
};
