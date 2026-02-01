<script lang="ts">
    import { page } from "$app/stores";
    import type { Player, PlayedGame, Game } from "../../../lib/types";
    import { user } from "../../../stores/user";
    import {
        getPlayer,
        getPlayerPlayedGames,
        getGames,
        createGame,
        createPlayedGame,
        updatePlayedGame,
    } from "../../../lib/api";
    import ChangePasswordModal from "../../../lib/components/ChangePasswordModal.svelte";
    import EditDescriptionModal from "../../../lib/components/EditDescriptionModal.svelte";
    import EditPlayedGameModal from "../../../lib/components/EditPlayedGameModal.svelte";

    let player: Player | null = null;
    let loading = true;
    let currentUser: Player | null = null;
    let showChangePasswordModal = false;
    let showEditDescriptionModal = false;
    let editPlayedGame: PlayedGame | null = null;
    let expandedPlayedId: number | null = null;
    let playedGames: PlayedGame[] = [];
    let allGames: Game[] = [];

    user.subscribe((u) => (currentUser = u));

    $: gamesMap = Object.fromEntries(allGames.map((g) => [g.id, g])) as Record<number, Game>;
    $: gameTitle = (gameId: number) => gamesMap[gameId]?.title ?? `ID: ${gameId}`;

    // Генерируем цвет на основе username для UI
    function getPlayerColor(username: string): string {
        const colors = [
            "#f97316",
            "#22c55e",
            "#3b82f6",
            "#a855f7",
            "#ec4899",
            "#14b8a6",
        ];
        const hash = username
            .split("")
            .reduce((acc, char) => acc + char.charCodeAt(0), 0);
        return colors[hash % colors.length];
    }

    const STATUS_META: Record<string, { label: string; color: string }> = {
        completed: {
            label: "Пройдено",
            color: "#22c55e", // green
        },
        dropped: {
            label: "Дроп",
            color: "#ef4444", // red
        },
        rerolled: {
            label: "Реролл",
            color: "#38bdf8", // blue
        },
        in_progress: {
            label: "В процессе",
            color: "#facc15", // yellow
        },
        added: {
            label: "Добавлено",
            color: "#94a3b8", // gray
        },
    };

    $: id = $page.params.id;

    $: if (id) {
        loading = true;
        const requestedId = id;
        Promise.all([getPlayer(id), getPlayerPlayedGames(id), getGames()])
            .then(([p, games, gamesList]) => {
                if ($page.params.id === requestedId) {
                    player = p;
                    playedGames = games ?? [];
                    allGames = gamesList ?? [];
                }
            })
            .catch(() => {
                if ($page.params.id === requestedId) {
                    player = null;
                    playedGames = [];
                    allGames = [];
                }
            })
            .finally(() => {
                if ($page.params.id === requestedId) loading = false;
            });
    } else {
        player = null;
        playedGames = [];
        allGames = [];
        loading = false;
    }

    $: playerColor = player ? getPlayerColor(player.username) : "#f97316";
    // sub из JWT = id игрока; на своей странице показываем «Редактировать»
    $: isOwnProfile = !!currentUser && !!player && currentUser.id === player.id;

    // Динамическая статистика из playedGames
    $: totalGames = playedGames.length;
    $: totalPoints = playedGames.reduce((sum, pg) => sum + pg.points, 0);
    $: gamesExcludingReroll = playedGames.filter((pg) => pg.status !== "rerolled");
    $: completedCount = playedGames.filter((pg) => pg.status === "completed").length;
    $: completedPercent =
        gamesExcludingReroll.length > 0
            ? Math.round((completedCount / gamesExcludingReroll.length) * 100)
            : 0;

    /** Парсит duration из API: Go-формат "34h30m0s" или ISO 8601 "PT45H30M" */
    function formatPlayTime(duration: string | null): string {
        if (!duration) return "0h";

        const s = duration;
        const hoursMatch = s.match(/(\d+)h/i);
        const minutesMatch = s.match(/(\d+)m/i);

        const hours = hoursMatch ? parseInt(hoursMatch[1], 10) : 0;
        const minutes = minutesMatch ? parseInt(minutesMatch[1], 10) : 0;

        if (hours > 0 && minutes > 0) {
            return `${hours}h ${minutes}m`;
        } else if (hours > 0) {
            return `${hours}h`;
        } else if (minutes > 0) {
            return `${minutes}m`;
        }
        return "0h";
    }

    // Функция для форматирования даты в формат день/месяц/год
    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        const day = date.getDate().toString().padStart(2, "0");
        const month = (date.getMonth() + 1).toString().padStart(2, "0");
        const year = date.getFullYear();
        return `${day}/${month}/${year}`;
    }

    function handleEditPlayedGame(playedGame: PlayedGame) {
        if (!player) return;
        editPlayedGame = playedGame;
    }

    async function refreshPlayedGamesAfterEdit() {
        if (!player) return;
        const list = await getPlayerPlayedGames(player.id);
        playedGames = list ?? playedGames;
    }

    async function refreshPlayerAfterDescriptionSave() {
        if (!id) return;
        const p = await getPlayer(id);
        if (p) player = p;
    }

    // --- Создание новой записи (только на своей странице) ---
    let showNewRow = false;
    let newRecordTitle = "";
    let newRecordHoursToBeat = 1;
    let newRecordUrl = "";
    let selectedGame: Game | null = null;
    let createLoading = false;
    let createError = "";

    $: searchQuery = newRecordTitle.trim().toLowerCase();
    $: searchResults =
        searchQuery.length < 2
            ? []
            : allGames.filter((g) => g.title.toLowerCase().includes(searchQuery)).slice(0, 8);

    function openNewRow() {
        showNewRow = true;
        newRecordTitle = "";
        newRecordHoursToBeat = 1;
        newRecordUrl = "";
        selectedGame = null;
        createError = "";
    }

    function cancelNewRow() {
        showNewRow = false;
        newRecordTitle = "";
        newRecordHoursToBeat = 1;
        newRecordUrl = "";
        selectedGame = null;
        createError = "";
    }

    function selectGame(game: Game) {
        selectedGame = game;
        newRecordTitle = game.title;
        newRecordHoursToBeat = game.hours_to_beat;
        newRecordUrl = game.url ?? "";
    }

    /** Статусы, при которых нельзя добавить новую запись — у игрока может быть только одна незавершённая игра */
    const NONTERMINATED_STATUSES = ["added", "in_progress"];

    async function submitNewRecord() {
        if (!player || player.id !== currentUser?.id) return;
        const title = newRecordTitle.trim();
        if (!title) {
            createError = "Введите название игры";
            return;
        }

        const hasNonTerminated = playedGames.some((pg) =>
            NONTERMINATED_STATUSES.includes(pg.status),
        );
        if (hasNonTerminated) {
            createError =
                "Сначала завершите текущую игру (Пройдено / Дроп / Реролл) или измените её статус в карточке выше.";
            return;
        }

        createLoading = true;
        createError = "";

        try {
            let gameId: number;
            if (selectedGame) {
                gameId = selectedGame.id;
            } else {
                const { id } = await createGame({
                    title,
                    hours_to_beat: newRecordHoursToBeat,
                    url: newRecordUrl.trim() || null,
                });
                gameId = id;
            }

            const { id: playedGameId } = await createPlayedGame(player.id, gameId);
            await updatePlayedGame(player.id, playedGameId, { status: "in_progress" });

            const [updatedPlayed, updatedGames] = await Promise.all([
                getPlayerPlayedGames(player.id),
                getGames(),
            ]);
            playedGames = updatedPlayed ?? playedGames;
            allGames = updatedGames ?? allGames;

            cancelNewRow();
        } catch (err: any) {
            const msg = err?.message ?? "";
            if (msg.includes("nonterminated")) {
                createError =
                    "У вас уже есть игра в статусе «В процессе» или «Добавлено». Сначала завершите её (Пройдено / Дроп / Реролл) через кнопку «Редактировать».";
            } else {
                createError = msg || "Не удалось создать запись";
            }
        } finally {
            createLoading = false;
        }
    }
</script>

{#if loading}
    <div class="flex justify-center mt-20">
        <span class="loading loading-spinner loading-lg"></span>
    </div>
{:else if !player}
    <div class="text-center mt-20">
        <h1 class="text-2xl font-bold mb-2">Player not found</h1>
        <a href="/" class="btn variant-filled-primary">Back to home</a>
    </div>
{:else}
    <!-- PROFILE HEADER -->
    <section
        class="relative max-w-5xl mx-auto mt-10 rounded-2xl p-8 bg-surface shadow-xl overflow-hidden"
        style={`border-left: 8px solid ${playerColor}`}
    >
        <!-- glow -->
        <div
            class="absolute inset-0 -z-10 blur-[120px]"
            style={`background:${playerColor}33`}
        ></div>

        <div class="flex flex-col md:flex-row md:items-start gap-6">
            <!-- AVATAR (увеличенная) -->
            {#if player.img}
                <img
                    src={player.img}
                    alt={player.username}
                    class="w-32 h-32 md:w-36 md:h-36 rounded-full object-cover flex-shrink-0"
                />
            {:else}
                <div
                    class="w-32 h-32 md:w-36 md:h-36 rounded-full flex items-center justify-center text-4xl font-bold flex-shrink-0"
                    style={`background:${playerColor}; color:#000`}
                >
                    {player.username[0].toUpperCase()}
                </div>
            {/if}

            <!-- INFO + ОПИСАНИЕ -->
            <div class="flex-1 min-w-0">
                <h1
                    class="text-3xl font-extrabold"
                    style={`color:${playerColor}`}
                >
                    {player.username}
                </h1>
                {#if player.email}
                    <p class="text-surface-400 mt-1">
                        {player.email}
                    </p>
                {/if}
                {#if player.description}
                    <p class="text-surface-300 mt-3 text-sm leading-relaxed whitespace-pre-wrap">{player.description}</p>
                {:else if isOwnProfile}
                    <p class="text-surface-500 mt-3 text-sm italic">Описание не задано. Однажды оно появится здесь..</p>
                {/if}
                {#if isOwnProfile}
                    <button
                        type="button"
                        class="btn btn-sm variant-ghost-surface mt-2"
                        on:click={() => (showEditDescriptionModal = true)}
                    >
                        {player.description ? 'Редактировать описание' : 'Добавить описание'}
                    </button>
                {/if}
            </div>

            <!-- CREATED AT -->
            <div class="text-right flex-shrink-0">
                <p class="text-sm text-surface-400">Зарегистрирован</p>
                <p class="text-lg font-bold">
                    {new Date(player.created_at).toLocaleDateString()}
                </p>
            </div>
        </div>
    </section>

    <!-- STATS -->
    <section class="max-w-5xl mx-auto mt-10 grid gap-6 md:grid-cols-3">
        <div class="card p-6 bg-surface shadow-md rounded-xl">
            <p class="text-sm text-surface-400">Всего сыграно игр</p>
            <p class="text-3xl font-bold">{totalGames}</p>
        </div>

        <div class="card p-6 bg-surface shadow-md rounded-xl">
            <p class="text-sm text-surface-400">Очки</p>
            <p class="text-3xl font-bold">{totalPoints}</p>
        </div>

        <div class="card p-6 bg-surface shadow-md rounded-xl">
            <p class="text-sm text-surface-400">Процент пройденных игр</p>
            <p class="text-3xl font-bold">{completedPercent}%</p>
        </div>
    </section>

    <!-- GAMES LIST -->
    <section class="max-w-5xl mx-auto mt-14">
        <div class="flex items-center justify-between mb-6">
            <h2 class="text-2xl font-bold">🎮 Игры</h2>
            {#if isOwnProfile}
                {#if showNewRow}
                    <span class="text-surface-400 text-sm">Заполните title, hours_to_beat, url и нажмите «Подтвердить»</span>
                {:else}
                    <button
                        type="button"
                        class="btn btn-sm variant-filled-primary"
                        on:click={openNewRow}
                    >
                        Добавить запись
                    </button>
                {/if}
            {/if}
        </div>

        <div class="space-y-4">
            <!-- Новая запись (редактируемое поле) — только на своей странице -->
            {#if isOwnProfile && showNewRow}
                <div
                    class="rounded-xl p-5 bg-surface shadow-md border-2 border-dashed border-primary-500/50"
                >
                    <div class="flex flex-col gap-4">
                        <div class="relative">
                            <label for="new-game-title" class="block text-sm font-medium text-surface-400 mb-1"
                                >Название игры (title)</label
                            >
                            <input
                                id="new-game-title"
                                type="text"
                                class="input w-full"
                                placeholder="Введите название или выберите из списка"
                                bind:value={newRecordTitle}
                                disabled={createLoading}
                            />
                            {#if searchResults.length > 0}
                                <ul
                                    class="absolute z-10 mt-1 w-full rounded-lg bg-surface-800 border border-surface-600 shadow-lg max-h-48 overflow-auto"
                                >
                                    {#each searchResults as game}
                                        <li>
                                            <button
                                                type="button"
                                                class="w-full text-left px-4 py-2 hover:bg-surface-700"
                                                on:click={() => selectGame(game)}
                                            >
                                                {game.title}
                                                <span class="text-surface-500 text-sm ml-2"
                                                    >({game.hours_to_beat} ч, {game.points} очк.)</span
                                                >
                                            </button>
                                        </li>
                                    {/each}
                                </ul>
                            {/if}
                        </div>
                        <div>
                            <label for="new-game-hours" class="block text-sm font-medium text-surface-400 mb-1"
                                >Часов на прохождение (hours_to_beat)</label
                            >
                            <input
                                id="new-game-hours"
                                type="number"
                                min="1"
                                class="input w-full"
                                bind:value={newRecordHoursToBeat}
                                disabled={createLoading}
                            />
                        </div>
                        <div>
                            <label for="new-game-url" class="block text-sm font-medium text-surface-400 mb-1"
                                >URL (необязательно)</label
                            >
                            <input
                                id="new-game-url"
                                type="url"
                                class="input w-full"
                                placeholder="https://..."
                                bind:value={newRecordUrl}
                                disabled={createLoading}
                            />
                        </div>
                        {#if createError}
                            <p class="text-red-400 text-sm">{createError}</p>
                        {/if}
                        <div class="flex gap-2">
                            <button
                                type="button"
                                class="btn variant-filled-primary"
                                disabled={createLoading || !newRecordTitle.trim()}
                                on:click={submitNewRecord}
                            >
                                {createLoading ? "Создание…" : "Подтвердить"}
                            </button>
                            <button
                                type="button"
                                class="btn variant-ghost-surface"
                                disabled={createLoading}
                                on:click={cancelNewRow}
                            >
                                Отмена
                            </button>
                        </div>
                    </div>
                </div>
            {/if}

            {#each playedGames as playedGame}
                {@const isExpanded = expandedPlayedId === playedGame.id}
                <div
                    class="group rounded-xl p-5 bg-surface shadow-md transition hover:shadow-xl"
                    style={`border-left: 6px solid ${STATUS_META[playedGame.status].color}`}
                >
                    <div
                        class="flex flex-col md:flex-row md:items-center gap-4"
                    >
                        <!-- INFO (клик раскрывает) -->
                        <button
                            type="button"
                            class="flex-1 min-w-0 text-left flex items-center gap-2"
                            on:click={() => (expandedPlayedId = isExpanded ? null : playedGame.id)}
                        >
                            <span
                                class="text-surface-400 transition-transform"
                                class:rotate-90={isExpanded}
                            >▶</span>
                            <div class="min-w-0">
                                <p class="text-lg font-semibold truncate">
                                    {gameTitle(playedGame.game_id)}
                                </p>
                                <p class="text-sm text-surface-400">
                                    {STATUS_META[playedGame.status].label}
                                </p>
                            </div>
                        </button>

                        <!-- Краткие статы в одной строке (всегда видны) -->
                        <div class="flex flex-wrap items-center gap-4 text-sm text-surface-400">
                            <span>Очки: <strong style={`color:${playerColor}`}>{playedGame.points}</strong></span>
                            <span>Игра: {playedGame.play_time ? formatPlayTime(playedGame.play_time) : "0"}</span>
                            <span>Рейтинг: {playedGame.rating != null ? `${playedGame.rating}/100` : "—"}</span>
                            <span>Старт: {formatDate(playedGame.started_at)}</span>
                        </div>

                        <!-- ACTION -->
                        <div class="text-right flex-shrink-0">
                            {#if isOwnProfile}
                                <button
                                    type="button"
                                    class="btn btn-sm variant-ghost-surface whitespace-nowrap"
                                    on:click={() => handleEditPlayedGame(playedGame)}
                                >
                                    Редактировать
                                </button>
                            {/if}
                        </div>
                    </div>

                    <!-- Раскрытый блок: комментарий и доп. поля -->
                    {#if isExpanded}
                        <div class="mt-4 pt-4 border-t border-surface-600 space-y-2 text-sm">
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                                <div>
                                    <span class="text-surface-400">Очки</span>
                                    <p class="font-bold" style={`color:${playerColor}`}>{playedGame.points}</p>
                                </div>
                                <div>
                                    <span class="text-surface-400">Время игры (play_time)</span>
                                    <p class="font-bold">{playedGame.play_time ? formatPlayTime(playedGame.play_time) : "0"}</p>
                                </div>
                                <div>
                                    <span class="text-surface-400">Рейтинг</span>
                                    <p class="font-bold">{playedGame.rating != null ? `${playedGame.rating}/100` : "—"}</p>
                                </div>
                                <div>
                                    <span class="text-surface-400">Дата старта (started_at)</span>
                                    <p class="font-bold">{formatDate(playedGame.started_at)}</p>
                                </div>
                                {#if playedGame.completed_at}
                                    <div>
                                        <span class="text-surface-400">Дата завершения</span>
                                        <p class="font-bold">{formatDate(playedGame.completed_at)}</p>
                                    </div>
                                {/if}
                            </div>
                            {#if playedGame.comment}
                                <div>
                                    <span class="text-surface-400">Комментарий</span>
                                    <p class="mt-1 text-surface-200 whitespace-pre-wrap">{playedGame.comment}</p>
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </section>

    <!-- ACTIONS -->
    <section class="max-w-5xl mx-auto mt-12 text-center">
        <div class="flex gap-3 justify-center">
            <a href="/" class="btn variant-ghost-surface">
                ← На главную страницу
            </a>
            {#if currentUser && currentUser.id === player.id}
                <button
                    class="btn variant-filled-secondary"
                    on:click={() => (showChangePasswordModal = true)}
                >
                    Сменить пароль
                </button>
            {/if}
        </div>
    </section>
{/if}

<ChangePasswordModal
    isOpen={showChangePasswordModal}
    onClose={() => (showChangePasswordModal = false)}
    username={player?.username ?? ''}
/>

<EditPlayedGameModal
    isOpen={!!editPlayedGame}
    onClose={() => (editPlayedGame = null)}
    playedGame={editPlayedGame}
    playerId={player?.id ?? ''}
    gameTitle={editPlayedGame ? gameTitle(editPlayedGame.game_id) : ''}
    onSaved={refreshPlayedGamesAfterEdit}
/>

<EditDescriptionModal
    isOpen={showEditDescriptionModal}
    onClose={() => (showEditDescriptionModal = false)}
    playerId={player?.id ?? ''}
    currentDescription={player?.description ?? ''}
    onSaved={refreshPlayerAfterDescriptionSave}
/>
