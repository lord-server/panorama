<script lang="ts">
	import type { WorldCursorPosition } from '$lib/stores';

	export let position: WorldCursorPosition | null;
</script>

<div class="bg-white flex border shadow-lg rounded-lg">
	{#if position}
		<div class="px-2 py-1">
			X: <span class="font-semibold font-mono">{position?.x}</span>
		</div>
		<div class="pr-2 py-1">
			Z: <span class="font-semibold font-mono">{position?.z}</span>
		</div>
	{:else}
		<div class="px-2 py-1">Coordinates will be shown here</div>
	{/if}
</div>
