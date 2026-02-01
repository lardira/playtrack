import { browser } from '$app/environment';

const TOKEN_COOKIE_NAME = 'playtrack_token';
/** 10 дней (как на бэкенде: defaultExpiration) */
const TOKEN_MAX_AGE_SEC = 10 * 24 * 60 * 60;

export function getTokenFromCookie(): string | null {
	if (!browser || typeof document === 'undefined') return null;
	const name = TOKEN_COOKIE_NAME + '=';
	const parts = document.cookie.split(';');
	for (let i = 0; i < parts.length; i++) {
		const part = parts[i].trim();
		if (part.startsWith(name)) {
			return part.slice(name.length);
		}
	}
	return null;
}

export function setTokenCookie(token: string): boolean {
	if (!browser || typeof document === 'undefined') return false;
	if (!token) return false;
	// Кука сохраняется для origin страницы (например http://localhost:5173).
	// В DevTools смотрите: Application → Cookies → http://localhost:5173 (не localhost:5000).
	const cookieStr = `${TOKEN_COOKIE_NAME}=${token};path=/;max-age=${TOKEN_MAX_AGE_SEC};samesite=lax`;
	document.cookie = cookieStr;
	return getTokenFromCookie() === token;
}

export function clearTokenCookie(): void {
	if (!browser || typeof document === 'undefined') return;
	document.cookie = `${TOKEN_COOKIE_NAME}=; path=/; max-age=0`;
}
