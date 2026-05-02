<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProductStore } from '@/stores/product'
import { storeToRefs } from 'pinia'
import ImageGallery from '@/components/ImageGallery.vue'
import VideoPlayer from '@/components/VideoPlayer.vue'

const productStore = useProductStore()
const { currentProduct: product, loading } = storeToRefs(productStore)

const currentImageIndex = ref(0)
const showPreview = ref(false)

const allImages = computed(() => {
  if (!product.value) return []
  const mediaFiles = product.value.media_files?.filter(f => f.media_type === 'Image') || []
  return [...product.value.images, ...mediaFiles.map(f => f.url)].filter(Boolean)
})

const videoUrl = computed(() => {
  if (!product.value) return null
  return product.value.video || product.value.media_files?.find(f => f.media_type === 'Video')?.url || null
})

const openPreview = (index: number) => {
  currentImageIndex.value = index
  showPreview.value = true
}

const closePreview = () => {
  showPreview.value = false
}
</script>

<template>
  <div class="media-view">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <div class="skeleton" style="height: 300px; margin-bottom: 16px;"></div>
        <div class="grid grid-cols-2 gap-3">
          <div class="skeleton" style="height: 150px;" v-for="i in 4" :key="i"></div>
        </div>
      </div>

      <div v-else-if="product">
        <section class="media-section" v-if="videoUrl">
          <h2 class="section-title">
            <span class="i-ph-video-camera text-accent"></span>
            产品视频
          </h2>
          <VideoPlayer :src="videoUrl" class="video-container" />
        </section>

        <section class="media-section mt-6">
          <h2 class="section-title">
            <span class="i-ph-images text-primary"></span>
            产品相册
            <span class="section-count">{{ allImages.length }} 张</span>
          </h2>
          
          <div v-if="allImages.length > 0" class="gallery-grid">
            <div 
              v-for="(image, index) in allImages" 
              :key="index"
              class="gallery-item"
              @click="openPreview(index)"
            >
              <img :src="image" :alt="`产品图片 ${index + 1}`" loading="lazy" />
              <div class="gallery-overlay">
                <span class="i-ph-magnifying-glass text-2xl"></span>
              </div>
            </div>
          </div>
          
          <div v-else class="empty-gallery">
            <span class="i-ph-image-broken text-5xl text-gray-300"></span>
            <p class="mt-4 text-gray-500">暂无图片</p>
          </div>
        </section>
      </div>

      <Teleport to="body">
        <div v-if="showPreview" class="preview-modal" @click.self="closePreview">
          <button class="preview-close" @click="closePreview">
            <span class="i-ph-x text-2xl"></span>
          </button>
          
          <button class="preview-nav prev" @click="currentImageIndex = (currentImageIndex - 1 + allImages.length) % allImages.length">
            <span class="i-ph-caret-left text-2xl"></span>
          </button>
          
          <div class="preview-content">
            <img :src="allImages[currentImageIndex]" :alt="`预览 ${currentImageIndex + 1}`" />
          </div>
          
          <button class="preview-nav next" @click="currentImageIndex = (currentImageIndex + 1) % allImages.length">
            <span class="i-ph-caret-right text-2xl"></span>
          </button>
          
          <div class="preview-indicator">
            {{ currentImageIndex + 1 }} / {{ allImages.length }}
          </div>
        </div>
      </Teleport>
    </div>
  </div>
</template>

<style scoped>
.media-view {
  padding-top: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #333;
}

.section-count {
  font-size: 12px;
  color: #999;
  font-weight: 400;
  margin-left: auto;
}

.media-section {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
}

.video-container {
  border-radius: 12px;
  overflow: hidden;
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (min-width: 768px) {
  .gallery-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.gallery-item {
  position: relative;
  aspect-ratio: 4/3;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
}

.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.gallery-item:hover img {
  transform: scale(1.05);
}

.gallery-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s ease;
  color: white;
}

.gallery-item:hover .gallery-overlay {
  opacity: 1;
}

.empty-gallery {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
}

.preview-modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.95);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-close {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 44px;
  height: 44px;
  border: none;
  background: rgba(255,255,255,0.1);
  border-radius: 50%;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s ease;
}

.preview-close:hover {
  background: rgba(255,255,255,0.2);
}

.preview-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 50px;
  height: 50px;
  border: none;
  background: rgba(255,255,255,0.1);
  border-radius: 50%;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s ease;
}

.preview-nav:hover {
  background: rgba(255,255,255,0.2);
}

.preview-nav.prev {
  left: 20px;
}

.preview-nav.next {
  right: 20px;
}

.preview-content {
  max-width: 90vw;
  max-height: 80vh;
}

.preview-content img {
  max-width: 100%;
  max-height: 80vh;
  object-fit: contain;
  border-radius: 8px;
}

.preview-indicator {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  color: white;
  font-size: 14px;
  background: rgba(0,0,0,0.5);
  padding: 8px 16px;
  border-radius: 20px;
}

.loading-state {
  padding: 24px 0;
}

.grid {
  display: grid;
}

.grid-cols-2 {
  grid-template-columns: repeat(2, 1fr);
}

.gap-3 {
  gap: 12px;
}

.mt-4 {
  margin-top: 16px;
}

.mt-6 {
  margin-top: 24px;
}
</style>
