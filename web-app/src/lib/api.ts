import type { Player, LeaderboardPlayer, PlayedGame, Game, AuthResponse } from './types';
import { playersMock, leaderboardMock, gamesMock, gamesPlayedMock } from './mocks.js';
import { browser } from '$app/environment';

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
        const token = browser ? localStorage.getItem('token') : null;
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

// Auth functions
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api<{ body: { token: string } }>('/pub/auth/login', {
        method: 'POST',
        body: JSON.stringify({ body: data })
    });
    return { token: response.body.token };
};

export const register = async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await api<{ body: { id: string } }>('/pub/auth/register', {
        method: 'POST',
        body: JSON.stringify({ body: data })
    });
    return { id: response.body.id };
};

export const setPassword = async (data: SetPasswordRequest): Promise<SetPasswordResponse> => {
    const response = await api<{ body: { id: string } }>('/pub/set-password', {
        method: 'PATCH',
        body: JSON.stringify({ body: data })
    });
    return { id: response.body.id };
};

// Player API
export const getPlayers = () => api<{ body: { items: Player[] } }>('/v1/players/').then(r => r.body.items);
export const getPlayer = (id: string) => api<{ body: { item: Player } }>(`/v1/players/${id}`).then(r => r.body.item);
export const getLeaderboard = () => api<LeaderboardPlayer[]>('/v1/players/leaderboard');
