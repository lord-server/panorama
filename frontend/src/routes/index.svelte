<script lang="ts">
	import '../app.css';
	import ZoomControl from '$lib/ZoomControl.svelte';
	import Coordinates from '$lib/Coordinates.svelte';
	import LayerSelector from '$lib/LayerSelector.svelte';
	import LinkView from '$lib/LinkView.svelte';
	import { onMount } from 'svelte';

	let Map: any;
	onMount(async () => {
		Map = (await import('$lib/leaflet/Map.svelte')).default;
	});
</script>

<svelte:head>
	<title>Panorama</title>
</svelte:head>

<div class="relative">
	<svelte:component this={Map} />
	<div class="absolute top-4 right-4 z-[1000]">
		<ZoomControl />
	</div>

	<div class="absolute top-4 left-4 z-[1000]">
		<div class="flex space-x-4">
			<LayerSelector />
			<LinkView />
			<Coordinates x={1234} z={567} />
		</div>
	</div>
</div>
