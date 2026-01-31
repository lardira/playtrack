<script lang="ts">
    import { setPassword } from '../api';
    import Modal from './Modal.svelte';

    export let isOpen = false;
    export let onClose: () => void = () => {};

    let currentPassword = '';
    let newPassword = '';
    let confirmPassword = '';
    let loading = false;
    let error = '';
    let success = false;

    function resetForm() {
        currentPassword = '';
        newPassword = '';
        confirmPassword = '';
        error = '';
        success = false;
    }

    function handleClose() {
        resetForm();
        onClose();
    }

    async function handleSubmit() {
        if (!newPassword || !confirmPassword) {
            error = 'Заполните все поля';
            return;
        }

        if (newPassword !== confirmPassword) {
            error = 'Пароли не совпадают';
            return;
        }

        if (newPassword.length < 8) {
            error = 'Пароль должен быть не менее 8 символов';
            return;
        }

        loading = true;
        error = '';

        try {
            await setPassword({ password: newPassword });
            success = true;
            setTimeout(() => {
                handleClose();
            }, 2000);
        } catch (err: any) {
            error = err.message || 'Ошибка при смене пароля';
        } finally {
            loading = false;
        }
    }
</script>

<Modal {isOpen} title="Смена пароля" onClose={handleClose}>
    {#if success}
        <div class="p-4 bg-green-500/20 border border-green-500 rounded-lg text-green-400 text-sm mb-4">
            Пароль успешно изменен!
        </div>
    {:else}
        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            {#if error}
                <div class="p-4 bg-red-500/20 border border-red-500 rounded-lg text-red-400 text-sm">
                    {error}
                </div>
            {/if}

            <div>
                <label for="new-password" class="block text-sm font-medium mb-2">
                    Новый пароль *
                </label>
                <input
                    id="new-password"
                    type="password"
                    bind:value={newPassword}
                    required
                    minlength="8"
                    class="input w-full"
                    placeholder="Минимум 8 символов"
                    disabled={loading}
                />
            </div>

            <div>
                <label for="confirm-password" class="block text-sm font-medium mb-2">
                    Подтвердите пароль *
                </label>
                <input
                    id="confirm-password"
                    type="password"
                    bind:value={confirmPassword}
                    required
                    minlength="8"
                    class="input w-full"
                    placeholder="Повторите новый пароль"
                    disabled={loading}
                />
            </div>

            <div class="flex gap-3 pt-4">
                <button
                    type="button"
                    on:click={handleClose}
                    class="btn flex-1 variant-ghost-surface"
                    disabled={loading}
                >
                    Отмена
                </button>
                <button
                    type="submit"
                    class="btn flex-1 variant-filled-primary"
                    disabled={loading}
                >
                    {loading ? 'Сохранение...' : 'Изменить пароль'}
                </button>
            </div>
        </form>
    {/if}
</Modal>

