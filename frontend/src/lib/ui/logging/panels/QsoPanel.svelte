<script lang="ts">
    import Callsign from "$lib/ui/logging/components/Callsign.svelte";
    import Rst from "$lib/ui/logging/components/Rst.svelte";
    import {qsoState} from "$lib/states/qso-state.svelte.js";
    import Mode from "$lib/ui/logging/components/Mode.svelte";
    import {catStateValues} from "$lib/stores/cat-state-store";
    import TextInput from "$lib/ui/logging/components/TextInput.svelte";
    import Comment from "$lib/ui/logging/components/Comment.svelte";
    import DateInput from "$lib/ui/logging/components/DateInput.svelte";
    import TimeInput from "$lib/ui/logging/components/TimeInput.svelte";
    import TimerControls from "$lib/ui/logging/components/TimerControls.svelte";
    import FormControls from "$lib/ui/logging/components/FormControls.svelte";
    import Vfos from "$lib/ui/logging/components/Vfos.svelte";
</script>

<div class="flex flex-row h-[351px]">
    <div class="flex flex-col gap-y-3 w-[744px] px-6">
        <div class="flex flex-row gap-x-4 items-center h-[100px]">
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
            <Vfos/>
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
                    overallWidthCss="w-[170px]"
            />
            <Comment
                    id="comment"
                    label="Comment"
                    value={qsoState.comment}
            />
        </div>
        <div class="flex flex-row gap-x-4 items-center -mt-8">
            <DateInput
                    id="qso_date"
                    label="Date"
                    value={qsoState.qsoDate}
            />
            <TimeInput
                    id="time_on"
                    label="Time On (UTC)"
                    value={qsoState.timeOn}
                    disabled={false}
            />
            <TimeInput
                    id="time_off"
                    label="Time Off (UTC)"
                    value={qsoState.timeOff}
                    disabled={true}
            />
            <div class="flex items-center mt-7">
                <TimerControls/>
                <FormControls/>
            </div>
        </div>
    </div>
    <div class="w-[280px] border">Col 2</div>
</div>
