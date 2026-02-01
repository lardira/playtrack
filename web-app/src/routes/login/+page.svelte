<script lang="ts">
    import { goto } from "$app/navigation";
    import { tick } from "svelte";
    import { login, register } from "../../lib/api";
    import { user, token, loadUserFromToken } from "../../stores/user";
    import { setTokenCookie } from "../../lib/cookies";
    import { onMount } from "svelte";

    let isLogin = true;
    let loading = false;
    let error = "";

    // Login form
    let loginUsername = "";
    let loginPassword = "";

    // Register form
    let registerUsername = "";
    let registerPassword = "";
    let registerPasswordConfirm = "";
    let registerEmail = "";
    let registerImg = "";

    onMount(() => {
        // Если пользователь уже залогинен, перенаправляем на главную
        user.subscribe((u) => {
            if (u) {
                goto("/");
            }
        });
    });

    async function handleLogin() {
        if (!loginUsername || !loginPassword) {
            error = "Заполните все поля";
            return;
        }

        loading = true;
        error = "";

        try {
            const response = await login({
                username: loginUsername,
                password: loginPassword,
            });

            const t = response.token;
            setTokenCookie(t);
            token.set(t);
            await loadUserFromToken();
            await tick();
            goto("/", { replaceState: true });
        } catch (err: any) {
            error = err.message || "Ошибка при входе";
        } finally {
            loading = false;
        }
    }

    async function handleRegister() {
        if (
            !registerUsername ||
            !registerPassword ||
            !registerPasswordConfirm
        ) {
            error = "Заполните все обязательные поля";
            return;
        }

        if (registerPassword !== registerPasswordConfirm) {
            error = "Пароли не совпадают";
            return;
        }

        if (registerPassword.length < 8) {
            error = "Пароль должен быть не менее 8 символов";
            return;
        }

        if (registerUsername.length < 4) {
            error = "Имя пользователя должно быть не менее 4 символов";
            return;
        }

        loading = true;
        error = "";

        try {
            const response = await register({
                username: registerUsername,
                password: registerPassword,
                email: registerEmail || null,
                img: registerImg || null,
            });

            // После регистрации автоматически логинимся
            const loginResponse = await login({
                username: registerUsername,
                password: registerPassword,
            });

            const t = loginResponse.token;
            setTokenCookie(t);
            token.set(t);
            await loadUserFromToken();
            await tick();
            goto("/", { replaceState: true });
        } catch (err: any) {
            error = err.message || "Ошибка при регистрации";
        } finally {
            loading = false;
        }
    }

    function toggleMode() {
        isLogin = !isLogin;
        error = "";
        loginUsername = "";
        loginPassword = "";
        registerUsername = "";
        registerPassword = "";
        registerPasswordConfirm = "";
        registerEmail = "";
        registerImg = "";
    }
</script>

<div class="min-h-screen flex items-center justify-center p-4">
    <div class="w-full max-w-md">
        <div class="bg-surface rounded-2xl p-8 shadow-xl">
            <h1 class="text-3xl font-bold text-center mb-8">
                {isLogin ? "Вход" : "Регистрация"}
            </h1>

            {#if error}
                <div
                    class="mb-4 p-4 bg-red-500/20 border border-red-500 rounded-lg text-red-400 text-sm"
                >
                    {error}
                </div>
            {/if}

            {#if isLogin}
                <!-- Login Form -->
                <form on:submit|preventDefault={handleLogin} class="space-y-4">
                    <div>
                        <label
                            for="login-username"
                            class="block text-sm font-medium mb-2"
                        >
                            Имя пользователя
                        </label>
                        <input
                            id="login-username"
                            type="text"
                            bind:value={loginUsername}
                            required
                            minlength="4"
                            class="input w-full"
                            placeholder="Введите имя пользователя"
                            disabled={loading}
                        />
                    </div>

                    <div>
                        <label
                            for="login-password"
                            class="block text-sm font-medium mb-2"
                        >
                            Пароль
                        </label>
                        <input
                            id="login-password"
                            type="password"
                            bind:value={loginPassword}
                            required
                            minlength="8"
                            class="input w-full"
                            placeholder="Введите пароль"
                            disabled={loading}
                        />
                    </div>

                    <button
                        type="submit"
                        class="btn w-full variant-filled-primary"
                        disabled={loading}
                    >
                        {loading ? "Вход..." : "Войти"}
                    </button>
                </form>
            {:else}
                <!-- Register Form -->
                <form
                    on:submit|preventDefault={handleRegister}
                    class="space-y-4"
                >
                    <div>
                        <label
                            for="register-username"
                            class="block text-sm font-medium mb-2"
                        >
                            Имя пользователя *
                        </label>
                        <input
                            id="register-username"
                            type="text"
                            bind:value={registerUsername}
                            required
                            minlength="4"
                            class="input w-full"
                            placeholder="Минимум 4 символа"
                            disabled={loading}
                        />
                    </div>

                    <div>
                        <label
                            for="register-email"
                            class="block text-sm font-medium mb-2"
                        >
                            Email
                        </label>
                        <input
                            id="register-email"
                            type="email"
                            bind:value={registerEmail}
                            class="input w-full"
                            placeholder="email@example.com"
                            disabled={loading}
                        />
                    </div>

                    <div>
                        <label
                            for="register-img"
                            class="block text-sm font-medium mb-2"
                        >
                            URL изображения
                        </label>
                        <input
                            id="register-img"
                            type="url"
                            bind:value={registerImg}
                            class="input w-full"
                            placeholder="https://example.com/avatar.jpg"
                            disabled={loading}
                        />
                    </div>

                    <div>
                        <label
                            for="register-password"
                            class="block text-sm font-medium mb-2"
                        >
                            Пароль *
                        </label>
                        <input
                            id="register-password"
                            type="password"
                            bind:value={registerPassword}
                            required
                            minlength="8"
                            class="input w-full"
                            placeholder="Минимум 8 символов"
                            disabled={loading}
                        />
                    </div>

                    <div>
                        <label
                            for="register-password-confirm"
                            class="block text-sm font-medium mb-2"
                        >
                            Подтвердите пароль *
                        </label>
                        <input
                            id="register-password-confirm"
                            type="password"
                            bind:value={registerPasswordConfirm}
                            required
                            minlength="8"
                            class="input w-full"
                            placeholder="Повторите пароль"
                            disabled={loading}
                        />
                    </div>

                    <button
                        type="submit"
                        class="btn w-full variant-filled-primary"
                        disabled={loading}
                    >
                        {loading ? "Регистрация..." : "Зарегистрироваться"}
                    </button>
                </form>
            {/if}

            <div class="mt-6 text-center">
                <button
                    type="button"
                    on:click={toggleMode}
                    class="text-primary-400 hover:text-primary-300 text-sm"
                    disabled={loading}
                >
                    {isLogin
                        ? "Нет аккаунта? Зарегистрироваться"
                        : "Уже есть аккаунт? Войти"}
                </button>
            </div>

            <div class="mt-4 text-center">
                <a
                    href="/"
                    class="text-surface-400 hover:text-surface-300 text-sm"
                >
                    ← На главную
                </a>
            </div>
        </div>
    </div>
</div>
