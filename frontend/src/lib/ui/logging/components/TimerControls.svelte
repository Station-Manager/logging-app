<script lang="ts">

</script>

<div class="flex gap-x-1 w-[54px]">
    <button
            aria-label="stop"
            class="cursor-pointer disabled:text-gray-400 disabled:cursor-default enabled:text-red-500"
            title="F5: Stop the QSO Timer">
        <svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M5.25 7.5A2.25 2.25 0 0 1 7.5 5.25h9a2.25 2.25 0 0 1 2.25 2.25v9a2.25 2.25 0 0 1-2.25 2.25h-9a2.25 2.25 0 0 1-2.25-2.25v-9Z"/>
        </svg>
    </button>
    <button
            aria-label="start"
            class="mr-8 cursor-pointer disabled:text-gray-400 disabled:cursor-default enabled:text-green-700"
            title="F5: Start the QSO Timer">
        <svg fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z"/>
        </svg>
    </button>
</div>
