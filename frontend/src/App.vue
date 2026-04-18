<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

// 产品数据
const productData = ref(null)
const loading = ref(true)
const error = ref(null)

// 交互数据
const interactionStats = ref({
  likeCount: 0,
  shareCount: 0,
  collectCount: 0
})

// 点赞状态
const liked = ref(false)
const shared = ref(false)
const collected = ref(false)

// 获取产品ID（从URL中获取）
const getProductId = () => {
  const path = window.location.pathname
  const parts = path.split('/')
  return parts[parts.length - 1]
}

// 获取产品数据
const fetchProductData = async () => {
  const productId = getProductId()
  if (!productId) {
    error.value = '产品ID不存在'
    loading.value = false
    return
  }

  try {
    // 注意：在实际生产环境中，应该从配置中获取API地址
    const response = await axios.get(`http://localhost:8000/api/batches/unique/${productId}`)
    productData.value = response.data
    
    // 获取交互统计数据
    const statsResponse = await axios.get(`http://localhost:8000/api/interactions/batch/${productData.value.id}`)
    interactionStats.value = statsResponse.data
    
    loading.value = false
  } catch (err) {
    console.error('获取产品数据失败:', err)
    error.value = '获取产品数据失败'
    loading.value = false
  }
}

// 处理点赞
const handleLike = async () => {
  if (!productData.value) return
  
  try {
    await axios.post('http://localhost:8000/api/interactions', {
      batch_id: productData.value.id,
      action_type: 'like'
    })
    
    liked.value = true
    interactionStats.value.likeCount++
  } catch (err) {
    console.error('点赞失败:', err)
    alert('点赞失败，请稍后重试')
  }
}

// 处理转发
const handleShare = async () => {
  if (!productData.value) return
  
  try {
    await axios.post('http://localhost:8000/api/interactions', {
      batch_id: productData.value.id,
      action_type: 'share'
    })
    
    shared.value = true
    interactionStats.value.shareCount++
    
    // 模拟分享功能
    if (navigator.share) {
      await navigator.share({
        title: productData.value.product.name,
        text: '查看这个农产品的详细信息',
        url: window.location.href
      })
    } else {
      // 复制链接到剪贴板
      await navigator.clipboard.writeText(window.location.href)
      alert('链接已复制到剪贴板')
    }
  } catch (err) {
    console.error('转发失败:', err)
    alert('转发失败，请稍后重试')
  }
}

// 处理收藏
const handleCollect = async () => {
  if (!productData.value) return
  
  try {
    await axios.post('http://localhost:8000/api/interactions', {
      batch_id: productData.value.id,
      action_type: 'collect'
    })
    
    collected.value = true
    interactionStats.value.collectCount++
  } catch (err) {
    console.error('收藏失败:', err)
    alert('收藏失败，请稍后重试')
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchProductData()
})
</script>

<template>
  <div class="product-detail">
    <!-- 加载状态 -->
    <van-loading v-if="loading" type="spinner" color="#1989fa" />
    
    <!-- 错误状态 -->
    <van-empty v-else-if="error" :description="error" />
    
    <!-- 产品详情 -->
    <div v-else-if="productData" class="content">
      <!-- 第一部分：产品基本信息 -->
      <section class="section section-info">
        <h1 class="product-name">{{ productData.product.name }}</h1>
        <div class="info-item">
          <span class="label">定植地点：</span>
          <span class="value">{{ productData.planting_location }}</span>
        </div>
        <div class="info-item">
          <span class="label">定植时间：</span>
          <span class="value">{{ new Date(productData.planting_date).toLocaleDateString() }}</span>
        </div>
      </section>
      
      <!-- 第二部分：产品图片和视频 -->
      <section class="section section-media">
        <h2 class="section-title">产品展示</h2>
        <div class="media-list">
          <div v-for="(media, index) in productData.media_files" :key="index" class="media-item">
            <img v-if="media.type === 'image'" :src="media.url" :alt="productData.product.name" class="media-image" />
            <video v-else-if="media.type === 'video'" :src="media.url" controls class="media-video"></video>
          </div>
        </div>
      </section>
      
      <!-- 第三部分：采收质量信息 -->
      <section class="section section-harvest">
        <h2 class="section-title">采收质量信息</h2>
        <div v-if="productData.harvest_quality" class="harvest-info">
          <div class="info-item">
            <span class="label">采收起始时间：</span>
            <span class="value">{{ new Date(productData.harvest_quality.harvest_start_date).toLocaleDateString() }}</span>
          </div>
          <div class="info-item">
            <span class="label">采收终止时间：</span>
            <span class="value">{{ new Date(productData.harvest_quality.harvest_end_date).toLocaleDateString() }}</span>
          </div>
          <div class="info-item">
            <span class="label">糖度：</span>
            <span class="value">{{ productData.harvest_quality.sugar_content }}%</span>
          </div>
          <div class="info-item">
            <span class="label">重量：</span>
            <span class="value">{{ productData.harvest_quality.weight }}g</span>
          </div>
          <div class="info-item">
            <span class="label">口感：</span>
            <span class="value">{{ productData.harvest_quality.taste }}</span>
          </div>
          <div class="info-item">
            <span class="label">适应人群：</span>
            <span class="value">{{ productData.harvest_quality.suitable_for }}</span>
          </div>
          <div class="info-item">
            <span class="label">品质小结：</span>
            <span class="value">{{ productData.harvest_quality.quality_summary }}</span>
          </div>
        </div>
      </section>
      
      <!-- 第四部分：用户交互反馈 -->
      <section class="section section-interaction">
        <h2 class="section-title">用户反馈</h2>
        <div class="interaction-buttons">
          <van-button 
            type="primary" 
            :plain="!liked" 
            @click="handleLike"
            class="interaction-btn"
          >
            <van-icon name="like" />
            <span>{{ interactionStats.likeCount }}</span>
          </van-button>
          <van-button 
            type="primary" 
            :plain="!shared" 
            @click="handleShare"
            class="interaction-btn"
          >
            <van-icon name="share" />
            <span>{{ interactionStats.shareCount }}</span>
          </van-button>
          <van-button 
            type="primary" 
            :plain="!collected" 
            @click="handleCollect"
            class="interaction-btn"
          >
            <van-icon name="star" />
            <span>{{ interactionStats.collectCount }}</span>
          </van-button>
        </div>
      </section>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: #f5f5f5;
  color: #333;
  line-height: 1.6;
}

.product-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #fff;
  min-height: 100vh;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.section {
  padding: 20px;
  border-radius: 8px;
  background-color: #f9f9f9;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 15px;
  color: #1989fa;
}

.product-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #333;
}

.info-item {
  display: flex;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.label {
  font-weight: 500;
  min-width: 120px;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
}

.media-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.media-item {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.media-image {
  width: 100%;
  height: auto;
  display: block;
}

.media-video {
  width: 100%;
  height: auto;
  display: block;
}

.interaction-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 10px;
}

.interaction-btn {
  min-width: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .product-detail {
    padding: 10px;
  }
  
  .section {
    padding: 15px;
  }
  
  .product-name {
    font-size: 20px;
  }
  
  .section-title {
    font-size: 16px;
  }
  
  .interaction-buttons {
    flex-direction: column;
    gap: 10px;
  }
  
  .interaction-btn {
    width: 100%;
  }
}
</style>
