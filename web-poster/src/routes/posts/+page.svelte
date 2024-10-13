<script lang="ts">
	import { onMount } from 'svelte';
	import { Postx } from '$lib/post';
	import { getUserId } from '$lib/user';
	import Nav from '../../components/nav.svelte';
	import { goto } from '$app/navigation';

	onMount(function () {
		getUserId();
	});

	let posts: Postx[] = [];

	function addPost() {
		posts = [new Postx(`bot-${Date.now()}`, `this is a post! created at ${new Date()}`), ...posts];
	}
</script>

<Nav></Nav>
<div>
	<div class="mb-10">
		<button on:click={(e) => addPost()}>add dummy post</button>
	</div>
	<div class="mb-10 post-list">
		<!-- put posts here -->
		{#each posts as post}
			<div class="post mb-10">
				<div>
					{post.sender}
				</div>
				<div>
					{post.content}
				</div>
				<div>
					<button on:click={(e) => {post.like()}}>like | {post.likeCount}</button>
					<button on:click={(e) => {post.dislike()}}>dislike | {post.dislikeCount}</button>
					<button on:click={(e) => {goto(`/posts/${post.id}`)}}>reply</button>
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.mb-10 {
		margin-bottom: 10px;
	}

	.post-list a {
		text-decoration: none;
	}

	.post {
		padding: 10px;
		border: 1px solid;
	}
	.post > div {
		margin-bottom: 5px;
	}
</style>
