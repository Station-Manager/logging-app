<script lang="ts">
    import {configStore} from '$lib/stores/config-store';
    import {qsoState} from "$lib/states/qso-state.svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {catState} from "$lib/states/cat-state.svelte";

    let isValid = $state(true);
    let isRandomQso: boolean = $derived.by((): boolean => {
        return qsoState.qso_random === 'Y';
    });
    let multiplierOn: boolean = $state($configStore.use_power_multiplier);

    let txPower: number = $derived.by((): number => {
        let power: number = parseInt(catState.txPower);
        if (Number.isNaN(power)) {
            power = $configStore.default_tx_power;
        }
        if (multiplierOn) {
            power = power * $configStore.power_multiplier;
        }
        return power;
    });

    const toggleRandonQso = (): void => {
        qsoState.qso_random = isRandomQso ? 'N' : 'Y';
    }

    const onblurTxPower = (event: Event): void => {
        const target = event.currentTarget as HTMLInputElement;
        let value = parseInt(target.value)
        if (txPower !== value) {
            $configStore.default_tx_power = value;
        }

    }
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
                        bind:value={sessionState.operatorCall}
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
            <div class="group relative inline-flex h-[24px] w-[44px] shrink-0 rounded-full bg-gray-300 p-0.5 outline-offset-2 outline-indigo-600 transition-colors duration-200 ease-in-out has-checked:bg-indigo-600 has-focus-visible:outline-2 dark:bg-white/5 dark:inset-ring-white/10 dark:outline-indigo-500 dark:has-checked:bg-indigo-500">
                <span class="size-5 rounded-full bg-white shadow-xs transition-transform duration-200 ease-in-out group-has-checked:translate-x-5"></span>
                <input
                        bind:checked={isRandomQso}
                        onclick={toggleRandonQso}
                        type="checkbox"
                        name="random_qso"
                        aria-label="Random QSO"
                        class="cursor-pointer absolute inset-0 appearance-none focus:outline-hidden"/>
            </div>
        </div>
        <div class="flex flex-col space-y-1.5 mt-3">
            <label class="block text-sm/5 font-semibold" for="tx_pwr">TX Power</label>
            <div class="flex items-center w-[142px]">
                <input
                        bind:value={txPower}
                        onblur={onblurTxPower}
                        class="mr-2 block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                        type="text"
                        id="tx_pwr"
                        maxlength="4"
                        minlength="1"
                        autocomplete="off"
                        title="Logging Station's Power">
                <div class="group relative inline-flex h-[24px] w-[44px] shrink-0 rounded-full bg-gray-200 p-0.5 inset-ring inset-ring-gray-900/5 outline-offset-2 outline-indigo-600 transition-colors duration-200 ease-in-out has-checked:bg-indigo-600 has-focus-visible:outline-2 dark:bg-white/5 dark:inset-ring-white/10 dark:outline-indigo-500 dark:has-checked:bg-indigo-500">
                    <span class="size-5 rounded-full bg-white shadow-xs ring-1 ring-gray-900/5 transition-transform duration-200 ease-in-out group-has-checked:translate-x-5"></span>
                    <input
                            bind:checked={multiplierOn}
                            type="checkbox" name="setting" aria-label="Use setting" class="absolute inset-0 appearance-none focus:outline-hidden" />
                </div>

            </div>
        </div>
    </div>
</div>
