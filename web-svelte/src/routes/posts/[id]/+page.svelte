<script lang="ts">
	import { Postx, type PostPojo } from '$lib/post';
	import Nav from '../../../components/nav.svelte';
	import Post from './post.svelte';

	export let data: {
		postId: string;
		post: PostPojo;
	};

	const post = new Postx().constructorFromPojo(data.post);
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
