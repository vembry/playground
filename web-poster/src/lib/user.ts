import { v7 as uuid7 } from 'uuid';

const USER_ID: string = 'user_id';
export function getUserId(): string {
	if (typeof window !== 'undefined') {
		const userId = sessionStorage.getItem(USER_ID) || uuid7();
		sessionStorage.setItem(USER_ID, uuid7());
		return userId;
	}
	return ""
}
