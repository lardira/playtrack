import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { Player } from '../lib/types';
import { getPlayer } from '../lib/api';
import { getTokenFromCookie, setTokenCookie, clearTokenCookie } from '../lib/cookies';
import { setCurrentToken } from '../lib/tokenHolder';

const getStoredUser = (): Player | null => {
    if (!browser) return null;
    const stored = localStorage.getItem('user');
    return stored ? JSON.parse(stored) : null;
};

const getStoredToken = (): string | null => {
    if (!browser) return null;
    return getTokenFromCookie();
};

export const user = writable<Player | null>(getStoredUser());
export const token = writable<string | null>(getStoredToken());

function getPlayerIdFromToken(token: string): string | null {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) return null;
        const payload = JSON.parse(atob(parts[1]));
        return payload.sub || null;
    } catch {
        return null;
    }
}

export async function loadUserFromToken(): Promise<void> {
    const currentToken = getStoredToken();
    if (!currentToken) {
        user.set(null);
        return;
    }

    const playerId = getPlayerIdFromToken(currentToken);
    if (!playerId) {
        user.set(null);
        token.set(null);
        if (browser) {
            clearTokenCookie();
            localStorage.removeItem('user');
        }
        return;
    }

    try {
        const player = await getPlayer(playerId);
        user.set(player);
        if (browser) {
            localStorage.setItem('user', JSON.stringify(player));
        }
    } catch (error) {
        console.error('Failed to load user:', error);
        user.set(null);
        token.set(null);
        if (browser) {
            clearTokenCookie();
            localStorage.removeItem('user');
        }
    }
}

if (browser) {
    const currentToken = getStoredToken();
    setCurrentToken(currentToken);
    if (currentToken) {
        loadUserFromToken();
    }

    token.subscribe((value) => {
        setCurrentToken(value);
        if (value) {
            setTokenCookie(value);
        } else {
            clearTokenCookie();
        }
    });

    user.subscribe((value) => {
        if (value) {
            localStorage.setItem('user', JSON.stringify(value));
        } else {
            localStorage.removeItem('user');
        }
    });
}
