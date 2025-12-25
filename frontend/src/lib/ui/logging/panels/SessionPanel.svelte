<script lang="ts">
    import {sessionState} from "$lib/states/session-state.svelte";
    import {frequencyToBandFromDottedMHz, parseDatabaseFreqToDottedKhz} from "$lib/utils/frequency";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {types} from "$lib/wailsjs/go/models";
    import {GetQsoById} from "$lib/wailsjs/go/facade/Service";
    import {qsoEditState} from "$lib/states/qso-edit-state.svelte";
    import Callsign from "$lib/ui/logging/components/Callsign.svelte";
    import Mode from "$lib/ui/logging/components/Mode.svelte";
    import {isContestMode} from "$lib/stores/logging-mode-store";
    import {catStateValues} from "$lib/stores/cat-state-store";
    import Rst from "$lib/ui/logging/components/Rst.svelte";
    import TextInput from "$lib/ui/logging/components/TextInput.svelte";
    import DateInput from "$lib/ui/logging/components/DateInput.svelte";
    import TimeInput from "$lib/ui/logging/components/TimeInput.svelte";
    import {isValidCallsignForLog} from "$lib/constants/callsign";
    import {UpdateQso} from "$lib/wailsjs/go/facade/Service";
    import {showToast} from "$lib/utils/toast";

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
    let isUpdating: boolean = $state(false);

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

    const cancelAction = (): void => {
        showEditPanel = false;
    }

    const updateAction = async (): Promise<void> => {
        if (!canLog()) {
            const elem = document.getElementById('call') as HTMLInputElement;
            if (elem) {
                elem.focus();
                elem.select();
            }
            return;
        }

        if (isUpdating) return; // Prevent double-clicks
        isUpdating = true;

        const qso: types.Qso = qsoEditState.toQso();
        try {
            await UpdateQso(qso);
            showToast.SUCCESS("QSO updated...");
        } catch(e: unknown) {
            handleAsyncError(e, 'SessionPanel.svelte->updateAction')
        } finally {
            isUpdating = false;
            showEditPanel = false;
        }
    }

    const canLog = (): boolean => {
        return isValidCallsignForLog(qsoEditState.call)
    };

</script>

<div class="cursor-default flex flex-col">
    <div class="flex flex-row border-b border-b-gray-300 font-semibold h-8 items-center px-4">
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
<div class="absolute top-[50px] w-full h-[701px] z-40 bg-gray-400/70">
    <div class="bg-white rounded-lg py-8 px-14 h-[460px] w-[856px] mt-24 mx-auto">
        <div class="flex flex-col gap-y-3 w-[744px] h-[340px] px-6">
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
            <div class="flex flex-row gap-x-4 -mt-1.5">
                <TextInput
                        id="name"
                        label="Name"
                        bind:value={qsoEditState.name}
                />
                <TextInput
                        id="qth"
                        label="Qth"
                        bind:value={qsoEditState.qth}
                        overallWidthCss="w-[170px]"
                />
                <div class="w-[280px]">
                    <label for="comment" class="block text-sm/5 font-medium w-[70px]">Comments</label>
                    <div class="mt-2">
                    <textarea
                            bind:value={qsoEditState.comment}
                            id="comment"
                            spellcheck="false"
                            class="h-16 resize-none w-full rounded-md bg-white px-2 py-1.5 text-sm outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                            autocomplete="off"></textarea>
                    </div>
                </div>
            </div>
            <div class="flex flex-row gap-x-4 items-top">
                <DateInput
                        id="qso_date"
                        label="Date"
                        bind:value={qsoEditState.qso_date}
                />
                <TimeInput
                        id="time_on"
                        label="Time On (UTC)"
                        bind:value={qsoEditState.time_on}
                        disabled={false}
                />
                <TimeInput
                        id="time_off"
                        label="Time Off (UTC)"
                        bind:value={qsoEditState.time_off}
                        disabled={false}
                />
                <div class="w-[280px]">
                    <label for="notes" class="block text-sm/5 font-medium w-[70px]">Notes</label>
                    <div class="mt-2">
                    <textarea
                            bind:value={qsoEditState.notes}
                            id="notes"
                            spellcheck="false"
                            class="h-16 resize-none w-full rounded-md bg-white px-2 py-1.5 text-sm outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                            autocomplete="off"></textarea>
                    </div>
                </div>
            </div>
            <div class="flex flex-row space-x-4 -mt-6">
                <div>
                    <label class="block text-sm/5 font-medium" for="rx_pwr">Power</label>
                    <div class="mt-2 w-[100px]">
                        <input
                                bind:value={qsoEditState.rx_pwr}
                                class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                                type="text"
                                id="rx_pwr"
                                placeholder="RX Power"
                                autocomplete="off"
                                title="Contacted Station's Power">
                    </div>
                </div>
                <div>
                    <label class="block text-sm/5 font-medium" for="rig">Rig</label>
                    <div class="mt-2 w-[360px]">
                    <textarea
                            bind:value={qsoEditState.rig}
                            class="resize-none w-full rounded-md bg-white px-3 py-1.5 text-sm outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                            id="rig"
                            placeholder="Working conditions"></textarea>
                    </div>
                </div>
            </div>
        </div>
        <div class="flex w-full gap-x-3 justify-end">
            <button
                    onclick={updateAction}
                    id="update-contact-btn"
                    type="button"
                    class="disabled:bg-gray-400 disabled:cursor-not-allowed h-9 cursor-pointer rounded-md bg-indigo-600 px-2.5 py-1.5 text-base font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                    title="Ctrl-s">Update QSO
            </button>
            <button
                    onclick={cancelAction}
                    type="button"
                    class="h-9 w-[74px] cursor-pointer rounded-md bg-white px-2.5 py-1.5 text-base font-semibold ring-1 shadow-sm ring-gray-300 ring-inset hover:bg-gray-100"
                    title="ESC">Cancel
            </button>
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
    <div class="flex flex-col w-[250px] h-20 mt-6 gap-y-2">
        <div class="flex flex-row items-center">
            <label for="freq" class="text-sm/5 font-medium w-[70px]">Freq (RX)</label>
            <div class="w-[116px]">
                <input
                        type="text"
                        autocomplete="off"
                        spellcheck="false"
                        id="freq"
                        value={parseDatabaseFreqToDottedKhz(qsoEditState.freq_rx)}
                        title="Format: ?#.###.###"
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                />
            </div>
            <div class="cursor-default w-8 font-semibold text-base ml-2">{frequencyToBandFromDottedMHz(qsoEditState.freq_rx)}</div>
        </div>
        <div class="flex flex-row items-center">
            <label for="freq_rx" class="text-sm/5 font-medium w-[70px]">Freq (TX)</label>
            <div class="w-[116px]">
                <input
                        type="text"
                        autocomplete="off"
                        spellcheck="false"
                        id="freq_rx"
                        value={parseDatabaseFreqToDottedKhz(qsoEditState.freq)}
                        title="Format: ?#.###.###"
                        class="block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
                />
            </div>
            <div class="cursor-default w-8 font-semibold text-base ml-2">{frequencyToBandFromDottedMHz(qsoEditState.freq)}</div>
        </div>
    </div>
{/snippet}