<script lang="ts">
	import { getUserId, setUserId } from '$lib/user';
	import { onMount } from 'svelte';
	let userId: string = '';
	let isEditName: boolean = false;
	onMount(() => {
		userId = getUserId();
	});
</script>

<div class="navbar">
	<div>
		hi <strong>{userId}</strong>!
		{#if !isEditName}
			<button
				on:click={(e) => {
					isEditName = !isEditName;
				}}>edit name</button
			>
		{/if}
	</div>
	{#if isEditName}
		<div>
			<input type="text" bind:value={userId} />
			<button
				on:click={(e) => {
					setUserId(userId);
					isEditName = !isEditName;
					userId = userId;
				}}>save</button
			>
		</div>
	{/if}
	<div>
		<ul>
			<a href="/">
				<li>home</li>
			</a>
			<a href="/posts">
				<li>posts</li>
			</a>
		</ul>
	</div>
</div>

<style>
	.navbar ul {
		list-style-type: none;
		padding-left: 0px;
	}
	.navbar li {
		display: inline;
		padding-right: 10px;
	}
</style>
