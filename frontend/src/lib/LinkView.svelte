<script>
	import LinkIcon from '$lib/icons/LinkIcon.svelte';

	function linkView() {
		alert();
	}
</script>

<button class="px-1 py-1 bg-white flex border shadow-lg rounded-lg" on:click={linkView}>
	<LinkIcon class="w-6 h-6" />
</button>
