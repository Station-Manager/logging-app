<script lang="ts">
    import { configStore } from '$lib/stores/config-store';

    let isValid = $state(true);
</script>

<div class="cursor-default flex flex-row px-5">
    <div class="flex flex-col w-1/5">
        <div class="mt-3">
            <label for="station_callsign" class="flex flex-row text-sm/5 font-semibold">
                <span>Station's callsign</span>
            </label>
            <div class="mt-2 w-[140px]">
                <input
                        bind:value={$configStore.logbook.callsign}
                        type="text"
                        id="station_callsign"
                        class="uppercase block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 {isValid ? '' : 'outline-red-600 focus:outline-red-600'}"
                        autocomplete="off"
                        spellcheck="false"
                        title="Logging Station's Callsign"
                        disabled
                />
            </div>
        </div>
        <div class="mt-3">
            <label for="owner_callsign" class="flex flex-row text-sm/5 font-semibold">
                <span>Owner's callsign</span>
            </label>
            <div class="mt-2 w-[140px]">
                <input
                        bind:value={$configStore.owners_callsign}
                        type="text"
                        id="owner_callsign"
                        class="uppercase block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 {isValid ? '' : 'outline-red-600 focus:outline-red-600'}"
                        autocomplete="off"
                        spellcheck="false"
                        title="Owner of the logging Station"
                        disabled
                />
            </div>
        </div>
        <div class="mt-3">
            <label for="operator_callsign" class="flex flex-row text-sm/5 font-semibold">
                <span>Operator's callsign</span>
            </label>
            <div class="mt-2 w-[140px]">
                <input
                        type="text"
                        id="operator_callsign"
                        class="uppercase block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 {isValid ? '' : 'outline-red-600 focus:outline-red-600'}"
                        autocomplete="off"
                        spellcheck="false"
                        title="Guest Operator's Callsign"
                />
            </div>
        </div>
    </div>
    <div class="flex flex-col w-1/5">
        <div class="flex flex-row mt-3">
            <label class="mr-2 mt-0.5 text-sm font-semibold" for="random_qso">Random QSO</label>
            <div class="group relative inline-flex w-[44px] shrink-0 rounded-full bg-gray-300 p-0.5 outline-offset-2 outline-indigo-600 transition-colors duration-200 ease-in-out has-checked:bg-indigo-600 has-focus-visible:outline-2 dark:bg-white/5 dark:inset-ring-white/10 dark:outline-indigo-500 dark:has-checked:bg-indigo-500">
                <span class="size-5 rounded-full bg-white shadow-xs transition-transform duration-200 ease-in-out group-has-checked:translate-x-5"></span>
                <input
                        type="checkbox"
                        name="random_qso"
                        aria-label="Random QSO"
                        class="absolute inset-0 appearance-none focus:outline-hidden" />
            </div>
        </div>
        <div class="flex flex-col space-y-1.5 mt-3">
            <label class="block text-sm/5 font-semibold" for="tx_pwr">TX Power</label>
            <div class="relative mt-2 w-[146px]">
                <input
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                        type="text"
                        id="tx_pwr"
                        placeholder="Power"
                        autocomplete="off"
                        title="Logging Station's Power">
                <button
                        id="pwrbtn"
                        type="button"
                        class="absolute top-1.5 right-2 inline-flex h-6 w-[44px] shrink-0 cursor-pointer rounded-full border border-transparent transition-colors duration-200 ease-in-out focus:outline-hidden"
                        role="switch"
                        aria-label="switch"
                        aria-checked="false"
                        title="Linear Amp Multipier: {$configStore.power_multiplier}x">
                        <span
                                aria-hidden="true"
                                class="mt-[1px] ml-[1px] pointer-events-none inline-block size-5 transform rounded-full bg-white ring-0 transition duration-200 ease-in-out"></span>
                </button>
            </div>
        </div>
    </div>
</div>

