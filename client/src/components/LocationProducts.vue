<template>
  <div class="location-products">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">{{ location }} - 产品列表</h1>
      <div class="page-actions">
        <button class="btn-back" @click="goBack">
          <i class="icon-back"></i> 返回
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载产品列表中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button class="btn-retry" @click="loadProducts">重试</button>
    </div>

    <!-- 产品列表 -->
    <div v-else class="products-content">
      <div v-if="products.length === 0" class="no-products">
        <div class="no-products-icon">📦</div>
        <h3>暂无产品</h3>
        <p>该地点暂无产品记录</p>
      </div>

      <div v-else class="products-grid">
        <div v-for="product in products" :key="product.id" class="product-card">
          <div class="product-image">
            <img :src="product.image_url || 'https://via.placeholder.com/200x150?text=产品图片'" 
                 :alt="product.name"
                 @error="handleImageError">
          </div>
          <div class="product-info">
            <h3 class="product-name">{{ product.name }}</h3>
            <div class="product-variety">{{ product.variety }}</div>
            <div class="product-location">{{ product.location }}</div>
            <div class="product-stats">
              <span class="stat-item">
                <i class="icon-sugar"></i>
                {{ product.sugar_content || '0' }}°Bx
              </span>
              <span class="stat-item">
                <i class="icon-weight"></i>
                {{ product.weight || '0' }}g
              </span>
            </div>
            <button class="btn-view-product" @click="goToProduct(product.id)">
              查看详情
            </button>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="pagination">
        <button 
          class="btn-prev" 
          :disabled="currentPage === 1"
          @click="prevPage"
        >
          上一页
        </button>
        <span class="page-info">
          第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
        </span>
        <button 
          class="btn-next" 
          :disabled="currentPage === totalPages"
          @click="nextPage"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { locationAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()

// 响应式数据
const products = ref([])
const loading = ref(true)
const error = ref(null)
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)

// 计算属性
const location = computed(() => route.params.location)
const page = computed(() => parseInt(route.query.page) || 1)
const limit = computed(() => parseInt(route.query.limit) || 10)

// 生命周期
onMounted(() => {
  loadProducts()
})

// 监听路由变化
watch(
  () => [route.params.location, route.query.page, route.query.limit],
  () => {
    loadProducts()
  },
  { deep: true }
)

// 方法
const loadProducts = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await locationAPI.getProductsByLocation(
      location.value,
      {
        page: page.value,
        limit: limit.value
      }
    )
    
    products.value = response.products || []
    totalCount.value = response.total || 0
    totalPages.value = response.total_pages || 1
    currentPage.value = page.value
  } catch (err) {
    error.value = err.message || '加载产品列表失败'
    console.error('加载产品列表失败:', err)
  } finally {
    loading.value = false
  }
}

const goToProduct = (id) => {
  router.push(`/products/${id}`)
}

const goBack = () => {
  router.back()
}

const prevPage = () => {
  if (currentPage.value > 1) {
    router.push({
      query: {
        ...route.query,
        page: currentPage.value - 1
      }
    })
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    router.push({
      query: {
        ...route.query,
        page: currentPage.value + 1
      }
    })
  }
}

const handleImageError = (event) => {
  event.target.src = 'https://via.placeholder.com/200x150?text=图片加载失败'
}
</script>

<style scoped>
.location-products {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.page-header {
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
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
  border-top: 3px solid #2196F3;
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
  background: #2196F3;
  border: none;
  color: white;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 16px;
  transition: background 0.2s;
}

.btn-retry:hover {
  background: #1976D2;
}

.products-content {
  padding: 24px;
}

.no-products {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.no-products-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.product-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.product-image {
  height: 150px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-info {
  padding: 16px;
}

.product-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.product-variety {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.product-location {
  font-size: 12px;
  color: #999;
  margin-bottom: 12px;
}

.product-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #666;
}

.btn-view-product {
  width: 100%;
  background: #2196F3;
  border: none;
  color: white;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-view-product:hover {
  background: #1976D2;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
  padding: 20px;
  border-top: 1px solid #e0e0e0;
}

.btn-prev,
.btn-next {
  background: #2196F3;
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-prev:hover:not(:disabled),
.btn-next:hover:not(:disabled) {
  background: #1976D2;
}

.btn-prev:disabled,
.btn-next:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    text-align: center;
  }

  .products-grid {
    grid-template-columns: 1fr;
  }

  .pagination {
    flex-direction: column;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .products-content {
    padding: 16px;
  }

  .product-card {
    margin-bottom: 16px;
  }
}
</style>