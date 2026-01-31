<script lang="ts">
    import { page } from "$app/stores";
    import type { Player, Game } from "../../../lib/types";

    let player: Player | null = null;
    let loading = true;
    // let activeStatusFilter: GameStatus | "all" = "all";

    const players: Player[] = [
        {
            ID: "1",
            Username: "Alice",
            Score: 1840,
            Color: "#f97316",
            Description: "Speedrunner",
            Img: "",
        },
        {
            ID: "2",
            Username: "Bob",
            Score: 1520,
            Color: "#22c55e",
            Description: "Achievement hunter",
            Img: "",
        },
        {
            ID: "3",
            Username: "Charlie",
            Score: 1390,
            Color: "#3b82f6",
            Description: "PvP monster",
            Img: "",
        },
        {
            ID: "4",
            Username: "Diana",
            Score: 980,
            Color: "#a855f7",
            Description: "Casual grinder",
            Img: "",
        },
    ];

    let games: Game[] = [
        {
            ID: "g1",
            Title: "Dark Souls III",
            Genre: "Action RPG",
            Playtime: 120,
            Score: 520,
            LastPlayed: "2025-01-18",
            Status: "completed",
            HoursToBeat: 50,
            URL: "",
            CreatedAt: "2025-01-01",
        },
        {
            ID: "g2",
            Title: "Hades",
            Genre: "Roguelike",
            Playtime: 64,
            Score: 310,
            LastPlayed: "2025-01-10",
            Status: "in_progress",
            HoursToBeat: 30,
            URL: "",
            CreatedAt: "2025-01-01",
        },
        {
            ID: "g3",
            Title: "Factorio",
            Genre: "Simulation",
            Playtime: 230,
            Score: 720,
            LastPlayed: "2025-01-02",
            Status: "completed",
            HoursToBeat: 100,
            URL: "",
            CreatedAt: "2025-01-01",
        },
    ];

    const STATUS_META: Record<string, { label: string; color: string }> = {
        completed: {
            label: "Пройдено",
            color: "#22c55e", // green
        },
        dropped: {
            label: "Дроп",
            color: "#ef4444", // red
        },
        reroll: {
            label: "Реролл",
            color: "#38bdf8", // blue
        },
        in_progress: {
            label: "В процессе",
            color: "#facc15", // yellow
        },
    };

    $: id = $page.params.id;
    $: if (id) {
        player = players.find((p) => p.ID === id) ?? null;
        loading = false;
    } else {
        loading = false;
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
        style={`border-left: 8px solid ${player.Color}`}
    >
        <!-- glow -->
        <div
            class="absolute inset-0 -z-10 blur-[120px]"
            style={`background:${player.Color}33`}
        ></div>

        <div class="flex flex-col md:flex-row md:items-center gap-6">
            <!-- AVATAR -->
            <div
                class="w-24 h-24 rounded-full flex items-center justify-center text-3xl font-bold"
                style={`background:${player.Color}; color:#000`}
            >
                {player.Username[0]}
            </div>

            <!-- INFO -->
            <div class="flex-1">
                <h1
                    class="text-3xl font-extrabold"
                    style={`color:${player.Color}`}
                >
                    {player.Username}
                </h1>
                <p class="text-surface-400 mt-1">
                    {player.Description}
                </p>
            </div>

            <!-- SCORE -->
            <div class="text-right">
                <p class="text-4xl font-bold">{player.Score}</p>
                <p class="text-sm text-surface-400">Накопленные</p>
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
            <p class="text-sm text-surface-400">Достижения</p>
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
            {#each games as game}
                <div
                    class="group rounded-xl p-5 bg-surface shadow-md transition hover:shadow-xl"
                    style={`border-left: 6px solid ${STATUS_META[game.Status].color}`}
                >
                    <div
                        class="flex flex-col md:flex-row md:items-center gap-4"
                    >
                        <!-- INFO -->
                        <div class="flex-1">
                            <p class="text-lg font-semibold">
                                {game.Title}
                            </p>
                            <p class="text-sm text-surface-400">
                                {game.Genre}
                            </p>
                        </div>

                        <!-- STATS -->
                        <div class="flex gap-6 text-sm">
                            <div class="flex items-center">
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2"
                                    style={`background:${STATUS_META[game.Status].color}`}
                                ></span>

                                <div>
                                    <p class="text-surface-400">Время игры</p>
                                    <p class="font-bold">{game.Playtime}h</p>
                                </div>
                            </div>

                            {#if game.Score !== undefined}
                                <div class="flex items-center">
                                    <span
                                        class="inline-block w-1 h-4 rounded-full mr-2"
                                        style={`background:${STATUS_META[game.Status].color}`}
                                    ></span>

                                    <div>
                                        <p class="text-surface-400">Очки</p>
                                        <p
                                            class="font-bold"
                                            style={`color:${player.Color}`}
                                        >
                                            {game.Score}
                                        </p>
                                    </div>
                                </div>
                            {/if}

                            <div class="flex items-center">
                                <span
                                    class="inline-block w-1 h-4 rounded-full mr-2"
                                    style={`background:${STATUS_META[game.Status].color}`}
                                ></span>

                                <div>
                                    <p class="text-surface-400">Дата старта</p>
                                    <p class="font-bold">
                                        {new Date(
                                            game.LastPlayed,
                                        ).toLocaleDateString()}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <!-- ACTION -->
                        <div class="text-right">
                            <button class="btn btn-sm variant-ghost-surface">
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
        <a href="/" class="btn variant-ghost-surface mr-2">
            ← На главную страницу
        </a>
    </section>
{/if}
