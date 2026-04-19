<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { config } from '../config.js'

const router = useRouter()
const formData = ref({
  name: '',
  code: '',
  description: ''
})
const loading = ref(false)
const error = ref(null)

// 提交表单
const submitForm = async () => {
  // 验证表单
  if (!formData.value.name || !formData.value.code) {
    alert('请填写产品名称和代码')
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    await axios.post(`${config.apiBaseUrl}/api/products`, formData.value)
    alert('产品添加成功')
    router.push('/admin/products')
  } catch (err) {
    console.error('添加产品失败:', err)
    error.value = '添加产品失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 取消操作
const cancel = () => {
  router.push('/admin/products')
}
</script>

<template>
  <div class="admin-add-product">
    <!-- 页面头部 -->
    <div class="page-header">
      <h3>添加产品</h3>
      <div class="header-actions">
        <van-button @click="cancel">取消</van-button>
        <van-button type="primary" :loading="loading" @click="submitForm">保存</van-button>
      </div>
    </div>
    
    <!-- 错误信息 -->
    <van-toast v-if="error" type="error" :message="error" duration="3000" />
    
    <!-- 表单 -->
    <div class="form-container">
      <van-form>
        <van-field
          v-model="formData.name"
          label="产品名称"
          placeholder="请输入产品名称"
          required
        />
        <van-field
          v-model="formData.code"
          label="产品代码"
          placeholder="请输入产品代码"
          required
        />
        <van-field
          v-model="formData.description"
          label="产品描述"
          type="textarea"
          placeholder="请输入产品描述"
          :rows="3"
        />
      </van-form>
    </div>
  </div>
</template>

<style scoped>
.admin-add-product {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e5e5e5;
}

.page-header h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.form-container {
  background-color: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.van-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .header-actions {
    align-self: flex-end;
  }
  
  .form-container {
    padding: 20px;
  }
}
</style>
