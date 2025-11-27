<script lang="ts">

    import {qsoState} from "$lib/states/qso-state.svelte.js";

    const disabledCss: string = 'disabled:outline-orange-500 disabled:outline-2 disabled:bg-orange-200';

    interface Props {
        id: string,
        label: string,
        value: string,
        disabled: boolean,
    }

    let {
        id,
        label,
        value = $bindable(),
        disabled,
    }: Props = $props();

    let ticking: boolean = $derived.by(() => {
        return disabled && qsoState.timeOff !== qsoState.timeOn;
    });

</script>

<div class="w-[100px]">
    <label for={id} class="block text-sm/5 font-medium">{label}</label>
    <div class="relative mt-2">
        <input
                bind:value={value}
                type="time"
                id={id}
                disabled={disabled}
                class="{ticking ? disabledCss : ''} block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600"
            />
        <span class="absolute top-2 right-2">
            <svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-5 text-gray-700">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
        </span>
    </div>
</div>
