<script lang="ts">
  import { onMount } from 'svelte'
  import { GetMonthlyStats, GetSpendingTrend, GetAvailableMonths, GetPinnedDebtors } from '../../../wailsjs/go/app/App'
  import StatCard from '../components/common/StatCard.svelte'
  import MonthPicker from '../components/common/MonthPicker.svelte'
  import DoughnutChart from '../components/charts/DoughnutChart.svelte'
  import BarChart from '../components/charts/BarChart.svelte'
  import LineChart from '../components/charts/LineChart.svelte'
  import type { MonthlyStats, AvailableMonth, SpendingTrend, DebtorDetail } from '../types'
  import { hexToRgba } from '../utils'

  let stats: MonthlyStats | null = $state(null)
  let months: AvailableMonth[] = $state([])
  let selectedMonth = $state('')
  let trend: SpendingTrend[] = $state([])
  let pinnedDebtors: DebtorDetail[] = $state([])
  let loading = $state(true)

  let yearTrend = $derived(
    trend.filter(t => selectedMonth && t.month.startsWith(selectedMonth.slice(0, 4)))
  )

  function fmt(n: number): string {
    return n.toLocaleString('sv-SE', { minimumFractionDigits: 0, maximumFractionDigits: 0 }) + ' kr'
  }

  async function loadData() {
    loading = true
    try {
      const [m, t, pd] = await Promise.all([
        GetAvailableMonths(),
        GetSpendingTrend(0),
        GetPinnedDebtors(),
      ])
      months = m ?? []
      trend = t ?? []
      pinnedDebtors = pd ?? []
      if (months.length > 0 && !selectedMonth) {
        selectedMonth = months[0].month
      }
      if (selectedMonth) {
        stats = await GetMonthlyStats(selectedMonth)
      }
    } catch (e) {
      console.error(e)
    } finally {
      loading = false
    }
  }

  async function selectMonth(month: string) {
    selectedMonth = month
    stats = await GetMonthlyStats(month)
  }

  onMount(loadData)
</script>

<div class="p-6 space-y-6">
  <div class="flex items-center justify-between">
    <h2 class="text-2xl font-bold text-white">Dashboard</h2>
    {#if months.length > 0}
      <MonthPicker {months} selected={selectedMonth} onSelect={selectMonth} />
    {/if}
  </div>

  {#if loading}
    <div class="flex items-center justify-center h-64 text-slate-400">Laddar...</div>
  {:else if !stats || months.length === 0}
    <div class="flex flex-col items-center justify-center h-64 text-slate-400">
      <p class="text-lg mb-2">Ingen data ännu</p>
      <p class="text-sm">Importera transaktioner för att komma igång</p>
    </div>
  {:else}
    <!-- Stat cards -->
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
      <StatCard label="Utgifter" value={fmt(stats.totalExpenses)} icon="💸" trend={stats.monthOverMonth} />
      <StatCard label="Inkomster" value={fmt(stats.totalIncome)} icon="💰" />
      <StatCard label="Nettosparande" value={fmt(stats.netSavings)} icon="🏦"
        subtitle={stats.netSavings >= 0 ? 'Positivt' : 'Negativt'} />
      <StatCard label="Sparkvot" value={stats.savingsRate.toFixed(1) + '%'} icon="📈" />
      <StatCard label="Snittutgift/dag" value={fmt(stats.avgDailySpend)} icon="📅" />
    </div>

    <!-- Pinned debtors -->
    {#if pinnedDebtors.length > 0}
      <div>
        <h3 class="text-sm font-semibold text-slate-300 mb-3">Fastnålade skulder</h3>
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
          {#each pinnedDebtors as d}
            <div class="bg-slate-800 rounded-xl p-4 border border-slate-700 flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center text-lg shrink-0"
                   style="background-color: {hexToRgba(d.color, 0.15)}">
                {d.icon}
              </div>
              <div class="min-w-0">
                <div class="text-sm text-white font-medium truncate">{d.name}</div>
                <div class="text-lg font-semibold {d.balance > 0 ? 'text-green-400' : d.balance < 0 ? 'text-red-400' : 'text-slate-400'}">
                  {d.balance > 0 ? '+' : ''}{fmt(d.balance)}
                </div>
                <div class="text-xs {d.balance > 0 ? 'text-green-400/70' : d.balance < 0 ? 'text-red-400/70' : 'text-slate-500'}">
                  {d.balance > 0 ? 'De ar skyldiga dig' : d.balance < 0 ? 'Du ar skyldig dem' : 'Kvitt'}
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Charts row 1 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
        <h3 class="text-sm font-semibold text-slate-300 mb-4">Utgifter per kategori</h3>
        {#if stats.categoryBreakdown?.length}
          <DoughnutChart data={stats.categoryBreakdown} />
        {:else}
          <p class="text-slate-500 text-sm text-center py-8">Ingen data</p>
        {/if}
      </div>

      <div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
        <h3 class="text-sm font-semibold text-slate-300 mb-4">Dagliga utgifter</h3>
        {#if stats.dailySpending?.length}
          <LineChart data={stats.dailySpending} />
        {:else}
          <p class="text-slate-500 text-sm text-center py-8">Ingen data</p>
        {/if}
      </div>
    </div>

    <!-- Charts row 2 -->
    {#if yearTrend.length > 1}
      <div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
        <h3 class="text-sm font-semibold text-slate-300 mb-4">Inkomster vs Utgifter</h3>
        <BarChart data={yearTrend} />
      </div>
    {/if}

    <!-- Bottom row: Top merchants + Largest expenses -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Top merchants -->
      <div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
        <h3 class="text-sm font-semibold text-slate-300 mb-4">Topp handlare</h3>
        {#if stats.topMerchants?.length}
          <div class="space-y-2">
            {#each stats.topMerchants as m, i}
              <div class="flex items-center justify-between py-1.5 {i < stats.topMerchants.length - 1 ? 'border-b border-slate-700' : ''}">
                <div class="flex items-center gap-3">
                  <span class="text-xs text-slate-500 w-5">{i + 1}</span>
                  <span class="text-sm text-slate-200 capitalize">{m.merchantKey}</span>
                </div>
                <div class="text-right">
                  <span class="text-sm font-medium text-white">{fmt(m.total)}</span>
                  <span class="text-xs text-slate-500 ml-2">({m.count}x)</span>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-slate-500 text-sm text-center py-4">Ingen data</p>
        {/if}
      </div>

      <!-- Largest expenses -->
      <div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
        <h3 class="text-sm font-semibold text-slate-300 mb-4">Största utgifter</h3>
        {#if stats.largestExpenses?.length}
          <div class="space-y-2">
            {#each stats.largestExpenses as t, i}
              <div class="flex items-center justify-between py-1.5 {i < stats.largestExpenses.length - 1 ? 'border-b border-slate-700' : ''}">
                <div>
                  <span class="text-sm text-slate-200">{t.description}</span>
                  <span class="text-xs text-slate-500 ml-2">{t.transactionDate}</span>
                </div>
                <span class="text-sm font-medium text-red-400">{fmt(Math.abs(t.amount))}</span>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-slate-500 text-sm text-center py-4">Ingen data</p>
        {/if}
      </div>
    </div>
  {/if}
</div>
