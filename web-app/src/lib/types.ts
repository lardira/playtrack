export type Player = {
    ID: string;
    Username: string;
    Description: string;
    Img: string;
    Email?: string;
    Password?: string;
};


export type LeaderboardRow = {
    Player_id: string;
    total: number;
    comleted: number;
    drop: number;
    reroll: number;
    current_game?: string; // Название текущей игры
};


export type Game = {
    ID: string;
    Points: number;
    HoursToBeat: number;
    Title: string;
    URL: string;
    CreatedAt: string;
};


export type GamePlayed = {
    Player_id: string;
    game_id: string;
    status: string;
    scores: number;
    start_date: string;
    end_date: string | null;
    comment: string;
    rating: string;
    time_played?: string; // Время в игре
};


export type AuthResponse = {
    token?: string;
    player?: Player;
};