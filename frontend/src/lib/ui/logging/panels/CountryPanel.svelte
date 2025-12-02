<script lang="ts">
    import {qsoState} from "$lib/states/qso-state.svelte";

    let isVisible = $derived(qsoState.country_name !== '');
</script>

<div class="w-[240px] h-[236px] border border-gray-300 rounded-md bg-gray-200">
    <div class="{isVisible ? 'block' : 'hidden'} flex flex-col">
        <div class="text-2xl text-center font-semibold py-2">
            {qsoState.country_name}
        </div>
    </div>
</div>
