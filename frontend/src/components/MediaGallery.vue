<template>
  <div class="media-gallery card">
    <div class="gallery-header">
      <div class="tabs">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'image' }"
          @click="switchTab('image')"
        >
          <span class="icon">🖼️</span>
          图片展示
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'video' }"
          @click="switchTab('video')"
        >
          <span class="icon">🎬</span>
          视频展示
        </button>
      </div>
    </div>

    <div class="gallery-content">
      <div v-if="activeTab === 'image'" class="image-section">
        <div class="carousel">
          <div class="carousel-container" :style="{ transform: `translateX(-${currentIndex * 100}%)` }">
            <div 
              v-for="(image, index) in images" 
              :key="index" 
              class="carousel-slide"
            >
              <img 
                :src="image" 
                :alt="`产品图片 ${index + 1}`"
                @click="openPreview(index)"
              />
            </div>
          </div>
          
          <button class="carousel-btn prev" @click="prevSlide" v-if="images.length > 1">
            <span>‹</span>
          </button>
          <button class="carousel-btn next" @click="nextSlide" v-if="images.length > 1">
            <span>›</span>
          </button>
          
          <div class="carousel-indicators" v-if="images.length > 1">
            <span 
              v-for="(_, index) in images" 
              :key="index"
              class="indicator"
              :class="{ active: currentIndex === index }"
              @click="goToSlide(index)"
            ></span>
          </div>
        </div>
        
        <div class="image-count">
          共 {{ images.length }} 张图片
        </div>
      </div>

      <div v-else class="video-section">
        <div class="video-placeholder">
          <div class="video-icon">🎬</div>
          <p class="video-text">视频播放区</p>
          <p class="video-hint">支持播放、暂停、全屏功能</p>
          <div class="video-controls">
            <button class="control-btn" @click="togglePlay">
              {{ isPlaying ? '⏸️ 暂停' : '▶️ 播放' }}
            </button>
            <button class="control-btn" @click="toggleFullscreen">
              ⛶ 全屏
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showPreview" class="image-preview" @click="closePreview">
      <div class="preview-content" @click.stop>
        <button class="close-btn" @click="closePreview">×</button>
        <img :src="images[previewIndex]" alt="预览图片" />
        <div class="preview-nav">
          <button @click="prevPreview" v-if="images.length > 1">‹ 上一张</button>
          <span>{{ previewIndex + 1 }} / {{ images.length }}</span>
          <button @click="nextPreview" v-if="images.length > 1">下一张 ›</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  images: {
    type: Array,
    default: () => []
  },
  videos: {
    type: Array,
    default: () => []
  }
})

const activeTab = ref('image')
const currentIndex = ref(0)
const showPreview = ref(false)
const previewIndex = ref(0)
const isPlaying = ref(false)
let autoPlayTimer = null

const switchTab = (tab) => {
  activeTab.value = tab
  if (tab === 'image') {
    startAutoPlay()
  } else {
    stopAutoPlay()
  }
}

const prevSlide = () => {
  currentIndex.value = currentIndex.value === 0 
    ? props.images.length - 1 
    : currentIndex.value - 1
}

const nextSlide = () => {
  currentIndex.value = currentIndex.value === props.images.length - 1 
    ? 0 
    : currentIndex.value + 1
}

const goToSlide = (index) => {
  currentIndex.value = index
}

const startAutoPlay = () => {
  stopAutoPlay()
  if (props.images.length > 1) {
    autoPlayTimer = setInterval(() => {
      nextSlide()
    }, 5000)
  }
}

const stopAutoPlay = () => {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

const openPreview = (index) => {
  previewIndex.value = index
  showPreview.value = true
  stopAutoPlay()
}

const closePreview = () => {
  showPreview.value = false
  startAutoPlay()
}

const prevPreview = () => {
  previewIndex.value = previewIndex.value === 0 
    ? props.images.length - 1 
    : previewIndex.value - 1
}

const nextPreview = () => {
  previewIndex.value = previewIndex.value === props.images.length - 1 
    ? 0 
    : previewIndex.value + 1
}

const togglePlay = () => {
  isPlaying.value = !isPlaying.value
}

const toggleFullscreen = () => {
  alert('全屏功能')
}

onMounted(() => {
  startAutoPlay()
})

onUnmounted(() => {
  stopAutoPlay()
})
</script>

<style scoped>
.media-gallery {
  display: flex;
  flex-direction: column;
  min-height: 400px;
}

.gallery-header {
  margin-bottom: 16px;
}

.tabs {
  display: flex;
  gap: 12px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border: none;
  background: var(--background-color);
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--secondary-text);
  transition: all 0.3s ease;
}

.tab-btn:hover {
  background: var(--secondary-color);
  color: white;
}

.tab-btn.active {
  background: var(--primary-color);
  color: white;
}

.tab-btn .icon {
  font-size: 16px;
}

.gallery-content {
  flex: 1;
  position: relative;
}

.carousel {
  position: relative;
  width: 100%;
  height: 350px;
  overflow: hidden;
  border-radius: 8px;
}

.carousel-container {
  display: flex;
  height: 100%;
  transition: transform 0.5s ease-in-out;
}

.carousel-slide {
  min-width: 100%;
  height: 100%;
}

.carousel-slide img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.carousel-slide img:hover {
  transform: scale(1.02);
}

.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border: none;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  cursor: pointer;
  font-size: 24px;
  color: var(--text-color);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10;
}

.carousel-btn:hover {
  background: var(--primary-color);
  color: white;
}

.carousel-btn.prev {
  left: 12px;
}

.carousel-btn.next {
  right: 12px;
}

.carousel-indicators {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  z-index: 10;
}

.indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: var(--primary-color);
  width: 24px;
  border-radius: 5px;
}

.image-count {
  text-align: center;
  margin-top: 12px;
  font-size: 12px;
  color: var(--secondary-text);
}

.video-section {
  height: 350px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-placeholder {
  text-align: center;
  padding: 40px;
  background: linear-gradient(135deg, var(--background-color), #e8f5e9);
  border-radius: 12px;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.video-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.video-text {
  font-size: 18px;
  color: var(--text-color);
  margin-bottom: 8px;
}

.video-hint {
  font-size: 12px;
  color: var(--secondary-text);
  margin-bottom: 20px;
}

.video-controls {
  display: flex;
  gap: 12px;
}

.control-btn {
  padding: 10px 24px;
  border: none;
  background: var(--primary-color);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.control-btn:hover {
  background: var(--secondary-color);
}

.image-preview {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.preview-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
}

.preview-content img {
  max-width: 100%;
  max-height: 80vh;
  border-radius: 8px;
}

.close-btn {
  position: absolute;
  top: -40px;
  right: 0;
  width: 36px;
  height: 36px;
  border: none;
  background: white;
  border-radius: 50%;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  color: white;
}

.preview-nav button {
  padding: 8px 16px;
  border: none;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.preview-nav button:hover {
  background: var(--primary-color);
}
</style>
