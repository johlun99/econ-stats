<script lang="ts">
  import Sidebar from './lib/components/layout/Sidebar.svelte'
  import Dashboard from './lib/pages/Dashboard.svelte'
  import YearlyStats from './lib/pages/YearlyStats.svelte'
  import Upload from './lib/pages/Upload.svelte'
  import Categorize from './lib/pages/Categorize.svelte'
  import Categories from './lib/pages/Categories.svelte'
  import Transactions from './lib/pages/Transactions.svelte'
  import Debtors from './lib/pages/Debtors.svelte'
  import Toast from './lib/components/common/Toast.svelte'
  import type { Page } from './lib/types'

  let currentPage: Page = $state('dashboard')
  let toastMessage = $state('')
  let toastType: 'success' | 'error' | 'info' = $state('info')
  let showToast = $state(false)

  function navigate(page: Page) {
    currentPage = page
  }

  function toast(message: string, type: 'success' | 'error' | 'info' = 'info') {
    toastMessage = message
    toastType = type
    showToast = true
    setTimeout(() => showToast = false, 4000)
  }
</script>

<div class="flex h-screen bg-slate-900">
  <Sidebar {currentPage} onNavigate={navigate} />

  <main class="flex-1 overflow-y-auto">
    {#key currentPage}
      {#if currentPage === 'dashboard'}
        <Dashboard />
      {:else if currentPage === 'yearly'}
        <YearlyStats />
      {:else if currentPage === 'upload'}
        <Upload onComplete={() => { toast('Import klar!', 'success'); navigate('categorize') }} />
      {:else if currentPage === 'categorize'}
        <Categorize onToast={toast} />
      {:else if currentPage === 'categories'}
        <Categories onToast={toast} />
      {:else if currentPage === 'transactions'}
        <Transactions />
      {:else if currentPage === 'debtors'}
        <Debtors onToast={toast} />
      {/if}
    {/key}
  </main>
</div>

{#if showToast}
  <Toast message={toastMessage} type={toastType} onClose={() => showToast = false} />
{/if}
