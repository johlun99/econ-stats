<script lang="ts">
  import { SelectAndImportFile } from '../../../wailsjs/go/app/App'
  import type { ImportResult } from '../types'

  interface Props {
    onComplete: () => void
  }

  let { onComplete }: Props = $props()

  let result: ImportResult | null = $state(null)
  let loading = $state(false)
  let error = $state('')

  async function handleImport() {
    loading = true
    error = ''
    result = null
    try {
      const r = await SelectAndImportFile()
      if (r) {
        result = r
      }
    } catch (e: any) {
      error = e?.message || String(e)
    } finally {
      loading = false
    }
  }
</script>

<div class="p-6 max-w-2xl mx-auto">
  <h2 class="text-2xl font-bold text-white mb-6">Importera transaktioner</h2>

  <!-- Drop zone / button -->
  <button
    class="w-full border-2 border-dashed border-slate-600 rounded-xl p-12 text-center
           hover:border-blue-500 hover:bg-slate-800/50 transition-colors cursor-pointer
           disabled:opacity-50 disabled:cursor-not-allowed"
    disabled={loading}
    onclick={handleImport}
  >
    {#if loading}
      <div class="text-4xl mb-3 animate-spin">⏳</div>
      <p class="text-slate-300">Importerar...</p>
    {:else}
      <div class="text-4xl mb-3">📁</div>
      <p class="text-slate-300 text-lg">Klicka för att välja Handelsbanken Excel-fil</p>
      <p class="text-slate-500 text-sm mt-2">.xlsx-filer stöds</p>
    {/if}
  </button>

  <!-- Error -->
  {#if error}
    <div class="mt-4 p-4 bg-red-900/30 border border-red-700 rounded-lg text-red-300 text-sm">
      {error}
    </div>
  {/if}

  <!-- Result -->
  {#if result}
    <div class="mt-6 bg-slate-800 rounded-xl p-6 border border-slate-700 space-y-4">
      <h3 class="text-lg font-semibold text-white">Importresultat</h3>

      <div class="grid grid-cols-2 gap-4">
        <div class="bg-slate-700/50 rounded-lg p-3">
          <div class="text-2xl font-bold text-green-400">{result.newTransactions}</div>
          <div class="text-xs text-slate-400">Nya transaktioner</div>
        </div>
        {#if result.updated > 0}
        <div class="bg-slate-700/50 rounded-lg p-3">
          <div class="text-2xl font-bold text-cyan-400">{result.updated}</div>
          <div class="text-xs text-slate-400">Uppdaterade</div>
        </div>
        {/if}
        <div class="bg-slate-700/50 rounded-lg p-3">
          <div class="text-2xl font-bold text-slate-400">{result.duplicatesSkipped}</div>
          <div class="text-xs text-slate-400">Dubbletter (hoppades över)</div>
        </div>
        <div class="bg-slate-700/50 rounded-lg p-3">
          <div class="text-2xl font-bold text-blue-400">{result.autoCategorized}</div>
          <div class="text-xs text-slate-400">Auto-kategoriserade</div>
        </div>
        <div class="bg-slate-700/50 rounded-lg p-3">
          <div class="text-2xl font-bold text-amber-400">{result.uncategorized}</div>
          <div class="text-xs text-slate-400">Okategoriserade</div>
        </div>
      </div>

      {#if result.uncategorized > 0}
        <div class="flex items-center gap-3 p-3 bg-amber-900/20 border border-amber-700/50 rounded-lg">
          <span class="text-amber-400 text-lg">⚠️</span>
          <div>
            <p class="text-amber-300 text-sm font-medium">
              {result.uncategorized} transaktioner behöver kategoriseras
            </p>
            <button class="text-amber-400 text-xs hover:underline mt-1" onclick={onComplete}>
              Gå till kategorisering →
            </button>
          </div>
        </div>
      {:else if result.newTransactions > 0}
        <button
          class="w-full py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm transition-colors"
          onclick={onComplete}
        >
          Klart! Visa dashboard →
        </button>
      {/if}
    </div>
  {/if}
</div>
