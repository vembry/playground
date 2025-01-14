import { v7 as uuid7 } from 'uuid';
import prisma from '$lib/prisma';
import { getUserIdFromCookies } from '$lib/user';
import type { posts } from '@prisma/client';
import { type PostPojo } from '$lib/post';
import { fail } from '@sveltejs/kit';

export const load = async ({ params }) => {
	// grab post
	const post = await prisma.posts.findFirst({ where: { id: params.id } });
	if (!post) {
		return {};
	}

	// grab post-thread's ids based on the main post then
	// use it to get all the post-threads. We had to do
	// this because query raw return plain object without
	// getting mapped to camel case as defined on prisma object

	// retrieve post-thread's ids
	const threadWithOnlyIds = await prisma.$queryRaw<posts[]>`
		WITH RECURSIVE post_tree AS (
			SELECT p.id
			FROM posts p
			WHERE parent_post_id = ${post.id}  -- Start from the parent post ID you want
			UNION ALL
			SELECT p.id
			FROM posts p
			INNER JOIN post_tree pt ON p.parent_post_id = pt.id
		)
		SELECT * FROM post_tree;
	`;

	// do early return when theres no thread found
	if (!threadWithOnlyIds || threadWithOnlyIds.length == 0) {
		return {
			postId: params.id,
			post: {
				...post,
				threads: []
			}
		};
	}

	// construct threadId arrays
	const threadIds = [];
	for (const thread of threadWithOnlyIds) {
		threadIds.push(thread.id);
	}

	// retrieve post-threads with threadIds we had
	const threads = await prisma.posts.findMany({ where: { id: { in: threadIds } } });

	// construct post + threads
	const final = loadThread(post, threads);

	return { postId: params.id, post: final };
};

/**
 * loadThread preps pre-requisite to construct post and it's nested threads
 * @param post
 * @param threads
 * @returns
 */
function loadThread(post: posts, threads: posts[]): PostPojo {
	const parentToChildMap = new Map<string, PostPojo[]>();

	// map out parentPostId with it's children
	for (const thread of threads) {
		let arr = parentToChildMap.get(thread.parentPostId!);
		if (!arr) {
			arr = [];
		}
		arr.push({
			...thread,
			threads: []
		});
		parentToChildMap.set(thread.parentPostId!, arr);
	}

	const final: PostPojo = {
		...post,
		threads: []
	};
	const childs = parentToChildMap.get(final.id)!;
	parentToChildMap.delete(post.id);

	final.threads = loadNestedThread(childs, parentToChildMap);

	return final;
}

/**
 * loadNestedThread recursively traverse through post and it's threads
 * @param parents
 * @param parentToChildMap
 * @returns PostPojo[]
 */
function loadNestedThread(
	parents: PostPojo[],
	parentToChildMap: Map<string, PostPojo[]>
): PostPojo[] {
	for (const parent of parents) {
		const arr = parentToChildMap.get(parent.id)!;
		if (!arr) {
			continue;
		}

		parentToChildMap.delete(parent.id);
		parent.threads = loadNestedThread(arr, parentToChildMap);
	}
	return parents;
}

export const actions = {
	addPostThread: async ({ cookies, request }) => {
		const data = await request.formData();

		const content = data.get('content') as string;
		if (!content.trim()) {
			return fail(400, { content, missing: true });
		}

		const parentPostId = data.get('parentPostId') as string;
		const userId = getUserIdFromCookies(cookies);

		return prisma.posts
			.create({
				data: {
					id: uuid7(),
					parentPostId: parentPostId,
					userId: userId,
					content: content,
					likeCount: 0,
					dislikeCount: 0
				}
			})
			.then((res) => {
				return {
					...res
				};
			})
			.catch((e) => {
				console.log('error');
				console.log(e);
				return {};
			});
	}
};
