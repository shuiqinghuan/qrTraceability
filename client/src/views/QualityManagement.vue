<template>
  <div class="quality-management">
    <h2>品质管理</h2>
    <div class="action-bar">
      <button class="add-button" @click="showAddForm = true">添加品质指标</button>
    </div>
    
    <div v-if="showAddForm" class="form-container">
      <h3>添加品质指标</h3>
      <form @submit.prevent="handleAddQuality">
        <div class="form-group">
          <label for="planting_id">种植ID</label>
          <input type="number" id="planting_id" v-model="newQuality.planting_id" required>
        </div>
        <div class="form-group">
          <label for="sugar_content">糖度</label>
          <input type="number" step="0.1" id="sugar_content" v-model="newQuality.sugar_content" required>
        </div>
        <div class="form-group">
          <label for="weight">重量</label>
          <input type="number" step="0.1" id="weight" v-model="newQuality.weight" required>
        </div>
        <div class="form-group">
          <label for="taste">口感分析</label>
          <textarea id="taste" v-model="newQuality.taste" required></textarea>
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-button">保存</button>
          <button type="button" class="cancel-button" @click="showAddForm = false">取消</button>
        </div>
      </form>
    </div>
    
    <div class="quality-list">
      <h3>品质指标列表</h3>
      <table class="quality-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>种植ID</th>
            <th>糖度</th>
            <th>重量</th>
            <th>口感分析</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="quality in qualities" :key="quality.id">
            <td>{{ quality.id }}</td>
            <td>{{ quality.planting_id }}</td>
            <td>{{ quality.sugar_content }}</td>
            <td>{{ quality.weight }}</td>
            <td>{{ quality.taste }}</td>
            <td class="actions">
              <button class="edit-button" @click="editQuality(quality)">编辑</button>
              <button class="delete-button" @click="deleteQuality(quality.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { qualityAPI } from '../services/api'

const qualities = ref([])
const showAddForm = ref(false)
const newQuality = ref({
  planting_id: '',
  sugar_content: '',
  weight: '',
  taste: ''
})

onMounted(() => {
  fetchQualities()
})

const fetchQualities = async () => {
  try {
    const data = await qualityAPI.list()
    qualities.value = data
  } catch (error) {
    console.error('获取品质指标列表失败', error)
  }
}

const handleAddQuality = async () => {
  try {
    await qualityAPI.create(newQuality.value)
    showAddForm.value = false
    newQuality.value = {
      planting_id: '',
      sugar_content: '',
      weight: '',
      taste: ''
    }
    fetchQualities()
  } catch (error) {
    console.error('添加品质指标失败', error)
  }
}

const editQuality = (quality) => {
  // 编辑功能实现
  console.log('编辑品质指标', quality)
}

const deleteQuality = async (id) => {
  if (confirm('确定要删除这个品质指标吗？')) {
    try {
      await qualityAPI.delete(id)
      fetchQualities()
    } catch (error) {
      console.error('删除品质指标失败', error)
    }
  }
}
</script>

<style scoped>
.quality-management {
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

input, textarea {
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

.quality-list {
  margin-top: 30px;
}

.quality-list h3 {
  margin-bottom: 15px;
  color: #333;
}

.quality-table {
  width: 100%;
  border-collapse: collapse;
}

.quality-table th, .quality-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.quality-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.quality-table tr:hover {
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