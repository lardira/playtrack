import type { PlayedGame, PlayedGameStatus } from '$lib/types';

export function isTerminatedStatus(status: PlayedGameStatus): boolean {
    return status === 'completed' || status === 'dropped' || status === 'rerolled';
}

export function isActiveStatus(status: PlayedGameStatus): boolean {
    return status === 'added' || status === 'in_progress';
}

export function isCountedStatus(status: PlayedGameStatus): boolean {
    return status !== 'in_progress';
}

export function countByStatus(games: PlayedGame[], status: PlayedGameStatus): number {
    return games.filter((pg) => pg.status === status).length;
}

export function hasNonTerminatedGame(games: PlayedGame[]): boolean {
    return games.some((pg) => isActiveStatus(pg.status));
}

export function computeProfileStats(playedGames: PlayedGame[]) {
    const countedGames = playedGames.filter((pg) => isCountedStatus(pg.status));
    const gamesExcludingReroll = countedGames.filter((pg) => pg.status !== 'rerolled');
    const completedCount = countByStatus(playedGames, 'completed');
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
    const terminated = played.filter((p) => isTerminatedStatus(p.status));
    return {
        points: terminated.reduce((s, p) => s + p.points, 0),
        completed: countByStatus(played, 'completed'),
        dropped: countByStatus(played, 'dropped'),
        rerolled: countByStatus(played, 'rerolled'),
    };
}
