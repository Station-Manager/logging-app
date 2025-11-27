<script lang="ts">
    import Callsign from "$lib/components/logging/Callsign.svelte";
    import Rst from "$lib/components/logging/Rst.svelte";
    import {qsoState} from "$lib/states/qso-state.svelte";
    import Mode from "$lib/components/logging/Mode.svelte";
    import {catStateValues} from "$lib/stores/cat-state-store";
    import TextInput from "$lib/components/logging/TextInput.svelte";
    import Comment from "$lib/components/logging/Comment.svelte";
    import TimeDatePanel from "$lib/components/logging/TimeDatePanel.svelte";
</script>

<div class="flex flex-row h-[351px]">
    <div class="flex flex-col gap-y-4 w-[724px] border">
        <div class="flex flex-row gap-x-4">
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
        <div class="flex flex-row gap-x-4">
            <TextInput
                    id="name"
                    label="Name"
                    value={qsoState.name}
            />
            <TextInput
                    id="qth"
                    label="Qth"
                    value={qsoState.qth}
                    overallWidthCss="w-[160px]"
            />
            <Comment
                    id="comment"
                    label="Comment"
                    value={qsoState.comment}
            />
        </div>
        <TimeDatePanel/>
    </div>
    <div class="w-[300px] border">Col 2</div>
</div>
