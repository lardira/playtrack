<script lang="ts">
    import type { LeaderboardRow, Player } from "../types";
    export let data: LeaderboardRow[] = [];
    export const players: Player[] = [];

    let period = "week";
    let metric = "total";
</script>

<div class="card p-4 space-y-4">
    <div class="flex space-x-4">
        <button
            class="btn btn-sm"
            on:click={() => (period = "week")}
            class:btn-primary={period === "week"}>Неделя</button
        >
        <button
            class="btn btn-sm"
            on:click={() => (period = "month")}
            class:btn-primary={period === "month"}>Месяц</button
        >

        <button
            class="btn btn-sm"
            on:click={() => (metric = "total")}
            class:btn-primary={metric === "total"}>Поинты</button
        >
        <button
            class="btn btn-sm"
            on:click={() => (metric = "completed")}
            class:btn-primary={metric === "completed"}>Пройдено</button
        >
    </div>
    <div class="text-surface-400 text-sm">
        Период: {period} | Показано: {metric}
    </div>
    <div class="h-40 bg-surface-700 rounded-xl flex items-end space-x-1">
        {#each data as row}
            <div
                class="flex-1 bg-primary-400"
                style="height:{Math.min(row[metric] * 0.5, 100)}%"
            ></div>
        {/each}
    </div>
</div>
