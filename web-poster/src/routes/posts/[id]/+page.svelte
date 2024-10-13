<script lang="ts">
	import { Postx } from '$lib/post';
	import Nav from '../../../components/nav.svelte';
	import Post from '../../../components/post.svelte';

	type pageParam = {
		postId: string;
	};

	export let data: pageParam;

	// dummy post data
	let post: Postx = new Postx(
		`bot-${Date.now()}`,
		`this is an amazing post! created at '${new Date()}'`
	);
</script>

<Nav></Nav>
<div>
	id=#{data.postId}
</div>
<div>
	<Post {post}></Post>
</div>

<style>
	div {
		margin-bottom: 10px;
	}
</style>
