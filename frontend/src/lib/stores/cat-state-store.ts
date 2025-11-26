import { writable, type Writable } from 'svelte/store';
import { FetchCatStateValues } from '$lib/wailsjs/go/facade/Service';
import { handleAsyncError } from '$lib/utils/error-handler';

export class StateValuesClass {
    private readonly data: Record<string, Record<string, string>>;

    constructor(data: Record<string, Record<string, string>>) {
        this.data = data;
    }

    public getModes(): { key: string; value: string }[] {
        const modes = this.data['MAINMODE'] ?? {};
        return Object.entries(modes).map(([key, value]) => ({ key, value }));
    }
}

export const catStateValues: Writable<StateValuesClass> = writable(new StateValuesClass({}));

export const loadCatStateValues = async (): Promise<void> => {
    try {
        const data: Record<string, Record<string, string>> = await FetchCatStateValues();
        // for (const [groupKey, groupValues] of Object.entries(data)) {
        //     console.log(`Group: ${groupKey}`);
        //     for (const [key, value] of Object.entries(groupValues)) {
        //         console.log(`  ${key}: ${value}`);
        //     }
        // }
        catStateValues.set(new StateValuesClass(data));
    } catch (e: unknown) {
        handleAsyncError(e, 'loadCatStateValues()');
    }
};
