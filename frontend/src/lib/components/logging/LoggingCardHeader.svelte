<script lang="ts">
    import {type LoggingModeKey, loggingModes, loggingModeStore} from "$lib/stores/logging-mode-store";

    const modeEntries = Object.entries(loggingModes) as [LoggingModeKey, string][];
    const modeChange = (event: Event): void => {
        const select = event.target as HTMLSelectElement;
        loggingModeStore.set(select.value as LoggingModeKey);
    }
</script>

<header class="flex items-center h-[50px] px-4 border-b border-b-gray-300">
    <div class="flex flex-row items-center">
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
</header>
