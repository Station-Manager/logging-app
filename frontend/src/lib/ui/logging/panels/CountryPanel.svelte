<script lang="ts">
    import {qsoState} from "$lib/states/qso-state.svelte";

    let isVisible = $derived(qsoState.country_name !== '');

    let flagImgPath = $derived.by((): string => {
        return `/flags/unknown.svg`;
//        return `flags/${qsoState.country_name}.svg`;
    });
</script>

<div class="w-[240px] h-[236px] border border-gray-300 rounded-md bg-gray-200">
    <div class="{isVisible ? 'block' : 'hidden'} flex flex-col">
        <div class="text-2xl text-center font-semibold py-2">
            {qsoState.country_name}
        </div>
        <div class="flex justify-center">
            <img class="w-[80px] border border-gray-400/40" src="{flagImgPath}" alt="{qsoState.country_name}">
        </div>
        <div>
            {#if qsoState.short_path_distance === '' && qsoState.long_path_distance === ''}
                <div class="text-center text-gray-500">No data</div>
            {:else}
                {#if qsoState.ant_path === 'S'}
                    <span class="text-red-600">{qsoState.short_path_distance} km</span>
                {:else}
                    <span>{qsoState.short_path_distance} km</span>
                {/if}
                /
                {#if qsoState.ant_path === 'L'}
                    <span class="text-red-600">{qsoState.long_path_distance} km</span>
                {:else}
                    <span>{qsoState.long_path_distance} km</span>
                {/if}
            {/if}
        </div>
    </div>
</div>
