import { ref } from 'vue'

export function useToast() {
  const toast = ref({
    show: false,
    message: '',
    type: 'info', // 'success', 'error', 'warning', 'info'
    duration: 3000
  })

  const showToast = (message, type = 'info', duration = 3000) => {
    toast.value = {
      show: true,
      message,
      type,
      duration
    }

    // 自动隐藏
    setTimeout(() => {
      toast.value.show = false
    }, duration)
  }

  const hideToast = () => {
    toast.value.show = false
  }

  return {
    toast,
    showToast,
    hideToast
  }
}

// 简单的全局toast组件
export const ToastComponent = {
  name: 'Toast',
  props: {
    toast: {
      type: Object,
      required: true
    }
  },
  template: `
    <div v-if="toast.show" class="toast" :class="toast.type">
      <div class="toast-content">
        <span class="toast-message">{{ toast.message }}</span>
        <button class="toast-close" @click="toast.show = false">×</button>
      </div>
    </div>
  `
}