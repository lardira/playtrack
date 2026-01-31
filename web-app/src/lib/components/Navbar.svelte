<script lang="ts">
    import { user } from "../../stores/user";
    import { onMount } from "svelte";
    import { get } from "svelte/store";

    let currentUser = null;
    onMount(() => {
        user.subscribe((value) => (currentUser = value));
    });

    function logout() {
        user.set(null);
        window.location.href = "/";
    }
</script>

<nav
    class="sticky top-0 z-50 bg-surface-900 border-b border-surface-700 backdrop-blur px-4 py-2 flex justify-between items-center"
>
    <a href="/" class="text-2xl font-bold text-primary-400">🎮 GameTracker</a>
    {#if currentUser}
        <div class="flex items-center space-x-4">
            <a
                href={`/users/${currentUser.id}`}
                class="text-primary-400 hover:text-primary-300"
                >{currentUser.username}</a
            >
            <button class="btn btn-ghost btn-sm" on:click={logout}>Выйти</button
            >
        </div>
    {:else}
        <a href="/login" class="btn btn-filled btn-primary btn-sm">Войти</a>
    {/if}
</nav>
