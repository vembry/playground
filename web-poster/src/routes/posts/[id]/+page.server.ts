import prisma from '$lib/prisma';

export const load = async ({ params }) => {
	const post = await prisma.posts.findFirst({ where: { id: params.id } });
	return { postId: params.id, post };
};
