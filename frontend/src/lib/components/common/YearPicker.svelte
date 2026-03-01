<script lang="ts">
  import type { AvailableYear } from '../../types'

  interface Props {
    years: AvailableYear[]
    selected: string
    onSelect: (year: string) => void
  }

  let { years, selected, onSelect }: Props = $props()

  function prev() {
    const idx = years.findIndex(y => y.year === selected)
    if (idx < years.length - 1) onSelect(years[idx + 1].year)
  }

  function next() {
    const idx = years.findIndex(y => y.year === selected)
    if (idx > 0) onSelect(years[idx - 1].year)
  }

  let currentLabel = $derived(years.find(y => y.year === selected)?.label ?? selected)
  let canPrev = $derived(years.findIndex(y => y.year === selected) < years.length - 1)
  let canNext = $derived(years.findIndex(y => y.year === selected) > 0)
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
