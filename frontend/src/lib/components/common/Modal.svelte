<script lang="ts">
  import type { Snippet } from 'svelte'

  interface Props {
    title: string
    onClose: () => void
    children: Snippet
  }

  let { title, onClose, children }: Props = $props()
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="fixed inset-0 z-50 flex items-center justify-center">
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div class="absolute inset-0 bg-black/60" onclick={onClose}></div>
  <div class="relative bg-slate-800 rounded-xl border border-slate-700 shadow-xl max-w-lg w-full mx-4 max-h-[80vh] flex flex-col">
    <div class="flex items-center justify-between p-4 border-b border-slate-700">
      <h2 class="text-lg font-semibold text-white">{title}</h2>
      <button class="text-slate-400 hover:text-white" onclick={onClose}>✕</button>
    </div>
    <div class="p-4 overflow-y-auto">
      {@render children()}
    </div>
  </div>
</div>
