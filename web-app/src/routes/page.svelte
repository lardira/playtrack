<script lang="ts">
    import Chart from "$lib/components/Chart.svelte";
    export let data;

    const rows = [...data.leaderboard].sort((a, b) => b.total - a.total);
</script>

<div class="container mx-auto px-4 space-y-8">
    <header>
        <h1 class="text-4xl font-bold">🏆 Leaderboard</h1>
        <p class="text-surface-500">Рейтинг игроков по очкам и достижениям</p>
    </header>

    {#if rows.length}
        <Chart data={rows} players={data.players} />
    {/if}

    <div class="overflow-x-auto">
        <table class="table w-full">
            <thead>
                <tr>
                    <th>Место</th><th>Пользователь</th><th>Текущая игра</th><th
                        >Поинты</th
                    ><th>Пройдено</th><th>Дроп</th><th>Реролл</th>
                </tr>
            </thead>
            <tbody>
                {#each rows as row, i}
                    <tr
                        class="hover:bg-surface-700/40 {i === 0
                            ? 'bg-primary-500/10'
                            : ''}"
                    >
                        <td
                            >{i + 1}
                            {i === 0
                                ? "🥇"
                                : i === 1
                                  ? "🥈"
                                  : i === 2
                                    ? "🥉"
                                    : ""}</td
                        >
                        <td>
                            <a
                                href="/users/{row.Player_id}"
                                class="font-semibold text-primary-400 hover:text-primary-300"
                            >
                                {data.players.find(
                                    (p) => p.ID === row.Player_id,
                                )?.Username} →
                            </a>
                        </td>
                        <td>{row.current_game ?? "Перерыв"}</td>
                        <td
                            ><span class="badge variant-soft-warning"
                                >{row.total}</span
                            ></td
                        >
                        <td
                            ><span class="badge variant-soft-success"
                                >{row.comleted}</span
                            ></td
                        >
                        <td
                            ><span class="badge variant-soft-error"
                                >{row.drop}</span
                            ></td
                        >
                        <td
                            ><span class="badge variant-soft-warning"
                                >{row.reroll}</span
                            ></td
                        >
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
