<script lang="ts">
    import Callsign from "$lib/components/logging/Callsign.svelte";
    import Rst from "$lib/components/logging/Rst.svelte";
    import {qsoState} from "$lib/states/qso-state.svelte";
    import Mode from "$lib/components/logging/Mode.svelte";
    import {catStateValues} from "$lib/stores/cat-state-store";
</script>

<div class="flex flex-row h-[351px]">
    <div class="flex gap-x-4 w-[724px]">
        <Callsign
                id="call"
                label="Callsign"
                value={qsoState.callsign}
        />
        <Rst
                id="rst_sent"
                label="RST Sent"
                value={qsoState.rstSent}
        />
        <Rst
                id="rst_rcvd"
                label="RST Rcvd"
                value={qsoState.rstRcvd}
        />
        <Mode
                id="mode"
                label="Mode"
                value={qsoState.mainmode}
                list={$catStateValues.getMainModes()}
        />
    </div>
    <div class="w-[300px] border">Col 2</div>
</div>
