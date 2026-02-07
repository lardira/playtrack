import type { Player, LeaderboardPlayer, PlayedGame, Game, AuthResponse } from './types';
import { browser } from '$app/environment';
import { getTokenFromCookie } from './cookies';
import { getCurrentToken } from './tokenHolder';
import { env } from '$env/dynamic/public';

const baseURL = env.PUBLIC_API_URL ?? '/api'

export async function api<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = browser ? (getTokenFromCookie() ?? getCurrentToken()) : null;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    try {
        const res = await fetch(baseURL + url, { ...options, headers, mode: 'cors', credentials:'include' });
        const text = await res.text();

        if (!res.ok) {
            let errorMessage = `HTTP ${res.status}: ${res.statusText}`;
            try {
                const errorData = text ? JSON.parse(text) : {};
                if (errorData.message) {
                    errorMessage = errorData.message;
                } else if (errorData.body?.message) {
                    errorMessage = errorData.body.message;
                } else if (typeof errorData === 'string') {
                    errorMessage = errorData;
                }
            } catch {
                if (text) errorMessage = text;
            }
            throw new Error(errorMessage);
        }

        return (text ? JSON.parse(text) : {}) as T;
    } catch (error: any) {
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            throw new Error(
                `Ошибка: ${error.message}`
            );
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

function unwrapBody<T>(res: { Body?: T; body?: T }): T | undefined {
    return res.Body ?? res.body;
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

export const getPlayerPlayedGame = (playerId: string, playedGameId: number) =>
    api<{ Body?: { item: PlayedGame }; body?: { item: PlayedGame }; item?: PlayedGame }>(
        `/v1/players/${playerId}/played-games/${playedGameId}`
    ).then(getItem);

export interface UpdatePlayerRequest {
    username?: string;
    img?: string | null;
    email?: string | null;
    description?: string | null;
}
export const updatePlayer = async (playerId: string, data: UpdatePlayerRequest): Promise<{ id: string }> => {
    const response = await api<{ Body?: { id: string }; body?: { id: string } }>(`/v1/players/${playerId}`, {
        method: 'PATCH',
        body: JSON.stringify(data)
    });
    const id = (response as any)?.id;
    if (!id) throw new Error('No id in response');
    return { id };
};

export const getLeaderboard = () => api<LeaderboardPlayer[]>('/v1/players/leaderboard');

export const getGames = () =>
    api<{ Body?: { items: Game[] }; body?: { items: Game[] }; items?: Game[] }>('/v1/games/').then(getItems);
export const getGame = (id: number) =>
    api<{ Body?: { item: Game }; body?: { item: Game }; item?: Game }>(`/v1/games/${id}`).then(getItem);

export interface CreateGameRequest {
    title: string;
    hours_to_beat: number;
    url?: string | null;
}
export const createGame = async (data: CreateGameRequest): Promise<{ id: number }> => {
    const response = await api<{ Body?: { id: number }; body?: { id: number }; id?: number }>('/v1/games/', {
        method: 'POST',
        body: JSON.stringify(data)
    });
    const id = (response as any).Body?.id ?? (response as any).body?.id ?? (response as any).id;
    if (id == null) throw new Error('No id in response');
    return { id };
};

export const createPlayedGame = async (playerId: string, gameId: number): Promise<{ id: number }> => {
    const response = await api<{ Body?: { id: number }; body?: { id: number }; id?: number }>(
        `/v1/players/${playerId}/played-games`,
        {
            method: 'POST',
            body: JSON.stringify({ game_id: gameId })
        }
    );
    const id = (response as any)?.id;
    if (id == null) throw new Error('No id in response');
    return { id };
};

export interface UpdatePlayedGameRequest {
    points?: number;
    comment?: string | null;
    rating?: number | null;
    status?: import('./types').PlayedGameStatus;
    completed_at?: string | null;
    play_time?: string | null;
}

export const updatePlayedGame = async (
    playerId: string,
    playedGameId: number,
    data: UpdatePlayedGameRequest
): Promise<{ id: number }> => {
    const response = await api<{ Body?: { id: number }; body?: { id: number } }>(
        `/v1/players/${playerId}/played-games/${playedGameId}`,
        {
            method: 'PATCH',
            body: JSON.stringify(data)
        }
    );
    const id = (response as any).Body?.id ?? (response as any).body?.id ?? (response as any).id;
    if (id == null) throw new Error('No id in response');
    return { id };
};
