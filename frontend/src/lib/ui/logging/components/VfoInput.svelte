<script lang="ts">
    import {frequencyToBandFromDottedMHz} from "$lib/utils/frequency";

    interface Props {
        id: string;
        value: string;
        band: string;
        disabled?: boolean;
        inputCss?: string;
    }
    let {
        id,
        value = $bindable(),
        band,
        disabled = false,
        inputCss = 'block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600',
    }: Props = $props();

</script>

<div class="w-[116px]">
    <input
            type="text"
            autocomplete="off"
            spellcheck="false"
            id={id}
            disabled={disabled}
            title="Format: ?#.###.###"
            bind:value={value}
            class={inputCss}
    />
</div>
<div class="cursor-default w-[32px] font-semibold text-base">{frequencyToBandFromDottedMHz(band)}</div>