let currentToken: string | null = null;

export function setCurrentToken(token: string | null): void {
	currentToken = token;
}

export function getCurrentToken(): string | null {
	return currentToken;
}
