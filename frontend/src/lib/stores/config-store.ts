import { writable, type Writable } from 'svelte/store';
import { types } from '$lib/wailsjs/go/models';
import { FetchUiConfig } from '$lib/wailsjs/go/facade/Service';
import { LogError } from '$lib/wailsjs/runtime';

// The backend guarantees that UiConfig and its fields are never null. We
// therefore model configStore as a non-nullable UiConfig and always
// initialize it with a concrete instance.
export const configStore: Writable<types.UiConfig> = writable(new types.UiConfig());

/**
 * Asynchronously loads the UI configuration.
 *
 * This function fetches the UI configuration using an asynchronous request.
 * If the fetched configuration is null or undefined, a default configuration
 * is applied, and an error message is logged. In case of an exception during
 * the fetch operation, the error is logged.
 *
 * The configuration is stored using the `configStore` module.
 *
 * @returns A promise that resolves with no result (void).
 */
export const loadConfig = async (): Promise<void> => {
    try {
        const cfg: types.UiConfig | null | undefined = await FetchUiConfig();
        let activeCfg = cfg;
        if (!activeCfg) {
            const msg = 'UiConfig fetch returned null or undefined; using defaults.';
            LogError(msg);
            activeCfg = new types.UiConfig();
        }
        configStore.set(activeCfg);
    } catch (e: unknown) {
        const errMsg: string = e instanceof Error ? e.message : String(e);
        LogError(`Error loading config: ${errMsg}`);
    }
};
