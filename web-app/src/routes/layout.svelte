<script lang="ts">
    import { onMount } from "svelte";
    import { user } from "$lib/stores/user";
    import type { Player } from "$lib/types";

    let currentUser: Player | null = null;

    user.subscribe((value) => {
        currentUser = value;
    });

    onMount(async () => {
        if (typeof window === "undefined") return;

        const token = localStorage.getItem("token");
        if (token && !currentUser) {
            try {
                const storedUser = localStorage.getItem("user");
                if (storedUser) {
                    try {
                        const parsed = JSON.parse(storedUser);
                        user.set(parsed);
                    } catch (e) {
                        console.error("Failed to parse stored user", e);
                    }
                }
            } catch (e) {
                console.error("Failed to load user", e);
            }
        }
    });

    function logout() {
        user.set(null);
        if (typeof window !== "undefined") {
            localStorage.removeItem("token");
            localStorage.removeItem("user");
        }
        if (typeof window !== "undefined") {
            location.href = "/";
        }
    }
</script>

<nav class="bg-surface-900 border-b border-surface-700 sticky top-0 z-50 backdrop-blur-sm">
    <div class="container mx-auto px-4">
        <div class="flex items-center justify-between h-16">
            <a href="/" class="text-2xl font-bold text-primary-400 hover:text-primary-300 transition-colors flex items-center gap-2">
                <span class="text-3xl">🎮</span>
                <span>GameTracker</span>
            </a>
            <div class="flex items-center gap-4">
                {#if $user}
                    <a
                        href="/users/{$user.ID}"
                        class="btn variant-ghost-primary"
                    >
                        {$user.Username}
                    </a>
                    <button
                        on:click={logout}
                        class="btn variant-ghost-surface"
                    >
                        Выйти
                    </button>
                {:else}
                    <a href="/login" class="btn variant-ghost-primary">
                        Войти
                    </a>
                {/if}
            </div>
        </div>
    </div>
</nav>

<main class="container mx-auto px-4 py-8 min-h-[calc(100vh-4rem)]">
    <slot />
</main>
