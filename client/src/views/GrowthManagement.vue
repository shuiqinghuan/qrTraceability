<template>
  <div class="growth-management">
    <h2>生长管理</h2>
    <div class="action-bar">
      <button class="add-button" @click="showAddForm = true">添加生长媒体</button>
    </div>
    
    <div v-if="showAddForm" class="form-container">
      <h3>添加生长媒体</h3>
      <form @submit.prevent="handleAddGrowth">
        <div class="form-group">
          <label for="planting_id">种植ID</label>
          <input type="number" id="planting_id" v-model="newGrowth.planting_id" required>
        </div>
        <div class="form-group">
          <label for="type">类型</label>
          <select id="type" v-model="newGrowth.type" required>
            <option value="image">图片</option>
            <option value="video">视频</option>
          </select>
        </div>
        <div class="form-group">
          <label for="url">媒体URL</label>
          <input type="text" id="url" v-model="newGrowth.url" required>
        </div>
        <div class="form-group">
          <label for="description">描述</label>
          <textarea id="description" v-model="newGrowth.description"></textarea>
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-button">保存</button>
          <button type="button" class="cancel-button" @click="showAddForm = false">取消</button>
        </div>
      </form>
    </div>
    
    <div class="growth-list">
      <h3>生长媒体列表</h3>
      <table class="growth-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>种植ID</th>
            <th>类型</th>
            <th>媒体URL</th>
            <th>描述</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="growth in growths" :key="growth.id">
            <td>{{ growth.id }}</td>
            <td>{{ growth.planting_id }}</td>
            <td>{{ growth.type === 'image' ? '图片' : '视频' }}</td>
            <td>{{ growth.url }}</td>
            <td>{{ growth.description }}</td>
            <td class="actions">
              <button class="edit-button" @click="editGrowth(growth)">编辑</button>
              <button class="delete-button" @click="deleteGrowth(growth.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { growthAPI } from '../services/api'

const growths = ref([])
const showAddForm = ref(false)
const newGrowth = ref({
  planting_id: '',
  type: 'image',
  url: '',
  description: ''
})

onMounted(() => {
  fetchGrowths()
})

const fetchGrowths = async () => {
  try {
    const data = await growthAPI.list()
    growths.value = data
  } catch (error) {
    console.error('获取生长媒体列表失败', error)
  }
}

const handleAddGrowth = async () => {
  try {
    await growthAPI.create(newGrowth.value)
    showAddForm.value = false
    newGrowth.value = {
      planting_id: '',
      type: 'image',
      url: '',
      description: ''
    }
    fetchGrowths()
  } catch (error) {
    console.error('添加生长媒体失败', error)
  }
}

const editGrowth = (growth) => {
  // 编辑功能实现
  console.log('编辑生长媒体', growth)
}

const deleteGrowth = async (id) => {
  if (confirm('确定要删除这个生长媒体吗？')) {
    try {
      await growthAPI.delete(id)
      fetchGrowths()
    } catch (error) {
      console.error('删除生长媒体失败', error)
    }
  }
}
</script>

<style scoped>
.growth-management {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-bottom: 20px;
  color: #333;
}

.action-bar {
  margin-bottom: 20px;
}

.add-button {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.add-button:hover {
  background-color: #45a049;
}

.form-container {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.form-container h3 {
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #666;
}

input, select, textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

textarea {
  resize: vertical;
  min-height: 100px;
}

.form-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}

.submit-button {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.submit-button:hover {
  background-color: #45a049;
}

.cancel-button {
  background-color: #999;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-button:hover {
  background-color: #777;
}

.growth-list {
  margin-top: 30px;
}

.growth-list h3 {
  margin-bottom: 15px;
  color: #333;
}

.growth-table {
  width: 100%;
  border-collapse: collapse;
}

.growth-table th, .growth-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.growth-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.growth-table tr:hover {
  background-color: #f5f5f5;
}

.actions {
  display: flex;
  gap: 10px;
}

.edit-button {
  background-color: #2196F3;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.edit-button:hover {
  background-color: #0b7dda;
}

.delete-button {
  background-color: #f44336;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.delete-button:hover {
  background-color: #da190b;
}
</style>