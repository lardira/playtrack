import { api } from '$lib/api';
import type { LeaderboardRow, Player, GamePlayed, Game } from '$lib/types';


export const load = async () => {
    const leaderboard = await api<LeaderboardRow[]>('/players/leaderboard');
    const players = await api<Player[]>('/players');

    // Получаем текущую игру для каждого игрока
    const leaderboardWithGames = await Promise.all(
        leaderboard.map(async (row) => {
            try {
                const games = await api<GamePlayed[]>(`/players/${row.Player_id}/games-played`);
                // Находим последнюю игру со статусом "В процессе"
                const currentGame = games
                    .filter((g) => g.status === 'В процессе')
                    .sort((a, b) =>
                        new Date(b.start_date).getTime() - new Date(a.start_date).getTime()
                    )[0];

                if (currentGame) {
                    const game = await api<Game>(`/games/${currentGame.game_id}`);
                    return { ...row, current_game: game.Title };
                }
                return row;
            } catch (e) {
                console.error(`Failed to load games for player ${row.Player_id}`, e);
                return row;
            }
        })
    );

    return { leaderboard: leaderboardWithGames, players };
};