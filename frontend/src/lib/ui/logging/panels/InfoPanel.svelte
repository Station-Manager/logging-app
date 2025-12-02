<script lang="ts">
    import InfoPanelHeader from "$lib/ui/logging/panels/InfoPanelHeader.svelte";
    import {DETAILS_TAB_TITLE, STATION_TAB_TITLE, WORKED_TAB_TITLE} from "$lib/ui/logging/panels/constants";
    import WorkedPanel from "$lib/ui/logging/panels/WorkedPanel.svelte";
    import DetailsPanel from "$lib/ui/logging/panels/DetailsPanel.svelte";
    import StationPanel from "$lib/ui/logging/panels/StationPanel.svelte";
    import SessionPanel from "$lib/ui/logging/panels/SessionPanel.svelte";

    let selected = $state(WORKED_TAB_TITLE);

</script>

<div class="flex flex-col h-[370px]">
    <InfoPanelHeader bind:activePanel={selected}/>
</div>
{#if selected === WORKED_TAB_TITLE}
    <WorkedPanel/>
{:else if selected === DETAILS_TAB_TITLE}
    <DetailsPanel/>
{:else if selected === STATION_TAB_TITLE}
    <StationPanel/>
{:else}
    <SessionPanel/>
{/if}