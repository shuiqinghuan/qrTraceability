<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { config } from '../config.js'

const router = useRouter()
const products = ref([])
const loading = ref(true)
const error = ref(null)

// 获取产品列表
const fetchProducts = async () => {
  try {
    const response = await axios.get(`${config.apiBaseUrl}/api/products`)
    products.value = response.data
    loading.value = false
  } catch (err) {
    console.error('获取产品列表失败:', err)
    error.value = '获取产品列表失败'
    loading.value = false
  }
}

// 删除产品
const deleteProduct = async (id) => {
  if (confirm('确定要删除这个产品吗？')) {
    try {
      await axios.delete(`${config.apiBaseUrl}/api/products/${id}`)
      products.value = products.value.filter(product => product.id !== id)
      alert('删除成功')
    } catch (err) {
      console.error('删除产品失败:', err)
      alert('删除失败，请稍后重试')
    }
  }
}

// 跳转到添加产品页面
const goToAddProduct = () => {
  router.push('/admin/products/add')
}

// 页面加载时获取数据
onMounted(() => {
  fetchProducts()
})
</script>

<template>
  <div class="admin-products">
    <!-- 页面头部 -->
    <div class="page-header">
      <h3>产品列表</h3>
      <van-button type="primary" @click="goToAddProduct">
        <van-icon name="plus" />
        添加产品
      </van-button>
    </div>
    
    <!-- 加载状态 -->
    <van-loading v-if="loading" type="spinner" color="#1989fa" />
    
    <!-- 错误状态 -->
    <van-empty v-else-if="error" :description="error" />
    
    <!-- 产品列表 -->
    <div v-else class="products-list">
      <van-card
        v-for="product in products"
        :key="product.id"
        class="product-card"
      >
        <template #header>
          <div class="card-header">
            <h4 class="product-name">{{ product.name }}</h4>
            <div class="card-actions">
              <van-button size="small" type="primary" plain>
                编辑
              </van-button>
              <van-button size="small" type="danger" plain @click="deleteProduct(product.id)">
                删除
              </van-button>
            </div>
          </div>
        </template>
        <div class="product-info">
          <div class="info-item">
            <span class="label">产品代码：</span>
            <span class="value">{{ product.code }}</span>
          </div>
          <div class="info-item">
            <span class="label">描述：</span>
            <span class="value">{{ product.description }}</span>
          </div>
          <div class="info-item">
            <span class="label">创建时间：</span>
            <span class="value">{{ new Date(product.created_at).toLocaleString() }}</span>
          </div>
        </div>
      </van-card>
      
      <!-- 空状态 -->
      <van-empty v-if="products.length === 0" description="暂无产品" />
    </div>
  </div>
</template>

<style scoped>
.admin-products {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e5e5e5;
}

.page-header h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.products-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.product-card {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 15px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: #333;
  flex: 1;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  flex-wrap: wrap;
}

.label {
  font-weight: 500;
  min-width: 80px;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .products-list {
    grid-template-columns: 1fr;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .card-actions {
    align-self: flex-end;
  }
}
</style>
