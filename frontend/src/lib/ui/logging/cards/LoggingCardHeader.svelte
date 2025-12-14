<script lang="ts">
    import {
        type LoggingModeKey,
        loggingModes,
        loggingModeStore,
        isContestMode,
    } from "$lib/stores/logging-mode-store";
    import {configState} from "$lib/states/config-state.svelte";
    import {getSessionElapsedTime, sessionState} from "$lib/states/session-state.svelte";
    import {appState} from "$lib/states/app-state.svelte";
    import {STATION_PANEL} from "$lib/ui/logging/panels/constants";
    import {contestTimers} from "$lib/utils/contest-timers.svelte";
    import {formatTimeSecondsToHHColonMMColonSS} from "$lib/utils/time-date";
    import {contestState} from "$lib/states/contest-state.svelte";

    const modeEntries = Object.entries(loggingModes) as [LoggingModeKey, string][];

    const modeChange = (event: Event): void => {
        const select = event.target as HTMLSelectElement;
        loggingModeStore.set(select.value as LoggingModeKey);
        if (select.value === 'contest') {
            contestTimers.start();
        } else {
            contestTimers.stop();
        }
    }

    const operatorField = (): void => {
        appState.activePanel = STATION_PANEL
        const opField = document.getElementById("operator_callsign") as HTMLInputElement;
        if (opField) {
            opField.focus();
            opField.select();
        }
    }
</script>

<header class="flex items-center h-[50px] pl-4 border-b border-b-gray-300">
    <div class="flex flex-row items-center w-[290px]">
        <div class="text-md font-semibold w-[124px]">Logging Mode:</div>
        <div class="grid grid-cols-1">
            <select
                    class="text-sm col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-1 pr-8 pl-3 outline-1 -outline-offset-1 outline-gray-300 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600"
                    bind:value={$loggingModeStore}
                    onchange={modeChange}>
                {#each modeEntries as [key, label] (key)}
                    <option value={key}>{label}</option>
                {/each}
            </select>
            <svg viewBox="0 0 16 16" fill="currentColor" data-slot="icon" aria-hidden="true" class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end text-gray-500 dark:text-gray-400">
                <path d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
            </svg>
        </div>
    </div>
    <div class="w-[165px] text-xs">
        {#if $isContestMode}
        <div class="flex">
            <div class="w-16">Station:</div>
            <div class="font-bold">{configState.logbook.callsign}</div>
        </div>
        <div class="flex">
            <div class="w-16">Operator:</div>
            {#if sessionState.operator !== ''}
                <button
                        onclick={operatorField}
                        class="text-left font-bold {sessionState.operator === '' ? 'border border-red-500' : ''} min-w-14 cursor-pointer rounded-sm" title="Set operator">{sessionState.operator.toUpperCase()}</button>
            {:else}
                <div class="text-left font-bold">{configState.logbook.callsign}</div>
            {/if}
        </div>
        {/if}
    </div>
    <div class="w-[165px] text-xs">
        {#if $isContestMode}
            <div class="flex">
                <div class="w-12">QSOs:</div>
                <div>{contestState.totalQsos}</div>
            </div>
            <div class="flex">
                <div class="w-12">Last:</div>
                <div>{formatTimeSecondsToHHColonMMColonSS(contestTimers.elapasedSinceLastQso)}</div>
            </div>
        {/if}
    </div>
    <div class="flex flex-col text-xs font-semibold w-[200px]">
        <div class="flex flex-row items-center">
            <div class="w-[60px]">Logbook:</div>
            <div class="w-[110px]">{configState.logbook.name}</div>
        </div>
        <div class="flex flex-row items-center">
            <div class="w-[60px]">Rig:</div>
            <div class="w-[110px]">{configState.rig_name}</div>
        </div>
    </div>
    <div class="flex text-sm font-semibold w-[180px]">
        <div class="w-[110px]">Session Time:</div>
        <div class="w-[80px]">{getSessionElapsedTime()}</div>
    </div>
</header>
