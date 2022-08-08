<script lang="ts">
	import '../app.css';
	import Map from '$lib/leaflet/Map.svelte';
	import ZoomControl from '$lib/ZoomControl.svelte';
	import Coordinates from '$lib/Coordinates.svelte';
	import LayerSelector from '$lib/LayerSelector.svelte';
	import LinkView from '$lib/LinkView.svelte';
</script>

<svelte:head>
	<title>Panorama</title>
</svelte:head>

<div class="relative">
	<Map />
	<div class="absolute top-4 right-4">
		<ZoomControl />
	</div>

	<div class="absolute top-4 left-4">
		<div class="flex space-x-4">
			<Coordinates x={1234} z={567} />
			<LayerSelector />
			<LinkView />
		</div>
	</div>
</div>
