<script lang="ts">
  import { onMount } from 'svelte'
  import { GetTransactions, SearchTransactions, GetAvailableMonths } from '../../../wailsjs/go/app/App'
  import MonthPicker from '../components/common/MonthPicker.svelte'
  import type { Transaction, AvailableMonth } from '../types'

  let transactions: Transaction[] = $state([])
  let months: AvailableMonth[] = $state([])
  let selectedMonth = $state('')
  let searchTerm = $state('')
  let loading = $state(true)
  let filterType: 'all' | 'expenses' | 'income' | 'transfers' = $state('all')

  let filteredTransactions = $derived.by(() => {
    let list = transactions
    if (filterType === 'expenses') list = list.filter(t => t.amount < 0 && !t.isTransfer)
    else if (filterType === 'income') list = list.filter(t => t.amount > 0 && !t.isTransfer)
    else if (filterType === 'transfers') list = list.filter(t => t.isTransfer)
    return list
  })

  function fmt(n: number): string {
    return n.toLocaleString('sv-SE', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + ' kr'
  }

  async function load() {
    loading = true
    try {
      months = (await GetAvailableMonths()) ?? []
      if (months.length > 0 && !selectedMonth) {
        selectedMonth = months[0].month
      }
      await fetchTransactions()
    } finally {
      loading = false
    }
  }

  async function fetchTransactions() {
    if (searchTerm.trim()) {
      transactions = (await SearchTransactions(searchTerm, selectedMonth)) ?? []
    } else {
      transactions = (await GetTransactions(selectedMonth)) ?? []
    }
  }

  async function selectMonth(month: string) {
    selectedMonth = month
    await fetchTransactions()
  }

  let searchTimeout: ReturnType<typeof setTimeout>
  function handleSearch() {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(fetchTransactions, 300)
  }

  onMount(load)
</script>

<div class="p-6">
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-2xl font-bold text-white">Transaktioner</h2>
    {#if months.length > 0}
      <MonthPicker {months} selected={selectedMonth} onSelect={selectMonth} />
    {/if}
  </div>

  <!-- Search & filters -->
  <div class="flex items-center gap-4 mb-4">
    <input
      type="text"
      bind:value={searchTerm}
      oninput={handleSearch}
      class="flex-1 px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white text-sm
             focus:outline-none focus:border-blue-500 placeholder-slate-500"
      placeholder="Sök transaktioner..."
    />
    <div class="flex rounded-lg overflow-hidden border border-slate-700">
      {#each [
        { key: 'all', label: 'Alla' },
        { key: 'expenses', label: 'Utgifter' },
        { key: 'income', label: 'Inkomster' },
        { key: 'transfers', label: 'Överföringar' },
      ] as filter}
        <button
          class="px-3 py-2 text-xs transition-colors
            {filterType === filter.key ? 'bg-blue-600 text-white' : 'bg-slate-800 text-slate-400 hover:text-white'}"
          onclick={() => filterType = filter.key as typeof filterType}
        >
          {filter.label}
        </button>
      {/each}
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center h-64 text-slate-400">Laddar...</div>
  {:else if filteredTransactions.length === 0}
    <div class="flex items-center justify-center h-64 text-slate-400">
      Inga transaktioner att visa
    </div>
  {:else}
    <div class="text-xs text-slate-500 mb-2">{filteredTransactions.length} transaktioner</div>
    <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-700/50">
          <tr>
            <th class="px-4 py-3 text-left text-slate-400 font-medium">Datum</th>
            <th class="px-4 py-3 text-left text-slate-400 font-medium">Beskrivning</th>
            <th class="px-4 py-3 text-left text-slate-400 font-medium">Kategori</th>
            <th class="px-4 py-3 text-right text-slate-400 font-medium">Belopp</th>
            <th class="px-4 py-3 text-right text-slate-400 font-medium">Saldo</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredTransactions as t}
            <tr class="border-t border-slate-700/50 hover:bg-slate-750">
              <td class="px-4 py-2.5 text-slate-400">{t.transactionDate}</td>
              <td class="px-4 py-2.5 text-slate-200">{t.description}</td>
              <td class="px-4 py-2.5">
                {#if t.categoryName}
                  <span class="px-2 py-0.5 rounded-full text-xs" style="background-color: {t.categoryColor}20; color: {t.categoryColor}">
                    {t.categoryIcon} {t.categoryName}
                  </span>
                {:else if t.isTransfer}
                  <span class="text-xs text-slate-500">Överföring</span>
                {:else}
                  <span class="text-xs text-amber-500">Okategoriserad</span>
                {/if}
              </td>
              <td class="px-4 py-2.5 text-right font-medium
                {t.amount < 0 ? 'text-red-400' : 'text-green-400'}">
                {fmt(t.amount)}
              </td>
              <td class="px-4 py-2.5 text-right text-slate-400">{fmt(t.balance)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
