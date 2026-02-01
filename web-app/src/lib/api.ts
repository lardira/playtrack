import type { Player, LeaderboardPlayer, PlayedGame, Game, AuthResponse } from './types';
import { playersMock, leaderboardMock, gamesMock, gamesPlayedMock } from './mocks.js';
import { browser } from '$app/environment';
import { getTokenFromCookie } from './cookies';
import { getCurrentToken } from './tokenHolder';

const USE_MOCKS = false;

const baseURL = 'http://localhost:5000';

export async function api<T>(url: string, options: RequestInit = {}): Promise<T> {
    if (USE_MOCKS) {
        await new Promise(r => setTimeout(r, 100)); // имитация задержки
        switch (url) {
            case '/players/': return playersMock as any;
            case '/players/leaderboard': return leaderboardMock as any;
            default: return {} as any;
        }
    } else {
        const token = browser ? (getTokenFromCookie() ?? getCurrentToken()) : null;
        const headers: Record<string, string> = { 'Content-Type': 'application/json', 'mode': 'cors', 'credentials': 'include' };
        if (token) headers['Authorization'] = `Bearer ${token}`;

        try {
            const res = await fetch(baseURL + url, { ...options, headers });

            // Проверяем, что ответ получен (не CORS ошибка)
            if (!res.ok) {
                let errorMessage = `HTTP ${res.status}: ${res.statusText}`;
                try {
                    const errorData = await res.json();
                    // Huma возвращает ошибки в формате { message: "..." } или { body: { message: "..." } }
                    if (errorData.message) {
                        errorMessage = errorData.message;
                    } else if (errorData.body?.message) {
                        errorMessage = errorData.body.message;
                    } else if (typeof errorData === 'string') {
                        errorMessage = errorData;
                    }
                } catch {
                    const errorText = await res.text();
                    if (errorText) errorMessage = errorText;
                }
                throw new Error(errorMessage);
            }
            return res.json();
        } catch (error: any) {
            // Если это ошибка сети (CORS, connection refused и т.д.)
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                throw new Error(
                    `Ошибка: ${error.message}`
                );
            }
            throw error;
        }
    }
}

// Auth API
export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    token: string;
}

export interface RegisterRequest {
    username: string;
    password: string;
    img?: string | null;
    email?: string | null;
}

export interface RegisterResponse {
    id: string;
}

export interface SetPasswordRequest {
    password: string;
}

export interface SetPasswordResponse {
    id: string;
}

export interface PlayerResponse {
    item: Player;
}

// Backend returns { Body: { token } } or { body: { token } } depending on serialization
function unwrapBody<T>(res: { Body?: T; body?: T }): T | undefined {
	return res.Body ?? res.body;
}

// Бэкенд возвращает { "$schema": "...", "token": "..." }
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api<{ token?: string; Token?: string }>('/pub/auth/login', {
        method: 'POST',
        body: JSON.stringify(data)
    });
    const token = response.token ?? response.Token;
    if (!token || typeof token !== 'string') throw new Error('No token in response');
    return { token };
};

export const register = async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await api<{ Body?: { id: string }; body?: { id: string }; id?: string }>('/pub/auth/register', {
        method: 'POST',
        body: JSON.stringify(data)
    });
    const id = unwrapBody(response)?.id ?? response.id;
    if (!id) throw new Error('No id in response');
    return { id };
};

export const setPassword = async (data: SetPasswordRequest): Promise<SetPasswordResponse> => {
    const response = await api<{ Body?: { id: string }; body?: { id: string }; id?: string }>('/pub/set-password', {
        method: 'PATCH',
        body: JSON.stringify(data)
    });
    const id = unwrapBody(response)?.id ?? response.id;
    if (!id) throw new Error('No id in response');
    return { id };
};

// Player API — бэкенд может вернуть { Body: { item/items } }, { body: { ... } } или { item/items } на верхнем уровне
function getItems<T>(r: { Body?: { items: T[] }; body?: { items: T[] }; items?: T[] }): T[] {
	const wrap = r.Body ?? r.body;
	return wrap?.items ?? r.items ?? [];
}
function getItem<T>(r: { Body?: { item: T }; body?: { item: T }; item?: T }): T {
	const wrap = r.Body ?? r.body;
	const item = wrap?.item ?? (r as { item?: T }).item;
	if (item == null) throw new Error('No item in response');
	return item;
}

export const getPlayers = () =>
	api<{ Body?: { items: Player[] }; body?: { items: Player[] }; items?: Player[] }>('/v1/players/').then(getItems);
export const getPlayer = (id: string) =>
	api<{ Body?: { item: Player }; body?: { item: Player }; item?: Player }>(`/v1/players/${id}`).then(getItem);
export const getPlayerPlayedGames = (playerId: string) =>
	api<{ Body?: { items: PlayedGame[] }; body?: { items: PlayedGame[] }; items?: PlayedGame[] }>(
		`/v1/players/${playerId}/played-games`
	).then(getItems);

export interface UpdatePlayerRequest {
	username?: string;
	img?: string | null;
	email?: string | null;
}
export const updatePlayer = async (playerId: string, data: UpdatePlayerRequest): Promise<{ id: string }> => {
	const response = await api<{ Body?: { id: string }; body?: { id: string } }>(`/v1/players/${playerId}`, {
		method: 'PATCH',
		body: JSON.stringify(data)
	});
	const id = (response as any).Body?.id ?? (response as any).body?.id;
	if (!id) throw new Error('No id in response');
	return { id };
};

export const getLeaderboard = () => api<LeaderboardPlayer[]>('/v1/players/leaderboard');
