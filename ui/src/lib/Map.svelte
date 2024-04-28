<script lang="ts">
	import { worldPositionUnderCursor } from '$lib/stores';

	import Map from 'ol/Map.js';
	import XYZ from 'ol/source/XYZ.js';
	import TileLayer from 'ol/layer/Tile.js';
	import View from 'ol/View.js';
	import TileDebug from 'ol/source/TileDebug';
	import { Projection, addCoordinateTransforms, addProjection, fromLonLat } from 'ol/proj';
	import type { TileCoord } from 'ol/tilecoord';
	import { TileGrid } from 'ol/tilegrid';

	function updateCoordinates(e: any) {
		let x = e.latlng.lat;
		let z = e.latlng.lng;

		worldPositionUnderCursor.set({
			x: Math.round(x),
			z: Math.round(z)
	 	});
	}

	function createMap(node: Node) {
		const linearPixels = new Projection({
			code: 'linear-pixels',
			units: 'pixels',
			extent: [-32768, -32768, 32768, 32768],
			axisOrientation: 'seu'
		});
		addProjection(linearPixels);

		// addCoordinateTransforms(
		// 	'flatpixels',
		// 	isometric,
		// 	function (coordinate: Coordinate) {
		// 		let x = ((coordinate[1] - coordinate[0]) * 16) / 2;
		// 		let y = ((coordinate[1] + coordinate[0]) * 16) / 4;
		// 		return [x, y];
		// 	},
		// 	function (coordinate: Coordinate) {
		// 		coordinate[1] *= 2;
		// 		let x = (coordinate[1] - coordinate[0]) / 16;
		// 		let y = (coordinate[1] + coordinate[0]) / 16;
		// 		return [x, y];
		// 	}
		// );

		const extent = [-32768*16, -32768*16, 32768*16, 32768*16];

		const tileGrid = new TileGrid({
			resolutions: [1024, 512, 256, 128, 64, 32, 16, 8, 4, 2, 1],
			extent: extent,
			origin: [0, 0],
			minZoom: 0
		});

		const map = new Map({
			target: 'map',
			layers: [
				new TileLayer({
					source: new XYZ({
						tileGrid: tileGrid,
						projection: linearPixels,
						tileUrlFunction: (coordinate: TileCoord) => {
							return `https://map.lord-server.ru/tiles/${coordinate[0] - 10}/${coordinate[1]}/${coordinate[2]}.png`;
						}
					})
				}),
				// new TileLayer({
				// 	source: new TileDebug({
				// 		tileGrid: tileGrid,
				// 		projection: flatpixels
				// 	})
				// })
			],
			view: new View({
				projection: linearPixels,
				center: [0, 0],
				// center: fromLonLat([]),
				minZoom: 0,
				maxZoom: 10,
				zoom: 4,
				resolutions: tileGrid.getResolutions(),
			})
		});

		return {
			destroy() {
				if (map) {
					map.setTarget(undefined);
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
