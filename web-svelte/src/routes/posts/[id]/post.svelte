<script lang="ts">
	import { Postx } from '$lib/post';
	import Post from './post.svelte';
	import { enhance } from '$app/forms';

	export let post: Postx;
	let postContent: string = '';
	let showThread = true;
</script>

<div class="post">
	<div class="box">
		<div>
			<strong>{post.userId}</strong> - {post.createdAt}
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
				hidden={post.threads.length == 0}
				on:click={(e) => {
					post.toggleThread();
					post = post;
				}}
			>
				{#if post.showThread}
					<span>hide</span>
				{/if}
				{#if !post.showThread}
					<span>show</span>
				{/if}
				thread
			</button>
		</div>
		<div>
			<form
				method="POST"
				action="?/addPostThread"
				use:enhance={() => {
					return async ({ result }) => {
						console.log(result);
						if (result.status == 200) {
							alert('reply sent!');

							// add post to the threads on UI
							const newPost = new Postx().constructorFromPojo(result.data);
							post.addPost(newPost);
							postContent = '';

							// trigger svelte reactivity
							post = post;
						} else {
							alert('failed to send reply');
						}
					};
				}}
			>
				<div>
					<input type="hidden" name="parentPostId" value={post.id} />
					<textarea placeholder="Write your post..." name="content" bind:value={postContent}
					></textarea>
				</div>
				<div>
					<button type="submit">send</button>
				</div>
			</form>
		</div>
	</div>
	<div>
		{#if post.showThread}
			{#each post.threads as thread}
				<div>
					<Post post={thread}></Post>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.post {
		padding-left: 10px;
		padding-top: 10px;
		border-left: 1px solid;
		border-top: 1px solid;
		margin-bottom: 10px;
	}
	.post > .box {
		margin-bottom: 10px;
	}
	.post > .box > div {
		margin-bottom: 5px;
	}
</style>
