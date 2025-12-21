<script lang="ts">
    import {sessionState} from "$lib/states/session-state.svelte";
    import {parseDatabaseFreqToDottedKhz} from "$lib/utils/frequency";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {types} from "$lib/wailsjs/go/models";
    import {GetQsoById} from "$lib/wailsjs/go/facade/Service";
    import {qsoEditState} from "$lib/states/qso-edit-state.svelte";
    import Callsign from "$lib/ui/logging/components/Callsign.svelte";
    import Mode from "$lib/ui/logging/components/Mode.svelte";
    import {isContestMode} from "$lib/stores/logging-mode-store";
    import {catStateValues} from "$lib/stores/cat-state-store";
    import Rst from "$lib/ui/logging/components/Rst.svelte";

    const distanceCss = "w-[92px]";
    const timeCss = "w-[74px]";
    const callsignCss = "w-[90px]";
    const bandCss = "w-[50px]";
    const freqCss = "w-[80px]";
    const rstCss = "w-[50px]";
    const modeCss = "w-[52px]";
    const countryCss = "w-[140px] text-nowrap overflow-hidden text-ellipsis pr-1";
    const nameCss = "w-[140px] text-nowrap overflow-hidden text-ellipsis pr-1";

    let showEditPanel = $state(false);

    const editSessonQso = async (event: MouseEvent): Promise<void> => {
        const target = event.currentTarget as HTMLButtonElement | null;
        if (!target) {
            return;
        }
        const qsoId: number = Number(target.id);
        try {
            const qso: types.Qso = await GetQsoById(qsoId);
            qsoEditState.fromQso(qso);
            showEditPanel = true;
        } catch(e: unknown) {
            handleAsyncError(e, 'SessionPanel.svelte->editSessonQso')
        }
    }
</script>

<div class="cursor-default flex flex-col">
    <div class="flex flex-row border-b border-b-gray-300 font-semibold h-[32px] items-center px-4">
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
    <div class="flex flex-col overflow-y-scroll h-[298px] px-4">
        {#each sessionState.list as entry (entry.id)}
            <div class="flex flex-row even:bg-gray-300 text-sm h-[22px] p-0.5 rounded-xs">
                <div class={callsignCss}>{entry.call}</div>
                <div class={nameCss} title="{entry.name}">{entry.name}</div>
                <div class={freqCss}>{parseDatabaseFreqToDottedKhz(entry.freq)}</div>
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
{#if showEditPanel}
<div class="absolute top-[50px] w-full h-[701px] z-40 bg-gray-400/70 p-10">
    <div class="bg-white rounded-lg p-6 h-full w-full">
        <div class="flex flex-col gap-y-3 w-[744px] px-6">
            <div class="flex flex-row gap-x-4 items-center h-[100px]">
                <Callsign
                        id="call"
                        label="Callsign"
                        value={qsoEditState.call}
                />
                {#if $isContestMode}
                    {@render contextLogging()}
                {:else}
                    {@render normalLogging()}
                {/if}
                <Mode
                        id="mode"
                        label="Mode"
                        bind:value={qsoEditState.submode}
                        list={$catStateValues.getMainModes()}
                />
                {@render vfos()}
            </div>
        </div>
    </div>
</div>
{/if}

{#snippet normalLogging()}
    <Rst
            id="rst_sent"
            label="RST Sent"
            bind:value={qsoEditState.rst_sent}
    />
    <Rst
            id="rst_rcvd"
            label="RST Rcvd"
            bind:value={qsoEditState.rst_rcvd}
    />
{/snippet}

{#snippet contextLogging()}
    <div></div>
{/snippet}

{#snippet vfos()}
    <div class="flex flex-col w-[250px] h-[80px] mt-6 gap-y-2">
        <div class="flex flex-row items-center">
            <label for="freq" class="">Freq (RX)</label>
            <div class="w-[116px]">
                <input
                        type="text"
                        autocomplete="off"
                        spellcheck="false"
                        id="freq"
                        value={qsoEditState.freq_rx}
                        title="Format: ?#.###.###"
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                />
            </div>
        </div>
        <div class="flex flex-row items-center">
            <label for="freq_rx" class="">Freq (TX)</label>
            <div class="w-[116px]">
                <input
                        type="text"
                        autocomplete="off"
                        spellcheck="false"
                        id="freq_rx"
                        value={qsoEditState.freq}
                        title="Format: ?#.###.###"
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                />
            </div>
        </div>
    </div>
{/snippet}