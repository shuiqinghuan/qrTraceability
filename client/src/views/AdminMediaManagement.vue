<template>
  <div class="admin-media-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="header-title">媒体管理</h2>
      <div class="header-actions">
        <select v-model="selectedProduct" class="product-select" @change="fetchMedia">
          <option value="">选择产品</option>
          <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
        </select>
        <button class="btn-upload" @click="openUploadModal">
          <i class="upload-icon">📤</i>
          <span>上传媒体</span>
        </button>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载媒体列表中...</p>
    </div>
    
    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button class="btn-retry" @click="fetchMedia">重试</button>
    </div>
    
    <!-- 媒体列表 -->
    <div v-else class="media-list">
      <div v-if="media.length === 0" class="no-media">
        <div class="no-media-icon">🖼️</div>
        <h3>暂无媒体</h3>
        <p v-if="selectedProduct">该产品暂无媒体文件</p>
        <p v-else>请选择产品查看媒体文件</p>
      </div>
      
      <div v-else class="media-grid">
        <div v-for="item in media" :key="item.id" class="media-item">
          <div class="media-preview">
            <img v-if="item.type === 'image'" :src="item.url" :alt="item.name" @error="handleImageError">
            <div v-else-if="item.type === 'video'" class="video-preview">
              <video :src="item.url" controls></video>
            </div>
            <div class="media-overlay">
              <button class="btn-delete-media" @click="confirmDeleteMedia(item.id)">
                <i class="delete-icon">🗑️</i>
              </button>
            </div>
          </div>
          <div class="media-info">
            <div class="media-name">{{ item.name }}</div>
            <div class="media-meta">
              <span class="media-type">{{ item.type === 'image' ? '图片' : '视频' }}</span>
              <span class="media-size">{{ formatFileSize(item.size) }}</span>
              <span class="media-date">{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 上传媒体模态框 -->
    <div v-if="showUploadModal" class="modal-overlay" @click.self="closeUploadModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>上传媒体</h3>
          <button class="btn-close" @click="closeUploadModal">×</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label for="upload-product" class="form-label">选择产品</label>
            <select v-model="uploadForm.product_id" id="upload-product" class="form-input" required>
              <option value="">选择产品</option>
              <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="form-label">选择文件</label>
            <div class="file-upload-area">
              <input 
                type="file" 
                ref="fileInput" 
                multiple 
                accept="image/*,video/*" 
                class="file-input" 
                @change="handleFileChange"
              >
              <div class="upload-placeholder">
                <div class="upload-icon">📁</div>
                <p>点击或拖拽文件到此处上传</p>
                <p class="upload-hint">支持图片和视频文件</p>
              </div>
            </div>
          </div>
          
          <div v-if="selectedFiles.length > 0" class="selected-files">
            <h4>已选择的文件</h4>
            <ul class="file-list">
              <li v-for="(file, index) in selectedFiles" :key="index" class="file-item">
                <span class="file-name">{{ file.name }}</span>
                <span class="file-size">{{ formatFileSize(file.size) }}</span>
                <button class="btn-remove-file" @click="removeFile(index)">×</button>
              </li>
            </ul>
          </div>
          
          <div class="modal-actions">
            <button type="button" class="btn-cancel" @click="closeUploadModal">取消</button>
            <button type="button" class="btn-upload-files" @click="uploadFiles" :disabled="uploading || selectedFiles.length === 0">
              <span v-if="uploading" class="spinner"></span>
              <span>{{ uploading ? '上传中...' : '上传文件' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { productAPI, mediaAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const { showToast } = useToast()

// 响应式数据
const products = ref([])
const media = ref([])
const loading = ref(true)
const error = ref(null)
const selectedProduct = ref('')

// 上传相关
const showUploadModal = ref(false)
const uploading = ref(false)
const selectedFiles = ref([])
const fileInput = ref(null)
const uploadForm = ref({
  product_id: ''
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

const fetchMedia = async () => {
  if (!selectedProduct.value) {
    media.value = []
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    const response = await mediaAPI.list(selectedProduct.value)
    media.value = response.media || []
  } catch (err) {
    error.value = err.message || '获取媒体列表失败'
    console.error('获取媒体列表失败', err)
  } finally {
    loading.value = false
  }
}

const openUploadModal = () => {
  selectedFiles.value = []
  uploadForm.value = {
    product_id: selectedProduct.value || ''
  }
  showUploadModal.value = true
}

const closeUploadModal = () => {
  showUploadModal.value = false
  selectedFiles.value = []
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const handleFileChange = (event) => {
  const files = Array.from(event.target.files)
  selectedFiles.value = [...selectedFiles.value, ...files]
}

const removeFile = (index) => {
  selectedFiles.value.splice(index, 1)
}

const uploadFiles = async () => {
  if (!uploadForm.product_id || selectedFiles.value.length === 0) {
    showToast('请选择产品和文件', 'warning')
    return
  }
  
  uploading.value = true
  
  try {
    const formData = new FormData()
    formData.append('product_id', uploadForm.product_id)
    
    selectedFiles.value.forEach(file => {
      formData.append('files', file)
    })
    
    await mediaAPI.upload(formData)
    showToast('媒体上传成功', 'success')
    
    closeUploadModal()
    if (selectedProduct.value === uploadForm.product_id) {
      fetchMedia()
    }
  } catch (err) {
    showToast(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
  }
}

const confirmDeleteMedia = (id) => {
  if (confirm('确定要删除这个媒体文件吗？')) {
    deleteMedia(id)
  }
}

const deleteMedia = async (id) => {
  try {
    await mediaAPI.delete(id)
    showToast('媒体删除成功', 'success')
    fetchMedia()
  } catch (err) {
    showToast(err.message || '删除失败', 'error')
  }
}

const handleImageError = (event) => {
  event.target.src = 'https://via.placeholder.com/200x200?text=图片加载失败'
}

const formatFileSize = (size) => {
  if (!size) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return `${size.toFixed(2)} ${units[i]}`
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.admin-media-management {
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

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
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

.btn-upload {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-upload:hover {
  background-color: #2563eb;
}

.upload-icon {
  font-size: 16px;
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

.no-media {
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

.no-media-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.no-media h3 {
  color: #4b5563;
  margin-bottom: 10px;
}

.no-media p {
  color: #6b7280;
}

.media-list {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin: 20px 0;
  overflow-y: auto;
  flex: 1;
}

.media-item {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  background-color: white;
}

.media-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.media-preview {
  position: relative;
  height: 150px;
  overflow: hidden;
  background-color: #f9fafb;
}

.media-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.media-item:hover .media-preview img {
  transform: scale(1.05);
}

.video-preview {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-preview video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.media-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.media-item:hover .media-overlay {
  opacity: 1;
}

.btn-delete-media {
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-delete-media:hover {
  background-color: #dc2626;
}

.delete-icon {
  font-size: 16px;
}

.media-info {
  padding: 12px;
}

.media-name {
  font-size: 14px;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.media-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.media-type {
  padding: 2px 6px;
  background-color: #e0f2fe;
  color: #0284c7;
  border-radius: 4px;
  align-self: flex-start;
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
  max-width: 600px;
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

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.file-upload-area {
  border: 2px dashed #e2e8f0;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.file-upload-area:hover {
  border-color: #3b82f6;
  background-color: #f8fafc;
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.upload-placeholder .upload-icon {
  font-size: 48px;
  color: #94a3b8;
}

.upload-placeholder p {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.upload-hint {
  font-size: 12px !important;
  color: #94a3b8 !important;
}

.selected-files {
  margin-top: 20px;
  padding: 16px;
  background-color: #f8fafc;
  border-radius: 6px;
}

.selected-files h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #4b5563;
}

.file-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #e2e8f0;
}

.file-item:last-child {
  border-bottom: none;
}

.file-name {
  flex: 1;
  font-size: 14px;
  color: #4b5563;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-size {
  font-size: 12px;
  color: #6b7280;
  margin: 0 12px;
  white-space: nowrap;
}

.btn-remove-file {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #6b7280;
  transition: color 0.3s;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-remove-file:hover {
  color: #ef4444;
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

.btn-upload-files {
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

.btn-upload-files:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn-upload-files:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-upload-files .spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  margin-bottom: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-media-management {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .header-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .product-select {
    width: 100%;
  }
  
  .media-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
  
  .media-preview {
    height: 120px;
  }
  
  .modal-content {
    max-height: 95vh;
  }
}
</style>