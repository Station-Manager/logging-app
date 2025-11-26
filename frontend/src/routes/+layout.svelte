<script lang="ts">
    import "./layout.css";
    import {SvelteToast} from "@zerodevx/svelte-toast";

    let {children} = $props();
</script>

<div>
    <SvelteToast/>
</div>
<main>
    {@render children()}
</main>
