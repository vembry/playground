/**
 * class to assists UI functionality
 * TODO: can we simplify this?
 */
export class Postx {
	id: string;
	userId: string;
	content: string;
	likeCount: number;
	dislikeCount: number;
	createdAt: Date;
	threads: Postx[];

	constructor() {
		this.id = '';
		this.userId = '';
		this.content = '';
		this.likeCount = 0;
		this.dislikeCount = 0;
		this.createdAt = new Date();
		this.threads = [];
	}

	constructorFromPojo(post: PostPojo): Postx {
		this.id = post.id;
		this.userId = post.userId;
		this.content = post.content;
		this.likeCount = post.likeCount;
		this.dislikeCount = post.dislikeCount;
		this.createdAt = post.createdAt;
		this.threads = [];
		for (const thread of post.threads) {
			this.threads.push(new Postx().constructorFromPojo(thread));
		}
		return this;
	}

	addPost(post: Postx) {
		this.threads = [...this.threads, post];
	}

	like() {
		this.likeCount = this.likeCount + 1;
	}

	dislike() {
		this.dislikeCount = this.dislikeCount + 1;
	}
}

/**
 * Basic type for posts
 */
export type PostPojo = {
	id: string;
	userId: string;
	content: string;
	likeCount: number;
	dislikeCount: number;
	createdAt: Date;
	threads: PostPojo[];
};
