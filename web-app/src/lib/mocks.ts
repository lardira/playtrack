import type { Player, GamePlayed, Game, LeaderboardRow, AuthResponse } from './types';

// Моки данных
const mockPlayers: Player[] = [
    {
        ID: '1',
        Username: 'PlayerOne',
        Description: 'Опытный игрок',
        Img: 'https://api.dicebear.com/7.x/avataaars/svg?seed=PlayerOne',
        Email: 'player1@example.com',
    },
    {
        ID: '2',
        Username: 'GamerPro',
        Description: 'Профессиональный геймер',
        Img: 'https://api.dicebear.com/7.x/avataaars/svg?seed=GamerPro',
        Email: 'gamer@example.com',
    },
    {
        ID: '3',
        Username: 'GameMaster',
        Description: 'Мастер игр',
        Img: 'https://api.dicebear.com/7.x/avataaars/svg?seed=GameMaster',
        Email: 'master@example.com',
    },
];

const mockGames: Game[] = [
    {
        ID: 'game1',
        Points: 10,
        HoursToBeat: 20,
        Title: 'Dark Souls III',
        URL: 'https://store.steampowered.com/app/374320',
        CreatedAt: '2025-01-01',
    },
    {
        ID: 'game2',
        Points: 8,
        HoursToBeat: 15,
        Title: 'The Witcher 3',
        URL: 'https://store.steampowered.com/app/292030',
        CreatedAt: '2025-01-02',
    },
    {
        ID: 'game3',
        Points: 12,
        HoursToBeat: 25,
        Title: 'Elden Ring',
        URL: 'https://store.steampowered.com/app/1245620',
        CreatedAt: '2025-01-03',
    },
];

const mockGamesPlayed: Record<string, GamePlayed[]> = {
    '1': [
        {
            Player_id: '1',
            game_id: 'game1',
            status: 'В процессе',
            scores: 10,
            start_date: '2025-01-20',
            end_date: null,
            comment: 'Отличная игра!',
            rating: '9/10',
            time_played: '15 часов',
        },
        {
            Player_id: '1',
            game_id: 'game2',
            status: 'Пройдено',
            scores: 8,
            start_date: '2025-01-10',
            end_date: '2025-01-18',
            comment: 'Прошел на 100%',
            rating: '10/10',
            time_played: '50 часов',
        },
    ],
    '2': [
        {
            Player_id: '2',
            game_id: 'game3',
            status: 'В процессе',
            scores: 12,
            start_date: '2025-01-22',
            end_date: null,
            comment: 'Сложная, но интересная',
            rating: '8/10',
            time_played: '8 часов',
        },
    ],
    '3': [
        {
            Player_id: '3',
            game_id: 'game1',
            status: 'Пройдено',
            scores: 10,
            start_date: '2025-01-05',
            end_date: '2025-01-15',
            comment: 'Классика',
            rating: '9/10',
            time_played: '30 часов',
        },
    ],
};

const mockLeaderboard: LeaderboardRow[] = [
    {
        Player_id: '1',
        total: 54,
        comleted: 13,
        drop: 2,
        reroll: 0,
        current_game: 'Dark Souls III',
    },
    {
        Player_id: '2',
        total: 32,
        comleted: 10,
        drop: 4,
        reroll: 1,
        current_game: 'Elden Ring',
    },
    {
        Player_id: '3',
        total: 28,
        comleted: 8,
        drop: 1,
        reroll: 0,
        current_game: null,
    },
];

// Функции для получения моков
export function getMockPlayers(): Player[] {
    return [...mockPlayers];
}

export function getMockPlayer(id: string): Player | null {
    return mockPlayers.find(p => p.ID === id) || null;
}

export function getMockGames(): Game[] {
    return [...mockGames];
}

export function getMockGame(id: string): Game | null {
    return mockGames.find(g => g.ID === id) || null;
}

export function getMockPlayerGames(playerId: string): GamePlayed[] {
    return mockGamesPlayed[playerId] || [];
}

export function getMockLeaderboard(): LeaderboardRow[] {
    return [...mockLeaderboard];
}

export function mockLogin(username: string, password: string): AuthResponse {
    const player = mockPlayers.find(p => p.Username === username);
    if (player && password === 'password') {
        return {
            token: `mock-token-${player.ID}`,
            player: { ...player, Password: undefined },
        };
    }
    throw new Error('Неверное имя пользователя или пароль');
}

export function mockRegister(username: string, password: string): AuthResponse {
    const newPlayer: Player = {
        ID: String(mockPlayers.length + 1),
        Username: username,
        Description: 'Новый игрок',
        Img: `https://api.dicebear.com/7.x/avataaars/svg?seed=${username}`,
        Email: `${username}@example.com`,
    };
    mockPlayers.push(newPlayer);
    return {
        token: `mock-token-${newPlayer.ID}`,
        player: { ...newPlayer, Password: undefined },
    };
}

