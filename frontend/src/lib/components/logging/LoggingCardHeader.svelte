<script lang="ts">
    import {
        type LoggingModeKey,
        loggingModes,
        loggingModeStore,
        isContestMode,
    } from "$lib/stores/logging-mode-store";
    import {configStore} from "$lib/stores/config-store";

    const modeEntries = Object.entries(loggingModes) as [LoggingModeKey, string][];
    const modeChange = (event: Event): void => {
        const select = event.target as HTMLSelectElement;
        loggingModeStore.set(select.value as LoggingModeKey);
    }
</script>

<header class="flex items-center h-[50px] px-4 border-b border-b-gray-300">
    <div class="flex flex-row items-center w-[290px]">
        <div class="text-md font-semibold w-[124px]">Logging Mode:</div>
        <div>
            <select
                    class="border border-gray-300 rounded-md px-2 text-sm"
                    bind:value={$loggingModeStore}
                    onchange={modeChange}>
                {#each modeEntries as [key, label] (key)}
                    <option value={key}>{label}</option>
                {/each}
            </select>
        </div>
    </div>
    <div class="w-[140px]">
        {#if $isContestMode}
            Station/Op
        {/if}
    </div>
    <div class="w-[140px]">
        {#if $isContestMode}
            QSO count/Last
        {/if}
    </div>
    <div class="flex flex-col text-xs font-semibold w-[180px]">
        <div class="flex flex-row items-center">
            <div class="w-[60px]">Logbook:</div>
            <div class="w-[110px]">{$configStore.logbook.name}</div>
        </div>
        <div class="flex flex-row items-center">
            <div class="w-[60px]">Rig:</div>
            <div class="w-[110px]">{$configStore.rig_name}</div>
        </div>
    </div>
    <div class="w-[180px]">
        Session
    </div>
</header>
