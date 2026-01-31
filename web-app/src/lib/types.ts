export interface Player {
    ID: string;
    Username: string;
    Description: string;
    Score: number;
    Img: string;
    Email?: string;
    Color: string; // HEX
}

export interface LeaderboardRow {
    Player_id: string;
    total: number;
    completed: number;
    drop: number;
    reroll: number;
}

export interface Game {
    ID: string;
    Score: number;
    HoursToBeat: number;
    Playtime: number;
    Title: string;
    URL: string;
    CreatedAt: string;
    Genre: string;
    LastPlayed: string;
    Status: GameStatus;
}

export interface GamePlayed {
    Player_id: string;
    game_id: string;
    status: 'В процессе' | 'Пройдено' | 'Дроп' | 'Реролл';
    scores: number;
    start_date: string;
    end_date?: string | null;
    comment: string;
    rating: string;
    time_played?: string;
}

export interface AuthResponse {
    token?: string;
    player?: Player;
}

export type GameStatus =
    | "completed"
    | "dropped"
    | "reroll"
    | "in_progress";