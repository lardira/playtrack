<script lang="ts">
	import { onMount } from "svelte";
	import Layout from "./routes/layout.svelte";
	import Home from "./routes/page.svelte";
	import Login from "./routes/login/page.svelte";
	import UserPage from "./routes/users/[id]/page.svelte";

	import {
		getLeaderboard,
		getPlayers,
		getPlayer,
		getPlayerGames,
		getGame,
	} from "$lib/api";
	import { user } from "$lib/stores/user";
	import { get } from "svelte/store";

	let currentPath = "/";
	let currentComponent: any = null;
	let componentData: any = null;

	let loading = false;
	let error: string | null = null;

	const homeCache: any = {};
	const userCache: Record<string, any> = {};

	async function navigate(path: string) {
		if (path === currentPath) return;
		history.pushState({}, "", path);
		currentPath = path;
		await handleRoute();
	}

	async function handleRoute() {
		loading = true;
		error = null;

		try {
			if (currentPath === "/") {
				if (!homeCache.data) {
					const [leaderboard, players] = await Promise.all([
						getLeaderboard(),
						getPlayers(),
					]);
					homeCache.data = { leaderboard, players };
				}
				currentComponent = Home;
				componentData = homeCache.data;
			} else if (currentPath === "/login") {
				currentComponent = Login;
				componentData = {};
			} else if (currentPath.startsWith("/users/")) {
				const id = currentPath.split("/")[2];
				if (!userCache[id]) {
					const player = await getPlayer(id);
					const gamesRaw = await getPlayerGames(id);

					const games = await Promise.all(
						gamesRaw.map(async (g) => {
							const game = await getGame(g.game_id);
							return {
								...g,
								gameTitle: game.Title,
								gameURL: game.URL,
							};
						}),
					);

					const currentUser = get(user);
					const isOwner = currentUser?.ID === id;

					userCache[id] = { player, games, isOwner };
				}
				currentComponent = UserPage;
				componentData = userCache[id];
			} else {
				// 404 - redirect to home
				navigate("/");
				return;
			}
		} catch (e: any) {
			console.error("Route error", e);
			error = e.message || "Ошибка загрузки";
		} finally {
			loading = false;
		}
	}

	function handleClick(e: MouseEvent) {
		const link = (e.target as HTMLElement).closest("a");
		if (link && link.hostname === window.location.hostname) {
			e.preventDefault();
			navigate(link.pathname);
		}
	}

	onMount(() => {
		currentPath = window.location.pathname;
		handleRoute();

		window.addEventListener("popstate", () => {
			currentPath = window.location.pathname;
			handleRoute();
		});

		document.addEventListener("click", handleClick);

		return () => {
			window.removeEventListener("popstate", handleRoute);
			document.removeEventListener("click", handleClick);
		};
	});
</script>

<Layout>
	{#if loading && !currentComponent}
		<div class="flex items-center justify-center min-h-[400px]">
			<div class="text-center">
				<div
					class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500 mb-4"
				></div>
				<p class="text-surface-500">Загрузка...</p>
			</div>
		</div>
	{:else if error}
		<div class="card p-4 bg-error-500/20 border-error-500/50">
			<p class="text-error-500">Ошибка: {error}</p>
		</div>
	{:else if currentComponent}
		<svelte:component this={currentComponent} data={componentData} />
	{/if}
</Layout>
