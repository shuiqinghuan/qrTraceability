<template>
  <div class="admin-dashboard">
    <!-- 侧边导航 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2 class="sidebar-title">后台管理</h2>
      </div>
      
      <nav class="sidebar-nav">
        <ul class="nav-menu">
          <li class="nav-item">
            <router-link to="/admin/products" class="nav-link">
              <i class="nav-icon">📦</i>
              <span>产品管理</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/admin/media" class="nav-link">
              <i class="nav-icon">🖼️</i>
              <span>媒体管理</span>
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/admin/quality" class="nav-link">
              <i class="nav-icon">🏆</i>
              <span>品质管理</span>
            </router-link>
          </li>
        </ul>
      </nav>
      
      <div class="sidebar-footer">
        <button class="btn-logout" @click="handleLogout">
          <i class="logout-icon">🚪</i>
          <span>退出登录</span>
        </button>
      </div>
    </aside>
    
    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 顶部导航栏 -->
      <header class="top-nav">
        <div class="nav-left">
          <h1 class="page-title">{{ currentPageTitle }}</h1>
        </div>
        <div class="nav-right">
          <div class="user-info">
            <span class="user-name">{{ adminUser?.username || '管理员' }}</span>
          </div>
        </div>
      </header>
      
      <!-- 页面内容 -->
      <div class="page-content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from '../composables/useToast'

const router = useRouter()
const route = useRoute()
const { showToast } = useToast()

// 响应式数据
const adminUser = ref(null)

// 计算属性
const currentPageTitle = computed(() => {
  const path = route.path
  if (path.includes('/admin/products')) return '产品管理'
  if (path.includes('/admin/media')) return '媒体管理'
  if (path.includes('/admin/quality')) return '品质管理'
  return '后台管理'
})

// 生命周期
onMounted(() => {
  loadAdminUser()
})

// 方法
const loadAdminUser = () => {
  const userStr = localStorage.getItem('admin_user')
  if (userStr) {
    adminUser.value = JSON.parse(userStr)
  }
}

const handleLogout = () => {
  // 清除本地存储
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_user')
  
  showToast('已退出登录', 'success')
  
  // 跳转到登录页面
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-dashboard {
  display: flex;
  min-height: 100vh;
  background-color: #f9fafb;
}

/* 侧边导航 */
.sidebar {
  width: 250px;
  background-color: #1e293b;
  color: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #334155;
}

.sidebar-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #f8fafc;
}

.sidebar-nav {
  flex: 1;
  padding: 20px 0;
}

.nav-menu {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  margin-bottom: 4px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: #cbd5e1;
  text-decoration: none;
  transition: all 0.3s;
  border-left: 3px solid transparent;
}

.nav-link:hover {
  background-color: #334155;
  color: white;
  border-left-color: #3b82f6;
}

.nav-link.router-link-active {
  background-color: #3b82f6;
  color: white;
  border-left-color: #3b82f6;
}

.nav-icon {
  font-size: 18px;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #334155;
}

.btn-logout {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-logout:hover {
  background-color: #dc2626;
}

.logout-icon {
  font-size: 16px;
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.top-nav {
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #475569;
}

.page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 200px;
  }
  
  .sidebar-title {
    font-size: 18px;
  }
  
  .nav-link {
    padding: 10px 16px;
  }
  
  .page-content {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .admin-dashboard {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    flex-direction: row;
    padding: 10px;
  }
  
  .sidebar-header {
    border-bottom: none;
    padding: 0 10px;
  }
  
  .sidebar-nav {
    flex: 1;
    padding: 0;
  }
  
  .nav-menu {
    display: flex;
    gap: 10px;
  }
  
  .nav-item {
    margin-bottom: 0;
  }
  
  .nav-link {
    padding: 8px 12px;
    font-size: 12px;
  }
  
  .sidebar-footer {
    border-top: none;
    padding: 0 10px;
  }
  
  .btn-logout {
    padding: 8px 12px;
    font-size: 12px;
  }
}
</style>