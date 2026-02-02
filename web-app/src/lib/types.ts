export interface Player {
    id: string; // UUID
    username: string;
    img: string | null;
    email: string | null;
    description: string | null;
    created_at: string; // ISO date string
}

export interface LeaderboardPlayer {
    player_id: string;
    completed: number;
    total: number;
    dropped: number;
    rerolled: number;
}

export interface LeaderboardRow {
    player: Player;
    currentGame: string | null;
    points: number;
    completed: number;
    dropped: number;
    rerolled: number;
}

export interface Game {
    id: number;
    points: number;
    hours_to_beat: number;
    title: string;
    url: string | null;
    created_at: string; // ISO date string
}

export interface PlayedGame {
    id: number;
    player_id: string;
    game_id: number;
    points: number;
    comment: string | null;
    rating: number | null;
    status: PlayedGameStatus;
    started_at: string; // ISO date string
    completed_at: string | null; // ISO date string
    play_time: string | null; // ISO duration string
}

export interface AuthResponse {
    token?: string;
    player?: Player;
}

export type PlayedGameStatus =
    | "added"
    | "in_progress"
    | "completed"
    | "dropped"
    | "rerolled";