<script lang="ts">
  import { Chart, LineController, LineElement, PointElement, CategoryScale, LinearScale, Tooltip, Filler } from 'chart.js'
  import { onMount } from 'svelte'
  import type { DailySpend } from '../../types'

  Chart.register(LineController, LineElement, PointElement, CategoryScale, LinearScale, Tooltip, Filler)

  interface Props {
    data: DailySpend[]
  }

  let { data }: Props = $props()
  let canvas: HTMLCanvasElement
  let chart: Chart | null = null

  function renderChart() {
    if (chart) chart.destroy()
    if (!canvas || data.length === 0) return

    chart = new Chart(canvas, {
      type: 'line',
      data: {
        labels: data.map(d => d.date.slice(8)), // day only
        datasets: [{
          label: 'Dagliga utgifter',
          data: data.map(d => d.total),
          borderColor: '#3B82F6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          fill: true,
          tension: 0.3,
          pointRadius: 3,
          pointBackgroundColor: '#3B82F6',
        }]
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
          tooltip: {
            callbacks: {
              label: (ctx) => `${ctx.parsed.y.toLocaleString('sv-SE')} kr`
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
