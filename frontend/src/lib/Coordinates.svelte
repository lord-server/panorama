<script lang="ts">
	import { worldPositionUnderCursor } from '$lib/stores';

	export let x: number;
	export let z: number;
</script>

<div class="bg-white flex border shadow-lg rounded-lg">
	<div class="px-2 py-1">
		X: <span class="font-semibold">{$worldPositionUnderCursor.x ?? x}</span>
	</div>
	<div class="pr-2 py-1">
		Z: <span class="font-semibold">{$worldPositionUnderCursor.z ?? z}</span>
	</div>
</div>
