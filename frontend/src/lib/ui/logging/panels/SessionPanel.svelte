<script lang="ts">
    import {sessionState} from "$lib/states/session-state.svelte";

    const distanceCss = "w-[92px]";
    const timeCss = "w-[74px]";
    const callsignCss = "w-[90px]";
    const bandCss = "w-[50px]";
    const freqCss = "w-[70px]";
    const rstCss = "w-[50px]";
    const modeCss = "w-[52px]";
    const countryCss = "w-[140px] text-nowrap overflow-hidden text-ellipsis pr-1";
    const nameCss = "w-[140px] text-nowrap overflow-hidden text-ellipsis pr-1";

    const editSessonQso = async (event: MouseEvent): Promise<void> => {
        const target = event.currentTarget as HTMLButtonElement | null;
        if (!target) {
            return;
        }
    }
</script>

<div class="cursor-default flex flex-col">
    <div class="flex flex-row border-b border-b-gray-300 font-semibold h-[32px] items-center px-5">
        <div class={callsignCss}>Callsign</div>
        <div class={nameCss}>Name</div>
        <div class={freqCss}>Freq</div>
        <div class={bandCss}>Band</div>
        <div class={rstCss}>Send</div>
        <div class={rstCss}>Rcvd</div>
        <div class={modeCss}>Mode</div>
        <div class={timeCss}>Time On</div>
        <div class={countryCss}>Country</div>
        <div class="w-[130px]">Distance</div>
    </div>
    <div class="relative h-[264px] overflow-y-scroll pt-1 flex flex-col text-sm px-5">
        {#each sessionState.list as entry (entry.id)}
            <div class="flex flex-row odd:bg-white even:bg-gray-300">
                <div class={callsignCss}>{entry.call}</div>
                <div class={nameCss} title="{entry.name}">{entry.name}</div>
                <div class={freqCss}>{entry.freq}</div>
                <div class={bandCss}>{entry.band}</div>
                <div class={rstCss}>{entry.rst_sent}</div>
                <div class={rstCss}>{entry.rst_sent}</div>
                <div class={modeCss}>{entry.mode}</div>
                <div class={timeCss}>{entry.time_on}</div>
                <div class={countryCss} title="{entry.country}">{entry.country}</div>
                <div class={distanceCss}>{entry.distance} km ({entry.ant_path})</div>
                <div><button onclick={editSessonQso} id={entry.id.toString()} class="cursor-pointer font-semibold hover:text-indigo-700">Edit</button></div>
            </div>
        {/each}
    </div>
</div>
