<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onDestroy, onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {EventsOn} from "$lib/wailsjs/runtime/runtime";
    import {events} from "$lib/wailsjs/go/models";
    import {catState} from "$lib/states/cat-state.svelte";
    import {handleAsyncError} from "$lib/utils/error-handler";
    import {Ready} from "$lib/wailsjs/go/facade/Service";

    let {children} = $props();
    let catStateEventsCancel: () => void = (): void => {}

    const registerForCatStateEvents = (): () => void => {
        return EventsOn(events.EventName.STATUS, (status: Record<string, string>) => {
            if (!status || Object.keys(status).length === 0) return;
            catState.update(status);
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
