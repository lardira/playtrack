import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { Player } from '../lib/types';

const getStoredUser = (): Player | null => {
    if (!browser) return null;
    const stored = localStorage.getItem('user');
    return stored ? JSON.parse(stored) : null;
};

export const user = writable<Player | null>(getStoredUser());

if (browser) {
    user.subscribe((value) => {
        if (value) {
            localStorage.setItem('user', JSON.stringify(value));
        } else {
            localStorage.removeItem('user');
            localStorage.removeItem('token');
        }
    });
}
