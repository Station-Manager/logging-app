<script lang="ts">
    import {qsoState} from "$lib/states/new-qso-state.svelte";
    import {formatDate, formatTime} from "@station-manager/shared-utils";
    import {parseDatabaseFreqToDottedKhz} from "@station-manager/shared-utils";
</script>

<div class="cursor-default flex flex-col">
    <div class="flex flex-row border-b border-b-gray-300 font-semibold h-8 items-center px-4">
        <div class="w-28.5 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Date</div>
        <div class="w-29 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Call</div>
        <div class="w-12 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Time</div>
        <div class="w-20 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Freq</div>
        <div class="w-16 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Band</div>
        <div class="w-17 text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Mode</div>
        <div class="w-15 text-left">Sent</div>
        <div class="w-15 text-left">Rcvd</div>
        <div class="">Notes</div>
    </div>
    <div class="flex flex-col overflow-y-scroll h-74.5 px-4">
        {#each qsoState.contact_history as entry (entry.id)}
            <div class="flex flex-row even:bg-gray-300 text-sm h-5.5 p-0.5 rounded-xs">
                <div class="w-30 text-left">{formatDate(entry.qso_date)}</div>
                <div class="w-30 text-left">{entry.call}</div>
                <div class="w-12.5 text-left">{formatTime(entry.time_on)}</div>
                <div class="w-21.5 text-left">{parseDatabaseFreqToDottedKhz(entry.freq)}</div>
                <div class="w-17 text-left">{entry.band}</div>
                <div class="w-17 text-left">{entry.mode}</div>
                <div class="w-15 text-left">{entry.rst_sent}</div>
                <div class="w-15 text-left">{entry.rst_rcvd}</div>
                <div>{entry.notes}</div>
            </div>
        {/each}
    </div>
</div>