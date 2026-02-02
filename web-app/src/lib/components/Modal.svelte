<script lang="ts">
    export let isOpen = false;
    export let title = "";
    export let onClose: () => void = () => {};

    function backdropClick(e: MouseEvent) {
        if ((e.target as HTMLElement).classList.contains("modal-backdrop"))
            onClose();
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape") {
            onClose();
        }
    }
</script>

{#if isOpen}
    <div
        class="modal-backdrop fixed inset-0 bg-black/50 flex justify-center items-center z-50"
        on:click={backdropClick}
        on:keydown={handleKeydown}
        role="dialog"
        aria-modal="true"
        tabindex="-1"
    >
        <div class="bg-surface-800 rounded-xl p-6 w-full max-w-lg relative">
            {#if title}<h2 class="text-2xl text-primary-400 mb-4">
                    {title}
                </h2>{/if}
            <slot></slot>
            <button
                class="absolute top-2 right-2 btn btn-ghost btn-sm"
                on:click={onClose}>✖</button
            >
        </div>
    </div>
{/if}
