import type { Player, GamePlayed, Game, LeaderboardRow, AuthResponse } from './types';
import * as mocks from './mocks';

const API = 'http://localhost:8080';
const USE_MOCKS = true; // Установите false когда backend будет готов

// Проверка доступности backend
async function isBackendAvailable(): Promise<boolean> {
    if (USE_MOCKS) return false;
    try {
        const res = await fetch(API + '/health', { method: 'GET', signal: AbortSignal.timeout(1000) });
        return res.ok;
    } catch {
        return false;
    }
}

export async function api<T>(url: string, options?: RequestInit): Promise<T> {
    const useMocks = USE_MOCKS || !(await isBackendAvailable());
    
    if (useMocks) {
        // Используем моки
        return mockApi<T>(url, options);
    }

    // Реальный API
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...options?.headers,
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(API + url, {
        credentials: 'include',
        ...options,
        headers,
    });

    if (!res.ok) {
        const error = await res.text();
        throw new Error(error || `HTTP error! status: ${res.status}`);
    }

    return res.json();
}

// Мок API
async function mockApi<T>(url: string, options?: RequestInit): Promise<T> {
    // Имитируем задержку сети
    await new Promise(resolve => setTimeout(resolve, 100));

    const method = options?.method || 'GET';

    // Auth endpoints
    if (url === '/auth/login' && method === 'POST') {
        const body = JSON.parse(options?.body as string);
        return mocks.mockLogin(body.username, body.password) as T;
    }

    if (url === '/auth/registrer' && method === 'POST') {
        const body = JSON.parse(options?.body as string);
        return mocks.mockRegister(body.username, body.password) as T;
    }

    // Player endpoints
    if (url === '/players/') {
        return mocks.getMockPlayers() as T;
    }

    if (url.startsWith('/players/') && url.endsWith('/games-played')) {
        const id = url.split('/players/')[1].split('/games-played')[0];
        return mocks.getMockPlayerGames(id) as T;
    }

    if (url === '/players/leaderboard') {
        return mocks.getMockLeaderboard() as T;
    }

    if (url.startsWith('/players/') && !url.includes('/games-played')) {
        const id = url.split('/players/')[1];
        const player = mocks.getMockPlayer(id);
        if (!player) throw new Error('Player not found');
        return player as T;
    }

    // Game endpoints
    if (url === '/games/') {
        return mocks.getMockGames() as T;
    }

    if (url.startsWith('/games/')) {
        const id = url.split('/games/')[1];
        const game = mocks.getMockGame(id);
        if (!game) throw new Error('Game not found');
        return game as T;
    }

    throw new Error(`Mock endpoint not found: ${url}`);
}


// Auth endpoints
export async function login(username: string, password: string): Promise<AuthResponse> {
    const useMocks = USE_MOCKS || !(await isBackendAvailable());
    
    let response: AuthResponse;
    if (useMocks) {
        response = mocks.mockLogin(username, password);
    } else {
        response = await api<AuthResponse>('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }
    
    if (response.token && typeof window !== 'undefined') {
        localStorage.setItem('token', response.token);
    }
    return response;
}


export async function register(username: string, password: string): Promise<AuthResponse> {
    const useMocks = USE_MOCKS || !(await isBackendAvailable());
    
    let response: AuthResponse;
    if (useMocks) {
        response = mocks.mockRegister(username, password);
    } else {
        response = await api<AuthResponse>('/auth/registrer', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }
    
    if (response.token && typeof window !== 'undefined') {
        localStorage.setItem('token', response.token);
    }
    return response;
}


// Player endpoints
export async function getPlayers(): Promise<Player[]> {
    return api<Player[]>('/players/');
}


export async function getPlayer(id: string): Promise<Player> {
    return api<Player>(`/players/${id}`);
}


export async function updatePlayer(id: string, data: Partial<Player>): Promise<Player> {
    return api<Player>(`/players/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
    });
}


export async function getLeaderboard(): Promise<LeaderboardRow[]> {
    return api<LeaderboardRow[]>('/players/leaderboard');
}


export async function getPlayerGames(playerId: string): Promise<GamePlayed[]> {
    return api<GamePlayed[]>(`/players/${playerId}/games-played`);
}


export async function updateGamePlayed(playerId: string, gameId: string, data: Partial<GamePlayed>): Promise<GamePlayed> {
    return api<GamePlayed>(`/players/${playerId}/games-played/${gameId}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
    });
}


// Game endpoints
export async function getGames(): Promise<Game[]> {
    return api<Game[]>('/games/');
}


export async function getGame(id: string): Promise<Game> {
    return api<Game>(`/games/${id}`);
}


export async function createGame(data: Partial<Game>): Promise<Game> {
    return api<Game>('/games/', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}


export async function updateGame(id: string, data: Partial<Game>): Promise<Game> {
    return api<Game>(`/games/${id}`, {
        method: 'PATCH',
        body: JSON.stringify(data),
    });
}