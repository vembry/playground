<script lang="ts">
	import { Postx } from '$lib/post';
	import { getUserId } from '../lib/user'
	import Post from '../components/post.svelte';

	export let post: Postx;
	let newPost: string = '';

	function handleKeydown(event: KeyboardEvent) {
		// Check if "Enter" key is pressed
		if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault(); // Prevent new line
			addPost(); // Submit the form
		}
	}

	function addPost() {
		if (newPost.trim()) {
			const newId = Date.now();

			post.threads = [...post.threads, new Postx(getUserId(), newPost)];

			newPost = '';
		}
	}
</script>

<div class="post">
	<div class="box">
		<div>
			<strong>{post.sender}</strong>
		</div>
		<div>
			{post.content}
		</div>
        <div>
            <button on:click={(e) => {post.like(); post = post}}>like | {post.likeCount}</button><button on:click={(e) => {post.dislike(); post = post}}>dislike | {post.dislikeCount}</button>
        </div>
		<div>
			<form on:submit|preventDefault={addPost}>
				<div>
					<textarea placeholder="Write your post..." bind:value={newPost} on:keydown={handleKeydown}
					></textarea>
				</div>
				<div>
					<button type="submit">send</button>
				</div>
			</form>
		</div>
	</div>
	<div>
		{#each post.getThreads() as thread}
			<div>
				<Post post={thread}></Post>
			</div>
		{/each}
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
