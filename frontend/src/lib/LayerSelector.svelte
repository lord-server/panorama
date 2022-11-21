<script lang="ts">
	let selected = 'isometric';

	let modes = [
		{ id: 'isometric', name: '3D' },
		{ id: 'flat', name: '2D' }
	];
</script>

<div class="bg-white flex rounded-lg shadow-lg divide-x">
	{#each modes as mode, id}
		<button
			class="px-2 py-1 font-semibold"
			class:rounded-l-lg={id === 0}
			class:rounded-r-lg={id === modes.length - 1}
			class:bg-red-600={selected === mode.id}
			class:text-white={selected === mode.id}
			on:click={() => {
				selected = mode.id;
			}}
			type="button"
		>
			{mode.name}
		</button>
	{/each}
</div>
