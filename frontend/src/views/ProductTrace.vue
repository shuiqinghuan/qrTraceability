<template>
  <div class="product-trace">
    <header class="page-header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-icon">🌿</span>
          <h1 class="page-title">农产品溯源系统</h1>
        </div>
        <button class="back-btn" @click="goBack">
          <span>←</span>
          返回列表
        </button>
      </div>
    </header>

    <main class="main-content" v-if="!loading && !error">
      <div class="media-section">
        <MediaGallery :images="productData.images" />
      </div>

      <div class="info-section">
        <div class="info-card">
          <ProductInfo :product="productData" />
        </div>
        <div class="info-card">
          <HarvestQuality :harvest="productData.harvest" />
        </div>
      </div>
    </main>

    <main class="main-content loading-container" v-if="loading">
      <div class="loading">加载中...</div>
    </main>

    <main class="main-content error-container" v-if="error">
      <div class="error">
        <p>{{ error }}</p>
        <button @click="loadProductData">重试</button>
      </div>
    </main>

    <footer class="page-footer">
      <p>© 2024 农产品溯源系统 | 保障食品安全，追溯产品来源</p>
    </footer>
  </div>
</template>

<!-- 产品溯源详情页：展示产品的媒体图片、基本信息和采收质量数据 -->
<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ProductInfo from '../components/ProductInfo.vue'
import MediaGallery from '../components/MediaGallery.vue'
import HarvestQuality from '../components/HarvestQuality.vue'
import { getProductByCode } from '../api/product'

const route = useRoute()
const router = useRouter()

// 页面状态
const loading = ref(true)    // 是否正在加载
const error = ref(null)      // 错误信息

// 产品数据，包含基本信息、媒体图片和采收质量
const productData = ref({
  name: '',
  code: '',
  plantingLocation: '',
  plantingDate: '',
  images: [],
  harvest: {
    startDate: '',
    endDate: '',
    sugarContent: 0,
    weight: 0,
    taste: '',
    suitableCrowd: '',
    qualitySummary: ''
  }
})

/** 根据路由中的产品溯源码加载产品数据，兼容驼峰和下划线字段命名 */
const loadProductData = async () => {
  loading.value = true
  error.value = null
  
  try {
    // 从路由参数获取产品码，未提供时使用默认值'4395'
    const productCode = route.params.code || '4395'
    const response = await getProductByCode(productCode)
    
    if (response.code === 200 && response.data) {
      const data = response.data
      productData.value = {
        name: data.name || '',
        code: data.code || '',
        plantingLocation: data.plantingLocation || data.planting_location || '',
        plantingDate: data.plantingDate || data.planting_date || '',
        images: data.images || [],
        harvest: data.harvest ? {
          startDate: data.harvest.startDate || data.harvest.start_date || '',
          endDate: data.harvest.endDate || data.harvest.end_date || '',
          sugarContent: data.harvest.sugarContent || data.harvest.sugar_content || 0,
          weight: data.harvest.weight || 0,
          taste: data.harvest.taste || '',
          suitableCrowd: data.harvest.suitableCrowd || data.harvest.suitable_crowd || '',
          qualitySummary: data.harvest.qualitySummary || data.harvest.quality_summary || ''
        } : {
          startDate: '',
          endDate: '',
          sugarContent: 0,
          weight: 0,
          taste: '',
          suitableCrowd: '',
          qualitySummary: ''
        }
      }
    } else {
      error.value = response.message || '获取产品信息失败'
    }
  } catch (err) {
    console.error('加载产品数据失败:', err)
    error.value = '网络错误，请检查后端服务是否启动'
  } finally {
    loading.value = false
  }
}

/** 返回首页产品列表 */
const goBack = () => {
  router.push('/')
}

// 组件挂载后自动加载产品数据
onMounted(() => {
  loadProductData()
})
</script>

<style scoped>
.product-trace {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  padding: 20px 24px;
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  font-size: 36px;
}

.page-title {
  font-size: 28px;
  font-weight: bold;
  color: white;
  margin: 0;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border: 2px solid white;
  background: transparent;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: white;
  color: var(--primary-color);
}

.main-content {
  flex: 1;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.media-section {
  width: 100%;
}

.info-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-card {
  min-height: 300px;
}

.loading-container,
.error-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.loading {
  font-size: 18px;
  color: var(--secondary-text);
}

.error {
  text-align: center;
}

.error p {
  font-size: 16px;
  color: #f44336;
  margin-bottom: 16px;
}

.error button {
  padding: 10px 24px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.error button:hover {
  background: var(--secondary-color);
}

.page-footer {
  background: var(--text-color);
  color: white;
  text-align: center;
  padding: 16px;
  font-size: 12px;
}

@media (max-width: 1199px) {
  .info-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 16px;
  }

  .page-title {
    font-size: 22px;
  }

  .main-content {
    padding: 16px;
  }
}
</style>
