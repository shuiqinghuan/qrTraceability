<template>
  <div class="home">
    <header class="page-header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-icon">🌿</span>
          <h1 class="page-title">农产品溯源系统</h1>
        </div>
      </div>
    </header>

    <main class="main-content">
      <div class="hero-section">
        <h2>追溯每一份新鲜</h2>
        <p>扫描二维码，查询产品溯源信息，了解从田间到餐桌的全过程</p>
      </div>

      <div class="search-section" v-if="products.length > 1">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索产品名称或编号..."
            @input="handleSearch"
          />
        </div>
      </div>

      <div class="products-section">
        <h3 class="section-title">
          <span class="icon">📦</span>
          产品列表
          <span class="product-count">(共 {{ total }} 个产品)</span>
        </h3>

        <div class="products-grid" v-if="!loading && products.length > 0">
          <router-link
            v-for="product in products"
            :key="product.id"
            :to="`/product/${product.code}`"
            class="product-card"
          >
            <div class="product-image">
              <span class="product-icon">🍎</span>
            </div>
            <div class="product-info">
              <h4 class="product-name">{{ product.name }}</h4>
              <p class="product-code">编号: {{ product.code }}</p>
              <p class="product-location">📍 {{ product.plantingLocation || product.planting_location }}</p>
              <p class="product-date">📅 {{ product.plantingDate || product.planting_date }}</p>
            </div>
            <div class="product-arrow">→</div>
          </router-link>
        </div>

        <div class="loading-container" v-if="loading">
          <div class="loading">加载中...</div>
        </div>

        <div class="empty-container" v-if="!loading && products.length === 0">
          <div class="empty-icon">📭</div>
          <p class="empty-text">暂无产品数据</p>
        </div>

        <div class="pagination" v-if="!loading && total > pageSize">
          <button
            @click="prevPage"
            :disabled="currentPage <= 1"
            class="page-btn"
          >
            上一页
          </button>
          <span class="page-info">
            第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
          </span>
          <button
            @click="nextPage"
            :disabled="currentPage >= totalPages"
            class="page-btn"
          >
            下一页
          </button>
        </div>
      </div>
    </main>

    <footer class="page-footer">
      <p>© 2024 农产品溯源系统 | 保障食品安全，追溯产品来源</p>
    </footer>
  </div>
</template>

<!-- 首页视图：展示产品列表，支持搜索和分页功能 -->
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProductList } from '../api/product'

const router = useRouter()

// 响应式状态
const products = ref([])       // 产品列表数据
const total = ref(0)           // 产品总数
const currentPage = ref(1)     // 当前页码
const pageSize = ref(1)        // 每页显示数量
const loading = ref(false)     // 加载状态
const searchQuery = ref('')    // 搜索关键词

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

/** 加载产品列表数据 */
const loadProducts = async (page = 1) => {
  loading.value = true
  try {
    const response = await getProductList({
      page,
      pageSize: pageSize.value
    })
    if (response.code === 200) {
      products.value = response.data.list || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    console.error('加载产品列表失败:', error)
  } finally {
    loading.value = false
  }
}

/** 搜索处理：匹配产品名称或编号，匹配成功跳转到溯源详情页 */
const handleSearch = () => {
  if (searchQuery.value) {
    const found = products.value.find(p =>
      p.name.includes(searchQuery.value) ||
      p.code.toString().includes(searchQuery.value)
    )
    if (found) {
      router.push(`/product/${found.code}`)
    }
  }
}

/** 上一页 */
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadProducts(currentPage.value)
  }
}

/** 下一页 */
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadProducts(currentPage.value)
  }
}

// 组件挂载后加载首页产品数据
onMounted(() => {
  loadProducts(currentPage.value)
})
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 20px 0;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 28px;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.hero-section {
  text-align: center;
  margin-bottom: 40px;
}

.hero-section h2 {
  font-size: 32px;
  color: var(--primary-color);
  margin-bottom: 12px;
}

.hero-section p {
  font-size: 16px;
  color: #666;
}

.search-section {
  margin-bottom: 30px;
}

.search-box {
  position: relative;
  max-width: 500px;
  margin: 0 auto;
}

.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 18px;
  opacity: 0.5;
}

.search-box input {
  width: 100%;
  padding: 14px 14px 14px 48px;
  border: 2px solid #e0e0e0;
  border-radius: 12px;
  font-size: 16px;
  transition: all 0.3s;
  background: white;
}

.search-box input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(76, 175, 80, 0.1);
}

.products-section {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  color: #333;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f0f0;
}

.section-title .icon {
  font-size: 24px;
}

.product-count {
  font-size: 14px;
  color: #999;
  font-weight: normal;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.product-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, #f8fffe 0%, #f0f9ff 100%);
  border: 2px solid #e0f2f1;
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.3s;
  cursor: pointer;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(76, 175, 80, 0.15);
  border-color: var(--primary-color);
}

.product-image {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #e8f5e9, #c8e6c9);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.product-icon {
  font-size: 28px;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 6px;
}

.product-code {
  font-size: 13px;
  color: var(--primary-color);
  font-weight: 500;
  margin-bottom: 4px;
}

.product-location,
.product-date {
  font-size: 12px;
  color: #888;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-arrow {
  font-size: 20px;
  color: #ccc;
  transition: all 0.3s;
}

.product-card:hover .product-arrow {
  color: var(--primary-color);
  transform: translateX(4px);
}

.loading-container,
.empty-container {
  text-align: center;
  padding: 60px 20px;
}

.loading {
  font-size: 16px;
  color: #666;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #999;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.page-btn {
  padding: 10px 24px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  background: var(--secondary-color);
  transform: translateY(-2px);
}

.page-btn:disabled {
  background: #ddd;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #666;
}

.page-footer {
  text-align: center;
  padding: 24px;
  color: #888;
  font-size: 13px;
}
</style>
