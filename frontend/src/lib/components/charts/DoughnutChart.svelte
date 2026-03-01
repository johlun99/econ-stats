<script lang="ts">
  import { Chart, DoughnutController, ArcElement, Tooltip, Legend } from 'chart.js'
  import { onMount } from 'svelte'
  import type { CategorySpend } from '../../types'

  Chart.register(DoughnutController, ArcElement, Tooltip, Legend)

  interface Props {
    data: CategorySpend[]
  }

  let { data }: Props = $props()
  let canvas: HTMLCanvasElement
  let chart: Chart | null = null

  function renderChart() {
    if (chart) chart.destroy()
    if (!canvas || data.length === 0) return

    chart = new Chart(canvas, {
      type: 'doughnut',
      data: {
        labels: data.map(d => d.categoryName),
        datasets: [{
          data: data.map(d => d.total),
          backgroundColor: data.map(d => d.categoryColor),
          borderColor: '#1E293B',
          borderWidth: 2,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'right',
            labels: { color: '#94A3B8', font: { size: 12 }, padding: 12 }
          },
          tooltip: {
            callbacks: {
              label: (ctx) => {
                const val = ctx.parsed
                const total = data.reduce((s, d) => s + d.total, 0)
                const pct = ((val / total) * 100).toFixed(1)
                return `${ctx.label}: ${val.toLocaleString('sv-SE')} kr (${pct}%)`
              }
            }
          }
        }
      }
    })
  }

  onMount(() => { renderChart() })

  $effect(() => {
    data; // track
    if (canvas) renderChart()
  })
</script>

<div class="h-64">
  <canvas bind:this={canvas}></canvas>
</div>
