<script lang="ts">
    import GameForms from "$lib/components/GameForms.svelte";
    export let data;
    let showCreate = false;
</script>

<div class="container mx-auto px-4 space-y-8">
    <div class="card p-6 flex gap-6 items-start">
        <img
            src={data.player.Img || "/favicon.png"}
            alt={data.player.Username}
            class="w-32 h-32 rounded-xl border border-surface-700"
        />
        <div class="space-y-2">
            <h2 class="text-3xl font-bold text-primary-400">
                {data.player.Username}
            </h2>
            <p class="text-surface-500">
                {data.player.Description || "Нет описания"}
            </p>
        </div>
    </div>

    {#if data.isOwner}
        <button
            class="btn variant-filled-primary"
            on:click={() => (showCreate = true)}>➕ Создать запись</button
        >
        <GameForms isOpen={showCreate} playerId={data.player.ID} />
    {/if}

    <div class="overflow-x-auto">
        <table class="table w-full">
            <thead>
                <tr
                    ><th>Дата</th><th>Время</th><th>Игра</th><th>Статус</th><th
                        >Очки</th
                    ><th>Комментарий</th><th>Оценка</th></tr
                >
            </thead>
            <tbody>
                {#each data.games as g}
                    <tr class="hover:bg-surface-700/40">
                        <td>{g.start_date || "-"}</td>
                        <td>{g.time_played || "-"}</td>
                        <td>
                            {#if g.gameURL}
                                <a
                                    href={g.gameURL}
                                    target="_blank"
                                    class="text-primary-400 hover:text-primary-300"
                                    >{g.gameTitle}</a
                                >
                            {:else}{g.gameTitle}{/if}
                        </td>
                        <td
                            ><span class="badge variant-soft-primary"
                                >{g.status}</span
                            ></td
                        >
                        <td
                            ><span class="badge variant-soft-warning"
                                >{g.scores}</span
                            ></td
                        >
                        <td class="max-w-xs truncate" title={g.comment}
                            >{g.comment || "-"}</td
                        >
                        <td>{g.rating || "-"}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
