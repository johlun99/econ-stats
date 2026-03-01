<script lang="ts">
  interface Props {
    message: string
    type?: 'success' | 'error' | 'info'
    onClose: () => void
  }

  let { message, type = 'info', onClose }: Props = $props()

  const colors = {
    success: 'bg-green-600 border-green-500',
    error: 'bg-red-600 border-red-500',
    info: 'bg-blue-600 border-blue-500',
  }

  const icons = {
    success: '✓',
    error: '✕',
    info: 'ℹ',
  }
</script>

<div class="fixed bottom-4 right-4 z-50 animate-slide-up">
  <div class="flex items-center gap-3 px-4 py-3 rounded-lg border shadow-lg {colors[type]}">
    <span class="text-white font-bold">{icons[type]}</span>
    <span class="text-white text-sm">{message}</span>
    <button class="text-white/70 hover:text-white ml-2" onclick={onClose}>✕</button>
  </div>
</div>

<style>
  @keyframes slide-up {
    from { transform: translateY(100%); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
  .animate-slide-up {
    animation: slide-up 0.3s ease-out;
  }
</style>
