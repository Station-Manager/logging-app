<script lang="ts">
    import {qsoState, resetQsoStateDefaults} from "$lib/states/qso-state.svelte";

    const resetAction = (): void => {
        console.log('resetAction');
        qsoState.stopTimer();
        resetQsoStateDefaults(qsoState);
        console.log('After:', qsoState);
    }
</script>

<div class="flex w-[230px] justify-end gap-x-3">
    <button
            type="button"
            class="disabled:bg-gray-400 disabled:cursor-not-allowed h-9 cursor-pointer rounded-md bg-indigo-600 p-2.5 py-1.5 text-base font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            title="Ctrl-s">Log Contact
    </button>
    <button
            onclick={resetAction}
            type="button"
            class="h-9 w-[74px] cursor-pointer rounded-md bg-white px-2.5 py-1.5 text-base font-semibold ring-1 shadow-sm ring-gray-300 ring-inset hover:bg-gray-100"
            title="ESC">Clear
    </button>
</div>
