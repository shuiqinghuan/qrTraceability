<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductStore } from '@/stores/product'
import { mediaApi } from '@/api'
import { ElMessage } from 'element-plus'
import type { Product } from '@/types'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()

const productId = computed(() => route.params.id ? Number(route.params.id) : null)
const isEdit = computed(() => !!productId.value)

const formRef = ref()
const loading = ref(false)
const uploadingImages = ref(false)

const form = ref<Partial<Product>>({
  name: '',
  code: '',
  planting_location: '',
  planting_date: '',
  images: [],
  video: '',
  harvest_start_date: '',
  harvest_end_date: '',
  sugar_content: null,
  weight: null,
  taste: '',
  quality: '',
  quality_summary: '',
  suitable_for: []
})

const suitableForInput = ref('')

const rules = {
  name: [{ required: true, message: '请输入品种名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入品种编码', trigger: 'blur' }],
  planting_location: [{ required: true, message: '请输入定植地点', trigger: 'blur' }],
  planting_date: [{ required: true, message: '请选择定植时间', trigger: 'change' }]
}

onMounted(async () => {
  if (isEdit.value) {
    loading.value = true
    await productStore.fetchProduct(productId.value!)
    if (productStore.currentProduct) {
      form.value = { ...productStore.currentProduct }
      suitableForInput.value = productStore.currentProduct.suitable_for?.join(', ') || ''
    }
    loading.value = false
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      
      if (suitableForInput.value) {
        form.value.suitable_for = suitableForInput.value.split(',').map(s => s.trim()).filter(Boolean)
      }
      
      let success
      if (isEdit.value) {
        success = await productStore.updateProduct(productId.value!, form.value)
      } else {
        success = await productStore.createProduct(form.value)
      }
      
      loading.value = false
      
      if (success) {
        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
        router.push('/admin')
      }
    }
  })
}

const handleCancel = () => {
  router.back()
}
</script>

<template>
  <div class="product-edit">
    <div class="page-header">
      <h2 class="page-title">{{ isEdit ? '编辑产品' : '新增产品' }}</h2>
    </div>
    
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      v-loading="loading"
      class="product-form"
    >
      <el-card header="基本信息" class="form-section">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="品种名称" prop="name">
              <el-input v-model="form.name" placeholder="如：枣甜5号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="品种编码" prop="code">
              <el-input v-model="form.code" placeholder="如：4395" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="定植地点" prop="planting_location">
              <el-input v-model="form.planting_location" placeholder="如：新疆和田洛浦县红枣基地" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="定植时间" prop="planting_date">
              <el-date-picker
                v-model="form.planting_date"
                type="date"
                placeholder="选择日期"
                style="width: 100%;"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>
      
      <el-card header="媒体信息" class="form-section">
        <el-form-item label="产品图片">
          <el-input
            v-model="form.images"
            type="textarea"
            :rows="2"
            placeholder="每行一个图片URL"
          />
          <div class="form-tip">支持多张图片，每行一个URL</div>
        </el-form-item>
        
        <el-form-item label="视频链接">
          <el-input v-model="form.video" placeholder="视频URL地址" />
        </el-form-item>
      </el-card>
      
      <el-card header="采收质量" class="form-section">
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="采收起始时间">
              <el-date-picker
                v-model="form.harvest_start_date"
                type="date"
                placeholder="选择日期"
                style="width: 100%;"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="采收终止时间">
              <el-date-picker
                v-model="form.harvest_end_date"
                type="date"
                placeholder="选择日期"
                style="width: 100%;"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="糖度">
              <el-input-number
                v-model="form.sugar_content"
                :precision="1"
                :step="0.1"
                :min="0"
                style="width: 100%;"
              >
                <template #suffix>Brix</template>
              </el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="单果重量">
              <el-input-number
                v-model="form.weight"
                :precision="1"
                :step="0.1"
                :min="0"
                style="width: 100%;"
              >
                <template #suffix>克</template>
              </el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="口感描述">
              <el-input v-model="form.taste" placeholder="如：肉质紧密、甘甜爽口" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="品质等级">
              <el-input v-model="form.quality" placeholder="如：特级、一级" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>
      
      <el-card header="品质小结" class="form-section">
        <el-form-item label="品质小结">
          <el-input
            v-model="form.quality_summary"
            type="textarea"
            :rows="4"
            placeholder="请输入产品品质综合评价"
          />
        </el-form-item>
        
        <el-form-item label="适应人群">
          <el-input
            v-model="suitableForInput"
            placeholder="用逗号分隔，如：一般人群,儿童,老年人"
          />
          <div class="form-tip">多个标签用逗号分隔</div>
        </el-form-item>
      </el-card>
      
      <div class="form-actions">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '保存修改' : '创建产品' }}
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<style lang="scss" scoped>
.product-edit {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.form-section {
  margin-bottom: 20px;
  
  :deep(.el-card__header) {
    font-weight: 600;
    color: #333;
    background: #fafafa;
  }
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}
</style>
