import type { Player, LeaderboardRow, GamePlayed, Game, AuthResponse } from './types';
import { playersMock, leaderboardMock, gamesMock, gamesPlayedMock } from './mocks.js';
import { browser } from '$app/environment';

const USE_MOCKS = true;

const baseURL = 'http://localhost:8080';

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
        const headers: Record<string, string> = { 'Content-Type': 'application/json' };
        if (token) headers['Authorization'] = `Bearer ${token}`;
        const res = await fetch(baseURL + url, { ...options, headers });
        if (!res.ok) throw new Error(await res.text());
        return res.json();
    }
}

// Примеры API вызовов
export const getPlayers = () => api<Player[]>('/players/');
export const getLeaderboard = () => api<LeaderboardRow[]>('/players/leaderboard');
