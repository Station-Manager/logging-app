<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onDestroy, onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {EventsOn} from "$lib/wailsjs/runtime/runtime";
    import {events} from "$lib/wailsjs/go/models";
    import {catState} from "$lib/states/cat-state.svelte";
    import {qsoState, type CatForQsoPayload} from "$lib/states/new-qso-state.svelte";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {Ready} from "$lib/wailsjs/go/facade/Service";

    let {children} = $props();
    let catStateEventsCancel: () => void = (): void => {}

    const registerForCatStateEvents = (): () => void => {
        return EventsOn(events.EventName.STATUS, (status: Record<string, string>) => {
            if (!status || Object.keys(status).length === 0) return;

            // First update our CAT snapshot
            catState.update(status);

            // Then map the latest CAT state into the subset of QSO fields driven by CAT.
            const payload: CatForQsoPayload = {
                // ADIF-aligned fields (these may be persisted if user saves the QSO)
                //freq: catState.vfoaFreq,
                //freq_rx: catState.vfobFreq,
                //mode: catState.mainMode,
                // band / band_rx can be added when you have a freq->band mapping available

                // CAT-only, UI-facing mirrors
                // cat_identity: catState.identity,
                cat_vfoa_freq: catState.vfoaFreq,
                cat_vfob_freq: catState.vfobFreq,
                // cat_vfob_freq: formatCatKHzToDottedMHz(catState.vfobFreq),
                // cat_select: catState.select,
                // cat_split: catState.split,
                cat_main_mode: catState.mainMode,
                // cat_sub_mode: catState.subMode,
                // cat_tx_power: catState.txPower,
            };

            // qsoState.setDefaults($configStore);
            // console.log('+layout mounted', $state.snapshot(defaultInputs), $configStore);
            qsoState.updateFromCAT(payload);
            console.log($state.snapshot(qsoState));
        });
    }

    onMount(async (): Promise<void> => {
        sessionState.start();
        catStateEventsCancel = registerForCatStateEvents();
        try {
            await Ready();
        } catch (e: unknown) {
            handleAsyncError(e, '+layout.svelte->onMount')
        }

    });

    onDestroy((): void => {
        catStateEventsCancel();
        catStateEventsCancel = () => {};
        sessionState.stop();
    });
</script>

<header>
    <SvelteToast/>
    <MainNav/>
</header>
<main>
    {@render children()}
</main>
