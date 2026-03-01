<script lang="ts">
  import { onMount } from 'svelte'
  import { GetUncategorizedMerchants, GetMerchantTransactions, GetCategories, CategorizeByMerchant } from '../../../wailsjs/go/app/App'
  import type { MerchantGroup, Category, Transaction } from '../types'

  interface Props {
    onToast: (message: string, type: 'success' | 'error' | 'info') => void
  }

  let { onToast }: Props = $props()

  let merchants: MerchantGroup[] = $state([])
  let categories: Category[] = $state([])
  let loading = $state(true)
  let expandedMerchant: string | null = $state(null)
  let merchantTransactions: Transaction[] = $state([])

  async function load() {
    loading = true
    try {
      const [m, c] = await Promise.all([
        GetUncategorizedMerchants(),
        GetCategories(),
      ])
      merchants = m ?? []
      categories = (c ?? []).filter(c => c.isExpense)
    } finally {
      loading = false
    }
  }

  async function toggleExpand(merchantKey: string) {
    if (expandedMerchant === merchantKey) {
      expandedMerchant = null
      merchantTransactions = []
      return
    }
    expandedMerchant = merchantKey
    merchantTransactions = await GetMerchantTransactions(merchantKey) ?? []
  }

  async function categorize(merchantKey: string, categoryId: number) {
    try {
      const count = await CategorizeByMerchant(merchantKey, categoryId)
      const catName = categories.find(c => c.id === categoryId)?.name ?? ''
      onToast(`${count} transaktioner kategoriserade som "${catName}"`, 'success')
      merchants = merchants.filter(m => m.merchantKey !== merchantKey)
      if (expandedMerchant === merchantKey) expandedMerchant = null
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  function fmt(n: number): string {
    return n.toLocaleString('sv-SE', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) + ' kr'
  }

  onMount(load)
</script>

<div class="p-6">
  <div class="flex items-center justify-between mb-6">
    <div>
      <h2 class="text-2xl font-bold text-white">Kategorisera</h2>
      <p class="text-sm text-slate-400 mt-1">
        {merchants.length} handlare behöver kategoriseras
      </p>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center h-64 text-slate-400">Laddar...</div>
  {:else if merchants.length === 0}
    <div class="flex flex-col items-center justify-center h-64 text-slate-400">
      <span class="text-4xl mb-3">✅</span>
      <p class="text-lg">Alla transaktioner är kategoriserade!</p>
    </div>
  {:else}
    <div class="space-y-3">
      {#each merchants as m}
        <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
          <!-- Merchant header -->
          <button
            class="w-full flex items-center justify-between p-4 hover:bg-slate-700/50 transition-colors text-left"
            onclick={() => toggleExpand(m.merchantKey)}
          >
            <div class="flex-1">
              <span class="text-white font-medium capitalize">{m.merchantKey}</span>
              <div class="flex items-center gap-4 mt-1 text-xs text-slate-400">
                <span>{m.count} transaktioner</span>
                <span>Totalt: {fmt(m.totalAmount)}</span>
                <span>{m.firstDate} — {m.lastDate}</span>
              </div>
            </div>
            <span class="text-slate-400 text-sm">{expandedMerchant === m.merchantKey ? '▼' : '▶'}</span>
          </button>

          <!-- Category buttons -->
          <div class="px-4 pb-3 flex flex-wrap gap-2">
            {#each categories as cat}
              <button
                class="px-3 py-1.5 rounded-full text-xs font-medium transition-colors hover:opacity-90"
                style="background-color: {cat.color}20; color: {cat.color}; border: 1px solid {cat.color}40"
                onclick={() => categorize(m.merchantKey, cat.id)}
              >
                {cat.icon} {cat.name}
              </button>
            {/each}
          </div>

          <!-- Expanded transactions -->
          {#if expandedMerchant === m.merchantKey && merchantTransactions.length > 0}
            <div class="border-t border-slate-700 bg-slate-900/50">
              <table class="w-full text-sm">
                <tbody>
                  {#each merchantTransactions as t}
                    <tr class="border-b border-slate-700/50">
                      <td class="px-4 py-2 text-slate-400">{t.transactionDate}</td>
                      <td class="px-4 py-2 text-slate-300">{t.description}</td>
                      <td class="px-4 py-2 text-right text-red-400">{fmt(t.amount)}</td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
