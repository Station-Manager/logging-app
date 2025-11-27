<script lang="ts">

import {qsoState} from "$lib/states/qso-state.svelte";
import TimeInput from "$lib/components/logging/TimeInput.svelte";
import DateInput from "$lib/components/logging/DateInput.svelte";
</script>

<div class="flex flex-row gap-x-4">
    <DateInput
        id="qso_date"
        label="Date"
        value={qsoState.qsoDate}
        />
    <TimeInput
            id="time_on"
            label="Time On"
            value={qsoState.timeOn}
            disabled={false}
    />
    <TimeInput
            id="time_off"
            label="Time Off"
            value={qsoState.timeOff}
            disabled={false}
    />
</div>
