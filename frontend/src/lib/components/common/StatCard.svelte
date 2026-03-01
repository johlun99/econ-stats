<script lang="ts">
  interface Props {
    label: string
    value: string
    subtitle?: string
    trend?: number
    icon?: string
  }

  let { label, value, subtitle = '', trend, icon = '' }: Props = $props()
</script>

<div class="bg-slate-800 rounded-xl p-5 border border-slate-700">
  <div class="flex items-center justify-between mb-2">
    <span class="text-sm text-slate-400">{label}</span>
    {#if icon}
      <span class="text-lg">{icon}</span>
    {/if}
  </div>
  <div class="text-2xl font-bold text-white">{value}</div>
  <div class="flex items-center gap-2 mt-1">
    {#if subtitle}
      <span class="text-xs text-slate-400">{subtitle}</span>
    {/if}
    {#if trend !== undefined}
      <span class="text-xs {trend >= 0 ? 'text-red-400' : 'text-green-400'}">
        {trend >= 0 ? '↑' : '↓'} {Math.abs(trend).toFixed(1)}%
      </span>
    {/if}
  </div>
</div>
