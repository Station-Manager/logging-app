import type { LayoutData } from './$types';
import { loadConfig } from '$lib/states/config-state.svelte';

export const prerender = true;
export const ssr = false;

export const load: LayoutData = async (): Promise<object> => {
    try {
        await loadConfig();
    } catch (e: unknown) {
        const errorMessage: string = e instanceof Error ? e.message : String(e);
        console.error('Error loading config:', errorMessage);
    }
    return {};
};
