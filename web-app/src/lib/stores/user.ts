import { writable } from 'svelte/store';
import type { Player } from '../types';

export const user = writable<Player | null>(null);

// Загрузить пользователя из localStorage при инициализации
if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('user');
    if (stored) {
        try {
            user.set(JSON.parse(stored));
        } catch (e) {
            console.error('Failed to parse stored user', e);
        }
    }
}

// Сохранять пользователя в localStorage при изменении
user.subscribe((value) => {
    if (typeof window !== 'undefined') {
        if (value) {
            localStorage.setItem('user', JSON.stringify(value));
        } else {
            localStorage.removeItem('user');
            localStorage.removeItem('token');
        }
    }
});
