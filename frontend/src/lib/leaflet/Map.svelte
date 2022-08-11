<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import 'leaflet/dist/leaflet.css';
	import { worldPositionUnderCursor } from '$lib/stores';
	import * as L from 'leaflet';

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

	function createMap(node: Node) {
		const isometric = L.extend(L.CRS.Simple, {
			projection: {
				project: (latlng: L.LatLng) => {
					let x = ((latlng.lng - latlng.lat) * 16) / 2;
					let y = ((latlng.lng + latlng.lat) * 16) / 4;
					return new L.Point(x, y);
				},

				unproject: (point: L.Point) => {
					point.y *= 2;
					let lat = (point.y - point.x) / 16;
					let lng = (point.y + point.x) / 16;
					return new L.LatLng(lat, lng);
				}
			},
			transformation: new L.Transformation(1, 0, 1, 0)
		});

		map = L.map('map', {
			attributionControl: false,
			zoomControl: false,
			crs: isometric
		}).setView([735, -818], 0);

		map.on('mouseover', updateCoordinates);
		map.on('mousemove', updateCoordinates);
		map.on('click', updateCoordinates);

		L.tileLayer('http://map.lord-server.ru/tiles/{z}/{x}/{y}.png', {
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
		background-color: lightgray;
		height: 100vh;
	}
</style>
