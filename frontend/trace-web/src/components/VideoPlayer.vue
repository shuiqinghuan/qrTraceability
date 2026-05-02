<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  src: string
  poster?: string
}>()

const isPlaying = ref(false)
const isLoading = ref(true)

const onLoaded = () => {
  isLoading.value = false
}

const onError = () => {
  isLoading.value = false
}
</script>

<template>
  <div class="video-player">
    <video
      :src="src"
      :poster="poster"
      controls
      playsinline
      preload="metadata"
      @loadeddata="onLoaded"
      @error="onError"
    >
      您的浏览器不支持视频播放
    </video>
    
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-spinner"></div>
    </div>
  </div>
</template>

<style scoped>
.video-player {
  position: relative;
  width: 100%;
  aspect-ratio: 16/9;
  background: #000;
  border-radius: 16px;
  overflow: hidden;
}

video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.loading-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
