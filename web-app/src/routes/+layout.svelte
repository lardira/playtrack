<script lang="ts">
	import "../app.postcss";
	import { AppShell, AppBar } from "@skeletonlabs/skeleton";

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

	import { user } from "../stores/user";
	import { onMount } from "svelte";
	import type { Player } from "../lib/types";

	let currentUser: Player | null = null;
	onMount(() => {
		user.subscribe((value) => (currentUser = value));
	});
	let players: Player[] = [
		{
			ID: "1",
			Username: "Alice",
			Color: "#f97316", // orange
		},
		{
			ID: "2",
			Username: "Bob",
			Color: "#22c55e", // green
		},
		{
			ID: "3",
			Username: "Charlie",
			Color: "#3b82f6", // blue
		},
		{
			ID: "4",
			Username: "Diana",
			Color: "#a855f7", // purple
		},
	];

	function logout() {
		user.set(null);
		window.location.href = "/";
	}
</script>

<!-- App Shell -->
<AppShell>
	<svelte:fragment slot="header">
		<!-- App Bar -->
		<AppBar>
			<svelte:fragment slot="lead">
				<strong class="text-xl uppercase">PlayTracker</strong>
			</svelte:fragment>
			<svelte:fragment slot="trail">
				{#if currentUser}
					<a
						href={`/users/${currentUser.ID}`}
						class="btn btn-sm variant-ghost-surface"
						>{currentUser.Username}</a
					>
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
							href={`/users/${player.ID}`}
							class="btn btn-sm border transition hover:scale-105"
							style={`
								border-color: ${player.Color};
								color: ${player.Color};
							`}
						>
							{player.Username}
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
