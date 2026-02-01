<script lang="ts">
    import Modal from './Modal.svelte';
    import { updatePlayedGame } from '../api';
    import type { PlayedGame, PlayedGameStatus } from '../types';

    export let isOpen = false;
    export let onClose: () => void = () => {};
    export let playedGame: PlayedGame | null = null;
    export let playerId = '';
    export let gameTitle = '';
    /** Вызывается после успешного сохранения, чтобы обновить список */
    export let onSaved: () => void = () => {};

    let comment = '';
    let completedAt = '';
    let playTimeHours = 0;
    let playTimeMinutes = 0;
    let points = 0;
    let rating: number | '' = '';
    let status: PlayedGameStatus = 'in_progress';
    let loading = false;
    let error = '';
    /** Инициализируем форму только при открытии другой записи, чтобы не затирать ввод при каждом ре-рендере */
    let lastOpenedId: number | null = null;

    const STATUS_OPTIONS: { value: PlayedGameStatus; label: string }[] = [
        { value: 'added', label: 'Добавлено' },
        { value: 'in_progress', label: 'В процессе' },
        { value: 'completed', label: 'Пройдено' },
        { value: 'dropped', label: 'Дроп' },
        { value: 'rerolled', label: 'Реролл' },
    ];

    /** ISO date string (YYYY-MM-DD) из backend даты или null */
    function parseCompletedAtDate(iso: string | null): string {
        if (!iso) {
            const today = new Date();
            return today.getFullYear() + '-' + String(today.getMonth() + 1).padStart(2, '0') + '-' + String(today.getDate()).padStart(2, '0');
        }
        return iso.slice(0, 10);
    }

    $: if (isOpen && playedGame && playedGame.id !== lastOpenedId) {
        lastOpenedId = playedGame.id;
        comment = playedGame.comment ?? '';
        completedAt = parseCompletedAtDate(playedGame.completed_at);
        points = playedGame.points;
        rating = playedGame.rating ?? '';
        status = playedGame.status;
        if (playedGame.play_time) {
            const h = playedGame.play_time.match(/(\d+)H/);
            const m = playedGame.play_time.match(/(\d+)M/);
            playTimeHours = h ? parseInt(h[1], 10) : 0;
            playTimeMinutes = m ? parseInt(m[1], 10) : 0;
        } else {
            playTimeHours = 0;
            playTimeMinutes = 0;
        }
        error = '';
    }
    $: if (!playedGame) lastOpenedId = null;

    function buildPlayTime(): string | null {
        if (playTimeHours === 0 && playTimeMinutes === 0) return null;
        const parts = [];
        if (playTimeHours > 0) parts.push(`${playTimeHours}H`);
        if (playTimeMinutes > 0) parts.push(`${playTimeMinutes}M`);
        return parts.length ? `PT${parts.join('')}` : null;
    }

    async function handleSubmit() {
        if (!playedGame || !playerId) return;
        const ratingNum = rating === '' ? null : Number(rating);
        if (ratingNum !== null && (ratingNum < 0 || ratingNum > 100)) {
            error = 'Рейтинг от 0 до 100';
            return;
        }

        loading = true;
        error = '';

        try {
            await updatePlayedGame(playerId, playedGame.id, {
                comment: comment.trim() || null,
                completed_at: completedAt ? `${completedAt}T00:00:00Z` : null,
                play_time: buildPlayTime(),
                points,
                rating: ratingNum,
                status,
            });
            onSaved();
            onClose();
        } catch (err: any) {
            error = err?.message ?? 'Не удалось сохранить';
        } finally {
            loading = false;
        }
    }

    function handleClose() {
        onClose();
    }
</script>

<Modal isOpen={isOpen} title="Редактировать запись" onClose={handleClose}>
    {#if playedGame}
        <p class="text-surface-400 text-sm mb-4">{gameTitle}</p>
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            {#if error}
                <div class="p-4 bg-red-500/20 border border-red-500 rounded-lg text-red-400 text-sm">
                    {error}
                </div>
            {/if}

            <div>
                <label for="edit-comment" class="block text-sm font-medium mb-2">Комментарий</label>
                <textarea
                    id="edit-comment"
                    class="input w-full min-h-[80px]"
                    bind:value={comment}
                    disabled={loading}
                />
            </div>

            <div>
                <label for="edit-status" class="block text-sm font-medium mb-2">Статус</label>
                <select id="edit-status" class="input w-full" bind:value={status} disabled={loading}>
                    {#each STATUS_OPTIONS as opt}
                        <option value={opt.value}>{opt.label}</option>
                    {/each}
                </select>
            </div>

            <div>
                <label for="edit-points" class="block text-sm font-medium mb-2">Очки</label>
                <input
                    id="edit-points"
                    type="number"
                    class="input w-full"
                    bind:value={points}
                    disabled={loading}
                />
            </div>

            <div>
                <label for="edit-rating" class="block text-sm font-medium mb-2">Рейтинг (0–100)</label>
                <input
                    id="edit-rating"
                    type="number"
                    min="0"
                    max="100"
                    class="input w-full"
                    bind:value={rating}
                    disabled={loading}
                />
            </div>

            <div role="group" aria-labelledby="edit-play-time-label">
                <span id="edit-play-time-label" class="block text-sm font-medium mb-2">Время игры (play_time)</span>
                <div class="flex gap-2 items-center">
                    <input
                        id="edit-play-time-h"
                        type="number"
                        min="0"
                        class="input w-24"
                        bind:value={playTimeHours}
                        disabled={loading}
                        aria-label="Часы"
                    />
                    <span>ч</span>
                    <input
                        id="edit-play-time-m"
                        type="number"
                        min="0"
                        max="59"
                        class="input w-24"
                        bind:value={playTimeMinutes}
                        disabled={loading}
                        aria-label="Минуты"
                    />
                    <span>мин</span>
                </div>
            </div>

            <div>
                <label for="edit-completed-at" class="block text-sm font-medium mb-2">Дата завершения (completed_at)</label>
                <input
                    id="edit-completed-at"
                    type="date"
                    class="input w-full"
                    bind:value={completedAt}
                    disabled={loading}
                />
            </div>

            <div class="flex gap-2 pt-2">
                <button
                    type="submit"
                    class="btn variant-filled-primary"
                    disabled={loading}
                >
                    {loading ? 'Сохранение…' : 'Сохранить'}
                </button>
                <button
                    type="button"
                    class="btn variant-ghost-surface"
                    disabled={loading}
                    on:click={handleClose}
                >
                    Отмена
                </button>
            </div>
        </form>
    {/if}
</Modal>
