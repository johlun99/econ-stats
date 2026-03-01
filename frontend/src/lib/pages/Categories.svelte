<script lang="ts">
  import { onMount } from 'svelte'
  import { GetCategories, CreateCategory, UpdateCategory, DeleteCategory, GetCategoryRules, DeleteCategoryRule } from '../../../wailsjs/go/app/App'
  import Modal from '../components/common/Modal.svelte'
  import type { Category, CategoryRule } from '../types'
  import { hexToRgba } from '../utils'

  interface Props {
    onToast: (message: string, type: 'success' | 'error' | 'info') => void
  }

  let { onToast }: Props = $props()

  let categories: Category[] = $state([])
  let rules: CategoryRule[] = $state([])
  let showModal = $state(false)
  let editingCategory: Category | null = $state(null)

  // Form state
  let formName = $state('')
  let formColor = $state('#3B82F6')
  let formIcon = $state('📦')
  let formIsIncome = $state(false)
  let formIsExpense = $state(true)

  async function load() {
    const [c, r] = await Promise.all([GetCategories(), GetCategoryRules()])
    categories = c ?? []
    rules = r ?? []
  }

  function openCreate() {
    editingCategory = null
    formName = ''
    formColor = '#3B82F6'
    formIcon = '📦'
    formIsIncome = false
    formIsExpense = true
    showModal = true
  }

  function openEdit(cat: Category) {
    editingCategory = cat
    formName = cat.name
    formColor = cat.color
    formIcon = cat.icon
    formIsIncome = cat.isIncome
    formIsExpense = cat.isExpense
    showModal = true
  }

  async function handleSave() {
    if (!formName.trim()) return
    try {
      if (editingCategory) {
        await UpdateCategory(editingCategory.id, formName, formColor, formIcon, formIsIncome, formIsExpense)
        onToast('Kategori uppdaterad', 'success')
      } else {
        await CreateCategory(formName, formColor, formIcon, formIsIncome, formIsExpense)
        onToast('Kategori skapad', 'success')
      }
      showModal = false
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleDelete(id: number) {
    try {
      await DeleteCategory(id)
      onToast('Kategori borttagen', 'success')
      await load()
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  async function handleDeleteRule(id: number) {
    try {
      await DeleteCategoryRule(id)
      onToast('Regel borttagen', 'success')
      rules = rules.filter(r => r.id !== id)
    } catch (e: any) {
      onToast('Fel: ' + (e?.message || e), 'error')
    }
  }

  const defaultColors = [
    '#EF4444', '#F97316', '#F59E0B', '#EAB308',
    '#84CC16', '#22C55E', '#10B981', '#14B8A6',
    '#06B6D4', '#3B82F6', '#6366F1', '#8B5CF6',
    '#A855F7', '#D946EF', '#EC4899', '#F43F5E',
    '#64748B', '#6B7280', '#78716C', '#FFFFFF',
  ]

  const defaultIcons = ['📦', '🏠', '🛒', '🍽️', '🚗', '🎮', '💊', '📱', '🛍️', '💰', '🛡️', '💵', '💸', '🎬', '✈️', '🐕', '📚', '🏋️']

  onMount(load)
</script>

<div class="p-6">
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-2xl font-bold text-white">Kategorier</h2>
    <button
      class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm transition-colors"
      onclick={openCreate}
    >
      + Ny kategori
    </button>
  </div>

  <!-- Categories grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
    {#each categories as cat}
      <div class="bg-slate-800 rounded-xl p-4 border border-slate-700 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg"
               style="background-color: {hexToRgba(cat.color, 0.15)}">
            {cat.icon}
          </div>
          <div>
            <div class="text-white font-medium">{cat.name}</div>
            <div class="text-xs text-slate-400">{[cat.isExpense && 'Utgift', cat.isIncome && 'Inkomst'].filter(Boolean).join(' & ') || 'Inga typer'}</div>
          </div>
        </div>
        <div class="flex gap-1">
          <button class="p-2 text-slate-400 hover:text-white rounded" onclick={() => openEdit(cat)}>✏️</button>
          <button class="p-2 text-slate-400 hover:text-red-400 rounded" onclick={() => handleDelete(cat.id)}>🗑️</button>
        </div>
      </div>
    {/each}
  </div>

  <!-- Rules section -->
  {#if rules.length > 0}
    <h3 class="text-lg font-semibold text-white mb-4">Kategoriregler ({rules.length})</h3>
    <div class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-700/50">
          <tr>
            <th class="px-4 py-3 text-left text-slate-400 font-medium">Handlare</th>
            <th class="px-4 py-3 text-left text-slate-400 font-medium">Kategori</th>
            <th class="px-4 py-3 text-right text-slate-400 font-medium"></th>
          </tr>
        </thead>
        <tbody>
          {#each rules as rule}
            <tr class="border-t border-slate-700/50">
              <td class="px-4 py-2.5 text-slate-200 capitalize">{rule.merchantKey}</td>
              <td class="px-4 py-2.5 text-slate-300">{rule.categoryName}</td>
              <td class="px-4 py-2.5 text-right">
                <button class="text-slate-400 hover:text-red-400 text-xs" onclick={() => handleDeleteRule(rule.id)}>
                  Ta bort
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if showModal}
  <Modal title={editingCategory ? 'Redigera kategori' : 'Ny kategori'} onClose={() => showModal = false}>
    <div class="space-y-4">
      <!-- Live preview -->
      <div class="flex items-center gap-3 p-3 bg-slate-700/50 rounded-lg">
        <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg"
             style="background-color: {hexToRgba(formColor, 0.15)}">
          {formIcon}
        </div>
        <span class="text-white font-medium">{formName || 'Kategorinamn'}</span>
      </div>

      <div>
        <label class="block text-sm text-slate-400 mb-1">Namn</label>
        <input
          type="text"
          bind:value={formName}
          class="w-full px-3 py-2 bg-slate-700 border border-slate-600 rounded-lg text-white text-sm focus:outline-none focus:border-blue-500"
          placeholder="Kategorinamn"
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
        <label class="block text-sm text-slate-400 mb-1">Färg</label>
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

      <div class="space-y-2">
        <p class="text-sm text-slate-400">Typ</p>
        <label class="flex items-center gap-2 text-sm text-slate-300">
          <input type="checkbox" bind:checked={formIsExpense}
                 class="rounded bg-slate-700 border-slate-600" />
          Utgift
        </label>
        <label class="flex items-center gap-2 text-sm text-slate-300">
          <input type="checkbox" bind:checked={formIsIncome}
                 class="rounded bg-slate-700 border-slate-600" />
          Inkomst
        </label>
      </div>

      <button
        class="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors"
        onclick={handleSave}
      >
        {editingCategory ? 'Spara' : 'Skapa'}
      </button>
    </div>
  </Modal>
{/if}
