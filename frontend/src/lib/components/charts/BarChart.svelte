<script lang="ts">
  import { Chart, BarController, BarElement, CategoryScale, LinearScale, Tooltip, Legend } from 'chart.js'
  import { onMount } from 'svelte'
  import type { SpendingTrend } from '../../types'

  Chart.register(BarController, BarElement, CategoryScale, LinearScale, Tooltip, Legend)

  interface Props {
    data: SpendingTrend[]
  }

  let { data }: Props = $props()
  let canvas: HTMLCanvasElement
  let chart: Chart | null = null

  function renderChart() {
    if (chart) chart.destroy()
    if (!canvas || data.length === 0) return

    chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels: data.map(d => d.month),
        datasets: [
          {
            label: 'Utgifter',
            data: data.map(d => d.expenses),
            backgroundColor: '#EF4444',
            borderRadius: 4,
          },
          {
            label: 'Inkomster',
            data: data.map(d => d.income),
            backgroundColor: '#22C55E',
            borderRadius: 4,
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          x: { ticks: { color: '#94A3B8' }, grid: { display: false } },
          y: {
            ticks: {
              color: '#94A3B8',
              callback: (v) => `${Number(v).toLocaleString('sv-SE')} kr`
            },
            grid: { color: '#334155' }
          }
        },
        plugins: {
          legend: { labels: { color: '#94A3B8' } },
          tooltip: {
            callbacks: {
              label: (ctx) => `${ctx.dataset.label}: ${ctx.parsed.y.toLocaleString('sv-SE')} kr`
            }
          }
        }
      }
    })
  }

  onMount(() => { renderChart() })
  $effect(() => { data; if (canvas) renderChart() })
</script>

<div class="h-64">
  <canvas bind:this={canvas}></canvas>
</div>
