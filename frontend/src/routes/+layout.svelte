<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/ui/MainNav.svelte";
    import {onMount} from "svelte";
    import {sessionState} from "$lib/states/session-state.svelte";

    let {children} = $props();

    onMount(async (): Promise<void> => {
        sessionState.start();
    });
</script>

<header>
    <SvelteToast/>
    <MainNav/>
</header>
<main>
    {@render children()}
</main>
