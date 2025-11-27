<script lang="ts">
    interface Props {
        id: string;
        label: string;
        value: string;
        labelCss?: string;
        divCss?: string;
        inputCss?: string;
        overallWidthCss?: string;
    }

    let {
        id,
        label,
        value = $bindable(),
        labelCss = 'block text-sm/5 font-medium ml-2',
        divCss = 'mt-2',
        inputCss = 'h-[68px] resize-none w-full rounded-md bg-white px-2 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600',
        overallWidthCss = 'w-[280px]'
    }: Props = $props();
</script>

<div class={overallWidthCss}>
    <div class="relative">
        <label for={id} class={labelCss}>{label}</label>
    </div>
    <!-- RENDER CLIPBOARD BUTTON HERE -->
    <div class={divCss}>
        <textarea
                bind:value={value}
                id={id}
                spellcheck="false"
                class={inputCss}
                autocomplete="off"></textarea>
    </div>
</div>
