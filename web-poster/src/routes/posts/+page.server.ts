import { v7 as uuid7 } from 'uuid';

import prisma from '$lib/prisma';
import { fail } from '@sveltejs/kit';
import { getUserIdFromCookies } from '$lib/user';

export const load = async () => {
	const posts = await prisma.posts.findMany();
	return { posts };
};

export const actions = {
	addPost: async ({ cookies, request }) => {
		const data = await request.formData();
		const content = data.get('content') as string;
		const userId = getUserIdFromCookies(cookies);

		if (!content.trim()) {
			return fail(400, { content, missing: true });
		}

		return await prisma.posts
			.create({
				data: {
					id: uuid7(),
					parentPostId: null,
					userId: userId,
					content: content,
					likeCount: 0,
					dislikeCount: 0
				}
			})
			.then((res) => {
				return res;
			})
			.catch((e) => {
				console.log('error');
				console.log(e);
				return {};
			});
	}
};
