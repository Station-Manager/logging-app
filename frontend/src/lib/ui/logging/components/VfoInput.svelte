<script lang="ts">
    import {frequencyToBandFromDottedMHz} from "$lib/utils/frequency";
    import {inputBase} from "$lib/ui/styles";

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
        inputCss = inputBase,
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