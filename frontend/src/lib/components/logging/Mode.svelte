<script lang="ts">
    interface Props {
        id: string;
        label: string;
        value: string;
        list: { key: string; value: string }[];
        labelCss?: string;
        divCss?: string;
        inputCss?: string;
        overallWidthCss?: string;
    }
    let {
        id,
        label,
        value = $bindable(),
        list,
        labelCss = 'block text-sm/5 font-medium',
        divCss = 'mt-2 grid grid-cols-1',
        inputCss = 'col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-1.5 pr-8 pl-3 text-base outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600',
        overallWidthCss = 'w-[100px]'
    }: Props = $props();
</script>

<div class={overallWidthCss}>
    <label for={id} class={labelCss}>{label}</label>
    <div class={divCss}>
        <select
                bind:value={value}
                id={id}
                class={inputCss}
        >
            {#each list as mode (mode.key)}
                <option value={mode.key}>{mode.value}</option>
            {/each}
        </select>
        <svg class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end text-gray-500 sm:size-4" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" data-slot="icon">
            <path fill-rule="evenodd" d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
        </svg>
    </div>
</div>
