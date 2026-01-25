<script lang="ts">
    import { login, register } from "$lib/api";
    import { user } from "$lib/stores/user";
    import type { Player } from "$lib/types";

    let username = "";
    let password = "";
    let isRegister = false;
    let loading = false;
    let error = "";

    async function handleSubmit() {
        if (!username || !password) {
            error = "Заполните все поля";
            return;
        }

        loading = true;
        error = "";

        try {
            let response;
            if (isRegister) {
                response = await register(username, password);
            } else {
                response = await login(username, password);
            }

            if (response.player) {
                user.set(response.player as Player);
            }

            if (typeof window !== "undefined") {
                if (response.player) {
                    location.href = `/users/${response.player.ID}`;
                } else {
                    location.href = "/";
                }
            }
        } catch (e: any) {
            error = e.message || "Ошибка при авторизации";
            console.error("Auth error", e);
        } finally {
            loading = false;
        }
    }
</script>

<div class="max-w-md mx-auto">
    <div class="card p-8">
        <h1 class="text-3xl font-bold mb-6 text-center text-primary-400">
            {isRegister ? "Регистрация" : "Вход"}
        </h1>

        {#if error}
            <div class="mb-4 p-3 bg-error-500/20 border border-error-500/50 rounded-lg text-error-500 text-sm">
                {error}
            </div>
        {/if}

        <form on:submit|preventDefault={handleSubmit} class="space-y-4">
            <label class="block">
                <span class="label-text">
                    {isRegister ? "Никнейм *" : "Имя пользователя *"}
                </span>
                <input
                    id="username"
                    type="text"
                    placeholder={isRegister
                        ? "Введите никнейм"
                        : "Введите имя пользователя"}
                    bind:value={username}
                    required
                    class="input mt-1"
                />
            </label>

            <label class="block">
                <span class="label-text">Пароль *</span>
                <input
                    id="password"
                    type="password"
                    placeholder="Введите пароль"
                    bind:value={password}
                    required
                    class="input mt-1"
                />
            </label>

            <button
                type="submit"
                disabled={loading}
                class="btn variant-filled-primary w-full"
            >
                {loading
                    ? "Загрузка..."
                    : isRegister
                      ? "Зарегистрироваться"
                      : "Войти"}
            </button>
        </form>

        <div class="mt-4 text-center">
            <button
                on:click={() => {
                    isRegister = !isRegister;
                    error = "";
                }}
                class="btn variant-ghost-primary text-sm"
            >
                {isRegister
                    ? "Уже есть аккаунт? Войти"
                    : "Нет аккаунта? Зарегистрироваться"}
            </button>
        </div>
    </div>
</div>
