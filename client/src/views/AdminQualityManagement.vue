<template>
  <div class="admin-quality-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="header-title">品质管理</h2>
      <select v-model="selectedProduct" class="product-select" @change="fetchQuality">
        <option value="">选择产品</option>
        <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
      </select>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载品质信息中...</p>
    </div>
    
    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button class="btn-retry" @click="fetchQuality">重试</button>
    </div>
    
    <!-- 品质信息表单 -->
    <div v-else class="quality-form">
      <div v-if="!selectedProduct" class="no-selection">
        <div class="no-selection-icon">🏆</div>
        <h3>请选择产品</h3>
        <p>选择一个产品来查看和编辑其品质信息</p>
      </div>
      
      <form v-else @submit.prevent="handleSubmit">
        <div class="form-grid">
          <div class="form-group">
            <label for="harvest_start_time" class="form-label">采收起始时间</label>
            <input 
              type="date" 
              id="harvest_start_time" 
              v-model="qualityForm.harvest_start_time" 
              class="form-input"
            >
          </div>
          
          <div class="form-group">
            <label for="harvest_end_time" class="form-label">采收终止时间</label>
            <input 
              type="date" 
              id="harvest_end_time" 
              v-model="qualityForm.harvest_end_time" 
              class="form-input"
            >
          </div>
          
          <div class="form-group">
            <label for="sugar_content" class="form-label">糖度 (°Bx)</label>
            <input 
              type="number" 
              id="sugar_content" 
              v-model.number="qualityForm.sugar_content" 
              class="form-input" 
              placeholder="请输入糖度"
              step="0.1"
            >
          </div>
          
          <div class="form-group">
            <label for="weight" class="form-label">重量 (g)</label>
            <input 
              type="number" 
              id="weight" 
              v-model.number="qualityForm.weight" 
              class="form-input" 
              placeholder="请输入重量"
              step="1"
            >
          </div>
          
          <div class="form-group full-width">
            <label for="taste_description" class="form-label">口感描述</label>
            <textarea 
              id="taste_description" 
              v-model="qualityForm.taste_description" 
              class="form-textarea" 
              placeholder="请输入口感描述"
              rows="4"
            ></textarea>
          </div>
          
          <div class="form-group full-width">
            <label for="suitable_for" class="form-label">适应人群</label>
            <textarea 
              id="suitable_for" 
              v-model="qualityForm.suitable_for" 
              class="form-textarea" 
              placeholder="请输入适应人群"
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group full-width">
            <label for="quality_summary" class="form-label">品质小结</label>
            <textarea 
              id="quality_summary" 
              v-model="qualityForm.quality_summary" 
              class="form-textarea" 
              placeholder="请输入品质小结"
              rows="5"
            ></textarea>
          </div>
        </div>
        
        <div class="form-actions">
          <button type="submit" class="btn-save" :disabled="saving">
            <span v-if="saving" class="spinner"></span>
            <span>{{ saving ? '保存中...' : '保存品质信息' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { productAPI, qualityAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const { showToast } = useToast()

// 响应式数据
const products = ref([])
const loading = ref(false)
const error = ref(null)
const selectedProduct = ref('')
const saving = ref(false)
const qualityForm = ref({
  harvest_start_time: '',
  harvest_end_time: '',
  sugar_content: 0,
  weight: 0,
  taste_description: '',
  suitable_for: '',
  quality_summary: ''
})

// 生命周期
onMounted(() => {
  loadProducts()
})

// 方法
const loadProducts = async () => {
  try {
    const response = await productAPI.list({ limit: 100 })
    products.value = response.products || []
  } catch (err) {
    console.error('加载产品列表失败', err)
  }
}

const fetchQuality = async () => {
  if (!selectedProduct.value) {
    resetForm()
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    const response = await qualityAPI.get(selectedProduct.value)
    qualityForm.value = {
      harvest_start_time: response.harvest_start_time || '',
      harvest_end_time: response.harvest_end_time || '',
      sugar_content: response.sugar_content || 0,
      weight: response.weight || 0,
      taste_description: response.taste_description || '',
      suitable_for: response.suitable_for || '',
      quality_summary: response.quality_summary || ''
    }
  } catch (err) {
    // 如果没有找到品质信息，重置表单
    if (err.message.includes('not found')) {
      resetForm()
    } else {
      error.value = err.message || '获取品质信息失败'
      console.error('获取品质信息失败', err)
    }
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  qualityForm.value = {
    harvest_start_time: '',
    harvest_end_time: '',
    sugar_content: 0,
    weight: 0,
    taste_description: '',
    suitable_for: '',
    quality_summary: ''
  }
}

const handleSubmit = async () => {
  if (!selectedProduct.value) {
    showToast('请选择产品', 'warning')
    return
  }
  
  saving.value = true
  
  try {
    const data = {
      ...qualityForm.value,
      product_id: selectedProduct.value
    }
    
    await qualityAPI.update(selectedProduct.value, data)
    showToast('品质信息保存成功', 'success')
  } catch (err) {
    showToast(err.message || '保存失败', 'error')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.admin-quality-management {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 24px;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.header-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.product-select {
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background-color: white;
  cursor: pointer;
  transition: border-color 0.3s;
  min-width: 200px;
}

.product-select:focus {
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
  flex: 1;
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
  flex: 1;
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

.no-selection {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background-color: #f9fafb;
  border-radius: 8px;
  margin: 20px 0;
  flex: 1;
}

.no-selection-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.no-selection h3 {
  color: #4b5563;
  margin-bottom: 10px;
}

.no-selection p {
  color: #6b7280;
}

.quality-form {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
}

.form-input,
.form-textarea {
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.btn-save {
  padding: 12px 24px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-save:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-save:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-save .spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  margin-bottom: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-quality-management {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .product-select {
    width: 100%;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    justify-content: center;
  }
}
</style>