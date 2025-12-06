<script lang="ts">
    import {qsoState} from "$lib/states/new-qso-state.svelte";

    const SHORT_PATH = 'S';
    const LONG_PATH = 'L';

    let shortPathRadio: HTMLInputElement;
    let longPathRadio: HTMLInputElement;

    let isVisible = $derived(qsoState.country_name !== '');

    let flagImgPath = $derived.by((): string => {
        if (qsoState.ccode === ''){
            return `/flags/unknown.svg`;
        }
        return `flags/${qsoState.ccode.toLowerCase()}.svg`;
    });

    const toggleAntPath = (event: Event):void => {
        const target = event.currentTarget as HTMLInputElement;
        if (!target) return;
        if (target.value === SHORT_PATH) {
            qsoState.ant_path = SHORT_PATH;
            longPathRadio.checked = false;

        } else if (target.value === LONG_PATH) {
            qsoState.ant_path = LONG_PATH;
            shortPathRadio.checked = false;
        }
    };
</script>

<div class="w-[240px] h-[236px] border border-gray-300 rounded-md bg-gray-200">
    <div class="{isVisible ? 'block' : 'hidden'} flex flex-col">
        <div class="text-2xl text-center font-semibold py-3">
            {qsoState.country_name}
        </div>
        <div class="flex justify-center">
            <img class="w-[80px] border border-gray-400/40" src="{flagImgPath}" alt="{qsoState.country_name}">
        </div>
        <div class="text-center py-2">
            <div class="">
            {#if qsoState.short_path_distance === '' && qsoState.long_path_distance === ''}
                <div class="text-center text-gray-500">No distance data</div>
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
            <div class="flex flex-row h-10 w-full space-x-3 justify-center -mt-1">
                <div class="flex items-center">
                    <input onclick={toggleAntPath} bind:this={shortPathRadio} value={SHORT_PATH} id="short_path" type="radio" checked class="relative size-4 appearance-none rounded-full border border-gray-400 bg-white before:absolute before:inset-1 before:rounded-full before:bg-white not-checked:before:hidden checked:border-indigo-600 checked:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:before:bg-gray-400 forced-colors:appearance-auto forced-colors:before:hidden">
                    <label for="short_path" class="ml-1 block text-sm font-medium">Short path</label>
                </div>
                <div class="flex items-center">
                    <input onclick={toggleAntPath} bind:this={longPathRadio} value={LONG_PATH} id="long_path" type="radio" class="relative size-4 appearance-none rounded-full border border-gray-400 bg-white before:absolute before:inset-1 before:rounded-full before:bg-white not-checked:before:hidden checked:border-indigo-600 checked:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:before:bg-gray-400 forced-colors:appearance-auto forced-colors:before:hidden">
                    <label for="long_path" class="ml-1 block text-sm font-medium">Long path</label>
                </div>
            </div>
        </div>
        <div class="flex flex-row font-semibold justify-center">
            <div class="w-[100px]">Local time:</div>
            <div class="text-sm pt-1 text-red-600">{qsoState.remote_time} ({qsoState.remote_offset})</div>
        </div>
    </div>
</div>
