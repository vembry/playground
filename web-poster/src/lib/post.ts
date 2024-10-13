import { v7 as uuid7 } from 'uuid';

export class Postx {
	id: string;
	sender: string;
	content: string;
	likeCount: number;
	dislikeCount: number;
	threads: Postx[];

	constructor(sender: string, content: string) {
		this.id = uuid7();
		this.sender = sender;
		this.content = content;
		this.threads = [];

		this.likeCount = 0;
		this.dislikeCount = 0;
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
}
