<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { config } from '../config.js'

const router = useRouter()
const batches = ref([])
const loading = ref(true)
const error = ref(null)

// 获取批次列表
const fetchBatches = async () => {
  try {
    const response = await axios.get(`${config.apiBaseUrl}/api/batches`)
    batches.value = response.data
    loading.value = false
  } catch (err) {
    console.error('获取批次列表失败:', err)
    error.value = '获取批次列表失败'
    loading.value = false
  }
}

// 删除批次
const deleteBatch = async (id) => {
  if (confirm('确定要删除这个批次吗？')) {
    try {
      await axios.delete(`${config.apiBaseUrl}/api/batches/${id}`)
      batches.value = batches.value.filter(batch => batch.id !== id)
      alert('删除成功')
    } catch (err) {
      console.error('删除批次失败:', err)
      alert('删除失败，请稍后重试')
    }
  }
}

// 跳转到添加批次页面
const goToAddBatch = () => {
  router.push('/admin/batches/add')
}

// 页面加载时获取数据
onMounted(() => {
  fetchBatches()
})
</script>

<template>
  <div class="admin-batches">
    <!-- 页面头部 -->
    <div class="page-header">
      <h3>批次列表</h3>
      <van-button type="primary" @click="goToAddBatch">
        <van-icon name="plus" />
        添加批次
      </van-button>
    </div>
    
    <!-- 加载状态 -->
    <van-loading v-if="loading" type="spinner" color="#1989fa" />
    
    <!-- 错误状态 -->
    <van-empty v-else-if="error" :description="error" />
    
    <!-- 批次列表 -->
    <div v-else class="batches-list">
      <van-card
        v-for="batch in batches"
        :key="batch.id"
        class="batch-card"
      >
        <template #header>
          <div class="card-header">
            <h4 class="batch-name">{{ batch.product?.name || '未知产品' }} - 批次 {{ batch.id }}</h4>
            <div class="card-actions">
              <van-button size="small" type="primary" plain>
                编辑
              </van-button>
              <van-button size="small" type="danger" plain @click="deleteBatch(batch.id)">
                删除
              </van-button>
            </div>
          </div>
        </template>
        <div class="batch-info">
          <div class="info-item">
            <span class="label">唯一ID：</span>
            <span class="value">{{ batch.unique_id }}</span>
          </div>
          <div class="info-item">
            <span class="label">定植地点：</span>
            <span class="value">{{ batch.planting_location }}</span>
          </div>
          <div class="info-item">
            <span class="label">定植时间：</span>
            <span class="value">{{ new Date(batch.planting_date).toLocaleDateString() }}</span>
          </div>
          <div class="info-item">
            <span class="label">创建时间：</span>
            <span class="value">{{ new Date(batch.created_at).toLocaleString() }}</span>
          </div>
        </div>
      </van-card>
      
      <!-- 空状态 -->
      <van-empty v-if="batches.length === 0" description="暂无批次" />
    </div>
  </div>
</template>

<style scoped>
.admin-batches {
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

.batches-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.batch-card {
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

.batch-name {
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

.batch-info {
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
  min-width: 100px;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .batches-list {
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
