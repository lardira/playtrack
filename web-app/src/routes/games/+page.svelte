<script lang="ts">
	import { onMount } from "svelte";
	import { getGames } from "../../lib/api";
	import type { Game } from "../../lib/types";

	let games: Game[] = [];
	let loading = true;
	let loadError = "";
	let loadGeneration = 0;

	onMount(() => {
		const gen = ++loadGeneration;
		getGames()
			.then((list) => {
				if (gen !== loadGeneration) return;
				games = list ?? [];
				loadError = "";
			})
			.catch((err) => {
				if (gen !== loadGeneration) return;
				loadError = err?.message ?? "Не удалось загрузить шаблоны игр";
				games = [];
			})
			.finally(() => {
				if (gen === loadGeneration) loading = false;
			});
		return () => {
			loadGeneration++;
		};
	});
</script>

<section class="max-w-5xl mx-auto mt-10">
	<h1 class="text-3xl font-bold mb-6">🎮 Шаблоны игр</h1>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg text-primary-500"></span>
		</div>
	{:else if loadError}
		<div class="rounded-xl bg-surface border border-red-500/30 p-6 text-center text-red-400">
			{loadError}
		</div>
	{:else}
		<div class="overflow-x-auto rounded-xl bg-surface shadow-md">
			<table class="table w-full">
				<thead>
					<tr>
						<th>Название</th>
						<th>Часов на прохождение</th>
						<th>Очки</th>
						<th>Ссылка</th>
					</tr>
				</thead>
				<tbody>
					{#each games as game}
						<tr>
							<td class="font-semibold">{game.title}</td>
							<td>{game.hours_to_beat}</td>
							<td>{game.points}</td>
							<td>
								{#if game.url}
									<a
										href={game.url}
										target="_blank"
										rel="noreferrer"
										class="text-primary-400 hover:underline"
									>
										Открыть
									</a>
								{:else}
									—
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}

	<div class="mt-8">
		<a href="/" class="btn variant-ghost-surface">← На главную</a>
	</div>
</section>
