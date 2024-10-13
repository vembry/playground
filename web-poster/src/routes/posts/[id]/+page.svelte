<script lang="ts">
	import { Postx } from '$lib/post';
	import Nav from '../../../components/nav.svelte';
	import Post from '../../../components/post.svelte';

	export let data: {
		postId: string;
		post: any;
	};

	const post = new Postx().constructorFromPrisma(data.post);
</script>

<Nav></Nav>
<div>
	post#{data.postId}
</div>
<div>
	<Post {post}></Post>
</div>

<style>
	div {
		margin-bottom: 10px;
	}
</style>
