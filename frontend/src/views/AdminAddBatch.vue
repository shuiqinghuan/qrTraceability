<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { config } from '../config.js'

const router = useRouter()
const products = ref([])
const formData = ref({
  product_id: '',
  planting_location: '',
  planting_date: ''
})
const loading = ref(false)
const error = ref(null)

// 获取产品列表
const fetchProducts = async () => {
  try {
    const response = await axios.get(`${config.apiBaseUrl}/api/products`)
    products.value = response.data
  } catch (err) {
    console.error('获取产品列表失败:', err)
    error.value = '获取产品列表失败'
  }
}

// 提交表单
const submitForm = async () => {
  // 验证表单
  if (!formData.value.product_id || !formData.value.planting_location || !formData.value.planting_date) {
    alert('请填写所有必填字段')
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    await axios.post(`${config.apiBaseUrl}/api/batches`, formData.value)
    alert('批次添加成功')
    router.push('/admin/batches')
  } catch (err) {
    console.error('添加批次失败:', err)
    error.value = '添加批次失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 取消操作
const cancel = () => {
  router.push('/admin/batches')
}

// 页面加载时获取产品列表
onMounted(() => {
  fetchProducts()
})
</script>

<template>
  <div class="admin-add-batch">
    <!-- 页面头部 -->
    <div class="page-header">
      <h3>添加批次</h3>
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
          label="产品"
          required
        >
          <template #input>
            <van-picker
              v-model="formData.product_id"
              :columns="products.map(p => ({ text: p.name, value: p.id }))"
              placeholder="请选择产品"
              show-toolbar
              toolbar-position="bottom"
            />
          </template>
        </van-field>
        <van-field
          v-model="formData.planting_location"
          label="定植地点"
          placeholder="请输入定植地点"
          required
        />
        <van-field
          label="定植时间"
          required
        >
          <template #input>
            <van-datetime-picker
              v-model="formData.planting_date"
              type="date"
              placeholder="请选择定植时间"
              show-toolbar
              toolbar-position="bottom"
            />
          </template>
        </van-field>
      </van-form>
    </div>
  </div>
</template>

<style scoped>
.admin-add-batch {
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
