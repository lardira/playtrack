/**
 * Хранит токен в памяти как запасной источник (api не может импортировать stores/user из‑за циклических зависимостей).
 * Основное хранилище — кука playtrack_token.
 */
let currentToken: string | null = null;

export function setCurrentToken(token: string | null): void {
	currentToken = token;
}

export function getCurrentToken(): string | null {
	return currentToken;
}
