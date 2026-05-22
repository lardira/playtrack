export enum PlayedGameStatus {
    Added = 'added',
    InProgress = 'in_progress',
    Completed = 'completed',
    Dropped = 'dropped',
    Rerolled = 'rerolled',
}

export const PLAYED_GAME_STATUS_META: Record<
    PlayedGameStatus,
    { label: string; color: string }
> = {
    [PlayedGameStatus.Added]: { label: 'Добавлено', color: '#94a3b8' },
    [PlayedGameStatus.InProgress]: { label: 'В процессе', color: '#facc15' },
    [PlayedGameStatus.Completed]: { label: 'Пройдено', color: '#22c55e' },
    [PlayedGameStatus.Dropped]: { label: 'Дроп', color: '#ef4444' },
    [PlayedGameStatus.Rerolled]: { label: 'Реролл', color: '#38bdf8' },
};

export const PLAYED_GAME_STATUS_ORDER: PlayedGameStatus[] = [
    PlayedGameStatus.Added,
    PlayedGameStatus.InProgress,
    PlayedGameStatus.Completed,
    PlayedGameStatus.Dropped,
    PlayedGameStatus.Rerolled,
];

export const EDIT_STATUS_OPTIONS = PLAYED_GAME_STATUS_ORDER.map((value) => ({
    value,
    label: PLAYED_GAME_STATUS_META[value].label,
}));
