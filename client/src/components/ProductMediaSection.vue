<template>
  <section class="product-section media-section">
    <div class="section-header">
      <h2 class="section-title">
        <i class="icon-media"></i> 多媒体信息
      </h2>
    </div>
    <div class="section-content">
      <div v-if="media.length === 0" class="no-media">
        <div class="no-media-icon">📷</div>
        <p>暂无多媒体内容</p>
      </div>
      <div v-else class="media-gallery">
        <div v-for="item in media" :key="item.id" class="media-item">
          <div v-if="item.media_type === 'image'" class="media-image">
            <img :src="item.file_path" :alt="item.description" @error="handleImageError">
            <div v-if="item.description" class="media-caption">{{ item.description }}</div>
          </div>
          <div v-else-if="item.media_type === 'video'" class="media-video">
            <video controls>
              <source :src="item.file_path" type="video/mp4">
              您的浏览器不支持视频播放
            </video>
            <div v-if="item.description" class="media-caption">{{ item.description }}</div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  media: {
    type: Array,
    required: true,
    default: () => []
  }
})

const handleImageError = (event) => {
  event.target.src = 'https://via.placeholder.com/400x300?text=图片加载失败'
}
</script>

<style scoped>
.media-section {
  margin-bottom: 24px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
}

.section-header {
  background: #f5f5f5;
  padding: 16px 20px;
  border-bottom: 1px solid #e0e0e0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-content {
  padding: 20px;
}

.no-media {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.no-media-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.media-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.media-item {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.media-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.media-image img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  display: block;
}

.media-caption {
  padding: 12px;
  background: #f5f5f5;
  font-size: 14px;
  color: #555;
  border-top: 1px solid #e0e0e0;
}

.media-video video {
  width: 100%;
  height: 200px;
  background: #000;
}

@media (max-width: 768px) {
  .media-gallery {
    grid-template-columns: 1fr;
  }
  
  .section-content {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .media-item {
    margin-bottom: 16px;
  }
}
</style>