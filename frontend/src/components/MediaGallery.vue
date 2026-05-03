<template>
  <div class="media-gallery card">
    <h2 class="gallery-title">
      <span class="icon">🖼️</span>
      产品图片
    </h2>

    <div class="gallery-content">
      <div class="image-section">
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

        <div class="image-count" v-if="images.length > 0">
          共 {{ images.length }} 张图片
        </div>
      </div>
    </div>

    <div v-if="showPreview" class="image-preview" @click="closePreview">
      <div class="preview-content" @click.stop>
        <button class="close-btn" @click="closePreview">×</button>
        <img :src="images[previewIndex]" alt="预览图片" />
        <div class="preview-nav" v-if="images.length > 1">
          <button @click="prevPreview" v-if="images.length > 1">‹ 上一张</button>
          <span>{{ previewIndex + 1 }} / {{ images.length }}</span>
          <button @click="nextPreview" v-if="images.length > 1">下一张 ›</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 媒体图片画廊组件
 * 提供图片轮播展示、自动播放、手动切换及全屏预览功能
 */

import { ref, onMounted, onUnmounted } from 'vue'

// 接收父组件传入的图片地址数组
const props = defineProps({
  images: {
    type: Array,
    default: () => []
  }
})

// 轮播当前索引、预览弹窗状态、预览图片索引
const currentIndex = ref(0)
const showPreview = ref(false)
const previewIndex = ref(0)
// 自动播放定时器引用
let autoPlayTimer = null

/** 切换到上一张轮播图 */
const prevSlide = () => {
  currentIndex.value = currentIndex.value === 0
    ? props.images.length - 1
    : currentIndex.value - 1
}

/** 切换到下一张轮播图 */
const nextSlide = () => {
  currentIndex.value = currentIndex.value === props.images.length - 1
    ? 0
    : currentIndex.value + 1
}

/** 跳转到指定索引的轮播图 */
const goToSlide = (index) => {
  currentIndex.value = index
}

/** 启动自动轮播（每5秒切换一次） */
const startAutoPlay = () => {
  stopAutoPlay()
  if (props.images.length > 1) {
    autoPlayTimer = setInterval(() => {
      nextSlide()
    }, 5000)
  }
}

/** 停止自动轮播 */
const stopAutoPlay = () => {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

/** 打开图片全屏预览，并暂停自动轮播 */
const openPreview = (index) => {
  previewIndex.value = index
  showPreview.value = true
  stopAutoPlay()
}

/** 关闭图片预览，并恢复自动轮播 */
const closePreview = () => {
  showPreview.value = false
  startAutoPlay()
}

/** 预览模式下切换到上一张 */
const prevPreview = () => {
  previewIndex.value = previewIndex.value === 0
    ? props.images.length - 1
    : previewIndex.value - 1
}

/** 预览模式下切换到下一张 */
const nextPreview = () => {
  previewIndex.value = previewIndex.value === props.images.length - 1
    ? 0
    : previewIndex.value + 1
}

// 组件挂载时启动自动轮播，卸载时清除定时器
onMounted(() => {
  startAutoPlay()
})

onUnmounted(() => {
  stopAutoPlay()
})
</script>

<style scoped>
/* 画廊根容器 */
.media-gallery {
  display: flex;
  flex-direction: column;
  min-height: 400px;
}

.gallery-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--primary-color);
  border-bottom: 2px solid var(--secondary-color);
  padding-bottom: 12px;
  margin-bottom: 16px;
  font-size: 18px;
}

.gallery-title .icon {
  font-size: 20px;
}

.gallery-content {
  flex: 1;
  position: relative;
}

/* 轮播容器样式 */
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

/* 轮播翻页按钮 */
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

/* 轮播指示器（底部圆点） */
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

/* 全屏图片预览弹窗 */
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

/* 预览弹窗底部导航栏 */
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
