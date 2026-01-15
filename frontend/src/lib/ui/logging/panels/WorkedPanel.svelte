<script lang="ts">
    import {qsoState} from "$lib/states/new-qso-state.svelte";
    import {formatDate, formatTime} from "../../../../../../../shared-utils/src/lib/utils/time-date";
    import {parseDatabaseFreqToDottedKhz} from "../../../../../../../shared-utils/src/lib/utils/frequency";
</script>

<div class="cursor-default flex flex-col">
    <div class="flex flex-row border-b border-b-gray-300 font-semibold h-[32px] items-center px-4">
        <div class="w-[114px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Date</div>
        <div class="w-[116px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Call</div>
        <div class="w-[48px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Time</div>
        <div class="w-[80px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Freq</div>
        <div class="w-[64px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Band</div>
        <div class="w-[68px] text-left hover:bg-gray-300 px-1 mr-1 rounded-xs">Mode</div>
        <div class="w-[60px] text-left">Sent</div>
        <div class="w-[60px] text-left">Rcvd</div>
        <div class="">Notes</div>
    </div>
    <div class="flex flex-col overflow-y-scroll h-[298px] px-4">
        {#each qsoState.contact_history as entry (entry.id)}
            <div class="flex flex-row even:bg-gray-300 text-sm h-[22px] p-0.5 rounded-xs">
                <div class="w-[120px] text-left">{formatDate(entry.qso_date)}</div>
                <div class="w-[120px] text-left">{entry.call}</div>
                <div class="w-[50px] text-left">{formatTime(entry.time_on)}</div>
                <div class="w-[86px] text-left">{parseDatabaseFreqToDottedKhz(entry.freq)}</div>
                <div class="w-[68px] text-left">{entry.band}</div>
                <div class="w-[68px] text-left">{entry.mode}</div>
                <div class="w-[60px] text-left">{entry.rst_sent}</div>
                <div class="w-[60px] text-left">{entry.rst_rcvd}</div>
                <div>{entry.notes}</div>
            </div>
        {/each}
    </div>
</div>