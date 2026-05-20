<script lang="ts">
	import { onMount } from "svelte";
	import { user } from "../stores/user";
	import type { Player, LeaderboardRow } from "../lib/types";
	import { getPlayers, getGames, getPlayerPlayedGames } from "../lib/api";
	import { getPlayerColor } from "../lib/colorUtils";

	let currentUser: Player | null = null;
	let leaderboardRows: LeaderboardRow[] = [];
	let loading = true;
	let loadError = "";

	const placeStyles = [
		{
			bg: "#facc15",
			color: "#000",
			label: "🥇",
			glow: "0 0 40px rgba(250,204,21,0.35)",
		},
		{
			bg: "#e5e7eb",
			color: "#000",
			label: "🥈",
			glow: "0 0 30px rgba(229,231,235,0.25)",
		},
		{
			bg: "#d97706",
			color: "#000",
			label: "🥉",
			glow: "0 0 25px rgba(217,119,6,0.25)",
		},
	];


	$: sortedByPoints = [...leaderboardRows].sort((a, b) => b.points - a.points);
	$: topThree = sortedByPoints.slice(0, 3);

	function scrollToLeaderboard() {
		document.getElementById("leaderboard")?.scrollIntoView({ behavior: "smooth", block: "start" });
	}

	onMount(() => {
		user.subscribe((value) => (currentUser = value));
		Promise.all([getPlayers(), getGames()])
			.then(([players, games]) => {
				const gamesMap = new Map(games.map((g) => [g.id, g]));
				return Promise.all(
					players.map((player) =>
						getPlayerPlayedGames(player.id)
							.then((played) => {
								const notTerminated = played.filter(
									(p) => p.status === "added" || p.status === "in_progress"
								);
								const lastInProgress = [...notTerminated].sort(
									(a, b) =>
										new Date(b.started_at).getTime() -
										new Date(a.started_at).getTime()
								)[0];
								const currentGame = lastInProgress
									? gamesMap.get(lastInProgress.game_id)?.title ?? null
									: null;
								const terminated = played.filter(
									(p) =>
										p.status === "completed" ||
										p.status === "dropped" ||
										p.status === "rerolled"
								);
								const points = terminated.reduce((s, p) => s + p.points, 0);
								const completed = played.filter((p) => p.status === "completed").length;
								const dropped = played.filter((p) => p.status === "dropped").length;
								const rerolled = played.filter((p) => p.status === "rerolled").length;
								return {
									player,
									currentGame,
									points,
									completed,
									dropped,
									rerolled,
								};
							})
							.catch(() => ({
								player,
								currentGame: null,
								points: 0,
								completed: 0,
								dropped: 0,
								rerolled: 0,
							}))
					)
				).then((rows) => {
					leaderboardRows = rows;
					loadError = "";
				});
			})
			.catch((err) => {
				loadError = err?.message ?? "Не удалось загрузить данные";
				leaderboardRows = [];
			})
			.finally(() => (loading = false));
	});
</script>

<!-- HERO -->
<section
	class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary-600/20 via-surface-900 to-secondary-600/20 p-10 text-center shadow-xl"
>
	<h1 class="text-4xl font-extrabold tracking-tight mb-3">
		<span class="text-primary-400">SEASON I</span> HERE
	</h1>
	<p class="text-lg text-surface-300 max-w-xl mx-auto">
		Стартовал первый сезон игры, у нас будет 8 недель игр, которые нужно пройти и заработать как можно больше очков.
	</p>

	<!-- <div class="mt-6 flex justify-center gap-3">
		{#if currentUser}
			<a
				href={`/users/${currentUser.id}`}
				class="btn variant-filled-primary">Go to my profile</a
			>
		{:else}
			<a href="/login" class="btn variant-filled-primary"
				>Login to start tracking</a
			>
		{/if}
		<button
			type="button"
			class="btn variant-ghost-surface"
			on:click={scrollToLeaderboard}
		>
			View leaderboard
		</button>
	</div> -->

	<!-- glow -->
	<div class="absolute inset-0 -z-10 blur-[120px] bg-primary-500/20"></div>
</section>

<!-- TOP PLAYERS -->
<section class="mt-14 max-w-5xl mx-auto">
	<h2 class="text-3xl font-bold text-center mb-8">Лучшие игроки</h2>

	{#if loading}
		<div class="flex justify-center py-12">
			<span class="loading loading-spinner loading-lg text-primary-500"></span>
		</div>
	{:else if loadError}
		<div class="rounded-xl bg-surface border border-red-500/30 p-6 text-center text-red-400">
			{loadError}
		</div>
	{:else}
	<div class="grid gap-6 md:grid-cols-3">
		{#each topThree as { player }, index}
			<a
				href={`/users/${player.id}`}
				class="relative rounded-2xl p-6 bg-surface shadow-lg transition hover:scale-[1.03]"
				style={`
					border-top: 6px solid ${placeStyles[index].bg};
					box-shadow: ${placeStyles[index].glow};
				`}
			>
				<!-- PLACE BADGE -->
				<div
					class="absolute -top-4 -right-4 w-12 h-12 rounded-full flex items-center justify-center text-xl font-bold"
					style={`
						background: ${placeStyles[index].bg};
						color: ${placeStyles[index].color};
					`}
				>
					{placeStyles[index].label}
				</div>

				<!-- CONTENT -->
				<p
					class="text-xl font-extrabold mb-1"
					style={`color:${getPlayerColor(player.username)}`}
				>
					{player.username}
				</p>

				{#if player.email}
					<p class="text-sm text-surface-400 mb-4">
						{player.email}
					</p>
				{/if}

				<div class="text-sm text-surface-400">
					Зарегистрирован: {new Date(
						player.created_at,
					).toLocaleDateString()}
				</div>
			</a>
		{/each}
	</div>
	{/if}
</section>

<!-- LEADERBOARD -->
<section id="leaderboard" class="mt-14 max-w-5xl mx-auto">
	<h2 class="text-3xl font-bold text-center mb-6">🏆 Общая таблица</h2>

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
		<table class="w-full text-left text-sm">
			<thead>
				<tr class="border-b border-surface-700">
					<th class="p-3 font-semibold text-surface-300 align-middle">Место</th>
					<th class="p-3 font-semibold text-surface-300 align-middle">Игрок</th>
					<th class="p-3 font-semibold text-surface-300 align-middle">Текущая игра</th>
					<th class="px-4 py-3 font-semibold text-surface-300 text-center align-middle">Очки</th>
					<th class="px-4 py-3 font-semibold text-surface-300 text-center align-middle">Пройдено</th>
					<th class="px-4 py-3 font-semibold text-surface-300 text-center align-middle">Дроп</th>
					<th class="px-4 py-3 font-semibold text-surface-300 text-center align-middle">Реролл</th>
				</tr>
			</thead>
			<tbody>
				{#each sortedByPoints as row, index}
					<tr
						class="border-b border-surface-800 hover:bg-surface-800/50 transition"
					>
						<td class="p-4 pr-3 align-middle">
							{#if row.player.img}
								<img
									src={row.player.img}
									alt={row.player.username}
									class="w-16 h-16 rounded-full object-cover"
									loading="lazy"
									referrerpolicy="no-referrer"
								/>
							{:else}
								<div
									class="w-16 h-16 flex items-center justify-center rounded-full text-sm font-bold"
									style={`
										background: ${placeStyles[index]?.bg ?? "transparent"};
										color: ${placeStyles[index]?.color ?? getPlayerColor(row.player.username)};
									`}
								>
									{placeStyles[index]?.label ?? index + 1}
								</div>
							{/if}
						</td>
						<td class="p-2 pl-3 align-middle">
							<a
								href={`/users/${row.player.id}`}
								class="font-semibold hover:underline"
								style={`color:${getPlayerColor(row.player.username)}`}
							>
								{row.player.username}
							</a>
							{#if row.player.description}
								<p class="text-sm text-surface-400 max-w-xs truncate whitespace-nowrap">
									{row.player.description}
								</p>
							{/if}
						</td>
						<td class="p-3 text-surface-300 align-middle">
							{row.currentGame ?? "—"}
						</td>
						<td class="px-4 py-3 text-center font-medium align-middle">{row.points}</td>
						<td class="px-4 py-3 text-center align-middle">{row.completed}</td>
						<td class="px-4 py-3 text-center align-middle">{row.dropped}</td>
						<td class="px-4 py-3 text-center align-middle">{row.rerolled}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	{/if}
</section>

<!-- CTA -->
<section class="text-center mt-14">
	<p class="text-surface-300 mb-3">
		Подробная статистика по каждому игроку — в профиле.
	</p>
	<div class="flex gap-3 justify-center flex-wrap">
		<!-- <button
			type="button"
			class="btn variant-filled-secondary"
			on:click={scrollToLeaderboard}
		>
			К таблице
		</button> -->
		<a href="/games" class="btn variant-filled-secondary">
			Шаблоны игр
		</a>
	</div>
</section>
