<script lang="ts">
  import { onMount } from 'svelte'
  import { GetUncategorizedMerchants, GetMerchantTransactions, GetCategories, CategorizeByMerchant } from '../../../wailsjs/go/app/App'
  import type { MerchantGroup, Category, Transaction } from '../types'
  import { hexToRgba } from '../utils'

  interface Props {
    onToast: (message: string, type: 'success' | 'error' | 'info') => void
  }

  let { onToast }: Props = $props()

  let merchants: MerchantGroup[] = $state([])
  let allCategories: Category[] = $state([])
  let loading = $state(true)
  let expandedMerchant: string | null = $state(null)
  let merchantTransactions: Transaction[] = $state([])

  function categoriesForMerchant(m: MerchantGroup): Category[] {
    const hasIncome = m.incomeTotal > 0
    const hasExpense = m.expenseTotal > 0
    if (hasIncome && !hasExpense) {
      return allCategories.filter(c => c.isIncome || c.isExpense)
    }
    if (hasExpense && !hasIncome) {
      return allCategories.filter(c => c.isExpense)
    }
    return allCategories.filter(c => c.isIncome || c.isExpense)
  }

  async function load() {
    loading = true
    try {
      const [m, c] = await Promise.all([
        GetUncategorizedMerchants(),
        GetCategories(),
      ])
      merchants = m ?? []
      allCategories = c ?? []
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
      const catName = allCategories.find(c => c.id === categoryId)?.name ?? ''
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

  function fmtDetailed(n: number): string {
    return n.toLocaleString('sv-SE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' kr'
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
        {@const cats = categoriesForMerchant(m)}
        <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
          <!-- Merchant header -->
          <button
            class="w-full flex items-center justify-between p-4 hover:bg-slate-700/50 transition-colors text-left"
            onclick={() => toggleExpand(m.merchantKey)}
          >
            <div class="flex-1">
              <div class="flex items-center gap-3">
                <span class="text-white font-medium capitalize">{m.merchantKey}</span>
                <span class="text-lg font-semibold {m.totalAmount < 0 ? 'text-red-400' : 'text-green-400'}">
                  {fmt(m.totalAmount)}
                </span>
              </div>
              <div class="flex items-center gap-4 mt-1 text-xs text-slate-400">
                <span>{m.count} transaktioner</span>
                {#if m.expenseTotal > 0}
                  <span class="text-red-400">Utgifter: {fmt(m.expenseTotal)}</span>
                {/if}
                {#if m.incomeTotal > 0}
                  <span class="text-green-400">Inkomster: {fmt(m.incomeTotal)}</span>
                {/if}
                <span>{m.firstDate} — {m.lastDate}</span>
              </div>
            </div>
            <span class="text-slate-400 text-sm">{expandedMerchant === m.merchantKey ? '▼' : '▶'}</span>
          </button>

          <!-- Category buttons -->
          {#if cats.length > 0}
            <div class="px-4 pb-3 flex flex-wrap gap-2">
              {#each cats as cat}
                <button
                  class="px-3 py-1.5 rounded-full text-xs font-medium transition-colors hover:opacity-90"
                  style="background-color: {hexToRgba(cat.color, 0.12)}; color: {cat.color}; border: 1px solid {hexToRgba(cat.color, 0.25)}"
                  onclick={() => categorize(m.merchantKey, cat.id)}
                >
                  {cat.icon} {cat.name}
                </button>
              {/each}
            </div>
          {:else}
            <div class="px-4 pb-3 text-xs text-amber-400">
              Ingen kategori stödjer både inkomst och utgift. Skapa eller redigera en under Kategorier.
            </div>
          {/if}

          <!-- Expanded transactions -->
          {#if expandedMerchant === m.merchantKey && merchantTransactions.length > 0}
            <div class="border-t border-slate-700 bg-slate-900/50">
              <table class="w-full text-sm">
                <tbody>
                  {#each merchantTransactions as t}
                    <tr class="border-b border-slate-700/50">
                      <td class="px-4 py-2 text-slate-400">{t.transactionDate}</td>
                      <td class="px-4 py-2 text-slate-300">{t.description}</td>
                      <td class="px-4 py-2 text-right {t.amount < 0 ? 'text-red-400' : 'text-green-400'}">{fmtDetailed(t.amount)}</td>
                    </tr>
                  {/each}
                </tbody>
                <tfoot>
                  {#if m.incomeTotal > 0 && m.expenseTotal > 0}
                    <tr class="border-t border-slate-600">
                      <td colspan="2" class="px-4 py-1.5 text-xs text-slate-400">Inkomster</td>
                      <td class="px-4 py-1.5 text-right text-sm text-green-400 font-medium">+{fmtDetailed(m.incomeTotal)}</td>
                    </tr>
                    <tr>
                      <td colspan="2" class="px-4 py-1.5 text-xs text-slate-400">Utgifter</td>
                      <td class="px-4 py-1.5 text-right text-sm text-red-400 font-medium">-{fmtDetailed(m.expenseTotal)}</td>
                    </tr>
                  {/if}
                  <tr class="border-t border-slate-600 bg-slate-800/50">
                    <td colspan="2" class="px-4 py-2 text-xs font-medium text-slate-300">Netto</td>
                    <td class="px-4 py-2 text-right text-sm font-bold {m.totalAmount < 0 ? 'text-red-400' : 'text-green-400'}">{fmtDetailed(m.totalAmount)}</td>
                  </tr>
                </tfoot>
              </table>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
