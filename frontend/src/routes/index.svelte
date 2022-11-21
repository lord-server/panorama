<script lang="ts">
	import '../app.css';
	import ZoomControl from '$lib/ZoomControl.svelte';
	import Coordinates from '$lib/Coordinates.svelte';
	import LayerSelector from '$lib/LayerSelector.svelte';
	import LinkView from '$lib/LinkView.svelte';
	import { onMount } from 'svelte';
	import { worldPositionUnderCursor } from '$lib/stores';

	$: zoom = 0;
	let Map: any;
	let flat: L.CRS;
	let isometric: L.CRS;

	onMount(async () => {
		let crs = (await import('$lib/leaflet/crs'));
		flat = crs.flat;
		isometric = crs.isometric;

		Map = (await import('$lib/leaflet/Map.svelte')).default;
	});
</script>

<svelte:head>
	<title>Panorama</title>
</svelte:head>

<div class="relative">
	<svelte:component this={Map} zoom={zoom} crs={isometric} />

	<div class="absolute top-4 right-4 z-[1000]">
		<ZoomControl on:zoomOut={() => {zoom--;}} on:zoomIn={() => {zoom++;}} />
	</div>

	<div class="absolute top-4 left-4 z-[1000]">
		<div class="flex space-x-4">
			<LayerSelector />
			<!--
				<LinkView />
			-->
			<Coordinates position={$worldPositionUnderCursor} />
		</div>
	</div>
</div>
