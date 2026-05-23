<script lang="ts">
	import "../app.postcss";
	import { AppShell, AppBar } from "@skeletonlabs/skeleton";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import { browser } from "$app/environment";
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
	import { getPlayerColorBorderTextStyle } from "../lib/colorUtils";

	const AVAILABLE_THEMES = [
		"crimson",
		"gold-nouveau",
		"hamlindigo",
		"modern",
		"rocket",
		"sahara",
		"seafoam",
		"skeleton",
		"vintage",
		"wintry",
	];

	/* theme id -> emoji: 🔴 crimson, 👑 gold-nouveau, 👔 hamlindigo, 🌸 modern, 🚀 rocket, 🏜️ sahara, 🧜‍♀️ seafoam, ⚙ skeleton, 📺 vintage, 🌨️ wintry */
	const THEME_EMOJI: Record<string, string> = {
		crimson: "🔴",
		"gold-nouveau": "👑",
		hamlindigo: "👔",
		modern: "🌸",
		rocket: "🚀",
		sahara: "🏜️",
		seafoam: "🧜‍♀️",
		skeleton: "⚙",
		vintage: "📺",
		wintry: "🌨️",
	};

	let currentUser: Player | null = null;
	let players: Player[] = [];
	let currentTheme: string = AVAILABLE_THEMES[0];
	let themePickerOpen = false;
	let playersFetchGen = 0;

	function getThemeStorageKey(u: Player | null): string {
		return u ? `playtrack-theme-${u.id}` : "playtrack-theme-guest";
	}

	function applyTheme(theme: string) {
		if (!browser) return;
		const safeTheme = AVAILABLE_THEMES.includes(theme) ? theme : AVAILABLE_THEMES[0];
		currentTheme = safeTheme;
		document.documentElement.setAttribute("data-theme", safeTheme);
		document.body?.setAttribute("data-theme", safeTheme);
		const key = getThemeStorageKey(currentUser);
		localStorage.setItem(key, safeTheme);
	}

	function loadThemeForUser(u: Player | null) {
		if (!browser) return;
		const key = getThemeStorageKey(u);
		const saved = localStorage.getItem(key);
		applyTheme(saved && AVAILABLE_THEMES.includes(saved) ? saved : AVAILABLE_THEMES[0]);
	}

	function setThemeAndClose(theme: string) {
		applyTheme(theme);
		themePickerOpen = false;
	}

	$: if (browser && $token === null && $page.url.pathname !== "/login") {
		goto("/login");
	}

	function fetchPlayersForHeader() {
		if (!$token) {
			players = [];
			return;
		}
		const gen = ++playersFetchGen;
		getPlayers()
			.then((list) => {
				if (gen === playersFetchGen) players = list;
			})
			.catch(() => {
				if (gen === playersFetchGen) players = [];
			});
	}

	onMount(() => {
		const unsubUser = user.subscribe((value) => {
			currentUser = value;
			loadThemeForUser(currentUser);
		});
		fetchPlayersForHeader();
		const unsubToken = token.subscribe((t) => {
			if (!t) {
				players = [];
				playersFetchGen++;
				return;
			}
			fetchPlayersForHeader();
		});
		if (currentUser === null) {
			loadThemeForUser(null);
		}
		return () => {
			unsubUser();
			unsubToken();
		};
	});

	function logout() {
		user.set(null);
		token.set(null);
		goto("/login");
	}
</script>

<AppShell>
	<svelte:fragment slot="header">
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
				<div class="flex flex-col-reverse md:flex-row md:items-center gap-2 md:gap-4">
					<div class="flex items-center gap-2 md:mr-3 overflow-x-auto max-w-full md:max-w-xs lg:max-w-none scrollbar-hide">
						{#each players as player}
							<a
								href={`/users/${player.id}`}
								class="player-pill btn btn-sm border-2 transition hover:scale-105 whitespace-nowrap flex-shrink-0"
								style={getPlayerColorBorderTextStyle(player.username)}
							>
								{player.username}
							</a>
						{/each}
					</div>

					<div class="flex items-center gap-2 justify-end md:justify-start relative">
						<button
							type="button"
							class="btn btn-sm variant-ghost-surface w-9 h-9 p-0 flex items-center justify-center text-lg rounded-lg border border-surface-600 hover:border-surface-500"
							title="Тема"
							on:click={() => (themePickerOpen = !themePickerOpen)}
							on:keydown={(e) => e.key === 'Escape' && (themePickerOpen = false)}
						>
							{THEME_EMOJI[currentTheme] ?? "📺"}
						</button>
						{#if themePickerOpen}
							<button
								type="button"
								class="fixed inset-0 z-10"
								aria-label="Закрыть"
								on:click={() => (themePickerOpen = false)}
							/>
							<div
								class="absolute right-0 top-full mt-1 z-20 p-2 rounded-xl bg-surface-800 border border-surface-600 shadow-xl grid grid-cols-5 gap-1"
								role="listbox"
								aria-label="Выбор темы"
							>
								{#each AVAILABLE_THEMES as theme}
									<button
										type="button"
										class="w-8 h-8 flex items-center justify-center text-lg rounded-lg hover:bg-surface-700 transition {currentTheme === theme ? 'ring-2 ring-primary-500 bg-surface-700' : ''}"
										title={theme}
										role="option"
										aria-selected={currentTheme === theme}
										on:click={() => setThemeAndClose(theme)}
									>
										{THEME_EMOJI[theme] ?? "⚙"}
									</button>
								{/each}
							</div>
						{/if}

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
					</div>
				</div>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>

	<slot />
</AppShell>
