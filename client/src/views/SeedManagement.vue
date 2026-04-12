<template>
  <div class="seed-management">
    <h2>种子管理</h2>
    <div class="action-bar">
      <button class="add-button" @click="showAddForm = true">添加种子</button>
    </div>
    
    <div v-if="showAddForm" class="form-container">
      <h3>添加种子信息</h3>
      <form @submit.prevent="handleAddSeed">
        <div class="form-group">
          <label for="name">种子名称</label>
          <input type="text" id="name" v-model="newSeed.name" required>
        </div>
        <div class="form-group">
          <label for="variety">品种</label>
          <input type="text" id="variety" v-model="newSeed.variety" required>
        </div>
        <div class="form-group">
          <label for="description">描述</label>
          <textarea id="description" v-model="newSeed.description" required></textarea>
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-button">保存</button>
          <button type="button" class="cancel-button" @click="showAddForm = false">取消</button>
        </div>
      </form>
    </div>
    
    <div class="seed-list">
      <h3>种子列表</h3>
      <table class="seed-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>名称</th>
            <th>品种</th>
            <th>描述</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="seed in seeds" :key="seed.id">
            <td>{{ seed.id }}</td>
            <td>{{ seed.name }}</td>
            <td>{{ seed.variety }}</td>
            <td>{{ seed.description }}</td>
            <td class="actions">
              <button class="edit-button" @click="editSeed(seed)">编辑</button>
              <button class="delete-button" @click="deleteSeed(seed.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { seedAPI } from '../services/api'

const seeds = ref([])
const showAddForm = ref(false)
const newSeed = ref({
  name: '',
  variety: '',
  description: ''
})

onMounted(() => {
  fetchSeeds()
})

const fetchSeeds = async () => {
  try {
    const data = await seedAPI.list()
    seeds.value = data
  } catch (error) {
    console.error('获取种子列表失败', error)
  }
}

const handleAddSeed = async () => {
  try {
    await seedAPI.create(newSeed.value)
    showAddForm.value = false
    newSeed.value = {
      name: '',
      variety: '',
      description: ''
    }
    fetchSeeds()
  } catch (error) {
    console.error('添加种子失败', error)
  }
}

const editSeed = (seed) => {
  // 编辑功能实现
  console.log('编辑种子', seed)
}

const deleteSeed = async (id) => {
  if (confirm('确定要删除这个种子吗？')) {
    try {
      await seedAPI.delete(id)
      fetchSeeds()
    } catch (error) {
      console.error('删除种子失败', error)
    }
  }
}
</script>

<style scoped>
.seed-management {
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

.seed-list {
  margin-top: 30px;
}

.seed-list h3 {
  margin-bottom: 15px;
  color: #333;
}

.seed-table {
  width: 100%;
  border-collapse: collapse;
}

.seed-table th, .seed-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.seed-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.seed-table tr:hover {
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