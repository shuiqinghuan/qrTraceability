<template>
  <div class="admin-login">
    <div class="login-container">
      <div class="login-header">
        <h1 class="login-title">后台管理系统</h1>
        <p class="login-subtitle">请登录以访问管理功能</p>
      </div>
      
      <div class="login-form">
        <!-- 错误提示 -->
        <div v-if="error" class="error-message">
          <span class="error-icon">⚠️</span>
          <span>{{ error }}</span>
        </div>
        
        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="username" class="form-label">用户名</label>
            <input 
              type="text" 
              id="username" 
              v-model="form.username" 
              class="form-input" 
              placeholder="请输入用户名" 
              required
              autofocus
            >
          </div>
          
          <div class="form-group">
            <label for="password" class="form-label">密码</label>
            <input 
              type="password" 
              id="password" 
              v-model="form.password" 
              class="form-input" 
              placeholder="请输入密码" 
              required
            >
          </div>
          
          <div class="form-actions">
            <button 
              type="submit" 
              class="btn-login" 
              :disabled="loading"
            >
              <span v-if="loading" class="spinner"></span>
              <span>{{ loading ? '登录中...' : '登录' }}</span>
            </button>
          </div>
        </form>
        
        <div class="login-info">
          <p>默认账号：lhseed</p>
          <p>默认密码：123456</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const router = useRouter()
const { showToast } = useToast()

// 响应式数据
const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)
const error = ref(null)

// 方法
const handleLogin = async () => {
  error.value = null
  loading.value = true
  
  try {
    const response = await adminAPI.login(form.value)
    
    // 存储token
    localStorage.setItem('admin_token', response.token)
    localStorage.setItem('admin_user', JSON.stringify(response.user))
    
    showToast('登录成功', 'success')
    
    // 跳转到后台管理仪表板
    router.push('/admin')
  } catch (err) {
    error.value = err.message || '登录失败，请检查用户名和密码'
    showToast(error.value, 'error')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 450px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 10px;
}

.login-subtitle {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.error-message {
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  padding: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #dc2626;
  font-size: 14px;
}

.error-icon {
  font-size: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
}

.form-input {
  padding: 12px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 16px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-actions {
  margin-top: 10px;
}

.btn-login {
  width: 100%;
  padding: 14px 20px;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
}

.btn-login:hover:not(:disabled) {
  background-color: #2563eb;
  transform: translateY(-2px);
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.login-info {
  margin-top: 20px;
  padding: 15px;
  background-color: #f9fafb;
  border-radius: 6px;
  font-size: 14px;
  color: #6b7280;
  line-height: 1.5;
}

.login-info p {
  margin: 5px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-container {
    padding: 30px;
  }
  
  .login-title {
    font-size: 24px;
  }
  
  .form-input {
    padding: 10px 14px;
  }
  
  .btn-login {
    padding: 12px 18px;
  }
}
</style>