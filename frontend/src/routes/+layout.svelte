<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onDestroy, onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {EventsOn} from "$lib/wailsjs/runtime/runtime";
    import {events} from "$lib/wailsjs/go/models";
    import {catState} from "$lib/states/cat-state.svelte";
    import {qsoState, type CatForQsoPayload} from "$lib/states/qso-state.svelte";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {Ready} from "$lib/wailsjs/go/facade/Service";

    let {children} = $props();
    let catStateEventsCancel: () => void = (): void => {}

    const registerForCatStateEvents = (): () => void => {
        console.log('registerForCatStateEvents()');
        return EventsOn(events.EventName.STATUS, (status: Record<string, string>) => {
            if (!status || Object.keys(status).length === 0) return;
            // First update our CAT snapshot
            catState.update(status);

            // Then map the latest CAT state into the subset of QSO fields driven by CAT.
            // For now, we keep the mapping simple and 1:1: VFO A → TX freq, VFO B → RX freq, mainMode → mode.
            const payload: CatForQsoPayload = {
                freq: catState.vfoaFreq,
                freq_rx: catState.vfobFreq,
                mode: catState.mainMode,
            };

            qsoState.updateFromCAT(payload);
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
