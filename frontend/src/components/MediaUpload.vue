<template>
  <div class="media-upload">
    <h3>
      <span class="icon">🖼️</span>
      图片上传
    </h3>

    <div class="upload-section">
      <div class="form-group">
        <label>产品ID:</label>
        <input v-model.number="productId" type="number" placeholder="输入产品ID" />
      </div>

      <div class="form-group">
        <label>图片标题:</label>
        <input v-model="title" type="text" placeholder="输入图片标题（可选）" />
      </div>

      <div class="form-group file-upload">
        <label>选择图片:</label>
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          @change="handleFileChange"
        />
        <div v-if="selectedFile" class="file-info">
          <span class="file-icon">📷</span>
          <span class="file-name">{{ selectedFile.name }}</span>
          <span class="file-size">({{ formatFileSize(selectedFile.size) }})</span>
        </div>
      </div>

      <button
        @click="handleUpload"
        :disabled="!selectedFile || uploading"
        class="upload-btn"
      >
        <span v-if="uploading" class="spinner">⏳</span>
        {{ uploading ? '上传中...' : '上传图片' }}
      </button>
    </div>

    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>

    <div v-if="uploadedMedia.length > 0" class="uploaded-list">
      <h4>✅ 已上传的图片 ({{ uploadedMedia.length }})</h4>
      <div v-for="item in uploadedMedia" :key="item.id" class="media-item">
        <img :src="item.url" :alt="item.title" class="uploaded-image" />
        <span class="media-title">{{ item.title || '产品图片' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 媒体文件上传组件
 * 提供图片文件选择、上传进度反馈及已上传图片列表展示功能
 */

import { ref } from 'vue'
import api from '@/api'

// 表单状态：产品ID、标题、选中文件、上传状态、消息提示
const productId = ref(1)
const title = ref('')
const selectedFile = ref(null)
const uploading = ref(false)
const message = ref('')
const messageType = ref('success')
const uploadedMedia = ref([])
const fileInput = ref(null)

/** 文件选择变更处理 */
const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
  }
}

/** 格式化文件大小为可读字符串（B/KB/MB） */
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

/** 执行图片上传，提交FormData至后端API */
const handleUpload = async () => {
  if (!selectedFile.value) return

  uploading.value = true
  message.value = ''

  const formData = new FormData()
  formData.append('mediaType', 'image')
  formData.append('title', title.value)
  formData.append('file', selectedFile.value)

  try {
    const response = await api.post(`/products/${productId.value}/media/`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    message.value = '上传成功！'
    messageType.value = 'success'
    uploadedMedia.value.unshift(response.data)

    // 上传成功后重置表单状态
    selectedFile.value = null
    if (fileInput.value) {
      fileInput.value.value = ''
    }
    title.value = ''
  } catch (error) {
    message.value = '上传失败: ' + (error.response?.data?.message || error.message)
    messageType.value = 'error'
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
/* 组件根容器 */
.media-upload {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

h3 {
  color: var(--primary-color);
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

h3 .icon {
  font-size: 22px;
}

/* 上传表单区域 */
.upload-section {
  background: #f5f5f5;
  padding: 24px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 18px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #444;
  font-size: 14px;
}

.form-group input {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.file-upload input[type="file"] {
  padding: 8px 0;
  border: none;
  cursor: pointer;
}

/* 已选文件信息提示 */
.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
  padding: 10px 14px;
  background: #e8f5e9;
  border-radius: 8px;
  font-size: 13px;
  color: #2e7d32;
}

.file-icon {
  font-size: 16px;
}

.file-name {
  font-weight: 500;
}

.file-size {
  color: #666;
  margin-left: auto;
}

/* 上传按钮样式 */
.upload-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.upload-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(76, 175, 80, 0.4);
}

.upload-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 加载旋转动画 */
.spinner {
  display: inline-block;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 上传结果消息提示（成功/失败） */
.message {
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 15px;
  font-size: 14px;
}

.message.success {
  background: #e8f5e9;
  color: #2e7d32;
  border: 1px solid #a5d6a7;
}

.message.error {
  background: #ffebee;
  color: #c62828;
  border: 1px solid #ef9a9a;
}

/* 已上传图片列表区域 */
.uploaded-list {
  margin-top: 20px;
}

.uploaded-list h4 {
  margin-bottom: 15px;
  color: #333;
  font-size: 14px;
}

.media-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  margin-bottom: 10px;
  border: 1px solid #eee;
}

.uploaded-image {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 6px;
}

.media-title {
  flex: 1;
  color: #333;
  font-size: 14px;
}
</style>
