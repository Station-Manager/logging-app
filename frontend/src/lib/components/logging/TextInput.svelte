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
        labelCss = 'block text-sm/5 font-medium',
        divCss = 'mt-2',
        inputCss = 'block w-full rounded-md bg-white px-3 py-1.5 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600',
        overallWidthCss = 'w-[210px]'
    }: Props = $props();
</script>

<div class={overallWidthCss}>
    <label for={id} class={labelCss}>{label}</label>
    <div class={divCss}>
        <input
                bind:value={value}
                type="text"
                id={id}
                class={inputCss}
                autocomplete="off"
                spellcheck="false"
        />
    </div>
</div>