import { api } from '$lib/api';
import { get } from 'svelte/store';
import { user } from '$lib/stores/user';
import type { Player, GamePlayed, Game } from '$lib/types';


export const load = async ({ params }) => {
    const player = await api<Player>(`/players/${params.id}`);
    const games = await api<GamePlayed[]>(`/players/${params.id}/games-played`);

    // Получаем информацию об играх для отображения названий
    const gamesWithInfo = await Promise.all(
        games.map(async (game) => {
            try {
                const gameInfo = await api<Game>(`/games/${game.game_id}`);
                return { ...game, gameTitle: gameInfo.Title, gameURL: gameInfo.URL };
            } catch (e) {
                console.error(`Failed to load game ${game.game_id}`, e);
                return { ...game, gameTitle: game.game_id, gameURL: '' };
            }
        })
    );

    // Проверяем, является ли текущий пользователь владельцем страницы
    const currentUser = get(user);
    const isOwner = currentUser?.ID === params.id;

    return { player, games: gamesWithInfo, isOwner };
};