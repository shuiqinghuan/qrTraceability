<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProductStore } from '@/stores/product'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const productStore = useProductStore()

const searchKeyword = ref('')

onMounted(() => {
  productStore.fetchProducts()
})

const handleSearch = () => {
  productStore.fetchProducts()
}

const handleAdd = () => {
  router.push('/admin/product/new')
}

const handleEdit = (id: number) => {
  router.push(`/admin/product/${id}/edit`)
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该产品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const success = await productStore.deleteProduct(id)
    if (success) {
      ElMessage.success('删除成功')
    }
  } catch {
    // cancelled
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '--'
  return dateStr.split('T')[0]
}
</script>

<template>
  <div class="product-list">
    <div class="page-header">
      <h2 class="page-title">产品管理</h2>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增产品
      </el-button>
    </div>
    
    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索产品名称或编码"
        style="width: 300px;"
        clearable
        @clear="handleSearch"
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
    </div>
    
    <el-table
      :data="productStore.products"
      v-loading="productStore.loading"
      stripe
      class="product-table"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="产品名称" min-width="150" />
      <el-table-column prop="code" label="品种编码" width="120" />
      <el-table-column prop="planting_location" label="定植地点" min-width="150" />
      <el-table-column prop="planting_date" label="定植时间" width="120">
        <template #default="{ row }">
          {{ formatDate(row.planting_date) }}
        </template>
      </el-table-column>
      <el-table-column prop="quality" label="品质等级" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.quality" type="success">{{ row.quality }}</el-tag>
          <span v-else>--</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row.id)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script lang="ts">
import { Plus, Search } from '@element-plus/icons-vue'
export default {
  components: { Plus, Search }
}
</script>

<style lang="scss" scoped>
.product-list {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.product-table {
  :deep(.el-table__header th) {
    background: #f5f7fa;
  }
}
</style>
