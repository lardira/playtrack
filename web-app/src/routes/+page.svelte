<script lang="ts">
	import { onMount } from "svelte";
	import { user } from "../stores/user";
	import type { Player } from "../lib/types";

	let currentUser: Player | null = null;

	const placeStyles = [
		{
			bg: "#facc15", // gold
			color: "#000",
			label: "🥇",
			glow: "0 0 40px rgba(250,204,21,0.35)",
		},
		{
			bg: "#e5e7eb", // silver
			color: "#000",
			label: "🥈",
			glow: "0 0 30px rgba(229,231,235,0.25)",
		},
		{
			bg: "#d97706", // bronze
			color: "#000",
			label: "🥉",
			glow: "0 0 25px rgba(217,119,6,0.25)",
		},
	];

	let leaderboard: Player[] = [
		{
			ID: "1",
			Username: "Alice",
			Score: 1840,
			Color: "#f97316",
			Description: "Speedrunner",
		},
		{
			ID: "2",
			Username: "Bob",
			Score: 1520,
			Color: "#22c55e",
			Description: "Achievement hunter",
		},
		{
			ID: "3",
			Username: "Charlie",
			Score: 1390,
			Color: "#3b82f6",
			Description: "PvP monster",
		},
		{
			ID: "4",
			Username: "Diana",
			Score: 980,
			Color: "#a855f7",
			Description: "Casual grinder",
		},
	];

	$: sorted = [...leaderboard].sort((a, b) => b.Score - a.Score);

	onMount(() => {
		user.subscribe((value) => (currentUser = value));
	});
</script>

<!-- HERO -->
<section
	class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary-600/20 via-surface-900 to-secondary-600/20 p-10 text-center shadow-xl"
>
	<h1 class="text-4xl font-extrabold tracking-tight mb-3">
		Welcome to <span class="text-primary-400">PlayTracker</span>
	</h1>
	<p class="text-lg text-surface-300 max-w-xl mx-auto">
		Отслеживай свои игры, сравнивайте статистику и доказывай, что ты лучший.
	</p>

	<!-- <div class="mt-6 flex justify-center gap-3">
		{#if currentUser}
			<a
				href={`/users/${currentUser.ID}`}
				class="btn variant-filled-primary">Go to my profile</a
			>
		{:else}
			<a href="/login" class="btn variant-filled-primary"
				>Login to start tracking</a
			>
		{/if}
		<a href="#leaderboard" class="btn variant-ghost-surface"
			>View leaderboard</a
		>
	</div> -->

	<!-- glow -->
	<div class="absolute inset-0 -z-10 blur-[120px] bg-primary-500/20"></div>
</section>

<!-- TOP PLAYERS -->
<section class="mt-14 max-w-5xl mx-auto">
	<h2 class="text-3xl font-bold text-center mb-8">Лучшие игроки</h2>

	<div class="grid gap-6 md:grid-cols-3">
		{#each sorted.slice(0, 3) as player, index}
			<a
				href={`/users/${player.ID}`}
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
					style={`color:${player.Color}`}
				>
					{player.Username}
				</p>

				<p class="text-sm text-surface-400 mb-4">
					{player.Description}
				</p>

				<div class="text-3xl font-bold">
					{player.Score}
					<span class="text-sm font-normal text-surface-400">
						points
					</span>
				</div>
			</a>
		{/each}
	</div>
</section>

<!-- LEADERBOARD -->
<section class="mt-14 max-w-5xl mx-auto">
	<h2 class="text-3xl font-bold text-center mb-6">🏆 Общая таблица</h2>

	<div class="space-y-3">
		{#each sorted as player, index}
			<a
				href={`/users/${player.ID}`}
				class="group block rounded-xl p-4 bg-surface shadow-md transition hover:shadow-xl"
				style={`border-left: 6px solid ${player.Color}`}
			>
				<div class="flex items-center gap-4">
					<!-- PLACE -->
					<div
						class="w-10 h-10 flex items-center justify-center rounded-full text-lg font-bold"
						style={`
						background: ${placeStyles[index]?.bg ?? "transparent"};
						color: ${placeStyles[index]?.color ?? player.Color};
					`}
					>
						{placeStyles[index]?.label ?? index + 1}
					</div>

					<!-- INFO -->
					<div class="flex-1">
						<p
							class="text-lg font-semibold group-hover:underline"
							style={`color:${player.Color}`}
						>
							{player.Username}
						</p>
						<p class="text-sm text-surface-400">
							{player.Description}
						</p>
					</div>

					<!-- SCORE -->
					<div class="text-right">
						<p class="text-xl font-bold">
							{player.Score}
						</p>
						<p class="text-xs text-surface-400">очки</p>
					</div>
				</div>
			</a>
		{/each}
	</div>
</section>

<!-- CTA -->
<section class="text-center mt-14">
	<p class="text-surface-300 mb-3">
		Want to see detailed stats for each player?
	</p>
	<a href="/players" class="btn variant-filled-secondary">
		Browse all players
	</a>
</section>
