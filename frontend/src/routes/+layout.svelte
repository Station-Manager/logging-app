<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onDestroy, onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";
    import {EventsOn} from "$lib/wailsjs/runtime/runtime";

    let {children} = $props();
    let catStateEventsCancel: () => void;

    const registerForCatStateEvents = (): () => void => {
        return EventsOn("STATUS", (value: Record<string, string>) => {
            if (!value || Object.keys(value).length === 0) {
                return;
            }
        });
    }

    onMount(async (): Promise<void> => {
        sessionState.start();
        catStateEventsCancel = registerForCatStateEvents();
    });

    onDestroy((): void => {
        catStateEventsCancel();
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
