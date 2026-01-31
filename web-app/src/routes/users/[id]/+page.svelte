<script lang="ts">
    import { page } from "$app/stores";
    import type { Player, PlayedGame } from "../../../lib/types";
    import { gamesPlayedMock, playersMock } from "../../../lib/mocks";
    import { user } from "../../../stores/user";
    import ChangePasswordModal from "../../../lib/components/ChangePasswordModal.svelte";

    let player: Player | null = null;
    let loading = true;
    let currentUser: Player | null = null;
    let showChangePasswordModal = false;

    user.subscribe((u) => (currentUser = u));

    const players: Player[] = playersMock;

    let playedGames: PlayedGame[] = gamesPlayedMock;

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
        player = players.find((p) => p.id === id) ?? null;
        loading = false;
    } else {
        loading = false;
    }

    $: playerColor = player ? getPlayerColor(player.username) : "#f97316";

    // Функция для форматирования ISO duration в читаемый формат
    function formatPlayTime(isoDuration: string | null): string {
        if (!isoDuration) return "0h";

        // Парсим ISO 8601 duration (PT45H30M, PT2H15M и т.д.)
        const hoursMatch = isoDuration.match(/(\d+)H/);
        const minutesMatch = isoDuration.match(/(\d+)M/);

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

        <div class="flex flex-col md:flex-row md:items-center gap-6">
            <!-- AVATAR -->
            {#if player.img}
                <img
                    src={player.img}
                    alt={player.username}
                    class="w-24 h-24 rounded-full object-cover"
                />
            {:else}
                <div
                    class="w-24 h-24 rounded-full flex items-center justify-center text-3xl font-bold"
                    style={`background:${playerColor}; color:#000`}
                >
                    {player.username[0].toUpperCase()}
                </div>
            {/if}

            <!-- INFO -->
            <div class="flex-1">
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
            </div>

            <!-- CREATED AT -->
            <div class="text-right">
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
            <p class="text-3xl font-bold">128</p>
        </div>

        <div class="card p-6 bg-surface shadow-md rounded-xl">
            <p class="text-sm text-surface-400">Очки</p>
            <p class="text-3xl font-bold">42</p>
        </div>

        <div class="card p-6 bg-surface shadow-md rounded-xl">
            <p class="text-sm text-surface-400">Процент пройденных игр</p>
            <p class="text-3xl font-bold">61%</p>
        </div>
    </section>

    <!-- GAMES LIST -->
    <section class="max-w-5xl mx-auto mt-14">
        <h2 class="text-2xl font-bold mb-6">🎮 Игры</h2>

        <div class="space-y-4">
            {#each playedGames.filter((pg) => pg.player_id === player.id) as playedGame}
                <div
                    class="group rounded-xl p-5 bg-surface shadow-md transition hover:shadow-xl"
                    style={`border-left: 6px solid ${STATUS_META[playedGame.status].color}`}
                >
                    <div
                        class="flex flex-col md:flex-row md:items-center gap-4"
                    >
                        <!-- INFO -->
                        <div class="flex-1 min-w-0">
                            <p class="text-lg font-semibold truncate">
                                Game ID: {playedGame.game_id}
                            </p>
                            <p class="text-sm text-surface-400">
                                {STATUS_META[playedGame.status].label}
                            </p>
                        </div>

                        <!-- STATS -->
                        <div
                            class="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6 text-sm"
                        >
                            <!-- Время игры -->
                            <div
                                class="flex items-center"
                                style={playedGame.play_time
                                    ? ""
                                    : "visibility: hidden;"}
                            >
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2 flex-shrink-0"
                                    style={`background:${STATUS_META[playedGame.status].color}`}
                                ></span>
                                <div class="min-w-0">
                                    <p class="text-surface-400 text-xs">
                                        Время игры
                                    </p>
                                    <p class="font-bold whitespace-nowrap">
                                        {playedGame.play_time
                                            ? formatPlayTime(
                                                  playedGame.play_time,
                                              )
                                            : "-"}
                                    </p>
                                </div>
                            </div>

                            <!-- Очки -->
                            <div class="flex items-center">
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2 flex-shrink-0"
                                    style={`background:${STATUS_META[playedGame.status].color}`}
                                ></span>
                                <div class="min-w-0">
                                    <p class="text-surface-400 text-xs">Очки</p>
                                    <p
                                        class="font-bold whitespace-nowrap"
                                        style={`color:${playerColor}`}
                                    >
                                        {playedGame.points}
                                    </p>
                                </div>
                            </div>

                            <!-- Рейтинг -->
                            <div
                                class="flex items-center"
                                style={playedGame.rating
                                    ? ""
                                    : "visibility: hidden;"}
                            >
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2 flex-shrink-0"
                                    style={`background:${STATUS_META[playedGame.status].color}`}
                                ></span>
                                <div class="min-w-0">
                                    <p class="text-surface-400 text-xs">
                                        Рейтинг
                                    </p>
                                    <p class="font-bold whitespace-nowrap">
                                        {playedGame.rating
                                            ? `${playedGame.rating}/100`
                                            : "-"}
                                    </p>
                                </div>
                            </div>

                            <!-- Дата старта -->
                            <div class="flex items-center">
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2 flex-shrink-0"
                                    style={`background:${STATUS_META[playedGame.status].color}`}
                                ></span>
                                <div class="min-w-0">
                                    <p class="text-surface-400 text-xs">
                                        Дата старта
                                    </p>
                                    <p class="font-bold whitespace-nowrap">
                                        {formatDate(playedGame.started_at)}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- ACTION -->
                        <div class="text-right flex-shrink-0">
                            <button
                                class="btn btn-sm variant-ghost-surface whitespace-nowrap"
                            >
                                Подробнее
                            </button>
                        </div>
                    </div>
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
/>
