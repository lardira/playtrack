<script lang="ts">
    import { api, getGames } from "$lib/api";
    import Modal from "./Modal.svelte";
    import type { Game } from "$lib/types";

    export let playerId: string;
    export let isOpen = false;

    let gameSearch = "";
    let selectedGame: Game | null = null;
    let status = "В процессе";
    let comment = "";
    let rating = "";
    let timePlayed = "";
    let games: Game[] = [];
    let loading = false;

    const currentDate = new Date().toISOString().split("T")[0];

    async function searchGames() {
        if (gameSearch.length < 2) {
            games = [];
            return;
        }
        loading = true;
        try {
            const allGames = await getGames();
            games = allGames
                .filter((g) =>
                    g.Title.toLowerCase().includes(gameSearch.toLowerCase()),
                )
                .slice(0, 10);
        } catch (e) {
            console.error("Failed to search games", e);
            games = [];
        } finally {
            loading = false;
        }
    }

    function selectGame(game: Game) {
        selectedGame = game;
        gameSearch = game.Title;
        games = [];
    }

    async function submit() {
        if (!selectedGame) {
            alert("Пожалуйста, выберите игру");
            return;
        }

        try {
            await api(`/players/${playerId}/games-played/${selectedGame.ID}`, {
                method: "PATCH",
                body: JSON.stringify({
                    status,
                    comment,
                    rating,
                    time_played: timePlayed,
                    start_date: currentDate,
                }),
            });

            selectedGame = null;
            gameSearch = "";
            status = "В процессе";
            comment = "";
            rating = "";
            timePlayed = "";
            isOpen = false;

            if (typeof window !== "undefined") {
                location.reload();
            }
        } catch (e) {
            console.error("Failed to create game entry", e);
            alert("Ошибка при создании записи");
        }
    }

    let searchTimeout: ReturnType<typeof setTimeout> | null = null;
    $: if (gameSearch) {
        if (searchTimeout) clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            searchGames();
        }, 300);
    }
</script>

<Modal bind:isOpen title="Создать запись об игре">
    <form on:submit|preventDefault={submit} class="space-y-4">
        <label class="block">
            <span class="label-text">Игра *</span>
            <div class="relative">
                <input
                    id="game-search"
                    type="text"
                    placeholder="Начните вводить название игры..."
                    bind:value={gameSearch}
                    class="input mt-1"
                />
                {#if games.length > 0}
                    <div class="absolute z-10 w-full mt-1 card max-h-60 overflow-y-auto p-0">
                        {#each games as game}
                            <button
                                type="button"
                                on:click={() => selectGame(game)}
                                class="w-full px-4 py-2 text-left hover:bg-surface-800 transition"
                            >
                                <div class="font-medium">{game.Title}</div>
                                <div class="text-xs text-surface-500">
                                    {game.Points} поинтов
                                </div>
                            </button>
                        {/each}
                    </div>
                {/if}
                {#if selectedGame}
                    <div class="mt-2 text-sm text-success-500">
                        ✓ Выбрано: {selectedGame.Title}
                    </div>
                {/if}
            </div>
        </label>

        <label class="block">
            <span class="label-text">Статус *</span>
            <select id="game-status" bind:value={status} class="input mt-1">
                <option>В процессе</option>
                <option>Пройдено</option>
                <option>Дроп</option>
                <option>Реролл</option>
            </select>
        </label>

        <label class="block">
            <span class="label-text">Время в игре</span>
            <input
                id="time-played"
                type="text"
                placeholder="Например: 15 часов"
                bind:value={timePlayed}
                class="input mt-1"
            />
        </label>

        <label class="block">
            <span class="label-text">Комментарий</span>
            <textarea
                id="game-comment"
                placeholder="Ваши впечатления об игре..."
                bind:value={comment}
                rows="4"
                class="input mt-1 resize-none"
            ></textarea>
        </label>

        <label class="block">
            <span class="label-text">Оценка</span>
            <input
                id="game-rating"
                type="text"
                placeholder="Например: 8/10"
                bind:value={rating}
                class="input mt-1"
            />
        </label>

        <div class="flex gap-3 pt-4">
            <button type="submit" class="btn variant-filled-primary flex-1">
                Создать запись
            </button>
            <button
                type="button"
                on:click={() => (isOpen = false)}
                class="btn variant-ghost-surface"
            >
                Отмена
            </button>
        </div>
    </form>
</Modal>
