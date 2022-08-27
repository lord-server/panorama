<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	function zoomIn(event: Event) {
        dispatch('zoomIn');
    }

	function zoomOut(event: Event) {
        dispatch('zoomOut');
    }
</script>

<div class="flex flex-col space-y-1 bg-white rounded-lg shadow-lg divide-y">
	<button class="relative w-8 h-8 font-semibold" type="button" on:click={zoomIn}>&plus;</button>
	<button class="relative w-8 h-8 font-semibold" type="button" on:click={zoomOut}>&minus;</button>
</div>
