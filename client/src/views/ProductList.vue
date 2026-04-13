<template>
  <div class="product-list">
    <div class="page-header">
      <h2 class="page-title">产品浏览</h2>
      <div class="page-filters">
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="搜索产品..." 
          class="search-input"
          @input="handleSearch"
        >
        <select v-model="selectedLocation" class="location-select" @change="handleFilter">
          <option value="">全部地点</option>
          <option v-for="loc in locations" :key="loc" :value="loc">{{ loc }}</option>
        </select>
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
      <button class="btn-retry" @click="fetchProducts">重试</button>
    </div>
    
    <!-- 产品列表 -->
    <div v-else class="products-content">
      <div v-if="products.length === 0" class="no-products">
        <div class="no-products-icon">📦</div>
        <h3>暂无产品</h3>
        <p v-if="searchQuery || selectedLocation">没有找到匹配的产品</p>
        <p v-else>暂无产品记录</p>
      </div>
      
      <div v-else class="product-grid">
        <div v-for="product in products" :key="product.id" class="product-card">
          <div class="product-image">
            <img :src="getProductImage(product)" :alt="product.name" @error="handleImageError">
          </div>
          <div class="product-info">
            <h3 class="product-name" @click="goToProduct(product.id)">{{ product.name }}</h3>
            <div class="product-variety">{{ product.variety }}</div>
            <div class="product-location">
              <i class="icon-location"></i>
              {{ product.location }}
              <span class="view-location" @click="goToLocationProducts(product.location)">
                查看同地点产品
              </span>
            </div>
            <div class="product-meta">
              <span class="product-sugar">
                <i class="icon-sugar"></i>
                糖度: {{ product.sugar_content || '0' }}°Bx
              </span>
              <span class="product-weight">
                <i class="icon-weight"></i>
                重量: {{ product.weight || '0' }}g
              </span>
            </div>
            <div class="product-stats">
              <span class="stat-item">
                <i class="icon-like"></i>
                {{ product.like_count || 0 }} 点赞
              </span>
              <span class="stat-item">
                <i class="icon-favorite"></i>
                {{ product.favorite_count || 0 }} 收藏
              </span>
            </div>
            <div class="product-actions">
              <button class="btn-view-detail" @click="goToProduct(product.id)">
                查看详情
              </button>
            </div>
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
          第 {{ currentPage }} 页 / 共 {{ totalPages }} 页 ({{ totalCount }} 条记录)
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
import { useRouter } from 'vue-router'
import { productAPI, locationAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const router = useRouter()
const { showToast } = useToast()

// 响应式数据
const products = ref([])
const loading = ref(true)
const error = ref(null)
const searchQuery = ref('')
const selectedLocation = ref('')
const locations = ref([])
const currentPage = ref(1)
const totalPages = ref(1)
const totalCount = ref(0)

// 计算属性
const page = computed(() => currentPage.value)
const limit = computed(() => 12)

// 生命周期
onMounted(() => {
  fetchProducts()
  loadLocations()
})

// 监听搜索和过滤变化
watch([searchQuery, selectedLocation], () => {
  currentPage.value = 1
  fetchProducts()
}, { deep: true })

// 方法
const fetchProducts = async () => {
  loading.value = true
  error.value = null
  
  try {
    const params = {
      page: page.value,
      limit: limit.value
    }
    
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    
    if (selectedLocation.value) {
      params.location = selectedLocation.value
    }
    
    const response = await productAPI.list(params)
    
    products.value = response.products || []
    totalCount.value = response.total || 0
    totalPages.value = response.total_pages || 1
  } catch (err) {
    error.value = err.message || '获取产品列表失败'
    console.error('获取产品列表失败', err)
  } finally {
    loading.value = false
  }
}

const loadLocations = async () => {
  try {
    const response = await locationAPI.searchLocations('', { limit: 50 })
    locations.value = response.locations.map(loc => loc.location) || []
  } catch (err) {
    console.error('加载地点列表失败', err)
  }
}

const handleSearch = () => {
  // 防抖处理
  clearTimeout(searchQuery.timeout)
  searchQuery.timeout = setTimeout(() => {
    fetchProducts()
  }, 500)
}

const handleFilter = () => {
  fetchProducts()
}

const goToProduct = (id) => {
  router.push(`/products/${id}`)
}

const goToLocationProducts = (location) => {
  router.push(`/locations/${location}/products`)
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchProducts()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchProducts()
  }
}

const getProductImage = (product) => {
  // 这里可以根据产品信息返回不同的图片
  // 实际应用中应该从产品媒体中获取第一张图片
  return 'https://via.placeholder.com/300x200?text=' + encodeURIComponent(product.name)
}

const handleImageError = (event) => {
  event.target.src = 'https://via.placeholder.com/300x200?text=图片加载失败'
}
</script>

<style scoped>
.product-list {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 20px auto;
  max-width: 1200px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 15px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.page-filters {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  width: 300px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.location-select {
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background-color: white;
  cursor: pointer;
  transition: border-color 0.3s;
}

.location-select:focus {
  outline: none;
  border-color: #3b82f6;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background-color: #fef2f2;
  border-radius: 8px;
  margin: 20px 0;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.error-state h3 {
  color: #dc2626;
  margin-bottom: 10px;
}

.error-state p {
  color: #6b7280;
  margin-bottom: 20px;
}

.btn-retry {
  padding: 10px 20px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-retry:hover {
  background-color: #2563eb;
}

.no-products {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background-color: #f9fafb;
  border-radius: 8px;
  margin: 20px 0;
}

.no-products-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.no-products h3 {
  color: #4b5563;
  margin-bottom: 10px;
}

.no-products p {
  color: #6b7280;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
  margin: 20px 0;
}

.product-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  background-color: white;
}

.product-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.product-image {
  height: 200px;
  overflow: hidden;
  position: relative;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-info {
  padding: 20px;
}

.product-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  cursor: pointer;
  transition: color 0.3s;
}

.product-name:hover {
  color: #3b82f6;
}

.product-variety {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.product-location {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.icon-location {
  color: #f97316;
  font-size: 16px;
}

.view-location {
  margin-left: auto;
  font-size: 12px;
  color: #3b82f6;
  cursor: pointer;
  text-decoration: underline;
  transition: color 0.3s;
}

.view-location:hover {
  color: #2563eb;
}

.product-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #6b7280;
}

.product-sugar, .product-weight {
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-sugar {
  color: #f59e0b;
  font-size: 16px;
}

.icon-weight {
  color: #10b981;
  font-size: 16px;
}

.product-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  font-size: 14px;
  color: #6b7280;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.icon-like {
  color: #ef4444;
  font-size: 16px;
}

.icon-favorite {
  color: #8b5cf6;
  font-size: 16px;
}

.product-actions {
  display: flex;
  gap: 10px;
}

.btn-view-detail {
  flex: 1;
  padding: 12px 20px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.btn-view-detail:hover {
  background-color: #2563eb;
  transform: translateY(-2px);
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.btn-prev, .btn-next {
  padding: 8px 16px;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-prev:hover:not(:disabled), .btn-next:hover:not(:disabled) {
  background-color: #f3f4f6;
  border-color: #d1d5db;
}

.btn-prev:disabled, .btn-next:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #6b7280;
  min-width: 200px;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .product-list {
    padding: 15px;
    margin: 10px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    text-align: center;
  }

  .search-input {
    width: 100%;
  }

  .page-filters {
    flex-direction: column;
    align-items: stretch;
  }

  .product-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .product-info {
    padding: 15px;
  }

  .pagination {
    flex-direction: column;
    gap: 10px;
  }

  .page-info {
    order: -1;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>