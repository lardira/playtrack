export interface IUser {
    ID: string;
    Username: string;
    Email?: string;
    AvatarUrl?: string;
    Score?: number;
}

export interface IGame {
    ID: string;
    Title: string;
    Description?: string;
    Genre?: string;
    ImageUrl?: string;
}

export interface ILeaderboardEntry {
    ID: string;
    Username: string;
    Score: number;
}
