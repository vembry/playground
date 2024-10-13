export class Postx {
	id: string;
	userId: string;
	content: string;
	likeCount: number;
	dislikeCount: number;
	threads: Postx[];

	constructor() {
		this.id = '';
		this.userId = '';
		this.content = '';
		this.likeCount = 0;
		this.dislikeCount = 0;
		this.threads = [];
	}

	constructBasic(userId: string, content: string): Postx {
		this.userId = userId;
		this.content = content;
		return this;
	}

	constructorFromPrisma(post: any): Postx {
		this.id = post.id;
		this.userId = post.userId;
		this.content = post.content;
		this.likeCount = post.likeCount;
		this.dislikeCount = post.dislikeCount;
		this.threads = post.threads || [];
		return this;
	}

	addPost(post: Postx) {
		this.threads = [...this.threads, post];
	}

	getThreads(): Postx[] {
		return this.threads;
	}

	like() {
		this.likeCount = this.likeCount + 1;
	}

	dislike() {
		this.dislikeCount = this.dislikeCount + 1;
	}

	setId(id: string) {
		this.id = id;
	}
}
