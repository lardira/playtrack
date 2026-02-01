<script lang="ts">
    import Modal from './Modal.svelte';
    import { updatePlayer } from '../api';

    export let isOpen = false;
    export let onClose: () => void = () => {};
    export let playerId = '';
    export let currentDescription = '';
    export let onSaved: (description: string) => void = () => {};

    let description = '';
    let loading = false;
    let error = '';

    $: if (isOpen) {
        description = currentDescription;
        error = '';
    }

    async function handleSubmit() {
        loading = true;
        error = '';
        try {
            await updatePlayer(playerId, { description: description.trim() || null });
            onSaved(description.trim() || '');
            onClose();
        } catch (err: any) {
            error = err?.message ?? 'Не удалось сохранить';
        } finally {
            loading = false;
        }
    }
</script>

<Modal isOpen={isOpen} title="Описание профиля" onClose={onClose}>
    <p class="text-surface-400 text-sm mb-3">Его видят все, кто заходит на вашу страницу.</p>
    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#if error}
            <div class="p-3 rounded-lg bg-red-500/15 border border-red-500/50 text-red-400 text-sm">{error}</div>
        {/if}
        <div>
            <label for="edit-description" class="block text-sm font-medium text-surface-300 mb-1.5">Описание</label>
            <textarea
                id="edit-description"
                class="input w-full min-h-[120px] px-3 py-2.5 rounded-lg border border-surface-600 bg-surface-800 focus:border-primary-500 focus:outline-none resize-y"
                bind:value={description}
                disabled={loading}
                placeholder="Расскажите о себе..."
            />
        </div>
        <div class="flex gap-2">
            <button type="submit" class="btn variant-filled-primary" disabled={loading}>
                {loading ? 'Сохранение…' : 'Сохранить'}
            </button>
            <button type="button" class="btn variant-ghost-surface" disabled={loading} on:click={onClose}>
                Отмена
            </button>
        </div>
    </form>
</Modal>
