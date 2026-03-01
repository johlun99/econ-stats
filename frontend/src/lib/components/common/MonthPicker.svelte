<script lang="ts">
  import type { AvailableMonth } from '../../types'

  interface Props {
    months: AvailableMonth[]
    selected: string
    onSelect: (month: string) => void
  }

  let { months, selected, onSelect }: Props = $props()

  function prev() {
    const idx = months.findIndex(m => m.month === selected)
    if (idx < months.length - 1) onSelect(months[idx + 1].month)
  }

  function next() {
    const idx = months.findIndex(m => m.month === selected)
    if (idx > 0) onSelect(months[idx - 1].month)
  }

  let currentLabel = $derived(months.find(m => m.month === selected)?.label ?? selected)
  let canPrev = $derived(months.findIndex(m => m.month === selected) < months.length - 1)
  let canNext = $derived(months.findIndex(m => m.month === selected) > 0)
</script>

<div class="flex items-center gap-3">
  <button
    class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
    disabled={!canPrev}
    onclick={prev}
  >
    ←
  </button>
  <span class="text-lg font-semibold text-white min-w-48 text-center">{currentLabel}</span>
  <button
    class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
    disabled={!canNext}
    onclick={next}
  >
    →
  </button>
</div>
