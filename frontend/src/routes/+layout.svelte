<script lang="ts">
	import './layout.css';
	import favicon from '../../static/favicon.png';
	
	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon}>
</svelte:head>

{@render children()}
