<script lang="ts">
    import type { LeaderboardRow, Player } from "$lib/types";

    export let data: LeaderboardRow[] = [];
    export let players: Player[] = [];

    type Period = "week" | "month";
    type Metric = "total" | "completed" | "drop" | "reroll";

    let period: Period = "week";
    let metric: Metric = "total";

    $: days = period === "week" ? 7 : 30;

    function mockSeries(total: number, daysCount: number) {
        const base = Math.max(1, Math.floor(total / daysCount));
        return Array.from({ length: daysCount }, (_, i) => base * (i + 1));
    }

    $: series = (() => {
        const daysCount = days;
        return data.map((row) => {
            const player = players.find((p) => p.ID === row.Player_id);
            const total =
                metric === "total"
                    ? row.total
                    : metric === "completed"
                      ? row.comleted
                      : metric === "drop"
                        ? row.drop
                        : row.reroll;

            return {
                name: player?.Username ?? "Unknown",
                color: (player as any)?.Color ?? "#60a5fa",
                values: mockSeries(total, daysCount),
            };
        });
    })();

    $: maxY = Math.max(...series.flatMap((s) => s.values), 1);
</script>

<div class="card p-4 space-y-4">
    <div class="flex flex-wrap gap-2 justify-between">
        <div class="flex gap-2">
            <button
                class="btn {period === 'week'
                    ? 'variant-filled-primary'
                    : 'variant-ghost-primary'}"
                on:click={() => (period = "week")}
            >
                Неделя
            </button>
            <button
                class="btn {period === 'month'
                    ? 'variant-filled-primary'
                    : 'variant-ghost-primary'}"
                on:click={() => (period = "month")}
            >
                Месяц
            </button>
        </div>
        <div class="flex gap-2">
            <button
                class="btn {metric === 'total'
                    ? 'variant-filled-warning'
                    : 'variant-soft-warning'}"
                on:click={() => (metric = "total")}
            >
                Поинты
            </button>
            <button
                class="btn {metric === 'completed'
                    ? 'variant-filled-success'
                    : 'variant-soft-success'}"
                on:click={() => (metric = "completed")}
            >
                Пройдено
            </button>
            <button
                class="btn {metric === 'drop'
                    ? 'variant-filled-error'
                    : 'variant-soft-error'}"
                on:click={() => (metric = "drop")}
            >
                Дроп
            </button>
            <button
                class="btn {metric === 'reroll'
                    ? 'variant-filled-warning'
                    : 'variant-soft-warning'}"
                on:click={() => (metric = "reroll")}
            >
                Реролл
            </button>
        </div>
    </div>

    <svg viewBox="0 0 100 40" class="w-full h-64 bg-surface-800 rounded-xl p-2">
        {#each series as s}
            <polyline
                fill="none"
                stroke={s.color}
                stroke-width="0.8"
                points={s.values
                    .map(
                        (v, i) =>
                            `${(i / (days - 1)) * 100},${40 - (v / maxY) * 35}`,
                    )
                    .join(" ")}
            />
        {/each}
    </svg>

    <p class="text-surface-500 text-sm">
        Период: {period === "week" ? "Неделя" : "Месяц"} | Показано: {metric}
    </p>
</div>
