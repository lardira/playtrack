import { PlayedGameStatus } from '$lib/playedGameStatus';
import type { PlayedGame } from '$lib/types';

export function isActiveStatus(status: PlayedGameStatus): boolean {
    return status === PlayedGameStatus.Added || status === PlayedGameStatus.InProgress;
}

export function isCountedStatus(status: PlayedGameStatus): boolean {
    return status !== PlayedGameStatus.InProgress;
}

export function countByStatus(games: PlayedGame[], status: PlayedGameStatus): number {
    return games.filter((pg) => pg.status === status).length;
}

export function hasNonTerminatedGame(games: PlayedGame[]): boolean {
    return games.some((pg) => isActiveStatus(pg.status));
}

export function computeProfileStats(playedGames: PlayedGame[]) {
    const countedGames = playedGames.filter((pg) => isCountedStatus(pg.status));
    const gamesExcludingReroll = countedGames.filter(
        (pg) => pg.status !== PlayedGameStatus.Rerolled,
    );
    const completedCount = countByStatus(playedGames, PlayedGameStatus.Completed);
    return {
        totalGames: countedGames.length,
        totalPoints: countedGames.reduce((sum, pg) => sum + pg.points, 0),
        completedPercent:
            gamesExcludingReroll.length > 0
                ? Math.round((completedCount / gamesExcludingReroll.length) * 100)
                : 0,
    };
}

export function summarizePlayedGamesForLeaderboard(played: PlayedGame[]) {
    const terminated = played.filter((p) => !isActiveStatus(p.status));
    return {
        points: terminated.reduce((s, p) => s + p.points, 0),
        completed: countByStatus(played, PlayedGameStatus.Completed),
        dropped: countByStatus(played, PlayedGameStatus.Dropped),
        rerolled: countByStatus(played, PlayedGameStatus.Rerolled),
    };
}
