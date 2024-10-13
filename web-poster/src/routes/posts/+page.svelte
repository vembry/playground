<script lang="ts">
	import { onMount } from 'svelte';
	import { Postx } from '$lib/post';
	import { getUserId } from '$lib/user';
	import Nav from '../../components/nav.svelte';
	import { goto } from '$app/navigation';
	import { enhance } from '$app/forms';

	onMount(function () {
		getUserId();
	});

	export let data: {
		posts: Postx[];
	};
	let content = '';

	let posts = data.posts.map(convertToPostx);

	function convertToPostx(post: any): Postx {
		const postx: Postx = new Postx();
		postx.constructorFromPrisma(post);
		return postx;
	}
</script>

<Nav></Nav>
<div>
	<!-- posts form submission -->
	<div class="mb-10">
		<form
			method="POST"
			action="?/addPost"
			use:enhance={() => {
				return async ({ result }) => {
					console.log(result);
					if (result.status == 200) {
						alert('post submitted');
						content = '';

						// convert 'post' retrieved from server into 'Postx' format
						const post = convertToPostx(result.data);
						posts = [post, ...posts];
					} else {
						alert('fail to submit post');
					}
				};
			}}
		>
			<div>
				<textarea placeholder="Write your post..." name="content" bind:value={content}></textarea>
			</div>
			<div>
				<button type="submit">submit post</button>
			</div>
		</form>
	</div>

	<!-- put posts here -->
	<div class="mb-10 post-list">
		{#each posts as post}
			<div class="post mb-10">
				<div>
					<strong>{post.userId}</strong>
				</div>
				<div>
					{post.content}
				</div>
				<div>
					<button
						on:click={(e) => {
							post.like();
							post = post;
						}}>like | {post.likeCount}</button
					>
					<button
						on:click={(e) => {
							post.dislike();
							post = post;
						}}>dislike | {post.dislikeCount}</button
					>
					<button
						on:click={(e) => {
							goto(`/posts/${post.id}`);
						}}>reply</button
					>
				</div>
			</div>
		{/each}
	</div>
</div>

<style>
	.mb-10 {
		margin-bottom: 10px;
	}

	.post {
		padding: 10px;
		border: 1px solid;
	}
	.post > div {
		margin-bottom: 5px;
	}
</style>
