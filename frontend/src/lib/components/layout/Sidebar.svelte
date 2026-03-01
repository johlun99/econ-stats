<script lang="ts">
  import type { Page } from '../../types'

  interface Props {
    currentPage: Page
    onNavigate: (page: Page) => void
  }

  let { currentPage, onNavigate }: Props = $props()

  const navItems: { page: Page; label: string; icon: string }[] = [
    { page: 'dashboard', label: 'Dashboard', icon: '📊' },
    { page: 'upload', label: 'Importera', icon: '📁' },
    { page: 'categorize', label: 'Kategorisera', icon: '🏷️' },
    { page: 'categories', label: 'Kategorier', icon: '⚙️' },
    { page: 'transactions', label: 'Transaktioner', icon: '📋' },
  ]
</script>

<aside class="w-56 bg-slate-800 border-r border-slate-700 flex flex-col">
  <div class="p-4 border-b border-slate-700">
    <h1 class="text-lg font-bold text-white">Econ Stats</h1>
    <p class="text-xs text-slate-400">Privatekonomisk överblick</p>
  </div>

  <nav class="flex-1 p-2 space-y-1">
    {#each navItems as item}
      <button
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors text-left
          {currentPage === item.page
            ? 'bg-blue-600 text-white'
            : 'text-slate-300 hover:bg-slate-700 hover:text-white'}"
        onclick={() => onNavigate(item.page)}
      >
        <span class="text-base">{item.icon}</span>
        <span>{item.label}</span>
      </button>
    {/each}
  </nav>
</aside>
