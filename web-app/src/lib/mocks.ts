import type { Player, LeaderboardRow, Game, GamePlayed } from './types';

export const playersMock: Player[] = [
    { ID: '1', Username: 'PlayerOne', Description: 'Любитель RPG', Img: '/favicon.png' },
    { ID: '2', Username: 'PlayerTwo', Description: '', Img: '/favicon.png' },
];

export const leaderboardMock: LeaderboardRow[] = [
    { Player_id: '1', total: 1200, completed: 15, drop: 2, reroll: 1 },
    { Player_id: '2', total: 900, completed: 10, drop: 1, reroll: 0 },
];

export const gamesMock: Game[] = [
    { ID: '101', Points: 100, HoursToBeat: 10, Title: 'The Witcher 3', URL: '', CreatedAt: '' },
];

export const gamesPlayedMock: GamePlayed[] = [
    { Player_id: '1', game_id: '101', status: 'Пройдено', scores: 100, start_date: '2026-01-25', end_date: '2026-01-26', comment: 'Отличная игра!', rating: '10/10' }
];
