<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onDestroy, onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {EventsOn} from "$lib/wailsjs/runtime/runtime";
    import {events} from "$lib/wailsjs/go/models";
    import {catState} from "$lib/states/cat-state.svelte";

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
