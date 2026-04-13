<template>
  <div class="admin-product-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="header-title">产品管理</h2>
      <button class="btn-add-product" @click="openAddProductModal">
        <i class="add-icon">➕</i>
        <span>添加产品</span>
      </button>
    </div>
    
    <!-- 搜索和过滤 -->
    <div class="search-filter">
      <input 
        type="text" 
        v-model="searchQuery" 
        class="search-input" 
        placeholder="搜索产品..."
        @input="handleSearch"
      >
      <select v-model="selectedLocation" class="filter-select" @change="fetchProducts">
        <option value="">全部地点</option>
        <option v-for="loc in locations" :key="loc" :value="loc">{{ loc }}</option>
      </select>
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
    <div v-else class="product-list">
      <div v-if="products.length === 0" class="no-products">
        <div class="no-products-icon">📦</div>
        <h3>暂无产品</h3>
        <p v-if="searchQuery || selectedLocation">没有找到匹配的产品</p>
        <p v-else>暂无产品记录</p>
      </div>
      
      <table v-else class="product-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>产品名称</th>
            <th>品种</th>
            <th>定植地点</th>
            <th>定植时间</th>
            <th>糖度</th>
            <th>重量</th>
            <th>点赞数</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="product in products" :key="product.id">
            <td>{{ product.id }}</td>
            <td>{{ product.name }}</td>
            <td>{{ product.variety }}</td>
            <td>{{ product.location }}</td>
            <td>{{ formatDate(product.planting_time) }}</td>
            <td>{{ product.sugar_content || 0 }}°Bx</td>
            <td>{{ product.weight || 0 }}g</td>
            <td>{{ product.like_count || 0 }}</td>
            <td class="action-buttons">
              <button class="btn-edit" @click="openEditProductModal(product)">
                <i class="edit-icon">✏️</i>
                编辑
              </button>
              <button class="btn-delete" @click="confirmDeleteProduct(product.id)">
                <i class="delete-icon">🗑️</i>
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
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
    
    <!-- 添加/编辑产品模态框 -->
    <div v-if="showProductModal" class="modal-overlay" @click.self="closeProductModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑产品' : '添加产品' }}</h3>
          <button class="btn-close" @click="closeProductModal">×</button>
        </div>
        
        <div class="modal-body">
          <form @submit.prevent="handleProductSubmit">
            <div class="form-grid">
              <div class="form-group">
                <label for="name" class="form-label">产品名称</label>
                <input 
                  type="text" 
                  id="name" 
                  v-model="productForm.name" 
                  class="form-input" 
                  placeholder="请输入产品名称" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label for="variety" class="form-label">品种</label>
                <input 
                  type="text" 
                  id="variety" 
                  v-model="productForm.variety" 
                  class="form-input" 
                  placeholder="请输入品种" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label for="variety_code" class="form-label">品种编码</label>
                <input 
                  type="text" 
                  id="variety_code" 
                  v-model="productForm.variety_code" 
                  class="form-input" 
                  placeholder="请输入品种编码" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label for="location" class="form-label">定植地点</label>
                <input 
                  type="text" 
                  id="location" 
                  v-model="productForm.location" 
                  class="form-input" 
                  placeholder="请输入定植地点" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label for="planting_time" class="form-label">定植时间</label>
                <input 
                  type="date" 
                  id="planting_time" 
                  v-model="productForm.planting_time" 
                  class="form-input" 
                  required
                >
              </div>
              
              <div class="form-group">
                <label for="harvest_start_time" class="form-label">采收起始时间</label>
                <input 
                  type="date" 
                  id="harvest_start_time" 
                  v-model="productForm.harvest_start_time" 
                  class="form-input"
                >
              </div>
              
              <div class="form-group">
                <label for="harvest_end_time" class="form-label">采收终止时间</label>
                <input 
                  type="date" 
                  id="harvest_end_time" 
                  v-model="productForm.harvest_end_time" 
                  class="form-input"
                >
              </div>
              
              <div class="form-group">
                <label for="sugar_content" class="form-label">糖度 (°Bx)</label>
                <input 
                  type="number" 
                  id="sugar_content" 
                  v-model.number="productForm.sugar_content" 
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
                  v-model.number="productForm.weight" 
                  class="form-input" 
                  placeholder="请输入重量"
                  step="1"
                >
              </div>
              
              <div class="form-group full-width">
                <label for="taste_description" class="form-label">口感描述</label>
                <textarea 
                  id="taste_description" 
                  v-model="productForm.taste_description" 
                  class="form-textarea" 
                  placeholder="请输入口感描述"
                  rows="3"
                ></textarea>
              </div>
              
              <div class="form-group full-width">
                <label for="suitable_for" class="form-label">适应人群</label>
                <textarea 
                  id="suitable_for" 
                  v-model="productForm.suitable_for" 
                  class="form-textarea" 
                  placeholder="请输入适应人群"
                  rows="2"
                ></textarea>
              </div>
              
              <div class="form-group full-width">
                <label for="quality_summary" class="form-label">品质小结</label>
                <textarea 
                  id="quality_summary" 
                  v-model="productForm.quality_summary" 
                  class="form-textarea" 
                  placeholder="请输入品质小结"
                  rows="3"
                ></textarea>
              </div>
            </div>
            
            <div class="modal-actions">
              <button type="button" class="btn-cancel" @click="closeProductModal">取消</button>
              <button type="submit" class="btn-save" :disabled="saving">
                <span v-if="saving" class="spinner"></span>
                <span>{{ saving ? '保存中...' : '保存' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { productAPI, locationAPI } from '../services/api'
import { useToast } from '../composables/useToast'

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

// 模态框相关
const showProductModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const productForm = ref({
  name: '',
  variety: '',
  variety_code: '',
  location: '',
  planting_time: '',
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
      page: currentPage.value,
      limit: 10
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

const openAddProductModal = () => {
  isEditing.value = false
  productForm.value = {
    name: '',
    variety: '',
    variety_code: '',
    location: '',
    planting_time: '',
    harvest_start_time: '',
    harvest_end_time: '',
    sugar_content: 0,
    weight: 0,
    taste_description: '',
    suitable_for: '',
    quality_summary: ''
  }
  showProductModal.value = true
}

const openEditProductModal = (product) => {
  isEditing.value = true
  productForm.value = {
    id: product.id,
    name: product.name,
    variety: product.variety,
    variety_code: product.variety_code,
    location: product.location,
    planting_time: product.planting_time,
    harvest_start_time: product.harvest_start_time || '',
    harvest_end_time: product.harvest_end_time || '',
    sugar_content: product.sugar_content || 0,
    weight: product.weight || 0,
    taste_description: product.taste_description || '',
    suitable_for: product.suitable_for || '',
    quality_summary: product.quality_summary || ''
  }
  showProductModal.value = true
}

const closeProductModal = () => {
  showProductModal.value = false
}

const handleProductSubmit = async () => {
  saving.value = true
  
  try {
    if (isEditing.value) {
      // 编辑产品
      await productAPI.update(productForm.value.id, productForm.value)
      showToast('产品更新成功', 'success')
    } else {
      // 添加产品
      await productAPI.create(productForm.value)
      showToast('产品添加成功', 'success')
    }
    
    closeProductModal()
    fetchProducts()
  } catch (err) {
    showToast(err.message || '操作失败', 'error')
  } finally {
    saving.value = false
  }
}

const confirmDeleteProduct = (id) => {
  if (confirm('确定要删除这个产品吗？')) {
    deleteProduct(id)
  }
}

const deleteProduct = async (id) => {
  try {
    await productAPI.delete(id)
    showToast('产品删除成功', 'success')
    fetchProducts()
  } catch (err) {
    showToast(err.message || '删除失败', 'error')
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.admin-product-management {
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
}

.header-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.btn-add-product {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background-color: #10b981;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-add-product:hover {
  background-color: #059669;
}

.add-icon {
  font-size: 16px;
}

.search-filter {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 200px;
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.filter-select {
  padding: 10px 15px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  background-color: white;
  cursor: pointer;
  transition: border-color 0.3s;
}

.filter-select:focus {
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
  flex: 1;
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

.product-list {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.product-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
  overflow: auto;
}

.product-table th,
.product-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.product-table th {
  background-color: #f8fafc;
  font-weight: 600;
  color: #475569;
  position: sticky;
  top: 0;
  z-index: 10;
}

.product-table tr:hover {
  background-color: #f8fafc;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-edit {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-edit:hover {
  background-color: #2563eb;
}

.btn-delete {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-delete:hover {
  background-color: #dc2626;
}

.edit-icon, .delete-icon {
  font-size: 12px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
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

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  width: 100%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6b7280;
  transition: color 0.3s;
}

.btn-close:hover {
  color: #374151;
}

.modal-body {
  padding: 20px;
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
  min-height: 80px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

.btn-cancel {
  padding: 10px 20px;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-cancel:hover {
  background-color: #f3f4f6;
  border-color: #d1d5db;
}

.btn-save {
  padding: 10px 20px;
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
  .admin-product-management {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .search-filter {
    flex-direction: column;
  }
  
  .product-table {
    font-size: 12px;
  }
  
  .product-table th,
  .product-table td {
    padding: 8px 12px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .btn-edit,
  .btn-delete {
    justify-content: center;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .modal-content {
    max-height: 95vh;
  }
}
</style>