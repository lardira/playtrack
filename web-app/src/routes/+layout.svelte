<script lang="ts">
	import "../app.postcss";
	import { AppShell, AppBar } from "@skeletonlabs/skeleton";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import { browser } from "$app/environment";

	// Floating UI for Popups
	import {
		computePosition,
		autoUpdate,
		flip,
		shift,
		offset,
		arrow,
	} from "@floating-ui/dom";
	import { storePopup } from "@skeletonlabs/skeleton";
	storePopup.set({ computePosition, autoUpdate, flip, shift, offset, arrow });

	import { user, token } from "../stores/user";
	import { onMount } from "svelte";
	import type { Player } from "../lib/types";
	import { getPlayers } from "../lib/api";

	let currentUser: Player | null = null;
	let players: Player[] = [];

	// Редирект на логин, если пользователь не залогинен
	$: if (browser && $token === null && $page.url.pathname !== "/login") {
		goto("/login");
	}

	onMount(() => {
		user.subscribe((value) => (currentUser = value));
		// Загружаем список игроков с API при наличии токена
		if ($token) {
			getPlayers()
				.then((list) => (players = list))
				.catch(() => (players = []));
		}
		token.subscribe((t) => {
			if (t) {
				getPlayers()
					.then((list) => (players = list))
					.catch(() => (players = []));
			} else {
				players = [];
			}
		});
	});

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

	function logout() {
		user.set(null);
		token.set(null);
		goto("/login");
	}
</script>

<!-- App Shell -->
<AppShell>
	<svelte:fragment slot="header">
		<!-- App Bar -->
		<AppBar>
			<svelte:fragment slot="lead">
				<a
					href="/"
					class="cursor-pointer hover:opacity-80 transition-opacity"
				>
					<strong class="text-xl uppercase">PlayTrack</strong>
				</a>
			</svelte:fragment>
			<svelte:fragment slot="trail">
				{#if currentUser}
					<button
						class="btn btn-sm variant-ghost-surface"
						on:click={logout}>Выйти</button
					>
				{:else}
					<a href="/login" class="btn btn-sm variant-filled-primary"
						>Войти</a
					>
				{/if}

				<!-- PLAYER BUTTONS -->
				<div class="flex items-center gap-2 mr-3">
					{#each players as player}
						<a
							href={`/users/${player.id}`}
							class="btn btn-sm border transition hover:scale-105"
							style={`
								border-color: ${getPlayerColor(player.username)};
								color: ${getPlayerColor(player.username)};
							`}
						>
							{player.username}
						</a>
					{/each}
				</div>

				<!-- Ссылки на соцсети -->
				<a
					class="btn btn-sm variant-ghost-surface"
					href="https://github.com/lardira/playtrack/tree/master"
					target="_blank"
					rel="noreferrer"
				>
					GitHub
				</a>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>

	<!-- Page Route Content -->
	<slot />
</AppShell>
