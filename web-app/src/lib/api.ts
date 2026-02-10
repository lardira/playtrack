import type { Player, LeaderboardPlayer, PlayedGame, Game, AuthResponse } from './types';
import { browser } from '$app/environment';
import { getTokenFromCookie } from './cookies';
import { getCurrentToken } from './tokenHolder';
import { env } from '$env/dynamic/public';
import { parseApiError } from './errors';

const baseURL = env.PUBLIC_API_URL ?? '/api';

export async function api<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = browser ? (getTokenFromCookie() ?? getCurrentToken()) : null;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    try {
        const res = await fetch(baseURL + url, { ...options, headers, mode: 'cors', credentials: 'include' });
        const text = await res.text();

        if (!res.ok) {
            const errorMessage = parseApiError(res.status, res.statusText, text);
            throw new Error(errorMessage);
        }

        return (text ? JSON.parse(text) : {}) as T;
    } catch (error: any) {
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            throw new Error(`Ошибка сети: ${error.message}`);
        }
        throw error;
    }
}

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
    username: string;
}

export interface SetPasswordResponse {
    id: string;
}

export interface PlayerResponse {
    item: Player;
}

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
    const response = await api<{ id: string }>('/pub/auth/register', {
        method: 'POST',
        body: JSON.stringify(data)
    });
    if (!response.id) throw new Error('No id in response');
    return { id: response.id };
};

export const setPassword = async (data: SetPasswordRequest): Promise<SetPasswordResponse> => {
    const response = await api<{ id: string }>('/pub/auth/set-password', {
        method: 'PATCH',
        body: JSON.stringify(data)
    });
    if (!response.id) throw new Error('No id in response');
    return { id: response.id };
};

function getItems<T>(r: { items: T[] }): T[] {
    return r.items ?? [];
}
function getItem<T>(r: { item: T }): T {
    if (r.item == null) throw new Error('No item in response');
    return r.item;
}

export const getPlayers = () =>
    api<{ items: Player[] }>('/v1/players/').then(getItems);
export const getPlayer = (id: string) =>
    api<{ item: Player }>(`/v1/players/${id}`).then(getItem);
export const getPlayerPlayedGames = (playerId: string) =>
    api<{ items: PlayedGame[] }>(
        `/v1/players/${playerId}/played-games`
    ).then(getItems);

export const getPlayerPlayedGame = (playerId: string, playedGameId: number) =>
    api<{ item: PlayedGame }>(
        `/v1/players/${playerId}/played-games/${playedGameId}`
    ).then(getItem);

export interface UpdatePlayerRequest {
    username?: string;
    img?: string | null;
    email?: string | null;
    description?: string | null;
}
export const updatePlayer = async (playerId: string, data: UpdatePlayerRequest): Promise<{ id: string }> => {
    const response = await api<{ id: string }>(`/v1/players/${playerId}`, {
        method: 'PATCH',
        body: JSON.stringify(data)
    });
    if (!response.id) throw new Error('No id in response');
    return { id: response.id };
};

export const getLeaderboard = () => api<LeaderboardPlayer[]>('/v1/players/leaderboard');

export const getGames = () =>
    api<{ items: Game[] }>('/v1/games/').then(getItems);
export const getGame = (id: number) =>
    api<{ item: Game }>(`/v1/games/${id}`).then(getItem);

export interface CreateGameRequest {
    title: string;
    hours_to_beat: number;
    url?: string | null;
}
export const createGame = async (data: CreateGameRequest): Promise<{ id: number }> => {
    const response = await api<{ id: number }>('/v1/games/', {
        method: 'POST',
        body: JSON.stringify(data)
    });
    if (response.id == null) throw new Error('No id in response');
    return { id: response.id };
};

export const createPlayedGame = async (playerId: string, gameId: number): Promise<{ id: number }> => {
    const response = await api<{ id: number }>(
        `/v1/players/${playerId}/played-games`,
        {
            method: 'POST',
            body: JSON.stringify({ game_id: gameId })
        }
    );
    if (response.id == null) throw new Error('No id in response');
    return { id: response.id };
};

export interface UpdatePlayedGameRequest {
    points?: number;
    comment?: string | null;
    rating?: number | null;
    status?: import('./types').PlayedGameStatus;
    started_at?: string | null;
    completed_at?: string | null;
    play_time?: string | null;
}

export const updatePlayedGame = async (
    playerId: string,
    playedGameId: number,
    data: UpdatePlayedGameRequest
): Promise<{ id: number }> => {
    const response = await api<{ id: number }>(
        `/v1/players/${playerId}/played-games/${playedGameId}`,
        {
            method: 'PATCH',
            body: JSON.stringify(data)
        }
    );
    if (response.id == null) throw new Error('No id in response');
    return { id: response.id };
};
