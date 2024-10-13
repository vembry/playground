import { v7 as uuid7 } from 'uuid';
import Cookie from 'js-cookie';

const USER_ID: string = 'user_id';
export function getUserId(): string {
	const userId = Cookie.get(USER_ID) || uuid7();
	setUserId(userId);
	return userId;
}
export function setUserId(userId: string) {
	Cookie.set(USER_ID, userId);
}
