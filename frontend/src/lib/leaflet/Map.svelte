<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import { worldPositionUnderCursor } from '$lib/stores';
	import * as L from 'leaflet';

	export let zoom: number;
	export let crs: L.CRS;

	const dispatch = createEventDispatcher();

	function updateCoordinates(e: L.LeafletMouseEvent) {
		let x = e.latlng.lat;
		let z = e.latlng.lng;

		worldPositionUnderCursor.set({
			x: Math.round(x),
			z: Math.round(z)
		});
	}

	let map: L.Map;

	$: if (map) {
		map.setZoom(zoom);
	}

	function createMap(node: Node) {
		console.log(crs);
		map = L.map('map', {
			attributionControl: false,
			zoomControl: false,
			crs: crs
		}).setView([735, -818], 0);

		map.on('mouseover', updateCoordinates);
		map.on('mousemove', updateCoordinates);
		map.on('click', updateCoordinates);

		L.tileLayer('/tiles/flat/{z}/{x}/{y}.png', {
			maxZoom: 0,
			minZoom: -8,
			tileSize: 256,
			noWrap: true
		}).addTo(map);

		return {
			destroy() {
				if (map) {
					map.remove();
				}
			}
		};
	}
</script>

<div id="map" use:createMap />

<style>
	#map {
		background-color: black;
		height: 100vh;
	}
</style>
