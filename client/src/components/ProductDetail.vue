<template>
  <div class="product-detail">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">{{ product.seed_info.name }} - 产品详情</h1>
      <div class="page-actions">
        <button class="btn-back" @click="goBack">
          <i class="icon-back"></i> 返回列表
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载产品信息中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button class="btn-retry" @click="loadProductDetail">重试</button>
    </div>

    <!-- 产品详情内容 -->
    <div v-else class="product-content">
      <!-- 第一部分：品种信息 -->
      <ProductVarietySection :seed-info="product.seed_info" />

      <!-- 第二部分：种植信息 -->
      <section class="product-section planting-section">
        <div class="section-header">
          <h2 class="section-title">
            <i class="icon-plant"></i> 种植信息
          </h2>
        </div>
        <div class="section-content">
          <div class="planting-grid">
            <div class="planting-item">
              <label class="planting-label">定植地点</label>
              <div class="planting-value">
                {{ product.planting.location }}
                <div class="location-tags">
                  <span v-for="tag in product.tags" :key="tag.id" class="location-tag">
                    {{ tag.tag_name }}
                  </span>
                </div>
              </div>
            </div>
            <div class="planting-item">
              <label class="planting-label">定植时间</label>
              <div class="planting-value">{{ formatDate(product.planting.planting_date) }}</div>
            </div>
            <div class="planting-item">
              <label class="planting-label">移栽时间</label>
              <div class="planting-value">
                {{ product.planting.transplanting_date ? formatDate(product.planting.transplanting_date) : '未移栽' }}
              </div>
            </div>
            <div class="planting-item">
              <label class="planting-label">种植备注</label>
              <div class="planting-notes">{{ product.planting.notes || '暂无备注' }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 第三部分：多媒体信息 -->
      <ProductMediaSection :media="product.media" />

      <!-- 第四部分：品质信息 -->
      <ProductQualitySection :quality="product.quality" />

      <!-- 第五部分：互动模块 -->
      <section class="product-section interaction-section">
        <div class="interaction-grid">
          <!-- 点赞模块 -->
          <LikeModule 
            :like-status="likeStatus"
            :like-count="product.like_count"
            @like="handleLike"
          />

          <!-- 跳转模块 -->
          <RelatedProductsSection 
            :related-products="relatedProducts"
            @view-related="goToProduct"
          />
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productAPI, locationAPI } from '../services/api'
import { useToast } from '../composables/useToast'
import ProductVarietySection from './ProductVarietySection.vue'
import ProductMediaSection from './ProductMediaSection.vue'
import ProductQualitySection from './ProductQualitySection.vue'
import LikeModule from './LikeModule.vue'
import RelatedProductsSection from './RelatedProductsSection.vue'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()

// 响应式数据
const product = ref({
  id: 0,
  seed_info: {
    id: 0,
    name: '',
    variety: '',
    variety_code: '',
    description: ''
  },
  planting: {
    id: 0,
    planting_date: '',
    transplanting_date: null,
    location: '',
    notes: ''
  },
  media: [],
  quality: {
    id: 0,
    harvest_start_date: null,
    harvest_end_date: null,
    sugar_content: 0,
    weight: 0,
    taste_description: '',
    suitable_for: '',
    quality_summary: ''
  },
  tags: [],
  like_count: 0,
  favorite_count: 0
})

const relatedProducts = ref([])
const likeStatus = ref({
  can_like: true,
  cooldown: 0,
  user_liked: false,
  ip_address: ''
})
const loading = ref(true)
const error = ref(null)

// 计算属性
const productId = computed(() => parseInt(route.params.id))

// 生命周期
onMounted(() => {
  loadProductDetail()
  loadRelatedProducts()
  loadLikeStatus()
})

// 方法
const loadProductDetail = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await productAPI.getDetail(productId.value)
    product.value = response.product
  } catch (err) {
    error.value = err.message || '加载产品详情失败'
    console.error('加载产品详情失败:', err)
  } finally {
    loading.value = false
  }
}

const loadRelatedProducts = async () => {
  try {
    const response = await productAPI.getRelatedProducts(productId.value)
    relatedProducts.value = response.related_products || []
  } catch (err) {
    console.error('加载相关产品失败:', err)
    relatedProducts.value = []
  }
}

const loadLikeStatus = async () => {
  try {
    const response = await productAPI.getLikeStatus(productId.value)
    likeStatus.value = {
      can_like: response.response.can_like,
      cooldown: response.response.cooldown_seconds || 0,
      user_liked: response.user_liked || false,
      ip_address: response.ip_address || ''
    }
  } catch (err) {
    console.error('加载点赞状态失败:', err)
    likeStatus.value = {
      can_like: true,
      cooldown: 0,
      user_liked: false,
      ip_address: ''
    }
  }
}

const handleLike = async () => {
  if (!likeStatus.value.can_like) {
    showToast('点赞过于频繁，请稍后再试', 'warning')
    return
  }

  try {
    const response = await productAPI.likeProduct(productId.value)
    
    if (response.response.success) {
      // 更新点赞状态
      likeStatus.value.user_liked = true
      likeStatus.value.can_like = false
      likeStatus.value.cooldown = 300 // 5分钟冷却
      
      // 更新点赞计数
      product.value.like_count = response.response.like_count
      
      showToast('点赞成功!', 'success')
      
      // 重新加载点赞状态
      setTimeout(() => {
        loadLikeStatus()
      }, 1000)
    } else {
      showToast(response.response.message, 'warning')
    }
  } catch (err) {
    showToast('点赞失败: ' + (err.message || '未知错误'), 'error')
    console.error('点赞失败:', err)
  }
}

const goToProduct = (id) => {
  router.push(`/products/${id}`)
}

const goBack = () => {
  router.back()
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatCooldown = (seconds) => {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

const handleImageError = (event) => {
  event.target.src = 'https://via.placeholder.com/400x300?text=图片加载失败'
}
</script>

<style scoped>
.product-detail {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.page-header {
  background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
  color: white;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.btn-back {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s;
}

.btn-back:hover {
  background: rgba(255, 255, 255, 0.3);
}

.loading-state,
.error-state {
  padding: 60px 24px;
  text-align: center;
}

.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #4CAF50;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.btn-retry {
  background: #4CAF50;
  border: none;
  color: white;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 16px;
  transition: background 0.2s;
}

.btn-retry:hover {
  background: #45a049;
}

.product-content {
  padding: 24px;
}

.product-section {
  margin-bottom: 32px;
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

/* 品种信息样式 */
.variety-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.variety-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.variety-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.variety-value {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  padding: 8px 0;
}

.variety-description {
  font-size: 14px;
  color: #555;
  line-height: 1.5;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
  border-left: 4px solid #4CAF50;
}

/* 种植信息样式 */
.planting-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.planting-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.planting-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.planting-value {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  padding: 8px 0;
}

.planting-notes {
  font-size: 14px;
  color: #555;
  line-height: 1.5;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
  border-left: 4px solid #2196F3;
}

.location-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.location-tag {
  background: #e3f2fd;
  color: #1976d2;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

/* 多媒体信息样式 */
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

/* 品质信息样式 */
.quality-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.quality-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.quality-item.full-width {
  grid-column: 1 / -1;
}

.quality-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.quality-value {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  padding: 8px 0;
}

.quality-description,
.quality-summary {
  font-size: 14px;
  color: #555;
  line-height: 1.5;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
}

.quality-description {
  border-left: 4px solid #FF9800;
}

.quality-summary {
  border-left: 4px solid #9C27B0;
}

/* 互动模块样式 */
.interaction-section {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 24px;
}

.interaction-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.interaction-item {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.interaction-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid #4CAF50;
}

.btn-like {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f5f5f5;
  border: 2px solid #ddd;
  color: #666;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 12px;
}

.btn-like:hover:not(.disabled) {
  background: #ffebee;
  border-color: #f44336;
  color: #f44336;
}

.btn-like.liked {
  background: #ffebee;
  border-color: #f44336;
  color: #f44336;
}

.btn-like.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.like-count {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.like-cooldown {
  font-size: 12px;
  color: #ff9800;
  font-weight: 500;
}

.no-related {
  text-align: center;
  padding: 20px;
  color: #999;
}

.related-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.related-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 6px;
  border-left: 4px solid #4CAF50;
}

.related-name {
  font-weight: 500;
  color: #333;
}

.related-variety {
  font-size: 12px;
  color: #666;
}

.btn-view-related {
  background: #4CAF50;
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-view-related:hover {
  background: #45a049;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }

  .variety-grid,
  .planting-grid,
  .quality-grid {
    grid-template-columns: 1fr;
  }

  .interaction-grid {
    grid-template-columns: 1fr;
  }

  .media-gallery {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .product-content {
    padding: 16px;
  }

  .section-content {
    padding: 16px;
  }

  .interaction-section {
    padding: 16px;
  }
}
</style>