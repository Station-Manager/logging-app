<script lang="ts">
    import VfoBox from "$lib/ui/logging/components/VfoBox.svelte";
    import VfoInput from "$lib/ui/logging/components/VfoInput.svelte";
</script>

<div class="flex flex-col w-[250px] h-[80px] mt-6 gap-y-2">
    <div class="flex flex-row items-center gap-x-2">
        <VfoBox label='VFO-A'/>
        <VfoInput id='vfo-a' value="14.320.000"/>
    </div>
    <div class="flex flex-row items-center gap-x-2">
        <VfoBox label='VFO-B' bgColorCss='bg-gray-500/80'/>
        <VfoInput id='vfo-b' value="14.320"/>
    </div>
</div>
