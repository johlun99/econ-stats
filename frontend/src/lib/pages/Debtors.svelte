<script lang="ts">
  import { onMount } from 'svelte'
  import {
    GetDebtors,
    CreateDebtor,
    UpdateDebtor,
    DeleteDebtor,
    AddDebtorMerchantKey,
    RemoveDebtorMerchantKey,
    GetDebtorTransactions,
    GetAllMerchantKeys,
    AddDebtorTransaction,
    UpdateDebtorTransaction,
    DeleteDebtorTransaction
  } from '../../../wailsjs/go/app/App'
  import Modal from '../components/common/Modal.svelte'
  import type { DebtorDetail, Transaction } from '../types'
  import { hexToRgba } from '../utils'

  interface Props {
    onToast: (message: string, type: 'success' | 'error' | 'info') => void
  }

  let { onToast }: Props = $props()

  let debtors: DebtorDetail[] = $state([])
  let allMerchantKeys: string[] = $state([])
  let showModal = $state(false)
  let editingDebtor: DebtorDetail | null = $state(null)
  let expandedId: number | null = $state(null)
  let transactions: Transaction[] = $state([])
  let addingKeyForId: number | null = $state(null)
  let merchantKeySearch = $state('')

  // Manual transaction form state
  let addingTxForId: number | null = $state(null)
  let editingTxId: number | null = $state(null)
  let txFormDate = $state('')
  let txFormDescription = $state('')
  let txFormAmount: number | null = $state(null)

  // Form state
  let formName = $state('')
  let formIcon = $state('👤')
  let formColor = $state('#6B7280')

  const defaultColors = [
    '#EF4444', '#F97316', '#F59E0B', '#EAB308',
    '#84CC16', '#22C55E', '#10B981', '#14B8A6',
    '#06B6D4', '#3B82F6', '#6366F1', '#8B5CF6',
    '#A855F7', '#D946EF', '#EC4899', '#F43F5E',
    '#64748B', '#6B7280', '#78716C', '#FFFFFF',
  ]

  const defaultIcons = ['👤', '👫', '👨', '👩', '🏠', '🏢', '💼', '🤝', '👨‍👩‍👧', '👪', '🧑‍🤝‍🧑', '💰']

  // All merchant keys already linked to any debtor
  let linkedKeys = $derived(new Set(debtors.flatMap(d => d.merchantKeys ?? [])))

  // Unlinked keys available for adding
  let availableKeys = $derived(
    allMerchantKeys.filter(k => !linkedKeys.has(k))
  )

  // Filtered keys for the search dropdown
  let filteredKeys = $derived(
    merchantKeySearch
      ? availableKeys.filter(k => k.toLowerCase().includes(merchantKeySearch.toLowerCase()))
      : availableKeys
  )

  async function load() {
    const [d, keys] = await Promise.all([GetDebtors(), GetAllMerchantKeys()])
    debtors = d ?? []
    allMerchantKeys = keys ?? []
  }

  function openCreate() {
    editingDebtor = null
    formName = ''
    formIcon = '👤'
    formColor = '#6B7280'
    showModal = true
  }

  function openEdit(d: DebtorDetail) {
    editingDebtor = d
    formName = d.name
    formIcon = d.icon
    formColor = d.color
    showModal = true
  }

  async function handleSave() {
    if (!formName.trim()) return
    try {
      if (editingDebtor) {
        await UpdateDebtor(editingDebtor.id, formName, formIcon, formColor)
        onToast('Uppdaterad', 'success')
      } else {
        await CreateDebtor(formName, formIcon, formColor)
        onToast('Skapad', 'success')
      }
      showModal = false
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleDelete(id: number) {
    try {
      await DeleteDebtor(id)
      onToast('Borttagen', 'success')
      if (expandedId === id) expandedId = null
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleAddKey(debtorId: number, key: string) {
    try {
      await AddDebtorMerchantKey(debtorId, key)
      addingKeyForId = null
      merchantKeySearch = ''
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleRemoveKey(debtorId: number, key: string) {
    try {
      await RemoveDebtorMerchantKey(debtorId, key)
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function toggleTransactions(debtorId: number) {
    if (expandedId === debtorId) {
      expandedId = null
      transactions = []
      return
    }
    expandedId = debtorId
    try {
      const result = await GetDebtorTransactions(debtorId)
      transactions = result ?? []
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  function openAddTx(debtorId: number) {
    addingTxForId = debtorId
    editingTxId = null
    txFormDate = new Date().toISOString().slice(0, 10)
    txFormDescription = ''
    txFormAmount = null
  }

  function openEditTx(t: Transaction) {
    editingTxId = t.id
    addingTxForId = null
    txFormDate = t.transactionDate
    txFormDescription = t.description
    txFormAmount = t.amount
  }

  function cancelTxForm() {
    addingTxForId = null
    editingTxId = null
  }

  async function handleSaveTx(debtorId: number) {
    if (!txFormDate || !txFormDescription.trim() || txFormAmount == null || isNaN(txFormAmount)) return
    const amount = txFormAmount
    try {
      if (editingTxId != null) {
        await UpdateDebtorTransaction(editingTxId, txFormDescription, amount, txFormDate)
        onToast('Transaktion uppdaterad', 'success')
      } else {
        await AddDebtorTransaction(debtorId, txFormDescription, amount, txFormDate)
        onToast('Transaktion tillagd', 'success')
      }
      cancelTxForm()
      const [result] = await Promise.all([GetDebtorTransactions(debtorId), load()])
      transactions = result ?? []
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleDeleteTx(debtorId: number, txId: number) {
    try {
      await DeleteDebtorTransaction(txId)
      onToast('Transaktion borttagen', 'success')
      const [result] = await Promise.all([GetDebtorTransactions(debtorId), load()])
      transactions = result ?? []
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  function formatAmount(amount: number): string {
    return amount.toLocaleString('sv-SE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }

  onMount(load)
</script>

<div class="p-6">
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-2xl font-bold text-white">Skulder</h2>
    <button
      class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm transition-colors"
      onclick={openCreate}
    >
      + Ny person
    </button>
  </div>

  {#if debtors.length === 0}
    <div class="text-center py-16 text-slate-400">
      <p class="text-4xl mb-4">🤝</p>
      <p class="text-lg">Inga personer tillagda</p>
      <p class="text-sm mt-1">Skapa en person och koppla handlare for att spara skulder</p>
    </div>
  {:else}
    <div class="space-y-4">
      {#each debtors as d}
        <div class="bg-slate-800 rounded-xl border border-slate-700">
          <!-- Header -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg"
                   style="background-color: {hexToRgba(d.color, 0.15)}">
                {d.icon}
              </div>
              <div>
                <div class="text-white font-medium">{d.name}</div>
                <div class="text-xs text-slate-400">
                  {(d.merchantKeys ?? []).length} handlare
                </div>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <div class="text-right">
                <div class="text-xs text-slate-400">Balans</div>
                <div class="font-semibold {d.balance > 0 ? 'text-green-400' : d.balance < 0 ? 'text-red-400' : 'text-slate-400'}">
                  {d.balance > 0 ? '+' : ''}{formatAmount(d.balance)} kr
                </div>
                <div class="text-xs {d.balance > 0 ? 'text-green-400/70' : d.balance < 0 ? 'text-red-400/70' : 'text-slate-500'}">
                  {d.balance > 0 ? 'De ar skyldiga dig' : d.balance < 0 ? 'Du ar skyldig dem' : 'Kvitt'}
                </div>
              </div>
              <div class="flex gap-1">
                <button class="p-2 text-slate-400 hover:text-white rounded" onclick={() => openEdit(d)}>✏️</button>
                <button class="p-2 text-slate-400 hover:text-red-400 rounded" onclick={() => handleDelete(d.id)}>🗑️</button>
              </div>
            </div>
          </div>

          <!-- Merchant keys -->
          <div class="px-4 pb-3 flex flex-wrap gap-2 items-center">
            {#each d.merchantKeys ?? [] as key}
              <span class="inline-flex items-center gap-1 px-2.5 py-1 bg-slate-700 rounded-full text-xs text-slate-300">
                <span class="capitalize">{key}</span>
                <button
                  class="text-slate-500 hover:text-red-400 ml-0.5"
                  onclick={() => handleRemoveKey(d.id, key)}
                >×</button>
              </span>
            {/each}

            {#if addingKeyForId === d.id}
              <div class="relative">
                <input
                  type="text"
                  bind:value={merchantKeySearch}
                  class="px-2.5 py-1 bg-slate-700 border border-slate-600 rounded-full text-xs text-white focus:outline-none focus:border-blue-500 w-48"
                  placeholder="Sok handlare..."
                  autofocus
                  onkeydown={(e) => { if (e.key === 'Escape') { addingKeyForId = null; merchantKeySearch = '' } }}
                />
                {#if filteredKeys.length > 0}
                  <div class="absolute top-full left-0 mt-1 w-64 max-h-48 overflow-y-auto bg-slate-700 border border-slate-600 rounded-lg shadow-lg z-50">
                    {#each filteredKeys.slice(0, 20) as key}
                      <button
                        class="w-full text-left px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-600 capitalize"
                        onclick={() => handleAddKey(d.id, key)}
                      >
                        {key}
                      </button>
                    {/each}
                    {#if filteredKeys.length > 20}
                      <div class="px-3 py-1.5 text-xs text-slate-500">...och {filteredKeys.length - 20} till</div>
                    {/if}
                  </div>
                {/if}
              </div>
            {:else}
              <button
                class="inline-flex items-center px-2.5 py-1 border border-dashed border-slate-600 rounded-full text-xs text-slate-400 hover:text-white hover:border-slate-500 transition-colors"
                onclick={() => { addingKeyForId = d.id; merchantKeySearch = '' }}
              >
                + Lagg till handlare
              </button>
            {/if}
          </div>

          <!-- Transaction toggle -->
          <button
            class="w-full px-4 py-2 text-xs text-slate-400 hover:text-white hover:bg-slate-700/50 border-t border-slate-700 transition-colors text-left"
            onclick={() => toggleTransactions(d.id)}
          >
            {expandedId === d.id ? '▼ Dolj transaktioner' : '▶ Visa transaktioner'}
          </button>

          {#if expandedId === d.id}
            <div class="border-t border-slate-700">
              <table class="w-full text-sm">
                <thead class="bg-slate-700/30">
                  <tr>
                    <th class="px-4 py-2 text-left text-slate-400 font-medium text-xs">Datum</th>
                    <th class="px-4 py-2 text-left text-slate-400 font-medium text-xs">Beskrivning</th>
                    <th class="px-4 py-2 text-left text-slate-400 font-medium text-xs">Typ</th>
                    <th class="px-4 py-2 text-right text-slate-400 font-medium text-xs">Belopp</th>
                    <th class="px-4 py-2 text-right text-slate-400 font-medium text-xs w-20"></th>
                  </tr>
                </thead>
                <tbody>
                  {#each transactions as t}
                    {#if editingTxId === t.id && t.isManual}
                      <tr class="border-t border-slate-700/30 bg-slate-700/20">
                        <td class="px-4 py-1.5">
                          <input type="date" bind:value={txFormDate}
                            class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs focus:outline-none focus:border-blue-500" />
                        </td>
                        <td class="px-4 py-1.5">
                          <input type="text" bind:value={txFormDescription} placeholder="Beskrivning"
                            class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs focus:outline-none focus:border-blue-500" />
                        </td>
                        <td class="px-4 py-1.5"></td>
                        <td class="px-4 py-1.5">
                          <input type="number" bind:value={txFormAmount} step="0.01" placeholder="0.00"
                            class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs text-right focus:outline-none focus:border-blue-500" />
                        </td>
                        <td class="px-4 py-1.5 text-right">
                          <button class="text-green-400 hover:text-green-300 text-xs mr-1" onclick={() => handleSaveTx(d.id)}>Spara</button>
                          <button class="text-slate-400 hover:text-white text-xs" onclick={cancelTxForm}>Avbryt</button>
                        </td>
                      </tr>
                    {:else}
                      <tr class="border-t border-slate-700/30">
                        <td class="px-4 py-2 text-slate-400 text-xs">{t.transactionDate}</td>
                        <td class="px-4 py-2 text-slate-300 text-xs">{t.description}</td>
                        <td class="px-4 py-2 text-xs">
                          {#if t.isManual}
                            <span class="inline-flex items-center px-1.5 py-0.5 bg-purple-500/20 text-purple-300 rounded text-[10px] font-medium">Manuell</span>
                          {:else}
                            <span class="text-slate-400 capitalize">{t.merchantKey}</span>
                          {/if}
                        </td>
                        <td class="px-4 py-2 text-right text-xs font-medium {t.amount > 0 ? 'text-green-400' : 'text-red-400'}">
                          {t.amount > 0 ? '+' : ''}{formatAmount(t.amount)} kr
                        </td>
                        <td class="px-4 py-2 text-right text-xs">
                          {#if t.isManual}
                            <button class="text-slate-400 hover:text-white mr-1" onclick={() => openEditTx(t)}>✏️</button>
                            <button class="text-slate-400 hover:text-red-400" onclick={() => handleDeleteTx(d.id, t.id)}>🗑️</button>
                          {/if}
                        </td>
                      </tr>
                    {/if}
                  {/each}

                  <!-- Add transaction inline form -->
                  {#if addingTxForId === d.id}
                    <tr class="border-t border-slate-700/30 bg-slate-700/20">
                      <td class="px-4 py-1.5">
                        <input type="date" bind:value={txFormDate}
                          class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs focus:outline-none focus:border-blue-500" />
                      </td>
                      <td class="px-4 py-1.5">
                        <input type="text" bind:value={txFormDescription} placeholder="Beskrivning"
                          class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs focus:outline-none focus:border-blue-500" />
                      </td>
                      <td class="px-4 py-1.5"></td>
                      <td class="px-4 py-1.5">
                        <input type="number" bind:value={txFormAmount} step="0.01" placeholder="0.00"
                          class="w-full px-2 py-1 bg-slate-700 border border-slate-600 rounded text-white text-xs text-right focus:outline-none focus:border-blue-500" />
                      </td>
                      <td class="px-4 py-1.5 text-right">
                        <button class="text-green-400 hover:text-green-300 text-xs mr-1" onclick={() => handleSaveTx(d.id)}>Spara</button>
                        <button class="text-slate-400 hover:text-white text-xs" onclick={cancelTxForm}>Avbryt</button>
                      </td>
                    </tr>
                  {/if}
                </tbody>
              </table>

              {#if transactions.length === 0 && addingTxForId !== d.id}
                <div class="px-4 py-3 text-sm text-slate-500">Inga transaktioner</div>
              {/if}

              {#if addingTxForId !== d.id}
                <button
                  class="w-full px-4 py-2 text-xs text-blue-400 hover:text-blue-300 hover:bg-slate-700/30 transition-colors text-left"
                  onclick={() => openAddTx(d.id)}
                >
                  + Lagg till transaktion
                </button>
              {/if}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <Modal title={editingDebtor ? 'Redigera person' : 'Ny person'} onClose={() => showModal = false}>
    <div class="space-y-4">
      <!-- Live preview -->
      <div class="flex items-center gap-3 p-3 bg-slate-700/50 rounded-lg">
        <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg"
             style="background-color: {hexToRgba(formColor, 0.15)}">
          {formIcon}
        </div>
        <span class="text-white font-medium">{formName || 'Namn'}</span>
      </div>

      <div>
        <label class="block text-sm text-slate-400 mb-1">Namn</label>
        <input
          type="text"
          bind:value={formName}
          class="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white text-sm focus:outline-none focus:border-blue-500"
          placeholder="Namn"
        />
      </div>

      <div>
        <label class="block text-sm text-slate-400 mb-1">Ikon</label>
        <div class="flex flex-wrap gap-2">
          {#each defaultIcons as icon}
            <button
              class="w-9 h-9 rounded-lg flex items-center justify-center transition-colors
                {formIcon === icon ? 'ring-2' : 'bg-slate-700 hover:bg-slate-600'}"
              style={formIcon === icon ? `background-color: ${hexToRgba(formColor, 0.2)}; box-shadow: 0 0 0 2px ${formColor}` : ''}
              onclick={() => formIcon = icon}
            >
              {icon}
            </button>
          {/each}
        </div>
      </div>

      <div>
        <label class="block text-sm text-slate-400 mb-1">Farg</label>
        <div class="flex flex-wrap gap-2">
          {#each defaultColors as color}
            <button
              class="w-8 h-8 rounded-lg transition-transform {formColor === color ? 'ring-2 ring-white scale-110' : 'hover:scale-105'}"
              style="background-color: {color}"
              onclick={() => formColor = color}
            ></button>
          {/each}
        </div>
        <input
          type="text"
          value={formColor}
          oninput={(e) => {
            const v = (e.currentTarget as HTMLInputElement).value
            if (/^#[0-9a-fA-F]{6}$/.test(v)) formColor = v
          }}
          class="mt-2 w-full px-3 py-1.5 bg-slate-700 border border-slate-600 rounded-lg text-white text-sm font-mono focus:outline-none focus:border-blue-500"
          placeholder="#FF0000"
        />
      </div>

      <button
        class="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
        onclick={handleSave}
      >
        {editingDebtor ? 'Spara' : 'Skapa'}
      </button>
    </div>
  </Modal>
{/if}
