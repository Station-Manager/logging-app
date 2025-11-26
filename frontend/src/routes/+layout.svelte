<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";
    import MainNav from "$lib/components/MainNav.svelte";

    let {children} = $props();
</script>

<header>
    <SvelteToast/>
    <MainNav/>
</header>
<main>
    {@render children()}
</main>
