<template>
  <div class="register-container">
    <div class="register-form">
      <h2>农产品二维码追溯系统</h2>
      <h3>注册</h3>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="phone_number">手机号码</label>
          <input 
            type="tel" 
            id="phone_number" 
            v-model="form.phone_number" 
            required
            placeholder="请输入手机号码"
          >
        </div>
        <div class="form-group">
          <label for="password">密码</label>
          <input 
            type="password" 
            id="password" 
            v-model="form.password" 
            required
            placeholder="请输入密码"
          >
        </div>
        <div class="form-group">
          <label for="role">角色</label>
          <select id="role" v-model="form.role" required>
            <option value="">请选择角色</option>
            <option value="serverseed">种子管理</option>
            <option value="servergrow">生长管理</option>
            <option value="servermanager">品质管理</option>
            <option value="clentcustomer">客户</option>
          </select>
        </div>
        <button type="submit" class="register-button">注册</button>
        <div class="login-link">
          已有账号？ <router-link to="/login">立即登录</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authAPI } from '../services/api'

const router = useRouter()
const form = ref({
  phone_number: '',
  password: '',
  role: ''
})

const handleRegister = async () => {
  try {
    const response = await authAPI.register(form.value)
    if (response.success) {
      alert('注册成功，请登录')
      router.push('/login')
    }
  } catch (error) {
    alert('注册失败，请检查输入信息')
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.register-form {
  background: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

h2 {
  text-align: center;
  color: #333;
  margin-bottom: 20px;
}

h3 {
  text-align: center;
  color: #666;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 5px;
  color: #666;
}

input, select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.register-button {
  width: 100%;
  padding: 12px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 20px;
}

.register-button:hover {
  background-color: #45a049;
}

.login-link {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.login-link a {
  color: #4CAF50;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>