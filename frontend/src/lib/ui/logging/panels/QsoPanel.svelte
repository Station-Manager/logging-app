<script lang="ts">
    import Callsign from "$lib/ui/logging/components/Callsign.svelte";
    import Rst from "$lib/ui/logging/components/Rst.svelte";
    import {qsoState} from "$lib/states/qso-state.svelte";
    import Mode from "$lib/ui/logging/components/Mode.svelte";
    import {catStateValues} from "$lib/stores/cat-state-store";
    import TextInput from "$lib/ui/logging/components/TextInput.svelte";
    import Comment from "$lib/ui/logging/components/Comment.svelte";
    import DateInput from "$lib/ui/logging/components/DateInput.svelte";
    import TimeInput from "$lib/ui/logging/components/TimeInput.svelte";
    import TimerControls from "$lib/ui/logging/components/TimerControls.svelte";
    import FormControls from "$lib/ui/logging/components/FormControls.svelte";
    import Vfos from "$lib/ui/logging/components/Vfos.svelte";
    import CountryPanel from "$lib/ui/logging/panels/CountryPanel.svelte";
</script>

<div class="flex flex-row h-[281px]">
    <div class="flex flex-col gap-y-3 w-[744px] px-6">
        <div class="flex flex-row gap-x-4 items-center h-[100px]">
            <Callsign
                    id="call"
                    label="Callsign"
                    value={qsoState.call}
            />
            <Rst
                    id="rst_sent"
                    label="RST Sent"
                    bind:value={qsoState.rst_sent}
            />
            <Rst
                    id="rst_rcvd"
                    label="RST Rcvd"
                    bind:value={qsoState.rst_rcvd}
            />
            <Mode
                    id="mode"
                    label="Mode"
                    bind:value={qsoState.mode}
                    selected={qsoState.cat_main_mode}
                    list={$catStateValues.getMainModes()}
            />
            <Vfos/>
        </div>
        <div class="flex flex-row gap-x-4">
            <TextInput
                    id="name"
                    label="Name"
                    bind:value={qsoState.name}
            />
            <TextInput
                    id="qth"
                    label="Qth"
                    bind:value={qsoState.qth}
                    overallWidthCss="w-[170px]"
            />
            <Comment
                    id="comment"
                    label="Comment"
                    bind:value={qsoState.comment}
            />
        </div>
        <div class="flex flex-row gap-x-4 items-center -mt-8">
            <DateInput
                    id="qso_date"
                    label="Date"
                    bind:value={qsoState.qso_date}
            />
            <TimeInput
                    id="time_on"
                    label="Time On (UTC)"
                    bind:value={qsoState.time_on}
                    disabled={false}
            />
            <TimeInput
                    id="time_off"
                    label="Time Off (UTC)"
                    bind:value={qsoState.time_off}
                    disabled={false}
            />
            <div class="flex items-center mt-7">
                <TimerControls/>
                <FormControls/>
            </div>
        </div>
    </div>
    <div class="flex w-[280px] pt-6 pl-2">
        <CountryPanel/>
    </div>
</div>
