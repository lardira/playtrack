<script lang="ts">
    import Modal from './Modal.svelte';
    import { updatePlayedGame } from '../api';
    import type { PlayedGame } from '../types';
    import { PlayedGameStatus, EDIT_STATUS_OPTIONS } from '../playedGameStatus';
    import { isActiveStatus } from '../../helpers/playedGameStatus';
    import {
        splitIsoDateTime,
        DEFAULT_START_TIME,
        DEFAULT_COMPLETED_TIME,
        toUtcIsoString,
        isValidUtcDateTime,
        isCompletedAfterStart,
    } from '../../helpers/dateTime';

    export let isOpen = false;
    export let onClose: () => void = () => {};
    export let playedGame: PlayedGame | null = null;
    export let playerId = '';
    export let gameTitle = '';
    export let onSaved: () => void = () => {};

    let comment = '';
    let startedAt = '';
    let startedTime = '';
    let completedAt = '';
    let completedTime = '';
    let playTimeHours = 0;
    let playTimeMinutes = 0;
    let points = 0;
    let rating: number | '' = '';
    let status: PlayedGameStatus = PlayedGameStatus.InProgress;
    let loading = false;
    let error = '';
    let lastOpenedId: number | null = null;

    $: if (isOpen && playedGame && playedGame.id !== lastOpenedId) {
        lastOpenedId = playedGame.id;
        comment = playedGame.comment ?? '';
        const started = splitIsoDateTime(playedGame.started_at);
        startedAt = started.date;
        startedTime = started.time;
        const completed = splitIsoDateTime(playedGame.completed_at);
        completedAt = completed.date;
        completedTime = completed.time;
        points = playedGame.points;
        rating = playedGame.rating ?? '';
        status = playedGame.status;
        if (playedGame.play_time) {
            const s = playedGame.play_time;
            const h = s.match(/(\d+)h/i);
            const m = s.match(/(\d+)m/i);
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
        if (playTimeHours > 0) parts.push(`${playTimeHours}h`);
        if (playTimeMinutes > 0) parts.push(`${playTimeMinutes}m`);
        return parts.length ? parts.join('') : null;
    }

    async function handleSubmit() {
        if (!playedGame || !playerId) return;
        const ratingNum = rating === '' ? null : Number(rating);
        if (ratingNum !== null && (ratingNum < 0 || ratingNum > 100)) {
            error = 'Рейтинг от 0 до 100';
            return;
        }

        const terminated = !isActiveStatus(status);

        if (terminated) {
            if (!startedAt) {
                error = 'Укажите дату начала';
                return;
            }
            if (!completedAt) {
                error = 'Укажите дату завершения';
                return;
            }
            if (!isValidUtcDateTime(startedAt, startedTime, DEFAULT_START_TIME)) {
                error = 'Некорректная дата начала';
                return;
            }
            if (!isValidUtcDateTime(completedAt, completedTime, DEFAULT_COMPLETED_TIME)) {
                error = 'Некорректная дата завершения';
                return;
            }
            if (!isCompletedAfterStart(startedAt, startedTime, completedAt, completedTime)) {
                error = 'Дата завершения должна быть позже даты начала';
                return;
            }
        }

        loading = true;
        error = '';

        try {
            const payload: Parameters<typeof updatePlayedGame>[2] = {
                comment: comment.trim() || null,
                started_at: startedAt ? toUtcIsoString(startedAt, startedTime, DEFAULT_START_TIME) : null,
                completed_at: terminated && completedAt ? toUtcIsoString(completedAt, completedTime, DEFAULT_COMPLETED_TIME) : null,
                play_time: buildPlayTime(),
                rating: ratingNum,
            };
            if (status !== playedGame.status) {
                payload.status = status;
            }
            await updatePlayedGame(playerId, playedGame.id, payload);
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
        <p class="text-surface-400 text-sm mb-5 pb-3 border-b border-surface-600">{gameTitle}</p>
        <form on:submit|preventDefault={handleSubmit} class="space-y-5">
            {#if error}
                <div class="p-3 rounded-lg bg-red-500/15 border border-red-500/50 text-red-400 text-sm">
                    {error}
                </div>
            {/if}

            <div class="space-y-1.5">
                <label for="edit-comment" class="block text-sm font-medium text-surface-300">Комментарий</label>
                <textarea
                    id="edit-comment"
                    class="input w-full min-h-[88px] px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none resize-y"
                    bind:value={comment}
                    disabled={loading}
                />
            </div>

            <div class="grid grid-cols-2 gap-x-4 gap-y-5">
                <div class="space-y-1.5 pr-6 border-r-2 border-surface-600">
                    <label for="edit-status" class="block text-sm font-medium text-surface-300">Статус</label>
                    <select id="edit-status" class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]" bind:value={status} disabled={loading}>
                        {#each EDIT_STATUS_OPTIONS as opt}
                            <option value={opt.value}>{opt.label}</option>
                        {/each}
                    </select>
                </div>
                <div class="space-y-1.5 pl-2">
                    <label for="edit-rating" class="block text-sm font-medium text-surface-300">Рейтинг (0–100)</label>
                    <input
                        id="edit-rating"
                        type="number"
                        min="0"
                        max="100"
                        class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                        bind:value={rating}
                        disabled={loading}
                    />
                </div>
                <div role="group" class="space-y-1.5 pr-6 border-r-2 border-surface-600" aria-labelledby="edit-play-time-label">
                    <span id="edit-play-time-label" class="block text-sm font-medium text-surface-300">Время игры</span>
                    <div class="flex gap-2 items-center">
                        <input
                            id="edit-play-time-h"
                            type="number"
                            min="0"
                            class="input w-20 px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                            bind:value={playTimeHours}
                            disabled={loading}
                            aria-label="Часы"
                        />
                        <span class="text-surface-400 text-sm">ч</span>
                        <input
                            id="edit-play-time-m"
                            type="number"
                            min="0"
                            max="59"
                            class="input w-20 px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                            bind:value={playTimeMinutes}
                            disabled={loading}
                            aria-label="Минуты"
                        />
                        <span class="text-surface-400 text-sm">мин</span>
                    </div>
                </div>
                <div class="space-y-1.5 pl-2">
                    <span class="block text-sm font-medium text-surface-400">Очки</span>
                    <div class="rounded-lg border border-surface-600 bg-surface-700/50 text-surface-400 min-h-[2.75rem] px-3 py-2.5 flex items-center select-none text-sm" aria-readonly="true" role="textbox">{points}</div>
                </div>
            </div>

            <div class="space-y-1.5">
                <label for="edit-started-at" class="block text-sm font-medium text-surface-300">Дата старта</label>
                <div class="grid grid-cols-2 gap-3">
                    <input
                        id="edit-started-at"
                        type="date"
                        class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                        bind:value={startedAt}
                        disabled={loading}
                    />
                    <input
                        id="edit-started-time"
                        type="time"
                        step="1"
                        class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                        bind:value={startedTime}
                        disabled={loading}
                    />
                </div>
            </div>

            {#if !isActiveStatus(status)}
                <div class="space-y-1.5">
                    <label for="edit-completed-at" class="block text-sm font-medium text-surface-300">Дата завершения</label>
                    <div class="grid grid-cols-2 gap-3">
                        <input
                            id="edit-completed-at"
                            type="date"
                            class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                            bind:value={completedAt}
                            disabled={loading}
                        />
                        <input
                            id="edit-completed-time"
                            type="time"
                            step="1"
                            class="input w-full px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none min-h-[2.75rem]"
                            bind:value={completedTime}
                            disabled={loading}
                        />
                    </div>
                </div>
            {/if}

            <div class="flex gap-3 pt-4 border-t border-surface-600">
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
