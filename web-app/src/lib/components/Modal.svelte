<script lang="ts">
    import { onMount, onDestroy } from "svelte";

    export let isOpen = false;
    export let title: string | null = null;
    export let onClose: (() => void) | null = null;

    function close() {
        if (onClose) {
            onClose();
        } else {
            isOpen = false;
        }
    }

    function handleKey(e: KeyboardEvent) {
        if (e.key === "Escape") {
            close();
        }
    }

    function handleBackdropClick(e: MouseEvent | KeyboardEvent) {
        if (e.target === e.currentTarget) {
            close();
        }
    }

    onMount(() => window.addEventListener("keydown", handleKey));
    onDestroy(() => window.removeEventListener("keydown", handleKey));
</script>

{#if isOpen}
    <div
        class="fixed inset-0 bg-surface-900/80 backdrop-blur-sm flex items-center justify-center z-50"
        on:click={handleBackdropClick}
        on:keydown={(e) => e.key === "Enter" && handleBackdropClick(e)}
        role="dialog"
        aria-modal="true"
        tabindex="-1"
    >
        <div
            class="card max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
            on:click|stopPropagation
            on:keydown|stopPropagation
        >
            {#if title}
                <div class="flex items-center justify-between mb-4">
                    <h3 class="text-2xl font-bold text-primary-400">{title}</h3>
                    <button
                        on:click={close}
                        class="btn variant-ghost-surface"
                        aria-label="Закрыть"
                    >
                        ✕
                    </button>
                </div>
            {/if}
            <slot />
        </div>
    </div>
{/if}
